package notification

import (
	"fmt"
	"reflect"

	contractslog "github.com/goravel/framework/contracts/log"
	contractsnotification "github.com/goravel/framework/contracts/notification"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/notification/events"
)

// NotificationSender handles sending notifications.
type NotificationSender struct {
	manager *ChannelManager
	queue   contractsqueue.Queue
	log     contractslog.Log
}

// NewNotificationSender creates a new notification sender instance.
func NewNotificationSender(
	manager *ChannelManager,
	queue contractsqueue.Queue,
	log contractslog.Log,
) *NotificationSender {
	return &NotificationSender{
		manager: manager,
		queue:   queue,
		log:     log,
	}
}

// Send sends the given notification to the given notifiable entities.
func (s *NotificationSender) Send(notifiables any, notification contractsnotification.Notification) error {
	// Check if notification should be queued
	if shouldQueue, ok := notification.(contractsnotification.ShouldQueue); ok && shouldQueue.ShouldQueue() {
		return s.queueNotification(notifiables, notification)
	}

	// Send immediately
	return s.SendNow(notifiables, notification)
}

// SendNow sends the given notification immediately.
func (s *NotificationSender) SendNow(notifiables any, notification contractsnotification.Notification, channels ...string) error {
	notifiablesList := s.formatNotifiables(notifiables)

	for _, notifiable := range notifiablesList {
		// Get channels to send through
		viaChannels := channels
		if len(viaChannels) == 0 {
			viaChannels = notification.Via(notifiable)
		}

		if len(viaChannels) == 0 {
			continue
		}

		// Send through each channel
		for _, channelName := range viaChannels {
			if err := s.sendToNotifiable(notifiable, notification, channelName); err != nil {
				s.log.Error(fmt.Sprintf("Failed to send notification via %s: %v", channelName, err))
				// Continue with other channels even if one fails
			}
		}
	}

	return nil
}

// sendToNotifiable sends a notification to a single notifiable entity via a specific channel.
func (s *NotificationSender) sendToNotifiable(
	notifiable any,
	notification contractsnotification.Notification,
	channelName string,
) error {
	// Fire sending event
	s.dispatchSendingEvent(notifiable, notification, channelName)

	// Get the channel
	channel, err := s.manager.Channel(channelName)
	if err != nil {
		s.dispatchFailedEvent(notifiable, notification, channelName, err)
		return err
	}

	// Send the notification
	if err := channel.Send(notifiable, notification); err != nil {
		s.dispatchFailedEvent(notifiable, notification, channelName, err)
		return err
	}

	// Fire sent event
	s.dispatchSentEvent(notifiable, notification, channelName)

	return nil
}

// dispatchSendingEvent dispatches the NotificationSending event.
func (s *NotificationSender) dispatchSendingEvent(notifiable any, notification contractsnotification.Notification, channelName string) {
	if s.manager.event != nil {
		sendingEvent := &events.NotificationSending{
			Notifiable:   notifiable,
			Notification: notification,
			Channel:      channelName,
		}
		if err := s.manager.event.Job(sendingEvent, nil).Dispatch(); err != nil {
			s.log.Warning(fmt.Sprintf("Failed to dispatch NotificationSending event: %v", err))
		}
	}
}

// dispatchSentEvent dispatches the NotificationSent event.
func (s *NotificationSender) dispatchSentEvent(notifiable any, notification contractsnotification.Notification, channelName string) {
	if s.manager.event != nil {
		sentEvent := &events.NotificationSent{
			Notifiable:   notifiable,
			Notification: notification,
			Channel:      channelName,
			Response:     nil, // Channels don't currently return responses
		}
		if err := s.manager.event.Job(sentEvent, nil).Dispatch(); err != nil {
			s.log.Warning(fmt.Sprintf("Failed to dispatch NotificationSent event: %v", err))
		}
	}
}

// dispatchFailedEvent dispatches the NotificationFailed event.
func (s *NotificationSender) dispatchFailedEvent(notifiable any, notification contractsnotification.Notification, channelName string, notifErr error) {
	if s.manager.event != nil {
		failedEvent := &events.NotificationFailed{
			Notifiable:   notifiable,
			Notification: notification,
			Channel:      channelName,
			Error:        notifErr,
		}
		if err := s.manager.event.Job(failedEvent, nil).Dispatch(); err != nil {
			s.log.Warning(fmt.Sprintf("Failed to dispatch NotificationFailed event: %v", err))
		}
	}
}

// queueNotification queues the notification for later sending.
func (s *NotificationSender) queueNotification(notifiables any, notification contractsnotification.Notification) error {
	notifiablesList := s.formatNotifiables(notifiables)

	for _, notifiable := range notifiablesList {
		channels := notification.Via(notifiable)
		if len(channels) == 0 {
			continue
		}

		// Queue a job for each notifiable
		job := NewSendQueuedNotificationJob(notifiable, notification, channels)
		if err := s.queue.Job(job, []contractsqueue.Arg{}).Dispatch(); err != nil {
			return fmt.Errorf("failed to queue notification: %w", err)
		}
	}

	return nil
}

// formatNotifiables converts the notifiables parameter to a slice.
func (s *NotificationSender) formatNotifiables(notifiables any) []any {
	val := reflect.ValueOf(notifiables)

	// If it's already a slice, convert it
	if val.Kind() == reflect.Slice {
		result := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = val.Index(i).Interface()
		}
		return result
	}

	// Otherwise, return as single-item slice
	return []any{notifiables}
}
