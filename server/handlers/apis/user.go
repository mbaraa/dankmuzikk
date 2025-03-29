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

func NewUserApi(usecases *actions.Actions) *userApi {
	return &userApi{
		usecases: usecases,
	}
}

func (u *userApi) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Header["Authorization"]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	profile, err := u.usecases.GetProfile(sessionToken[0])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(profile)
}
