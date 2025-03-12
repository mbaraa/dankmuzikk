package auth

import (
	"context"
	"dankmuzikk-web/config"
	"dankmuzikk-web/entities"
	"dankmuzikk-web/handlers/middlewares/contenttype"
	"dankmuzikk-web/services/requests"
	"net/http"
	"slices"
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

var noAuthPaths = []string{"/login", "/signup"}

type mw struct {
}

// New returns a new auth middle ware instance.
// Using a GORMDBGetter because this is supposed to be a light fetch,
// Where BaseDB doesn't provide column selection yet :(
func New() *mw {
	return &mw{}
}

// AuthPage authenticates a page's handler.
func (a *mw) AuthPage(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmxRedirect := contenttype.IsNoLayoutPage(r)
		profile, err := a.authenticate(r)
		authed := err == nil
		ctx := context.WithValue(r.Context(), ProfileIdKey, profile.Id)

		switch {
		case authed && slices.Contains(noAuthPaths, r.URL.Path):
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
		case !authed && slices.Contains(noAuthPaths, r.URL.Path):
			h(w, r.WithContext(ctx))
		case !authed && htmxRedirect:
			w.Header().Set("HX-Redirect", "/login")
		case !authed && !htmxRedirect:
			http.Redirect(w, r, config.Env().Hostname+"/login", http.StatusTemporaryRedirect)
		default:
			h(w, r.WithContext(ctx))
		}
	}
}

// OptionalAuthPage authenticates a page's handler optionally (without redirection).
func (a *mw) OptionalAuthPage(h http.HandlerFunc) http.HandlerFunc {
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

func (a *mw) authenticate(r *http.Request) (entities.Profile, error) {
	sessionToken, err := r.Cookie(SessionTokenKey)
	if err != nil {
		return entities.Profile{}, err
	}

	user, err := requests.GetRequestAuth[entities.Profile]("/v1/profile", sessionToken.Value)
	if err != nil {
		return entities.Profile{}, err
	}

	return user, nil
}
