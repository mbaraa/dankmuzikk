package auth

import (
	"context"
	"dankmuzikk/actions"
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"net/http"
)

// Context keys
const (
	AccountKey = "account"
)

type mw struct {
	usecases *actions.Actions
}

// New returns a new auth middleware instance.
func New(usecases *actions.Actions) *mw {
	return &mw{
		usecases: usecases,
	}
}

func (a *mw) AuthHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, err := a.authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), AccountKey, account)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthApi authenticates an API's handler.
func (a *mw) AuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := a.authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), AccountKey, account)
		h(w, r.WithContext(ctx))
	}
}

// OptionalAuthApi authenticates an API's handler optionally (without 401).
func (a *mw) OptionalAuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := a.authenticate(r)
		if err != nil {
			h(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), AccountKey, account)
		h(w, r.WithContext(ctx))
	}
}

func (a *mw) authenticate(r *http.Request) (models.Account, error) {
	sessionToken, ok := r.Header["Authorization"]
	if !ok {
		return models.Account{}, &app.ErrInvalidSessionToken{}
	}

	return a.usecases.AuthenticateAccount(sessionToken[0])
}
