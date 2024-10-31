package usecase

// import (
// 	"context"
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/Na322Pr/unimates/internal/dto"
// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func (uc *ProfileUsecase) GetOfferStatus(ctx context.Context, userID int64) (dto.OfferStatus, error) {
// 	op := "ProfileUsecase.GetOfferStatus"

// 	offerDTO, err := uc.repo.GetOffer(ctx, userID)
// 	if err != nil {
// 		return "", fmt.Errorf("%s: %w", op, err)
// 	}

// 	if offerDTO.Text.String == "" {
// 		return dto.OfferStatusText, nil
// 	}

// 	if offerDTO.Interest.String == "" {
// 		return dto.OfferStatusInterest, nil
// 	}

// 	return dto.OfferStatusReady, nil
// }

// func (uc *ProfileUsecase) CreateOffer(ctx context.Context, userID int64) error {
// 	op := "ProfileUsecase.GetOfferStatus"

// 	if err := uc.repo.CreateOffer(ctx, userID); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (uc *ProfileUsecase) AddOfferText(ctx context.Context, userID int64, text string) error {
// 	op := "ProfileUsecase.GetOfferStatus"

// 	if err := uc.repo.UpdateOfferText(ctx, userID, text); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (uc *ProfileUsecase) AddOfferInterest(ctx context.Context, userID int64, interest string) error {
// 	op := "ProfileUsecase.GetOfferStatus"

// 	if err := uc.repo.UpdateOfferInterest(ctx, userID, strings.ToLower(interest)); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (uc *ProfileUsecase) DeleteOffer(ctx context.Context, userID int64) error {
// 	op := "ProfileUsecase.DeleteOffer"

// 	if err := uc.repo.DeletOffer(ctx, userID); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (uc *ProfileUsecase) GetMatch(ctx context.Context, userID int64) error {
// 	op := "ProfileUsecase.GetMatch"

// 	offerDTO, err := uc.repo.GetOffer(ctx, userID)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	userDTO, err := uc.repo.GetUser(ctx, userID)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	matchUserIDs, err := uc.repo.GetMatch(
// 		ctx,
// 		offerDTO.Interest.String,
// 		userDTO.Interests,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	for _, matchUserID := range matchUserIDs {
// 		if matchUserID == userID {
// 			continue
// 		}
// 		if err := uc.sendOffer(ctx, userID, matchUserID, offerDTO.Text.String); err != nil {
// 			fmt.Printf("%s: %v", op, err)
// 		}
// 	}

// 	return nil
// }

// func (uc *ProfileUsecase) sendOffer(ctx context.Context, senderID, userID int64, offerText string) error {
// 	op := "ProfileUsecase.sendOffer"
// 	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Для вас новое предложение\n%s\nИнтересно?", offerText))

// 	OfferReplyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("Да", strconv.FormatInt(senderID, 10)),
// 			tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
// 		),
// 	)

// 	msg.ReplyMarkup = OfferReplyKeyboard
// 	if _, err := uc.bot.Send(msg); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }
