package dto

import "database/sql"

type UserStatus string

const (
	UserStatusUnkown  UserStatus = "unknown"
	UserStatusFree    UserStatus = "free"
	UserStatusProfile UserStatus = "profile"
	UserStatusOffer   UserStatus = "offer"
)

type OfferStatus string

const (
	OfferStatusText     OfferStatus = "text"
	OfferStatusInterest OfferStatus = "interest"
	OfferStatusReady    OfferStatus = "ready"
)

type UserDTO struct {
	ID        int64      `db:"id"`
	Username  string     `db:"username"`
	Interests []string   `db:"interests"`
	Status    UserStatus `db:"status"`
}

type OfferDTO struct {
	UserID   int64          `db:"user_id"`
	Text     sql.NullString `db:"interest"`
	Interest sql.NullString `db:"interest"`
}
