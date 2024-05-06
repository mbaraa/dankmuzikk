package pages

import (
	"context"
	"dankmuzikk/views/pages"
	"net/http"
)

func HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	pages.Signup(isMobile(r), getTheme(r)).Render(context.Background(), w)
}
