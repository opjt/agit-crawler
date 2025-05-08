package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type SampleJob struct{}

func NewSampleJob() Job {
	return &SampleJob{}
}

func (j *SampleJob) Register(c *cron.Cron) {
	c.AddFunc("@every 10s", func() {
		fmt.Println("[SampleJob] Running every 10 seconds at", time.Now())
	})
}
