package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
)

func (r *Requests) SearchYouTube(query string) ([]actions.Song, error) {
	return makeRequest[any, []actions.Song](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/search",
		queryParams: map[string]string{
			"query": query,
		},
	})
}

func (r *Requests) SearchYouTubeSuggestions(query string) ([]string, error) {
	return makeRequest[any, []string](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/search/suggestions",
		queryParams: map[string]string{
			"query": query,
		},
	})
}
