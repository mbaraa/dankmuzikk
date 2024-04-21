package handlers

import (
	"dankmuzikk/components/pages"
	"net/http"

	"github.com/a-h/templ"
)

func HandleHomePage(hand *http.ServeMux) {
	hand.Handle("/", templ.Handler(pages.Index()))
}
