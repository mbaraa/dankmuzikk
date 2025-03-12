package history

import (
	"dankmuzikk-web/entities"
	"dankmuzikk-web/services/requests"
	"strconv"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (h *Service) Get(authToken string, page uint) ([]entities.Song, error) {
	return requests.GetRequestAuth[[]entities.Song]("/v1/history/"+strconv.Itoa(int(page)), authToken)
}
