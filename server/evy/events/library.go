package events

type SongAddedToFavorites struct {
	SongPublicId string `json:"song_public_id"`
}

func (s SongAddedToFavorites) Topic() string {
	return "song-added-to-favorites"
}
