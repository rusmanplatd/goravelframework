package route

import (
	"github.com/rusmanplatd/goravelframework/contracts/binding"
	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/errors"
	routeconsole "github.com/rusmanplatd/goravelframework/route/console"
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Route,
		},
		Dependencies: binding.Bindings[binding.Route].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(binding.Route, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()
		if config == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleRoute)
		}

		return NewRoute(config)
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	app.MakeArtisan().Register([]console.Command{
		routeconsole.NewList(app.MakeRoute()),
	})
}

func (r *ServiceProvider) Runners(app foundation.Application) []foundation.Runner {
	return []foundation.Runner{NewRouteRunner(app.MakeConfig(), app.MakeRoute())}
}
