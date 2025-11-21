package broadcast

import (
	"fmt"
	"time"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
)

// This file contains practical examples of using the enhanced broadcast features.

// ExampleEvent demonstrates a comprehensive broadcastable event.
type ExampleEvent struct {
	UserID    string
	ChatRoom  string
	Message   string
	Timestamp time.Time
}

// BroadcastOn defines the channels for this event.
func (e *ExampleEvent) BroadcastOn() []contractsbroadcast.Channel {
	return []contractsbroadcast.Channel{
		NewChannel("general-messages"),
		NewPrivateChannel(fmt.Sprintf("user.%s", e.UserID)),
		NewPresenceChannel(fmt.Sprintf("chat.room.%s", e.ChatRoom)),
	}
}

// BroadcastAs defines a custom event name.
func (e *ExampleEvent) BroadcastAs() string {
	return "chat.message"
}

// BroadcastWith defines custom payload data.
func (e *ExampleEvent) BroadcastWith() map[string]any {
	return map[string]any{
		"user_id":    e.UserID,
		"chat_room":  e.ChatRoom,
		"message":    e.Message,
		"timestamp":  e.Timestamp.Unix(),
		"message_id": fmt.Sprintf("msg_%d_%s", e.Timestamp.Unix(), e.UserID),
	}
}

// BroadcastConnections specifies which broadcasters to use.
func (e *ExampleEvent) BroadcastConnections() []string {
	return []string{"redis", "ably"}
}

// Example of setting up comprehensive channel authentication
func setupChannelAuthentication(manager contractsbroadcast.Manager) {
	// Get the default broadcaster
	broadcaster, _ := manager.Connection()

	// Public channel - no authentication needed
	broadcaster.Channel("general-messages", func(user any, params ...any) any {
		return true // Allow all users
	})

	// Private channel - user-specific notifications
	broadcaster.Channel("private-user.{id}", func(user any, params ...any) any {
		userID := params[0].(string)

		// Example: Check if user can access their own channel
		if authUser, ok := user.(Authenticatable); ok {
			return authUser.GetAuthIdentifier() == userID
		}

		// Fallback: check user map or struct
		if userMap, ok := user.(map[string]any); ok {
			if id, exists := userMap["id"]; exists {
				return fmt.Sprintf("%v", id) == userID
			}
		}

		return false
	})

	// Presence channel - chat room with user data
	broadcaster.Channel("presence-chat.room.{roomID}", func(user any, params ...any) any {
		roomID := params[0].(string)

		// Verify user has access to chat room
		if !userCanAccessChatRoom(user, roomID) {
			return false
		}

		// Return user presence data
		return extractUserPresenceData(user)
	})

	// Multi-parameter channel example
	broadcaster.Channel("private-order.{userID}.{orderID}", func(user any, params ...any) any {
		userID := params[0].(string)
		orderID := params[1].(string)

		// Check if user owns this order
		return userOwnsOrder(user, userID, orderID)
	})
}

// Example user authentication setup
func setupUserAuthentication(manager contractsbroadcast.Manager) {
	// Set up user authentication callback
	broadcaster, _ := manager.Connection()

	broadcaster.ResolveAuthenticatedUserUsing(UserAuthenticationCallback(func(request any) any {
		// Example: Extract user from different request types
		switch req := request.(type) {
		case map[string]any:
			// Extract from map-based request
			return extractUserFromMap(req)
		case string:
			// Extract from token string
			return extractUserFromToken(req)
		default:
			return nil
		}
	}))
}

// Helper functions for the examples
func userCanAccessChatRoom(user any, roomID string) bool {
	// Implement your chat room access logic here
	return true // Placeholder
}

func extractUserPresenceData(user any) map[string]interface{} {
	// Extract user data for presence channels
	if authUser, ok := user.(Authenticatable); ok {
		return map[string]interface{}{
			"id":   authUser.GetAuthIdentifier(),
			"name": getUserName(user),
		}
	}

	userID, _ := GetBroadcastIdentifier(user)
	return map[string]interface{}{
		"id":   userID,
		"name": "Anonymous",
	}
}

func userOwnsOrder(user any, userID, orderID string) bool {
	// Implement order ownership logic
	return true // Placeholder
}

func getUserName(user any) string {
	// Extract user name from various user types
	if userMap, ok := user.(map[string]any); ok {
		if name, exists := userMap["name"]; exists {
			return fmt.Sprintf("%v", name)
		}
	}
	return "User"
}

func extractUserFromMap(req map[string]any) any {
	// Extract user from map-based request
	if user, exists := req["user"]; exists {
		return user
	}
	return nil
}

func extractUserFromToken(token string) any {
	// Extract user from JWT token or similar
	// This is where you'd implement your token validation logic
	return nil
}

// ExampleChannelPatterns demonstrates different channel pattern examples
var ExampleChannelPatterns = struct {
	PrivateUser      string
	PrivateOrder     string
	PrivateOrderUser string
	PresenceChat     string
	PresenceGame     string
	MultiParam       string
}{
	PrivateUser:      "private-user.{id}",
	PrivateOrder:     "private-order.{orderID}",
	PrivateOrderUser: "private-order.{userID}.{orderID}",
	PresenceChat:     "presence-chat.room.{roomID}",
	PresenceGame:     "presence-game.session.{sessionID}",
	MultiParam:       "private-team.{teamID}.channel.{channelID}",
}

// ExamplePatternMatching demonstrates how the enhanced pattern matching works
func ExamplePatternMatching() {
	broadcaster := NewBaseBroadcaster()

	// Test different channel patterns
	testCases := []struct {
		channel string
		pattern string
		expect  bool
	}{
		{"private-user.123", "private-user.{id}", true},
		{"private-order.456", "private-order.{orderID}", true},
		{"private-order.user123.456", "private-order.{userID}.{orderID}", true},
		{"presence-chat.room.789", "presence-chat.room.{roomID}", true},
		{"invalid-channel", "private-user.{id}", false},
		{"wrong-pattern", "presence-game.session.{sessionID}", false},
	}

	for _, tc := range testCases {
		matches := broadcaster.channelNameMatchesPattern(tc.channel, tc.pattern)
		fmt.Printf("Channel: %s, Pattern: %s, Matches: %v (Expected: %v)\n",
			tc.channel, tc.pattern, matches, tc.expect)

		// Extract parameters if matches
		if matches {
			params := broadcaster.extractParametersFromPattern(tc.channel, tc.pattern)
			fmt.Printf("  Extracted parameters: %v\n", params)
		}
	}
}

// ExamplePresenceChannelMember demonstrates creating presence channel members
func ExamplePresenceChannelMember() {
	// Create a presence channel member
	member := NewPresenceChannelMember("user-123", map[string]interface{}{
		"name": "John Doe",
		"email": "john@example.com",
		"avatar": "https://example.com/avatars/john.jpg",
		"role": "member",
	})

	// Convert to JSON for WebSocket transmission
	memberJSON, err := member.ToJSON()
	if err != nil {
		fmt.Printf("Error serializing member: %v\n", err)
		return
	}

	fmt.Printf("Member JSON: %s\n", string(memberJSON))

	// Create presence channel data for authentication
	presenceData := NewPresenceChannelData("user-123", map[string]interface{}{
		"name":      "John Doe",
		"role":      "member",
		"joined_at": time.Now().Unix(),
	})

	dataJSON, err := presenceData.ToJSON()
	if err != nil {
		fmt.Printf("Error serializing presence data: %v\n", err)
		return
	}

	fmt.Printf("Presence Data JSON: %s\n", string(dataJSON))
}