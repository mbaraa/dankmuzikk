package server

import (
	"dankmuzikk/config"
	"dankmuzikk/handlers/apis"
	"dankmuzikk/handlers/pages"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube"
	"embed"
	"net/http"
)

func StartServer(staticFS embed.FS) error {
	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/static/", http.FileServer(http.FS(staticFS)))
	pagesHandler.Handle("/music/", http.StripPrefix("/music", http.FileServer(http.Dir(config.Env().YouTube.MusicDir))))

	pagesHandler.HandleFunc("/", pages.Handler(pages.HandleHomePage))
	pagesHandler.HandleFunc("/login", pages.Handler(pages.HandleLoginPage))
	pagesHandler.HandleFunc("/profile", pages.Handler(pages.HandleProfilePage))
	pagesHandler.HandleFunc("/about", pages.Handler(pages.HandleAboutPage))
	pagesHandler.HandleFunc("/playlists", pages.Handler(pages.HandlePlaylistsPage))
	pagesHandler.HandleFunc("/privacy", pages.Handler(pages.HandlePrivacyPage))
	pagesHandler.HandleFunc("/tos", pages.Handler(pages.HandleTOSPage))
	pagesHandler.HandleFunc("/search", pages.Handler(pages.HandleSearchResultsPage(&youtube.YouTubeScraperSearch{})))

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("/login/google", apis.HandleGoogleOAuthLogin)
	apisHandler.HandleFunc("/login/google/callback", apis.HandleGoogleOAuthLoginCallback)
	apisHandler.HandleFunc("/search-suggession", apis.HandleSearchSugessions)
	apisHandler.HandleFunc("/song/download/{youtube_video_id}", apis.HandleDownloadSong)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Info("Starting http server at port " + config.Env().Port)
	return http.ListenAndServe(":"+config.Env().Port, applicationHandler)
}
