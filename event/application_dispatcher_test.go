package event

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goravel/framework/contracts/event"
	mocksevent "github.com/goravel/framework/mocks/event"
	mocksqueue "github.com/goravel/framework/mocks/queue"
)

// TestListen tests the Listen method
func TestListen(t *testing.T) {
	mockQueue := mocksqueue.NewQueue(t)
	app := NewApplication(mockQueue)

	// Test listening to string event
	err := app.Listen("user.created", func(name string, user any) error {
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, app.HasListeners("user.created"))

	// Test listening to wildcard event
	err = app.Listen("notification.*", func(name string, payload any) error {
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, app.HasListeners("notification.sent"))
	assert.True(t, app.HasListeners("notification.failed"))

	// Test multiple listeners for same event
	err = app.Listen("order.placed", func() error { return nil }, func() error { return nil })
	assert.NoError(t, err)
	listeners := app.getListenersForEvent("order.placed")
	assert.Len(t, listeners, 2)

	// Test error cases
	err = app.Listen(nil, func() error { return nil })
	assert.Error(t, err)

	err = app.Listen("test.event")
	assert.Error(t, err)
}

// TestDispatch tests the Dispatch method
func TestDispatch_Sync(t *testing.T) {
	mockQueue := mocksqueue.NewQueue(t)
	app := NewApplication(mockQueue)

	// Test successful dispatch
	called := false
	err := app.Listen("user.registered", func(user string) error {
		called = true
		assert.Equal(t, "john", user)
		return nil
	})
	assert.NoError(t, err)

	responses, err := app.Dispatch("user.registered", "john")
	assert.NoError(t, err)
	assert.True(t, called)
	assert.NotNil(t, responses)

	// Test dispatch with multiple listeners
	counter := 0
	err = app.Listen("order.created", func() error {
		counter++
		return nil
	}, func() error {
		counter++
		return nil
	})
	assert.NoError(t, err)

	_, err = app.Dispatch("order.created")
	assert.NoError(t, err)
	assert.Equal(t, 2, counter)

	// Test dispatch with error
	err = app.Listen("error.event", func() error {
		return fmt.Errorf("listener error")
	})
	assert.NoError(t, err)

	_, err = app.Dispatch("error.event")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "listener error")
}

// TestDispatch_WithEventObject tests dispatching Event interface objects
func TestDispatch_WithEventObject(t *testing.T) {
	mockQueue := mocksqueue.NewQueue(t)
	app := NewApplication(mockQueue)

	mockEvent := mocksevent.NewEvent(t)
	mockEvent.EXPECT().Handle([]event.Arg{}).Return([]event.Arg{}, nil).Maybe()

	called := false
	err := app.Listen("Event", func(evt event.Event) error {
		called = true
		return nil
	})
	assert.NoError(t, err)

	_, err = app.Dispatch(mockEvent)
	assert.NoError(t, err)
	assert.True(t, called)
}

// TestUntil tests the Until method
func TestUntil(t *testing.T) {
	mockQueue := mocksqueue.NewQueue(t)
	app := NewApplication(mockQueue)

	// Test until with first non-nil response
	err := app.Listen("check.permission",
		func() (bool, error) { return false, nil },
		func() (bool, error) { return true, nil },
		func() (bool, error) { return false, nil },
	)
	assert.NoError(t, err)

	result, err := app.Until("check.permission")
	assert.NoError(t, err)
	assert.Equal(t, false, result)

	// Test until with all nil responses
	err = app.Listen("no.response", func() error { return nil })
	assert.NoError(t, err)

	result, err = app.Until("no.response")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// TestWildcardListeners tests wildcard pattern matching
func TestWildcardListeners(t *testing.T) {
	mockQueue := mocksqueue.NewQueue(t)
	app := NewApplication(mockQueue)

	counter := 0
	err := app.Listen("user.*", func() error {
		counter++
		return nil
	})
	assert.NoError(t, err)

	// Test various user events
	_, err = app.Dispatch("user.created")
	assert.NoError(t, err)
	assert.Equal(t, 1, counter)

	_, err = app.Dispatch("user.updated")
	assert.NoError(t, err)
	assert.Equal(t, 2, counter)

	_, err = app.Dispatch("user.deleted")
	assert.NoError(t, err)
	assert.Equal(t, 3, counter)

	// Non-matching event should not trigger
	_, err = app.Dispatch("order.created")
	assert.NoError(t, err)
	assert.Equal(t, 3, counter)
}

// TestSubscribe tests the Subscribe method
func TestSubscribe(t *testing.T) {
	mockQueue := mocksqueue.NewQueue(t)
	app := NewApplication(mockQueue)

	subscriber := &testSubscriber{}
	err := app.Subscribe(subscriber)
	assert.NoError(t, err)

	// Verify listeners were registered
	assert.True(t, app.HasListeners("user.created"))
	assert.True(t, app.HasListeners("user.updated"))

	// Test dispatching events
	_, err = app.Dispatch("user.created", "john")
	assert.NoError(t, err)
	assert.Equal(t, 1, subscriber.createdCount)

	_, err = app.Dispatch("user.updated", "jane")
	assert.NoError(t, err)
	assert.Equal(t, 1, subscriber.updatedCount)
}

// TestForget tests the Forget method
func TestForget(t *testing.T) {
	mockQueue := mocksqueue.NewQueue(t)
	app := NewApplication(mockQueue)

	err := app.Listen("test.event", func() error { return nil })
	assert.NoError(t, err)
	assert.True(t, app.HasListeners("test.event"))

	app.Forget("test.event")
	assert.False(t, app.HasListeners("test.event"))

	// Test forgetting wildcard
	err = app.Listen("user.*", func() error { return nil })
	assert.NoError(t, err)
	assert.True(t, app.HasListeners("user.created"))

	app.Forget("user.*")
	assert.False(t, app.HasListeners("user.created"))
}

// TestPushAndFlush tests the Push and Flush methods
func TestPushAndFlush(t *testing.T) {
	mockQueue := mocksqueue.NewQueue(t)
	app := NewApplication(mockQueue)

	counter := 0
	err := app.Listen("deferred.event", func(value int) error {
		counter += value
		return nil
	})
	assert.NoError(t, err)

	// Push events
	app.Push("deferred.event", 1)
	app.Push("deferred.event", 2)
	app.Push("deferred.event", 3)

	// Counter should still be 0
	assert.Equal(t, 0, counter)

	// Flush events
	err = app.Flush("deferred.event")
	assert.NoError(t, err)
	assert.Equal(t, 6, counter)

	// Flushing again should do nothing
	err = app.Flush("deferred.event")
	assert.NoError(t, err)
	assert.Equal(t, 6, counter)
}

// testSubscriber is a test implementation of event.Subscriber
type testSubscriber struct {
	createdCount int
	updatedCount int
}

func (s *testSubscriber) Subscribe(dispatcher event.Instance) map[any][]any {
	return map[any][]any{
		"user.created": {
			func(name string) error {
				s.createdCount++
				return nil
			},
		},
		"user.updated": {
			func(name string) error {
				s.updatedCount++
				return nil
			},
		},
	}
}
