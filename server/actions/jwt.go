package actions

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Subject represents the token's subject
type Subject = string

const (
	// JwtSessionToken used to verify that the user is logged in correctly and can access the good stuff.
	JwtSessionToken Subject = "SESSION_TOKEN"
	// JwtVerificationToken used to verify verification code (on sign-in or sign-up).
	JwtVerificationToken Subject = "OTP_TOKEN"
)

// JwtClaims is iondsa, it's just JWT claims blyat!
type JwtClaims[T any] struct {
	jwt.RegisteredClaims
	Payload T `json:"payload"`
}

// JwtSigner is a wrapper to JWT signing method using the set JWT secret,
// claims are set(mostly unique) in each implementation of the thing
type JwtSigner[T any] interface {
	Sign(data T, subject Subject, expTime time.Duration) (string, error)
}

// JwtValidator is a wrapper to JWT validation stuff, also uses the claims for that current implementation
type JwtValidator interface {
	Validate(token string, subject Subject) error
}

// JwtDecoder is a wrapper to JWT decoding stuff, based on the implementation's claims,
// this interface is usually implemented with the other two(Signer and Validator), because reasons...
type JwtDecoder[T any] interface {
	Decode(token string, subject Subject) (JwtClaims[T], error)
}

// JwtManager is a wrapper to JWT operations, so I don't do much shit each time I work with JWT
type JwtManager[T any] interface {
	JwtSigner[T]
	JwtValidator
	JwtDecoder[T]
}
