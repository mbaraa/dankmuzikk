package jwt

import (
	"net/http"
)

type ErrEmptyToken struct{}

func (e ErrEmptyToken) Error() string {
	return "empty-token"
}

func (e ErrEmptyToken) ClientStatusCode() int {
	return http.StatusBadRequest
}

func (e ErrEmptyToken) ExtraData() map[string]any {
	return nil
}

func (e ErrEmptyToken) ExposeToClients() bool {
	return true
}

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

type ErrInvalidSubject struct{}

func (e ErrInvalidSubject) Error() string {
	return "invalid-subject"
}

func (e ErrInvalidSubject) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrInvalidSubject) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidSubject) ExposeToClients() bool {
	return true
}
