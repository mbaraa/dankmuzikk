package apis

import (
	"context"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube"
	"dankmuzikk/views/components/search"
	"net/http"
)

func HandleSearchSugessions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	if q == "" {
		w.Write(nil)
		return
	}

	suggessions, err := youtube.SearchSuggestions(q)
	if err != nil {
		log.Warningln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = search.SearchSuggestions(suggessions, q).Render(context.Background(), w)
	if err != nil {
		log.Warningln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
