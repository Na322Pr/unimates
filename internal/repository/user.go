package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/pkg/postgres"
)

type UserRepository struct {
	*postgres.Postgres
}

func NewUserRepository(pg *postgres.Postgres) *UserRepository {
	return &UserRepository{pg}
}

func (r *UserRepository) CreateUser(ctx context.Context, userDTO dto.UserDTO) error {
	op := "UserRepository.CreateUser"
	query := `INSERT INTO users(id, username, status) VALUES($1, $2, $3)`

	_, err := r.Conn.Exec(ctx, query,
		userDTO.ID,
		userDTO.Username,
		userDTO.Status,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, userID int64) (*dto.UserDTO, error) {
	op := "UserRepository.GetUser"
	query := `SELECT id, username, status, role, created_at, modified_at FROM users WHERE id = $1`

	var userDTO dto.UserDTO
	err := r.Conn.QueryRow(ctx, query, userID).Scan(
		&userDTO.ID,
		&userDTO.Username,
		&userDTO.Status,
		&userDTO.Role,
		&userDTO.CreatedAt,
		&userDTO.ModifiedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &userDTO, nil
}

func (r *UserRepository) GetUserUsername(ctx context.Context, userID int64) (string, error) {
	op := "UserRepository.GetUserStatus"
	query := `SELECT username FROM users WHERE id = $1`

	var username string
	err := r.Conn.QueryRow(ctx, query, userID).Scan(&username)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return username, nil
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

func (r *UserRepository) GetAdminUserList(ctx context.Context, userID int64) ([]dto.UserDTO, error) {
	op := "UserRepository.GetAdminUserList"
	query := `SELECT id, username, status, role, created_at, modified_at FROM users WHERE id = $1 and role = $2`

	rows, err := r.Conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	usersDTOs := make([]dto.UserDTO, 0)

	for rows.Next() {
		var userDTO dto.UserDTO
		if err := rows.Scan(
			&userDTO.ID,
			&userDTO.Username,
			&userDTO.Status,
			&userDTO.Role,
			&userDTO.CreatedAt,
			&userDTO.ModifiedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		usersDTOs = append(usersDTOs, userDTO)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return usersDTOs, nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, userID int64, status dto.UserStatus) error {
	op := "UserRepository.UpdateStatus"
	query := `UPDATE users SET status = $2, modified_at = $3 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, status, time.Now())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID int64, role dto.UserRole) error {
	op := "UserRepository.UpdateRole"
	query := `UPDATE users SET role = $2, modified_at = $3 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, role, time.Now())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

var _ User = (*UserRepository)(nil)
