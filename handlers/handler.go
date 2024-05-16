package handlers

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
	"net/http"
	"slices"
	"strings"
)

var noAuthPaths = []string{"/login", "/signup"}

// Handler is handler for pages and APIs, where it wraps the common stuff in one place.
type Handler struct {
	profileRepo db.GORMDBGetter
	jwtUtil     jwt.Decoder[any]
}

// NewHandler returns a new AuthHandler instance.
// Using a GORMDBGetter because this is supposed to be a light fetch,
// Where BaseDB doesn't provide column selection yet :(
func NewHandler(
	accountRepo db.GORMDBGetter,
	jwtUtil jwt.Decoder[any],
) *Handler {
	return &Handler{accountRepo, jwtUtil}
}

// OptionalAuthPage authenticates a page's handler optionally (without redirection).
func (a *Handler) OptionalAuthPage(h http.HandlerFunc) http.HandlerFunc {
	return a.NoAuthPage(func(w http.ResponseWriter, r *http.Request) {
		profileId, err := a.authenticate(r)
		if err != nil {
			h(w, r)
			return
		}
		h(w, r.WithContext(context.WithValue(r.Context(), ProfileIdKey, profileId)))
	})
}

// AuthPage authenticates a page's handler.
func (a *Handler) AuthPage(h http.HandlerFunc) http.HandlerFunc {
	return a.NoAuthPage(func(w http.ResponseWriter, r *http.Request) {
		htmxRedirect := IsNoReloadPage(r)
		profileId, err := a.authenticate(r)
		authed := err == nil

		switch {
		case authed && slices.Contains(noAuthPaths, r.URL.Path):
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
		case !authed && slices.Contains(noAuthPaths, r.URL.Path):
			h(w, r.WithContext(context.WithValue(r.Context(), ProfileIdKey, profileId)))
		case !authed && htmxRedirect:
			w.Header().Set("HX-Redirect", "/login")
		case !authed && !htmxRedirect:
			http.Redirect(w, r, config.Env().Hostname+"/login", http.StatusTemporaryRedirect)
		default:
			h(w, r.WithContext(context.WithValue(r.Context(), ProfileIdKey, profileId)))
		}
	})
}

// NoAuthPage returns a page's handler after setting Content-Type to text/html, and some context values.
func (a *Handler) NoAuthPage(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		ctx := context.WithValue(r.Context(), ThemeKey, getTheme(r))
		ctx = context.WithValue(ctx, IsMobileKey, isMobile(r))
		h(w, r.WithContext(ctx))
	}
}

// AuthApi authenticates an API's handler.
func (a *Handler) AuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileId, err := a.authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(w, r.WithContext(context.WithValue(r.Context(), ProfileIdKey, profileId)))
	}
}

// NoAuthApi returns a page's handler after setting Content-Type to application/json.
func (a *Handler) NoAuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}

func (a *Handler) authenticate(r *http.Request) (uint, error) {
	sessionToken, err := r.Cookie(SessionTokenKey)
	if err != nil {
		return 0, err
	}
	theThing, err := a.jwtUtil.Decode(sessionToken.Value, jwt.SessionToken)
	if err != nil {
		return 0, err
	}
	payload, valid := theThing.Payload.(map[string]any)
	if !valid || payload == nil {
		return 0, err
	}
	username, validUsername := theThing.Payload.(map[string]any)["username"].(string)
	if !validUsername || username == "" {
		return 0, err
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
		return 0, err
	}

	return profile.Id, nil
}

func isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}

func getTheme(r *http.Request) string {
	themeCookie, err := r.Cookie(ThemeName)
	if err != nil || themeCookie == nil || themeCookie.Value == "" {
		return "default"
	}
	switch themeCookie.Value {
	case "black":
		return "black"
	case "default":
		fallthrough
	default:
		return "default"
	}
}

// IsNoReloadPage checks if the requested page requires a no reload or not.
func IsNoReloadPage(r *http.Request) bool {
	noReload, exists := r.URL.Query()["no_reload"]
	return exists && noReload[0] == "true"
}
