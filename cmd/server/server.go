package server

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/handlers/apis"
	"dankmuzikk/handlers/pages"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
	"dankmuzikk/services/login"
	"dankmuzikk/services/youtube/download"
	"dankmuzikk/services/youtube/search"
	"embed"
	"net/http"
)

func StartServer(staticFS embed.FS) error {
	dbConn, err := db.Connector()
	if err != nil {
		return err
	}

	accountRepo := db.NewBaseDB[models.Account](dbConn)
	profileRepo := db.NewBaseDB[models.Profile](dbConn)
	otpRepo := db.NewBaseDB[models.EmailVerificationCode](dbConn)
	songRepo := db.NewBaseDB[models.Song](dbConn)
	// playlistRepo := db.NewBaseDB[models.Playlist](dbConn)

	jwtUtil := jwt.NewJWTImpl()

	///////////// Pages and files /////////////
	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/static/", http.FileServer(http.FS(staticFS)))
	pagesHandler.Handle("/music/", http.StripPrefix("/music", http.FileServer(http.Dir(config.Env().YouTube.MusicDir))))

	pagesRouter := pages.NewPagesHandler(profileRepo, jwtUtil)
	pagesHandler.HandleFunc("/", pagesRouter.Handler(pagesRouter.HandleHomePage))
	pagesHandler.HandleFunc("/signup", pagesRouter.AuthHandler(pagesRouter.HandleSignupPage))
	pagesHandler.HandleFunc("/login", pagesRouter.AuthHandler(pagesRouter.HandleLoginPage))
	pagesHandler.HandleFunc("/profile", pagesRouter.AuthHandler(pagesRouter.HandleProfilePage))
	pagesHandler.HandleFunc("/about", pagesRouter.Handler(pagesRouter.HandleAboutPage))
	pagesHandler.HandleFunc("/playlists", pagesRouter.AuthHandler(pagesRouter.HandlePlaylistsPage))
	pagesHandler.HandleFunc("/privacy", pagesRouter.Handler(pagesRouter.HandlePrivacyPage))
	pagesHandler.HandleFunc("/search", pagesRouter.Handler(pagesRouter.HandleSearchResultsPage(&search.ScraperSearch{})))

	///////////// APIs /////////////

	emailLoginApi := apis.NewEmailLoginApi(login.NewEmailLoginService(accountRepo, profileRepo, otpRepo, jwtUtil))
	googleLoginApi := apis.NewGoogleLoginApi(login.NewGoogleLoginService(accountRepo, profileRepo, otpRepo, jwtUtil))
	songDownloadApi := apis.NewDownloadHandler(*download.New(songRepo))

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("POST /login/email", emailLoginApi.HandleEmailLogin)
	apisHandler.HandleFunc("POST /signup/email", emailLoginApi.HandleEmailSignup)
	apisHandler.HandleFunc("POST /verify-otp", emailLoginApi.HandleEmailOTPVerification)
	apisHandler.HandleFunc("GET /login/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("GET /signup/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("/login/google/callback", googleLoginApi.HandleGoogleOAuthLoginCallback)
	apisHandler.HandleFunc("GET /logout", apis.HandleLogout)
	apisHandler.HandleFunc("GET /search-suggestion", apis.HandleSearchSuggestions)
	apisHandler.HandleFunc("GET /song/download", songDownloadApi.HandleDownloadSong)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Info("Starting http server at port " + config.Env().Port)
	return http.ListenAndServe(":"+config.Env().Port, applicationHandler)
}
