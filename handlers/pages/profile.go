package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	pages.Profile(isMobile(r)).Render(context.Background(), w)
}
