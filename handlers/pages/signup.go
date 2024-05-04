package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
)

func HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	pages.Signup(isMobile(r)).Render(context.Background(), w)
}
