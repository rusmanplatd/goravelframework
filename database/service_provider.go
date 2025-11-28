package database

import (
	"context"
	"fmt"

	contractsbinding "github.com/rusmanplatd/goravelframework/contracts/binding"
	contractsconsole "github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/database/driver"
	"github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/database/console"
	consolemigration "github.com/rusmanplatd/goravelframework/database/console/migration"
	"github.com/rusmanplatd/goravelframework/database/db"
	"github.com/rusmanplatd/goravelframework/database/migration"
	databaseorm "github.com/rusmanplatd/goravelframework/database/orm"
	databaseschema "github.com/rusmanplatd/goravelframework/database/schema"
	databaseseeder "github.com/rusmanplatd/goravelframework/database/seeder"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/support/binding"
	"github.com/rusmanplatd/goravelframework/support/color"
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() contractsbinding.Relationship {
	bindings := []string{
		contractsbinding.Orm,
		contractsbinding.DB,
		contractsbinding.Schema,
		contractsbinding.Seeder,
	}

	return contractsbinding.Relationship{
		Bindings:     bindings,
		Dependencies: binding.Dependencies(bindings...),
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(contractsbinding.Orm, func(app foundation.Application) (any, error) {
		ctx := context.Background()
		config := app.MakeConfig()
		if config == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleOrm)
		}

		log := app.MakeLog()
		if log == nil {
			return nil, errors.LogFacadeNotSet.SetModule(errors.ModuleOrm)
		}

		connection := config.GetString("database.default")
		if connection == "" {
			return nil, nil
		}

		orm, err := databaseorm.BuildOrm(ctx, config, connection, log, app.Fresh)
		if err != nil {
			color.Warningln(errors.OrmInitConnection.Args(connection, err).SetModule(errors.ModuleOrm))

			return orm, nil
		}

		return orm, nil
	})

	app.Singleton(contractsbinding.DB, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()
		if config == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleDB)
		}

		log := app.MakeLog()
		if log == nil {
			return nil, errors.LogFacadeNotSet.SetModule(errors.ModuleDB)
		}

		connection := config.GetString("database.default")
		if connection == "" {
			return nil, nil
		}

		db, err := db.BuildDB(context.Background(), config, log, connection)
		if err != nil {
			color.Warningln(errors.OrmInitConnection.Args(connection, err).SetModule(errors.ModuleDB))

			return nil, nil
		}

		return db, nil
	})

	app.Singleton(contractsbinding.Schema, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()
		if config == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleSchema)
		}

		log := app.MakeLog()
		if log == nil {
			return nil, errors.LogFacadeNotSet.SetModule(errors.ModuleSchema)
		}

		orm := app.MakeOrm()
		if orm == nil {
			// The Orm module will print the error message, so it's safe to return an empty schema.
			return &databaseschema.Schema{}, nil
		}

		driverCallback, exist := config.Get(fmt.Sprintf("database.connections.%s.via", orm.Name())).(func() (driver.Driver, error))
		if !exist {
			return nil, errors.DatabaseConfigNotFound
		}

		driver, err := driverCallback()
		if err != nil {
			return nil, err
		}

		return databaseschema.NewSchema(config, log, orm, driver, nil)
	})
	app.Singleton(contractsbinding.Seeder, func(app foundation.Application) (any, error) {
		return databaseseeder.NewSeederFacade(), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	r.registerCommands(app)
}

func (r *ServiceProvider) registerCommands(app foundation.Application) {
	artisan := app.MakeArtisan()
	config := app.MakeConfig()
	log := app.MakeLog()
	schema := app.MakeSchema()
	seeder := app.MakeSeeder()

	if artisan != nil && config != nil && log != nil && schema != nil && seeder != nil {
		migrator := migration.NewMigrator(artisan, schema, config.GetString("database.migrations.table"))
		artisan.Register([]contractsconsole.Command{
			consolemigration.NewMigrateMakeCommand(app, migrator),
			consolemigration.NewMigrateCommand(migrator),
			consolemigration.NewMigrateRollbackCommand(migrator),
			consolemigration.NewMigrateResetCommand(migrator),
			consolemigration.NewMigrateRefreshCommand(artisan),
			consolemigration.NewMigrateFreshCommand(artisan, migrator),
			consolemigration.NewMigrateStatusCommand(migrator),
			console.NewModelMakeCommand(artisan, schema),
			console.NewObserverMakeCommand(),
			console.NewSeedCommand(config, seeder),
			console.NewSeederMakeCommand(app),
			console.NewFactoryMakeCommand(),
			console.NewTableCommand(config, schema),
			console.NewShowCommand(config, schema),
			console.NewWipeCommand(config, schema),
		})
	}
}
