package main

import (
	"dankmuzikk/handlers"
	"dankmuzikk/log"
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
	handlers.HandleSearchResultsPage(applicationHandler)
	handlers.HandleSearchSugessions(applicationHandler)

	log.Info("Starting http server at port 8080")
	log.Fatalln(log.ErrorLevel, http.ListenAndServe(":8080", applicationHandler))
}
