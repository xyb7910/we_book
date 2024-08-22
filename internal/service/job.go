package service

import (
	"context"
	"time"
	"we_book/internal/domain"
	"we_book/internal/repository"
	"we_book/pkg/logger"
)

type JobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, j domain.Job) error
}

type cronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	l               logger.V1
}

func (c *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		for range ticker.C {
			c.refresh(job.Id)
		}
	}()
	// 续约之后，要考虑释放的问题
	job.CancelFunc = func() error {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return c.repo.Release(ctx, job.Id)
	}
	return job, err
}

func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 更新 job 的时间
	err := c.repo.UpdateUtime(ctx, id)
	if err != nil {
		c.l.Error("续约失败",
			logger.Error(err),
			logger.Int64("job_id", id),
		)
	}
}

func (c *cronJobService) ResetNextTime(ctx context.Context, j domain.Job) error {
	next := j.NextTime()
	if next.IsZero() {
		return c.repo.Stop(ctx, j.Id)
	}
	return c.repo.UpdateNextUtime(ctx, j.Id, next)
}

//func NewCronJobService(repo repository.JobRepository, refreshInterval time.Duration, l logger.V1) JobService {
//	return &cronJobService{
//		repo:            repo,
//		refreshInterval: refreshInterval,
//		l:               l,
//	}
//}
