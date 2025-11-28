package notification

import (
	"encoding/json"

	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
)

// SendQueuedNotificationJob is a queue job for sending notifications.
type SendQueuedNotificationJob struct {
	notifiable   any
	notification contractsnotification.Notification
	channels     []string
	sender       *NotificationSender // Injected dependency
}

// NewSendQueuedNotificationJob creates a new queued notification job.
func NewSendQueuedNotificationJob(
	notifiable any,
	notification contractsnotification.Notification,
	channels []string,
	sender *NotificationSender,
) *SendQueuedNotificationJob {
	return &SendQueuedNotificationJob{
		notifiable:   notifiable,
		notification: notification,
		channels:     channels,
		sender:       sender,
	}
}

// Signature returns the job signature.
func (j *SendQueuedNotificationJob) Signature() string {
	return "send_queued_notification"
}

// Handle handles the job.
func (j *SendQueuedNotificationJob) Handle(...any) error {
	// Use the injected sender instead of accessing through foundation
	return j.sender.SendNow(j.notifiable, j.notification, j.channels...)
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
