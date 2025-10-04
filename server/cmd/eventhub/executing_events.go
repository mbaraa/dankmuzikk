package main

import (
	"crypto/md5"
	"dankmuzikk/evy"
	"encoding/hex"
	"fmt"
	"sync"
)

type executingEvents struct {
	currentEvents map[string]struct{}
	mu            sync.RWMutex
}

func eventId(event evy.EventPayload) string {
	hasher := md5.New()
	hasher.Write([]byte(event.Body))
	eventBodyHash := hex.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf("%s-%s", event.Topic, eventBodyHash)
}

func (e *executingEvents) Add(event evy.EventPayload) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.currentEvents[eventId(event)] = struct{}{}
}

func (e *executingEvents) Exists(event evy.EventPayload) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	_, exists := e.currentEvents[eventId(event)]
	return exists
}

func (e *executingEvents) Delete(event evy.EventPayload) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.currentEvents, eventId(event))
}
