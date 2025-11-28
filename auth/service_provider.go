package auth

import (
	"context"

	"github.com/rusmanplatd/goravelframework/auth/access"
	"github.com/rusmanplatd/goravelframework/auth/console"
	contractsbinding "github.com/rusmanplatd/goravelframework/contracts/binding"
	"github.com/rusmanplatd/goravelframework/contracts/cache"
	"github.com/rusmanplatd/goravelframework/contracts/config"
	contractconsole "github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/database/orm"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/contracts/http"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/support/binding"
)

var (
	cacheFacade  cache.Cache
	configFacade config.Config
	ormFacade    orm.Orm
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() contractsbinding.Relationship {
	bindings := []string{
		contractsbinding.Auth,
		contractsbinding.Gate,
	}

	return contractsbinding.Relationship{
		Bindings:     bindings,
		Dependencies: binding.Dependencies(bindings...),
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.BindWith(contractsbinding.Auth, func(app foundation.Application, parameters map[string]any) (any, error) {
		configFacade = app.MakeConfig()
		if configFacade == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleAuth)
		}

		log := app.MakeLog()
		if log == nil {
			return nil, errors.LogFacadeNotSet.SetModule(errors.ModuleAuth)
		}

		ctx, ok := parameters["ctx"]
		if ok {
			return NewAuth(ctx.(http.Context), configFacade, log)
		}

		// ctx is optional when calling facades.Auth().Extend()
		return NewAuth(nil, configFacade, log)
	})
	app.Singleton(contractsbinding.Gate, func(app foundation.Application) (any, error) {
		return access.NewGate(context.Background()), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	cacheFacade = app.MakeCache()
	ormFacade = app.MakeOrm()

	r.registerCommands(app)
}

func (r *ServiceProvider) registerCommands(app foundation.Application) {
	app.Commands([]contractconsole.Command{
		console.NewJwtSecretCommand(app.MakeConfig()),
		console.NewPolicyMakeCommand(),
	})
}
