package app

// DankError is implemented for every error around here :)
type DankError interface {
	error
	// ClientStatusCode the HTTP status for clients.
	ClientStatusCode() int
	// ExtraData any data that will be helpful for clients for better UX context.
	ExtraData() map[string]any
	// ExposeToClients reports whether to expose this error to clients or not.
	ExposeToClients() bool
}
