package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	pages.Profile(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
