package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Na322Pr/unimates/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler struct {
	uc  *usecase.Usecase
	bot *tgbotapi.BotAPI
}

func NewCallbackHandler(bot *tgbotapi.BotAPI, uc *usecase.Usecase) *CallbackHandler {
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
		h.SendOfferAnswerNo(ctx, update)
	default:
		h.SendOfferAnswerYes(ctx, update)
	}
}

func (h *CallbackHandler) SendOfferAnswerNo(ctx context.Context, update tgbotapi.Update) {
	op := "CallbackHandler.SendOfferAnswerNo"

	msg := tgbotapi.NewMessage(
		int64(update.CallbackQuery.From.ID),
		"Спасибо за обратную связь",
	)

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *CallbackHandler) SendOfferAnswerYes(ctx context.Context, update tgbotapi.Update) {
	op := "CallbackHandler.SendOfferAnswerYes"

	fmt.Println(update.CallbackQuery.Data)

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &data); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	userID := int64(data["sender_id"].(float64))
	offerRespID := update.CallbackQuery.From.ID

	username, err := h.uc.User.GetUserUsername(
		ctx,
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

	username, err = h.uc.User.GetUserUsername(
		context.Background(),
		userID,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg = tgbotapi.NewMessage(
		int64(offerRespID),
		fmt.Sprintf("Вот чел, который это предложил: @%s", username),
	)

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}
