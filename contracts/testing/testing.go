package testing

import (
	"github.com/rusmanplatd/goravelframework/contracts/testing/docker"
)

type Testing interface {
	// Docker get the Docker instance.
	Docker() Docker
}

type Docker interface {
	// Cache gets a cache connection instance.
	Cache(store ...string) (docker.CacheDriver, error)
	// Database gets a database connection instance.
	Database(connection ...string) (docker.Database, error)
	// Image gets a image instance.
	Image(image docker.Image) docker.ImageDriver
}
