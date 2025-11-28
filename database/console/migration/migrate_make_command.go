package migration

import (
	"fmt"

	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/console/command"
	"github.com/rusmanplatd/goravelframework/contracts/database/migration"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/packages/match"
	"github.com/rusmanplatd/goravelframework/packages/modify"
	"github.com/rusmanplatd/goravelframework/support"
	supportconsole "github.com/rusmanplatd/goravelframework/support/console"
	"github.com/rusmanplatd/goravelframework/support/env"
	"github.com/rusmanplatd/goravelframework/support/str"
)

type MigrateMakeCommand struct {
	app      foundation.Application
	migrator migration.Migrator
}

func NewMigrateMakeCommand(app foundation.Application, migrator migration.Migrator) *MigrateMakeCommand {
	return &MigrateMakeCommand{app: app, migrator: migrator}
}

// Signature The name and signature of the console command.
func (r *MigrateMakeCommand) Signature() string {
	return "make:migration"
}

// Description The console command description.
func (r *MigrateMakeCommand) Description() string {
	return "Create a new migration file"
}

// Extend The console command extend.
func (r *MigrateMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
	}
}

// Handle Execute the console command.
func (r *MigrateMakeCommand) Handle(ctx console.Context) error {
	make, err := supportconsole.NewMake(ctx, "migration", ctx.Argument(0), support.Config.Paths.Migration)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	fileName, err := r.migrator.Create(make.GetName())
	if err != nil {
		ctx.Error(errors.MigrationCreateFailed.Args(err).Error())
		return nil
	}

	ctx.Success(fmt.Sprintf("Created Migration: %s", make.GetName()))

	structName := str.Of(fileName).Prepend("m_").Studly().String()
	if env.IsBootstrapSetup() {
		err = modify.AddMigration(make.GetPackageImportPath(), fmt.Sprintf("&%s.%s{}", make.GetPackageName(), structName))
	} else {
		err = r.registerInKernel(make.GetPackageImportPath(), structName)
	}

	if err != nil {
		ctx.Error(errors.MigrationRegisterFailed.Args(err).Error())
		return nil
	}

	ctx.Success("Migration registered successfully")

	return nil
}

// DEPRECATED: The kernel file will be removed in future versions.
func (r *MigrateMakeCommand) registerInKernel(pkg, structName string) error {
	return modify.GoFile(r.app.DatabasePath("kernel.go")).
		Find(match.Imports()).Modify(modify.AddImport(pkg)).
		Find(match.Migrations()).Modify(modify.Register(fmt.Sprintf("&migrations.%s{}", structName))).
		Apply()
}
