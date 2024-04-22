package handlers

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
	"strings"
)

func HandleHomePage(hand *http.ServeMux) {
	hand.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		isMobile := strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
		pages.Index(isMobile).Render(context.Background(), w)
	})
}
