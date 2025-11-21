package broadcast

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
)

// BaseBroadcaster provides common functionality for all broadcasters.
type BaseBroadcaster struct {
	channels                  map[string]contractsbroadcast.ChannelAuthCallback
	options                   map[string]contractsbroadcast.ChannelOption
	authenticatedUserCallback func(request any) map[string]any
	mu                        sync.RWMutex
}

// NewBaseBroadcaster creates a new base broadcaster.
func NewBaseBroadcaster() *BaseBroadcaster {
	return &BaseBroadcaster{
		channels: make(map[string]contractsbroadcast.ChannelAuthCallback),
		options:  make(map[string]contractsbroadcast.ChannelOption),
	}
}

// Channel registers a channel authenticator callback.
func (b *BaseBroadcaster) Channel(channel string, callback contractsbroadcast.ChannelAuthCallback, opts ...contractsbroadcast.ChannelOption) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.channels[channel] = callback

	if len(opts) > 0 {
		b.options[channel] = opts[0]
	}
}

// formatChannels converts Channel interfaces to string names.
func (b *BaseBroadcaster) formatChannels(channels []contractsbroadcast.Channel) []string {
	formatted := make([]string, len(channels))
	for i, channel := range channels {
		formatted[i] = channel.GetName()
	}
	return formatted
}

// isGuardedChannel checks if a channel requires authentication.
func (b *BaseBroadcaster) isGuardedChannel(channelName string) bool {
	return strings.HasPrefix(channelName, "private-") || strings.HasPrefix(channelName, "presence-")
}

// normalizeChannelName removes the prefix from channel names.
func (b *BaseBroadcaster) normalizeChannelName(channelName string) string {
	channelName = strings.TrimPrefix(channelName, "private-")
	channelName = strings.TrimPrefix(channelName, "presence-")
	return channelName
}

// getChannelCallback retrieves the callback for a channel pattern.
func (b *BaseBroadcaster) getChannelCallback(channel string) (contractsbroadcast.ChannelAuthCallback, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Try exact match first
	if callback, exists := b.channels[channel]; exists {
		return callback, true
	}

	// Try pattern matching (simple wildcard support)
	for pattern, callback := range b.channels {
		if b.channelNameMatchesPattern(channel, pattern) {
			return callback, true
		}
	}

	return nil, false
}

// channelNameMatchesPattern checks if a channel name matches a pattern.
func (b *BaseBroadcaster) channelNameMatchesPattern(channel, pattern string) bool {
	// Exact match for simple patterns
	if !strings.Contains(pattern, "{") {
		return channel == pattern
	}

	// Convert pattern to regex for better matching
	regexPattern := b.patternToRegex(pattern)
	matched, err := regexp.MatchString("^"+regexPattern+"$", channel)
	if err != nil {
		// Fallback to simple matching if regex fails
		return b.simplePatternMatch(channel, pattern)
	}

	return matched
}

// patternToRegex converts a channel pattern to regex pattern.
func (b *BaseBroadcaster) patternToRegex(pattern string) string {
	// Escape special regex characters
	regexPattern := regexp.QuoteMeta(pattern)

	// Convert {param} to regex capture groups
	regexPattern = regexp.MustCompile(`\\\{([^}]+)\\\}`).ReplaceAllString(regexPattern, `([^\.]+)`)

	// Replace escaped dots with literal dots
	regexPattern = strings.ReplaceAll(regexPattern, `\.`, `\.`)

	return regexPattern
}

// simplePatternMatch provides fallback pattern matching.
func (b *BaseBroadcaster) simplePatternMatch(channel, pattern string) bool {
	// Replace {id} style placeholders with wildcard
	pattern = strings.ReplaceAll(pattern, "{id}", "*")
	pattern = strings.ReplaceAll(pattern, "{name}", "*")

	// Simple wildcard matching
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(channel, prefix)
	}

	return false
}

// getChannelOptions retrieves the options for a channel.
func (b *BaseBroadcaster) getChannelOptions(channel string) contractsbroadcast.ChannelOption {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if opts, exists := b.options[channel]; exists {
		return opts
	}

	return contractsbroadcast.ChannelOption{}
}

// ResolveAuthenticatedUser resolves the authenticated user payload for connection requests.
func (b *BaseBroadcaster) ResolveAuthenticatedUser(request any) (map[string]any, error) {
	if b.authenticatedUserCallback != nil {
		return b.authenticatedUserCallback(request), nil
	}
	return nil, nil
}

// ResolveAuthenticatedUserUsing registers the user retrieval callback for authentication.
func (b *BaseBroadcaster) ResolveAuthenticatedUserUsing(callback func(request any) map[string]any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.authenticatedUserCallback = callback
}

// verifyUserCanAccessChannel verifies if a user can access a channel.
func (b *BaseBroadcaster) verifyUserCanAccessChannel(request any, channel string, retrieveUser func(request any, channel string) any) (any, error) {
	callback, exists := b.getChannelCallback(channel)
	if !exists {
		return nil, fmt.Errorf("channel not registered: %s", channel)
	}

	// Extract parameters from channel pattern matching
	params := b.extractChannelParameters(channel)

	user := retrieveUser(request, channel)
	result := callback(user, params...)

	if result == false {
		return nil, fmt.Errorf("access denied to channel: %s", channel)
	}

	return result, nil
}

// extractChannelParameters extracts parameters from channel name based on registered patterns.
func (b *BaseBroadcaster) extractChannelParameters(channel string) []any {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for pattern := range b.channels {
		if params := b.extractParametersFromPattern(channel, pattern); len(params) > 0 {
			return params
		}
	}
	return []any{}
}

// extractParametersFromPattern extracts parameters from a channel name based on a pattern.
func (b *BaseBroadcaster) extractParametersFromPattern(channel, pattern string) []any {
	// Simple parameter extraction for {param} style patterns
	if !strings.Contains(pattern, "{") {
		return []any{}
	}

	// Check if pattern matches first
	if !b.channelNameMatchesPattern(channel, pattern) {
		return []any{}
	}

	// Extract parameter names from pattern
	paramNames := b.extractParameterNames(pattern)
	if len(paramNames) == 0 {
		return []any{}
	}

	// Create regex to extract values
	regexPattern := b.patternToRegex(pattern)
	re, err := regexp.Compile("^" + regexPattern + "$")
	if err != nil {
		// Fallback to simple extraction
		return b.simpleParameterExtraction(channel, pattern)
	}

	// Extract values using regex
	matches := re.FindStringSubmatch(channel)
	if len(matches) <= 1 { // matches[0] is the full match
		return []any{}
	}

	// Convert matches to any slice (skip the full match)
	params := make([]any, 0, len(matches)-1)
	for i := 1; i < len(matches); i++ {
		params = append(params, matches[i])
	}

	return params
}

// extractParameterNames extracts parameter names from a pattern.
func (b *BaseBroadcaster) extractParameterNames(pattern string) []string {
	var names []string
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	for _, match := range matches {
		if len(match) > 1 {
			names = append(names, match[1])
		}
	}

	return names
}

// simpleParameterExtraction provides fallback parameter extraction.
func (b *BaseBroadcaster) simpleParameterExtraction(channel, pattern string) []any {
	parts := strings.Split(pattern, "}")
	channelParts := strings.Split(channel, "-")

	if len(parts) != len(channelParts) {
		return []any{}
	}

	var params []any
	for i, part := range parts {
		if strings.Contains(part, "{") {
			if i < len(channelParts) {
				// Remove the prefix part before the parameter
				value := channelParts[i]
				if prefix := strings.Split(part, "{")[0]; prefix != "" {
					value = strings.TrimPrefix(value, prefix)
				}
				params = append(params, value)
			}
		}
	}

	return params
}
