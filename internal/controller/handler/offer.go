package handler

import (
	"context"
	"fmt"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OfferHandler struct {
	uc  *usecase.ProfileUsecase
	bot *tgbotapi.BotAPI
}

func NewOfferHandler(bot *tgbotapi.BotAPI, uc *usecase.ProfileUsecase) *OfferHandler {
	return &OfferHandler{
		bot: bot,
		uc:  uc,
	}
}

func (h *OfferHandler) Handle(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.Handle"

	status, err := h.uc.GetOfferStatus(ctx, update.Message.From.ID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	switch status {
	case dto.OfferStatusText:
		h.AddText(ctx, update)
	case dto.OfferStatusInterest:
		h.AddInterest(ctx, update)
	case dto.OfferStatusReady:
		h.SendOffer(ctx, update)
	}
}

func (h *OfferHandler) AddText(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.AddText"

	err := h.uc.AddOfferText(
		ctx,
		update.Message.From.ID,
		update.Message.Text,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Хэштег")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *OfferHandler) AddInterest(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.AddInterest"

	err := h.uc.AddOfferInterest(
		ctx,
		update.Message.From.ID,
		update.Message.Text,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Отправить предложение?")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *OfferHandler) SendOffer(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.SendOffer"

	err := h.uc.SetStatus(
		ctx,
		update.Message.From.ID,
		dto.UserStatusFree,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	if err = h.uc.GetMatch(ctx, update.Message.From.ID); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Предложение отправлено пользователям")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	if err := h.uc.DeleteOffer(ctx, update.Message.From.ID); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}
