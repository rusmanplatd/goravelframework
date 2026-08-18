package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goravel/framework/errors"
	mocksconsole "github.com/goravel/framework/mocks/console"
	mocksconfig "github.com/goravel/framework/mocks/config"
	mocksmigration "github.com/goravel/framework/mocks/database/migration"
)

func TestMigrateCommand(t *testing.T) {
	var (
		mockConfig   *mocksconfig.Config
		mockContext  *mocksconsole.Context
		mockMigrator *mocksmigration.Migrator
	)

	beforeEach := func() {
		mockConfig = mocksconfig.NewConfig(t)
		mockContext = mocksconsole.NewContext(t)
		mockMigrator = mocksmigration.NewMigrator(t)
	}

	tests := []struct {
		name  string
		setup func()
	}{
		{
			name: "Happy path - no flags",
			setup: func() {
				mockContext.EXPECT().Option("path").Return("").Once()
				mockContext.EXPECT().Option("schema").Return("").Once()
				mockMigrator.EXPECT().Run().Return(nil).Once()
				mockContext.EXPECT().Success("Migration success").Once()
			},
		},
		{
			name: "Happy path - with schema flag",
			setup: func() {
				mockContext.EXPECT().Option("path").Return("").Once()
				mockContext.EXPECT().Option("schema").Return("myschema").Once()
				mockConfig.EXPECT().GetString("database.default").Return("postgres").Once()
				mockConfig.EXPECT().GetString("database.connections.postgres.schema").Return("public").Once()
				mockConfig.EXPECT().Add("database.connections.postgres.schema", "myschema").Once()
				mockMigrator.EXPECT().Run().Return(nil).Once()
				mockConfig.EXPECT().Add("database.connections.postgres.schema", "public").Once()
				mockContext.EXPECT().Success("Migration success").Once()
			},
		},
		{
			name: "Sad path - run failed",
			setup: func() {
				mockContext.EXPECT().Option("path").Return("").Once()
				mockContext.EXPECT().Option("schema").Return("").Once()
				mockMigrator.EXPECT().Run().Return(assert.AnError).Once()
				mockContext.EXPECT().Error(errors.MigrationMigrateFailed.Args(assert.AnError).Error()).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			beforeEach()
			test.setup()

			command := NewMigrateCommand(mockMigrator, mockConfig)
			err := command.Handle(mockContext)

			assert.NoError(t, err)
		})
	}
}
