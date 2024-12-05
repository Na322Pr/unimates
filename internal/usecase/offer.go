package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/keyboard/reply"
	"github.com/Na322Pr/unimates/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OfferUsecase struct {
	bot          *tgbotapi.BotAPI
	repo         repository.Repository
	activeOffers map[int64]int64
}

func NewOfferUsecase(bot *tgbotapi.BotAPI, repo repository.Repository) *OfferUsecase {
	return &OfferUsecase{bot: bot, repo: repo, activeOffers: make(map[int64]int64)}
}

func (uc *OfferUsecase) GetUserOffers(ctx context.Context, userID int64) ([]dto.OfferDTO, error) {
	op := "OfferUsecase.GetUserOffers"

	offers, err := uc.repo.GetUserOffers(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return offers, nil
}

func (uc *OfferUsecase) CreateOffer(ctx context.Context, userID int64) error {
	op := "OfferUsecase.CreateOffer"

	offerID, err := uc.repo.CreateOffer(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	uc.activeOffers[userID] = offerID
	return nil
}

func (uc *OfferUsecase) SelectActiveOffer(ctx context.Context, userID, offerID int64) error {
	// op := "OfferUsecase.SelectActiveOffer"
	uc.activeOffers[userID] = offerID
	return nil
}

func (uc *OfferUsecase) DeleteOffer(ctx context.Context, userID int64) error {
	op := "OfferUsecase.DeleteOffer"
	if err := uc.repo.DeletOffer(ctx, uc.activeOffers[userID]); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *OfferUsecase) CreateOfferAcceptance(ctx context.Context, userID, offerID int64) error {
	op := "OfferUsecase.CreateOfferAcceptance"

	if err := uc.repo.CreateUserAcceptedOffer(ctx, userID, offerID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *OfferUsecase) GetOfferAcceptances(ctx context.Context, userID int64, offerText string) error {
	op := "OfferUsecase.GetOfferAcceptances"

	offerDTO, err := uc.repo.GetOfferByText(ctx, userID, offerText)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	uc.SelectActiveOffer(ctx, userID, offerDTO.ID)

	users, err := uc.repo.GetUserAcceptedOffer(ctx, offerDTO.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	msgText := "Пока что никто не откликнулся на это предложение"
	if len(users) != 0 {
		msgText = "Вот кто откликнулся:\n@" + strings.Join(users, "\n@")
	}

	msg := tgbotapi.NewMessage(userID, msgText)
	msg.ReplyMarkup = reply.EditOfferKeyboard

	if _, err := uc.bot.Send(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *OfferUsecase) GetOfferStatus(ctx context.Context, userID int64) (dto.OfferStatus, error) {
	op := "OfferUsecase.GetOfferStatus"

	offerID := uc.activeOffers[userID]

	offerDTO, err := uc.repo.GetOfferByID(ctx, offerID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !offerDTO.Text.Valid {
		return dto.OfferStatusText, nil
	}

	if !offerDTO.InterestID.Valid {
		return dto.OfferStatusInterest, nil
	}

	return dto.OfferStatusReady, nil
}

func (uc *OfferUsecase) AddOfferText(ctx context.Context, userID int64, text string) error {
	op := "OfferUsecase.AddOfferText"

	if err := uc.repo.UpdateOfferText(ctx, uc.activeOffers[userID], text); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *OfferUsecase) AddOfferInterest(ctx context.Context, userID int64, newInterest string) error {
	op := "OfferUsecase.AddOfferInterest"

	interests, err := uc.repo.GetInterests(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newInterest = strings.ToLower(newInterest)

	newInterestID := -1

	for _, interest := range interests {
		if interest.Name == newInterest {
			newInterestID = interest.ID
			break
		}
	}

	if newInterestID == -1 {
		msgText := "У нас нет похожих интересов...\nПопробуйте ввести что-то еще"
		msg := tgbotapi.NewMessage(userID, msgText)

		if _, err := uc.bot.Send(msg); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return ErrInvalidInterestName
	}

	if err := uc.repo.UpdateOfferInterest(ctx, uc.activeOffers[userID], newInterestID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(userID, newInterestID)
	fmt.Println(userID, newInterestID)

	matchUserIDs, err := uc.repo.GetMatch(ctx, userID, newInterestID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	fmt.Printf("%s: %d", op, matchUserIDs)

	offerDTO, err := uc.repo.Offer.GetOfferByID(ctx, uc.activeOffers[userID])
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, matchUserID := range matchUserIDs {
		if matchUserID == userID {
			continue
		}
		if err := uc.sendOffer(ctx, userID, matchUserID, offerDTO.ID, offerDTO.Text.String); err != nil {
			fmt.Printf("%s: %v", op, err)
		}
	}

	return nil
}

func (uc *OfferUsecase) sendOffer(ctx context.Context, senderID, userID, offerID int64, offerText string) error {
	op := "OfferUsecase.sendOffer"
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Для вас новое предложение\n%s\nИнтересно?", offerText))

	data := make(map[string]int64)
	data["sender_id"] = senderID
	data["offer_id"] = offerID

	jsonBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	OfferReplyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", string(jsonBytes)),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
		),
	)

	msg.ReplyMarkup = OfferReplyKeyboard
	if _, err := uc.bot.Send(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
