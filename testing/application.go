package testing

import (
	contractscache "github.com/rusmanplatd/goravelframework/contracts/cache"
	contractsconfig "github.com/rusmanplatd/goravelframework/contracts/config"
	contractsconsole "github.com/rusmanplatd/goravelframework/contracts/console"
	contractsorm "github.com/rusmanplatd/goravelframework/contracts/database/orm"
	"github.com/rusmanplatd/goravelframework/contracts/testing"
	"github.com/rusmanplatd/goravelframework/testing/docker"
)

type Application struct {
	artisan contractsconsole.Artisan
	cache   contractscache.Cache
	config  contractsconfig.Config
	orm     contractsorm.Orm
}

func NewApplication(artisan contractsconsole.Artisan, cache contractscache.Cache, config contractsconfig.Config, orm contractsorm.Orm) *Application {
	return &Application{
		artisan: artisan,
		cache:   cache,
		config:  config,
		orm:     orm,
	}
}

func (r *Application) Docker() testing.Docker {
	return docker.NewDocker(r.artisan, r.cache, r.config, r.orm)
}
