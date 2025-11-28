package events

import (
	contractsevent "github.com/rusmanplatd/goravelframework/contracts/event"
	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
)

// NotificationSent is fired after a notification is successfully sent.
type NotificationSent struct {
	Notifiable   any
	Notification contractsnotification.Notification
	Channel      string
	Response     any
}

// Handle handles the event.
func (e *NotificationSent) Handle(args []contractsevent.Arg) ([]contractsevent.Arg, error) {
	return args, nil
}
