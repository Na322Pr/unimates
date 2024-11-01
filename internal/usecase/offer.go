package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OfferUsecase struct {
	bot          *tgbotapi.BotAPI
	repo         repository.Offer
	activeOffers map[int64]int64
}

func NewOfferUsecase(bot *tgbotapi.BotAPI, repo repository.Offer) *OfferUsecase {
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

func (uc *OfferUsecase) DeleteOffer(ctx context.Context, userID int64, orderID int64) {

}

func (uc *OfferUsecase) GetOfferAcceptances(ctx context.Context, userID int64, offerText string) error {
	op := "OfferUsecase.GetOfferAcceptances"

	offerDTO, err := uc.repo.GetOfferByText(ctx, userID, offerText)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	users, err := uc.repo.GetUserAcceptedOffer(ctx, offerDTO.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	msgText := "Пока что никто не откликнулся на это предложение"
	if len(users) != 0 {
		msgText = "Вот кто откликнулся:\n@" + strings.Join(users, "\n@")
	}

	msg := tgbotapi.NewMessage(userID, msgText)

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
	op := "OfferUsecase.GetOfferStatus"

	if err := uc.repo.UpdateOfferText(ctx, uc.activeOffers[userID], text); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *OfferUsecase) AddOfferInterest(ctx context.Context, userID int64, interest string) error {
	op := "OfferUsecase.GetOfferStatus"

	if err := uc.repo.UpdateOfferInterest(ctx, uc.activeOffers[userID], 1); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
