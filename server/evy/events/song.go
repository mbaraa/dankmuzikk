package events

type SongPlayed struct {
	ProfileId     uint   `json:"profile_id"`
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
	PlaylistPubId string `json:"playlist_pub_id"`
	SongYtId      string `json:"song_yt_id"`
}

func (s SongAddedToPlaylist) Topic() string {
	return "song-added-to-playlist"
}

type SongRemovedFromPlaylist struct {
	PlaylistPubId string `json:"playlist_pub_id"`
	SongYtId      string `json:"song_yt_id"`
}

func (s SongRemovedFromPlaylist) Topic() string {
	return "song-removed-from-playlist"
}
