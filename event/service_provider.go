package event

import (
	"github.com/rusmanplatd/goravelframework/contracts/binding"
	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/errors"
	eventConsole "github.com/rusmanplatd/goravelframework/event/console"
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Event,
		},
		Dependencies: binding.Bindings[binding.Event].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(binding.Event, func(app foundation.Application) (any, error) {
		queueFacade := app.MakeQueue()
		if queueFacade == nil {
			return nil, errors.QueueFacadeNotSet.SetModule(errors.ModuleEvent)
		}

		return NewApplication(queueFacade), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	r.registerCommands(app)
}

func (r *ServiceProvider) registerCommands(app foundation.Application) {
	app.Commands([]console.Command{
		&eventConsole.EventMakeCommand{},
		&eventConsole.ListenerMakeCommand{},
	})
}
