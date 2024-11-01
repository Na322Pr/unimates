package dto

import "database/sql"

type OfferStatus string

const (
	OfferStatusText     OfferStatus = "text"
	OfferStatusInterest OfferStatus = "interest"
	OfferStatusReady    OfferStatus = "ready"
)

type OfferDTO struct {
	ID         int64          `db:"id"`
	UserID     int64          `db:"user_id"`
	Text       sql.NullString `db:"text"`
	InterestID sql.NullInt32  `db:"interest_id"`
	Notify     sql.NullBool   `db:"notify"`
	InactiveAt sql.NullTime   `db:"inactive_at"`
}
