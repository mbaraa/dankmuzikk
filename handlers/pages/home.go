package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	if isNoReload(r) {
		pages.IndexNoReload().Render(context.Background(), w)
		return
	}
	pages.Index(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
