package docker

import (
	contractscache "github.com/rusmanplatd/goravelframework/contracts/cache"
	contractsconfig "github.com/rusmanplatd/goravelframework/contracts/config"
	contractsconsole "github.com/rusmanplatd/goravelframework/contracts/console"
	contractsorm "github.com/rusmanplatd/goravelframework/contracts/database/orm"
	"github.com/rusmanplatd/goravelframework/contracts/testing/docker"
	"github.com/rusmanplatd/goravelframework/errors"
)

type Docker struct {
	artisan contractsconsole.Artisan
	cache   contractscache.Cache
	config  contractsconfig.Config
	orm     contractsorm.Orm
}

func NewDocker(artisan contractsconsole.Artisan, cache contractscache.Cache, config contractsconfig.Config, orm contractsorm.Orm) *Docker {
	return &Docker{
		artisan: artisan,
		cache:   cache,
		config:  config,
		orm:     orm,
	}
}

func (r *Docker) Cache(store ...string) (docker.CacheDriver, error) {
	if r.config == nil {
		return nil, errors.ConfigFacadeNotSet
	}

	if len(store) == 0 {
		store = append(store, r.config.GetString("cache.default"))
	}

	return r.cache.Store(store[0]).Docker()
}

func (r *Docker) Database(connection ...string) (docker.Database, error) {
	if len(connection) == 0 {
		return NewDatabase(r.artisan, r.config, r.orm, "")
	} else {
		return NewDatabase(r.artisan, r.config, r.orm, connection[0])
	}
}

func (r *Docker) Image(image docker.Image) docker.ImageDriver {
	return NewImageDriver(image)
}
