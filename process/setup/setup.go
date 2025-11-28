package main

import (
	"os"

	"github.com/rusmanplatd/goravelframework/contracts/facades"
	"github.com/rusmanplatd/goravelframework/packages"
	"github.com/rusmanplatd/goravelframework/packages/modify"
	"github.com/rusmanplatd/goravelframework/support/path"
)

func main() {
	stubs := Stubs{}
	processFacadePath := path.Facades("process.go")
	modulePath := packages.GetModulePath()
	processServiceProvider := "&process.ServiceProvider{}"

	packages.Setup(os.Args).
		Install(
			// Add the process service provider to the providers array in bootstrap/providers.go
			modify.AddProviderApply(modulePath, processServiceProvider),

			// Add the Process facade
			modify.WhenFacade(facades.Process, modify.File(processFacadePath).Overwrite(stubs.ProcessFacade())),
		).
		Uninstall(
			modify.WhenNoFacades([]string{facades.Process},
				// Remove the process service provider from the providers array in bootstrap/providers.go
				modify.RemoveProviderApply(modulePath, processServiceProvider),
			),

			// Remove the Process facade
			modify.WhenFacade(facades.Process, modify.File(processFacadePath).Remove()),
		).
		Execute()
}
