package apis

import (
	"dankmuzikk/log"
	"dankmuzikk/services/youtube/search/suggestions"
	"encoding/json"
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

	_ = json.NewEncoder(w).Encode(sug)
}
