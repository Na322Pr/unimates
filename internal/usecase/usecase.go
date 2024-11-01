package usecase

import (
	"github.com/Na322Pr/unimates/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Usecase struct {
	User     UserUsecase
	Interest InterestUsecase
	Offer    OfferUsecase
}

func NewUsecase(bot *tgbotapi.BotAPI, repo *repository.Repository) *Usecase {
	return &Usecase{
		User:     *NewUserUsecase(bot, repo.User),
		Interest: *NewInterestUsecase(bot, repo.Interest),
		Offer:    *NewOfferUsecase(bot, repo.Offer),
	}
}
