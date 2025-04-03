package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/config"
	dankerrors "dankmuzikk-web/errors"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/status"
	"errors"
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
	payload, err := g.usecases.LoginUsingGoogle()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, payload.RedirectUrl, http.StatusTemporaryRedirect)
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

	payload, err := g.usecases.FinishLoginUsingGoogle(actions.FinishLoginUsingGoogleParams{
		Code:  code,
		State: state,
	})
	if errors.Is(err, dankerrors.ErrDifferentLoginMethodUsed) {
		log.Errorf("[GOOGLE LOGIN API]: Failed to login user:  error: %s\n", err.Error())
		status.
			BugsBunnyError("This account uses Email to login!").
			Render(r.Context(), w)
		return
	}
	if err != nil {
		log.Errorln("[GOOGLE LOGIN API]: ", err)
		status.
			GenericError("Account doesn't exist").
			Render(r.Context(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionTokenKey,
		Value:    payload.SessionToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
