package broadcast

import (
	"fmt"
	"reflect"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
	contractsqueue "github.com/goravel/framework/contracts/queue"
)

// BroadcastEvent wraps an event for queued broadcasting.
type BroadcastEvent struct {
	event contractsbroadcast.ShouldBroadcast
}

// NewBroadcastEvent creates a new broadcast event wrapper.
func NewBroadcastEvent(event contractsbroadcast.ShouldBroadcast) *BroadcastEvent {
	return &BroadcastEvent{
		event: event,
	}
}

// Signature returns the unique identifier for the job.
func (b *BroadcastEvent) Signature() string {
	return fmt.Sprintf("broadcast:%s", reflect.TypeOf(b.event).String())
}

// Handle executes the broadcast.
func (b *BroadcastEvent) Handle(args ...any) error {
	// Extract manager from args
	if len(args) == 0 {
		return fmt.Errorf("broadcast manager not provided")
	}

	manager, ok := args[0].(contractsbroadcast.Manager)
	if !ok {
		return fmt.Errorf("first argument must be broadcast manager")
	}

	// Get event name
	eventName := b.getEventName()

	// Get channels
	channels := b.event.BroadcastOn()
	if len(channels) == 0 {
		return nil
	}

	// Get connections to use
	connections := b.getConnections()

	// Get payload
	payload := b.getPayload()

	// Broadcast on each connection
	for _, connection := range connections {
		broadcaster, err := manager.Connection(connection)
		if err != nil {
			return fmt.Errorf("failed to get broadcaster for connection %s: %w", connection, err)
		}

		channelsForConnection := b.getChannelsForConnection(channels, connection)
		payloadForConnection := b.getPayloadForConnection(payload, connection)

		if err := broadcaster.Broadcast(channelsForConnection, eventName, payloadForConnection); err != nil {
			return fmt.Errorf("failed to broadcast on connection %s: %w", connection, err)
		}
	}

	return nil
}

// getEventName returns the event name for broadcasting.
func (b *BroadcastEvent) getEventName() string {
	if named, ok := b.event.(contractsbroadcast.BroadcastAs); ok {
		return named.BroadcastAs()
	}

	// Use the struct type name
	eventType := reflect.TypeOf(b.event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}

	return eventType.Name()
}

// getConnections returns the connections to broadcast on.
func (b *BroadcastEvent) getConnections() []string {
	if conns, ok := b.event.(contractsbroadcast.BroadcastConnections); ok {
		connections := conns.BroadcastConnections()
		if len(connections) > 0 {
			return connections
		}
	}

	return []string{""}
}

// getPayload returns the broadcast payload.
func (b *BroadcastEvent) getPayload() map[string]any {
	if withPayload, ok := b.event.(contractsbroadcast.BroadcastWith); ok {
		return withPayload.BroadcastWith()
	}

	// Use reflection to get public fields
	payload := make(map[string]any)
	eventValue := reflect.ValueOf(b.event)
	if eventValue.Kind() == reflect.Ptr {
		eventValue = eventValue.Elem()
	}

	eventType := eventValue.Type()
	for i := 0; i < eventValue.NumField(); i++ {
		field := eventType.Field(i)
		if field.IsExported() {
			payload[field.Name] = eventValue.Field(i).Interface()
		}
	}

	return payload
}

// getChannelsForConnection returns channels for a specific connection.
func (b *BroadcastEvent) getChannelsForConnection(channels []contractsbroadcast.Channel, connection string) []contractsbroadcast.Channel {
	// For now, return all channels for all connections
	// This can be enhanced to support per-connection channel filtering
	return channels
}

// getPayloadForConnection returns payload for a specific connection.
func (b *BroadcastEvent) getPayloadForConnection(payload map[string]any, connection string) map[string]any {
	// For now, return the same payload for all connections
	// This can be enhanced to support per-connection payload customization
	return payload
}

// Queue returns queue configuration for the broadcast job.
func (b *BroadcastEvent) Queue(args ...any) contractsqueue.Args {
	queueArgs := contractsqueue.Args{
		Connection: "",
		Queue:      "",
	}

	if queueConfig, ok := b.event.(contractsbroadcast.BroadcastQueue); ok {
		queueArgs.Queue = queueConfig.BroadcastQueue()
	}

	if connConfig, ok := b.event.(contractsbroadcast.BroadcastConnection); ok {
		queueArgs.Connection = connConfig.BroadcastConnection()
	}

	return queueArgs
}
