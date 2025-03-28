package main

import (
	"dankmuzikk/actions"
	"dankmuzikk/app"
	"dankmuzikk/blobs"
	"dankmuzikk/config"
	"dankmuzikk/handlers/events"
	"dankmuzikk/jwt"
	"dankmuzikk/log"
	"dankmuzikk/mailer"
	"dankmuzikk/mariadb"
	"dankmuzikk/redis"
	"dankmuzikk/youtube"
	"dankmuzikk/zip"
	"net/http"
)

func init() {
	mariadbRepo, err := mariadb.New()
	if err != nil {
		log.Fatalln(err)
	}

	repo = mariadbRepo

	cache := redis.New()

	app := app.New(mariadbRepo, cache)
	zipArchiver := zip.New()
	blobstorage := blobs.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	mailer := mailer.New()
	yt := youtube.New()

	usecases := actions.New(
		app,
		cache,
		&eventHub{},
		zipArchiver,
		blobstorage,
		jwtUtil,
		mailer,
		yt,
		nil,
	)

	handlers = events.New(usecases)
}

func main() {
	go fetchAndExecuteEventsAsync()

	http.HandleFunc("/emit", handleEventEmitted)
	http.ListenAndServe(":"+config.Env().EventHubPort, nil)
}
