package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rusmanplatd/goravelframework/contracts/mail"
	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
	mocksevent "github.com/rusmanplatd/goravelframework/mocks/event"
	mockslog "github.com/rusmanplatd/goravelframework/mocks/log"
	mocksnotification "github.com/rusmanplatd/goravelframework/mocks/notification"
	mocksqueue "github.com/rusmanplatd/goravelframework/mocks/queue"
)

// TestAnonymousNotifiable tests the anonymous notifiable functionality
func TestAnonymousNotifiable(t *testing.T) {
	t.Run("creates anonymous notifiable with routes", func(t *testing.T) {
		notifiable := NewAnonymousNotifiable().
			Route("mail", "user@example.com").
			Route("sms", "+1234567890")

		assert.Equal(t, "user@example.com", notifiable.RouteNotificationFor("mail"))
		assert.Equal(t, "+1234567890", notifiable.RouteNotificationFor("sms"))
	})

	t.Run("returns nil for unknown channel", func(t *testing.T) {
		notifiable := NewAnonymousNotifiable().
			Route("mail", "user@example.com")

		assert.Nil(t, notifiable.RouteNotificationFor("sms"))
	})

	t.Run("sends notification to anonymous notifiable", func(t *testing.T) {
		mockQueue := mocksqueue.NewQueue(t)
		mockLog := mockslog.NewLog(t)
		mockEvent := mocksevent.NewInstance(t)
		mockChannel := mocksnotification.NewChannel(t)

		manager := &ChannelManager{
			channels: map[string]contractsnotification.Channel{
				"mail": mockChannel,
			},
			event: mockEvent,
		}

		sender := NewNotificationSender(manager, mockQueue, mockLog)

		notifiable := NewAnonymousNotifiable().Route("mail", "user@example.com")
		notification := &testAnonymousNotification{}

		// Expect Until to be called for NotificationSending
		mockEvent.EXPECT().Until(mock.Anything).Return(nil, nil).Once()

		// Expect channel to send
		mockChannel.EXPECT().Send(notifiable, notification).Return(nil).Once()

		// Expect Dispatch to be called for NotificationSent
		mockEvent.EXPECT().Dispatch(mock.Anything).Return([]any{}, nil).Once()

		err := sender.SendNow(notifiable, notification)
		assert.NoError(t, err)
	})
}

// TestLocaleSupport tests locale handling in notifications
func TestLocaleSupport(t *testing.T) {
	t.Run("uses notification locale if available", func(t *testing.T) {
		mockQueue := mocksqueue.NewQueue(t)
		mockLog := mockslog.NewLog(t)
		mockEvent := mocksevent.NewInstance(t)
		mockChannel := mocksnotification.NewChannel(t)

		manager := &ChannelManager{
			channels: map[string]contractsnotification.Channel{
				"test": mockChannel,
			},
			event: mockEvent,
		}

		sender := NewNotificationSender(manager, mockQueue, mockLog)

		notifiable := &testNotifiableWithLocale{locale: "en"}
		notification := &testNotificationWithLocale{locale: "fr"}

		// Expect Until to be called
		mockEvent.EXPECT().Until(mock.Anything).Return(nil, nil).Once()

		// Expect channel to send
		mockChannel.EXPECT().Send(notifiable, notification).Return(nil).Once()

		// Expect Dispatch to be called
		mockEvent.EXPECT().Dispatch(mock.Anything).Return([]any{}, nil).Once()

		// Get preferred locale - should prefer notification locale
		locale := sender.getPreferredLocale(notifiable, notification)
		assert.Equal(t, "fr", locale)

		err := sender.SendNow(notifiable, notification)
		assert.NoError(t, err)
	})

	t.Run("uses notifiable locale if notification has no locale", func(t *testing.T) {
		mockQueue := mocksqueue.NewQueue(t)
		mockLog := mockslog.NewLog(t)
		mockEvent := mocksevent.NewInstance(t)

		manager := &ChannelManager{
			event: mockEvent,
		}

		sender := NewNotificationSender(manager, mockQueue, mockLog)

		notifiable := &testNotifiableWithLocale{locale: "es"}
		notification := &testAnonymousNotification{}

		locale := sender.getPreferredLocale(notifiable, notification)
		assert.Equal(t, "es", locale)
	})

	t.Run("returns empty string if no locale available", func(t *testing.T) {
		mockQueue := mocksqueue.NewQueue(t)
		mockLog := mockslog.NewLog(t)
		mockEvent := mocksevent.NewInstance(t)

		manager := &ChannelManager{
			event: mockEvent,
		}

		sender := NewNotificationSender(manager, mockQueue, mockLog)

		notifiable := "user@example.com"
		notification := &testAnonymousNotification{}

		locale := sender.getPreferredLocale(notifiable, notification)
		assert.Equal(t, "", locale)
	})
}

// testAnonymousNotification is a simple notification for testing
type testAnonymousNotification struct{}

func (n *testAnonymousNotification) Via(notifiable any) []string {
	return []string{"mail"}
}

func (n *testAnonymousNotification) ToDatabase(notifiable any) *contractsnotification.DatabaseMessage {
	return nil
}

func (n *testAnonymousNotification) ToMail(notifiable any) mail.Mailable {
	return nil
}

func (n *testAnonymousNotification) ToArray(notifiable any) map[string]any {
	return nil
}

// testNotificationWithLocale implements HasLocale
type testNotificationWithLocale struct {
	locale string
}

func (n *testNotificationWithLocale) Via(notifiable any) []string {
	return []string{"test"}
}

func (n *testNotificationWithLocale) ToDatabase(notifiable any) *contractsnotification.DatabaseMessage {
	return nil
}

func (n *testNotificationWithLocale) ToMail(notifiable any) mail.Mailable {
	return nil
}

func (n *testNotificationWithLocale) ToArray(notifiable any) map[string]any {
	return nil
}

func (n *testNotificationWithLocale) Locale() string {
	return n.locale
}

// testNotifiableWithLocale implements HasLocalePreference
type testNotifiableWithLocale struct {
	locale string
}

func (n *testNotifiableWithLocale) PreferredLocale() string {
	return n.locale
}
