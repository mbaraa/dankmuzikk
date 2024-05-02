package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	pages.Privacy(isMobile(r)).Render(context.Background(), w)
}
