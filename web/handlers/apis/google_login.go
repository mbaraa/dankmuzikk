package apis

import (
	"context"
	"dankmuzikk-web/config"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/services/requests"
	"dankmuzikk-web/views/components/status"
	"errors"
	"net/http"
	"time"
)

type googleLoginApi struct {
}

func NewGoogleLoginApi() *googleLoginApi {
	return &googleLoginApi{}
}

func (g *googleLoginApi) HandleGoogleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	resp, err := requests.GetRequest[map[string]string]("/v1/login/google")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := resp["redirect_url"]
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	resp, err := requests.PostRequest[map[string]string, map[string]string]("/v1/login/google/callback", map[string]string{
		"code":  code,
		"state": state,
	})
	if errors.Is(err, requests.ErrDifferentLoginMethodUsed) {
		log.Errorf("[GOOGLE LOGIN API]: Failed to login user:  error: %s\n", err.Error())
		status.
			BugsBunnyError("This account uses Email to login!").
			Render(context.Background(), w)
		return
	}
	if err != nil {
		log.Errorln("[GOOGLE LOGIN API]: ", err)
		status.
			GenericError("Account doesn't exist").
			Render(context.Background(), w)
		return
	}

	sessionToken := resp["session_token"]
	if sessionToken == "" {
		log.Errorln("[GOOGLE LOGIN API]: ", err)
		status.
			GenericError("Something went wrong").
			Render(context.Background(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionTokenKey,
		Value:    sessionToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
