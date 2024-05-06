package jwt

import (
	"dankmuzikk/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Json map[string]any

// JWTImpl implements JWTManager to verify session tokens
type JWTImpl[T any] struct{}

// NewJWTImpl returns a new JWTImpl instance,
// and since session tokens are to validate users the working type is models.User
func NewJWTImpl() Manager[any] {
	return &JWTImpl[any]{}
}

// Sign returns a JWT string(which will be the session token) based on the set JWT secret,
// using HS256 algorithm, and the given validity
// and an occurring error
func (s *JWTImpl[T]) Sign(data T, subject Subject, expTime time.Duration) (string, error) {
	expirationDate := jwt.NumericDate{Time: time.Now().UTC().Add(expTime)}
	currentTime := jwt.NumericDate{Time: time.Now().UTC()}

	claims := Claims[T]{
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
func (s *JWTImpl[T]) Validate(token string, subject Subject) error {
	_, err := s.Decode(token, subject)
	if err != nil {
		return err
	}

	return nil
}

// Decode decodes the given token using the set JWT secret
func (s *JWTImpl[T]) Decode(token string, subject Subject) (Claims[T], error) {
	if len(token) == 0 {
		return Claims[T]{}, ErrEmptyToken
	}

	claims := Claims[T]{}

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if claims.Subject != subject {
			return nil, ErrInvalidSubject
		}
		if claims.ExpiresAt.Time.Before(time.Now().UTC()) {
			return nil, ErrExpiredToken
		}

		return []byte(config.Env().JwtSecret), nil
	})

	if err != nil {
		return Claims[T]{}, err
	}

	return claims, nil
}
