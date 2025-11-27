package event

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation"
)

// EventServiceProvider is a base service provider for event registration.
// Applications can embed this struct and override the Listen() and Subscribe() methods
// to define their event-listener mappings.
//
// Example usage:
//
//	type AppEventServiceProvider struct {
//	    event.EventServiceProvider
//	}
//
//	func (a *AppEventServiceProvider) Listen() map[any][]any {
//	    return map[any][]any{
//	        &events.UserRegistered{}: {
//	            &listeners.SendWelcomeEmail{},
//	            &listeners.CreateUserProfile{},
//	        },
//	        "notification.*": {
//	            &listeners.LogNotification{},
//	        },
//	    }
//	}
type EventServiceProvider struct {
	app foundation.Application
}

// NewEventServiceProvider creates a new EventServiceProvider instance.
func NewEventServiceProvider(app foundation.Application) *EventServiceProvider {
	return &EventServiceProvider{
		app: app,
	}
}

// Register registers the service provider.
// This is called during the application bootstrap process.
func (e *EventServiceProvider) Register(app foundation.Application) error {
	e.app = app
	return nil
}

// Boot boots the service provider.
// This is where event listeners and subscribers are registered.
func (e *EventServiceProvider) Boot(app foundation.Application) error {
	eventInstance := app.MakeEvent()

	// Register event listeners
	listeners := e.Listen()
	for evt, listenerList := range listeners {
		if err := eventInstance.Listen(evt, listenerList...); err != nil {
			return err
		}
	}

	// Register event subscribers
	subscribers := e.Subscribe()
	for _, subscriber := range subscribers {
		if err := eventInstance.Subscribe(subscriber); err != nil {
			return err
		}
	}

	return nil
}

// Listen returns a map of events to listeners.
// Override this method in your application's EventServiceProvider to define event-listener mappings.
//
// Example:
//
//	func (a *AppEventServiceProvider) Listen() map[any][]any {
//	    return map[any][]any{
//	        &events.UserRegistered{}: {
//	            &listeners.SendWelcomeEmail{},
//	            &listeners.CreateUserProfile{},
//	        },
//	        "notification.*": {
//	            &listeners.LogNotification{},
//	        },
//	    }
//	}
func (e *EventServiceProvider) Listen() map[any][]any {
	return make(map[any][]any)
}

// Subscribe returns a slice of event subscribers.
// Override this method in your application's EventServiceProvider to define subscribers.
//
// Example:
//
//	func (a *AppEventServiceProvider) Subscribe() []event.Subscriber {
//	    return []event.Subscriber{
//	        &subscribers.UserEventSubscriber{},
//	    }
//	}
func (e *EventServiceProvider) Subscribe() []event.Subscriber {
	return make([]event.Subscriber, 0)
}
