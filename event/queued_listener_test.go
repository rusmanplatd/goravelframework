package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocksqueue "github.com/goravel/framework/mocks/queue"
)

// TestQueuedListener tests queued listener functionality
func TestQueuedListener(t *testing.T) {
	t.Run("listener is queued when ShouldQueue returns true", func(t *testing.T) {
		mockQueue := mocksqueue.NewQueue(t)
		mockPendingJob := mocksqueue.NewPendingJob(t)
		app := NewApplication(mockQueue)

		listener := &testQueuedListener{shouldQueue: true}

		// Expect Job to be called
		mockQueue.EXPECT().Job(mock.AnythingOfType("*event.QueuedListenerJob"), mock.Anything).
			Return(mockPendingJob).Once()

		// Expect Dispatch to be called
		mockPendingJob.EXPECT().Dispatch().Return(nil).Once()

		err := app.Listen("user.created", listener)
		assert.NoError(t, err)

		_, err = app.Dispatch("user.created", "john")
		assert.NoError(t, err)

		// Listener should not have been called synchronously
		assert.False(t, listener.called)
	})

	t.Run("listener executes synchronously when ShouldQueue returns false", func(t *testing.T) {
		mockQueue := mocksqueue.NewQueue(t)
		app := NewApplication(mockQueue)

		listener := &testQueuedListener{shouldQueue: false}

		err := app.Listen("user.updated", listener)
		assert.NoError(t, err)

		_, err = app.Dispatch("user.updated", "jane")
		assert.NoError(t, err)

		// Listener should have been called synchronously
		assert.True(t, listener.called)
	})

	t.Run("queued listener with custom queue configuration", func(t *testing.T) {
		mockQueue := mocksqueue.NewQueue(t)
		mockPendingJob := mocksqueue.NewPendingJob(t)
		app := NewApplication(mockQueue)

		listener := &testQueueableListener{
			shouldQueue: true,
			connection:  "redis",
			queueName:   "notifications",
			delay:       5,
		}

		// Expect Job to be called
		mockQueue.EXPECT().Job(mock.AnythingOfType("*event.QueuedListenerJob"), mock.Anything).
			Return(mockPendingJob).Once()

		// Expect queue configuration methods to be called
		mockPendingJob.EXPECT().OnConnection("redis").Return(mockPendingJob).Once()
		mockPendingJob.EXPECT().OnQueue("notifications").Return(mockPendingJob).Once()
		mockPendingJob.EXPECT().Delay(mock.AnythingOfType("time.Time")).Return(mockPendingJob).Once()
		mockPendingJob.EXPECT().Dispatch().Return(nil).Once()

		err := app.Listen("order.placed", listener)
		assert.NoError(t, err)

		_, err = app.Dispatch("order.placed", "order-123")
		assert.NoError(t, err)
	})

	t.Run("QueuedListenerJob executes listener when handled", func(t *testing.T) {
		listener := &testQueuedListener{shouldQueue: false}
		job := NewQueuedListenerJob(listener, "test.event", []any{"payload"})

		assert.Equal(t, "queued_listener:test.event", job.Signature())

		err := job.Handle()
		assert.NoError(t, err)
		assert.True(t, listener.called)
	})
}

// testQueuedListener is a test listener that implements ShouldQueue
type testQueuedListener struct {
	shouldQueue bool
	called      bool
}

func (l *testQueuedListener) ShouldQueue() bool {
	return l.shouldQueue
}

func (l *testQueuedListener) Handle(args ...any) error {
	l.called = true
	return nil
}

// testQueueableListener implements both ShouldQueue and QueueableListener
type testQueueableListener struct {
	shouldQueue bool
	connection  string
	queueName   string
	delay       int
	called      bool
}

func (l *testQueueableListener) ShouldQueue() bool {
	return l.shouldQueue
}

func (l *testQueueableListener) ViaConnection() string {
	return l.connection
}

func (l *testQueueableListener) ViaQueue() string {
	return l.queueName
}

func (l *testQueueableListener) WithDelay() int {
	return l.delay
}

func (l *testQueueableListener) Handle(args ...any) error {
	l.called = true
	return nil
}

// TestShouldQueueListener tests the shouldQueueListener helper
func TestShouldQueueListener(t *testing.T) {
	t.Run("returns true for listener implementing ShouldQueue", func(t *testing.T) {
		listener := &testQueuedListener{shouldQueue: true}
		assert.True(t, shouldQueueListener(listener, []any{}))
	})

	t.Run("returns false for listener not implementing ShouldQueue", func(t *testing.T) {
		listener := func() error { return nil }
		assert.False(t, shouldQueueListener(listener, []any{}))
	})

	t.Run("returns false when ShouldQueue returns false", func(t *testing.T) {
		listener := &testQueuedListener{shouldQueue: false}
		assert.False(t, shouldQueueListener(listener, []any{}))
	})
}

// TestGetListenerQueueConfig tests queue configuration extraction
func TestGetListenerQueueConfig(t *testing.T) {
	t.Run("extracts queue configuration from QueueableListener", func(t *testing.T) {
		listener := &testQueueableListener{
			connection: "redis",
			queueName:  "high-priority",
			delay:      10,
		}

		conn, queue, delay := getListenerQueueConfig(listener)
		assert.Equal(t, "redis", conn)
		assert.Equal(t, "high-priority", queue)
		assert.Equal(t, 10*time.Second, delay)
	})

	t.Run("returns empty values for non-queueable listener", func(t *testing.T) {
		listener := func() error { return nil }

		conn, queue, delay := getListenerQueueConfig(listener)
		assert.Equal(t, "", conn)
		assert.Equal(t, "", queue)
		assert.Equal(t, time.Duration(0), delay)
	})
}
