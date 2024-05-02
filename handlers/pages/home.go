package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
	"strings"

	_ "github.com/a-h/templ"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	pages.Index(isMobile(r)).Render(context.Background(), w)
}

func Handler(hand http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		hand(w, r)
	}
}

func isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}
