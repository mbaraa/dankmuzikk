package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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

type YouTubeScraperSearch struct{}

func (y *YouTubeScraperSearch) Search(query string) (results []SearchResult, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/search?q=%s", os.Getenv("YOUTUBE_SCAPER_URL"), url.QueryEscape(query)))
	if err != nil {
		return
	}

	var apiResults scrapSearchResult
	err = json.NewDecoder(resp.Body).Decode(&apiResults)
	if err != nil {
		return
	}

	for _, res := range apiResults.Results {
		duration := strings.Split(res.Video.Duration, ":")
		if len(duration[0]) == 1 {
			duration[0] = "0" + duration[0]
		}

		results = append(results, SearchResult{
			Id:           res.Video.Id,
			Title:        res.Video.Title,
			ChannelTitle: res.Uploader.Username,
			Description:  "",
			ThumbnailUrl: res.Video.ThumbnailUrl,
			Duration:     strings.Join(duration, ":"),
		})
	}

	return
}
