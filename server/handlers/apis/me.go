package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/log"
	"encoding/json"
	"net/http"
)

type meApi struct {
	usecases *actions.Actions
}

func NewMeApi(usecases *actions.Actions) *meApi {
	return &meApi{
		usecases: usecases,
	}
}

func (u *meApi) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	profile, err := u.usecases.GetProfile(ctx.Account.Email)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(profile)
}

func (u *meApi) HandleAuthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
}

func (m *meApi) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Header["Authorization"]
	if !ok {
		return
	}
	_ = m.usecases.InvalidateAuthenticatedAccount(sessionToken[0])
}
