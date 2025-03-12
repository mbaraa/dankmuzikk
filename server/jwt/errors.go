package jwt

import (
	"net/http"
)

type ErrInvalidToken struct{}

func (e ErrInvalidToken) Error() string {
	return "invalid-token"
}

func (e ErrInvalidToken) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrInvalidToken) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidToken) ExposeToClients() bool {
	return true
}

type ErrExpiredToken struct{}

func (e ErrExpiredToken) Error() string {
	return "expired-token"
}

func (e ErrExpiredToken) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrExpiredToken) ExtraData() map[string]any {
	return nil
}

func (e ErrExpiredToken) ExposeToClients() bool {
	return true
}
