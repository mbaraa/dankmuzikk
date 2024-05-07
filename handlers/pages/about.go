package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	if isNoReload(r) {
		pages.AboutNoReload().Render(context.Background(), w)
		return
	}
	pages.About(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
