package controller

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Na322Pr/misinder/internal/controller/handler"
	"github.com/Na322Pr/misinder/internal/dto"
	"github.com/Na322Pr/misinder/internal/usecase"
)

type Controller struct {
	bot            *tgbotapi.BotAPI
	uc             *usecase.ProfileUsecase
	commandHandler *handler.CommandHandler
	profileHandler *handler.ProfileHandler
}

func NewController(bot *tgbotapi.BotAPI, uc *usecase.ProfileUsecase) *Controller {
	controller := &Controller{
		bot:            bot,
		uc:             uc,
		commandHandler: handler.NewCommandHandler(bot, uc),
		profileHandler: handler.NewProfileHandler(bot, uc),
	}

	return controller
}

func (c *Controller) HandleUpdates(ctx context.Context) {
	op := "Controller.HandleUpdates"

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)
	for update := range updates {

		if update.Message.IsCommand() {
			c.commandHandler.Handle(ctx, update)
			continue
		}

		status, err := c.uc.GetUserStatus(ctx, update.Message.From.ID)
		if err != nil {
			fmt.Printf("%s: %v", op, err)
		}

		if status == dto.UserStatusProfile {
			c.profileHandler.Handle(ctx, update)
			continue
		}
	}
}
