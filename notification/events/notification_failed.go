package events

import (
	contractsevent "github.com/rusmanplatd/goravelframework/contracts/event"
	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
)

// NotificationFailed is fired when a notification fails to send.
type NotificationFailed struct {
	Notifiable   any
	Notification contractsnotification.Notification
	Channel      string
	Error        error
}

// Handle handles the event.
func (e *NotificationFailed) Handle(args []contractsevent.Arg) ([]contractsevent.Arg, error) {
	return args, nil
}
