package broadcast

import (
	"fmt"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
)

// Controller handles broadcast authentication requests.
// This is similar to Laravel's BroadcastController but adapted for Go/Goravel.
type Controller struct {
	manager contractsbroadcast.Manager
}

// NewController creates a new broadcast controller.
func NewController(manager contractsbroadcast.Manager) *Controller {
	return &Controller{
		manager: manager,
	}
}

// AuthRequest represents a broadcast authentication request.
type AuthRequest struct {
	ChannelName string `json:"channel_name"`
	SocketID    string `json:"socket_id"`
}

// Authenticate handles channel authentication requests.
func (c *Controller) Authenticate(request AuthRequest) (any, error) {
	// Get the default broadcaster
	broadcaster, err := c.manager.Connection()
	if err != nil {
		return nil, fmt.Errorf("failed to get broadcaster: %w", err)
	}

	// Use the broadcaster to authenticate
	return broadcaster.Auth(request)
}

// AuthenticateUser handles user authentication requests for presence channels.
func (c *Controller) AuthenticateUser(request AuthRequest) (any, error) {
	// Get the default broadcaster
	broadcaster, err := c.manager.Connection()
	if err != nil {
		return nil, fmt.Errorf("failed to get broadcaster: %w", err)
	}

	// Check if the broadcaster supports user authentication
	if result, err := broadcaster.ResolveAuthenticatedUser(request); err == nil && result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("user authentication not supported by current broadcaster")
}