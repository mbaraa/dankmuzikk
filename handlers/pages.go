package handlers

import (
	"dankmuzikk/components/pages"
	"net/http"

	"github.com/a-h/templ"
)

func HandleHomePage() http.Handler {
	hand := http.NewServeMux()
	hand.Handle("/", templ.Handler(pages.Index()))
	return hand
}
