package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	pages.Privacy(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
