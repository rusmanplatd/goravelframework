package notification

import (
	"github.com/rusmanplatd/goravelframework/contracts/binding"
	contractsconsole "github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/notification/console"
)

// ServiceProvider provides notification services.
type ServiceProvider struct {
}

// Relationship returns the service provider relationship.
func (s *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Notification,
		},
		Dependencies: []string{
			binding.Config,
			binding.Event,
			binding.Log,
			binding.Mail,
			binding.Orm,
			binding.Queue,
		},
		ProvideFor: []string{},
	}
}

// Register registers the service provider.
func (s *ServiceProvider) Register(app foundation.Application) {
	app.Bind(binding.Notification, func(app foundation.Application) (any, error) {
		// Get dependencies
		config := app.MakeConfig()
		if config == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleNotification)
		}

		event := app.MakeEvent()
		if event == nil {
			return nil, errors.EventFacadeNotSet.SetModule(errors.ModuleNotification)
		}

		log := app.MakeLog()
		if log == nil {
			return nil, errors.LogFacadeNotSet.SetModule(errors.ModuleNotification)
		}

		mail := app.MakeMail()
		if mail == nil {
			return nil, errors.MailFacadeNotSet.SetModule(errors.ModuleNotification)
		}

		orm := app.MakeOrm()
		if orm == nil {
			return nil, errors.OrmFacadeNotSet.SetModule(errors.ModuleNotification)
		}

		queue := app.MakeQueue()
		if queue == nil {
			return nil, errors.QueueFacadeNotSet.SetModule(errors.ModuleNotification)
		}

		// Create and return the channel manager
		return NewChannelManager(
			config,
			event,
			log,
			mail,
			orm,
			queue,
		), nil
	})
}

// Boot boots the service provider.
func (s *ServiceProvider) Boot(app foundation.Application) {
	// Register console commands
	app.Commands([]contractsconsole.Command{
		console.NewNotificationMakeCommand(),
	})
}
