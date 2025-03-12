package jwt

import (
	"dankmuzikk/actions"
	"dankmuzikk/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Jwt implements JWTManager to verify session tokens
type Jwt[T any] struct{}

// NewJWTImpl returns a new JWTImpl instance,
// and since session tokens are to validate users the working type is models.User
func New[T any]() *Jwt[T] {
	return &Jwt[T]{}
}

// Sign returns a JWT string(which will be the session token) based on the set JWT secret,
// using HS256 algorithm, and the given validity
// and an occurring error
func (s *Jwt[T]) Sign(data T, subject actions.Subject, expTime time.Duration) (string, error) {
	expirationDate := jwt.NumericDate{Time: time.Now().UTC().Add(expTime)}
	currentTime := jwt.NumericDate{Time: time.Now().UTC()}

	claims := actions.JwtClaims[T]{
		Payload: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &expirationDate,
			Issuer:    "DankMuzikk",
			Subject:   subject,
			NotBefore: &currentTime,
			IssuedAt:  &currentTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	// error is ignored here, because it's either json marshaling or signing method(using library's built-in method) errors
	// which are correct as of this version of golang-jwt, god knows what could go wrong in the future :)
	tokenString, _ := token.SignedString([]byte(config.Env().JwtSecret))

	return tokenString, nil
}

// Validate checks the validity of the JWT string, and returns an occurring error
func (s *Jwt[T]) Validate(token string, subject actions.Subject) error {
	_, err := s.Decode(token, subject)
	if err != nil {
		return err
	}

	return nil
}

// Decode decodes the given token using the set JWT secret
func (s *Jwt[T]) Decode(token string, subject actions.Subject) (actions.JwtClaims[T], error) {
	if len(token) == 0 {
		return actions.JwtClaims[T]{}, &ErrInvalidToken{}
	}

	claims := actions.JwtClaims[T]{}

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if claims.Subject != subject {
			return nil, &ErrInvalidToken{}
		}
		if claims.ExpiresAt.Time.Before(time.Now().UTC()) {
			return nil, &ErrExpiredToken{}
		}

		return []byte(config.Env().JwtSecret), nil
	})

	if err != nil {
		return actions.JwtClaims[T]{}, err
	}

	return claims, nil
}
