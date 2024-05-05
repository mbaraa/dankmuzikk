package pages

import (
	"dankmuzikk/config"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
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
		sessionToken, err := r.Cookie(handlers.SessionTokenKey)
		if err != nil {
			log.Errorln("[AUTH]:", err)
			if slices.Contains(noAuthPaths, r.URL.Path) {
				hand(w, r)
				return
			}
			http.Redirect(w, r, config.Env().Hostname+"/login", http.StatusTemporaryRedirect)
			return
		}

		err = jwtUtil.Validate(sessionToken.Value, jwt.SessionToken)
		if err != nil {
			log.Errorln("[AUTH]:", err)
			if slices.Contains(noAuthPaths, r.URL.Path) {
				hand(w, r)
				return
			}
			http.Redirect(w, r, config.Env().Hostname+"/login", http.StatusTemporaryRedirect)
			return
		}

		if slices.Contains(noAuthPaths, r.URL.Path) {
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
			return
		}

		hand(w, r)
	})
}

func isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}
