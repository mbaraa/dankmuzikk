package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandlePlaylistsPage(hand *http.ServeMux) {
	hand.HandleFunc("/playlists", func(w http.ResponseWriter, r *http.Request) {
		pages.Playlists(isMobile(r)).Render(context.Background(), w)
	})
}
