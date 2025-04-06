package main

import (
	"dankmuzikk/actions"
	"dankmuzikk/app"
	"dankmuzikk/blobs"
	"dankmuzikk/config"
	"dankmuzikk/evy"
	"dankmuzikk/genius"
	"dankmuzikk/handlers/apis"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/handlers/middlewares/contenttype"
	"dankmuzikk/handlers/middlewares/logger"
	"dankmuzikk/jwt"
	"dankmuzikk/log"
	"dankmuzikk/mailer"
	"dankmuzikk/mariadb"
	"dankmuzikk/redis"
	"dankmuzikk/youtube"
	"dankmuzikk/zip"
	"net/http"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
)

func main() {
	mariadbRepo, err := mariadb.New()
	if err != nil {
		log.Fatalln(err)
	}

	cache := redis.New()

	app := app.New(mariadbRepo, cache)
	eventhub := evy.New()
	zipArchiver := zip.New()
	blobstorage := blobs.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	mailer := mailer.New()
	yt := youtube.New()
	lyrics := genius.New()

	usecases := actions.New(
		app,
		cache,
		eventhub,
		zipArchiver,
		blobstorage,
		jwtUtil,
		mailer,
		yt,
		lyrics,
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
	accountApi := apis.NewAccountApi(usecases)
	libraryApi := apis.NewLibraryApi(usecases)

	v1ApisHandler := http.NewServeMux()
	v1ApisHandler.HandleFunc("POST /login/email", emailLoginApi.HandleEmailLogin)
	v1ApisHandler.HandleFunc("POST /signup/email", emailLoginApi.HandleEmailSignup)
	v1ApisHandler.HandleFunc("POST /verify-otp", emailLoginApi.HandleEmailOTPVerification)
	v1ApisHandler.HandleFunc("GET /login/google", googleLoginApi.HandleGoogleOAuthLogin)
	v1ApisHandler.HandleFunc("POST /login/google/callback", googleLoginApi.HandleGoogleOAuthLoginCallback)
	v1ApisHandler.HandleFunc("GET /search/suggestions", searchApi.HandleSearchSuggestions)
	v1ApisHandler.HandleFunc("GET /search", searchApi.HandleSearchResults)

	v1ApisHandler.HandleFunc("GET /song/play", authMw.OptionalAuthApi(songApi.HandlePlaySong))
	v1ApisHandler.HandleFunc("GET /song/single", authMw.OptionalAuthApi(songApi.HandleGetSong))
	v1ApisHandler.HandleFunc("PUT /song/playlist", authMw.AuthApi(songApi.HandleToggleSongInPlaylist))
	v1ApisHandler.HandleFunc("PUT /song/playlist/upvote", authMw.AuthApi(songApi.HandleUpvoteSongPlaysInPlaylist))
	v1ApisHandler.HandleFunc("PUT /song/playlist/downvote", authMw.AuthApi(songApi.HandleDownvoteSongPlaysInPlaylist))
	v1ApisHandler.HandleFunc("GET /song/lyrics", songApi.HandleGetSongLyrics)

	v1ApisHandler.HandleFunc("GET /playlist", authMw.AuthApi(playlistsApi.HandleGetPlaylist))
	v1ApisHandler.HandleFunc("POST /playlist", authMw.AuthApi(playlistsApi.HandleCreatePlaylist))
	v1ApisHandler.HandleFunc("DELETE /playlist", authMw.AuthApi(playlistsApi.HandleDeletePlaylist))
	v1ApisHandler.HandleFunc("GET /playlist/songs/mapped", authMw.AuthApi(playlistsApi.HandleGetPlaylistsForPopover))
	v1ApisHandler.HandleFunc("GET /playlist/all", authMw.AuthApi(playlistsApi.HandleGetPlaylists))
	v1ApisHandler.HandleFunc("PUT /playlist/public", authMw.AuthApi(playlistsApi.HandleTogglePublicPlaylist))
	v1ApisHandler.HandleFunc("PUT /playlist/join", authMw.AuthApi(playlistsApi.HandleToggleJoinPlaylist))
	v1ApisHandler.HandleFunc("GET /playlist/zip", authMw.AuthApi(playlistsApi.HandleDonwnloadPlaylist))

	v1ApisHandler.HandleFunc("GET /history", authMw.AuthApi(historyApi.HandleGetHistoryItems))

	v1ApisHandler.HandleFunc("POST /library/favorite/song", authMw.AuthApi(libraryApi.HandleAddSongToFavorites))
	v1ApisHandler.HandleFunc("DELETE /library/favorite/song", authMw.AuthApi(libraryApi.HandleRemoveSongFromFavorites))
	v1ApisHandler.HandleFunc("GET /library/favorite/songs", authMw.AuthApi(libraryApi.HandleGetFavoriteSongs))

	v1ApisHandler.HandleFunc("GET /me/profile", authMw.AuthApi(accountApi.HandleGetProfile))
	v1ApisHandler.HandleFunc("GET /me/auth", authMw.AuthApi(accountApi.HandleAuthCheck))
	v1ApisHandler.HandleFunc("GET /me/logout", func(w http.ResponseWriter, r *http.Request) {
		sessionToken, ok := r.Header["Authorization"]
		if !ok {
			return
		}
		_ = cache.InvalidateAuthenticatedAccount(sessionToken[0])
	})

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/v1/", http.StripPrefix("/v1", contenttype.Json(v1ApisHandler)))

	log.Info("Starting http server at port " + config.Env().Port)
	if config.Env().GoEnv == "dev" || config.Env().GoEnv == "beta" {
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, logger.Handler(applicationHandler)))
	}

	log.Fatalln(http.ListenAndServe(":"+config.Env().Port, minfyer.Middleware(applicationHandler)))
}
