package apis

import (
	"dankmuzikk/config"
	"dankmuzikk/log"
	"dankmuzikk/services/google"
	"net/http"
)

func HandleGoogleOAuthLogin(hand *http.ServeMux) {
	hand.HandleFunc("/api/login/google", func(w http.ResponseWriter, r *http.Request) {
		url := config.GoogleOAuthConfig().AuthCodeURL(google.CurrentRandomState())
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

func HandleGoogleOAuthLoginCallback(hand *http.ServeMux) {
	hand.HandleFunc("/api/login/google/callback", func(w http.ResponseWriter, r *http.Request) {
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

		err := google.CompleteLoginWithGoogle(state, code)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Errorln("[GOOGLE LOGIN API]: ", err)
			return
		}
	})
}
