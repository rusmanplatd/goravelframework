package migration

import (
	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/console/command"
	"github.com/rusmanplatd/goravelframework/contracts/database/migration"
	"github.com/rusmanplatd/goravelframework/errors"
)

type MigrateCommand struct {
	migrator migration.Migrator
}

func NewMigrateCommand(migrator migration.Migrator) *MigrateCommand {
	return &MigrateCommand{
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
	return command.Extend{}
}

// Handle Execute the console command.
func (r *MigrateCommand) Handle(ctx console.Context) error {
	if err := r.migrator.Run(); err != nil {
		ctx.Error(errors.MigrationMigrateFailed.Args(err).Error())
		return nil
	}

	ctx.Success("Migration success")

	return nil
}
