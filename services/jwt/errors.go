package jwt

import (
	"errors"
)

var (
	ErrEmptyToken     = errors.New("jwt: empty token")
	ErrInvalidToken   = errors.New("jwt: invalid token")
	ErrExpiredToken   = errors.New("jwt: expired token")
	ErrInvalidSubject = errors.New("jwt: invalid token's subject")
)
