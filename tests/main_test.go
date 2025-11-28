package tests

import (
	"os"
	"testing"

	"github.com/goravel/mysql"
	"github.com/goravel/sqlite"
	"github.com/goravel/sqlserver"
	"github.com/rusmanplatd/goravelframework/foundation/json"
	mocksfoundation "github.com/rusmanplatd/goravelframework/mocks/foundation"
	"github.com/rusmanplatd/goravelpostgres"
)

func TestMain(m *testing.M) {
	mockApp := &mocksfoundation.Application{}
	mockApp.EXPECT().GetJson().Return(json.New())
	goravelpostgres.App = mockApp
	mysql.App = mockApp
	sqlite.App = mockApp
	sqlserver.App = mockApp

	os.Exit(m.Run())
}
