package events

import "time"

type SongsSearched struct {
	Songs []struct {
		YouTubeId    string        `json:"youtube_id"`
		Title        string        `json:"title"`
		Artist       string        `json:"artist"`
		ThumbnailUrl string        `json:"thumbnail_url"`
		Duration     time.Duration `json:"duration"`
	} `json:"songs"`
}

func (s SongsSearched) Topic() string {
	return "songs-searched"
}

type SongPlayedEntryPoint string

const (
	SingleSongEntryPoint SongPlayedEntryPoint = "single"
	// TODO: remove in favor of a standalone action
	PlayPlaylistEntryPoint SongPlayedEntryPoint = "play-playlist"
	FromPlaylistEntryPoint SongPlayedEntryPoint = "from-playlist"
	FavoriteSongEntryPoint SongPlayedEntryPoint = "favorite-song"
	QueueSongEntryPoint    SongPlayedEntryPoint = "from-queue"
)

type SongPlayed struct {
	AccountId        uint                 `json:"account_id"`
	ClientHash       string               `json:"client_hash"`
	SongPublicId     string               `json:"song_public_id"`
	PlaylistPublicId string               `json:"playlist_public_id"`
	EntryPoint       SongPlayedEntryPoint `json:"entry_point"`
}

func (s SongPlayed) Topic() string {
	return "song-played"
}

type SongDownloaded struct {
	SongPublicId string `json:"song_public_id"`
}

func (s SongDownloaded) Topic() string {
	return "song-downloaded"
}

type SongAddedToPlaylist struct {
	AccountId     uint   `json:"account_id"`
	PlaylistPubId string `json:"playlist_pub_id"`
	SongPublicId  string `json:"song_public_id"`
}

func (s SongAddedToPlaylist) Topic() string {
	return "song-added-to-playlist"
}

type SongRemovedFromPlaylist struct {
	AccountId     uint   `json:"account_id"`
	PlaylistPubId string `json:"playlist_pub_id"`
	SongPublicId  string `json:"song_public_id"`
}

func (s SongRemovedFromPlaylist) Topic() string {
	return "song-removed-from-playlist"
}
