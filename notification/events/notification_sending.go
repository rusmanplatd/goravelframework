package events

import (
	contractsevent "github.com/rusmanplatd/goravelframework/contracts/event"
	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
)

// NotificationSending is fired before a notification is sent.
type NotificationSending struct {
	Notifiable   any
	Notification contractsnotification.Notification
	Channel      string
}

// Handle handles the event.
func (e *NotificationSending) Handle(args []contractsevent.Arg) ([]contractsevent.Arg, error) {
	return args, nil
}
