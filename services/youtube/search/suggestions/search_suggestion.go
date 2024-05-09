package suggestions

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func SearchSuggestions(query string) (suggestions []string, err error) {
	resp, err := http.Get("https://suggestqueries.google.com/complete/search?client=firefox&ds=yt&q=" +
		url.QueryEscape(query))
	if err != nil {
		panic(err)
	}

	var results []any
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return
	}

	for i, res := range results[1].([]any) {
		if i >= 9 { // max displayed suggestions is 10
			break
		}
		suggestions = append(suggestions, res.(string))
	}

	return
}
