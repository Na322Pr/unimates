package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/pkg/postgres"
)

type OfferRepository struct {
	*postgres.Postgres
}

func NewOfferRepository(pg *postgres.Postgres) *OfferRepository {
	return &OfferRepository{pg}
}

func (r *OfferRepository) CreateOffer(ctx context.Context, userID int64) (int64, error) {
	op := "OfferRepository.CreateOffer"
	query := `INSERT INTO offers(user_id) VALUES($1) RETURNING id`

	var offerID int64
	err := r.Conn.QueryRow(ctx, query, userID).Scan(&offerID)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return offerID, nil
}

func (r *OfferRepository) GetOfferByID(ctx context.Context, offerID int64) (*dto.OfferDTO, error) {
	op := "OfferRepository.GetUserOffers"
	query := `SELECT id, user_id, text, interest_id, inactive_at FROM offers WHERE id = $1`

	var offerDTO dto.OfferDTO
	err := r.Conn.QueryRow(ctx, query, offerID).Scan(
		&offerDTO.ID,
		&offerDTO.UserID,
		&offerDTO.Text,
		&offerDTO.InterestID,
		&offerDTO.InactiveAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &offerDTO, nil
}

func (r *OfferRepository) GetOfferByText(ctx context.Context, userID int64, offerText string) (*dto.OfferDTO, error) {
	op := "OfferRepository.GetOfferByText"
	query := `SELECT id, user_id, text, interest_id, inactive_at FROM offers WHERE user_id = $1 AND text = $2 limit 1`

	var offerDTO dto.OfferDTO
	err := r.Conn.QueryRow(ctx, query, userID, offerText).Scan(
		&offerDTO.ID,
		&offerDTO.UserID,
		&offerDTO.Text,
		&offerDTO.InterestID,
		&offerDTO.InactiveAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrOfferNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &offerDTO, nil
}

func (r *OfferRepository) GetUserOffers(ctx context.Context, userID int64) ([]dto.OfferDTO, error) {
	op := "OfferRepository.GetUserOffers"
	query := `SELECT id, user_id, text, interest_id, inactive_at FROM offers WHERE user_id = $1`

	rows, err := r.Conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var offersDTOs []dto.OfferDTO
	for rows.Next() {
		var offerDTO dto.OfferDTO
		if err := rows.Scan(
			&offerDTO.ID,
			&offerDTO.UserID,
			&offerDTO.Text,
			&offerDTO.InterestID,
			&offerDTO.InactiveAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		offersDTOs = append(offersDTOs, offerDTO)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return offersDTOs, nil
}

func (r *OfferRepository) CreateUserAcceptedOffer(ctx context.Context, userID, offerID int64) error {
	op := "OfferRepository.CreateUserAcceptedOffer"
	query := `INSERT INTO offer_acceptances(user_id, offer_id) VALUES($1, $2)`

	if _, err := r.Conn.Exec(ctx, query, userID, offerID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OfferRepository) GetUserAcceptedOffer(ctx context.Context, offerID int64) ([]string, error) {
	op := "OfferRepository.GetUserAcceptedOffer"
	query := `SELECT username FROM offers
		JOIN offer_acceptances ON offers.id = offer_acceptances.offer_id
		JOIN users ON offer_acceptances.user_id = users.id
		WHERE offers.id = $1
	`

	rows, err := r.Conn.Query(ctx, query, offerID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	usernames := make([]string, 0)
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		usernames = append(usernames, username)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return usernames, nil
}

func (r *OfferRepository) UpdateOfferText(ctx context.Context, offerID int64, text string) error {
	op := "OfferRepository.UpdateOfferText"
	query := `UPDATE offers SET "text" = $2 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, offerID, text)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OfferRepository) UpdateOfferInterest(ctx context.Context, offerID int64, interestID int) error {
	op := "OfferRepository.UpdateOfferMainInterest"
	query := `UPDATE offers SET interest_id = $2 WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, offerID, interestID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OfferRepository) DeletOffer(ctx context.Context, offerID int64) error {
	op := "OfferRepository.DeleteOffer"
	query := `DELETE FROM offers WHERE id = $1`

	_, err := r.Conn.Exec(ctx, query, offerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OfferRepository) GetMatch(ctx context.Context, userID int64, interestID int) ([]int64, error) {
	op := "OfferRepository.GetMatch"
	query := `
		SELECT main_eq.user_id FROM 
		(SELECT user_id FROM user_interests WHERE interest_id = $2) AS main_eq
		JOIN
		(SELECT sui.user_id FROM user_interests AS fui
		JOIN user_interests AS sui on fui.interest_id = sui.interest_id
		WHERE fui.user_id = $1 AND fui.user_id != sui.user_id 
		GROUP BY sui.user_id
		HAVING COUNT(*) >= 3) as dop_eq
		ON main_eq.user_id = dop_eq.user_id
	`

	rows, err := r.Conn.Query(ctx, query, userID, interestID)
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
