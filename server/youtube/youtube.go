package youtube

import (
	"dankmuzikk/actions"
	"dankmuzikk/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mbaraa/ytscrape"
)

type YouTube struct{}

func New() *YouTube {
	return &YouTube{}
}

func (y *YouTube) Search(query string) (results []actions.YouTubeSong, err error) {
	searchHits, err := ytscrape.Search(query)
	if err != nil {
		return nil, err
	}

	for _, res := range searchHits {
		if res.Id == "" || res.Title == "" || res.ThumbnailUrl == "" || res.Uploader.Title == "" || res.Duration == 0 {
			continue
		}

		if res.Duration > time.Hour*2 {
			continue
		}

		results = append(results, actions.YouTubeSong{
			YtId:         res.Id,
			Title:        res.Title,
			Artist:       res.Uploader.Title,
			ThumbnailUrl: res.ThumbnailUrl,
			Duration:     res.Duration,
		})
	}

	return
}

func (y *YouTube) SearchSuggestions(query string) (suggestions []string, err error) {
	resp, err := http.Get("https://suggestqueries.google.com/complete/search?client=firefox&ds=yt&q=" +
		url.QueryEscape(query))
	if err != nil {
		return nil, err
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

// DownloadYoutubeSong downloads a YouTube muzikk file into the path specified by the environment variable
// BLOBS_DIR, where the file name will be <video_id.mp3> to be served under /muzikk/{id}
// and returns an occurring error
//
// Used when playing a new song (usually from search).
func (y *YouTube) DownloadYoutubeSong(songYtId string) error {
	resp, err := http.Get(fmt.Sprintf("%s/download/%s", config.Env().YouTubeDownloaderAddress, songYtId))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("something went wrong when downloading a song; id: " + songYtId)
	}

	respBody := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	if respBody["error"] != "" {
		return errors.New(respBody["error"])
	}

	return nil
}
