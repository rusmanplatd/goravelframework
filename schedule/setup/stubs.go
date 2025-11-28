package main

type Stubs struct{}

func (s Stubs) ScheduleFacade() string {
	return `package facades

import (
	"github.com/rusmanplatd/goravelframework/contracts/schedule"
)

func Schedule() schedule.Schedule {
	return App().MakeSchedule()
}
`
}
