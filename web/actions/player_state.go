package actions

import (
	"dankmuzikk-web/requests"
	"net/http"
	"strconv"
)

type PlayerState struct {
	Shuffled         bool   `json:"shuffled"`
	CurrentSongIndex int    `json:"current_song_index"`
	LoopMode         string `json:"loop_mode"`
	Songs            []Song `json:"songs"`
}

type GetPlayerStatePayload struct {
	PlayerState PlayerState `json:"player_state"`
}

func (a *Actions) GetPlayerState(ctx ActionContext) (GetPlayerStatePayload, error) {
	return requests.Do[any, GetPlayerStatePayload](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/player",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})
}

func (a *Actions) SetPlayerShuffleOn(ctx ActionContext) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPost,
		Endpoint: "/v1/player/shuffle",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})

	return err
}

func (a *Actions) SetPlayerShuffleOff(ctx ActionContext) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodDelete,
		Endpoint: "/v1/player/shuffle",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})

	return err
}

func (a *Actions) SetPlayerLoopOff(ctx ActionContext) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/player/loop/off",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})

	return err
}

func (a *Actions) SetPlayerLoopOnce(ctx ActionContext) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/player/loop/once",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})

	return err
}

func (a *Actions) SetPlayerLoopAll(ctx ActionContext) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/player/loop/all",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})

	return err
}

type AddSongToQueueNextParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) AddSongToQueueNext(params AddSongToQueueNextParams) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPost,
		Endpoint: "/v1/player/queue/song/next",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id": params.SongPublicId,
		},
	})

	return err
}

type AddSongToQueueAtLastParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) AddSongToQueueAtLast(params AddSongToQueueNextParams) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPost,
		Endpoint: "/v1/player/queue/song/last",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id": params.SongPublicId,
		},
	})

	return err
}

type RemoveSongFromQueueParams struct {
	ActionContext
	SongIndex int
}

func (a *Actions) RemoveSongFromQueue(params RemoveSongFromQueueParams) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodDelete,
		Endpoint: "/v1/player/queue/song",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"index": strconv.Itoa(params.SongIndex),
		},
	})

	return err
}

type AddPlaylistToQueueNextParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) AddPlaylistToQueueNext(params AddPlaylistToQueueNextParams) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPost,
		Endpoint: "/v1/player/queue/playlist/next",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id": params.PlaylistPublicId,
		},
	})

	return err
}

type AddPlaylistToQueueAtLastParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) AddPlaylistToQueueAtLast(params AddPlaylistToQueueAtLastParams) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodPost,
		Endpoint: "/v1/player/queue/playlist/last",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
			"X-Client-Hash": params.ClientHash,
		},
		QueryParams: map[string]string{
			"id": params.PlaylistPublicId,
		},
	})

	return err
}

type GetNextSongInQueuePayload struct {
	Song             Song `json:"song"`
	CurrentSongIndex int  `json:"current_song_index"`
	EndOfQueue       bool `json:"end_of_queue"`
}

func (a *Actions) GetNextSongInQueue(ctx ActionContext) (GetNextSongInQueuePayload, error) {
	return requests.Do[any, GetNextSongInQueuePayload](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/player/song/next",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})
}

type GetPreviousSongInQueuePayload struct {
	Song             Song `json:"song"`
	CurrentSongIndex int  `json:"current_song_index"`
	EndOfQueue       bool `json:"end_of_queue"`
}

func (a *Actions) GetPreviousSongInQueue(ctx ActionContext) (GetPreviousSongInQueuePayload, error) {
	return requests.Do[any, GetPreviousSongInQueuePayload](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/player/song/previous",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})
}

func (a *Actions) GetPlayingSongLyrics(ctx ActionContext) (GetLyricsForSongPayload, error) {
	return requests.Do[any, GetLyricsForSongPayload](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/player/song/lyrics",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
			"X-Client-Hash": ctx.ClientHash,
		},
		QueryParams: map[string]string{},
	})
}
