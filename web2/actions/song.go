package actions

type Song struct {
	YtId            string `json:"yt_id"`
	Title           string `json:"title"`
	Artist          string `json:"artist"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	Duration        string `json:"duration"`
	PlayTimes       int    `json:"play_times"`
	Votes           int    `json:"votes"`
	AddedAt         string `json:"added_at"`
	FullyDownloaded bool   `json:"fully_downloaded"`
}

func (a *Actions) PlaySong(sessionToken, songPublicId, playlistPublicId string) (string, error) {
	return a.requests.PlaySong(sessionToken, songPublicId, playlistPublicId)
}

func (a *Actions) GetSongMetadata(sessionToken, songPublicId string) (Song, error) {
	return a.requests.GetSongMetadata(sessionToken, songPublicId)
}

func (a *Actions) ToggleSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (added bool, err error) {
	return a.requests.ToggleSongInPlaylist(sessionToken, songPublicId, playlistPublicId)
}

func (a *Actions) UpvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (votesCount int, err error) {
	return a.requests.UpvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId)
}

func (a *Actions) DownvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (votesCount int, err error) {
	return a.requests.DownvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId)
}

type GetLyricsForSongPayload struct {
	SongTitle string   `json:"song_title"`
	Lyrics    []string `json:"lyrics"`
}

func (a *Actions) GetSongLyrics(songPublicId string) (GetLyricsForSongPayload, error) {
	return a.requests.GetSongLyrics(songPublicId)
}
