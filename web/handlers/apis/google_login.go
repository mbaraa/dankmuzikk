package apis

import (
	"dankmuzikk-web/config"
	"dankmuzikk-web/log"
	"net/http"
)

type googleLoginApi struct {
}

func NewGoogleLoginApi() *googleLoginApi {
	return &googleLoginApi{}
}

func (g *googleLoginApi) HandleGoogleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	url := config.GetRequestUrl("/v1/login/google") // config.GoogleOAuthConfig().AuthCodeURL(g.service.CurrentRandomState())
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (g *googleLoginApi) HandleGoogleOAuthLoginCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Warningln("[GOOGLE LOGIN API]: Failed to login with Google due to empty state")
		return
	}
	code := r.FormValue("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Warningln("[GOOGLE LOGIN API]: Failed to login with Google due to empty code")
		return
	}

	url := config.GetRequestUrl("/v1/login/google/callback") // config.GoogleOAuthConfig().AuthCodeURL(g.service.CurrentRandomState())
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	return
}
