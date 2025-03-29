package events

type SongsSearched struct {
	Songs []struct {
		YouTubeId    string `json:"youtube_id"`
		Title        string `json:"title"`
		Artist       string `json:"artist"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Duration     string `json:"duration"`
	} `json:"songs"`
}

func (s SongsSearched) Topic() string {
	return "songs-searched"
}

type SongPlayed struct {
	AccountId     uint   `json:"account_id"`
	SongYtId      string `json:"song_yt_id"`
	PlaylistPubId string `json:"playlist_pub_id"`
}

func (s SongPlayed) Topic() string {
	return "song-played"
}

type SongDownloaded struct {
	SongYtId string `json:"song_yt_id"`
}

func (s SongDownloaded) Topic() string {
	return "song-downloaded"
}

type SongAddedToPlaylist struct {
	AccountId     uint   `json:"account_id"`
	PlaylistPubId string `json:"playlist_pub_id"`
	SongYtId      string `json:"song_yt_id"`
}

func (s SongAddedToPlaylist) Topic() string {
	return "song-added-to-playlist"
}

type SongRemovedFromPlaylist struct {
	AccountId     uint   `json:"account_id"`
	PlaylistPubId string `json:"playlist_pub_id"`
	SongYtId      string `json:"song_yt_id"`
}

func (s SongRemovedFromPlaylist) Topic() string {
	return "song-removed-from-playlist"
}
