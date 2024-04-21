package main

import (
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube"
	"embed"
	"fmt"
	"net/http"
)

//go:embed static/*
var static embed.FS

//go:generate npx tailwindcss build -i static/css/style.css -o static/css/tailwind.css -m

func main() {
	vids, err := youtube.Search("lana del rey raise me up")
	if err != nil {
		panic(err)
	}

	for _, vid := range vids {
		fmt.Printf("%+v\n", vid)
	}

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/static/", http.FileServer(http.FS(static)))
	handlers.HandleHomePage(applicationHandler)

	log.Info("Starting http server at port 8080")
	log.Fatalln(log.ErrorLevel, http.ListenAndServe(":8080", applicationHandler))
}
