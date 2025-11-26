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
