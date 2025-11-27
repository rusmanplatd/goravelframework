package notification

import (
	"fmt"
)

// AnonymousNotifiable represents a notifiable entity without a persistent identity.
// It allows sending notifications to arbitrary addresses without requiring a database model.
type AnonymousNotifiable struct {
	routes map[string]any
}

// NewAnonymousNotifiable creates a new anonymous notifiable.
func NewAnonymousNotifiable() *AnonymousNotifiable {
	return &AnonymousNotifiable{
		routes: make(map[string]any),
	}
}

// Route sets the route for a given channel.
// For example: Route("mail", "user@example.com") or Route("sms", "+1234567890")
func (a *AnonymousNotifiable) Route(channel string, route any) *AnonymousNotifiable {
	a.routes[channel] = route
	return a
}

// RouteNotificationFor returns the notification routing information for the given channel.
// This method is called by notification channels to determine where to send the notification.
func (a *AnonymousNotifiable) RouteNotificationFor(channel string) any {
	if route, ok := a.routes[channel]; ok {
		return route
	}
	return nil
}

// String returns a string representation of the anonymous notifiable.
func (a *AnonymousNotifiable) String() string {
	return fmt.Sprintf("AnonymousNotifiable{routes: %v}", a.routes)
}
