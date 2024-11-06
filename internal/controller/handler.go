package controller

import (
	"context"
	"fmt"

	"github.com/Na322Pr/unimates/internal/controller/handler"
	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Controller struct {
	bot             *tgbotapi.BotAPI
	uc              *usecase.Usecase
	commandHandler  *handler.CommandHandler
	callbackHandler *handler.CallbackHandler
	interestHandler *handler.InterestHandler
	offerHandler    *handler.OfferHandler
}

func NewController(bot *tgbotapi.BotAPI, uc *usecase.Usecase) *Controller {
	controller := &Controller{
		bot:             bot,
		uc:              uc,
		commandHandler:  handler.NewCommandHandler(bot, uc),
		callbackHandler: handler.NewCallbackHandler(bot, uc),
		interestHandler: handler.NewInterestHandler(bot, uc),
		offerHandler:    handler.NewOfferHandler(bot, uc),
	}

	return controller
}

func (c *Controller) HandleUpdates(ctx context.Context) {
	op := "Controller.HandleUpdates"

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)
	for update := range updates {

		if update.CallbackQuery != nil {
			c.callbackHandler.Handle(ctx, update)
			continue
		}

		if update.Message.IsCommand() {
			c.commandHandler.Handle(ctx, update)
			continue
		}

		status, err := c.uc.User.GetUserStatus(ctx, update.Message.From.ID)
		if err != nil {
			fmt.Printf("%s: %v", op, err)
		}

		switch status {
		case dto.UserStatusInterest, dto.UserStatusInterestAdd, dto.UserStatusInterestDelete:
			c.interestHandler.Handle(ctx, update)
		case dto.UserStatusOffer, dto.UserStatusOfferNew, dto.UserStatusOfferEdit:
			c.offerHandler.Handle(ctx, update)
		}
	}
}
