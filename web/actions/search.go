package actions

import (
	"dankmuzikk-web/requests"
	"net/http"
)

func (a *Actions) SearchYouTube(query string) ([]Song, error) {
	return requests.Do[any, []Song](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/search",
		QueryParams: map[string]string{
			"query": query,
		},
	})
}

func (a *Actions) SearchYouTubeSuggestions(query string) ([]string, error) {
	return requests.Do[any, []string](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/search/suggestions",
		QueryParams: map[string]string{
			"query": query,
		},
	})
}
