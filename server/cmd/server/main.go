package main

import (
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/handlers/apis"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/handlers/middlewares/logger"
	"dankmuzikk/log"
	"dankmuzikk/services/archive"
	"dankmuzikk/services/history"
	"dankmuzikk/services/jwt"
	"dankmuzikk/services/login"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/playlists/songs"
	"dankmuzikk/services/youtube/download"
	"net/http"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
)

//go:generate templ generate -path ../../views/

func main() {
	err := StartServer()
	if err != nil {
		log.Fatalln(log.ErrorLevel, err)
	}
}

func StartServer() error {
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

	minfyer := minify.New()
	minfyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minfyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)

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
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Info("Starting http server at port " + config.Env().Port)
	if config.Env().GoEnv == "dev" || config.Env().GoEnv == "beta" {
		return http.ListenAndServe(":"+config.Env().Port, logger.Handler(applicationHandler))
	}
	return http.ListenAndServe(":"+config.Env().Port, minfyer.Middleware(applicationHandler))
}
