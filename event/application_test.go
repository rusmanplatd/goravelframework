package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rusmanplatd/goravelframework/contracts/event"
	"github.com/rusmanplatd/goravelframework/contracts/queue"
	mocksevent "github.com/rusmanplatd/goravelframework/mocks/event"
	mocksqueue "github.com/rusmanplatd/goravelframework/mocks/queue"
)

func TestApplication_Register(t *testing.T) {
	var (
		mockQueue *mocksqueue.Queue
	)

	tests := []struct {
		name   string
		events func() map[event.Event][]event.Listener
	}{
		{
			name: "MultipleEvents",
			events: func() map[event.Event][]event.Listener {
				event1 := mocksevent.NewEvent(t)
				event2 := mocksevent.NewEvent(t)
				listener1 := mocksevent.NewListener(t)
				listener1.EXPECT().Signature().Return("listener1").Twice()
				listener2 := mocksevent.NewListener(t)
				listener2.EXPECT().Signature().Return("listener2").Times(3)

				mockQueue.EXPECT().Register(mock.MatchedBy(func(listeners []queue.Job) bool {
					return assert.ElementsMatch(t, []queue.Job{
						listener1,
						listener2,
					}, listeners)
				})).Once()

				return map[event.Event][]event.Listener{
					event1: {
						listener1,
						listener2,
					},
					event2: {
						listener2,
					},
				}
			},
		},
		{
			name: "NoEvents",
			events: func() map[event.Event][]event.Listener {
				mockQueue.EXPECT().Register([]queue.Job(nil)).Once()

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueue = mocksqueue.NewQueue(t)
			app := NewApplication(mockQueue)

			events := tt.events()
			app.Register(events)

			assert.Equal(t, len(events), len(app.GetEvents()))
			for e, listeners := range events {
				assert.ElementsMatch(t, listeners, app.GetEvents()[e])
			}
		})
	}
}
