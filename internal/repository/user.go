package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *UserRepository) UpdateUser(ctx context.Context, userDTO dto.UserDTO) error {
	op := "UserRepository.UpdateUserInterests"
	query := `UPDATE users SET status = $2, interests = $3 WHERE id = $1`

	fmt.Println("Repo interests: ", userDTO.Interests)

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

func (r *UserRepository) CreateOffer(ctx context.Context, userID int64) error {
	op := "UserRepository.CreateOffer"
	query := `INSERT INTO offers(user_id) VALUES($1)`

	_, err := r.Conn.Exec(ctx, query, userID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) GetOffer(ctx context.Context, userID int64) (*dto.OfferDTO, error) {
	op := "UserRepository.GetOffer"
	query := `SELECT user_id, "text", interest FROM offers WHERE user_id = $1`

	var offerDTO dto.OfferDTO
	err := r.Conn.QueryRow(ctx, query, userID).Scan(
		&offerDTO.UserID,
		&offerDTO.Text,
		&offerDTO.Interest,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &offerDTO, nil

}

func (r *UserRepository) UpdateOfferText(ctx context.Context, userID int64, text string) error {
	op := "UserRepository.UpdateOfferText"
	query := `UPDATE offers SET "text" = $2 WHERE user_id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, text)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) UpdateOfferInterest(ctx context.Context, userID int64, interest string) error {
	op := "UserRepository.UpdateOfferInterest"
	query := `UPDATE offers SET interest = $2 WHERE user_id = $1`

	_, err := r.Conn.Exec(ctx, query, userID, interest)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) DeletOffer(ctx context.Context, userID int64) error {
	op := "UserRepository.DeleteOffer"
	query := `DELETE FROM offers WHERE user_id = $1`

	_, err := r.Conn.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) GetMatch(ctx context.Context, mainInterest string, interests []string) ([]int64, error) {
	op := "UserRepository.GetMatch"
	query := `
		SELECT id
		FROM (
			SELECT id, UNNEST(interests::varchar[]) AS interest FROM (
				SELECT * FROM users WHERE $1 = ANY(interests)
			) AS main_interest
		) AS cur_interests
		JOIN 
		(SELECT UNNEST($2::varchar[]) AS interest) AS given_interests 
		ON cur_interests.interest = given_interests.interest
		GROUP BY id HAVING COUNT(*) >= 3;
	`

	rows, err := r.Conn.Query(ctx, query, mainInterest, interests)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	userIDs := make([]int64, 0, 10)

	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userIDs, nil
}
