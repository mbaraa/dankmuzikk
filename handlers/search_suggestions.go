package handlers

import (
	"dankmuzikk/services/youtube"
	"encoding/json"
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
		_ = json.NewEncoder(w).Encode(suggessions)
	})
}
