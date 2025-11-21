package broadcast

// Broadcaster defines the interface for broadcasting events to channels.
type Broadcaster interface {
	// Channel registers a channel authenticator callback.
	Channel(channel string, callback ChannelAuthCallback, options ...ChannelOption) Broadcaster

	// Broadcast sends an event to the specified channels.
	Broadcast(channels []Channel, event string, payload map[string]any) error

	// Auth authenticates the incoming request for channel access.
	Auth(request any) (any, error)

	// ValidAuthenticationResponse returns the valid authentication response.
	ValidAuthenticationResponse(request any, result any) (any, error)

	// ResolveAuthenticatedUser resolves the authenticated user payload for connection requests.
	ResolveAuthenticatedUser(request any) (map[string]any, error)

	// ResolveAuthenticatedUserUsing registers the user retrieval callback for authentication.
	ResolveAuthenticatedUserUsing(callback func(request any) map[string]any)
}

// Manager defines the interface for managing broadcast connections.
type Manager interface {
	// Connection gets a broadcaster instance by name.
	// If no name is provided, the default connection is returned.
	Connection(name ...string) (Broadcaster, error)

	// Driver is an alias for Connection.
	Driver(name ...string) (Broadcaster, error)
}

// Channel represents a broadcast channel.
type Channel interface {
	// GetName returns the channel name.
	GetName() string

	// IsPrivate returns true if the channel is private.
	IsPrivate() bool

	// IsPresence returns true if the channel is a presence channel.
	IsPresence() bool
}

// ChannelAuthCallback is a function that authenticates a user for a channel.
// It receives the authenticated user and any channel parameters.
// Returns true/false for simple authorization, or a map for presence channels.
type ChannelAuthCallback func(user any, params ...any) any

// ChannelOption defines options for channel registration.
type ChannelOption struct {
	// Guards specifies which authentication guards to use.
	Guards []string
}
