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

type Repository struct {
	User
	Interest
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		User:     NewUserRepository(pg),
		Interest: NewInterestRepository(pg),
	}
}
