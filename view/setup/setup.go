package main

import (
	"os"

	"github.com/rusmanplatd/goravelframework/packages"
	"github.com/rusmanplatd/goravelframework/packages/modify"
	"github.com/rusmanplatd/goravelframework/support/path"
)

func main() {
	stubs := Stubs{}
	viewServiceProvider := "&view.ServiceProvider{}"
	modulePath := packages.GetModulePath()
	viewFacade := "View"
	viewFacadePath := path.Facades("view.go")

	packages.Setup(os.Args).
		Install(
			// Add the view service provider to the providers array in bootstrap/providers.go
			modify.AddProviderApply(modulePath, viewServiceProvider),

			// Add the View facade
			modify.WhenFacade(viewFacade, modify.File(viewFacadePath).Overwrite(stubs.ViewFacade())),
		).
		Uninstall(
			modify.WhenNoFacades([]string{viewFacade},
				// Remove the view service provider from the providers array in bootstrap/providers.go
				modify.RemoveProviderApply(modulePath, viewServiceProvider),
			),
			modify.WhenFacade(viewFacade, modify.File(viewFacadePath).Remove()),
		).
		Execute()
}
