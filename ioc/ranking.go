package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/robfig/cron/v3"
	"time"
	"we_book/internal/job"
	"we_book/internal/service"
	"we_book/pkg/logger"
)

func InitRankingJob(svc service.RankingService,
	rlockClient *rlock.Client,
	l logger.V1) *job.RankingJob {
	return job.NewRankingJob(svc, rlockClient, l, time.Second*30)
}

func InitJobs(l logger.V1, rankingJob *job.RankingJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	cbd := job.NewCronJobBuilder(l)
	_, err := res.AddJob("0 */3 * * * ?", cbd.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return res
}
