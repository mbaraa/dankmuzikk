package youtube

import (
	"dankmuzikk/actions"
	"dankmuzikk/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type YouTube struct{}

func New() *YouTube {
	return &YouTube{}
}

var (
	keyPattern   = regexp.MustCompile(`"innertubeApiKey":"([^"]*)`)
	dataPattern  = regexp.MustCompile(`ytInitialData[^{]*(.*?);\s*<\/script>`)
	dataPattern2 = regexp.MustCompile(`ytInitialData"[^{]*(.*);\s*window\["ytInitialPlayerResponse"\]`)
)

type videoResult struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	Duration     string `json:"duration"`
	ThumbnailUrl string `json:"thumbnail_src"`
	Uploader     string `json:"username"`
}

type req struct {
	Version          string `json:"version"`
	Parser           string `json:"parser"`
	Key              string `json:"key"`
	EstimatedResults string `json:"estimatedResults"`
}

type videoRenderer struct {
	VideoId string `json:"videoId"`
	Title   struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"title"`
	LengthText struct {
		SimpleText string `json:"simpleText"`
	} `json:"lengthText"`
	Thumbnail struct {
		Thumbnails []struct {
			URL string `json:"url"`
		} `json:"thumbnails"`
	} `json:"thumbnail"`
	OwnerText struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"ownerText"`
}

type ytSearchData struct {
	EstimatedResults string `json:"estimatedResults"`
	Contents         struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct { // sectionList
						ItemSectionRenderer struct {
							Contents []struct {
								ChannelRenderer any `json:"channelRenderer"`

								VideoRenderer videoRenderer `json:"videoRenderer"`

								RadioRenderer    any `json:"radioRenderer"`
								PlaylistRenderer any `json:"playlistRenderer"`
							} `json:"contents"`
						} `json:"itemSectionRenderer"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

func search(q string) ([]videoResult, error) {
	// get ze results
	url := "https://www.youtube.com/results?q=" + url.QueryEscape(q)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	jojo := req{
		Version: "0.1.5",
		Parser:  "json_format",
		Key:     "",
	}
	key := keyPattern.FindSubmatch(respBody)
	jojo.Key = string(key[1])

	matches := dataPattern.FindSubmatch(respBody)
	if len(matches) > 1 {
		jojo.Parser += ".object_var"
	} else {
		jojo.Parser += ".original"
		matches = dataPattern2.FindSubmatch(respBody)
	}
	data := ytSearchData{}
	err = json.Unmarshal(matches[1], &data)
	if err != nil {
		return nil, err
	}
	jojo.EstimatedResults = data.EstimatedResults

	// parse JSON data

	resSuka := make([]videoResult, 0)
	for _, sectionList := range data.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents {
		for _, content := range sectionList.ItemSectionRenderer.Contents {
			_ = content
			if content.VideoRenderer.VideoId == "" {
				continue
			}
			resSuka = append(resSuka, videoResult{
				Id:           content.VideoRenderer.VideoId,
				Title:        content.VideoRenderer.Title.Runs[0].Text,
				Duration:     content.VideoRenderer.LengthText.SimpleText,
				ThumbnailUrl: content.VideoRenderer.Thumbnail.Thumbnails[len(content.VideoRenderer.Thumbnail.Thumbnails)-1].URL,
				Uploader:     content.VideoRenderer.OwnerText.Runs[0].Text,
			})
		}
	}

	return resSuka, nil
}

func (y *YouTube) Search(query string) (results []actions.YouTubeSong, err error) {
	searchHits, err := search(query)
	if err != nil {
		return nil, err
	}

	for _, res := range searchHits {
		if res.Id == "" || res.Title == "" || res.ThumbnailUrl == "" || res.Uploader == "" || res.Duration == "" {
			continue
		}
		duration := strings.Split(res.Duration, ":")
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

		results = append(results, actions.YouTubeSong{
			YtId:         res.Id,
			Title:        res.Title,
			Artist:       res.Uploader,
			ThumbnailUrl: res.ThumbnailUrl,
			Duration:     strings.Join(duration, ":"),
		})
	}

	return
}

// TODO: use me waa waa
func getTime(isoDuration string) (time.Duration, error) {
	startIdx := 0
	for i, chr := range isoDuration {
		if unicode.IsDigit(chr) {
			startIdx = i
			break
		}
	}
	duration, err := time.ParseDuration(strings.ToLower(isoDuration[startIdx:]))
	if err != nil {
		return 0, err
	}

	return duration, nil
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
// YOUTUBE_MUSIC_DOWNLOAD_PATH, where the file name will be <video_id.mp3> to be served under /muzikk/{id}
// and returns an occurring error
//
// Used when playing a new song (usually from search).
func (y *YouTube) DownloadYoutubeSong(songYtId string) error {
	resp, err := http.Get(fmt.Sprintf("%s/download/%s", config.Env().YouTube.DownloaderUrl, songYtId))
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
