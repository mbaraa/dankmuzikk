package main

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/config"
	"dankmuzikk-web/handlers/apis"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/handlers/middlewares/contenttype"
	"dankmuzikk-web/handlers/middlewares/ismobile"
	"dankmuzikk-web/handlers/middlewares/logger"
	"dankmuzikk-web/handlers/middlewares/theme"
	"dankmuzikk-web/handlers/pages"
	"dankmuzikk-web/log"
	"dankmuzikk-web/requests"
	"embed"
	"net/http"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

//go:generate templ generate

var (
	//go:embed static/*
	static embed.FS

	minifyer *minify.M

	usecases       *actions.Actions
	authMiddleware *auth.Middleware
)

func init() {
	minifyer = minify.New()
	minifyer.AddFunc("text/css", css.Minify)
	minifyer.AddFunc("text/html", html.Minify)
	minifyer.AddFunc("image/svg+xml", svg.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	reqs := requests.New()
	usecases = actions.New(reqs)
	authMiddleware = auth.New(usecases)
}

func main() {
	pagesHandler := http.NewServeMux()
	switch config.Env().GoEnv {
	case "prod":
		pagesHandler.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case strings.HasPrefix(r.URL.Path, "/static/js"):
				w.Header().Set("Cache-Control", "public, max-age=7200, stale-while-revalidate=5")
			case strings.HasPrefix(r.URL.Path, "/static/css"):
				w.Header().Set("Cache-Control", "public, max-age=7200, stale-while-revalidate=5")
			default:
				w.Header().Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=5")
			}

			minifyer.Middleware(http.FileServer(http.FS(static))).
				ServeHTTP(w, r)
		}))
	case "beta":
		pagesHandler.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Env().GoEnv != "dev" {
				switch {
				case strings.HasPrefix(r.URL.Path, "/static/js"):
					w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=5")
				case strings.HasPrefix(r.URL.Path, "/static/css"):
					w.Header().Set("Cache-Control", "public, max-age=600, stale-while-revalidate=5")
				default:
					w.Header().Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=5")
				}
			}

			http.FileServer(http.FS(static)).
				ServeHTTP(w, r)
		}))
	default:
		pagesHandler.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.FileServer(http.FS(static)).
				ServeHTTP(w, r)
		}))
	}

	pagesHandler.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		robotsFile, _ := static.ReadFile("static/robots.txt")
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write(robotsFile)
	})

	pages := pages.New(usecases)
	pagesHandler.HandleFunc("/", contenttype.Html(authMiddleware.OptionalAuthPage(pages.HandleHomePage)))
	pagesHandler.HandleFunc("GET /signup", contenttype.Html(authMiddleware.AuthPage(pages.HandleSignupPage)))
	pagesHandler.HandleFunc("GET /login", contenttype.Html(authMiddleware.AuthPage(pages.HandleLoginPage)))
	pagesHandler.HandleFunc("GET /profile", contenttype.Html(authMiddleware.AuthPage(pages.HandleProfilePage)))
	pagesHandler.HandleFunc("GET /about", contenttype.Html(pages.HandleAboutPage))
	pagesHandler.HandleFunc("GET /playlists", contenttype.Html(authMiddleware.AuthPage(pages.HandlePlaylistsPage)))
	pagesHandler.HandleFunc("GET /playlist/{playlist_id}", contenttype.Html(authMiddleware.AuthPage(pages.HandleSinglePlaylistPage)))
	pagesHandler.HandleFunc("GET /song/{song_id}", contenttype.Html(authMiddleware.OptionalAuthPage(pages.HandleSingleSongPage)))
	pagesHandler.HandleFunc("GET /privacy", contenttype.Html(pages.HandlePrivacyPage))
	pagesHandler.HandleFunc("GET /search", contenttype.Html(authMiddleware.OptionalAuthPage(pages.HandleSearchResultsPage)))

	emailLoginApi := apis.NewEmailLoginApi(usecases)
	googleLoginApi := apis.NewGoogleLoginApi(usecases)
	songApi := apis.NewDownloadHandler(usecases)
	playlistsApi := apis.NewPlaylistApi(usecases)
	historyApi := apis.NewHistoryApi(usecases)
	logoutApi := apis.NewLogoutApi(usecases)
	searchSuggestionsApi := apis.NewSearchSiggestionsApi(usecases)

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("POST /login/email", emailLoginApi.HandleEmailLogin)
	apisHandler.HandleFunc("POST /signup/email", emailLoginApi.HandleEmailSignup)
	apisHandler.HandleFunc("POST /verify-otp", emailLoginApi.HandleEmailOTPVerification)
	apisHandler.HandleFunc("GET /login/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("GET /signup/google", googleLoginApi.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("/login/google/callback", googleLoginApi.HandleGoogleOAuthLoginCallback)
	apisHandler.HandleFunc("GET /logout", authMiddleware.AuthApi(logoutApi.HandleLogout))
	apisHandler.HandleFunc("GET /search-suggestion", searchSuggestionsApi.HandleSearchSuggestions)
	apisHandler.HandleFunc("GET /song", authMiddleware.OptionalAuthApi(songApi.HandlePlaySong))
	apisHandler.HandleFunc("GET /song/single", authMiddleware.OptionalAuthApi(songApi.HandleGetSong))
	apisHandler.HandleFunc("PUT /song/playlist", authMiddleware.AuthApi(songApi.HandleToggleSongInPlaylist))
	apisHandler.HandleFunc("PUT /song/playlist/upvote", authMiddleware.AuthApi(songApi.HandleUpvoteSongPlaysInPlaylist))
	apisHandler.HandleFunc("PUT /song/playlist/downvote", authMiddleware.AuthApi(songApi.HandleDownvoteSongPlaysInPlaylist))
	apisHandler.HandleFunc("GET /song/lyrics", songApi.HandleGetSongLyrics)
	apisHandler.HandleFunc("GET /playlist/all", authMiddleware.AuthApi(playlistsApi.HandleGetPlaylistsForPopover))
	apisHandler.HandleFunc("GET /playlist", authMiddleware.AuthApi(playlistsApi.HandleGetPlaylist))
	apisHandler.HandleFunc("POST /playlist", authMiddleware.AuthApi(playlistsApi.HandleCreatePlaylist))
	apisHandler.HandleFunc("PUT /playlist/public", authMiddleware.AuthApi(playlistsApi.HandleTogglePublicPlaylist))
	apisHandler.HandleFunc("PUT /playlist/join", authMiddleware.AuthApi(playlistsApi.HandleToggleJoinPlaylist))
	apisHandler.HandleFunc("DELETE /playlist", authMiddleware.AuthApi(playlistsApi.HandleDeletePlaylist))
	apisHandler.HandleFunc("GET /playlist/zip", authMiddleware.AuthApi(playlistsApi.HandleDonwnloadPlaylist))
	apisHandler.HandleFunc("GET /history/{page}", authMiddleware.AuthApi(historyApi.HandleGetMoreHistoryItems))

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", ismobile.Handler(theme.Handler(pagesHandler)))
	applicationHandler.Handle("/api/", ismobile.Handler(theme.Handler(http.StripPrefix("/api", apisHandler))))

	log.Info("Starting http server at port " + config.Env().Port)
	if config.Env().GoEnv == "dev" || config.Env().GoEnv == "beta" {
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, logger.Handler(ismobile.Handler(theme.Handler(applicationHandler)))))
	}
	log.Fatalln(http.ListenAndServe(":"+config.Env().Port, ismobile.Handler(theme.Handler(minifyer.Middleware(applicationHandler)))))
}
