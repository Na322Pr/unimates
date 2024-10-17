package handler

import (
	"context"
	"fmt"

	"github.com/Na322Pr/misinder/internal/keyboard/reply"
	"github.com/Na322Pr/misinder/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	bot *tgbotapi.BotAPI
	uc  *usecase.ProfileUsecase
}

func NewCommandHandler(bot *tgbotapi.BotAPI, uc *usecase.ProfileUsecase) *CommandHandler {
	return &CommandHandler{
		bot: bot,
		uc:  uc,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		h.Start(update)
	case "rules":
		h.Rules(update)
	case "anketa":
		h.Profile(update)
	case "myanketa":
		h.MyProfile(update)
	case "offer":
		h.Offer(update)
	}
}

func (h *CommandHandler) Start(update tgbotapi.Update) {
	op := "CommandHandler.StartHandler"

	err := h.uc.CreateUser(
		context.Background(),
		update.Message.From.ID,
		update.Message.From.UserName,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Это бот для поиска людей с похожими интересами, заполните профиль")
	msg.ReplyMarkup = reply.StartFillProfileKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) Rules(update tgbotapi.Update) {

}

func (h *CommandHandler) Profile(update tgbotapi.Update) {

}

func (h *CommandHandler) MyProfile(update tgbotapi.Update) {

}

func (h *CommandHandler) Offer(update tgbotapi.Update) {

}
