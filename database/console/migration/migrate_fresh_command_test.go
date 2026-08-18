package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goravel/framework/errors"
	mocksconsole "github.com/goravel/framework/mocks/console"
	mocksconfig "github.com/goravel/framework/mocks/config"
	mocksmigration "github.com/goravel/framework/mocks/database/migration"
)

func TestMigrateFreshCommand(t *testing.T) {
	var (
		mockArtisan  *mocksconsole.Artisan
		mockConfig   *mocksconfig.Config
		mockContext  *mocksconsole.Context
		mockMigrator *mocksmigration.Migrator
	)

	beforeEach := func() {
		mockArtisan = mocksconsole.NewArtisan(t)
		mockConfig = mocksconfig.NewConfig(t)
		mockContext = mocksconsole.NewContext(t)
		mockMigrator = mocksmigration.NewMigrator(t)
	}

	tests := []struct {
		name  string
		setup func()
	}{
		{
			name: "Happy path",
			setup: func() {
				mockContext.EXPECT().Option("path").Return("").Once()
				mockContext.EXPECT().Option("schema").Return("").Once()
				mockMigrator.EXPECT().Fresh().Return(nil).Once()
				mockContext.EXPECT().OptionBool("seed").Return(true).Once()
				mockContext.EXPECT().OptionSlice("seeder").Return([]string{"UserSeeder", "AgentSeeder"}).Once()
				mockArtisan.EXPECT().Call("db:seed --seeder UserSeeder,AgentSeeder").Return(nil).Once()
				mockContext.EXPECT().Success("Migration fresh success").Once()
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
				mockMigrator.EXPECT().Fresh().Return(nil).Once()
				mockContext.EXPECT().OptionBool("seed").Return(false).Once()
				mockConfig.EXPECT().Add("database.connections.postgres.schema", "public").Once()
				mockContext.EXPECT().Success("Migration fresh success").Once()
			},
		},
		{
			name: "Sad path - fresh failed",
			setup: func() {
				mockContext.EXPECT().Option("path").Return("").Once()
				mockContext.EXPECT().Option("schema").Return("").Once()
				mockMigrator.EXPECT().Fresh().Return(assert.AnError).Once()
				mockContext.EXPECT().Error(errors.MigrationFreshFailed.Args(assert.AnError).Error()).Once()
			},
		},
		{
			name: "Sad path - call db:seed failed",
			setup: func() {
				mockContext.EXPECT().Option("path").Return("").Once()
				mockContext.EXPECT().Option("schema").Return("").Once()
				mockMigrator.EXPECT().Fresh().Return(nil).Once()
				mockContext.EXPECT().OptionBool("seed").Return(true).Once()
				mockContext.EXPECT().OptionSlice("seeder").Return([]string{"UserSeeder", "AgentSeeder"}).Once()
				mockArtisan.EXPECT().Call("db:seed --seeder UserSeeder,AgentSeeder").Return(assert.AnError).Once()
				mockContext.EXPECT().Error(errors.MigrationFreshFailed.Args(assert.AnError).Error()).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			beforeEach()
			test.setup()

			command := NewMigrateFreshCommand(mockArtisan, mockMigrator, mockConfig)
			err := command.Handle(mockContext)

			assert.NoError(t, err)
		})
	}
}
