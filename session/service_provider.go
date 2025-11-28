package session

import (
	"github.com/rusmanplatd/goravelframework/contracts/binding"
	"github.com/rusmanplatd/goravelframework/contracts/config"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/contracts/session"
	"github.com/rusmanplatd/goravelframework/errors"
)

var (
	SessionFacade session.Manager
	ConfigFacade  config.Config
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Session,
		},
		Dependencies: binding.Bindings[binding.Session].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(binding.Session, func(app foundation.Application) (any, error) {
		c := app.MakeConfig()
		if c == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleSession)
		}

		j := app.GetJson()
		if j == nil {
			return nil, errors.JSONParserNotSet.SetModule(errors.ModuleSession)
		}

		return NewManager(c, j), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	SessionFacade = app.MakeSession()
	ConfigFacade = app.MakeConfig()
}
