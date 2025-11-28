package packages

import (
	"github.com/rusmanplatd/goravelframework/contracts/packages/modify"
)

type Setup interface {
	Install(modifiers ...modify.Apply) Setup
	Uninstall(modifiers ...modify.Apply) Setup
	Execute()
}
