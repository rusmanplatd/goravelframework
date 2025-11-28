package console

import (
	"fmt"
	"strings"

	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/console/command"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/packages/modify"
	"github.com/rusmanplatd/goravelframework/support"
	supportconsole "github.com/rusmanplatd/goravelframework/support/console"
	"github.com/rusmanplatd/goravelframework/support/env"
	"github.com/rusmanplatd/goravelframework/support/file"
)

type MiddlewareMakeCommand struct {
}

// Signature The name and signature of the console command.
func (r *MiddlewareMakeCommand) Signature() string {
	return "make:middleware"
}

// Description The console command description.
func (r *MiddlewareMakeCommand) Description() string {
	return "Create a new middleware class"
}

// Extend The console command extend.
func (r *MiddlewareMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the middleware even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *MiddlewareMakeCommand) Handle(ctx console.Context) error {
	make, err := supportconsole.NewMake(ctx, "middleware", ctx.Argument(0), support.Config.Paths.Middleware)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(make.GetFilePath(), r.populateStub(r.getStub(), make.GetPackageName(), make.GetStructName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Middleware created successfully")

	if env.IsBootstrapSetup() {
		if err := modify.AddMiddleware(make.GetPackageImportPath(), fmt.Sprintf("%s.%s()", make.GetPackageName(), make.GetStructName())); err != nil {
			ctx.Error(errors.MiddlewareRegisterFailed.Args(make.GetStructName(), err).Error())
			return nil
		}

		ctx.Success("Middleware registered successfully")
	}

	return nil
}

func (r *MiddlewareMakeCommand) getStub() string {
	return Stubs{}.Middleware()
}

// populateStub Populate the place-holders in the command stub.
func (r *MiddlewareMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyMiddleware", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
