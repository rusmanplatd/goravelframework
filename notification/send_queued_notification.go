package notification

import (
	"encoding/json"
	"fmt"

	contractsnotification "github.com/goravel/framework/contracts/notification"
	"github.com/goravel/framework/foundation"
)

// SendQueuedNotificationJob is a queue job for sending notifications.
type SendQueuedNotificationJob struct {
	notifiable   any
	notification contractsnotification.Notification
	channels     []string
}

// NewSendQueuedNotificationJob creates a new queued notification job.
func NewSendQueuedNotificationJob(
	notifiable any,
	notification contractsnotification.Notification,
	channels []string,
) *SendQueuedNotificationJob {
	return &SendQueuedNotificationJob{
		notifiable:   notifiable,
		notification: notification,
		channels:     channels,
	}
}

// Signature returns the job signature.
func (j *SendQueuedNotificationJob) Signature() string {
	return "send_queued_notification"
}

// Handle handles the job.
func (j *SendQueuedNotificationJob) Handle(...any) error {
	// Get the notification factory from the application
	notificationFactory := foundation.App.MakeNotification()
	if notificationFactory == nil {
		return fmt.Errorf("notification facade not available")
	}

	// Send the notification immediately (already queued, so don't queue again)
	return notificationFactory.SendNow(j.notifiable, j.notification, j.channels...)
}

// Marshal serializes the job data.
func (j *SendQueuedNotificationJob) Marshal() ([]byte, error) {
	data := map[string]any{
		"notifiable":   j.notifiable,
		"notification": j.notification,
		"channels":     j.channels,
	}
	return json.Marshal(data)
}

// Unmarshal deserializes the job data.
func (j *SendQueuedNotificationJob) Unmarshal(data []byte) error {
	var jobData map[string]any
	if err := json.Unmarshal(data, &jobData); err != nil {
		return err
	}

	// Note: Proper unmarshaling would require type information
	// This is a simplified version
	if notifiable, ok := jobData["notifiable"]; ok {
		j.notifiable = notifiable
	}
	if channels, ok := jobData["channels"].([]string); ok {
		j.channels = channels
	}

	return nil
}
