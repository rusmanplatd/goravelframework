package broadcast

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	contractsbroadcast "github.com/rusmanplatd/goravelframework/contracts/broadcast"
)

// Ensure interface implementation
var _ contractsbroadcast.Broadcaster = (*Redis)(nil)

// Redis is a broadcaster that uses Redis pub/sub.
type Redis struct {
	*BaseBroadcaster
	redis      any // Redis instance (uses reflection to call methods)
	connection string
	prefix     string
}

// NewRedis creates a new Redis broadcaster.
func NewRedis(redis any, connection string, prefix string) *Redis {
	return &Redis{
		BaseBroadcaster: NewBaseBroadcaster(),
		redis:           redis,
		connection:      connection,
		prefix:          prefix,
	}
}

// Channel registers a channel authenticator callback.
func (r *Redis) Channel(channel string, callback contractsbroadcast.ChannelAuthCallback, options ...contractsbroadcast.ChannelOption) contractsbroadcast.Broadcaster {
	r.BaseBroadcaster.Channel(channel, callback, options...)
	return r
}

// Broadcast publishes the event to Redis channels.
func (r *Redis) Broadcast(channels []contractsbroadcast.Channel, event string, payload map[string]any) error {
	if len(channels) == 0 {
		return nil
	}

	channelNames := r.formatChannels(channels)

	// Extract socket ID if present
	socket, _ := payload["socket"].(string)
	delete(payload, "socket")

	// Create broadcast payload
	broadcastPayload := map[string]any{
		"event":  event,
		"data":   payload,
		"socket": socket,
	}

	payloadJSON, err := json.Marshal(broadcastPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Use reflection to call redis.Connection(r.connection)
	redisValue := reflect.ValueOf(r.redis)
	connectionMethod := redisValue.MethodByName("Connection")
	if !connectionMethod.IsValid() {
		return fmt.Errorf("redis instance does not have Connection method")
	}

	connectionResults := connectionMethod.Call([]reflect.Value{reflect.ValueOf(r.connection)})
	if len(connectionResults) == 0 {
		return fmt.Errorf("Connection method did not return a value")
	}

	connection := connectionResults[0].Interface()
	ctx := context.Background()

	// Use reflection to call connection.Publish(ctx, channel, payload)
	for _, channel := range channelNames {
		channelName := r.prefix + channel

		connValue := reflect.ValueOf(connection)
		publishMethod := connValue.MethodByName("Publish")
		if !publishMethod.IsValid() {
			return fmt.Errorf("connection does not have Publish method")
		}

		results := publishMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(channelName),
			reflect.ValueOf(string(payloadJSON)),
		})

		// Check if there was an error (last return value)
		if len(results) > 0 && !results[len(results)-1].IsNil() {
			err := results[len(results)-1].Interface().(error)
			return fmt.Errorf("failed to publish to channel %s: %w", channelName, err)
		}
	}

	return nil
}

// Auth authenticates the incoming request for channel access.
func (r *Redis) Auth(request any) (any, error) {
	// This would need to be implemented based on the HTTP request structure
	// For now, return the interface compliance
	return nil, fmt.Errorf("auth method requires HTTP request implementation")
}

// ValidAuthenticationResponse returns the valid authentication response for Redis.
func (r *Redis) ValidAuthenticationResponse(request any, result any) (any, error) {
	// This would need to be implemented based on the HTTP request structure
	// For now, return the interface compliance
	return nil, fmt.Errorf("valid authentication response requires HTTP request implementation")
}

// ResolveAuthenticatedUser resolves the authenticated user payload for connection requests.
func (r *Redis) ResolveAuthenticatedUser(request any) (map[string]any, error) {
	return r.BaseBroadcaster.ResolveAuthenticatedUser(request)
}

// ResolveAuthenticatedUserUsing registers the user retrieval callback for authentication.
func (r *Redis) ResolveAuthenticatedUserUsing(callback func(request any) map[string]any) {
	r.BaseBroadcaster.ResolveAuthenticatedUserUsing(callback)
}
