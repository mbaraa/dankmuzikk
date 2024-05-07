package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandlePlaylistsPage(w http.ResponseWriter, r *http.Request) {
	if isNoReload(r) {
		pages.PlaylistsNoReload().Render(context.Background(), w)
		return
	}
	pages.Playlists(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
