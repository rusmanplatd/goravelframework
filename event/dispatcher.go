package event

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rusmanplatd/goravelframework/contracts/event"
)

// parseEventAndPayload parses the given event and payload and prepares them for dispatching.
// If event is an Event interface, returns its type name and wraps it in payload.
// If event is a string, returns it as-is with the provided payload.
func parseEventAndPayload(evt any, payload []any) (string, []any, error) {
	if evt == nil {
		return "", nil, fmt.Errorf("event cannot be nil")
	}

	// Check if it's an Event interface
	if e, ok := evt.(event.Event); ok {
		eventName := reflect.TypeOf(e).String()
		// Remove package prefix if present (e.g., "*events.NotificationSent" -> "NotificationSent")
		parts := strings.Split(eventName, ".")
		if len(parts) > 1 {
			eventName = parts[len(parts)-1]
		}
		// Remove pointer prefix if present
		eventName = strings.TrimPrefix(eventName, "*")
		return eventName, []any{e}, nil
	}

	// Check if it's a string event name
	if eventName, ok := evt.(string); ok {
		return eventName, payload, nil
	}

	// For other types, use the type name
	eventName := reflect.TypeOf(evt).String()
	parts := strings.Split(eventName, ".")
	if len(parts) > 1 {
		eventName = parts[len(parts)-1]
	}
	eventName = strings.TrimPrefix(eventName, "*")
	return eventName, []any{evt}, nil
}

// matchWildcard checks if an event name matches a wildcard pattern.
// Supports patterns like "user.*", "notification.*", etc.
func matchWildcard(pattern, eventName string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == eventName
	}

	// Simple wildcard matching
	parts := strings.Split(pattern, "*")
	if len(parts) == 0 {
		return false
	}

	// Check prefix
	if parts[0] != "" && !strings.HasPrefix(eventName, parts[0]) {
		return false
	}

	// Check suffix
	if len(parts) > 1 && parts[len(parts)-1] != "" {
		if !strings.HasSuffix(eventName, parts[len(parts)-1]) {
			return false
		}
	}

	return true
}

// invokeListener calls a listener with the given event name and payload.
// Handles different listener types: closures, Listener interfaces, and structs with Handle method.
func invokeListener(listener any, eventName string, payload []any) (any, error) {
	if listener == nil {
		return nil, fmt.Errorf("listener cannot be nil")
	}

	// Check if it's a Listener interface
	if l, ok := listener.(event.Listener); ok {
		err := l.Handle(payload...)
		return nil, err
	}

	// Check if it's a closure/function
	listenerValue := reflect.ValueOf(listener)
	if listenerValue.Kind() == reflect.Func {
		return invokeFunction(listenerValue, eventName, payload)
	}

	// Check if it has a Handle method
	handleMethod := listenerValue.MethodByName("Handle")
	if handleMethod.IsValid() && handleMethod.Kind() == reflect.Func {
		return invokeFunction(handleMethod, eventName, payload)
	}

	return nil, fmt.Errorf("listener must be a function, Listener interface, or have a Handle method")
}

// invokeFunction invokes a function with the appropriate arguments.
func invokeFunction(fn reflect.Value, eventName string, payload []any) (any, error) {
	args := prepareArguments(fn, eventName, payload)
	results := fn.Call(args)
	return processResults(results)
}

// prepareArguments prepares the arguments for a function call.
func prepareArguments(fn reflect.Value, eventName string, payload []any) []reflect.Value {
	fnType := fn.Type()
	numIn := fnType.NumIn()
	var args []reflect.Value

	// Simply map payload to function parameters
	for i := 0; i < numIn; i++ {
		if i < len(payload) {
			if payload[i] == nil {
				args = append(args, reflect.Zero(fnType.In(i)))
			} else {
				pValue := reflect.ValueOf(payload[i])
				targetType := fnType.In(i)

				// Try to convert if types don't match exactly
				if pValue.Type().AssignableTo(targetType) {
					args = append(args, pValue)
				} else if pValue.Type().ConvertibleTo(targetType) {
					args = append(args, pValue.Convert(targetType))
				} else {
					args = append(args, pValue)
				}
			}
		} else {
			// Fill remaining parameters with zero values
			args = append(args, reflect.Zero(fnType.In(i)))
		}
	}

	return args
}

// processResults processes the results from a function call.
func processResults(results []reflect.Value) (any, error) {
	if len(results) == 0 {
		return nil, nil
	}

	// Check if last result is an error
	lastResult := results[len(results)-1]
	if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !lastResult.IsNil() {
			return nil, lastResult.Interface().(error)
		}
		// If there's only error return, return nil
		if len(results) == 1 {
			return nil, nil
		}
		// Return first result if there are multiple
		if results[0].IsValid() && results[0].CanInterface() {
			return results[0].Interface(), nil
		}
		return nil, nil
	}

	// Return first result
	if results[0].IsValid() && results[0].CanInterface() {
		return results[0].Interface(), nil
	}

	return nil, nil
}

// makeListenerCallable converts various listener types into a callable function.
func makeListenerCallable(listener any) (func(string, []any) (any, error), error) {
	if listener == nil {
		return nil, fmt.Errorf("listener cannot be nil")
	}

	return func(eventName string, payload []any) (any, error) {
		return invokeListener(listener, eventName, payload)
	}, nil
}
