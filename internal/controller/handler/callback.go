package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Na322Pr/unimates/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler struct {
	uc  *usecase.ProfileUsecase
	bot *tgbotapi.BotAPI
}

func NewCallbackHandler(bot *tgbotapi.BotAPI, uc *usecase.ProfileUsecase) *CallbackHandler {
	return &CallbackHandler{
		bot: bot,
		uc:  uc,
	}
}

func (h *CallbackHandler) Handle(ctx context.Context, update tgbotapi.Update) {
	// op := "OfferHandler.Handle"

	data := update.CallbackQuery.Data

	switch data {
	case "no":
		return
	default:
		h.SendOfferAnswer(ctx, update)
	}
}

func (h *CallbackHandler) SendOfferAnswer(ctx context.Context, update tgbotapi.Update) {
	op := "CallbackHandler.SendOfferAnswer"

	userID, err := strconv.ParseInt(update.CallbackQuery.Data, 10, 64)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}
	offerRespID := update.CallbackQuery.From.ID

	username, err := h.uc.GetUserUsername(
		context.Background(),
		offerRespID,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(
		int64(userID),
		fmt.Sprintf("Ваше предложение интересно @%s", username),
	)

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	username, err = h.uc.GetUserUsername(
		context.Background(),
		userID,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg = tgbotapi.NewMessage(
		int64(offerRespID),
		fmt.Sprintf("Вот чел который это предложил... @%s", username),
	)

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}
