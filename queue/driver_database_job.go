package queue

import (
	contractsdb "github.com/rusmanplatd/goravelframework/contracts/database/db"
	contractsfoundation "github.com/rusmanplatd/goravelframework/contracts/foundation"
	contractsqueue "github.com/rusmanplatd/goravelframework/contracts/queue"
	"github.com/rusmanplatd/goravelframework/queue/models"
	"github.com/rusmanplatd/goravelframework/queue/utils"
)

type DatabaseReservedJob struct {
	db        contractsdb.DB
	job       *models.Job
	jobsTable string
	task      contractsqueue.Task
}

func NewDatabaseReservedJob(job *models.Job, db contractsdb.DB, jobStorer contractsqueue.JobStorer, json contractsfoundation.Json, jobsTable string) (*DatabaseReservedJob, error) {
	task, err := utils.JsonToTask(job.Payload, jobStorer, json)
	if err != nil {
		return nil, err
	}

	return &DatabaseReservedJob{
		db:        db,
		job:       job,
		jobsTable: jobsTable,
		task:      task,
	}, nil
}

func (r *DatabaseReservedJob) Delete() error {
	_, err := r.db.Table(r.jobsTable).Where("id", r.job.ID).Delete()

	return err
}

func (r *DatabaseReservedJob) Task() contractsqueue.Task {
	return r.task
}
