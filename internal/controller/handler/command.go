package handler

import (
	"context"
	"fmt"
	"log"
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
	// op := "CommandHandler.Handle"

	switch update.Message.Command() {
	case "start":
		h.Start(ctx, update)
	case "rules":
		h.Rules(ctx, update)
	case "howitworks":
		h.HowItWorks(ctx, update)
	case "profile":
		h.EditInterests(ctx, update)
	case "myprofile":
		h.MyInterests(ctx, update)
	case "myoffers":
		h.MyOffers(ctx, update)
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
		log.Printf("%s: %v", op, err)
	}

	msgText := `
Приветствуем вас в Um-боте!🎉

🌟Мы рады видеть вас в нашем приложении для знакомств, где основное внимание уделяется вашим интересам и увлечениям. Здесь вы сможете найти людей, которые разделяют ваши хобби и страсти, что делает общение более естественным и увлекательным!

🤲 Мы хотим создать сообщество для всех: от хоббихорсеров до философов, от диджеев до коллекционеров! 

💡 Уникальность UM - это широкий выбор увлечений. Не нашёл свое?`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}

	msgText = `
👫Кого вы можете найти с помощью Um-бота?

- Новых знакомых с похожими увлечениями и интересами.

- Людей для совместного времяпрепровождения и развлечений.

- Новых друзей для общения и обмена опытом.

- Партнеров для бизнеса или учебного проекта.`

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}

	msgText = `❤️Не упустите возможность найти друзей и единомышленников в мире интересов! Мы уверены, что Um-бот поможет вам создать незабываемые связи. Удачи и приятного общения!💬✨`

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) Rules(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.Rules"

	msgText := `
👥Чтобы в боте было комфортно всем, просим тебя соблюдать несколько простых правил:

• Будь вежлив и уважай других участников.
• Не делись личной информацией, которую не хочешь делать публичной.
• Избегай спама, оскорблений и контента, который может быть рассчитан как незаконный или оскорбительный.
• Мы рекомендуем проводить первую встречу в стенах университета ради безопасности.`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) HowItWorks(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.HowItWorks"

	msgText := `
🧑‍💻Как работает UM?
	
1. Заполни анкету, выбрав свои увлечения: Укажи свои интересы
2. Просмотри свой профиль: После заполнения анкеты ты сможешь посмотреть свой профиль и увидеть выбранные тобой интересы.
3. Получи рекомендации: Мы будем подбирать тебе других пользователей с похожими интересами. Чем больше совпадений по интересам, тем выше шанс встретить своего идеального собеседника! 
4. После того, как потенциальный собеседник найден, у пользователя есть возможность написать еще пару слов о себе и запросить от собеседника тоже пару слов о себе, либо сразу получить ссылку на аккаунт человека.
Также у нас есть телеграмм канал для обратной связи и ваших вопросов: @unimateschannel`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) EditInterests(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.EditInterests"

	userID := update.Message.From.ID

	if err := h.uc.User.SetStatus(ctx, userID, dto.UserStatusInterest); err != nil {
		log.Printf("%s: %v", op, err)
	}

	interests, err := h.uc.Interest.GetUserInterests(ctx, userID)
	if err != nil {
		log.Printf("%s: %v", op, err)
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
		log.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) MyOffers(ctx context.Context, update tgbotapi.Update) {
	op := "OfferHandler.MyOffers"
	userID := update.Message.From.ID

	if err := h.uc.User.SetStatus(ctx, userID, dto.UserStatusOffer); err != nil {
		log.Printf("%s: %v", op, err)
	}

	offers, err := h.uc.Offer.GetUserOffers(ctx, userID)
	if err != nil {
		log.Printf("%s: %v", op, err)
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
		log.Printf("%s: %v", op, err)
	}
}

func (h *CommandHandler) MyInterests(ctx context.Context, update tgbotapi.Update) {
	op := "CommandHandler.MyInterests"

	interests, err := h.uc.Interest.GetUserInterests(ctx, update.Message.From.ID)
	if err != nil {
		log.Printf("%s: %v", op, err)
	}

	var msgText string
	if len(interests) == 0 {
		msgText = "У вас пока нет добавленных интересов"
	} else {
		msgText = fmt.Sprintf("Список ваших интересов:\n%s", strings.Join(interests, "\n"))
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("%s: %v", op, err)
	}
}
