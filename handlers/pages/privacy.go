package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandlePrivacyPage(hand *http.ServeMux) {
	hand.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		pages.Privacy(isMobile(r)).Render(context.Background(), w)
	})
}
