package main

import (
	"dankmuzikk/config"
	"dankmuzikk/evy"
	"dankmuzikk/evy/events"
	"dankmuzikk/handlers/middlewares/logger"
	"dankmuzikk/log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// mariadbRepo, err := mariadb.New()
	// if err != nil {
	// log.Fatalln(err)
	// }
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

	applicationHandler.Handle("/muzikkx/", http.StripPrefix("/muzikkx", http.FileServer(http.Dir(muzikkxDir))))
	applicationHandler.Handle("/muzikkx-raw/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment")
		http.
			StripPrefix("/muzikkx-raw", http.FileServer(http.Dir(muzikkxDir))).
			ServeHTTP(w, r)
	}))

	applicationHandler.Handle("/pix/", http.StripPrefix("/pix", http.FileServer(http.Dir(pixDir))))
	applicationHandler.Handle("/playlists/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment")

		title := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/playlists/"), ".zip")
		eventhub.Publish(events.PlaylistDownloaded{
			PlaylistTitle: title,
			DeleteAt:      time.Now().Add(time.Minute * 5),
		})

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
