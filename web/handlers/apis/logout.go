package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/config"
	"dankmuzikk-web/handlers/middlewares/auth"
	"net/http"
)

type logoutApi struct {
	usecases *actions.Actions
}

func NewLogoutApi(usecases *actions.Actions) *logoutApi {
	return &logoutApi{
		usecases: usecases,
	}
}

func (l *logoutApi) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	_ = l.usecases.Logout(sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:   auth.SessionTokenKey,
		Value:  "",
		Path:   "/",
		Domain: config.Env().Hostname,
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
