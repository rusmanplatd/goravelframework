# Broadcast

The Goravel broadcast feature enables real-time event broadcasting to various channels using multiple drivers, providing a comprehensive solution similar to Laravel's broadcasting system.

## Installation

The broadcast feature is included in Goravel by default. No additional installation is required.

## Configuration

Broadcasting configuration is located in `config/broadcasting.go`. You can configure the default driver and connection settings:

```go
// config/broadcasting.go
func Broadcasting() map[string]any {
    return map[string]any{
        "default": "BROADCAST_DRIVER", // null, log, redis, ably, pusher, reverb

        "connections": map[string]any{
            "reverb": map[string]any{
                "driver": "reverb",
                "key":    "REVERB_APP_KEY",
                "secret": "REVERB_APP_SECRET",
                "app_id": "REVERB_APP_ID",
                "options": map[string]any{
                    "host":   "REVERB_HOST",
                    "port":   "REVERB_PORT",
                    "scheme": "REVERB_SCHEME",
                    "useTLS": true,
                },
            },

            "pusher": map[string]any{
                "driver": "pusher",
                "key":    "PUSHER_APP_KEY",
                "secret": "PUSHER_APP_SECRET",
                "app_id": "PUSHER_APP_ID",
                "options": map[string]any{
                    "cluster":   "PUSHER_APP_CLUSTER",
                    "host":      "PUSHER_HOST",
                    "port":      "PUSHER_PORT",
                    "scheme":    "PUSHER_SCHEME",
                    "encrypted": true,
                },
            },

            "ably": map[string]any{
                "driver": "ably",
                "key":    "ABLY_KEY",
            },

            "redis": map[string]any{
                "driver":     "redis",
                "connection": "BROADCAST_REDIS_CONNECTION",
                "prefix":     "BROADCAST_REDIS_PREFIX",
            },

            "log": map[string]any{
                "driver": "log",
            },

            "null": map[string]any{
                "driver": "null",
            },
        },
    }
}
```

### Environment Variables

Set your broadcast driver in `.env`:

```env
BROADCAST_DRIVER=null  # Options: null, log, redis, ably, pusher, reverb
```

#### For Redis

```env
BROADCAST_DRIVER=redis
BROADCAST_REDIS_CONNECTION=default
BROADCAST_REDIS_PREFIX=broadcast:
```

#### For Ably

```env
BROADCAST_DRIVER=ably
ABLY_KEY=your-ably-api-key:your-ably-secret
```

#### For Pusher (requires `github.com/pusher/pusher-http-go/v5`)

```env
BROADCAST_DRIVER=pusher
PUSHER_APP_KEY=your-key
PUSHER_APP_SECRET=your-secret
PUSHER_APP_ID=your-app-id
PUSHER_APP_CLUSTER=mt1
```

#### For Reverb (Laravel's WebSocket server)

```env
BROADCAST_DRIVER=reverb
REVERB_APP_KEY=your-reverb-key
REVERB_APP_SECRET=your-reverb-secret
REVERB_APP_ID=your-app-id
REVERB_HOST=localhost
REVERB_PORT=8080
REVERB_SCHEME=http
```

## Available Drivers

### Null Driver (Default)

No-operation driver for testing or disabling broadcasting.

```go
facades.Broadcast().Connection("null") // or let default handle
```

### Log Driver

Logs broadcast events to Goravel's logger for debugging and development.

```go
facades.Broadcast().Connection("log")
```

### Redis Driver

Uses Redis pub/sub for real-time broadcasting. Requires Redis to be configured.

```go
facades.Broadcast().Connection("redis")
```

### Ably Driver

Real-time messaging using Ably's infrastructure. Requires Ably Go SDK.

```go
// To enable Ably, add to go.mod:
// go get github.com/ably/ably-go/v2

facades.Broadcast().Connection("ably")
```

### Pusher Driver (Optional)

Integrates with Pusher service for real-time events.

```go
// To enable Pusher:
// 1. Add to go.mod: go get github.com/pusher/pusher-http-go/v5
// 2. Uncomment the pusher_broadcaster.go file

facades.Broadcast().Connection("pusher")
```

### Reverb Driver

Laravel's WebSocket server for real-time communication.

```go
facades.Broadcast().Connection("reverb")
```

## Usage

### Creating Broadcastable Events

Implement the `ShouldBroadcast` interface on your events to enable broadcasting:

```go
package events

import (
    "fmt"
    "github.com/rusmanplatd/goravelframework/broadcast"
    contractsbroadcast "github.com/rusmanplatd/goravelframework/contracts/broadcast"
)

type OrderShipped struct {
    OrderID        string
    TrackingNumber string
    CustomerEmail  string
}

// BroadcastOn returns the channels to broadcast on
func (e *OrderShipped) BroadcastOn() []contractsbroadcast.Channel {
    return []contractsbroadcast.Channel{
        broadcast.NewChannel("orders"),                                   // Public channel
        broadcast.NewPrivateChannel(fmt.Sprintf("order.%s", e.OrderID)), // Private channel
    }
}
```

### Enhanced Event Features

#### Custom Event Name

```go
func (e *OrderShipped) BroadcastAs() string {
    return "order.shipped" // Custom event name instead of struct name
}
```

#### Custom Broadcast Payload

```go
func (e *OrderShipped) BroadcastWith() map[string]any {
    return map[string]any{
        "order_id":        e.OrderID,
        "tracking_number": e.TrackingNumber,
        "customer_email":  e.CustomerEmail,
        "timestamp":       time.Now().Unix(),
    }
}
```

#### Specify Broadcast Connections

```go
func (e *OrderShipped) BroadcastConnections() []string {
    return []string{"redis", "ably"} // Broadcast to multiple services
}
```

#### Synchronous Broadcasting

```go
// Implement ShouldBroadcastNow for immediate broadcasting
var _ contractsbroadcast.ShouldBroadcastNow = (*UrgentNotification)(nil)

type UrgentNotification struct {
    Message string
    UserID  string
}

func (e *UrgentNotification) BroadcastOn() []contractsbroadcast.Channel {
    return []contractsbroadcast.Channel{
        broadcast.NewPrivateChannel(fmt.Sprintf("user.%s", e.UserID)),
    }
}
```

#### Exclude Current User

```go
// Implement DontBroadcastToCurrentUser to exclude sender
var _ contractsbroadcast.DontBroadcastToCurrentUser = (*ChatMessage)(nil)

type ChatMessage struct {
    Message  string
    UserID   string
    RoomID   string
}

func (e *ChatMessage) BroadcastOn() []contractsbroadcast.Channel {
    return []contractsbroadcast.Channel{
        broadcast.NewPresenceChannel(fmt.Sprintf("chat.room.%s", e.RoomID)),
    }
}

func (e *ChatMessage) DontBroadcastToCurrentUser() {
    // This prevents the message from being sent back to the sender
}
```

### Broadcasting Events

#### Automatic Event Broadcasting

Events are automatically broadcast when dispatched:

```go
// Dispatch event - automatically queued and broadcast
facades.Event().Dispatch(&OrderShipped{
    OrderID:        "12345",
    TrackingNumber: "TRACK123",
    CustomerEmail:  "customer@example.com",
})
```

#### Manual Broadcasting

Directly broadcast without creating events:

```go
broadcaster, err := facades.Broadcast().Connection()
if err != nil {
    log.Fatal("Failed to get broadcaster:", err)
}

channels := []contractsbroadcast.Channel{
    broadcast.NewChannel("notifications"),
    broadcast.NewPrivateChannel("user.123"),
}

payload := map[string]any{
    "message":  "Hello, World!",
    "type":     "info",
    "timestamp": time.Now().Unix(),
}

err = broadcaster.Broadcast(channels, "notification.sent", payload)
if err != nil {
    log.Fatal("Failed to broadcast:", err)
}
```

### Channel Types

#### Public Channels

Open channels that anyone can subscribe to:

```go
broadcast.NewChannel("orders")           // Broadcasts to: orders
broadcast.NewChannel("news")              // Broadcasts to: news
broadcast.NewChannel("stock-updates")     // Broadcasts to: stock-updates
```

#### Private Channels

Require authentication and are accessible only to authorized users:

```go
broadcast.NewPrivateChannel("order.123")        // Broadcasts to: private-order.123
broadcast.NewPrivateChannel("user.456")         // Broadcasts to: private-user.456
broadcast.NewPrivateChannel("notifications.789") // Broadcasts to: private-notifications.789
```

#### Presence Channels

Private channels that show who is online and require user presence data:

```go
broadcast.NewPresenceChannel("chat.room.1")           // Broadcasts to: presence-chat.room.1
broadcast.NewPresenceChannel("online-users")          // Broadcasts to: presence-online-users
broadcast.NewPresenceChannel("game.session.abc")     // Broadcasts to: presence-game.session.abc
```

### Channel Authentication

#### Setting Up Authentication

Register broadcast routes in your application:

```go
// In your route registration file
import "github.com/rusmanplatd/goravelframework/broadcast"

// Register broadcast authentication routes
broadcast.RegisterBroadcastRoutes(router, facades.Broadcast())
```

#### Private Channel Authorization

```go
// Register channel authorization callbacks
routes := broadcast.NewBroadcastRoutes(facades.Broadcast())

// Private channel for user notifications
routes.PrivateChannel("user.{id}", func(user any, params ...any) any {
    userID := params[0].(string)

    if userModel, ok := user.(*models.User); ok {
        // User can only access their own notification channel
        return userModel.ID == userID
    }

    return false
})

// Private channel for order access
routes.PrivateChannel("order.{id}", func(user any, params ...any) any {
    orderID := params[0].(string)

    if userModel, ok := user.(*models.User); ok {
        // Check if user owns this order or has permission
        return userModel.CanAccessOrder(orderID)
    }

    return false
}, contractsbroadcast.ChannelOption{
    Guards: []string{"api", "web"}, // Use specific auth guards
})
```

#### Presence Channel Authorization

```go
// Presence channel for chat rooms
routes.PresenceChannel("chat.room.{roomID}", func(user any, params ...any) any {
    roomID := params[0].(string)

    if userModel, ok := user.(*models.User); ok {
        // Check if user has access to this chat room
        if !userModel.HasAccessToChatRoom(roomID) {
            return false
        }

        // Return user data for presence channel
        return map[string]interface{}{
            "id":   userModel.ID,
            "name": userModel.Name,
            "avatar": userModel.AvatarURL,
            "role": userModel.GetRoleInRoom(roomID),
        }
    }

    return false
})

// Presence channel for online status
routes.PresenceChannel("status.{type}", func(user any, params ...any) any {
    statusType := params[0].(string)

    if userModel, ok := user.(*models.User); ok {
        return map[string]interface{}{
            "id":   userModel.ID,
            "name": userModel.Name,
            "status": userModel.CurrentStatus,
            "last_seen": userModel.LastSeen,
        }
    }

    return false
})
```

### User Authentication Setup

Configure user authentication for broadcasting:

```go
import (
    "github.com/rusmanplatd/goravelframework/broadcast"
)

// Set up user authentication callback
facades.Broadcast().ResolveAuthenticatedUserUsing(
    broadcast.UserAuthenticationCallback(func(request any) any {
        // Extract user from request (session, token, etc.)
        // This example shows extracting from HTTP request
        if httpReq, ok := request.(*http.Request); ok {
            return extractUserFromRequest(httpReq)
        }
        return nil
    }),
)

// Example user extraction function
func extractUserFromRequest(req *http.Request) any {
    // Extract token from header or cookie
    token := req.Header.Get("Authorization")
    if token == "" {
        token = extractTokenFromCookie(req)
    }

    // Validate token and return user
    if user, err := validateTokenAndGetUser(token); err == nil {
        return user
    }

    return nil
}
```

### Presence Channel Features

#### Creating Presence Channel Members

```go
import "github.com/rusmanplatd/goravelframework/broadcast"

// Create presence channel member
member := broadcast.NewPresenceChannelMember("user-123", map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
    "avatar": "https://example.com/avatars/john.jpg",
})

// Convert to JSON for WebSocket transmission
memberJSON, err := member.ToJSON()
if err != nil {
    log.Fatal("Failed to serialize member:", err)
}

// Create presence channel data for authentication
presenceData := broadcast.NewPresenceChannelData("user-123", map[string]interface{}{
    "name": "John Doe",
    "role": "member",
    "joined_at": time.Now().Unix(),
})

dataJSON, err := presenceData.ToJSON()
if err != nil {
    log.Fatal("Failed to serialize presence data:", err)
}
```

### Advanced Usage

#### Custom Broadcast Drivers

Create custom broadcast drivers:

```go
type CustomBroadcaster struct {
    *broadcast.BaseBroadcaster
    client *YourCustomClient
}

func NewCustomBroadcaster(client *YourCustomClient) *CustomBroadcaster {
    return &CustomBroadcaster{
        BaseBroadcaster: broadcast.NewBaseBroadcaster(),
        client:          client,
    }
}

func (c *CustomBroadcaster) Channel(channel string, callback contractsbroadcast.ChannelAuthCallback, options ...contractsbroadcast.ChannelOption) contractsbroadcast.Broadcaster {
    c.BaseBroadcaster.Channel(channel, callback, options...)
    return c
}

func (c *CustomBroadcaster) Broadcast(channels []contractsbroadcast.Channel, event string, payload map[string]any) error {
    formattedChannels := c.formatChannels(channels)

    for _, channel := range formattedChannels {
        if err := c.client.Publish(channel, event, payload); err != nil {
            return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
        }
    }

    return nil
}

func (c *CustomBroadcaster) Auth(request any) (any, error) {
    // Implement authentication logic
    return nil, nil
}

func (c *CustomBroadcaster) ValidAuthenticationResponse(request any, result any) (any, error) {
    // Implement authentication response logic
    return nil, nil
}

func (c *CustomBroadcaster) ResolveAuthenticatedUser(request any) (map[string]any, error) {
    return c.BaseBroadcaster.ResolveAuthenticatedUser(request)
}

func (c *CustomBroadcaster) ResolveAuthenticatedUserUsing(callback func(request any) map[string]any) {
    c.BaseBroadcaster.ResolveAuthenticatedUserUsing(callback)
}

// Register custom driver in config
func Broadcasting() map[string]any {
    return map[string]any{
        "default": "custom",
        "connections": map[string]any{
            "custom": map[string]any{
                "driver": "custom",
                "via": func() (contractsbroadcast.Broadcaster, error) {
                    return NewCustomBroadcaster(yourClient), nil
                },
            },
        },
    }
}
```

### Testing

#### Test with Null Driver

```env
BROADCAST_DRIVER=null
```

#### Test with Log Driver

```env
BROADCAST_DRIVER=log
```

#### Unit Testing Examples

```go
func TestEventBroadcasting(t *testing.T) {
    // Create test event
    event := &OrderShipped{
        OrderID:        "test-123",
        TrackingNumber: "TRACK-TEST",
        CustomerEmail:  "test@example.com",
    }

    // Verify channels
    channels := event.BroadcastOn()
    if len(channels) != 2 {
        t.Errorf("Expected 2 channels, got %d", len(channels))
    }

    // Verify channel types
    for _, channel := range channels {
        switch c := channel.(type) {
        case *broadcast.Channel:
            if c.GetName() != "orders" {
                t.Errorf("Expected 'orders' channel, got %s", c.GetName())
            }
        case *broadcast.PrivateChannel:
            expected := "private-order.test-123"
            if c.GetName() != expected {
                t.Errorf("Expected %s channel, got %s", expected, c.GetName())
            }
        default:
            t.Errorf("Unexpected channel type: %T", channel)
        }
    }
}
```

## Architecture Overview

The broadcast system follows Goravel's manager pattern with these key components:

- **Contracts**: `contracts/broadcast/` - Interfaces for broadcasters, events, and channels
- **Manager**: `broadcast/manager.go` - Manages connections and driver resolution
- **Drivers**: `broadcast/*_broadcaster.go` - Implementations for different services
- **Channels**: `broadcast/channels.go` - Channel type implementations with presence support
- **Authentication**: `broadcast/auth.go` - User authentication and authorization
- **Routes**: `broadcast/routes.go` - HTTP routes for channel authentication
- **Controller**: `broadcast/controller.go` - HTTP request handlers

## Best Practices

1. **Use appropriate channel types**:

   - Public channels for general information
   - Private channels for user-specific data
   - Presence channels for real-time collaboration

2. **Implement proper authentication**:

   - Always validate user permissions for private channels
   - Return minimal user data in presence channels
   - Use rate limiting for authentication endpoints

3. **Handle errors gracefully**:

   - Check broadcast errors
   - Implement retry logic for critical events
   - Use logging for debugging broadcast issues

4. **Optimize payload size**:
   - Keep broadcast payloads small
   - Use IDs instead of full objects when possible
   - Consider data compression for large payloads

## License

The Goravel framework is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
