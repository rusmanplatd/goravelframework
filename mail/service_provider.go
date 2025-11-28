package mail

import (
	"github.com/rusmanplatd/goravelframework/contracts/binding"
	contractsconsole "github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	contractsqueue "github.com/rusmanplatd/goravelframework/contracts/queue"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/mail/console"
	"github.com/rusmanplatd/goravelframework/support/color"
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Mail,
		},
		Dependencies: binding.Bindings[binding.Mail].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Bind(binding.Mail, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()
		if config == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleMail)
		}

		queue := app.MakeQueue()
		if queue == nil {
			return nil, errors.QueueFacadeNotSet.SetModule(errors.ModuleMail)
		}

		return NewApplication(config, queue)
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	app.Commands([]contractsconsole.Command{
		console.NewMailMakeCommand(),
	})

	r.registerJobs(app)
}

func (r *ServiceProvider) registerJobs(app foundation.Application) {
	queueFacade := app.MakeQueue()
	if queueFacade == nil {
		color.Warningln("Queue Facade is not initialized. Skipping job registration.")
		return
	}

	configFacade := app.MakeConfig()
	if configFacade == nil {
		color.Warningln("Config Facade is not initialized. Skipping job registration.")
		return
	}

	queueFacade.Register([]contractsqueue.Job{
		NewSendMailJob(configFacade),
	})
}
