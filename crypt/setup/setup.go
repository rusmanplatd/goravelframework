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
	cryptFacadePath := path.Facades("crypt.go")
	cryptServiceProvider := "&crypt.ServiceProvider{}"
	modulePath := packages.GetModulePath()

	packages.Setup(os.Args).
		Install(
			// Add the crypt service provider to the providers array in bootstrap/providers.go
			modify.AddProviderApply(modulePath, cryptServiceProvider),

			// Add the Crypt facade
			modify.WhenFacade(facades.Crypt, modify.File(cryptFacadePath).Overwrite(stubs.CryptFacade())),
		).
		Uninstall(
			modify.WhenNoFacades([]string{facades.Crypt},
				// Remove the crypt service provider from the providers array in bootstrap/providers.go
				modify.RemoveProviderApply(modulePath, cryptServiceProvider),
			),

			// Remove the Crypt facade
			modify.WhenFacade(facades.Crypt, modify.File(cryptFacadePath).Remove()),
		).
		Execute()
}
