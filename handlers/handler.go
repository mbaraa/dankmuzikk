package handlers

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
	"net/http"
	"slices"
)

var noAuthPaths = []string{"/login", "/signup"}

// Handler is handler for pages and APIs, where it wraps the common stuff in one place.
type Handler struct {
	accountRepo db.GORMDBGetter
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

// AuthPage authenticates a page's handler.
func (a *Handler) AuthPage(h http.HandlerFunc) http.HandlerFunc {
	return a.NoAuthPage(func(w http.ResponseWriter, r *http.Request) {
		htmxRedirect := a.isNoReloadPage(r)
		authed := a.isAuthed(r)

		switch {
		case authed && slices.Contains(noAuthPaths, r.URL.Path):
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
		case !authed && slices.Contains(noAuthPaths, r.URL.Path):
			h(w, r)
		case !authed && htmxRedirect:
			w.Header().Set("HX-Redirect", "/login")
		case !authed && !htmxRedirect:
			http.Redirect(w, r, config.Env().Hostname+"/login", http.StatusTemporaryRedirect)
		default:
			h(w, r)
		}

	})
}

// NoAuthPage returns a page's handler after setting Content-Type to text/html.
func (a *Handler) NoAuthPage(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h(w, r)
	}
}

func (a *Handler) isAuthed(r *http.Request) bool {
	sessionToken, err := r.Cookie(SessionTokenKey)
	if err != nil {
		return false
	}
	theThing, err := a.jwtUtil.Decode(sessionToken.Value, jwt.SessionToken)
	if err != nil {
		return false
	}
	payload, valid := theThing.Payload.(map[string]any)
	if !valid || payload == nil {
		return false
	}
	userEmail, validEmail := theThing.Payload.(map[string]any)["email"].(string)
	if !validEmail || userEmail == "" {
		return false
	}

	var act models.Account

	return a.
		accountRepo.
		GetDB().
		Model(&act).
		Select("id").
		Where("email = ?", userEmail).
		First(&act).
		Error == nil
}

func (a *Handler) isNoReloadPage(r *http.Request) bool {
	noReload, exists := r.URL.Query()["no_reload"]
	return exists && noReload[0] == "true"
}
