package broadcast

import (
	"fmt"
	"reflect"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
)

// Authenticatable represents a user that can be authenticated for broadcasting.
type Authenticatable interface {
	// GetAuthIdentifier returns the unique identifier for the user.
	GetAuthIdentifier() string

	// GetAuthIdentifierForBroadcasting returns the identifier used for broadcasting.
	GetAuthIdentifierForBroadcasting() string
}

// DefaultAuthenticator provides default authentication functionality.
type DefaultAuthenticator struct {
	retrieveUser func(request any) any
}

// NewDefaultAuthenticator creates a new default authenticator.
func NewDefaultAuthenticator(retrieveUser func(request any) any) *DefaultAuthenticator {
	return &DefaultAuthenticator{
		retrieveUser: retrieveUser,
	}
}

// AuthenticateUser authenticates a user from a request.
func (a *DefaultAuthenticator) AuthenticateUser(request any) (any, error) {
	if a.retrieveUser == nil {
		return nil, fmt.Errorf("user retrieval function not configured")
	}

	user := a.retrieveUser(request)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetBroadcastIdentifier extracts the broadcast identifier from a user.
func GetBroadcastIdentifier(user any) (string, error) {
	// Check if user implements Authenticatable interface
	if authUser, ok := user.(Authenticatable); ok {
		return authUser.GetAuthIdentifierForBroadcasting(), nil
	}

	// Try map-based user
	if id, found := getIDFromMap(user); found {
		return id, nil
	}

	// Try method-based user
	if id, found := getIDFromMethods(user); found {
		return id, nil
	}

	// Try field-based user
	if id, found := getIDFromFields(user); found {
		return id, nil
	}

	return "", fmt.Errorf("unable to extract broadcast identifier from user")
}

// getIDFromMap tries to extract ID from map-based user
func getIDFromMap(user any) (string, bool) {
	userMap, ok := user.(map[string]any)
	if !ok {
		return "", false
	}

	if id, exists := userMap["id"]; exists {
		if idStr, ok := id.(string); ok {
			return idStr, true
		}
	}

	if id, exists := userMap["ID"]; exists {
		if idStr, ok := id.(string); ok {
			return idStr, true
		}
	}

	return "", false
}

// getIDFromMethods tries to extract ID using reflection methods
func getIDFromMethods(user any) (string, bool) {
	userValue := reflect.ValueOf(user)
	userType := userValue.Type()

	methodNames := []string{"GetID", "GetAuthIdentifier", "GetAuthIdentifierForBroadcasting"}

	for _, methodName := range methodNames {
		if method, found := userType.MethodByName(methodName); found {
			results := method.Func.Call([]reflect.Value{userValue})
			if len(results) > 0 {
				if id := results[0].Interface(); id != nil {
					if idStr, ok := id.(string); ok {
						return idStr, true
					}
				}
			}
		}
	}

	return "", false
}

// getIDFromFields tries to extract ID from struct fields
func getIDFromFields(user any) (string, bool) {
	userValue := reflect.ValueOf(user)
	if userValue.Kind() != reflect.Struct {
		return "", false
	}

	fieldNames := []string{"ID", "Id"}

	for _, fieldName := range fieldNames {
		if field := userValue.FieldByName(fieldName); field.IsValid() && field.Kind() == reflect.String {
			return field.String(), true
		}
	}

	return "", false
}

// UserAuthenticationCallback creates a user authentication callback for broadcasting.
func UserAuthenticationCallback(retrieveUser func(request any) any) func(request any) map[string]any {
	return func(request any) map[string]any {
		user := retrieveUser(request)
		if user == nil {
			return nil
		}

		identifier, err := GetBroadcastIdentifier(user)
		if err != nil {
			return nil
		}

		return map[string]any{
			"identifier": identifier,
			"user":       user,
		}
	}
}

// ChannelAuthorizationChecker checks if a user can access a channel.
type ChannelAuthorizationChecker struct {
	authenticator *DefaultAuthenticator
}

// NewChannelAuthorizationChecker creates a new channel authorization checker.
func NewChannelAuthorizationChecker(authenticator *DefaultAuthenticator) *ChannelAuthorizationChecker {
	return &ChannelAuthorizationChecker{
		authenticator: authenticator,
	}
}

// CheckAuthorization checks if a user can access a channel.
func (c *ChannelAuthorizationChecker) CheckAuthorization(request any, channel string, callback contractsbroadcast.ChannelAuthCallback) (any, error) {
	user, err := c.authenticator.AuthenticateUser(request)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Extract parameters from channel pattern
	params := extractChannelParameters(channel)

	// Call the authorization callback
	result := callback(user, params...)

	if result == false {
		return nil, fmt.Errorf("access denied to channel: %s", channel)
	}

	return result, nil
}

// extractChannelParameters extracts parameters from a channel name.
// This is a simplified version - in production, you'd want more sophisticated pattern matching.
func extractChannelParameters(channel string) []any {
	// Simple parameter extraction for common patterns
	// Example: "private-user.123" -> ["123"]
	// Example: "presence-chat.1.room.5" -> ["1", "5"]

	var params []any
	parts := []rune(channel)

	for i := 0; i < len(parts); i++ {
		if parts[i] == '.' || parts[i] == '-' {
			// Look ahead to find the parameter value
			start := i + 1
			if start >= len(parts) {
				continue
			}

			// Find the end of the parameter
			end := start
			for end < len(parts) && parts[end] != '.' && parts[end] != '-' {
				end++
			}

			if end > start {
				param := string(parts[start:end])
				params = append(params, param)
				i = end - 1 // Skip to the end of the parameter
			}
		}
	}

	return params
}