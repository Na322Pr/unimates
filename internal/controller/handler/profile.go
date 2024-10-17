package handler

import (
	"context"
	"fmt"

	"github.com/Na322Pr/misinder/internal/dto"
	"github.com/Na322Pr/misinder/internal/keyboard/reply"
	"github.com/Na322Pr/misinder/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProfileHandler struct {
	uc  *usecase.ProfileUsecase
	bot *tgbotapi.BotAPI
}

func NewProfileHandler(bot *tgbotapi.BotAPI, uc *usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{
		bot: bot,
		uc:  uc,
	}
}

func (h *ProfileHandler) Handle(ctx context.Context, update tgbotapi.Update) {
	switch update.Message.Text {
	case "Заполнить профиль":
		h.StartProfile(update)
	case "Закончить":
		h.EndProfile(update)
	default:
		h.Interest(update)
	}
}

func (h *ProfileHandler) StartProfile(update tgbotapi.Update) {
	op := "ProfileHandler.StartProfile"

	err := h.uc.SetStatus(
		context.Background(),
		update.Message.From.ID,
		dto.UserStatusProfile,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Введите ваши увлечения")
	msg.ReplyMarkup = reply.EndFillProfileKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *ProfileHandler) EndProfile(update tgbotapi.Update) {
	op := "ProfileHandler.EndProfile"

	err := h.uc.SetStatus(
		context.Background(),
		update.Message.From.ID,
		dto.UserStatusFree,
	)

	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Профиль заполнен")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *ProfileHandler) Interest(update tgbotapi.Update) {
	op := "ProfileHandler.Interest"

	err := h.uc.AddInterest(
		context.Background(),
		update.Message.From.ID,
		update.Message.Text,
	)

	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хэштег добавлен")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}
