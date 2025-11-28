package binding

import (
	"github.com/rusmanplatd/goravelframework/contracts/binding"
	"github.com/rusmanplatd/goravelframework/support/collect"
)

func Dependencies(bindings ...string) []string {
	var deps []string
	for _, bind := range bindings {
		deps = append(deps, binding.Bindings[bind].Dependencies...)
	}

	return collect.Diff(collect.Unique(deps), bindings)
}
