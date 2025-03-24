package events

import "time"

type PlaylistDownloaded struct {
	PlaylistId string    `json:"playlist_id"`
	DeleteAt   time.Time `json:"delete_at"`
}

func (p PlaylistDownloaded) Topic() string {
	return "playlist-downloaded"
}
