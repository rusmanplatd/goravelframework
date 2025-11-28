package broadcast

import (
	contractsbroadcast "github.com/rusmanplatd/goravelframework/contracts/broadcast"
)

// Ensure interface implementation
var _ contractsbroadcast.Broadcaster = (*Null)(nil)

// Null is a no-op broadcaster for testing or disabled broadcasting.
type Null struct {
	*BaseBroadcaster
}

// NewNull creates a new null broadcaster.
func NewNull() *Null {
	return &Null{
		BaseBroadcaster: NewBaseBroadcaster(),
	}
}

// Channel registers a channel authenticator callback.
func (n *Null) Channel(channel string, callback contractsbroadcast.ChannelAuthCallback, options ...contractsbroadcast.ChannelOption) contractsbroadcast.Broadcaster {
	n.BaseBroadcaster.Channel(channel, callback, options...)
	return n
}

// Broadcast does nothing (no-op).
func (n *Null) Broadcast(channels []contractsbroadcast.Channel, event string, payload map[string]any) error {
	// No-op: do nothing
	return nil
}

// Auth authenticates the incoming request for channel access.
func (n *Null) Auth(request any) (any, error) {
	// No-op: always allow access for null broadcaster
	return true, nil
}

// ValidAuthenticationResponse returns the valid authentication response for null broadcaster.
func (n *Null) ValidAuthenticationResponse(request any, result any) (any, error) {
	// No-op: return simple true response
	return true, nil
}

// ResolveAuthenticatedUser resolves the authenticated user payload for connection requests.
func (n *Null) ResolveAuthenticatedUser(request any) (map[string]any, error) {
	return n.BaseBroadcaster.ResolveAuthenticatedUser(request)
}

// ResolveAuthenticatedUserUsing registers the user retrieval callback for authentication.
func (n *Null) ResolveAuthenticatedUserUsing(callback func(request any) map[string]any) {
	n.BaseBroadcaster.ResolveAuthenticatedUserUsing(callback)
}
