package events

import "time"

type PlaylistDownloaded struct {
	PlaylistTitle string    `json:"playlist_title"`
	DeleteAt      time.Time `json:"delete_at"`
}

func (p PlaylistDownloaded) Topic() string {
	return "playlist-downloaded"
}
