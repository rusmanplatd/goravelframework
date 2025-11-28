package event

import (
	"fmt"
	"strings"

	"github.com/rusmanplatd/goravelframework/contracts/event"
)

// Listen registers an event listener with the dispatcher.
func (app *Application) Listen(evt any, listeners ...any) error {
	if evt == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if len(listeners) == 0 {
		return fmt.Errorf("at least one listener is required")
	}

	// Parse event name
	eventName, err := app.getEventName(evt)
	if err != nil {
		return err
	}

	// Check if it's a wildcard pattern
	if strings.Contains(eventName, "*") {
		app.wildcards[eventName] = append(app.wildcards[eventName], listeners...)
		// Clear wildcard cache since wildcards changed
		app.wildcardsCache = make(map[string][]any)
	} else {
		app.listeners[eventName] = append(app.listeners[eventName], listeners...)
	}

	return nil
}

// HasListeners determines if a given event has listeners.
func (app *Application) HasListeners(eventName string) bool {
	// Check direct listeners
	if len(app.listeners[eventName]) > 0 {
		return true
	}

	// Check wildcard listeners
	for pattern := range app.wildcards {
		if matchWildcard(pattern, eventName) {
			return true
		}
	}

	return false
}

// Dispatch fires an event and calls the listeners synchronously.
func (app *Application) Dispatch(evt any, payload ...any) ([]any, error) {
	eventName, parsedPayload, err := parseEventAndPayload(evt, payload)
	if err != nil {
		return nil, err
	}

	return app.invokeListeners(eventName, parsedPayload, false)
}

// Until dispatches an event until the first non-null response is returned.
func (app *Application) Until(evt any, payload ...any) (any, error) {
	eventName, parsedPayload, err := parseEventAndPayload(evt, payload)
	if err != nil {
		return nil, err
	}

	responses, err := app.invokeListeners(eventName, parsedPayload, true)
	if err != nil {
		return nil, err
	}

	if len(responses) > 0 {
		return responses[0], nil
	}

	return nil, nil
}

// Subscribe registers an event subscriber with the dispatcher.
func (app *Application) Subscribe(subscriber event.Subscriber) error {
	if subscriber == nil {
		return fmt.Errorf("subscriber cannot be nil")
	}

	eventMap := subscriber.Subscribe(app)
	for evt, listenerList := range eventMap {
		// listenerList is already []any from the map
		if err := app.Listen(evt, listenerList...); err != nil {
			return err
		}
	}

	return nil
}

// Forget removes a set of listeners from the dispatcher.
func (app *Application) Forget(eventName string) {
	if strings.Contains(eventName, "*") {
		delete(app.wildcards, eventName)
		// Clear wildcard cache since wildcards changed
		app.wildcardsCache = make(map[string][]any)
	} else {
		delete(app.listeners, eventName)
	}
}

// Push registers an event and payload to be fired later.
func (app *Application) Push(eventName string, payload ...any) {
	if app.pushedEvents[eventName] == nil {
		app.pushedEvents[eventName] = make([]any, 0)
	}
	// Store each payload item separately so Flush can dispatch them individually
	app.pushedEvents[eventName] = append(app.pushedEvents[eventName], payload...)
}

// Flush flushes a set of pushed events.
func (app *Application) Flush(eventName string) error {
	payloads, ok := app.pushedEvents[eventName]
	if !ok {
		return nil
	}

	delete(app.pushedEvents, eventName)

	for _, payload := range payloads {
		if _, err := app.Dispatch(eventName, payload); err != nil {
			return err
		}
	}

	return nil
}

// invokeListeners invokes all listeners for a given event.
// If halt is true, stops at the first non-nil response.
func (app *Application) invokeListeners(eventName string, payload []any, halt bool) ([]any, error) {
	var responses []any

	// Get all listeners for this event
	allListeners := app.getListenersForEvent(eventName)

	for _, listener := range allListeners {
		// Check if listener should be queued
		if shouldQueueListener(listener, payload) {
			// Queue the listener instead of executing it synchronously
			if err := queueListener(app.queue, listener, eventName, payload); err != nil {
				return nil, fmt.Errorf("failed to queue listener: %w", err)
			}
			// Queued listeners don't return responses
			continue
		}

		callable, err := makeListenerCallable(listener)
		if err != nil {
			return nil, err
		}

		response, err := callable(eventName, payload)
		if err != nil {
			return nil, err
		}

		// If halt is enabled and we got a non-nil response, return immediately
		if halt && response != nil {
			return []any{response}, nil
		}

		// If response is false, stop propagation
		if boolResponse, ok := response.(bool); ok && !boolResponse {
			break
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// getListenersForEvent returns all listeners for a given event name.
func (app *Application) getListenersForEvent(eventName string) []any {
	var allListeners []any

	// Add direct listeners
	if listeners, ok := app.listeners[eventName]; ok {
		allListeners = append(allListeners, listeners...)
	}

	// Add wildcard listeners (use cache if available)
	if cached, ok := app.wildcardsCache[eventName]; ok {
		allListeners = append(allListeners, cached...)
	} else {
		// Build and cache wildcard listeners for this event
		var wildcardListeners []any
		for pattern, listeners := range app.wildcards {
			if matchWildcard(pattern, eventName) {
				wildcardListeners = append(wildcardListeners, listeners...)
			}
		}
		// Cache for future use
		app.wildcardsCache[eventName] = wildcardListeners
		allListeners = append(allListeners, wildcardListeners...)
	}

	return allListeners
}

// getEventName extracts the event name from various event types.
func (app *Application) getEventName(evt any) (string, error) {
	if evt == nil {
		return "", fmt.Errorf("event cannot be nil")
	}

	// If it's a string, return as-is
	if eventName, ok := evt.(string); ok {
		return eventName, nil
	}

	// If it's an Event interface, parse it
	eventName, _, err := parseEventAndPayload(evt, nil)
	return eventName, err
}
