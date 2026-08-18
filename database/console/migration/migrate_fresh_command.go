package migration

import (
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/migration"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/color"
)

type MigrateFreshCommand struct {
	artisan  console.Artisan
	config   config.Config
	migrator migration.Migrator
}

func NewMigrateFreshCommand(artisan console.Artisan, migrator migration.Migrator, config config.Config) *MigrateFreshCommand {
	return &MigrateFreshCommand{
		artisan:  artisan,
		config:   config,
		migrator: migrator,
	}
}

// Signature The name and signature of the console command.
func (r *MigrateFreshCommand) Signature() string {
	return "migrate:fresh"
}

// Description The console command description.
func (r *MigrateFreshCommand) Description() string {
	return "Drop all tables and re-run all migrations"
}

// Extend The console command extend.
func (r *MigrateFreshCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:  "seed",
				Usage: "seed the database after running migrations",
			},
			&command.StringSliceFlag{
				Name:  "seeder",
				Usage: "specify the seeder(s) to use for seeding the database",
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
func (r *MigrateFreshCommand) Handle(ctx console.Context) error {
	path := ctx.Option("path")
	schema := ctx.Option("schema")

	// Override schema in config so it propagates to the "migrate" sub-call inside Fresh().
	if schema != "" {
		restore := overrideSchema(r.config, schema)
		defer restore()
	}

	if path != "" {
		color.Infoln("Using migration path:", path)
	}

	if err := r.migrator.Fresh(); err != nil {
		ctx.Error(errors.MigrationFreshFailed.Args(err).Error())
		return nil
	}

	// Seed the database if the "seed" flag is provided
	if ctx.OptionBool("seed") {
		seeders := ctx.OptionSlice("seeder")
		seederFlag := ""
		if len(seeders) > 0 {
			seederFlag = " --seeder " + strings.Join(seeders, ",")
		}

		if err := r.artisan.Call("db:seed" + seederFlag); err != nil {
			ctx.Error(errors.MigrationFreshFailed.Args(err).Error())
			return nil
		}
	}

	ctx.Success("Migration fresh success")

	return nil
}
