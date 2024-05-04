package jwt

import (
	"dankmuzikk/errors"
)

var JwtErrNamespace = errors.DankMuzikkErrNamespace.NewSubNamespace("jwt error")

var (
	ErrEmptyToken     = JwtErrNamespace.NewType("empty token")
	ErrInvalidToken   = JwtErrNamespace.NewType("invalid token")
	ErrExpiredToken   = JwtErrNamespace.NewType("expired token")
	ErrInvalidSubject = JwtErrNamespace.NewType("invalid token's subject")
)
