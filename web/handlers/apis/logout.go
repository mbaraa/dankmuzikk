package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/config"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/views/components/status"
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
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}
	_ = l.usecases.Logout(ctx)

	http.SetCookie(w, &http.Cookie{
		Name:   auth.SessionTokenKey,
		Value:  "",
		Path:   "/",
		Domain: config.Env().Hostname,
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
