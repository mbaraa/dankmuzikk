package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandlePlaylistsPage(w http.ResponseWriter, r *http.Request) {
	pages.Playlists(isMobile(r)).Render(context.Background(), w)
}
