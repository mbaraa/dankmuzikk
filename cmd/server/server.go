package server

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/handlers"
	"dankmuzikk/handlers/apis"
	"dankmuzikk/handlers/pages"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
	"dankmuzikk/services/login"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/playlists/songs"
	"dankmuzikk/services/youtube/download"
	"dankmuzikk/services/youtube/search"
	"embed"
	"net/http"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
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
	playlistRepo := db.NewBaseDB[models.Playlist](dbConn)
	playlistOwnersRepo := db.NewBaseDB[models.PlaylistOwner](dbConn)
	playlistSongsRepo := db.NewBaseDB[models.PlaylistSong](dbConn)

	downloadService := download.New(songRepo)
	playlistsService := playlists.New(playlistRepo, playlistOwnersRepo, playlistSongsRepo)
	songsService := songs.New(playlistSongsRepo, playlistOwnersRepo, songRepo, playlistRepo, downloadService)

	jwtUtil := jwt.NewJWTImpl()

	gHandler := handlers.NewHandler(profileRepo, jwtUtil)

	///////////// Pages and files /////////////
	pagesHandler := http.NewServeMux()

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	pagesHandler.Handle("/static/", m.Middleware(http.FileServer(http.FS(staticFS))))
	pagesHandler.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		robotsFile, _ := staticFS.ReadFile("static/robots.txt")
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write(robotsFile)
	})
	pagesHandler.Handle("/muzikkx/", http.StripPrefix("/muzikkx", http.FileServer(http.Dir(config.Env().YouTube.MusicDir))))

	pagesRouter := pages.NewPagesHandler(profileRepo, playlistsService, jwtUtil, &search.ScraperSearch{}, downloadService)
	pagesHandler.HandleFunc("/", gHandler.OptionalAuthPage(pagesRouter.HandleHomePage))
	pagesHandler.HandleFunc("/signup", gHandler.AuthPage(pagesRouter.HandleSignupPage))
	pagesHandler.HandleFunc("/login", gHandler.AuthPage(pagesRouter.HandleLoginPage))
	pagesHandler.HandleFunc("/profile", gHandler.AuthPage(pagesRouter.HandleProfilePage))
	pagesHandler.HandleFunc("/about", gHandler.NoAuthPage(pagesRouter.HandleAboutPage))
	pagesHandler.HandleFunc("/playlists", gHandler.AuthPage(pagesRouter.HandlePlaylistsPage))
	pagesHandler.HandleFunc("/playlist/{playlist_id}", gHandler.AuthPage(pagesRouter.HandleSinglePlaylistPage))
	pagesHandler.HandleFunc("/privacy", gHandler.NoAuthPage(pagesRouter.HandlePrivacyPage))
	pagesHandler.HandleFunc("/search", gHandler.OptionalAuthPage(pagesRouter.HandleSearchResultsPage))

	///////////// APIs /////////////

	emailLoginApi := apis.NewEmailLoginApi(login.NewEmailLoginService(accountRepo, profileRepo, otpRepo, jwtUtil))
	googleLoginApi := apis.NewGoogleLoginApi(login.NewGoogleLoginService(accountRepo, profileRepo, otpRepo, jwtUtil))
	songDownloadApi := apis.NewDownloadHandler(downloadService, songsService)
	playlistsApi := apis.NewPlaylistApi(playlistsService, songsService)

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("POST /login/email", emailLoginApi.HandleEmailLogin)
	apisHandler.HandleFunc("POST /signup/email", emailLoginApi.HandleEmailSignup)
	apisHandler.HandleFunc("POST /verify-otp", emailLoginApi.HandleEmailOTPVerification)
	apisHandler.HandleFunc("GET /login/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("GET /signup/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("/login/google/callback", googleLoginApi.HandleGoogleOAuthLoginCallback)
	apisHandler.HandleFunc("GET /logout", apis.HandleLogout)
	apisHandler.HandleFunc("GET /search-suggestion", apis.HandleSearchSuggestions)
	apisHandler.HandleFunc("GET /song", songDownloadApi.HandlePlaySong)
	apisHandler.HandleFunc("POST /playlist", gHandler.AuthApi(playlistsApi.HandleCreatePlaylist))
	apisHandler.HandleFunc("PUT /toggle-song-in-playlist", gHandler.AuthApi(playlistsApi.HandleToggleSongInPlaylist))
	apisHandler.HandleFunc("PUT /increment-song-plays", gHandler.AuthApi(songDownloadApi.HandleIncrementSongPlaysInPlaylist))

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Info("Starting http server at port " + config.Env().Port)
	return http.ListenAndServe(":"+config.Env().Port, applicationHandler)
}
