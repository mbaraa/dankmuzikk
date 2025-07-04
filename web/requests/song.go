package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
)

func (r *Requests) GetSongMetadata(sessionToken, clientHash, songPublicId string) (actions.Song, error) {
	return makeRequest[any, actions.Song](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/song",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})
}

func (r *Requests) PlaySong(sessionToken, clientHash, songPublicId string) (actions.Song, error) {
	return makeRequest[any, actions.Song](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/song/play",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})
}

func (r *Requests) PlaySongFromPlaylist(sessionToken, clientHash, songPublicId, playlistPublicId string) (actions.Song, error) {
	return makeRequest[any, actions.Song](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/song/play/playlist",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id":          songPublicId,
			"playlist-id": playlistPublicId,
		},
	})
}

func (r *Requests) PlaySongFromFavorites(sessionToken, clientHash, songPublicId string) (actions.Song, error) {
	return makeRequest[any, actions.Song](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/song/play/favorites",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})
}

func (r *Requests) PlaySongFromQueue(sessionToken, clientHash, songPublicId string) (actions.Song, error) {
	return makeRequest[any, actions.Song](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/song/play/queue",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})
}
