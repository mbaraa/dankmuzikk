package evy

import (
	"bytes"
	"dankmuzikk/config"
	"dankmuzikk/evy/events"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type EventPayload struct {
	Id uint `gorm:"primaryKey;autoIncrement" json:"-"`

	Topic string `json:"topic"`
	Body  string `json:"body"`
}

type Evy struct {
	httpClient *http.Client
}

func New() *Evy {
	return &Evy{
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
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

	req, err := http.NewRequest(http.MethodPost, config.Env().EventHubAddress+"/emit", body)
	if err != nil {
		return err
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("event hub request failed")
	}

	return nil
}
