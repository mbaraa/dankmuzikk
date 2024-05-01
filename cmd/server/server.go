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
	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/static/", http.FileServer(http.FS(staticFS)))
	pages.HandleHomePage(applicationHandler)
	pages.HandleSearchResultsPage(applicationHandler, &youtube.YouTubeScraperSearch{})
	pages.HandlePrivacyPage(applicationHandler)
	pages.HandleTOSPage(applicationHandler)
	pages.HandleAboutPage(applicationHandler)
	pages.HandleProfilePage(applicationHandler)
	pages.HandlePlaylistsPage(applicationHandler)
	pages.HandleLoginPage(applicationHandler)

	apis.HandleSearchSugessions(applicationHandler)
	apis.HandleServeSongs(applicationHandler)
	apis.HandleDownloadSong(applicationHandler)
	apis.HandleGoogleOAuthLogin(applicationHandler)
	apis.HandleGoogleOAuthLoginCallback(applicationHandler)

	log.Info("Starting http server at port " + config.Env().Port)
	return http.ListenAndServe(":"+config.Env().Port, applicationHandler)
}
