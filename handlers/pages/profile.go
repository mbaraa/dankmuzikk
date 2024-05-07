package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	if isNoReload(r) {
		pages.ProfileNoReload().Render(context.Background(), w)
		return
	}
	pages.Profile(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
