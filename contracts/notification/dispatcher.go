package notification

// Dispatcher represents a notification dispatcher.
type Dispatcher interface {
	// Send sends the given notification to the given notifiable entities.
	// If the notification implements ShouldQueue, it will be queued.
	Send(notifiables any, notification Notification) error

	// SendNow sends the given notification immediately to the given notifiable entities.
	// The notification will not be queued even if it implements ShouldQueue.
	SendNow(notifiables any, notification Notification, channels ...string) error
}
