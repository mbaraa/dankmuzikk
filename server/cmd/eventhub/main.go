package main

import (
	"dankmuzikk/actions"
	"dankmuzikk/app"
	"dankmuzikk/config"
	"dankmuzikk/handlers/events"
	"dankmuzikk/jwt"
	"dankmuzikk/log"
	"dankmuzikk/mailer"
	"dankmuzikk/mariadb"
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

	app := app.New(mariadbRepo)
	zipArchiver := zip.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	mailer := mailer.New()
	yt := youtube.New()

	usecases := actions.New(
		app,
		&eventHub{},
		zipArchiver,
		jwtUtil,
		mailer,
		yt,
	)

	handlers = events.New(usecases)
}

func main() {
	go fetchAndExecuteEventsAsync()

	http.HandleFunc("/emit", handleEventEmitted)
	http.ListenAndServe(":"+config.Env().EventHubPort, nil)
}
