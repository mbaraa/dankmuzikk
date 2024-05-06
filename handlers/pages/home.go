package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	pages.Index(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
