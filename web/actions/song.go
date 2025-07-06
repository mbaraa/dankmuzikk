package actions

import (
	"dankmuzikk-web/requests"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
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

type PlaySongParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) PlaySong(params PlaySongParams) (Song, error) {
	return requests.Do[any, Song](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/song/play",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id": params.SongPublicId,
		},
	})
}

type PlaySongFromPlaylistParams struct {
	ActionContext
	SongPublicId     string
	PlaylistPublicId string
}

func (a *Actions) PlaySongFromPlaylist(params PlaySongFromPlaylistParams) (Song, error) {
	return requests.Do[any, Song](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/song/play/playlist",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id":          params.SongPublicId,
			"playlist-id": params.PlaylistPublicId,
		},
	})
}

type PlaySongFromFavoritesParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) PlaySongFromFavorites(params PlaySongFromFavoritesParams) (Song, error) {
	return requests.Do[any, Song](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/song/play/favorites",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id": params.SongPublicId,
		},
	})
}

type PlaySongFromQueueParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) PlaySongFromQueue(params PlaySongFromQueueParams) (Song, error) {
	return requests.Do[any, Song](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/song/play/queue",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id": params.SongPublicId,
		},
	})
}

type GetSongMetadataParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) GetSongMetadata(params GetSongMetadataParams) (Song, error) {
	return requests.Do[any, Song](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/song",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id": params.SongPublicId,
		},
	})
}

type ToggleSongInPlaylistParams struct {
	ActionContext
	SongPublicId     string
	PlaylistPublicId string
}

func (a *Actions) ToggleSongInPlaylist(params ToggleSongInPlaylistParams) (added bool, err error) {
	resp, err := requests.Do[any, map[string]bool](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/playlist/song",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"song-id":     params.SongPublicId,
			"playlist-id": params.PlaylistPublicId,
		},
	})
	if err != nil {
		return false, err
	}

	return resp["added"], nil
}

type UpvoteSongInPlaylistParams struct {
	ActionContext
	SongPublicId     string
	PlaylistPublicId string
}

type UpvoteSongInPlaylistPayload struct {
	VotesCount int `json:"votes_count"`
}

func (a *Actions) UpvoteSongInPlaylist(params UpvoteSongInPlaylistParams) (UpvoteSongInPlaylistPayload, error) {
	return requests.Do[any, UpvoteSongInPlaylistPayload](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/playlist/song/upvote",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"song-id":     params.SongPublicId,
			"playlist-id": params.PlaylistPublicId,
		},
	})
}

type DownvoteSongInPlaylistParams struct {
	ActionContext
	SongPublicId     string
	PlaylistPublicId string
}

type DownvoteSongInPlaylistPayload struct {
	VotesCount int `json:"votes_count"`
}

func (a *Actions) DownvoteSongInPlaylist(params DownvoteSongInPlaylistParams) (DownvoteSongInPlaylistPayload, error) {
	return requests.Do[any, DownvoteSongInPlaylistPayload](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/playlist/song/downvote",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"song-id":     params.SongPublicId,
			"playlist-id": params.PlaylistPublicId,
		},
	})
}

type GetLyricsForSongPayload struct {
	SongTitle string            `json:"song_title"`
	Lyrics    []string          `json:"lyrics"`
	Synced    map[string]string `json:"synced"`
}

func (l GetLyricsForSongPayload) SyncedPairs() []struct {
	K string
	V string
} {
	pairs := make([]struct {
		K string
		V string
	}, 0, len(l.Synced))
	for ts, part := range l.Synced {
		pairs = append(pairs, struct {
			K string
			V string
		}{ts, part})
	}

	slices.SortFunc(pairs, func(pairI, pairJ struct {
		K string
		V string
	}) int {
		return strings.Compare(pairI.K, pairJ.K)
	})

	return pairs
}

func (a *Actions) GetSongLyrics(songPublicId string) (GetLyricsForSongPayload, error) {
	return requests.Do[any, GetLyricsForSongPayload](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/song/lyrics",
		QueryParams: map[string]string{
			"id": songPublicId,
		},
	})
}

type GetFavoritesParams struct {
	ActionContext
	PageIndex uint
}

type GetFavoritesPayload struct {
	Songs []Song `json:"songs"`
}

func (a *Actions) GetFavorites(params GetFavoritesParams) (GetFavoritesPayload, error) {
	return requests.Do[any, GetFavoritesPayload](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/library/favorite/songs",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"page": strconv.Itoa(int(params.PageIndex)),
		},
	})
}

type AddSongToFavoritesParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) AddSongToFavorites(params AddSongToFavoritesParams) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPost,
		Endpoint: "/v1/library/favorite/song",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"id": params.SongPublicId,
		},
	})

	return err
}

type RemoveSongFromFavoritesParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) RemoveSongFromFavorites(params AddSongToFavoritesParams) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodDelete,
		Endpoint: "/v1/library/favorite/song",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"id": params.SongPublicId,
		},
	})

	return err
}
