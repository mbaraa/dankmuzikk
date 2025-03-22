package main

import (
	"dankmuzikk/evy"
	dankevents "dankmuzikk/evy/events"
	"dankmuzikk/handlers/events"
	"dankmuzikk/log"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

const (
	eventsBatchItems     = 25
	fetchWaitTimeSeconds = 5
)

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

func executeEvents(events []evy.EventPayload) error {
	wg := sync.WaitGroup{}
	wg.Add(len(events))

	for _, e := range events {
		log.Warningln("handling event", e.Topic)

		switch e.Topic {
		case "song-played":
			var body dankevents.SongPlayed
			err := json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				log.Errorf("failed unmarshalling event's json: %v\n", err)
				continue
			}

			go func() {
				err := errors.Join(
					handlers.HandleDownloadSongOnPlay(body),
					handlers.HandleAddSongToHistory(body),
					handlers.HandleIncrementSongPlaysInPlaylist(body),
				)
				if err != nil {
					log.Errorln("song-played", err)
				}

				wg.Done()
			}()
		case "song-downloaded":
			var body dankevents.SongDownloaded
			err := json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				log.Errorf("failed unmarshalling event's json: %v\n", err)
				continue
			}

			go func() {
				err := handlers.HandleMarkSongAsDownloaded(body)
				if err != nil {
					log.Errorln("song-downloaded", err)
				}

				wg.Done()
			}()
		case "song-added-to-playlist":
			var body dankevents.SongAddedToPlaylist
			err := json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				log.Errorf("failed unmarshalling event's json: %v\n", err)
				continue
			}

			go func() {
				err := errors.Join(
					handlers.HandleIncrementPlaylistSongsCount(body),
					handlers.HandleDownloadSongOnAddingToPlaylist(body),
				)
				if err != nil {
					log.Errorln("song-added-to-playlist", err)
				}

				wg.Done()
			}()
		case "song-removed-from-playlist":
			var body dankevents.SongRemovedFromPlaylist
			err := json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				log.Errorf("failed unmarshalling event's json: %v\n", err)
				continue
			}

			go func() {
				err := handlers.HandleDecrementPlaylistSongsCount(body)
				if err != nil {
					log.Errorln("song-removed-from-playlist", err)
				}

				wg.Done()
			}()
		case "songs-searched":
			var body dankevents.SongsSearched
			err := json.Unmarshal([]byte(e.Body), &body)
			if err != nil {
				log.Errorf("failed unmarshalling event's json: %v\n", err)
				continue
			}

			go func() {
				err := handlers.HandleSaveSongsMetadataOnSearchBatch(body)
				if err != nil {
					log.Errorln("song-searched", err)
				}

				wg.Done()
			}()
		}
		err := repo.DeleteEvent(e.Id)
		if err != nil {
			log.Errorf("Failed deleting event: %+v, error: %v\n", e, err)
			return err
		}
	}

	wg.Wait()
	return nil
}

func fetchAndExecuteEventsAsync() {
	timer := time.NewTicker(time.Second * fetchWaitTimeSeconds)
	wg := sync.WaitGroup{}
	for range timer.C {
		events, err := repo.GetEventsBatch(eventsBatchItems)
		if err != nil {
			continue
		}

		wg.Add(1)
		go func() {
			err = executeEvents(events)
			if err != nil {
				log.Errorln("Failed executing events batch", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func handleEventEmitted(w http.ResponseWriter, r *http.Request) {
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
}
