package console

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/console/command"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/packages/match"
	"github.com/rusmanplatd/goravelframework/packages/modify"
	"github.com/rusmanplatd/goravelframework/support"
	supportconsole "github.com/rusmanplatd/goravelframework/support/console"
	"github.com/rusmanplatd/goravelframework/support/file"
	"github.com/rusmanplatd/goravelframework/support/str"
)

type JobMakeCommand struct {
}

// Signature The name and signature of the console command.
func (r *JobMakeCommand) Signature() string {
	return "make:job"
}

// Description The console command description.
func (r *JobMakeCommand) Description() string {
	return "Create a new job class"
}

// Extend The console command extend.
func (r *JobMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the job even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *JobMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "job", ctx.Argument(0), support.Config.Paths.Job)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(m.GetFilePath(), r.populateStub(r.getStub(), m.GetPackageName(), m.GetStructName(), m.GetSignature())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Job created successfully")

	if err = modify.GoFile(filepath.Join("app", "providers", "queue_service_provider.go")).
		Find(match.Imports()).Modify(modify.AddImport(m.GetPackageImportPath())).
		Find(match.Jobs()).Modify(modify.Register(fmt.Sprintf("&%s.%s{}", m.GetPackageName(), m.GetStructName()))).
		Apply(); err != nil {
		ctx.Warning(errors.QueueJobRegisterFailed.Args(err).Error())
		return nil
	}

	ctx.Success("Job registered successfully")

	return nil
}

func (r *JobMakeCommand) getStub() string {
	return JobStubs{}.Job()
}

// populateStub Populate the place-holders in the command stub.
func (r *JobMakeCommand) populateStub(stub string, packageName, structName, signature string) string {
	stub = strings.ReplaceAll(stub, "DummyJob", structName)
	stub = strings.ReplaceAll(stub, "DummySignature", str.Of(signature).Snake().String())
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
