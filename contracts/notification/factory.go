package notification

// Factory represents a notification factory.
type Factory interface {
	Dispatcher

	// Channel gets a channel instance by name.
	Channel(name string) (Channel, error)
}
