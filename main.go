package main

import (
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
	// fmt.Println(youtube.ScrapeSearch("lana del rey"))
	// return
	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/static/", http.FileServer(http.FS(static)))
	handlers.HandleHomePage(applicationHandler)
	handlers.HandleSearchResultsPage(applicationHandler, &youtube.YouTubeScraperSearch{})
	handlers.HandleSearchSugessions(applicationHandler)

	log.Info("Starting http server at port 3000")
	log.Fatalln(log.ErrorLevel, http.ListenAndServe(":3000", applicationHandler))
}
