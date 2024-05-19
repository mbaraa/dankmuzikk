package login

import "errors"

var (
	ErrAccountNotFound         = errors.New("account was not found")
	ErrProfileNotFound         = errors.New("profile was not found")
	ErrAccountExists           = errors.New("an account with the associated email already exists")
	ErrExpiredVerificationCode = errors.New("expired verification code")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
)
