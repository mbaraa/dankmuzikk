package apis

import (
	"dankmuzikk-web/config"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/services/requests"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		return
	}
	_ = requests.GetRequestAuthNoRespBody("/v1/logout", sessionToken.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   auth.SessionTokenKey,
		Value:  "",
		Path:   "/",
		Domain: config.Env().Hostname,
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
