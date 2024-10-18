package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProfileUsecase struct {
	bot  *tgbotapi.BotAPI
	repo *repository.UserRepository
}

func NewUserUsecase(bot *tgbotapi.BotAPI, repo *repository.UserRepository) *ProfileUsecase {
	return &ProfileUsecase{
		bot:  bot,
		repo: repo,
	}
}

func (uc *ProfileUsecase) CreateUser(ctx context.Context, userID int64, username string) error {
	op := "ProfileUsecase.CreateUser"

	userDTO := dto.UserDTO{
		ID:       userID,
		Username: username,
		Status:   dto.UserStatusProfile,
	}

	if err := uc.repo.CreateUser(ctx, userDTO); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) RecreateUser(ctx context.Context, userID int64) error {
	op := "ProfileUsecase.RecreateUser"

	userDTO, err := uc.repo.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userDTO.Status = dto.UserStatusProfile
	userDTO.Interests = []string{}

	if err := uc.repo.UpdateUser(ctx, *userDTO); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) GetUserUsername(ctx context.Context, userID int64) (string, error) {
	op := "ProfileUsecase.GetUserUsername"

	username, err := uc.repo.GetUserUsername(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return username, nil
}

func (uc *ProfileUsecase) GetUserStatus(ctx context.Context, userID int64) (dto.UserStatus, error) {
	op := "ProfileUsecase.CreatGetUserStatuseUser"

	status, err := uc.repo.GetUserStatus(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (uc *ProfileUsecase) GetUserInterests(ctx context.Context, userID int64) ([]string, error) {
	op := "ProfileUsecase.GetUserInterests"

	userDTO, err := uc.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userDTO.Interests, nil
}

func (uc *ProfileUsecase) AddInterest(ctx context.Context, userID int64, interest string) error {
	op := "ProfileUsecase.AddInterest"

	userDTO, err := uc.repo.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	interest = strings.ToLower(interest)

	for _, val := range userDTO.Interests {
		if val == interest {
			return fmt.Errorf("%s: %w", op, ErrInterestAlreadyExist)
		}
	}

	userDTO.Interests = append(userDTO.Interests, interest)

	if err := uc.repo.UpdateUser(ctx, *userDTO); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) SetStatus(ctx context.Context, userID int64, status dto.UserStatus) error {
	op := "ProfileUsecase.SetStatus"

	if err := uc.repo.UpdateStatus(ctx, userID, status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) GetOfferStatus(ctx context.Context, userID int64) (dto.OfferStatus, error) {
	op := "ProfileUsecase.GetOfferStatus"

	offerDTO, err := uc.repo.GetOffer(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if offerDTO.Text.String == "" {
		return dto.OfferStatusText, nil
	}

	if offerDTO.Interest.String == "" {
		return dto.OfferStatusInterest, nil
	}

	return dto.OfferStatusReady, nil
}

func (uc *ProfileUsecase) CreateOffer(ctx context.Context, userID int64) error {
	op := "ProfileUsecase.GetOfferStatus"

	if err := uc.repo.CreateOffer(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) AddOfferText(ctx context.Context, userID int64, text string) error {
	op := "ProfileUsecase.GetOfferStatus"

	if err := uc.repo.UpdateOfferText(ctx, userID, text); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) AddOfferInterest(ctx context.Context, userID int64, interest string) error {
	op := "ProfileUsecase.GetOfferStatus"

	if err := uc.repo.UpdateOfferInterest(ctx, userID, strings.ToLower(interest)); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) DeleteOffer(ctx context.Context, userID int64) error {
	op := "ProfileUsecase.DeleteOffer"

	if err := uc.repo.DeletOffer(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) GetMatch(ctx context.Context, userID int64) error {
	op := "ProfileUsecase.GetMatch"

	offerDTO, err := uc.repo.GetOffer(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userDTO, err := uc.repo.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	matchUserIDs, err := uc.repo.GetMatch(
		ctx,
		offerDTO.Interest.String,
		userDTO.Interests,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, matchUserID := range matchUserIDs {
		if matchUserID == userID {
			continue
		}
		if err := uc.sendOffer(ctx, userID, matchUserID, offerDTO.Text.String); err != nil {
			fmt.Printf("%s: %v", op, err)
		}
	}

	return nil
}

func (uc *ProfileUsecase) sendOffer(ctx context.Context, senderID, userID int64, offerText string) error {
	op := "ProfileUsecase.sendOffer"
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Для вас новое предложение\n%s\nИнтересно?", offerText))

	OfferReplyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", strconv.FormatInt(senderID, 10)),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
		),
	)

	msg.ReplyMarkup = OfferReplyKeyboard
	if _, err := uc.bot.Send(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
