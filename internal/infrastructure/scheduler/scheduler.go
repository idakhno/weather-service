package scheduler

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	s gocron.Scheduler
}

func New() (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &Scheduler{s: s}, nil
}

func (sch *Scheduler) Every(d time.Duration, job func()) error {
	_, err := sch.s.NewJob(
		gocron.DurationJob(d),
		gocron.NewTask(job),
	)
	return err
}

func (sch *Scheduler) Start() {
	sch.s.Start()
}
