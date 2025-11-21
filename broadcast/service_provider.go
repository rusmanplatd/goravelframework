package broadcast

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
)

type ServiceProvider struct{}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Broadcast,
		},
		Dependencies: binding.Bindings[binding.Broadcast].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(binding.Broadcast, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()
		json := app.GetJson()
		log := app.MakeLog()

		// Redis is optional and provided by external package
		// Set to nil for now - users can configure custom driver if needed
		var redis any = nil

		return NewManager(config, json, log, redis), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	// Boot logic if needed
}
