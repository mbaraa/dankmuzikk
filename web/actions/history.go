package actions

import (
	"dankmuzikk-web/requests"
	"net/http"
	"strconv"
)

type GetHistoryParams struct {
	ActionContext
	PageIndex uint
}

func (a *Actions) GetHistory(params GetHistoryParams) ([]Song, error) {
	return requests.Do[any, []Song](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/history",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"page": strconv.Itoa(int(params.PageIndex)),
		},
	})
}
