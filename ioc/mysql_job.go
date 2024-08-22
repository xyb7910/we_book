package ioc

import (
	"context"
	"time"
	"we_book/internal/domain"
	"we_book/internal/job"
	"we_book/internal/service"
	"we_book/pkg/logger"
)

func InitScheduler(l logger.V1, local *job.LocalFuncExecutor, svc service.JobService) *job.Scheduler {
	res := job.NewScheduler(svc, l)
	res.RegisterExecutor(local)
	return res
}

func InitLocalFuncExecutor(svc service.RankingService) *job.LocalFuncExecutor {
	res := job.NewLocalFuncExecutor()
	res.RegisterFunc("ranking", func(ctx context.Context, j domain.Job) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*20)
		defer cancel()
		return svc.TopN(ctx)
	})
	return res
}
