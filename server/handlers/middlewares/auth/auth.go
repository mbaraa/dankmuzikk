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
	GuestKey   = "guest"
)

type Middleware struct {
	usecases *actions.Actions
}

// New returns a new auth middleware instance.
func New(usecases *actions.Actions) *Middleware {
	return &Middleware{
		usecases: usecases,
	}
}

func (a *Middleware) AuthHandler(h http.Handler) http.Handler {
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
func (a *Middleware) AuthApi(h http.HandlerFunc) http.HandlerFunc {
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
func (a *Middleware) OptionalAuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := a.authenticate(r)
		ctx := r.Context()
		if err != nil {
			// TODO: set client hash from web :)
			clientHash, ok := r.Header["X-Client-Hash"]
			if !ok {
				h(w, r)
				return
			}
			ctx = context.WithValue(ctx, GuestKey, clientHash[0])
		}
		ctx = context.WithValue(ctx, AccountKey, account)
		h(w, r.WithContext(ctx))
	}
}

func (a *Middleware) authenticate(r *http.Request) (models.Account, error) {
	sessionToken, ok := r.Header["Authorization"]
	if !ok {
		return models.Account{}, &app.ErrInvalidSessionToken{}
	}

	return a.usecases.AuthenticateAccount(sessionToken[0])
}
