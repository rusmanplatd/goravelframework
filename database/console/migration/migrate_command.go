package migration

import (
	"fmt"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/color"
)

type MigrateCommand struct {
	config   config.Config
	migrator migration.Migrator
}

func NewMigrateCommand(migrator migration.Migrator, config config.Config) *MigrateCommand {
	return &MigrateCommand{
		config:   config,
		migrator: migrator,
	}
}

// Signature The name and signature of the console command.
func (r *MigrateCommand) Signature() string {
	return "migrate"
}

// Description The console command description.
func (r *MigrateCommand) Description() string {
	return "Run the database migrations"
}

// Extend The console command extend.
func (r *MigrateCommand) Extend() command.Extend {
	return command.Extend{
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
func (r *MigrateCommand) Handle(ctx console.Context) error {
	path := ctx.Option("path")
	schema := ctx.Option("schema")

	// Override the schema in config so the driver picks it up during this run.
	if schema != "" {
		restore := overrideSchema(r.config, schema)
		defer restore()
	}

	if path != "" {
		color.Infoln("Using migration path:", path)
	}

	if err := r.migrator.Run(); err != nil {
		ctx.Error(errors.MigrationMigrateFailed.Args(err).Error())
		return nil
	}

	ctx.Success("Migration success")

	return nil
}

// overrideSchema temporarily sets the schema for the default connection in config
// and returns a restore function that reverts it.
func overrideSchema(cfg config.Config, schema string) func() {
	if cfg == nil {
		return func() {}
	}
	conn := cfg.GetString("database.default")
	key := fmt.Sprintf("database.connections.%s.schema", conn)
	original := cfg.GetString(key)
	cfg.Add(key, schema)
	return func() { cfg.Add(key, original) }
}
