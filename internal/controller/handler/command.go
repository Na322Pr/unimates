package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/keyboard/reply"
	"github.com/Na322Pr/unimates/internal/usecase"

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
		h.Start(ctx, update)
	case "rules":
		h.Rules(ctx, update)
	case "profile":
		h.Profile(ctx, update)
	case "myprofile":
		h.MyProfile(ctx, update)
	case "offer":
		h.Offer(ctx, update)
	}
}

func (h *CommandHandler) Start(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.Start"

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

func (h *CommandHandler) Rules(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.Rules"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Тут пока пусто")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) Profile(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.Profile"

	err := h.uc.RecreateUser(ctx, update.Message.From.ID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Заполните профиль")
	msg.ReplyMarkup = reply.StartFillProfileKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) MyProfile(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.MyProfile"

	interests, err := h.uc.GetUserInterests(
		ctx,
		update.Message.From.ID,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msgText := fmt.Sprintf("Список ваших интересов:\n%s", strings.Join(interests, ", "))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) Offer(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.Offer"

	err := h.uc.SetStatus(
		ctx,
		update.Message.From.ID,
		dto.UserStatusOffer,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	err = h.uc.CreateOffer(
		ctx,
		update.Message.From.ID,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Текст вашего предложения")

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}
