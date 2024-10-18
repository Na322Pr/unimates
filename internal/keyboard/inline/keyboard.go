package inline

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var OfferReplyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "Что-то"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "Что-то"),
	),
)
