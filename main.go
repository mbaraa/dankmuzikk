package main

import (
	"dankmuzikk/config"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube"
	"embed"
	"net/http"
)

//go:embed static/*
var static embed.FS

//go:generate npx tailwindcss build -i static/css/style.css -o static/css/tailwind.css -m

func main() {
	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/static/", http.FileServer(http.FS(static)))
	handlers.HandleHomePage(applicationHandler)
	handlers.HandleSearchResultsPage(applicationHandler, &youtube.YouTubeScraperSearch{})
	handlers.HandleSearchSugessions(applicationHandler)
	handlers.HandleServeSongs(applicationHandler)
	handlers.HandleDownloadSong(applicationHandler)

	log.Info("Starting http server at port " + config.Vals().Port)
	log.Fatalln(log.ErrorLevel, http.ListenAndServe(":"+config.Vals().Port, applicationHandler))
}
