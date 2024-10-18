package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Na322Pr/unimates/internal/config"
	"github.com/Na322Pr/unimates/internal/controller"
	"github.com/Na322Pr/unimates/internal/repository"
	"github.com/Na322Pr/unimates/internal/usecase"
	"github.com/Na322Pr/unimates/pkg/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	bot, err := tgbotapi.NewBotAPI(cfg.TG.Token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запуск бота"},
		{Command: "rules", Description: "Правила использования"},
		{Command: "profile", Description: "Заполнение профиля"},
		{Command: "myprofile", Description: "Посмотреть свой профиль"},
		{Command: "offer", Description: "Создать предложение"},
	}

	cmdCfg := tgbotapi.NewSetMyCommands(commands...)
	_, err = bot.Request(cmdCfg)
	if err != nil {
		log.Fatalf("Failed to set commands: %v", err)
	}

	pg, err := postgres.Connection(psqlDSN(cfg))
	if err != nil {
		panic(err)
	}
	defer pg.Close()

	repository := repository.NewUserRepository(pg)
	usecase := usecase.NewUserUsecase(bot, repository)
	cntr := controller.NewController(bot, usecase)

	cntr.HandleUpdates(ctx)
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
