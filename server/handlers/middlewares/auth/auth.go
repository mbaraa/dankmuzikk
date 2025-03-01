package auth

import (
	"context"
	"dankmuzikk/app/entities"
	"dankmuzikk/app/models"
	"dankmuzikk/db"
	"dankmuzikk/services/jwt"
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
	profileRepo db.GORMDBGetter
	jwtUtil     jwt.Decoder[jwt.Json]
}

// New returns a new auth middleware instance.
// Using a GORMDBGetter because this is supposed to be a light fetch,
// Where BaseDB doesn't provide column selection yet :(
func New(
	accountRepo db.GORMDBGetter,
	jwtUtil jwt.Decoder[jwt.Json],
) *mw {
	return &mw{accountRepo, jwtUtil}
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
