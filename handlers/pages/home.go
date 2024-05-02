package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"net/http"
	"strings"

	_ "github.com/a-h/templ"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	pages.Index(isMobile(r)).Render(context.Background(), w)
}

func isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}
