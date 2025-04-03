package actions

import (
	"fmt"
	"time"
)

type Song struct {
	PublicId        string        `json:"public_id"`
	Title           string        `json:"title"`
	Artist          string        `json:"artist"`
	ThumbnailUrl    string        `json:"thumbnail_url"`
	RealDuration    time.Duration `json:"duration"`
	PlayTimes       int           `json:"play_times"`
	Votes           int           `json:"votes"`
	AddedAt         string        `json:"added_at"`
	FullyDownloaded bool          `json:"fully_downloaded"`
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
