package apis

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/login"
	"dankmuzikk/views/components/status"
	"errors"
	"net/http"
	"time"
)

type googleLoginApi struct {
	service *login.GoogleLoginService
}

func NewGoogleLoginApi(service *login.GoogleLoginService) *googleLoginApi {
	return &googleLoginApi{service}
}

func (g *googleLoginApi) HandleGoogleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOAuthConfig().AuthCodeURL(g.service.CurrentRandomState())
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

	sessionToken, err := g.service.Login(state, code)
	if err != nil && errors.Is(err, login.ErrDifferentLoginMethod) {
		log.Errorf("[EMAIL LOGIN API]: Failed to login, error: %s\n", err.Error())
		// w.WriteHeader(http.StatusInternalServerError)
		status.
			BugsBunnyError("This account uses Email to login!").
			Render(context.Background(), w)
		return
	} else if err != nil {
		// w.WriteHeader(http.StatusUnauthorized)
		status.
			GenericError("Account doesn't exist").
			Render(context.Background(), w)
		log.Errorln("[GOOGLE LOGIN API]: ", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     handlers.SessionTokenKey,
		Value:    sessionToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 30),
	})
	http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
}
