package pages

import (
	"dankmuzikk/log"
	"net/http"
	"strings"

	_ "github.com/a-h/templ"
)

func Handler(hand http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		hand(w, r)
	}
}

func AuthHandler(hand http.HandlerFunc) http.HandlerFunc {
	return Handler(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Cookie("token"))
		hand(w, r)
	})
}

func isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}
