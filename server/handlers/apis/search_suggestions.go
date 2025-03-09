package apis

import (
	"dankmuzikk/actions"
	"encoding/json"
	"net/http"
)

type searchApi struct {
	usecases *actions.Actions
}

func NewYouTubeSearchApi(usecases *actions.Actions) *searchApi {
	return &searchApi{
		usecases: usecases,
	}
}

func (s *searchApi) HandleSearchSuggestions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	if q == "" {
		_, _ = w.Write(nil)
		return
	}

	sug, err := s.usecases.SearchSuggestions(q)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	if len(sug) == 0 {
		_, _ = w.Write(nil)
		return
	}

	_ = json.NewEncoder(w).Encode(sug)
}

func (s *searchApi) HandleSearchResults(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results, err := s.usecases.SearchYouTube(query)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(results)
}
