package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
	"strconv"
)

func (r *Requests) GetFavorites(sessionToken string, pageIndex uint) (actions.GetFavoritesPayload, error) {
	return makeRequest[any, actions.GetFavoritesPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/library/favorite/songs",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"page": strconv.Itoa(int(pageIndex)),
		},
	})
}

func (r *Requests) AddSongToFavorites(sessionToken string, songPublicId string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodPost,
		endpoint: "/v1/library/favorite/song",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})

	return err
}

func (r *Requests) RemoveSongFromFavorites(sessionToken string, songPublicId string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodDelete,
		endpoint: "/v1/library/favorite/song",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})

	return err
}
