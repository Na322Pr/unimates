package usecase

import (
	"context"
	"fmt"
	"strings"
	"unicode"

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

	for i := 0; i < len(interests); i++ {
		runes := []rune(interests[i])
		runes[0] = unicode.ToUpper(runes[0])
		interests[i] = string(runes)
	}

	return interests, nil
}

func (uc *InterestUsecase) CreateUserInterest(ctx context.Context, userID int64, newInterest string) error {
	op := "InterestUsecase.CreateUserInterest"

	if len(newInterest) == 0 {
		return fmt.Errorf("%s: %w", op, ErrEmptyInterest)
	}

	newInterest = strings.ToLower(newInterest)

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

	interestID, err := uc.repo.CreateCustomInterest(ctx, newInterest)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := uc.repo.CreateUserInterest(ctx, userID, interestID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *InterestUsecase) DeleteUserInterest(ctx context.Context, userID int64, delInterest string) error {
	op := "InterestUsecase.DeleteUserInterest"

	if len(delInterest) == 0 {
		return fmt.Errorf("%s: %w", op, ErrEmptyInterest)
	}

	delInterest = strings.ToLower(delInterest)

	interests, err := uc.repo.GetUserInterestsDTOs(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, interest := range interests {
		if interest.Name == delInterest {
			if err := uc.repo.DeleteUserInterest(ctx, userID, interest.ID); err != nil {
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
		if levenshtein(delInterest, interest.Name) <= 2 {
			counter++
			near = append(near, interest.Name)
		}
	}

	for i := 0; i < len(near); i++ {
		runes := []rune(near[i])
		runes[0] = unicode.ToUpper(runes[0])
		near[i] = string(runes)
	}

	msgText := "У вас нет похожих интересов...\nПопробуйте ввести что-то еще"
	if len(near) != 0 {
		msgText = "Похожие варианты:\n" + strings.Join(near, "\n")
	}

	msg := tgbotapi.NewMessage(userID, msgText)

	if _, err := uc.bot.Send(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return ErrInvalidInterestName
}
