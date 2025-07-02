package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
	"strconv"
)

func (r *Requests) GetHistory(sessionToken string, pageIndex uint) ([]actions.Song, error) {
	return makeRequest[any, []actions.Song](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/history",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"page": strconv.Itoa(int(pageIndex)),
		},
	})
}
