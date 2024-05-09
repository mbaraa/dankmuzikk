package apis

import (
	"context"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube/search/suggestions"
	"dankmuzikk/views/components/search"
	"net/http"
)

func HandleSearchSuggestions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	if q == "" {
		_, _ = w.Write(nil)
		return
	}

	sug, err := suggestions.SearchSuggestions(q)
	if err != nil {
		log.Warningln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(sug) == 0 {
		_, _ = w.Write(nil)
		return
	}

	err = search.SearchSuggestions(sug, q).Render(context.Background(), w)
	if err != nil {
		log.Warningln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
