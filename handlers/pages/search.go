package pages

import (
	"context"
	"dankmuzikk/components/pages"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube"
	"net/http"
)

func HandleSearchResultsPage(ytSearch youtube.YouTubeSearcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		results, err := ytSearch.Search(query)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			log.Errorln(err)
			return
		}
		pages.SearchResults(isMobile(r), results).Render(context.Background(), w)
	}
}
