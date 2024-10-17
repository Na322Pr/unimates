package dto

type UserStatus string

const (
	UserStatusUnkown  UserStatus = "unknown"
	UserStatusFree    UserStatus = "free"
	UserStatusProfile UserStatus = "profile"
	UserStatusOffer   UserStatus = "offer"
)

type UserDTO struct {
	ID        int64      `db:"user_id"`
	Username  string     `db:"username"`
	Interests []string   `db:"interests"`
	Status    UserStatus `db:"status"`
}
