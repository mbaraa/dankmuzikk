package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	pages.Login(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
