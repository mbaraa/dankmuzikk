package server

import (
	"dankmuzikk/config"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube"
	"embed"
	"net/http"
)

func StartServer(staticFS embed.FS) error {
	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/static/", http.FileServer(http.FS(staticFS)))
	handlers.HandleHomePage(applicationHandler)
	handlers.HandleSearchResultsPage(applicationHandler, &youtube.YouTubeScraperSearch{})
	handlers.HandleSearchSugessions(applicationHandler)
	handlers.HandleServeSongs(applicationHandler)
	handlers.HandleDownloadSong(applicationHandler)

	log.Info("Starting http server at port " + config.Vals().Port)
	return http.ListenAndServe(":"+config.Vals().Port, applicationHandler)
}
