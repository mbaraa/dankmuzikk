package auth

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/handlers/middlewares/contenttype"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
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
	profileRepo db.GORMDBGetter
	jwtUtil     jwt.Decoder[jwt.Json]
}

// New returns a new auth middle ware instance.
// Using a GORMDBGetter because this is supposed to be a light fetch,
// Where BaseDB doesn't provide column selection yet :(
func New(
	accountRepo db.GORMDBGetter,
	jwtUtil jwt.Decoder[jwt.Json],
) *mw {
	return &mw{accountRepo, jwtUtil}
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
	theThing, err := a.jwtUtil.Decode(sessionToken.Value, jwt.SessionToken)
	if err != nil {
		return entities.Profile{}, err
	}
	username, validUsername := theThing.Payload["username"].(string)
	if !validUsername || username == "" {
		return entities.Profile{}, err
	}

	var profile models.Profile

	err = a.
		profileRepo.
		GetDB().
		Model(&profile).
		Select("id").
		Where("username = ?", username).
		First(&profile).
		Error

	if err != nil {
		return entities.Profile{}, err
	}

	return entities.Profile{
		Id:   profile.Id,
		Name: profile.Name,
	}, nil
}
