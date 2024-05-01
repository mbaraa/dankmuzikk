package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandleTOSPage(hand *http.ServeMux) {
	hand.HandleFunc("/tos", func(w http.ResponseWriter, r *http.Request) {
		pages.TOS(isMobile(r)).Render(context.Background(), w)
	})
}
