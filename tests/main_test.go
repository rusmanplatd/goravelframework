package tests

import (
	"os"
	"testing"

	"github.com/goravel/mysql"
	"github.com/goravel/postgres"
	"github.com/goravel/sqlite"
	"github.com/goravel/sqlserver"
	"github.com/rusmanplatd/goravelframework/foundation/json"
	mocksfoundation "github.com/rusmanplatd/goravelframework/mocks/foundation"
)

func TestMain(m *testing.M) {
	mockApp := &mocksfoundation.Application{}
	mockApp.EXPECT().GetJson().Return(json.New())
	postgres.App = mockApp
	mysql.App = mockApp
	sqlite.App = mockApp
	sqlserver.App = mockApp

	os.Exit(m.Run())
}
