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
	validationFacadePath := path.Facades("validation.go")
	validationServiceProvider := "&validation.ServiceProvider{}"
	modulePath := packages.GetModulePath()

	packages.Setup(os.Args).
		Install(
			// Add the validation service provider to the providers array in bootstrap/providers.go
			modify.AddProviderApply(modulePath, validationServiceProvider),

			// Add the Validation facade
			modify.WhenFacade(facades.Validation, modify.File(validationFacadePath).Overwrite(stubs.ValidationFacade())),
		).
		Uninstall(
			modify.WhenNoFacades([]string{facades.Validation},
				// Remove the validation service provider from the providers array in bootstrap/providers.go
				modify.RemoveProviderApply(modulePath, validationServiceProvider),
			),

			// Remove the Validation facade
			modify.WhenFacade(facades.Validation, modify.File(validationFacadePath).Remove()),
		).
		Execute()
}
