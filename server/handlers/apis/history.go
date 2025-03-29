package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/log"
	"encoding/json"
	"net/http"
	"strconv"
)

type historyApi struct {
	usecases *actions.Actions
}

func NewHistoryApi(usecases *actions.Actions) *historyApi {
	return &historyApi{
		usecases: usecases,
	}
}

func (h *historyApi) HandleGetMoreHistoryItems(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	page, err := strconv.Atoi(r.PathValue("page"))
	if err != nil || page <= 0 {
		handleErrorResponse(w, &ErrBadRequest{
			FieldName: "page",
		})
		return
	}

	params := actions.GetHistoryItemsParams{
		ActionContext: ctx,
		PageIndex:     uint(page),
	}

	recentPlays, err := h.usecases.GetHistoryItems(params)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(recentPlays)
}
