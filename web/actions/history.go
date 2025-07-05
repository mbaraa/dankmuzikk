package actions

type GetHistoryParams struct {
	ActionContext
	PageIndex uint
}

func (a *Actions) GetHistory(params GetHistoryParams) ([]Song, error) {
	return a.requests.GetHistory(params.SessionToken, params.PageIndex)
}
