package actions

import "dankmuzikk/app/models"

type YouTubeSong struct {
	YtId         string `json:"yt_id"`
	Title        string `json:"title"`
	Artist       string `json:"artist"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Duration     string `json:"duration"`
	PlayTimes    int    `json:"play_times"`
	Votes        int    `json:"votes"`
	AddedAt      string `json:"added_at"`
}

type YouTube interface {
	Search(query string) (results []YouTubeSong, err error)
	SearchSuggestions(query string) (suggestions []string, err error)

	DownloadYoutubeSong(songYtId string) error
}

func (a *Actions) SearchSuggestions(q string) ([]string, error) {
	return a.youtube.SearchSuggestions(q)
}

func (a *Actions) SearchYouTube(q string) ([]YouTubeSong, error) {
	results, err := a.youtube.Search(q)
	if err != nil {
		return nil, err
	}

	for _, newSong := range results {
		_, _ = a.app.CreateSong(models.Song{
			YtId:            newSong.YtId,
			Title:           newSong.Title,
			Artist:          newSong.Artist,
			ThumbnailUrl:    newSong.ThumbnailUrl,
			Duration:        newSong.Duration,
			FullyDownloaded: false,
		})
	}

	return results, nil
}
