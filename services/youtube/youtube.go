package youtube

import (
	"context"
	"dankmuzikk/log"
	"os"

	"google.golang.org/api/youtube/v3"
)

func init() {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatalln(log.ErrorLevel, "[YOUTUBE SERVICE] Missing Google API Service Account File")
	}
}

type SearchResult struct {
	Title        string
	Url          string
	ThumbnailUrl string
}

func Search(query string) (res []SearchResult, err error) {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx)
	if err != nil {
		return
	}

	response, err := youtubeService.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(5).Do()
	if err != nil {
		return
	}

	for _, item := range response.Items {
		if item.Id.Kind != "youtube#video" {
			continue
		}
		res = append(res, SearchResult{
			Title:        item.Snippet.Title,
			Url:          item.Id.VideoId,
			ThumbnailUrl: item.Snippet.Thumbnails.Default.Url,
		})
	}
	return
}
