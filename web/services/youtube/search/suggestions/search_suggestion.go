package suggestions

import (
	"dankmuzikk-web/services/requests"
	"net/url"
)

func SearchSuggestions(query string) ([]string, error) {
	respBody, err := requests.GetRequest[[]string]("/v1/search/suggestions?query=" + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
