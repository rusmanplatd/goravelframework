package broadcast

import (
	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
	contractshttp "github.com/goravel/framework/contracts/http"
	contractroute "github.com/goravel/framework/contracts/route"
)

// RegisterBroadcastRoutes registers the broadcast authentication routes.
// This should be called in your routes registration file.
func RegisterBroadcastRoutes(router contractroute.Router, manager contractsbroadcast.Manager) {
	controller := NewController(manager)

	// POST /broadcasting/auth
	router.Post("/broadcasting/auth", func(ctx contractshttp.Context) contractshttp.Response {
		var request AuthRequest
		if err := ctx.Request().Bind(&request); err != nil {
			return ctx.Response().Json(400, contractshttp.Json{
				"error": "Invalid request format",
			})
		}

		response, err := controller.Authenticate(request)
		if err != nil {
			return ctx.Response().Json(403, contractshttp.Json{
				"error": err.Error(),
			})
		}

		return ctx.Response().Json(200, response)
	})

	// POST /broadcasting/auth/user
	router.Post("/broadcasting/auth/user", func(ctx contractshttp.Context) contractshttp.Response {
		var request AuthRequest
		if err := ctx.Request().Bind(&request); err != nil {
			return ctx.Response().Json(400, contractshttp.Json{
				"error": "Invalid request format",
			})
		}

		response, err := controller.AuthenticateUser(request)
		if err != nil {
			return ctx.Response().Json(403, contractshttp.Json{
				"error": err.Error(),
			})
		}

		return ctx.Response().Json(200, response)
	})
}

// BroadcastRoutes provides a helper to register broadcast channels.
// This is similar to Laravel's Broadcast::channel() method.
type BroadcastRoutes struct {
	manager contractsbroadcast.Manager
}

// NewBroadcastRoutes creates a new broadcast routes instance.
func NewBroadcastRoutes(manager contractsbroadcast.Manager) *BroadcastRoutes {
	return &BroadcastRoutes{
		manager: manager,
	}
}

// Channel registers a new broadcast channel.
func (b *BroadcastRoutes) Channel(channel string, callback contractsbroadcast.ChannelAuthCallback, options ...contractsbroadcast.ChannelOption) {
	broadcaster, err := b.manager.Connection()
	if err != nil {
		// Log error or handle appropriately
		return
	}

	broadcaster.Channel(channel, callback, options...)
}

// PrivateChannel registers a new private broadcast channel.
func (b *BroadcastRoutes) PrivateChannel(channel string, callback contractsbroadcast.ChannelAuthCallback, options ...contractsbroadcast.ChannelOption) {
	fullChannel := "private-" + channel
	b.Channel(fullChannel, callback, options...)
}

// PresenceChannel registers a new presence broadcast channel.
func (b *BroadcastRoutes) PresenceChannel(channel string, callback contractsbroadcast.ChannelAuthCallback, options ...contractsbroadcast.ChannelOption) {
	fullChannel := "presence-" + channel
	b.Channel(fullChannel, callback, options...)
}
