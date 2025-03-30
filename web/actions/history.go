package actions

type GetHistoryParams struct {
	RequestContext `json:"-"`
	PageIndex      uint `json:"page_index"`
}

func (a *Actions) GetHistory(sessionToken string, pageIndex uint) ([]Song, error) {
	return a.requests.GetHistory(sessionToken, pageIndex)
}
