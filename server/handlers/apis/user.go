package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/log"
	"encoding/json"
	"net/http"
)

type userApi struct {
	usecases *actions.Actions
}

func NewAccountApi(usecases *actions.Actions) *userApi {
	return &userApi{
		usecases: usecases,
	}
}

func (u *userApi) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
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

func (u *userApi) HandleAuthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
}
