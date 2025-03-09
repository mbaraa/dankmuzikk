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

	payload, err := g.usecases.LoginWithGoogle(actions.LoginWithGoogleParams{
		Code:  code,
		State: state,
	})
	if err != nil {
		log.Errorln("[GOOGLE LOGIN API]: ", err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)

	//	// TODO: idk something
	//	http.SetCookie(w, &http.Cookie{
	//		Name:     auth.SessionTokenKey,
	//		Value:    sessionToken,
	//		HttpOnly: true,
	//		Path:     "/",
	//		Domain:   config.Env().Hostname,
	//		Expires:  time.Now().UTC().Add(time.Hour * 24 * 30),
	//	})
	//	http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
}
