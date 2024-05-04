package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"

	_ "github.com/a-h/templ"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	pages.Index(isMobile(r)).Render(context.Background(), w)
}
