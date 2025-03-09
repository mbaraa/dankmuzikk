package auth

import (
	"context"
	"dankmuzikk/actions"
	"dankmuzikk/app"
	"net/http"
)

// Cookie keys
const (
	VerificationTokenKey = "verification-token"
	SessionTokenKey      = "token"
)

// Context keys
const (
	ProfileIdKey       = "profile-id"
	PlaylistPermission = "playlist-permission"
)

type mw struct {
	usecases *actions.Actions
}

// New returns a new auth middleware instance.
// Using a GORMDBGetter because this is supposed to be a light fetch,
// Where BaseDB doesn't provide column selection yet :(
func New(usecases *actions.Actions) *mw {
	return &mw{
		usecases: usecases,
	}
}

// AuthApi authenticates an API's handler.
func (a *mw) AuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, err := a.authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ProfileIdKey, profile.Id)
		h(w, r.WithContext(ctx))
	}
}

// OptionalAuthApi authenticates a page's handler optionally (without 401).
func (a *mw) OptionalAuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, err := a.authenticate(r)
		if err != nil {
			h(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), ProfileIdKey, profile.Id)
		h(w, r.WithContext(ctx))
	}
}

func (a *mw) authenticate(r *http.Request) (actions.AuthenticateUserPayload, error) {
	sessionToken, ok := r.Header["Authorization"]
	if !ok {
		return actions.AuthenticateUserPayload{}, &app.ErrInvalidVerificationToken{}
	}

	return a.usecases.AuthenticateUser(sessionToken[0])
}
