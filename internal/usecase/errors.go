package usecase

import "errors"

var (
	ErrUserAlreadyExist     = errors.New("user already exist")
	ErrInterestAlreadyExist = errors.New("interest already exist")
)
