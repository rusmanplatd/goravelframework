package notification

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rusmanplatd/goravelframework/contracts/mail"
	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
	mocksevent "github.com/rusmanplatd/goravelframework/mocks/event"
	mockslog "github.com/rusmanplatd/goravelframework/mocks/log"
	mocksnotification "github.com/rusmanplatd/goravelframework/mocks/notification"
	mocksqueue "github.com/rusmanplatd/goravelframework/mocks/queue"
	"github.com/rusmanplatd/goravelframework/notification/events"
)

// TestNotificationSenderWithEventIntegration tests the notification sender with the new event system
func TestNotificationSenderWithEventIntegration(t *testing.T) {
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

	// Test 1: Successful notification with event dispatching
	t.Run("successful notification dispatches events", func(t *testing.T) {
		notification := mocksnotification.NewNotification(t)
		notification.EXPECT().Via(mock.Anything).Return([]string{"test"})

		// Expect Until to be called for NotificationSending
		mockEvent.EXPECT().Until(mock.MatchedBy(func(evt any) bool {
			_, ok := evt.(*events.NotificationSending)
			return ok
		})).Return(nil, nil).Once()

		// Expect channel to send
		mockChannel.EXPECT().Send(mock.Anything, notification).Return(nil).Once()

		// Expect Dispatch to be called for NotificationSent
		mockEvent.EXPECT().Dispatch(mock.MatchedBy(func(evt any) bool {
			_, ok := evt.(*events.NotificationSent)
			return ok
		})).Return([]any{}, nil).Once()

		err := sender.SendNow("user@example.com", notification)
		assert.NoError(t, err)
	})

	// Test 2: Event cancellation via Until returning false
	t.Run("notification cancelled by event listener", func(t *testing.T) {
		notification := mocksnotification.NewNotification(t)
		notification.EXPECT().Via(mock.Anything).Return([]string{"test"})

		// Until returns false to cancel
		mockEvent.EXPECT().Until(mock.MatchedBy(func(evt any) bool {
			_, ok := evt.(*events.NotificationSending)
			return ok
		})).Return(false, nil).Once()

		// Channel should NOT be called
		// Sent event should NOT be dispatched

		err := sender.SendNow("user@example.com", notification)
		assert.NoError(t, err)
	})

	// Test 3: ShouldSend callback cancels notification
	t.Run("notification cancelled by ShouldSend callback", func(t *testing.T) {
		notification := &testNotificationWithShouldSend{
			channels:   []string{"test"},
			shouldSend: false,
		}

		// Until should NOT be called because ShouldSend returned false
		// Channel should NOT be called
		// Sent event should NOT be dispatched

		err := sender.SendNow("user@example.com", notification)
		assert.NoError(t, err)
	})

	// Test 4: Failed notification dispatches failed event
	t.Run("failed notification dispatches failed event", func(t *testing.T) {
		notification := mocksnotification.NewNotification(t)
		notification.EXPECT().Via(mock.Anything).Return([]string{"test"})

		mockEvent.EXPECT().Until(mock.Anything).Return(nil, nil).Once()

		// Channel returns error
		sendErr := fmt.Errorf("send failed")
		mockChannel.EXPECT().Send(mock.Anything, notification).Return(sendErr).Once()

		// Expect log error to be called
		mockLog.EXPECT().Error(mock.Anything).Maybe()

		// Expect Dispatch to be called for NotificationFailed
		mockEvent.EXPECT().Dispatch(mock.MatchedBy(func(evt any) bool {
			failedEvt, ok := evt.(*events.NotificationFailed)
			return ok && failedEvt.Error == sendErr
		})).Return([]any{}, nil).Once()

		err := sender.SendNow("user@example.com", notification)
		assert.NoError(t, err) // SendNow logs errors but doesn't return them
	})
}

// testNotificationWithShouldSend is a test notification that implements ShouldSend
type testNotificationWithShouldSend struct {
	channels   []string
	shouldSend bool
}

func (n *testNotificationWithShouldSend) Via(notifiable any) []string {
	return n.channels
}

func (n *testNotificationWithShouldSend) ToDatabase(notifiable any) *contractsnotification.DatabaseMessage {
	return nil
}

func (n *testNotificationWithShouldSend) ToMail(notifiable any) mail.Mailable {
	return nil
}

func (n *testNotificationWithShouldSend) ToArray(notifiable any) map[string]any {
	return nil
}

func (n *testNotificationWithShouldSend) ShouldSend(notifiable any, channel string) bool {
	return n.shouldSend
}
