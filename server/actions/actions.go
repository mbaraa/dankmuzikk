package actions

import "dankmuzikk/app"

type Actions struct {
	app      *app.App
	eventhub EventHub
	archiver Archiver
	jwt      JwtManager[TokenPayload]
	mailer   Mailer
	youtube  YouTube
}

func New(
	app *app.App,
	eventhub EventHub,
	archiver Archiver,
	jwt JwtManager[TokenPayload],
	mailer Mailer,
	youtube YouTube,
) *Actions {
	return &Actions{
		app:      app,
		eventhub: eventhub,
		archiver: archiver,
		jwt:      jwt,
		mailer:   mailer,
		youtube:  youtube,
	}
}
