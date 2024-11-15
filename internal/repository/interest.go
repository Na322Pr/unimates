package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Na322Pr/unimates/internal/dto"
	"github.com/Na322Pr/unimates/pkg/postgres"
)

type InterestRepository struct {
	*postgres.Postgres
}

func NewInterestRepository(pg *postgres.Postgres) *InterestRepository {
	return &InterestRepository{pg}
}

func (r *InterestRepository) PreloadInterests(ctx context.Context, interests []string) error {
	op := "InterestRepository.PreloadInterests"
	query := `INSERT INTO interests(name) values`

	var values []string
	var args []any

	for id, name := range interests {
		values = append(values, fmt.Sprintf("($%d)", id+1))
		args = append(args, name)
	}

	query += strings.Join(values, ", ")

	_, err := r.Conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = "UPDATE interests SET name = LOWER(name)"
	_, err = r.Conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *InterestRepository) CreateCustomInterest(ctx context.Context, interest string) error {
	op := "InterestRepository.CreateCustomInterest"
	query := "INSERT INTO interests(name) VALUES($1);"

	_, err := r.Conn.Exec(ctx, query, interest)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *InterestRepository) GetInterests(ctx context.Context) ([]dto.InterestDTO, error) {
	op := "InterestRepository.GetInterests"
	query := "SELECT id, name FROM interests"

	rows, err := r.Conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	interestsDTOs := make([]dto.InterestDTO, 0, 100)

	for rows.Next() {
		var interestDTO dto.InterestDTO
		if err := rows.Scan(&interestDTO.ID, &interestDTO.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		interestsDTOs = append(interestsDTOs, interestDTO)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return interestsDTOs, nil
}

func (r *InterestRepository) GetUserInterests(ctx context.Context, userID int64) ([]string, error) {
	op := "InterestRepository.GetUserInterests"
	query := `SELECT name FROM user_interests 
		JOIN interests ON user_interests.interest_id = interests.id
		WHERE user_id = $1`

	rows, err := r.Conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var interests []string
	for rows.Next() {
		var interest string
		if err := rows.Scan(&interest); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		interests = append(interests, interest)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return interests, nil
}

func (r *InterestRepository) GetUserInterestsDTOs(ctx context.Context, userID int64) ([]dto.InterestDTO, error) {
	op := "InterestRepository.GetUserInterestsDTOs"
	query := `SELECT id, name FROM user_interests 
		JOIN interests ON user_interests.interest_id = interests.id
		WHERE user_id = $1`

	rows, err := r.Conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	interestsDTOs := make([]dto.InterestDTO, 0, 100)
	for rows.Next() {
		var interestDTO dto.InterestDTO
		if err := rows.Scan(&interestDTO.ID, &interestDTO.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		interestsDTOs = append(interestsDTOs, interestDTO)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return interestsDTOs, nil
}

func (r *InterestRepository) CreateUserInterest(ctx context.Context, userID int64, interestID int) error {
	op := "InterestRepository.CreateUserInterest"
	query := "INSERT INTO user_interests(user_id, interest_id) VALUES($1, $2)"

	_, err := r.Conn.Exec(ctx, query, userID, interestID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *InterestRepository) DeleteUserInterest(ctx context.Context, userID int64, interestID int) error {
	op := "InterestRepository.DeleteUserInterest"
	query := "DELETE FROM user_interests WHERE user_id = $1 AND interest_id = $2"

	_, err := r.Conn.Exec(ctx, query, userID, interestID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
