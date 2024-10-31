package usecase

import (
	"context"
	"fmt"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserUsecase struct {
	bot  *tgbotapi.BotAPI
	repo repository.User
}

func NewUserUsecase(bot *tgbotapi.BotAPI, repo repository.User) *UserUsecase {
	return &UserUsecase{
		bot:  bot,
		repo: repo,
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, userID int64, username string) error {
	op := "UserUsecase.CreateUser"

	userDTO := dto.UserDTO{
		ID:       userID,
		Username: username,
		Status:   dto.UserStatusInterest,
	}

	if err := uc.repo.CreateUser(ctx, userDTO); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// func (uc *UserUsecase) RecreateUser(ctx context.Context, userID int64) error {
// 	op := "UserUsecase.RecreateUser"

// 	userDTO, err := uc.repo.GetUser(ctx, userID)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	userDTO.Status = dto.UserStatusInterest

// 	if err := uc.repo.UpdateUser(ctx, *userDTO); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

func (uc *UserUsecase) GetUserUsername(ctx context.Context, userID int64) (string, error) {
	op := "UserUsecase.GetUserUsername"

	username, err := uc.repo.GetUserUsername(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return username, nil
}

func (uc *UserUsecase) GetUserStatus(ctx context.Context, userID int64) (dto.UserStatus, error) {
	op := "UserUsecase.CreatGetUserStatuseUser"

	status, err := uc.repo.GetUserStatus(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

// func (uc *UserUsecase) GetUserInterests(ctx context.Context, userID int64) ([]string, error) {
// 	op := "UserUsecase.GetUserInterests"

// 	userDTO, err := uc.repo.GetUser(ctx, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return userDTO.Interests, nil
// }

// func (uc *UserUsecase) AddInterest(ctx context.Context, userID int64, interest string) error {
// 	op := "UserUsecase.AddInterest"

// 	userDTO, err := uc.repo.GetUser(ctx, userID)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	interest = strings.ToLower(interest)

// 	for _, val := range userDTO.Interests {
// 		if val == interest {
// 			return fmt.Errorf("%s: %w", op, ErrInterestAlreadyExist)
// 		}
// 	}

// 	userDTO.Interests = append(userDTO.Interests, interest)

// 	if err := uc.repo.UpdateUser(ctx, *userDTO); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

func (uc *UserUsecase) SetStatus(ctx context.Context, userID int64, status dto.UserStatus) error {
	op := "UserUsecase.SetStatus"

	if err := uc.repo.UpdateStatus(ctx, userID, status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
