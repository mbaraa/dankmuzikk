package actions

import (
	"dankmuzikk/evy/events"
	"time"
)

type YouTubeSong struct {
	YtId         string        `json:"yt_id"`
	Title        string        `json:"title"`
	Artist       string        `json:"artist"`
	ThumbnailUrl string        `json:"thumbnail_url"`
	Duration     time.Duration `json:"duration"`
	PlayTimes    int           `json:"play_times"`
	Votes        int           `json:"votes"`
	AddedAt      string        `json:"added_at"`
}

type YouTube interface {
	Search(query string) (results []YouTubeSong, err error)
	SearchSuggestions(query string) (suggestions []string, err error)

	DownloadYoutubeSong(songYtId string) error
}

func (a *Actions) SearchSuggestions(q string) ([]string, error) {
	return a.youtube.SearchSuggestions(q)
}

func (a *Actions) SearchYouTube(q string) ([]Song, error) {
	results, err := a.youtube.Search(q)
	if err != nil {
		return nil, err
	}

	songs := make([]struct {
		YouTubeId    string        `json:"youtube_id"`
		Title        string        `json:"title"`
		Artist       string        `json:"artist"`
		ThumbnailUrl string        `json:"thumbnail_url"`
		Duration     time.Duration `json:"duration"`
	}, 0, len(results))
	for _, newSong := range results {
		songs = append(songs, struct {
			YouTubeId    string        `json:"youtube_id"`
			Title        string        `json:"title"`
			Artist       string        `json:"artist"`
			ThumbnailUrl string        `json:"thumbnail_url"`
			Duration     time.Duration `json:"duration"`
		}{
			YouTubeId:    newSong.YtId,
			Title:        newSong.Title,
			Artist:       newSong.Artist,
			ThumbnailUrl: newSong.ThumbnailUrl,
			Duration:     newSong.Duration,
		})
	}

	err = a.eventhub.Publish(events.SongsSearched{
		Songs: songs,
	})
	if err != nil {
		return nil, err
	}

	foundSongs := make([]Song, 0, len(results))
	for _, song := range results {
		foundSongs = append(foundSongs, Song{
			PublicId:        song.YtId,
			Title:           song.Title,
			Artist:          song.Artist,
			ThumbnailUrl:    song.ThumbnailUrl,
			Duration:        song.Duration,
			FullyDownloaded: false,
		})
	}

	return foundSongs, nil
}
