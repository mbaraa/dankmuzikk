package requests

import (
	"dankmuzikk-web/actions"
)

type Requests struct {
}

func New() *Requests {
	return &Requests{}
}

type createPlaylistResponse struct {
	NewPlaylist actions.Playlist `json:"new_playlist"`
}
