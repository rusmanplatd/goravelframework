package notification

import (
	"github.com/google/uuid"

	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
)

// BaseNotification provides a base implementation for notifications.
type BaseNotification struct {
	id     string
	locale string
}

// NewBaseNotification creates a new base notification.
func NewBaseNotification() *BaseNotification {
	return &BaseNotification{
		id: uuid.New().String(),
	}
}

// ID returns the notification ID.
func (n *BaseNotification) ID() string {
	if n.id == "" {
		n.id = uuid.New().String()
	}
	return n.id
}

// SetID sets the notification ID.
func (n *BaseNotification) SetID(id string) {
	n.id = id
}

// Locale returns the notification locale.
func (n *BaseNotification) Locale() string {
	return n.locale
}

// SetLocale sets the notification locale.
func (n *BaseNotification) SetLocale(locale string) {
	n.locale = locale
}

// Via returns the default channels (must be overridden).
func (n *BaseNotification) Via(notifiable any) []string {
	return []string{}
}

// ToDatabase returns the default database message (must be overridden).
func (n *BaseNotification) ToDatabase(notifiable any) *contractsnotification.DatabaseMessage {
	return contractsnotification.NewDatabaseMessage()
}

// ToMail returns nil by default (must be overridden if using mail channel).
func (n *BaseNotification) ToMail(notifiable any) any {
	return nil
}

// ToArray returns the default array representation (must be overridden).
func (n *BaseNotification) ToArray(notifiable any) map[string]any {
	return make(map[string]any)
}
