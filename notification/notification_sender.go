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

		// Determine locale for this notifiable
		locale := s.getPreferredLocale(notifiable, notification)

		// Send through each channel with locale context
		s.sendWithLocale(locale, func() {
			for _, channelName := range viaChannels {
				if err := s.sendToNotifiable(notifiable, notification, channelName); err != nil {
					s.log.Error(fmt.Sprintf("Failed to send notification via %s: %v", channelName, err))
					// Continue with other channels even if one fails
				}
			}
		})
	}

	return nil
}

// getPreferredLocale determines the preferred locale for sending a notification.
func (s *NotificationSender) getPreferredLocale(notifiable any, notification contractsnotification.Notification) string {
	// Check if notification has a locale
	if hasLocale, ok := notification.(contractsnotification.HasLocale); ok {
		if locale := hasLocale.Locale(); locale != "" {
			return locale
		}
	}

	// Check if notifiable has a locale preference
	if hasPreference, ok := notifiable.(contractsnotification.HasLocalePreference); ok {
		if locale := hasPreference.PreferredLocale(); locale != "" {
			return locale
		}
	}

	return ""
}

// sendWithLocale executes a function with a specific locale context.
// This is a placeholder for locale handling - in a real implementation,
// this would set the application locale temporarily.
func (s *NotificationSender) sendWithLocale(locale string, fn func()) {
	// TODO: Implement actual locale switching when translation support is available
	// For now, just execute the function
	fn()
}

// sendToNotifiable sends a notification to a single notifiable entity via a specific channel.
func (s *NotificationSender) sendToNotifiable(
	notifiable any,
	notification contractsnotification.Notification,
	channelName string,
) error {
	// Check if notification should be sent
	if !s.shouldSendNotification(notifiable, notification, channelName) {
		return nil
	}

	// Get the channel
	channel, err := s.manager.Channel(channelName)
	if err != nil {
		s.dispatchFailedEvent(notifiable, notification, channelName, nil, err)
		// SendNow logs errors but doesn't return them, so we don't return here either.
		return nil
	}

	// Send the notification
	if err := channel.Send(notifiable, notification); err != nil {
		s.dispatchFailedEvent(notifiable, notification, channelName, nil, err)
		// SendNow logs errors but doesn't return them, so we don't return here either.
		return nil
	}

	// Fire sent event (no response from channel currently)
	s.dispatchSentEvent(notifiable, notification, channelName, nil)

	return nil
}

// shouldSendNotification determines if the notification should be sent.
// It checks the notification's ShouldSend method and dispatches the NotificationSending event.
func (s *NotificationSender) shouldSendNotification(
	notifiable any,
	notification contractsnotification.Notification,
	channelName string,
) bool {
	// Check if notification has ShouldSend method
	if shouldSender, ok := notification.(contractsnotification.ShouldSend); ok {
		if !shouldSender.ShouldSend(notifiable, channelName) {
			return false
		}
	}

	// Fire NotificationSending event and check if any listener returns false
	if s.manager.event != nil {
		sendingEvent := &events.NotificationSending{
			Notifiable:   notifiable,
			Notification: notification,
			Channel:      channelName,
		}

		// Use Until to allow listeners to cancel the notification
		result, err := s.manager.event.Until(sendingEvent)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Error during NotificationSending event: %v", err))
			return true // Continue sending on error
		}

		// If any listener returned false, cancel the notification
		if result == false {
			return false
		}
	}

	return true
}

// dispatchSentEvent dispatches the NotificationSent event synchronously.
func (s *NotificationSender) dispatchSentEvent(notifiable any, notification contractsnotification.Notification, channelName string, response any) {
	if s.manager.event != nil {
		sentEvent := &events.NotificationSent{
			Notifiable:   notifiable,
			Notification: notification,
			Channel:      channelName,
			Response:     response,
		}
		// Use synchronous Dispatch instead of queued Job
		if _, err := s.manager.event.Dispatch(sentEvent); err != nil {
			s.log.Warning(fmt.Sprintf("Error dispatching NotificationSent event: %v", err))
		}
	}
}

// dispatchFailedEvent dispatches the NotificationFailed event synchronously.
func (s *NotificationSender) dispatchFailedEvent(notifiable any, notification contractsnotification.Notification, channelName string, response any, notifErr error) {
	if s.manager.event != nil {
		failedEvent := &events.NotificationFailed{
			Notifiable:   notifiable,
			Notification: notification,
			Channel:      channelName,
			Error:        notifErr,
		}
		// Use synchronous Dispatch instead of queued Job
		if _, err := s.manager.event.Dispatch(failedEvent); err != nil {
			s.log.Warning(fmt.Sprintf("Error dispatching NotificationFailed event: %v", err))
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

		// Queue a job for each notifiable, passing sender reference
		job := NewSendQueuedNotificationJob(notifiable, notification, channels, s)
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
