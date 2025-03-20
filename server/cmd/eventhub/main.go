package main

import (
	"dankmuzikk/actions"
	"dankmuzikk/app"
	"dankmuzikk/config"
	"dankmuzikk/evy"
	dankevents "dankmuzikk/evy/events"
	"dankmuzikk/handlers/events"
	"dankmuzikk/jwt"
	"dankmuzikk/log"
	"dankmuzikk/mailer"
	"dankmuzikk/mariadb"
	"dankmuzikk/youtube"
	"dankmuzikk/zip"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const eventsBatchSize = 25

var (
	repo     evy.Repository
	handlers *events.EventHandlers
)

type eventHub struct{}

func (e *eventHub) Publish(event dankevents.Event) error {
	eventBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	fullEvent := evy.EventPayload{
		Topic: event.Topic(),
		Body:  string(eventBody),
	}

	return repo.CreateEvent(fullEvent)
}

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

func executeEvents(events []evy.EventPayload) error {
	for _, e := range events {
		log.Warningln("handling event", e.Topic)

		var err error
		switch e.Topic {
		case "song-played":
			var body dankevents.SongPlayed
			err = json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				return err
			}

			err = errors.Join(
				handlers.HandleDownloadSongOnPlay(body),
				handlers.HandleAddSongToHistory(body),
				handlers.HandleIncrementSongPlaysInPlaylist(body),
			)
		case "song-downloaded":
			var body dankevents.SongDownloaded
			err = json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				return err
			}

			err = handlers.HandleMarkSongAsDownloaded(body)
		case "song-added-to-playlist":
			var body dankevents.SongAddedToPlaylist
			err = json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				return err
			}

			err = errors.Join(
				handlers.HandleIncrementPlaylistSongsCount(body),
				handlers.HandleDownloadSongOnAddingToPlaylist(body),
			)
		case "song-removed-from-playlist":
			var body dankevents.SongRemovedFromPlaylist
			err = json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				return err
			}

			err = handlers.HandleDecrementPlaylistSongsCount(body)
		}
		err2 := repo.DeleteEvent(e.Id)
		if err2 != nil {
			log.Errorf("Failed deleting event: %+v, error: %v\n", e, err2)
			return err2
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// TODO: make this concurrent :)
	timer := time.NewTicker(time.Second * 3)
	go func() {
		for range timer.C {
			events, err := repo.GetEventsBatch(eventsBatchSize)
			if err != nil {
				log.Errorln("Failed getting events", err)
				continue
			}

			err = executeEvents(events)
			if err != nil {
				log.Errorln("Failed executing events batch", err)
				continue
			}
		}
	}()

	http.HandleFunc("/emit", func(w http.ResponseWriter, r *http.Request) {
		var body evy.EventPayload
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Errorln("Failed marshalling event", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = repo.CreateEvent(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorln(err)
			return
		}
	})

	http.ListenAndServe(":"+config.Env().EventHubPort, nil)
}
