package actions

func (a *Actions) SearchYouTube(query string) ([]Song, error) {
	return a.requests.SearchYouTube(query)
}

func (a *Actions) SearchYouTubeSuggestions(query string) ([]string, error) {
	return a.requests.SearchYouTubeSuggestions(query)
}
