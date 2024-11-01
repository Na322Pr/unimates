package repository

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrOfferNotFound = errors.New("offer not found")
)
