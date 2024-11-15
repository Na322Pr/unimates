package handler

import (
	"context"
	"log"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/keyboard/reply"
	"github.com/Na322Pr/unimates/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InterestHandler struct {
	uc  *usecase.Usecase
	bot *tgbotapi.BotAPI
}

func NewInterestHandler(bot *tgbotapi.BotAPI, uc *usecase.Usecase) *InterestHandler {
	return &InterestHandler{
		bot: bot,
		uc:  uc,
	}
}

func (h *InterestHandler) Handle(ctx context.Context, update tgbotapi.Update) {
	// op := "InterestHandler.Handle"

	switch update.Message.Text {
	case "Заполнить интересы":
		h.StartInterest(ctx, update)
	case "Добавить существующий интерес":
		h.AddInterest(ctx, update)
	case "Создать новый интерес":
		h.AddCustomInterest(ctx, update)
	case "Удалить":
		h.DeleteInterest(ctx, update)
	case "Закончить":
		h.EndInterest(ctx, update)
	case "Сохранить":
		h.SaveInterest(ctx, update)
	default:
		h.Interest(ctx, update)
	}
}

func (h *InterestHandler) StartInterest(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.StartInterest"

	err := h.uc.User.SetStatus(
		ctx,
		update.Message.From.ID,
		dto.UserStatusInterest,
	)
	if err != nil {
		log.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Введите ваши увлечения")
	msg.ReplyMarkup = reply.EditInterestKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *InterestHandler) SaveInterest(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.EndInterest"

	err := h.uc.User.SetStatus(
		ctx,
		update.Message.From.ID,
		dto.UserStatusFree,
	)

	if err != nil {
		log.Printf("%s: %v", op, err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Готово")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *InterestHandler) Interest(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.Interest"
	userID := update.Message.From.ID

	status, err := h.uc.User.GetUserStatus(ctx, userID)
	if err != nil {
		log.Printf("%s: %v", op, err)
	}

	log.Printf("Status: %s\n", status)

	switch status {
	case dto.UserStatusInterestAdd:
		h.addInterestHandler(ctx, update)
	case dto.UserStatusInterestAddCustom:
		h.addCustomInterestHandler(ctx, update)
	case dto.UserStatusInterestDelete:
		h.deleteInterestHandler(ctx, update)
	}
}

func (h *InterestHandler) addInterestHandler(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.addInterestHandler"

	if err := h.uc.Interest.CreateUserInterest(
		ctx,
		update.Message.From.ID,
		update.Message.Text,
	); err != nil {
		if err == usecase.ErrInvalidInterestName {
			return
		} else {
			log.Printf("%s: %v", op, err)
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Интерес добавлен\nДобавить что-то еще?")

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *InterestHandler) addCustomInterestHandler(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.addCustomInterestHandler"

	msgText := "Готово"

	if update.Message.Text != "Отменить" {
		msgText = "Ваш интерес добавлен. Теперь другие пользователи смогут его использовать"
		if err := h.uc.Interest.CreateCustomInterest(ctx, update.Message.Text); err != nil {
			log.Printf("%s: %v", op, err)
			return
		}
	}

	if err := h.uc.User.SetStatus(ctx, update.Message.From.ID, dto.UserStatusInterest); err != nil {
		log.Printf("%s: %v", op, err)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyMarkup = reply.EditInterestKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *InterestHandler) deleteInterestHandler(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.deleteInterestHandler"

	if err := h.uc.Interest.DeleteUserInterest(
		ctx,
		update.Message.From.ID,
		update.Message.Text,
	); err != nil {
		if err == usecase.ErrInvalidInterestName {
			return
		} else {
			log.Printf("%s: %v", op, err)
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Интерес удален\nУдалить что-то еще?")

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *InterestHandler) AddInterest(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.AddInterest"

	if err := h.uc.User.SetStatus(ctx, update.Message.Chat.ID, dto.UserStatusInterestAdd); err != nil {
		log.Printf("%s: %v", op, err)
	}

	msgText := "Добавьте что-нибудь"

	msg := tgbotapi.NewMessage(update.Message.From.ID, msgText)
	msg.ReplyMarkup = reply.EndFillInterestKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *InterestHandler) AddCustomInterest(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.AddCustomInterest"

	if err := h.uc.User.SetStatus(ctx, update.Message.From.ID, dto.UserStatusInterestAddCustom); err != nil {
		log.Printf("%s: %v", op, err)
	}

	msgText := "Какой интерес хотите добавить?"

	msg := tgbotapi.NewMessage(update.Message.From.ID, msgText)
	msg.ReplyMarkup = reply.CancelAddCustomInterestKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *InterestHandler) DeleteInterest(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.DeleteInterest"

	if err := h.uc.User.SetStatus(ctx, update.Message.Chat.ID, dto.UserStatusInterestDelete); err != nil {
		log.Printf("%s: %v", op, err)
	}

	msgText := "Удалите что-нибудь"

	msg := tgbotapi.NewMessage(update.Message.From.ID, msgText)
	msg.ReplyMarkup = reply.EndFillInterestKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *InterestHandler) EndInterest(ctx context.Context, update tgbotapi.Update) {
	op := "InterestHandler.EndInterest"

	if err := h.uc.User.SetStatus(ctx, update.Message.Chat.ID, dto.UserStatusInterest); err != nil {
		log.Printf("%s: %v", op, err)
	}

	msgText := "Изменить что-то еще?"

	msg := tgbotapi.NewMessage(update.Message.From.ID, msgText)
	msg.ReplyMarkup = reply.EditInterestKeyboard

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}
