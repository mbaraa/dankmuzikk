package evy

import (
	"bytes"
	"dankmuzikk/config"
	"dankmuzikk/evy/events"
	"encoding/json"
	"errors"
	"net/http"
)

type EventPayload struct {
	Id uint `gorm:"primaryKey;autoIncrement" json:"-"`

	Topic string `json:"topic"`
	Body  string `json:"body"`
}

type Evy struct {
}

func New() *Evy {
	return &Evy{}
}

func (e *Evy) Publish(event events.Event) error {
	eventBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	fullEvent := EventPayload{
		Topic: event.Topic(),
		Body:  string(eventBody),
	}

	body := bytes.NewBuffer(nil)
	err = json.NewEncoder(body).Encode(fullEvent)
	if err != nil {
		return err
	}

	resp, err := http.Post(config.Env().EventHubAddress+"/emit", "application/json", body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("event hub request failed")
	}

	return nil
}
