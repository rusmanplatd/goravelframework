package broadcast

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
)

// AblyBroadcaster provides broadcasting functionality using Ably.
// Note: This requires the Ably Go SDK to be added to go.mod:
// github.com/ably/ably-go/v2
type AblyBroadcaster struct {
	*BaseBroadcaster
	client any // ably.Client or ably.RealtimeClient - requires ably-go SDK
	key    string
}

// AblyMessage represents an Ably message for broadcasting.
type AblyMessage struct {
	Name          string      `json:"name"`
	Data          interface{} `json:"data"`
	ConnectionKey string      `json:"connectionKey,omitempty"`
}

// AblyAuthResponse represents the authentication response for Ably.
type AblyAuthResponse struct {
	Auth        string `json:"auth"`
	ChannelData string `json:"channel_data,omitempty"`
}

// NewAbly creates a new Ably broadcaster instance.
func NewAbly(key string) *AblyBroadcaster {
	return &AblyBroadcaster{
		BaseBroadcaster: NewBaseBroadcaster(),
		key:             key,
	}
}

// Channel registers a channel authenticator callback.
func (a *AblyBroadcaster) Channel(channel string, callback contractsbroadcast.ChannelAuthCallback, options ...contractsbroadcast.ChannelOption) contractsbroadcast.Broadcaster {
	a.BaseBroadcaster.Channel(channel, callback, options...)
	return a
}

// NewAblyWithClient creates a new Ably broadcaster with a pre-configured client.
func NewAblyWithClient(client any, key string) *AblyBroadcaster {
	return &AblyBroadcaster{
		BaseBroadcaster: NewBaseBroadcaster(),
		client:          client,
		key:             key,
	}
}

// Broadcast sends an event to the specified channels via Ably.
func (a *AblyBroadcaster) Broadcast(channels []contractsbroadcast.Channel, event string, payload map[string]any) error {
	if a.client == nil {
		return fmt.Errorf("ably client not initialized - add ably-go SDK and configure client")
	}

	// Implementation would require ably-go SDK
	// This is a placeholder showing the expected interface
	/*
		for _, channel := range a.formatChannels(channels) {
			channel := a.client.Channels.Get(channel)
			message := a.buildAblyMessage(event, payload)
			err := channel.Publish(context.Background(), message.Name, message.Data)
			if err != nil {
				return fmt.Errorf("ably broadcast error: %w", err)
			}
		}
	*/

	return fmt.Errorf("ably broadcaster not fully implemented - requires ably-go SDK")
}

// Auth authenticates the incoming request for channel access.
func (a *AblyBroadcaster) Auth(request any) (any, error) {
	// This would need to be implemented based on the HTTP request structure
	// For now, return the interface compliance
	return nil, fmt.Errorf("auth method requires HTTP request implementation")
}

// ValidAuthenticationResponse returns the valid authentication response for Ably.
func (a *AblyBroadcaster) ValidAuthenticationResponse(request any, result any) (any, error) {
	// This would need to be implemented based on the HTTP request structure
	// For now, return the interface compliance
	return nil, fmt.Errorf("valid authentication response requires HTTP request implementation")
}

// ResolveAuthenticatedUser resolves the authenticated user payload for connection requests.
func (a *AblyBroadcaster) ResolveAuthenticatedUser(request any) (map[string]any, error) {
	return a.BaseBroadcaster.ResolveAuthenticatedUser(request)
}

// ResolveAuthenticatedUserUsing registers the user retrieval callback for authentication.
func (a *AblyBroadcaster) ResolveAuthenticatedUserUsing(callback func(request any) map[string]any) {
	a.BaseBroadcaster.ResolveAuthenticatedUserUsing(callback)
}

// buildAblyMessage builds an Ably message object for broadcasting.
func (a *AblyBroadcaster) buildAblyMessage(event string, payload map[string]any) *AblyMessage {
	message := &AblyMessage{
		Name: event,
		Data: payload,
	}

	if connectionKey, ok := payload["socket"]; ok {
		if connectionKeyStr, ok := connectionKey.(string); ok {
			message.ConnectionKey = connectionKeyStr
		}
	}

	return message
}

// formatChannels formats the channel names for Ably.
func (a *AblyBroadcaster) formatChannels(channels []contractsbroadcast.Channel) []string {
	formatted := make([]string, len(channels))
	for i, channel := range channels {
		channelName := channel.GetName()

		if strings.HasPrefix(channelName, "private-") {
			formatted[i] = strings.Replace(channelName, "private-", "private:", 1)
		} else if strings.HasPrefix(channelName, "presence-") {
			formatted[i] = strings.Replace(channelName, "presence-", "presence:", 1)
		} else {
			formatted[i] = "public:" + channelName
		}
	}
	return formatted
}

// generateAblySignature generates the signature needed for Ably authentication.
func (a *AblyBroadcaster) generateAblySignature(channelName, socketID string, userData []byte) string {
	keyParts := strings.Split(a.key, ":")
	if len(keyParts) != 2 {
		return ""
	}

	privateKey := keyParts[1]

	data := socketID + ":" + channelName
	if len(userData) > 0 {
		data += ":" + string(userData)
	}

	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// getPublicToken gets the public token value from the Ably key.
func (a *AblyBroadcaster) getPublicToken() string {
	keyParts := strings.Split(a.key, ":")
	if len(keyParts) >= 1 {
		return keyParts[0]
	}
	return ""
}

// getPrivateToken gets the private token value from the Ably key.
func (a *AblyBroadcaster) getPrivateToken() string {
	keyParts := strings.Split(a.key, ":")
	if len(keyParts) >= 2 {
		return keyParts[1]
	}
	return ""
}

// GetAbly returns the underlying Ably client.
func (a *AblyBroadcaster) GetAbly() any {
	return a.client
}

// SetAbly sets the underlying Ably client.
func (a *AblyBroadcaster) SetAbly(client any) {
	a.client = client
}