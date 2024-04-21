package main

import (
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube"
	"embed"
	"encoding/json"
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
	applicationHandler.Handle("/", handlers.HandleHomePage())
	applicationHandler.Handle("/static/", http.FileServer(http.FS(static)))

	applicationHandler.HandleFunc("/api/search-suggession", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("query")
		suggessions, err := youtube.SearchSuggestions(q)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(suggessions)
	})

	log.Info("Starting http server at port 8080")
	log.Fatalln(log.ErrorLevel, http.ListenAndServe(":8080", applicationHandler))
}
