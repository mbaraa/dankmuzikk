package youtube

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func SearchSuggestions(query string) (sugesstions []string, err error) {
	resp, err := http.Get("http://suggestqueries.google.com/complete/search?client=firefox&ds=yt&q=" +
		url.QueryEscape(query))
	if err != nil {
		panic(err)
	}

	var results []any
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		panic(err)
	}

	for i, res := range results[1].([]any) {
		if i >= 9 { // max displayed suggestions is 10
			break
		}
		sugesstions = append(sugesstions, res.(string))
	}

	return
}
