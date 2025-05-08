package cron

import (
	"agit-crawler/app/pkg"

	"github.com/robfig/cron/v3"
)

type AgitJob struct {
	crawler *pkg.Crawler
}

func NewAgitJob(
	crawler *pkg.Crawler,
) Job {
	return &AgitJob{
		crawler,
	}
}

func (j *AgitJob) Register(c *cron.Cron) {

	c.AddFunc("@every 5s", j.crawler.GetPosts)
}
