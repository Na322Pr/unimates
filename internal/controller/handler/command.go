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
	uc  *usecase.Usecase
}

func NewCommandHandler(bot *tgbotapi.BotAPI, uc *usecase.Usecase) *CommandHandler {
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
		h.EditInterests(ctx, update)
	case "myprofile":
		h.MyInterests(ctx, update)
	case "myoffers":
		h.MyOffers(ctx, update)
	case "offer":
		// h.Offer(ctx, update)
	}
}

func (h *CommandHandler) Start(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.Start"

	err := h.uc.User.CreateUser(
		ctx,
		update.Message.From.ID,
		update.Message.From.UserName,
	)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Это бот для поиска людей с похожими интересами, заполните профиль")
	// msg.ReplyMarkup = reply.StartFillProfileKeyboard

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

func (h *CommandHandler) EditInterests(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.EditInterests"

	userID := update.Message.From.ID

	if err := h.uc.User.SetStatus(ctx, userID, dto.UserStatusInterest); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	interests, err := h.uc.Interest.GetUserInterests(ctx, userID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	var msgText string
	if len(interests) == 0 {
		msgText = "У вас пока нет добавленных интересов"
	} else {
		msgText = fmt.Sprintf("Список ваших интересов:\n%s", strings.Join(interests, "\n"))
	}

	msg := tgbotapi.NewMessage(userID, msgText)
	msg.ReplyMarkup = reply.EditInterestKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) MyOffers(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.MyOffers"
	userID := update.Message.From.ID

	if err := h.uc.User.SetStatus(ctx, userID, dto.UserStatusOffer); err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	offers, err := h.uc.Offer.GetUserOffers(ctx, userID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	var keyboardRows [][]tgbotapi.KeyboardButton

	row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Новое предложение"))
	keyboardRows = append(keyboardRows, row)

	for _, offer := range offers {
		row = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(offer.Text.String))
		keyboardRows = append(keyboardRows, row)
	}

	row = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Сохранить"))
	keyboardRows = append(keyboardRows, row)

	keyboard := tgbotapi.NewReplyKeyboard(keyboardRows...)

	msgText := "Меню предложений\n"

	msg := tgbotapi.NewMessage(userID, msgText)
	msg.ReplyMarkup = keyboard

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) MyInterests(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.MyInterests"

	interests, err := h.uc.Interest.GetUserInterests(ctx, update.Message.From.ID)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
	}

	var msgText string
	if len(interests) == 0 {
		msgText = "У вас пока нет добавленных интересов"
	} else {
		msgText = fmt.Sprintf("Список ваших интересов:\n%s", strings.Join(interests, "\n"))
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("%s: %v", op, err)
	}
}

// func (h *CommandHandler) Offer(ctx context.Context, update tgbotapi.Update) {
// 	op := "CommandHandler.Offer"

// 	err := h.uc.SetStatus(
// 		ctx,
// 		update.Message.From.ID,
// 		dto.UserStatusOffer,
// 	)
// 	if err != nil {
// 		fmt.Printf("%s: %v", op, err)
// 	}

// 	err = h.uc.CreateOffer(
// 		ctx,
// 		update.Message.From.ID,
// 	)
// 	if err != nil {
// 		fmt.Printf("%s: %v", op, err)
// 	}

// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
// 		"Текст вашего предложения")

// 	if _, err := h.bot.Send(msg); err != nil {
// 		fmt.Printf("%s: %v", op, err)
// 	}
// }
