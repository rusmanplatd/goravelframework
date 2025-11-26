package notification

// Channel represents a notification delivery channel.
type Channel interface {
	// Send sends the given notification to the given notifiable entity.
	Send(notifiable any, notification Notification) error
}
