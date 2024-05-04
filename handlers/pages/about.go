package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	pages.About(isMobile(r)).Render(context.Background(), w)
}
