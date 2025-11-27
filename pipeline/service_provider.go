package pipeline

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
)

type ServiceProvider struct{}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(binding.Pipeline, func(app foundation.Application) (any, error) {
		return NewPipeline(app), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	// No boot logic needed for Pipeline
}
