package queue

import (
	contractsdb "github.com/rusmanplatd/goravelframework/contracts/database/db"
	contractsfoundation "github.com/rusmanplatd/goravelframework/contracts/foundation"
	contractslog "github.com/rusmanplatd/goravelframework/contracts/log"
	contractsqueue "github.com/rusmanplatd/goravelframework/contracts/queue"
	"github.com/rusmanplatd/goravelframework/errors"
)

type DriverCreator struct {
	config    contractsqueue.Config
	db        contractsdb.DB
	jobStorer contractsqueue.JobStorer
	json      contractsfoundation.Json
	log       contractslog.Log
}

func NewDriverCreator(config contractsqueue.Config, db contractsdb.DB, jobStorer contractsqueue.JobStorer, json contractsfoundation.Json, log contractslog.Log) *DriverCreator {
	return &DriverCreator{
		config:    config,
		db:        db,
		jobStorer: jobStorer,
		json:      json,
		log:       log,
	}
}

func (r *DriverCreator) Create(connection string) (contractsqueue.Driver, error) {
	driver := r.config.Driver(connection)

	switch driver {
	case contractsqueue.DriverSync:
		return NewSync(), nil
	case contractsqueue.DriverDatabase:
		if r.db == nil {
			return nil, errors.QueueInvalidDatabaseConnection.Args(connection)
		}

		return NewDatabase(r.config, r.db, r.jobStorer, r.json, connection)
	case contractsqueue.DriverCustom:
		custom := r.config.Via(connection)
		if driver, ok := custom.(contractsqueue.Driver); ok {
			return driver, nil
		}
		if driver, ok := custom.(func() (contractsqueue.Driver, error)); ok {
			return driver()
		}
		return nil, errors.QueueDriverInvalid.Args(connection)
	default:
		return nil, errors.QueueDriverNotSupported.Args(driver)
	}
}
