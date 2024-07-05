package server

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/handlers/apis"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/handlers/middlewares/contenttype"
	"dankmuzikk/handlers/middlewares/ismobile"
	"dankmuzikk/handlers/middlewares/theme"
	"dankmuzikk/handlers/pages"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/archive"
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

	zipService := archive.NewService()
	downloadService := download.New(songRepo)
	playlistsService := playlists.New(playlistRepo, playlistOwnersRepo, playlistSongsRepo, zipService)
	songsService := songs.New(playlistSongsRepo, playlistOwnersRepo, songRepo, playlistRepo, playlistVotersRepo, downloadService)
	historyService := history.New(historyRepo, songRepo)

	jwtUtil := jwt.NewJWTImpl()

	authMw := auth.New(profileRepo, jwtUtil)

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
	pagesHandler.HandleFunc("/", contenttype.Html(authMw.OptionalAuthPage(pagesRouter.HandleHomePage)))
	pagesHandler.HandleFunc("GET /signup", contenttype.Html(authMw.AuthPage(pagesRouter.HandleSignupPage)))
	pagesHandler.HandleFunc("GET /login", contenttype.Html(authMw.AuthPage(pagesRouter.HandleLoginPage)))
	pagesHandler.HandleFunc("GET /profile", contenttype.Html(authMw.AuthPage(pagesRouter.HandleProfilePage)))
	pagesHandler.HandleFunc("GET /about", contenttype.Html(pagesRouter.HandleAboutPage))
	pagesHandler.HandleFunc("GET /playlists", contenttype.Html(authMw.AuthPage(pagesRouter.HandlePlaylistsPage)))
	pagesHandler.HandleFunc("GET /playlist/{playlist_id}", contenttype.Html(authMw.AuthPage(pagesRouter.HandleSinglePlaylistPage)))
	pagesHandler.HandleFunc("GET /song/{song_id}", contenttype.Html(authMw.OptionalAuthPage(pagesRouter.HandleSingleSongPage)))
	pagesHandler.HandleFunc("GET /privacy", contenttype.Html(pagesRouter.HandlePrivacyPage))
	pagesHandler.HandleFunc("GET /search", contenttype.Html(authMw.OptionalAuthPage(pagesRouter.HandleSearchResultsPage)))

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
	apisHandler.HandleFunc("GET /song", authMw.OptionalAuthApi(songApi.HandlePlaySong))
	apisHandler.HandleFunc("GET /song/single", authMw.OptionalAuthApi(songApi.HandleGetSong))
	apisHandler.HandleFunc("PUT /song/playlist", authMw.AuthApi(playlistsApi.HandleToggleSongInPlaylist))
	apisHandler.HandleFunc("PUT /song/playlist/plays", authMw.AuthApi(songApi.HandleIncrementSongPlaysInPlaylist))
	apisHandler.HandleFunc("PUT /song/playlist/upvote", authMw.AuthApi(songApi.HandleUpvoteSongPlaysInPlaylist))
	apisHandler.HandleFunc("PUT /song/playlist/downvote", authMw.AuthApi(songApi.HandleDownvoteSongPlaysInPlaylist))
	apisHandler.HandleFunc("GET /playlist/all", authMw.AuthApi(playlistsApi.HandleGetPlaylistsForPopover))
	apisHandler.HandleFunc("GET /playlist", authMw.AuthApi(playlistsApi.HandleGetPlaylist))
	apisHandler.HandleFunc("POST /playlist", authMw.AuthApi(playlistsApi.HandleCreatePlaylist))
	apisHandler.HandleFunc("PUT /playlist/public", authMw.AuthApi(playlistsApi.HandleTogglePublicPlaylist))
	apisHandler.HandleFunc("PUT /playlist/join", authMw.AuthApi(playlistsApi.HandleToggleJoinPlaylist))
	apisHandler.HandleFunc("DELETE /playlist", authMw.AuthApi(playlistsApi.HandleDeletePlaylist))
	apisHandler.HandleFunc("GET /playlist/zip", authMw.AuthApi(playlistsApi.HandleDonwnloadPlaylist))
	apisHandler.HandleFunc("GET /history/{page}", authMw.AuthApi(historyApi.HandleGetMoreHistoryItems))

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", ismobile.Handler(theme.Handler(pagesHandler)))
	applicationHandler.Handle("/api/", ismobile.Handler(theme.Handler(http.StripPrefix("/api", apisHandler))))

	log.Info("Starting http server at port " + config.Env().Port)
	return http.ListenAndServe(":"+config.Env().Port, ismobile.Handler(theme.Handler(m.Middleware(applicationHandler))))
}
