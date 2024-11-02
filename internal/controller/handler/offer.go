package handler

import (
	"context"
	"fmt"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OfferHandler struct {
	uc  *usecase.Usecase
	bot *tgbotapi.BotAPI
}

func NewOfferHandler(bot *tgbotapi.BotAPI, uc *usecase.Usecase) *OfferHandler {
	return &OfferHandler{
		bot: bot,
		uc:  uc,
	}
}

func (h *OfferHandler) Handle(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.Handle"

	status, err := h.uc.User.GetUserStatus(ctx, update.Message.From.ID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	switch status {
	case dto.UserStatusOfferNew:
		h.CreateOfferHandler(ctx, update)
	case dto.UserStatusOffer:
		h.OfferHandler(ctx, update)
	}
}

func (h *OfferHandler) CreateOfferHandler(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.CreateOfferHandler"

	status, err := h.uc.Offer.GetOfferStatus(ctx, update.Message.From.ID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
		return
	}

	switch status {
	case dto.OfferStatusText:
		h.AddText(ctx, update)
	case dto.OfferStatusInterest:
		h.AddInterest(ctx, update)
		// case dto.OfferStatusReady:
		// h.SendOffer(ctx, update)
	}
}

func (h *OfferHandler) OfferHandler(ctx context.Context, update tgbotapi.Update) {

	switch update.Message.Text {
	case "Новое предложение":
		h.AddOffer(ctx, update)
	case "Сохранить":
		h.CloseMenu(ctx, update)
	default:
		h.OfferSettings(ctx, update)
	}
}

func (h *OfferHandler) OfferSettings(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.OfferSettings"

	if err := h.uc.Offer.GetOfferAcceptances(
		ctx,
		update.Message.From.ID,
		update.Message.Text,
	); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *OfferHandler) AddOffer(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.AddOffer"

	if err := h.uc.User.SetStatus(
		ctx,
		update.Message.From.ID,
		dto.UserStatusOfferNew,
	); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	if err := h.uc.Offer.CreateOffer(ctx, update.Message.From.ID); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Введите текст предложения")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *OfferHandler) AddText(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.AddText"

	err := h.uc.Offer.AddOfferText(
		ctx,
		update.Message.From.ID,
		update.Message.Text,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Введите какой интерес объединяет людей")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *OfferHandler) AddInterest(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.AddInterest"

	err := h.uc.Offer.AddOfferInterest(
		ctx,
		update.Message.From.ID,
		update.Message.Text,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
		return
	}

	if err := h.uc.User.SetStatus(
		ctx,
		update.Message.From.ID,
		dto.UserStatusOffer,
	); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Предложение отправлено, ожидайте ответа)")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *OfferHandler) CloseMenu(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.CloseMenu"

	if err := h.uc.User.SetStatus(
		ctx,
		update.Message.From.ID,
		dto.UserStatusFree,
	); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Готово")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

// func (h *OfferHandler) SendOffer(ctx context.Context, update tgbotapi.Update) {
// 	op := "OfferHandler.SendOffer"

// 	err := h.uc.User.SetStatus(
// 		ctx,
// 		update.Message.From.ID,
// 		dto.UserStatusFree,
// 	)
// 	if err != nil {
// 		fmt.Printf("%s: %v", op, err)
// 	}

// 	if err = h.uc.GetMatch(ctx, update.Message.From.ID); err != nil {
// 		fmt.Printf("%s: %v", op, err)
// 	}

// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
// 		"Предложение отправлено пользователям")

// 	if _, err := h.bot.Send(msg); err != nil {
// 		fmt.Printf("%s: %v", op, err)
// 	}

// 	if err := h.uc.DeleteOffer(ctx, update.Message.From.ID); err != nil {
// 		fmt.Printf("%s: %v", op, err)
// 	}
// }
