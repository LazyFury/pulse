package framework

import "github.com/robfig/cron/v3"

type IJob interface {
	Job()
	Spec() string
}

type Cron struct {
	*cron.Cron
}

func NewCorn() *Cron {
	return &Cron{
		Cron: cron.New(),
	}
}

func (c *Cron) Start() {
	c.Cron.Run()
}

func (c *Cron) AddJob(job IJob) {
	c.Cron.AddFunc(job.Spec(), job.Job)
}

func (c *Cron) Stop() {
	c.Cron.Stop()
}
