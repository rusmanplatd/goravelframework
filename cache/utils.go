package cache

import (
	"github.com/rusmanplatd/goravelframework/contracts/config"
)

func prefix(config config.Config) string {
	return config.GetString("cache.prefix") + ":"
}
