package search

import (
	"dankmuzikk/config"
	"dankmuzikk/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type scrapSearchResult struct {
	Results []struct {
		Video struct {
			Id           string `json:"id"`
			Title        string `json:"title"`
			Url          string `json:"url"`
			Duration     string `json:"duration"`
			ThumbnailUrl string `json:"thumbnail_src"`
		} `json:"video"`
		Uploader struct {
			Username string `json:"username"`
			Url      string `json:"url"`
		} `json:"uploader"`
	} `json:"results"`
}

// ScraperSearch is a scrapper enabled YouTube search, using the search service under ~/ytscraper
type ScraperSearch struct{}

func (y *ScraperSearch) Search(query string) (results []entities.Song, err error) {
	// TODO: write a proper scraper instead of this hacky node js api
	resp, err := http.Get(fmt.Sprintf("%s/api/search?q=%s", config.Env().YouTube.ScraperUrl, url.QueryEscape(query)))
	if err != nil {
		return
	}

	var apiResults scrapSearchResult
	err = json.NewDecoder(resp.Body).Decode(&apiResults)
	if err != nil {
		return
	}

	for _, res := range apiResults.Results {
		if res.Video.Id == "" || res.Video.Title == "" || res.Video.ThumbnailUrl == "" || res.Uploader.Username == "" {
			continue
		}
		duration := strings.Split(res.Video.Duration, ":")
		if len(duration[0]) == 1 {
			duration[0] = "0" + duration[0]
		}
		if len(duration) == 3 {
			hoursNum, err := strconv.Atoi(duration[0])
			if err != nil {
				continue
			}
			minsNum, err := strconv.Atoi(duration[1])
			if err != nil {
				continue
			}
			if hoursNum >= 1 && minsNum > 30 {
				continue
			}
		}
		if len(duration) > 3 {
			continue
		}

		results = append(results, entities.Song{
			YtId:         res.Video.Id,
			Title:        res.Video.Title,
			Artist:       res.Uploader.Username,
			ThumbnailUrl: res.Video.ThumbnailUrl,
			Duration:     strings.Join(duration, ":"),
		})
	}

	return
}
