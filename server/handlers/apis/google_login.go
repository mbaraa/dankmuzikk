package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/config"
	"dankmuzikk/log"
	"encoding/json"
	"net/http"
)

type googleLoginApi struct {
	usecases *actions.Actions
}

func NewGoogleLoginApi(usecases *actions.Actions) *googleLoginApi {
	return &googleLoginApi{
		usecases: usecases,
	}
}

func (g *googleLoginApi) HandleGoogleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOAuthConfig().AuthCodeURL(g.usecases.CurrentRandomState())

	_ = json.NewEncoder(w).Encode(map[string]string{
		"redirect_url": url,
	})
}

func (g *googleLoginApi) HandleGoogleOAuthLoginCallback(w http.ResponseWriter, r *http.Request) {
	var params actions.LoginWithGoogleParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warningln(err)
		return
	}

	payload, err := g.usecases.LoginWithGoogle(params)
	if err != nil {
		log.Errorln("[GOOGLE LOGIN API]: ", err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
