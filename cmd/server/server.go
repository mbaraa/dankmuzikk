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
	"dankmuzikk/services/youtube"
	"embed"
	"net/http"
)

func StartServer(staticFS embed.FS) error {
	///////////// Pages and files /////////////
	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/static/", http.FileServer(http.FS(staticFS)))
	pagesHandler.Handle("/music/", http.StripPrefix("/music", http.FileServer(http.Dir(config.Env().YouTube.MusicDir))))

	jwtUtil := jwt.NewJWTImpl()

	pagesHandler.HandleFunc("/", pages.Handler(pages.HandleHomePage))
	pagesHandler.HandleFunc("/signup", pages.AuthHandler(pages.HandleSignupPage, jwtUtil))
	pagesHandler.HandleFunc("/login", pages.AuthHandler(pages.HandleLoginPage, jwtUtil))
	pagesHandler.HandleFunc("/profile", pages.AuthHandler(pages.HandleProfilePage, jwtUtil))
	pagesHandler.HandleFunc("/about", pages.Handler(pages.HandleAboutPage))
	pagesHandler.HandleFunc("/playlists", pages.AuthHandler(pages.HandlePlaylistsPage, jwtUtil))
	pagesHandler.HandleFunc("/privacy", pages.Handler(pages.HandlePrivacyPage))
	pagesHandler.HandleFunc("/search", pages.Handler(pages.HandleSearchResultsPage(&youtube.YouTubeScraperSearch{})))

	///////////// APIs /////////////
	dbConn, err := db.Connector()
	if err != nil {
		log.Fatalln(log.ErrorLevel, err)
	}

	accountRepo := db.NewBaseDB[models.Account](dbConn)
	profileRepo := db.NewBaseDB[models.Profile](dbConn)
	otpRepo := db.NewBaseDB[models.EmailVerificationCode](dbConn)

	emailLoginApi := apis.NewEmailLoginApi(login.NewEmailLoginService(accountRepo, profileRepo, otpRepo, jwtUtil))
	googleLoginApi := apis.NewGoogleLoginApi(login.NewGoogleLoginService(accountRepo, profileRepo, otpRepo, jwtUtil))

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("POST /login/email", emailLoginApi.HandleEmailLogin)
	apisHandler.HandleFunc("POST /signup/email", emailLoginApi.HandleEmailSignup)
	apisHandler.HandleFunc("POST /verify-otp", emailLoginApi.HandleEmailOTPVerification)
	apisHandler.HandleFunc("GET /login/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("/login/google/callback", googleLoginApi.HandleGoogleOAuthLoginCallback)
	apisHandler.HandleFunc("GET /search-suggession", apis.HandleSearchSugessions)
	apisHandler.HandleFunc("GET /song/download/{youtube_video_id}", apis.HandleDownloadSong)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Info("Starting http server at port " + config.Env().Port)
	return http.ListenAndServe(":"+config.Env().Port, applicationHandler)
}
