package server

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/handlers"
	"dankmuzikk/handlers/apis"
	"dankmuzikk/handlers/pages"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/history"
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
	historyRepo := db.NewBaseDB[models.History](dbConn)
	playlistVotersRepo := db.NewBaseDB[models.PlaylistSongVoter](dbConn)

	downloadService := download.New(songRepo)
	playlistsService := playlists.New(playlistRepo, playlistOwnersRepo, playlistSongsRepo)
	songsService := songs.New(playlistSongsRepo, playlistOwnersRepo, songRepo, playlistRepo, playlistVotersRepo, downloadService)
	historyService := history.New(historyRepo, songRepo)

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

	pagesRouter := pages.NewPagesHandler(profileRepo, playlistsService, jwtUtil, &search.ScraperSearch{}, downloadService, historyService, songsService)
	pagesHandler.HandleFunc("/", gHandler.OptionalAuthPage(pagesRouter.HandleHomePage))
	pagesHandler.HandleFunc("GET /signup", gHandler.AuthPage(pagesRouter.HandleSignupPage))
	pagesHandler.HandleFunc("GET /login", gHandler.AuthPage(pagesRouter.HandleLoginPage))
	pagesHandler.HandleFunc("GET /profile", gHandler.AuthPage(pagesRouter.HandleProfilePage))
	pagesHandler.HandleFunc("GET /about", gHandler.NoAuthPage(pagesRouter.HandleAboutPage))
	pagesHandler.HandleFunc("GET /playlists", gHandler.AuthPage(pagesRouter.HandlePlaylistsPage))
	pagesHandler.HandleFunc("GET /playlist/{playlist_id}", gHandler.AuthPage(pagesRouter.HandleSinglePlaylistPage))
	pagesHandler.HandleFunc("GET /song/{song_id}", gHandler.OptionalAuthPage(pagesRouter.HandleSingleSongPage))
	pagesHandler.HandleFunc("GET /privacy", gHandler.NoAuthPage(pagesRouter.HandlePrivacyPage))
	pagesHandler.HandleFunc("GET /search", gHandler.OptionalAuthPage(pagesRouter.HandleSearchResultsPage))

	///////////// APIs /////////////

	emailLoginApi := apis.NewEmailLoginApi(login.NewEmailLoginService(accountRepo, profileRepo, otpRepo, jwtUtil))
	googleLoginApi := apis.NewGoogleLoginApi(login.NewGoogleLoginService(accountRepo, profileRepo, otpRepo, jwtUtil))
	songApi := apis.NewDownloadHandler(downloadService, songsService, historyService)
	playlistsApi := apis.NewPlaylistApi(playlistsService, songsService)
	historyApi := apis.NewHistoryApi(historyService)

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("POST /login/email", emailLoginApi.HandleEmailLogin)
	apisHandler.HandleFunc("POST /signup/email", emailLoginApi.HandleEmailSignup)
	apisHandler.HandleFunc("POST /verify-otp", emailLoginApi.HandleEmailOTPVerification)
	apisHandler.HandleFunc("GET /login/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("GET /signup/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("/login/google/callback", googleLoginApi.HandleGoogleOAuthLoginCallback)
	apisHandler.HandleFunc("GET /logout", apis.HandleLogout)
	apisHandler.HandleFunc("GET /search-suggestion", apis.HandleSearchSuggestions)
	apisHandler.HandleFunc("GET /song", gHandler.OptionalAuthApi(songApi.HandlePlaySong))
	apisHandler.HandleFunc("GET /song/single", gHandler.OptionalAuthApi(songApi.HandleGetSong))
	apisHandler.HandleFunc("PUT /song/playlist", gHandler.AuthApi(playlistsApi.HandleToggleSongInPlaylist))
	apisHandler.HandleFunc("PUT /song/playlist/plays", gHandler.AuthApi(songApi.HandleIncrementSongPlaysInPlaylist))
	apisHandler.HandleFunc("PUT /song/playlist/upvote", gHandler.AuthApi(songApi.HandleUpvoteSongPlaysInPlaylist))
	apisHandler.HandleFunc("PUT /song/playlist/downvote", gHandler.AuthApi(songApi.HandleDownvoteSongPlaysInPlaylist))
	apisHandler.HandleFunc("GET /playlist/all", gHandler.AuthApi(playlistsApi.HandleGetPlaylistsForPopover))
	apisHandler.HandleFunc("GET /playlist", gHandler.AuthApi(playlistsApi.HandleGetPlaylist))
	apisHandler.HandleFunc("POST /playlist", gHandler.AuthApi(playlistsApi.HandleCreatePlaylist))
	apisHandler.HandleFunc("PUT /playlist/public", gHandler.AuthApi(playlistsApi.HandleTogglePublicPlaylist))
	apisHandler.HandleFunc("PUT /playlist/join", gHandler.AuthApi(playlistsApi.HandleToggleJoinPlaylist))
	apisHandler.HandleFunc("DELETE /playlist", gHandler.AuthApi(playlistsApi.HandleDeletePlaylist))
	apisHandler.HandleFunc("GET /history/{page}", gHandler.AuthApi(historyApi.HandleGetMoreHistoryItems))

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Info("Starting http server at port " + config.Env().Port)
	return http.ListenAndServe(":"+config.Env().Port, applicationHandler)
}
