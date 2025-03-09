package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/handlers/middlewares/auth"
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
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	page, err := strconv.Atoi(r.PathValue("page"))
	if err != nil {
		page = 1
	}
	if page <= 0 {
		page *= -1
	}
	recentPlays, err := h.usecases.GetHistoryItems(profileId, uint(page))
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(recentPlays)
}
