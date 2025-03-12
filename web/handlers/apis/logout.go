package apis

import (
	"dankmuzikk-web/config"
	"dankmuzikk-web/handlers/middlewares/auth"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   auth.SessionTokenKey,
		Value:  "",
		Path:   "/",
		Domain: config.Env().DomainName,
		MaxAge: -1,
	})
	http.Redirect(w, r, config.Env().DomainName, http.StatusTemporaryRedirect)
}
