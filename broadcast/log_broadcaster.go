package broadcast

import (
	"encoding/json"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
	contractslog "github.com/goravel/framework/contracts/log"
)

// Ensure interface implementation
var _ contractsbroadcast.Broadcaster = (*Log)(nil)

// Log is a broadcaster that logs events instead of broadcasting them.
type Log struct {
	*BaseBroadcaster
	log contractslog.Log
}

// NewLog creates a new log broadcaster.
func NewLog(log contractslog.Log) *Log {
	return &Log{
		BaseBroadcaster: NewBaseBroadcaster(),
		log:             log,
	}
}

// Channel registers a channel authenticator callback.
func (l *Log) Channel(channel string, callback contractsbroadcast.ChannelAuthCallback, options ...contractsbroadcast.ChannelOption) contractsbroadcast.Broadcaster {
	l.BaseBroadcaster.Channel(channel, callback, options...)
	return l
}

// Broadcast logs the event instead of broadcasting it.
func (l *Log) Broadcast(channels []contractsbroadcast.Channel, event string, payload map[string]any) error {
	channelNames := l.formatChannels(channels)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		payloadJSON = []byte("{}")
	}

	l.log.Info("Broadcasting event", map[string]any{
		"event":    event,
		"channels": channelNames,
		"payload":  string(payloadJSON),
	})

	return nil
}

// Auth authenticates the incoming request for channel access.
func (l *Log) Auth(request any) (any, error) {
	// For log broadcaster, we'll allow all access and log the attempt
	l.log.Info("Broadcast authentication attempt", map[string]any{
		"request": request,
	})
	return true, nil
}

// ValidAuthenticationResponse returns the valid authentication response for log broadcaster.
func (l *Log) ValidAuthenticationResponse(request any, result any) (any, error) {
	// Log the authentication response
	l.log.Info("Broadcast authentication response", map[string]any{
		"request": request,
		"result":  result,
	})
	return true, nil
}

// ResolveAuthenticatedUser resolves the authenticated user payload for connection requests.
func (l *Log) ResolveAuthenticatedUser(request any) (map[string]any, error) {
	return l.BaseBroadcaster.ResolveAuthenticatedUser(request)
}

// ResolveAuthenticatedUserUsing registers the user retrieval callback for authentication.
func (l *Log) ResolveAuthenticatedUserUsing(callback func(request any) map[string]any) {
	l.BaseBroadcaster.ResolveAuthenticatedUserUsing(callback)
}
