package handlers

import (
	"context"
	"dankmuzikk/components/ui/search"
	"dankmuzikk/services/youtube"
	"net/http"
)

func HandleSearchSugessions(hand *http.ServeMux) {
	hand.HandleFunc("/api/search-suggession", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("query")
		suggessions, err := youtube.SearchSuggestions(q)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = search.SearchSuggestions(suggessions, q).Render(context.Background(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
