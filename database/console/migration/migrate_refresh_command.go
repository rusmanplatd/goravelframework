package migration

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/errors"
)

type MigrateRefreshCommand struct {
	artisan console.Artisan
	config  config.Config
}

func NewMigrateRefreshCommand(artisan console.Artisan, config config.Config) *MigrateRefreshCommand {
	return &MigrateRefreshCommand{
		artisan: artisan,
		config:  config,
	}
}

// Signature The name and signature of the console command.
func (r *MigrateRefreshCommand) Signature() string {
	return "migrate:refresh"
}

// Description The console command description.
func (r *MigrateRefreshCommand) Description() string {
	return "Reset and re-run all migrations"
}

// Extend The console command extend.
func (r *MigrateRefreshCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
		Flags: []command.Flag{
			&command.IntFlag{
				Name:  "step",
				Value: 0,
				Usage: "refresh steps",
			},
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
func (r *MigrateRefreshCommand) Handle(ctx console.Context) error {
	schema := ctx.Option("schema")
	path := ctx.Option("path")

	// Override schema in config; since all sub-commands run in-process they will
	// inherit the patched config automatically.
	if schema != "" {
		restore := overrideSchema(r.config, schema)
		defer restore()
	}

	// Build extra flags to propagate --path to sub-commands (schema propagates via config).
	pathFlag := ""
	if path != "" {
		pathFlag = " --path " + path
	}

	if step := ctx.OptionInt("step"); step == 0 {
		if err := r.artisan.Call("migrate:reset" + pathFlag); err != nil {
			ctx.Error(errors.MigrationRefreshFailed.Args(err).Error())
			return nil
		}
	} else {
		if err := r.artisan.Call(fmt.Sprintf("migrate:rollback --step %d", step) + pathFlag); err != nil {
			ctx.Error(errors.MigrationRefreshFailed.Args(err).Error())
			return nil
		}
	}

	if err := r.artisan.Call("migrate" + pathFlag); err != nil {
		ctx.Error(errors.MigrationRefreshFailed.Args(err).Error())
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
			ctx.Error(errors.MigrationRefreshFailed.Args(err).Error())
			return nil
		}
	}
	ctx.Success("Migration refresh success")

	return nil
}
