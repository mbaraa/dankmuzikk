package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	pages.About(isMobile(r)).Render(context.Background(), w)
}
