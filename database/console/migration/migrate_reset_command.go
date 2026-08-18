package migration

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/migration"
)

type MigrateResetCommand struct {
	config   config.Config
	migrator migration.Migrator
}

func NewMigrateResetCommand(migrator migration.Migrator, config config.Config) *MigrateResetCommand {
	return &MigrateResetCommand{
		config:   config,
		migrator: migrator,
	}
}

// Signature The name and signature of the console command.
func (r *MigrateResetCommand) Signature() string {
	return "migrate:reset"
}

// Description The console command description.
func (r *MigrateResetCommand) Description() string {
	return "Rollback all database migrations"
}

// Extend The console command extend.
func (r *MigrateResetCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
		Flags: []command.Flag{
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
func (r *MigrateResetCommand) Handle(ctx console.Context) error {
	schema := ctx.Option("schema")
	if schema != "" {
		restore := overrideSchema(r.config, schema)
		defer restore()
	}

	if err := r.migrator.Reset(); err != nil {
		ctx.Error(err.Error())
	}

	ctx.Success("Migration reset success")

	return nil
}
