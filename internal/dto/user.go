package dto

import (
	"time"
)

type UserStatus string

const (
	UserStatusUnkown            UserStatus = "unknown"
	UserStatusFree              UserStatus = "empty"
	UserStatusInterest          UserStatus = "interest"
	UserStatusInterestAdd       UserStatus = "interest_add"
	UserStatusInterestAddCustom UserStatus = "interest_add_custom"
	UserStatusInterestDelete    UserStatus = "interest_delete"
	UserStatusOffer             UserStatus = "offer"
	UserStatusOfferNew          UserStatus = "offer_new"
	UserStatusOfferEdit         UserStatus = "offer_edit"
)

type UserRole string

const (
	UserRoleBase  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type UserDTO struct {
	ID         int64      `db:"id"`
	Username   string     `db:"username"`
	Role       string     `db:"role"`
	Status     UserStatus `db:"status"`
	CreatedAt  time.Time  `db:"created_at"`
	ModifiedAt time.Time  `db:"modified_at"`
}
