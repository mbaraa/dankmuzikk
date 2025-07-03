package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
)

func (r *Requests) GetSongMetadata(sessionToken, clientHash, songPublicId string) (actions.Song, error) {
	return makeRequest[any, actions.Song](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/song/single",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})
}

func (r *Requests) PlaySong(sessionToken, clientHash, songPublicId, playlistPublicId string) (string, error) {
	resp, err := makeRequest[any, map[string]string](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/song/play",
		headers: map[string]string{
			"Authorization": sessionToken,
			"X-Client-Hash": clientHash,
		},
		queryParams: map[string]string{
			"id":          songPublicId,
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return "", err
	}

	return resp["media_url"], nil
}
