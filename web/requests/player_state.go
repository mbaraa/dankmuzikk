package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
)

func (r *Requests) GetPlayerState(sessionToken, clientHash string) (actions.GetPlayerStatePayload, error) {
	return makeRequest[any, actions.GetPlayerStatePayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/player",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{},
	})
}

func (r *Requests) SetPlayerShuffleOn(sessionToken, clientHash string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPost,
		endpoint: "/v1/player/shuffle",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{},
	})

	return err
}

func (r *Requests) SetPlayerShuffleOff(sessionToken, clientHash string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodDelete,
		endpoint: "/v1/player/shuffle",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{},
	})

	return err
}

func (r *Requests) SetPlayerLoopOff(sessionToken, clientHash string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/player/loop/off",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{},
	})

	return err

}

func (r *Requests) SetPlayerLoopOnce(sessionToken, clientHash string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/player/loop/once",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{},
	})

	return err
}

func (r *Requests) SetPlayerLoopAll(sessionToken, clientHash string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/player/loop/all",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{},
	})

	return err
}

func (r *Requests) GetNextSongInQueue(sessionToken, clientHash string) (actions.GetNextSongInQueuePayload, error) {
	return makeRequest[any, actions.GetNextSongInQueuePayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/player/song/next",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{},
	})
}

func (r *Requests) GetPreviousSongInQueue(sessionToken, clientHash string) (actions.GetPreviousSongInQueuePayload, error) {
	return makeRequest[any, actions.GetPreviousSongInQueuePayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/player/song/previous",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{},
	})
}

func (r *Requests) AddSongToQueueNext(sessionToken, clientHash, songPublicId string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPost,
		endpoint: "/v1/player/queue/song/next",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})

	return err
}

func (r *Requests) AddSongToQueueAtLast(sessionToken, clientHash, songPublicId string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPost,
		endpoint: "/v1/player/queue/song/last",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})

	return err
}

func (r *Requests) AddPlaylistToQueueNext(sessionToken, clientHash, playlistPublicId string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPost,
		endpoint: "/v1/player/queue/playlist/next",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": playlistPublicId,
		},
	})

	return err
}

func (r *Requests) AddPlaylistToQueueAtLast(sessionToken, clientHash, playlistPublicId string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPost,
		endpoint: "/v1/player/queue/playlist/last",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": playlistPublicId,
		},
	})

	return err
}
