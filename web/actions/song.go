package actions

import (
	"encoding/json"
	"fmt"
	"time"
)

type Song struct {
	PublicId     string        `json:"public_id"`
	Title        string        `json:"title"`
	Artist       string        `json:"artist"`
	ThumbnailUrl string        `json:"thumbnail_url"`
	RealDuration time.Duration `json:"duration"`
	PlayTimes    int           `json:"play_times"`
	Votes        int           `json:"votes,omitempty"`
	AddedAt      string        `json:"added_at,omitempty"`
	Favorite     bool          `json:"favorite"`
	MediaUrl     string        `json:"media_url"`
}

type fakeSong struct {
	PublicId     string  `json:"public_id"`
	Title        string  `json:"title"`
	Artist       string  `json:"artist"`
	ThumbnailUrl string  `json:"thumbnail_url"`
	Duration     float64 `json:"duration"`
	PlayTimes    int     `json:"play_times"`
	Votes        int     `json:"votes"`
	AddedAt      string  `json:"added_at"`
	Favorite     bool    `json:"favorite"`
	MediaUrl     string  `json:"media_url"`
}

func (s *Song) UnmarshalJSON(data []byte) error {
	helper := fakeSong{}
	if err := json.Unmarshal(data, &helper); err != nil {
		return err
	}
	s.PublicId = helper.PublicId
	s.Title = helper.Title
	s.Artist = helper.Artist
	s.ThumbnailUrl = helper.ThumbnailUrl
	s.PlayTimes = helper.PlayTimes
	s.Votes = helper.Votes
	s.AddedAt = helper.AddedAt
	s.RealDuration = time.Duration(helper.Duration) * time.Nanosecond
	s.Favorite = helper.Favorite
	s.MediaUrl = helper.MediaUrl
	return nil
}

func (song Song) Duration() string {
	s := int64(song.RealDuration.Seconds())
	if s < 60 {
		return fmt.Sprintf("00:%02d", s)
	}
	m := s / 60
	s %= 60
	if m < 60 {
		if s == 0 {
			return fmt.Sprintf("%02d:00", m)
		}
		return fmt.Sprintf("%02d:%02d", m, s)
	}
	h := m / 60
	m %= 60
	if h < 24 {
		if m == 0 && s == 0 {
			return fmt.Sprintf("%02d:00:00", h)
		} else if s == 0 {
			return fmt.Sprintf("%02d:%02d:00", h, m)
		}
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	days := h / 24
	h %= 24
	if h == 0 && m == 0 && s == 0 {
		return fmt.Sprintf("%02d:00:00:00", days)
	} else if m == 0 && s == 0 {
		return fmt.Sprintf("%02d:%02d:00:00", days, h)
	} else if s == 0 {
		return fmt.Sprintf("%02d:%02d:%02d:00", days, h, m)
	}

	return fmt.Sprintf("%02d:%02d:%02d:%02d:00", days, h, m, s)
}

func (a *Actions) PlaySong(sessionToken, clientHash, songPublicId string) (Song, error) {
	return a.requests.PlaySong(sessionToken, clientHash, songPublicId)
}

func (a *Actions) PlaySongFromPlaylist(sessionToken, clientHash, songPublicId, playlistPublicId string) (Song, error) {
	return a.requests.PlaySongFromPlaylist(sessionToken, clientHash, songPublicId, playlistPublicId)
}

func (a *Actions) PlaySongFromFavorites(sessionToken, clientHash, songPublicId string) (Song, error) {
	return a.requests.PlaySongFromFavorites(sessionToken, clientHash, songPublicId)
}

func (a *Actions) PlaySongFromQueue(sessionToken, clientHash, songPublicId string) (Song, error) {
	return a.requests.PlaySongFromQueue(sessionToken, clientHash, songPublicId)
}

func (a *Actions) GetSongMetadata(sessionToken, clientHash, songPublicId string) (Song, error) {
	return a.requests.GetSongMetadata(sessionToken, clientHash, songPublicId)
}

func (a *Actions) ToggleSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (added bool, err error) {
	return a.requests.ToggleSongInPlaylist(sessionToken, songPublicId, playlistPublicId)
}

type UpvoteSongInPlaylistPayload struct {
	VotesCount int `json:"votes_count"`
}

func (a *Actions) UpvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (UpvoteSongInPlaylistPayload, error) {
	return a.requests.UpvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId)
}

type DownvoteSongInPlaylistPayload struct {
	VotesCount int `json:"votes_count"`
}

func (a *Actions) DownvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (DownvoteSongInPlaylistPayload, error) {
	return a.requests.DownvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId)
}

type GetLyricsForSongPayload struct {
	SongTitle string            `json:"song_title"`
	Lyrics    []string          `json:"lyrics"`
	Synced    map[string]string `json:"synced"`
}

func (a *Actions) GetSongLyrics(songPublicId string) (GetLyricsForSongPayload, error) {
	return a.requests.GetSongLyrics(songPublicId)
}

type GetFavoritesPayload struct {
	Songs []Song `json:"songs"`
}

func (a *Actions) GetFavorites(sessionToken string, pageIndex uint) (GetFavoritesPayload, error) {
	return a.requests.GetFavorites(sessionToken, pageIndex)
}

func (a *Actions) AddSongToFavorites(sessionToken string, songPublicId string) error {
	return a.requests.AddSongToFavorites(sessionToken, songPublicId)
}

func (a *Actions) RemoveSongFromFavorites(sessionToken string, songPublicId string) error {
	return a.requests.RemoveSongFromFavorites(sessionToken, songPublicId)
}
