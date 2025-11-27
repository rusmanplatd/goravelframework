package event

type Instance interface {
	// Register event listeners to the application.
	Register(map[Event][]Listener)
	// Job create a new event task.
	Job(event Event, args []Arg) Task
	// GetEvents gets all registered events.
	GetEvents() map[Event][]Listener
	// Listen registers an event listener with the dispatcher.
	// The event parameter can be a string event name, Event interface, or wildcard pattern (e.g., "user.*").
	// The listeners parameter can be closures, listener structs, or listener names.
	Listen(event any, listeners ...any) error
	// HasListeners determines if a given event has listeners.
	HasListeners(event string) bool
	// Dispatch fires an event and calls the listeners synchronously.
	// Returns a slice of responses from all listeners.
	Dispatch(event any, payload ...any) ([]any, error)
	// Until dispatches an event until the first non-null response is returned.
	Until(event any, payload ...any) (any, error)
	// Subscribe registers an event subscriber with the dispatcher.
	Subscribe(subscriber Subscriber) error
	// Forget removes a set of listeners from the dispatcher.
	Forget(event string)
	// Push registers an event and payload to be fired later.
	Push(event string, payload ...any)
	// Flush flushes a set of pushed events.
	Flush(event string) error
}

type Event interface {
	// Handle the event.
	Handle(args []Arg) ([]Arg, error)
}

type Listener interface {
	// Signature returns the unique identifier for the listener.
	Signature() string
	// Queue configure the event queue options.
	Queue(args ...any) Queue
	// Handle the event.
	Handle(args ...any) error
}

type Task interface {
	// Dispatch an event and call the listeners.
	Dispatch() error
}

type Arg struct {
	Value any
	Type  string
}

type Queue struct {
	Connection string
	Queue      string
	Enable     bool
}

// Subscriber represents an event subscriber that can subscribe to multiple events.
type Subscriber interface {
	// Subscribe returns a map of events to listeners.
	// The keys can be event names (string) or Event interfaces.
	// The values can be listener method names (string), closures, or Listener interfaces.
	Subscribe(dispatcher Instance) map[any][]any
}

// ShouldQueue indicates that a listener should be queued.
// When a listener implements this interface, it will be executed asynchronously via the queue system.
type ShouldQueue interface {
	// ShouldQueue returns true if the listener should be queued.
	ShouldQueue() bool
}

// QueueableListener provides queue configuration for a listener.
// Listeners can implement this interface to customize queue behavior.
type QueueableListener interface {
	// ViaConnection returns the queue connection name to use.
	ViaConnection() string
	// ViaQueue returns the queue name to use.
	ViaQueue() string
	// WithDelay returns the delay before the job should be processed.
	WithDelay() int
}

// ShouldBroadcast indicates that an event should be broadcast.
// When an event implements this interface, it will be broadcast to connected clients.
type ShouldBroadcast interface {
	// BroadcastOn returns the channels the event should broadcast on.
	BroadcastOn() []string
	// BroadcastWhen returns true if the event should be broadcast.
	BroadcastWhen() bool
}

// ShouldDispatchAfterCommit indicates that an event should be dispatched after database transaction commits.
// This is useful for events that should only fire if the database transaction succeeds.
type ShouldDispatchAfterCommit interface {
	// ShouldDispatchAfterCommit returns true if the event should wait for transaction commit.
	ShouldDispatchAfterCommit() bool
}
