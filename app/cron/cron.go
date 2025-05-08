package cron

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

type Job interface {
	Register(*cron.Cron)
}

type CronManager struct {
	cron *cron.Cron
}

func (c CronManager) Start() {
	c.cron.Start()
}
func (c CronManager) Stop() {
	c.cron.Stop()
}

func NewCronManager(p struct {
	fx.In
	Cron *cron.Cron
	Jobs []Job `group:"cron.jobs"`
}) CronManager {
	for _, job := range p.Jobs {
		job.Register(p.Cron)
	}
	return CronManager{p.Cron}
}

var Module = fx.Options(
	fx.Provide(
		func() *cron.Cron {
			return cron.New()
		},
	),
	fx.Provide(
		fx.Annotate(
			NewSampleJob,
			fx.ResultTags(`group:"cron.jobs"`),
		),
		fx.Annotate(
			NewAgitJob,
			fx.ResultTags(`group:"cron.jobs"`),
		),
	),
	fx.Provide(NewCronManager),
)
