package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/Na322Pr/unimates/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InterestUsecase struct {
	bot  *tgbotapi.BotAPI
	repo repository.Interest
}

func NewInterestUsecase(bot *tgbotapi.BotAPI, repo repository.Interest) *InterestUsecase {
	return &InterestUsecase{bot: bot, repo: repo}
}

func (uc *InterestUsecase) GetUserInterests(ctx context.Context, userID int64) ([]string, error) {
	op := "InterestUsecase.GetUserInterests"

	interests, err := uc.repo.GetUserInterests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return interests, nil
}

func (uc *InterestUsecase) CreateUserInterest(ctx context.Context, userID int64, newInterest string) error {
	op := "InterestUsecase.CreateUserInterest"

	interests, err := uc.repo.GetInterests(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, interest := range interests {
		if interest.Name == newInterest {
			if err := uc.repo.CreateUserInterest(ctx, userID, interest.ID); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			return nil
		}
	}

	counter := 0
	near := make([]string, 0)

	for _, interest := range interests {
		if counter > 3 {
			break
		}
		if levenshtein(newInterest, interest.Name) <= 2 {
			counter++
			near = append(near, interest.Name)
		}
	}

	msgText := "У нас нет похожих интересов...\nПопробуйте ввести что-то еще"
	if len(near) != 0 {
		msgText = "Похожие варианты:\n" + strings.Join(near, "\n")
	}

	msg := tgbotapi.NewMessage(userID, msgText)

	if _, err := uc.bot.Send(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return ErrInvalidInterestName
}

func (uc *InterestUsecase) DeleteUserInterest(ctx context.Context, userID int64, newInterest string) error {
	return nil
}
