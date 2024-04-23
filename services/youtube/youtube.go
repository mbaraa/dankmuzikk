package youtube

import (
	"context"
	"dankmuzikk/log"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/api/youtube/v3"
)

const maxSearchResults = 10

func init() {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatalln(log.ErrorLevel, "[YOUTUBE SERVICE] Missing Google API Service Account File")
	}
}

type SearchResult struct {
	Title        string
	Url          string
	ThumbnailUrl string
	Duration     string
}

type apiSearchResult struct {
	Items []struct {
		Id      string `json:"id"`
		Kind    string `json:"kind"`
		Snippet struct {
			Title      string `json:"title"`
			Thumbnails struct {
				Default struct {
					Url string `json:"url"`
				} `json:"default"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

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

func Search(query string) (results []SearchResult, err error) {
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

		results = append(results, SearchResult{
			Title:        responseObj.Items[0].Snippet.Title,
			Url:          responseObj.Items[0].Id,
			ThumbnailUrl: responseObj.Items[0].Snippet.Thumbnails.Default.Url,
			Duration:     duration,
		})
	}

	return
}

func getTime(isoDuration string) (string, error) {
	duration, err := time.ParseDuration(strings.ToLower(isoDuration[2:]))
	if err != nil {
		return "", err
	}
	days, hours, mins, secs :=
		duration/(time.Hour*24), (duration / time.Hour), duration/time.Minute, duration/time.Second

	builder := strings.Builder{}
	if days > 0 {
		builder.WriteString(formatNumber(int(days)) + ":")
	}
	if hours > 0 {
		builder.WriteString(formatNumber(int(hours%60)) + ":")
	}
	if mins > 0 {
		builder.WriteString(formatNumber(int(mins%60)) + ":")
	}
	if secs > 0 {
		if days == 0 && hours == 0 && mins == 0 {
			builder.WriteString(formatNumber(int(secs%60)) + "s")
		} else {
			builder.WriteString(formatNumber(int(secs % 60)))
		}
	}

	return builder.String(), nil
}

func formatNumber(n int) string {
	if n < 10 {
		return fmt.Sprintf("0%d", n)
	}
	return fmt.Sprint(n)
}
