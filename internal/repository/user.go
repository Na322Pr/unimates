package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Na322Pr/misinder/internal/dto"
	"github.com/Na322Pr/misinder/pkg/postgres"
)

type UserRepository struct {
	*postgres.Postgres
}

func NewUserRepository(pg *postgres.Postgres) *UserRepository {
	return &UserRepository{pg}
}

func (r *UserRepository) CreateUser(ctx context.Context, userDTO dto.UserDTO) error {
	op := "UserRepository.CreateUser"
	query := `INSERT INTO users(id, username, status, interests) VALUES($1, $2, $3, $4)`

	_, err := r.Conn.Exec(ctx, query,
		userDTO.ID,
		userDTO.Username,
		userDTO.Status,
		userDTO.Interests,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, userID int64) (*dto.UserDTO, error) {
	op := "UserRepository.GetUser"
	query := `SELECT id, username, status, interests FROM users WHERE id = $1`

	var userDTO dto.UserDTO
	err := r.Conn.QueryRow(ctx, query, userID).Scan(
		&userDTO.ID,
		&userDTO.Username,
		&userDTO.Status,
		&userDTO.Interests,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &userDTO, nil
}

func (r *UserRepository) GetUserStatus(ctx context.Context, userID int64) (dto.UserStatus, error) {
	op := "UserRepository.GetUserStatus"
	query := `SELECT status FROM users WHERE id = $1`

	var status dto.UserStatus
	err := r.Conn.QueryRow(ctx, query, userID).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userDTO dto.UserDTO) error {
	op := "UserRepository.UpdateUserInterests"
	query := `UPDATE users SET status = $2, interests = $3 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, userDTO.ID, userDTO.Status, userDTO.Interests)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, userID int64, status dto.UserStatus) error {
	op := "UserRepository.UpdateUserInterests"
	query := `UPDATE users SET status = $2 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// func (r *UserRepository) GetMatch(ctx context.Context, userID int64) ([]*dto.UserDTO, error) {
// 	op := "UserRepository.GetMatch"

// 	query := `select id, username, interests FROM `

// 	return
// }
