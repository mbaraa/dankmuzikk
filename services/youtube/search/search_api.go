package search

import (
	"context"
	"dankmuzikk/log"
	"encoding/json"

	"google.golang.org/api/youtube/v3"
)

type apiSearchResult struct {
	Items []struct {
		Id      string `json:"id"`
		Kind    string `json:"kind"`
		Snippet struct {
			Title        string `json:"title"`
			ChannelTitle string `json:"channelTitle"`
			Description  string `json:"description"`
			Thumbnails   struct {
				Default struct {
					Url string `json:"url"`
				} `json:"medium"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// ApiSearch is a YouTube enabled YouTube search, using the official YouTube API.
type ApiSearch struct{}

func (y *ApiSearch) Search(query string) (results []Result, err error) {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx)
	if err != nil {
		return
	}

	response, err := youtubeService.Search.
		List([]string{"id"}).
		Q(query).
		MaxResults(maxSearchResults).Do()
	if err != nil {
		return
	}

	for _, item := range response.Items {
		vid, err := youtubeService.Videos.
			List([]string{"snippet", "contentDetails"}).
			Id(item.Id.VideoId).
			Do()
		if err != nil {
			log.Warningf("[YOUTUBE SERVICE] Fething search results: %s\n", err.Error())
			continue
		}

		var responseObj apiSearchResult
		resJson, _ := vid.MarshalJSON()
		err = json.Unmarshal(resJson, &responseObj)
		if err != nil {
			log.Warningf("[YOUTUBE SERVICE] Unmarshelling the response: %s\n", err.Error())
			continue
		}

		if len(responseObj.Items) == 0 || responseObj.Items[0].Kind != "youtube#video" {
			continue
		}
		duration, err := getTime(responseObj.Items[0].ContentDetails.Duration)
		if err != nil {
			log.Warningf("[YOUTUBE SERVICE] Parsing ISO duration: %s\n", err.Error())
			continue
		}

		results = append(results, Result{
			Title:        responseObj.Items[0].Snippet.Title,
			ChannelTitle: responseObj.Items[0].Snippet.ChannelTitle,
			Description:  responseObj.Items[0].Snippet.Description,
			Id:           responseObj.Items[0].Id,
			ThumbnailUrl: responseObj.Items[0].Snippet.Thumbnails.Default.Url,
			Duration:     duration,
		})
	}

	return
}
