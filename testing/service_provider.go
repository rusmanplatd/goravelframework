package testing

import (
	"github.com/rusmanplatd/goravelframework/contracts/binding"
	contractsconsole "github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	contractsroute "github.com/rusmanplatd/goravelframework/contracts/route"
	contractsession "github.com/rusmanplatd/goravelframework/contracts/session"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/support/color"
)

var (
	json          foundation.Json
	artisanFacade contractsconsole.Artisan
	routeFacade   contractsroute.Route
	sessionFacade contractsession.Manager
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Testing,
		},
		Dependencies: binding.Bindings[binding.Testing].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(binding.Testing, func(app foundation.Application) (any, error) {
		return NewApplication(app.MakeArtisan(), app.MakeCache(), app.MakeConfig(), app.MakeOrm()), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	artisanFacade = app.MakeArtisan()
	if artisanFacade == nil {
		color.Errorln(errors.ConsoleFacadeNotSet.SetModule(errors.ModuleTesting))
	}

	routeFacade = app.MakeRoute()
	if routeFacade == nil {
		color.Errorln(errors.RouteFacadeNotSet.SetModule(errors.ModuleTesting))
	}

	sessionFacade = app.MakeSession()
	if sessionFacade == nil {
		color.Errorln(errors.SessionFacadeNotSet.SetModule(errors.ModuleTesting))
	}

	json = app.GetJson()
}
