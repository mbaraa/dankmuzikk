package entities

import "dankmuzikk-web/models"

type Playlist struct {
	PublicId    string                     `json:"public_id"`
	Title       string                     `json:"title"`
	SongsCount  int                        `json:"songs_count"`
	Songs       []Song                     `json:"songs"`
	IsPublic    bool                       `json:"is_public"`
	Permissions models.PlaylistPermissions `json:"permissions"`
}
