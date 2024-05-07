package pages

import (
	"dankmuzikk/config"
	"dankmuzikk/handlers"
	"dankmuzikk/services/jwt"
	"net/http"
	"slices"
	"strings"

	_ "github.com/a-h/templ"
)

var noAuthPaths = []string{"/login", "/signup"}

func Handler(hand http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		hand(w, r)
	}
}

func AuthHandler(hand http.HandlerFunc, jwtUtil jwt.Manager[any]) http.HandlerFunc {
	return Handler(func(w http.ResponseWriter, r *http.Request) {
		htmxRedirect := isNoReload(r)
		authed := isAuthed(r, jwtUtil)

		switch {
		case authed && slices.Contains(noAuthPaths, r.URL.Path):
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
		case !authed && slices.Contains(noAuthPaths, r.URL.Path):
			hand(w, r)
		case !authed && htmxRedirect:
			w.Header().Set("HX-Redirect", "/login")
		case !authed && !htmxRedirect:
			http.Redirect(w, r, config.Env().Hostname+"/login", http.StatusTemporaryRedirect)
		default:
			hand(w, r)
		}

	})
}

func isAuthed(r *http.Request, jwtUtil jwt.Manager[any]) bool {
	sessionToken, err := r.Cookie(handlers.SessionTokenKey)
	if err != nil {
		return false
	}
	err = jwtUtil.Validate(sessionToken.Value, jwt.SessionToken)
	if err != nil {
		return false
	}

	return true
}

func isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}

func getTheme(r *http.Request) string {
	themeCookie, err := r.Cookie(handlers.ThemeName)
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

func isNoReload(r *http.Request) bool {
	noReload, exists := r.URL.Query()["no_reload"]
	return exists && noReload[0] == "true"
}
