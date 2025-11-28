package console

import (
	"strings"

	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/console/command"
	"github.com/rusmanplatd/goravelframework/packages"
	"github.com/rusmanplatd/goravelframework/support"
	supportconsole "github.com/rusmanplatd/goravelframework/support/console"
	"github.com/rusmanplatd/goravelframework/support/file"
)

type TestMakeCommand struct {
}

func NewTestMakeCommand() *TestMakeCommand {
	return &TestMakeCommand{}
}

// Signature The name and signature of the console command.
func (r *TestMakeCommand) Signature() string {
	return "make:test"
}

// Description The console command description.
func (r *TestMakeCommand) Description() string {
	return "Create a new test class"
}

// Extend The console command extend.
func (r *TestMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the test even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *TestMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "test", ctx.Argument(0), support.Config.Paths.Test)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	stub := r.getStub()

	if err := file.PutContent(m.GetFilePath(), r.populateStub(stub, m.GetPackageName(), m.GetStructName(), packages.GetModuleName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Test created successfully")

	return nil
}

func (r *TestMakeCommand) getStub() string {
	return Stubs{}.Test()
}

// populateStub Populate the place-holders in the command stub.
func (r *TestMakeCommand) populateStub(stub string, packageName, structName string, moduleName string) string {
	stub = strings.ReplaceAll(stub, "DummyTest", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)
	stub = strings.ReplaceAll(stub, "DummyModule", moduleName)

	return stub
}
