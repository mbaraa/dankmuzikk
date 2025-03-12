package main

import (
	"dankmuzikk/actions"
	"dankmuzikk/app"
	"dankmuzikk/config"
	"dankmuzikk/handlers/apis"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/handlers/middlewares/logger"
	"dankmuzikk/jwt"
	"dankmuzikk/log"
	"dankmuzikk/mailer"
	"dankmuzikk/mariadb"
	"dankmuzikk/youtube"
	"dankmuzikk/zip"
	"net/http"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
)

//go:generate templ generate -path ../../mailer/

func main() {
	err := StartServer()
	if err != nil {
		log.Fatalln(err)
	}
}

func StartServer() error {
	mariadbRepo, err := mariadb.New()
	if err != nil {
		return err
	}

	app := app.New(mariadbRepo)
	zipArchiver := zip.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	mailer := mailer.New()
	yt := youtube.New()

	usecases := actions.New(
		app,
		zipArchiver,
		jwtUtil,
		mailer,
		yt,
	)

	authMw := auth.New(usecases)

	minfyer := minify.New()
	minfyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minfyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)

	emailLoginApi := apis.NewEmailLoginApi(usecases)
	googleLoginApi := apis.NewGoogleLoginApi(usecases)
	searchApi := apis.NewYouTubeSearchApi(usecases)
	songApi := apis.NewSongsHandler(usecases)
	playlistsApi := apis.NewPlaylistApi(usecases)
	historyApi := apis.NewHistoryApi(usecases)
	userApi := apis.NewUserApi(usecases)

	v1ApisHandler := http.NewServeMux()
	v1ApisHandler.HandleFunc("POST /login/email", emailLoginApi.HandleEmailLogin)
	v1ApisHandler.HandleFunc("POST /signup/email", emailLoginApi.HandleEmailSignup)
	v1ApisHandler.HandleFunc("POST /verify-otp", emailLoginApi.HandleEmailOTPVerification)
	v1ApisHandler.HandleFunc("GET /login/google", googleLoginApi.HandleGoogleOAuthLogin)
	v1ApisHandler.HandleFunc("GET /signup/google", googleLoginApi.HandleGoogleOAuthLogin)
	v1ApisHandler.HandleFunc("POST /login/google/callback", googleLoginApi.HandleGoogleOAuthLoginCallback)
	v1ApisHandler.HandleFunc("GET /search/suggestions", searchApi.HandleSearchSuggestions)
	v1ApisHandler.HandleFunc("GET /search", searchApi.HandleSearchResults)
	v1ApisHandler.HandleFunc("GET /song/play", authMw.OptionalAuthApi(songApi.HandlePlaySong))
	v1ApisHandler.HandleFunc("GET /song/single", authMw.OptionalAuthApi(songApi.HandleGetSong))
	v1ApisHandler.HandleFunc("PUT /song/playlist", authMw.AuthApi(playlistsApi.HandleToggleSongInPlaylist))
	v1ApisHandler.HandleFunc("PUT /song/playlist/plays", authMw.AuthApi(songApi.HandleIncrementSongPlaysInPlaylist))
	v1ApisHandler.HandleFunc("PUT /song/playlist/upvote", authMw.AuthApi(songApi.HandleUpvoteSongPlaysInPlaylist))
	v1ApisHandler.HandleFunc("PUT /song/playlist/downvote", authMw.AuthApi(songApi.HandleDownvoteSongPlaysInPlaylist))
	v1ApisHandler.HandleFunc("GET /playlist/songs/mapped", authMw.AuthApi(playlistsApi.HandleGetPlaylistsForPopover))
	v1ApisHandler.HandleFunc("GET /playlist/all", authMw.AuthApi(playlistsApi.HandleGetPlaylists))
	v1ApisHandler.HandleFunc("GET /playlist", authMw.AuthApi(playlistsApi.HandleGetPlaylist))
	v1ApisHandler.HandleFunc("POST /playlist", authMw.AuthApi(playlistsApi.HandleCreatePlaylist))
	v1ApisHandler.HandleFunc("PUT /playlist/public", authMw.AuthApi(playlistsApi.HandleTogglePublicPlaylist))
	v1ApisHandler.HandleFunc("PUT /playlist/join", authMw.AuthApi(playlistsApi.HandleToggleJoinPlaylist))
	v1ApisHandler.HandleFunc("DELETE /playlist", authMw.AuthApi(playlistsApi.HandleDeletePlaylist))
	v1ApisHandler.HandleFunc("GET /playlist/zip", authMw.AuthApi(playlistsApi.HandleDonwnloadPlaylist))
	v1ApisHandler.HandleFunc("GET /history/{page}", authMw.AuthApi(historyApi.HandleGetMoreHistoryItems))
	v1ApisHandler.HandleFunc("GET /profile", userApi.HandleGetProfile)

	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/muzikkx/", http.StripPrefix("/muzikkx", http.FileServer(http.Dir(config.Env().YouTube.MuzikkDir))))

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/v1/", http.StripPrefix("/v1", v1ApisHandler))

	log.Info("Starting http server at port " + config.Env().Port)
	if config.Env().GoEnv == "dev" || config.Env().GoEnv == "beta" {
		return http.ListenAndServe(":"+config.Env().Port, logger.Handler(applicationHandler))
	}
	return http.ListenAndServe(":"+config.Env().Port, minfyer.Middleware(applicationHandler))
}
