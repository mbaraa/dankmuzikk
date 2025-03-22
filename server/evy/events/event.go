package events

// Event holds data to be sent to the event hub to be executed in the background.
type Event interface {
	// Topic returns a string which is the topic of the event.
	Topic() string
}
