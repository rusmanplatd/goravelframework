package tests

import (
	"context"

	contractsorm "github.com/rusmanplatd/goravelframework/contracts/database/orm"
	databaseorm "github.com/rusmanplatd/goravelframework/database/orm"
	"github.com/rusmanplatd/goravelframework/database/schema"
	"github.com/rusmanplatd/goravelframework/testing/utils"
)

const (
	testDatabase = "goravel"
	testUsername = "goravel"
	testPassword = "Framework!123"
	testSchema   = "goravel"
)

func newSchema(testQuery *TestQuery, connectionToTestQuery map[string]*TestQuery) *schema.Schema {
	queries := make(map[string]contractsorm.Query)
	for connection, testQuery := range connectionToTestQuery {
		queries[connection] = testQuery.Query()
	}

	log := utils.NewTestLog()
	dbConfig := testQuery.Driver().Pool().Writers[0]
	orm := databaseorm.NewOrm(context.Background(), testQuery.Config(), dbConfig.Connection, dbConfig, testQuery.Query(), queries, log, nil, nil)

	schema, err := schema.NewSchema(testQuery.Config(), log, orm, testQuery.Driver(), nil)
	if err != nil {
		log.Panic(err.Error())
	}

	return schema
}
