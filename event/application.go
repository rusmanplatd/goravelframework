package event

import (
	"slices"

	"github.com/rusmanplatd/goravelframework/contracts/event"
	"github.com/rusmanplatd/goravelframework/contracts/queue"
)

type Application struct {
	events         map[event.Event][]event.Listener
	listeners      map[string][]any // string event name -> listeners
	wildcards      map[string][]any // wildcard patterns -> listeners
	wildcardsCache map[string][]any // cached prepared wildcard listeners per event
	pushedEvents   map[string][]any // pushed events -> payloads
	queue          queue.Queue
}

func NewApplication(queue queue.Queue) *Application {
	return &Application{
		events:         make(map[event.Event][]event.Listener),
		listeners:      make(map[string][]any),
		wildcards:      make(map[string][]any),
		wildcardsCache: make(map[string][]any),
		pushedEvents:   make(map[string][]any),
		queue:          queue,
	}
}

func (app *Application) Register(events map[event.Event][]event.Listener) {
	var (
		jobs     []queue.Job
		jobNames []string
	)

	if app.events == nil {
		app.events = map[event.Event][]event.Listener{}
	}

	for e, listeners := range events {
		app.events[e] = listeners
		for _, listener := range listeners {
			if !slices.Contains(jobNames, listener.Signature()) {
				jobs = append(jobs, listener)
				jobNames = append(jobNames, listener.Signature())
			}
		}
	}

	app.queue.Register(jobs)
}

func (app *Application) GetEvents() map[event.Event][]event.Listener {
	return app.events
}

func (app *Application) Job(e event.Event, args []event.Arg) event.Task {
	listeners, ok := app.events[e]
	if !ok {
		listeners = make([]event.Listener, 0)
	}

	return NewTask(app.queue, args, e, listeners)
}
