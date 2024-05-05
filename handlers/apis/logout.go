package apis

import (
	"dankmuzikk/config"
	"dankmuzikk/handlers"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   handlers.SessionTokenKey,
		Value:  "",
		Path:   "/",
		Domain: config.Env().Hostname,
		MaxAge: -1,
	})
	http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
}
