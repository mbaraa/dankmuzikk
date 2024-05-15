package entities

type Playlist struct {
	PublicId   string `json:"public_id"`
	Title      string `json:"title"`
	SongsCount int    `json:"songs_count"`
	Songs      []Song `json:"songs"`
}
