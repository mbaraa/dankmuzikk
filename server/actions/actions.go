package actions

import "dankmuzikk/app"

type Actions struct {
	app         *app.App
	eventhub    EventHub
	archiver    Archiver
	blobstorage BlobStorage
	jwt         JwtManager[TokenPayload]
	mailer      Mailer
	youtube     YouTube
	lyrics      Lyrics
}

func New(
	app *app.App,
	eventhub EventHub,
	archiver Archiver,
	blobstorage BlobStorage,
	jwt JwtManager[TokenPayload],
	mailer Mailer,
	youtube YouTube,
	lyrics Lyrics,
) *Actions {
	return &Actions{
		app:         app,
		eventhub:    eventhub,
		archiver:    archiver,
		blobstorage: blobstorage,
		jwt:         jwt,
		mailer:      mailer,
		youtube:     youtube,
		lyrics:      lyrics,
	}
}
