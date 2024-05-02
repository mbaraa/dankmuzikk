package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandleTOSPage(w http.ResponseWriter, r *http.Request) {
	pages.TOS(isMobile(r)).Render(context.Background(), w)
}
