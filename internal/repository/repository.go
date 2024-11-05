package repository

import (
	"context"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, userDTO dto.UserDTO) error
	GetUser(ctx context.Context, userID int64) (*dto.UserDTO, error)
	GetUserUsername(ctx context.Context, userID int64) (string, error)
	GetUserStatus(ctx context.Context, userID int64) (dto.UserStatus, error)
	GetAdminUserList(ctx context.Context, userID int64) ([]dto.UserDTO, error)
	UpdateStatus(ctx context.Context, userID int64, status dto.UserStatus) error
	UpdateRole(ctx context.Context, userID int64, role dto.UserRole) error
}

var _ User = (*UserRepository)(nil)

type Interest interface {
	PreloadInterests(ctx context.Context, interests []string) error
	GetInterests(ctx context.Context) ([]dto.InterestDTO, error)
	GetUserInterests(ctx context.Context, userID int64) ([]string, error)
	CreateUserInterest(ctx context.Context, userID int64, interestID int) error
	DeleteUserInterest(ctx context.Context, userID int64, interestID int) error
}

var _ Interest = (*InterestRepository)(nil)

type Offer interface {
	CreateOffer(ctx context.Context, userID int64) (int64, error)
	GetOfferByID(ctx context.Context, offer int64) (*dto.OfferDTO, error)
	GetOfferByText(ctx context.Context, userID int64, offerText string) (*dto.OfferDTO, error)
	GetUserOffers(ctx context.Context, userID int64) ([]dto.OfferDTO, error)
	CreateUserAcceptedOffer(ctx context.Context, userID, offerID int64) error
	GetUserAcceptedOffer(ctx context.Context, offerID int64) ([]string, error)
	UpdateOfferText(ctx context.Context, offerID int64, text string) error
	UpdateOfferInterest(ctx context.Context, offerID int64, interestID int) error
	DeletOffer(ctx context.Context, offerID int64) error
	GetMatch(ctx context.Context, userID int64, interestID int) ([]int64, error)
}

type Repository struct {
	User
	Interest
	Offer
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		User:     NewUserRepository(pg),
		Interest: NewInterestRepository(pg),
		Offer:    NewOfferRepository(pg),
	}
}
