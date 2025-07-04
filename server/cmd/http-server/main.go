package main

import (
	"dankmuzikk/actions"
	"dankmuzikk/app"
	"dankmuzikk/blobs"
	"dankmuzikk/config"
	"dankmuzikk/danklyrics"
	"dankmuzikk/evy"
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

var (
	minifyer       *minify.M
	usecases       *actions.Actions
	authMiddleware *auth.Middleware
)

func init() {
	mariadbRepo, err := mariadb.New()
	if err != nil {
		log.Fatalln(err)
	}

	cache := redis.New()
	playerCache := redis.NewPlayerCache()

	app := app.New(mariadbRepo, cache, playerCache)
	eventhub := evy.New()
	zipArchiver := zip.New()
	blobstorage := blobs.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	mailer := mailer.New()
	yt := youtube.New()
	lyrics := danklyrics.New()

	usecases = actions.New(
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

	authMiddleware = auth.New(usecases)

	minifyer = minify.New()
	minifyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
}

func main() {
	emailLoginApi := apis.NewEmailLoginApi(usecases)
	googleLoginApi := apis.NewGoogleLoginApi(usecases)
	searchApi := apis.NewYouTubeSearchApi(usecases)
	songApi := apis.NewSongsHandler(usecases)
	playlistsApi := apis.NewPlaylistApi(usecases)
	historyApi := apis.NewHistoryApi(usecases)
	accountApi := apis.NewMeApi(usecases)
	libraryApi := apis.NewLibraryApi(usecases)
	playerStateApi := apis.NewPlayerStateApi(usecases)

	v1ApisHandler := http.NewServeMux()
	v1ApisHandler.HandleFunc("POST /login/email", emailLoginApi.HandleEmailLogin)
	v1ApisHandler.HandleFunc("POST /signup/email", emailLoginApi.HandleEmailSignup)
	v1ApisHandler.HandleFunc("POST /verify-otp", emailLoginApi.HandleEmailOTPVerification)
	v1ApisHandler.HandleFunc("GET /login/google", googleLoginApi.HandleGoogleOAuthLogin)
	v1ApisHandler.HandleFunc("POST /login/google/callback", googleLoginApi.HandleGoogleOAuthLoginCallback)
	v1ApisHandler.HandleFunc("GET /search/suggestions", searchApi.HandleSearchSuggestions)
	v1ApisHandler.HandleFunc("GET /search", searchApi.HandleSearchResults)

	v1ApisHandler.HandleFunc("GET /song", authMiddleware.OptionalAuthApi(songApi.HandleGetSong))
	v1ApisHandler.HandleFunc("PUT /song/play", authMiddleware.OptionalAuthApi(songApi.HandlePlaySong))
	v1ApisHandler.HandleFunc("PUT /song/play/playlist", authMiddleware.OptionalAuthApi(songApi.HandlePlaySongFromPlaylist))
	v1ApisHandler.HandleFunc("PUT /song/play/favorite", authMiddleware.OptionalAuthApi(songApi.HandlePlaySongFromFavorites))
	v1ApisHandler.HandleFunc("PUT /song/play/queue", authMiddleware.OptionalAuthApi(songApi.HandlePlaySongFromQueue))
	v1ApisHandler.HandleFunc("GET /song/lyrics", songApi.HandleGetSongLyrics)

	v1ApisHandler.HandleFunc("GET /playlist", authMiddleware.AuthApi(playlistsApi.HandleGetPlaylist))
	v1ApisHandler.HandleFunc("POST /playlist", authMiddleware.AuthApi(playlistsApi.HandleCreatePlaylist))
	v1ApisHandler.HandleFunc("DELETE /playlist", authMiddleware.AuthApi(playlistsApi.HandleDeletePlaylist))
	v1ApisHandler.HandleFunc("PUT /playlist/song", authMiddleware.AuthApi(songApi.HandleToggleSongInPlaylist))
	v1ApisHandler.HandleFunc("PUT /playlist/song/upvote", authMiddleware.AuthApi(songApi.HandleUpvoteSongPlaysInPlaylist))
	v1ApisHandler.HandleFunc("PUT /playlist/song/downvote", authMiddleware.AuthApi(songApi.HandleDownvoteSongPlaysInPlaylist))
	v1ApisHandler.HandleFunc("GET /playlist/songs/mapped", authMiddleware.AuthApi(playlistsApi.HandleGetPlaylistsForPopover))
	v1ApisHandler.HandleFunc("GET /playlist/all", authMiddleware.AuthApi(playlistsApi.HandleGetPlaylists))
	v1ApisHandler.HandleFunc("PUT /playlist/public", authMiddleware.AuthApi(playlistsApi.HandleTogglePublicPlaylist))
	v1ApisHandler.HandleFunc("PUT /playlist/join", authMiddleware.AuthApi(playlistsApi.HandleToggleJoinPlaylist))
	v1ApisHandler.HandleFunc("GET /playlist/zip", authMiddleware.AuthApi(playlistsApi.HandleDonwnloadPlaylist))

	v1ApisHandler.HandleFunc("GET /player", authMiddleware.OptionalAuthApi(playerStateApi.HandleGetPlayerState))
	v1ApisHandler.HandleFunc("POST /player/shuffle", authMiddleware.OptionalAuthApi(playerStateApi.HandleSetShuffleOn))
	v1ApisHandler.HandleFunc("DELETE /player/shuffle", authMiddleware.OptionalAuthApi(playerStateApi.HandleSetShuffleOff))
	v1ApisHandler.HandleFunc("GET /player/song/next", authMiddleware.OptionalAuthApi(playerStateApi.HandleGetNextSongInQueue))
	v1ApisHandler.HandleFunc("GET /player/song/previous", authMiddleware.OptionalAuthApi(playerStateApi.HandleGetPreviousSongInQueue))
	v1ApisHandler.HandleFunc("GET /player/song/lyrics", authMiddleware.OptionalAuthApi(playerStateApi.HandleGetPlayingSongLyrics))
	v1ApisHandler.HandleFunc("PUT /player/loop/off", authMiddleware.OptionalAuthApi(playerStateApi.HandleSetLoopOff))
	v1ApisHandler.HandleFunc("PUT /player/loop/once", authMiddleware.OptionalAuthApi(playerStateApi.HandleSetLoopOnce))
	v1ApisHandler.HandleFunc("PUT /player/loop/all", authMiddleware.OptionalAuthApi(playerStateApi.HandleSetLoopAll))
	v1ApisHandler.HandleFunc("POST /player/queue/song/next", authMiddleware.OptionalAuthApi(playerStateApi.HandleAddSongToQueueNext))
	v1ApisHandler.HandleFunc("POST /player/queue/song/last", authMiddleware.OptionalAuthApi(playerStateApi.HandleAddSongToQueueAtLast))
	v1ApisHandler.HandleFunc("DELETE /player/queue/song", authMiddleware.OptionalAuthApi(playerStateApi.HandleRemoveSongFromQueue))
	v1ApisHandler.HandleFunc("POST /player/queue/playlist/next", authMiddleware.OptionalAuthApi(playerStateApi.HandleAddPlaylistToQueueNext))
	v1ApisHandler.HandleFunc("POST /player/queue/playlist/last", authMiddleware.OptionalAuthApi(playerStateApi.HandleAddPlaylistToQueueAtLast))

	v1ApisHandler.HandleFunc("GET /history", authMiddleware.AuthApi(historyApi.HandleGetHistoryItems))

	v1ApisHandler.HandleFunc("POST /library/favorite/song", authMiddleware.AuthApi(libraryApi.HandleAddSongToFavorites))
	v1ApisHandler.HandleFunc("DELETE /library/favorite/song", authMiddleware.AuthApi(libraryApi.HandleRemoveSongFromFavorites))
	v1ApisHandler.HandleFunc("GET /library/favorite/songs", authMiddleware.AuthApi(libraryApi.HandleGetFavoriteSongs))

	v1ApisHandler.HandleFunc("GET /me/profile", authMiddleware.AuthApi(accountApi.HandleGetProfile))
	v1ApisHandler.HandleFunc("GET /me/auth", authMiddleware.AuthApi(accountApi.HandleAuthCheck))
	v1ApisHandler.HandleFunc("GET /me/logout", authMiddleware.AuthApi(accountApi.HandleLogout))

	if config.Env().GoEnv == config.GoEnvTest {
		v1ApisHandler.HandleFunc("GET /test/otp", nil)
	}

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/v1/", http.StripPrefix("/v1", contenttype.Json(v1ApisHandler)))

	log.Info("Starting http server at port " + config.Env().Port)
	switch config.Env().GoEnv {
	case config.GoEnvBeta, config.GoEnvDev, config.GoEnvTest:
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, logger.Handler(applicationHandler)))
	case config.GoEnvProd:
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, minifyer.Middleware(applicationHandler)))
	}
}
