package reply

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var StartFillProfileKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Заполнить профиль"),
	),
)

var EndFillProfileKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Закончить"),
	),
)
