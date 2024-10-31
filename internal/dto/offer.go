package dto

import "database/sql"

type OfferStatus string

const (
	OfferStatusText     OfferStatus = "text"
	OfferStatusInterest OfferStatus = "interest"
	OfferStatusReady    OfferStatus = "ready"
)

type OfferDTO struct {
	UserID   int64          `db:"user_id"`
	Text     sql.NullString `db:"interest"`
	Interest sql.NullString `db:"interest"`
}
