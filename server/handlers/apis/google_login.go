package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/config"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/log"
	"net/http"
	"time"
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

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionTokenKey,
		Value:    payload.SessionToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   config.Env().DomainName,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	http.Redirect(w, r, config.WebUrl(), http.StatusTemporaryRedirect)
}
