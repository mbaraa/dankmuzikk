package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
)

func (r *Requests) GetSongLyrics(songPublicId string) (actions.GetLyricsForSongPayload, error) {
	return makeRequest[any, actions.GetLyricsForSongPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/song/lyrics",
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})
}
