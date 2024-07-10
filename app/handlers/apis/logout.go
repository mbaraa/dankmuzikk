package apis

import (
	"dankmuzikk/config"
	"dankmuzikk/handlers/middlewares/auth"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   auth.SessionTokenKey,
		Value:  "",
		Path:   "/",
		Domain: config.Env().Hostname,
		MaxAge: -1,
	})
	http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
}
