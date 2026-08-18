package migration

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/errors"
)

type MigrateRollbackCommand struct {
	config   config.Config
	migrator migration.Migrator
}

func NewMigrateRollbackCommand(migrator migration.Migrator, config config.Config) *MigrateRollbackCommand {
	return &MigrateRollbackCommand{
		config:   config,
		migrator: migrator,
	}
}

// Signature The name and signature of the console command.
func (r *MigrateRollbackCommand) Signature() string {
	return "migrate:rollback"
}

// Description The console command description.
func (r *MigrateRollbackCommand) Description() string {
	return "Rollback the database migrations"
}

// Extend The console command extend.
func (r *MigrateRollbackCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
		Flags: []command.Flag{
			&command.IntFlag{
				Name:  "step",
				Value: 0,
				Usage: "rollback steps",
			},
			&command.IntFlag{
				Name:  "batch",
				Value: 0,
				Usage: "rollback batch number (only can be used in default driver)",
			},
			&command.StringFlag{
				Name:  "path",
				Usage: "the path to the migrations folder (overrides config default)",
			},
			&command.StringFlag{
				Name:  "schema",
				Usage: "the database schema to use (overrides config default)",
			},
		},
	}
}

// Handle Execute the console command.
func (r *MigrateRollbackCommand) Handle(ctx console.Context) error {
	schema := ctx.Option("schema")
	if schema != "" {
		restore := overrideSchema(r.config, schema)
		defer restore()
	}

	step := ctx.OptionInt("step")
	batch := ctx.OptionInt("batch")

	// Validate inputs
	if step < 0 {
		ctx.Error("The step option should be a positive integer")
		return nil
	}
	if batch < 0 {
		ctx.Error("The batch option should be a positive integer")
		return nil
	}
	if step > 0 && batch > 0 {
		ctx.Error("The step and batch options cannot be used together")
		return nil
	}

	if err := r.migrator.Rollback(step, batch); err != nil {
		ctx.Error(errors.MigrationMigrateFailed.Args(err).Error())
		return nil
	}

	ctx.Success("Migration rollback success")

	return nil
}
