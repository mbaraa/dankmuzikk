package apis

import (
	"context"
	"dankmuzikk-web/actions"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/search"
	"net/http"
)

type searchSuggestionsApi struct {
	usecases *actions.Actions
}

func NewSearchSiggestionsApi(usecases *actions.Actions) *searchSuggestionsApi {
	return &searchSuggestionsApi{
		usecases: usecases,
	}
}

func (s *searchSuggestionsApi) HandleSearchSuggestions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	if q == "" {
		_, _ = w.Write(nil)
		return
	}

	sug, err := s.usecases.SearchYouTubeSuggestions(q)
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
