package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandleProfilePage(hand *http.ServeMux) {
	hand.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		pages.Profile(isMobile(r)).Render(context.Background(), w)
	})
}
