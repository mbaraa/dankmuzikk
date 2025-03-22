package actions

import "dankmuzikk/evy/events"

// EventHub handles events publishing.
type EventHub interface {
	// Publish publishes the given event to the given eventhub implementation.
	Publish(event events.Event) error
}
