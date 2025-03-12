package search

import (
	"dankmuzikk-web/entities"
	"dankmuzikk-web/services/requests"
	"net/url"
)

// Service is an interface that represents a YouTube search behavior.
type Service interface {
	// Search searches YouTube for the given query,
	// and returns a SearchResult slice, and an occurring error.
	Search(query string) (results []entities.Song, err error)
}

// SearchImpl is a scrapper enabled YouTube search.
type SearchImpl struct{}

func (y *SearchImpl) Search(query string) ([]entities.Song, error) {
	songs, err := requests.GetRequest[[]entities.Song]("/v1/search?query=" + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}

	return songs, nil
}
