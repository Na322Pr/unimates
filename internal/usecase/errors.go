package usecase

import "errors"

var (
	ErrUserAlreadyExist     = errors.New("user already exist")
	ErrInterestAlreadyExist = errors.New("interest already exist")
	ErrEmptyInterest        = errors.New("interest is empty")
	ErrInvalidInterestName  = errors.New("invalid interest name")
)
