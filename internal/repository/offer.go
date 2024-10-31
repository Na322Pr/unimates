package repository

// func (r *UserRepository) CreateOffer(ctx context.Context, userID int64) error {
// 	op := "UserRepository.CreateOffer"
// 	query := `INSERT INTO offers(user_id) VALUES($1)`

// 	_, err := r.Conn.Exec(ctx, query, userID)

// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (r *UserRepository) GetOffer(ctx context.Context, userID int64) (*dto.OfferDTO, error) {
// 	op := "UserRepository.GetOffer"
// 	query := `SELECT user_id, "text", interest FROM offers WHERE user_id = $1`

// 	var offerDTO dto.OfferDTO
// 	err := r.Conn.QueryRow(ctx, query, userID).Scan(
// 		&offerDTO.UserID,
// 		&offerDTO.Text,
// 		&offerDTO.Interest,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
// 		}
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return &offerDTO, nil

// }

// func (r *UserRepository) UpdateOfferText(ctx context.Context, userID int64, text string) error {
// 	op := "UserRepository.UpdateOfferText"
// 	query := `UPDATE offers SET "text" = $2 WHERE user_id = $1`

// 	_, err := r.Conn.Exec(ctx, query, userID, text)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (r *UserRepository) UpdateOfferInterest(ctx context.Context, userID int64, interest string) error {
// 	op := "UserRepository.UpdateOfferInterest"
// 	query := `UPDATE offers SET interest = $2 WHERE user_id = $1`

// 	_, err := r.Conn.Exec(ctx, query, userID, interest)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (r *UserRepository) DeletOffer(ctx context.Context, userID int64) error {
// 	op := "UserRepository.DeleteOffer"
// 	query := `DELETE FROM offers WHERE user_id = $1`

// 	_, err := r.Conn.Exec(ctx, query, userID)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (r *UserRepository) GetMatch(ctx context.Context, mainInterest string, interests []string) ([]int64, error) {
// 	op := "UserRepository.GetMatch"
// 	query := `
// 		SELECT id
// 		FROM (
// 			SELECT id, UNNEST(interests::varchar[]) AS interest FROM (
// 				SELECT * FROM users WHERE $1 = ANY(interests)
// 			) AS main_interest
// 		) AS cur_interests
// 		JOIN
// 		(SELECT UNNEST($2::varchar[]) AS interest) AS given_interests
// 		ON cur_interests.interest = given_interests.interest
// 		GROUP BY id HAVING COUNT(*) >= 3;
// 	`

// 	rows, err := r.Conn.Query(ctx, query, mainInterest, interests)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	defer rows.Close()

// userIDs := make([]int64, 0, 10)

// 	for rows.Next() {
// 		var userID int64
// 		if err := rows.Scan(&userID); err != nil {
// 			return nil, fmt.Errorf("%s: %w", op, err)
// 		}

// userIDs = append(userIDs, userID)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return userIDs, nil
// }
