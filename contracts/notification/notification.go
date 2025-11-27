package notification

import "github.com/goravel/framework/contracts/mail"

// Notification represents a notification that can be sent through various channels.
type Notification interface {
	// Via returns the channels the notification should be sent through.
	Via(notifiable any) []string

	// ToDatabase returns the database representation of the notification.
	// This method is called when the notification is sent via the database channel.
	ToDatabase(notifiable any) *DatabaseMessage

	// ToMail returns the mail representation of the notification.
	// This method is called when the notification is sent via the mail channel.
	ToMail(notifiable any) mail.Mailable

	// ToArray returns the array representation of the notification.
	// This is a fallback method used when ToDatabase is not implemented.
	ToArray(notifiable any) map[string]any
}

// ShouldQueue indicates that a notification should be queued.
type ShouldQueue interface {
	// ShouldQueue returns true if the notification should be queued.
	ShouldQueue() bool
}

// HasLocale indicates that a notification has a locale preference.
type HasLocale interface {
	// Locale returns the locale for the notification.
	Locale() string
}

// ShouldSend allows notifications to conditionally send based on notifiable and channel.
// This is called before the NotificationSending event is dispatched.
type ShouldSend interface {
	// ShouldSend returns true if the notification should be sent to the given notifiable via the given channel.
	ShouldSend(notifiable any, channel string) bool
}

// HasLocalePreference indicates that a notifiable has a locale preference.
// This is used to determine the locale for sending notifications.
type HasLocalePreference interface {
	// PreferredLocale returns the preferred locale for the notifiable.
	PreferredLocale() string
}

// ViaConnections allows notifications to specify different queue connections per channel.
type ViaConnections interface {
	// ViaConnections returns a map of channel names to queue connection names.
	ViaConnections() map[string]string
}

// ViaQueues allows notifications to specify different queue names per channel.
type ViaQueues interface {
	// ViaQueues returns a map of channel names to queue names.
	ViaQueues() map[string]string
}

// WithDelay allows notifications to specify a delay before sending.
type WithDelay interface {
	// WithDelay returns the delay in seconds before the notification should be sent.
	// Can return a map of channel names to delays for per-channel delays.
	WithDelay(notifiable any, channel string) int
}
