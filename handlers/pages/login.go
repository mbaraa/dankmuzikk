package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandleLoginPage(hand *http.ServeMux) {
	hand.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		pages.Login(isMobile(r)).Render(context.Background(), w)
	})
}
