package usecase

import (
	"context"
	"fmt"

	"github.com/Na322Pr/misinder/internal/dto"
	"github.com/Na322Pr/misinder/internal/repository"
)

type ProfileUsecase struct {
	repo *repository.UserRepository
}

func NewUserUsecase(repo *repository.UserRepository) *ProfileUsecase {
	return &ProfileUsecase{repo: repo}
}

func (uc *ProfileUsecase) CreateUser(ctx context.Context, userID int64, username string) error {
	op := "ProfileUsecase.CreateUser"

	userDTO := dto.UserDTO{
		ID:       userID,
		Username: username,
		Status:   dto.UserStatusProfile,
	}

	fmt.Println(userDTO)

	if err := uc.repo.CreateUser(ctx, userDTO); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ProfileUsecase) GetUserStatus(ctx context.Context, userID int64) (dto.UserStatus, error) {
	op := "ProfileUsecase.CreatGetUserStatuseUser"

	status, err := uc.repo.GetUserStatus(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (uc *ProfileUsecase) AddInterest(ctx context.Context, userID int64, interest string) error {
	op := "ProfileUsecase.AddInterest"

	userDTO, err := uc.repo.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

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

func (uc *ProfileUsecase) SetStatusFree(ctx context.Context, userID int64) error {
	op := "ProfileUsecase.SetStatusFree"

	if err := uc.repo.UpdateStatus(ctx, userID, dto.UserStatusFree); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
