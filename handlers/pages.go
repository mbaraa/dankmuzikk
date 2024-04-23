package handlers

import (
	"context"
	"dankmuzikk/components/pages"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube"
	"net/http"
	"strings"

	_ "github.com/a-h/templ"
)

func HandleHomePage(hand *http.ServeMux) {
	hand.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Index(isMobile(r)).Render(context.Background(), w)
	})
}

func HandleSearchResultsPage(hand *http.ServeMux) {
	hand.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		results, err := youtube.Search(query)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			log.Errorln(err)
			return
		}
		pages.SearchResults(isMobile(r), results).Render(context.Background(), w)
	})
}

func isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}
