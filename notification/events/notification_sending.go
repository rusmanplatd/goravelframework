package events

import (
	contractsevent "github.com/goravel/framework/contracts/event"
	contractsnotification "github.com/goravel/framework/contracts/notification"
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
