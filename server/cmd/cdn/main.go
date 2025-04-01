package main

import (
	"dankmuzikk/config"
	"dankmuzikk/evy"
	"dankmuzikk/evy/events"
	"dankmuzikk/handlers/middlewares/logger"
	"dankmuzikk/log"
	"dankmuzikk/mariadb"
	"net/http"
	"strings"
	"time"
)

func main() {
	mariadbRepo, err := mariadb.New()
	if err != nil {
		log.Fatalln(err)
	}
	// app := app.New(mariadbRepo)
	eventhub := evy.New()
	// jwtUtil := jwt.New[actions.TokenPayload]()
	// usecases := actions.New(
	// 	app,
	// 	eventhub,
	// 	nil,
	// 	jwtUtil,
	// 	nil,
	// 	nil,
	// )
	// authMw := auth.New(usecases)
	applicationHandler := http.NewServeMux()

	muzikkxDir := config.Env().BlobsDir + "/muzikkx/"
	pixDir := config.Env().BlobsDir + "/pix/"
	playlistsDir := config.Env().BlobsDir + "/playlists/"

	applicationHandler.Handle("/muzikkx/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=5")
		w.Header().Set("Icy-Metadata", "1")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Disposition", "inline")
		http.
			StripPrefix("/muzikkx", http.FileServer(http.Dir(muzikkxDir))).
			ServeHTTP(w, r)
	}))

	applicationHandler.Handle("/muzikkx-raw/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/muzikkx-raw/"), ".mp3")
		song, err := mariadbRepo.GetSongByYouTubeId(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=5")
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+song.Title+".mp3")
		http.
			StripPrefix("/muzikkx-raw", http.FileServer(http.Dir(muzikkxDir))).
			ServeHTTP(w, r)
	}))

	applicationHandler.Handle("/pix/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=5")
		http.
			StripPrefix("/pix", http.FileServer(http.Dir(pixDir))).
			ServeHTTP(w, r)
	}))

	applicationHandler.Handle("/playlists/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/playlists/"), ".zip")
		eventhub.Publish(events.PlaylistDownloaded{
			PlaylistId: id,
			DeleteAt:   time.Now().Add(time.Hour),
		})

		pl, err := mariadbRepo.GetPlaylistByPublicId(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+pl.Title+".zip")
		http.
			StripPrefix("/playlists", http.FileServer(http.Dir(playlistsDir))).
			ServeHTTP(w, r)
	}))

	log.Info("Starting http cdn server at port " + config.Env().CdnPort)
	if config.Env().GoEnv == "dev" || config.Env().GoEnv == "beta" {
		log.Fatalln(http.ListenAndServe(":"+config.Env().CdnPort, logger.Handler(applicationHandler)))
	}
	log.Fatalln(http.ListenAndServe(":"+config.Env().CdnPort, applicationHandler))
}
