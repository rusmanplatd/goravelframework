package broadcast

// ShouldBroadcast defines an event that should be broadcast.
type ShouldBroadcast interface {
	// BroadcastOn returns the channels the event should broadcast on.
	BroadcastOn() []Channel
}

// ShouldBroadcastNow defines an event that should be broadcast synchronously.
// Events implementing this interface will be broadcast immediately instead of being queued.
type ShouldBroadcastNow interface {
	ShouldBroadcast
}

// BroadcastAs allows an event to customize its broadcast name.
type BroadcastAs interface {
	// BroadcastAs returns the name of the event for broadcasting.
	// If not implemented, the event's struct name will be used.
	BroadcastAs() string
}

// BroadcastWith allows an event to customize its broadcast payload.
type BroadcastWith interface {
	// BroadcastWith returns the data to broadcast with the event.
	// If not implemented, all public fields will be broadcast.
	BroadcastWith() map[string]any
}

// BroadcastConnections allows an event to specify which broadcast connections to use.
type BroadcastConnections interface {
	// BroadcastConnections returns the broadcast connections to use.
	// If not implemented, the default connection will be used.
	BroadcastConnections() []string
}

// BroadcastQueue allows an event to customize its queue configuration.
type BroadcastQueue interface {
	// BroadcastQueue returns the queue name to use for broadcasting.
	BroadcastQueue() string
}

// BroadcastConnection allows an event to customize its queue connection.
type BroadcastConnection interface {
	// BroadcastConnection returns the queue connection to use for broadcasting.
	BroadcastConnection() string
}

// HasBroadcastChannel allows a model to define its broadcast channel.
type HasBroadcastChannel interface {
	// BroadcastChannelRoute returns the route for the model's broadcast channel.
	BroadcastChannelRoute() string
}

// BroadcastVia allows an event to specify which broadcaster to use.
type BroadcastVia interface {
	// BroadcastVia returns the broadcaster connection name to use.
	BroadcastVia() string
}

// DontBroadcastToCurrentUser allows an event to exclude the current user from recipients.
type DontBroadcastToCurrentUser interface {
	// DontBroadcastToCurrentUser excludes current user from broadcast.
	DontBroadcastToCurrentUser()
}

// InteractsWithBroadcasting provides common broadcasting functionality for events.
type InteractsWithBroadcasting interface {
	// BroadcastVia sets the broadcaster connection to use.
	BroadcastVia(connection string) InteractsWithBroadcasting

	// BroadcastConnections returns the broadcaster connections to use.
	BroadcastConnections() []string
}
