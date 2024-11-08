package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Na322Pr/unimates/internal/config"
	"github.com/Na322Pr/unimates/internal/controller"
	"github.com/Na322Pr/unimates/internal/repository"
	"github.com/Na322Pr/unimates/internal/usecase"
	"github.com/Na322Pr/unimates/pkg/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("no .env file found")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.MustLoad()
	ctx := context.Background()

	bot, err := tgbotapi.NewBotAPI(cfg.TG.Token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	// bot.Debug = true

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запуск бота"},
		{Command: "rules", Description: "Правила использования"},
		{Command: "howitworks", Description: "Как работать с ботом?"},
		{Command: "profile", Description: "Заполнение профиля"},
		{Command: "myprofile", Description: "Посмотреть свой профиль"},
		{Command: "myoffers", Description: "Меню предложений"},
	}

	cmdCfg := tgbotapi.NewSetMyCommands(commands...)
	_, err = bot.Request(cmdCfg)
	if err != nil {
		log.Fatalf("Failed to set commands: %v", err)
	}

	pg, err := postgres.Connection(psqlDSN(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	defer pg.Close()

	repository := repository.NewRepository(pg)

	if err := preloadInterests(ctx, repository.Interest); err != nil {
		fmt.Printf("Failed to load interests: %v", err)
	}

	usecase := usecase.NewUsecase(bot, repository)
	cntr := controller.NewController(bot, usecase)

	go func() {
		cntr.HandleUpdates(ctx)
	}()

	if err := NotifyOnStartUp(bot, *cfg); err != nil {
		fmt.Printf("Failed to notify admins: %v", err)
	}

	<-stop
	fmt.Println("\nShutting down...")
	os.Exit(0)
}

func NotifyOnStartUp(bot *tgbotapi.BotAPI, cfg config.Config) error {
	op := "NotifyOnStartUp"

	msgText := "Бот запущен"

	for _, adminID := range cfg.TG.AdminIDs {
		msg := tgbotapi.NewMessage(adminID, msgText)

		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func psqlDSN(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.DB,
	)
}

func preloadInterests(ctx context.Context, repo repository.Interest) error {
	op := "PreloadInterests"

	seen := make(map[string]struct{})
	var interests []string

	file, err := os.Open("./config/interests.txt")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if _, ok := seen[scanner.Text()]; ok {
			continue
		}
		seen[scanner.Text()] = struct{}{}
		interests = append(interests, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := repo.PreloadInterests(ctx, interests); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
