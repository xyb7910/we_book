package repository

import (
	"context"
	"time"
	"we_book/internal/domain"
	"we_book/internal/repository/dao"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
	UpdateNextUtime(ctx context.Context, id int64, next time.Time) error
}

type PreemptCronJobRepository struct {
	dao dao.JobDAO
}

func (p *PreemptCronJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		Id:       job.Id,
		Name:     job.Name,
		Executor: job.Executor,

		Cfg: job.Cfg,
	}, nil
}

func (p *PreemptCronJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return p.dao.UpdateUtime(ctx, id)
}

func (p *PreemptCronJobRepository) Release(ctx context.Context, id int64) error {
	return p.dao.Release(ctx, id)
}

func (p *PreemptCronJobRepository) Stop(ctx context.Context, id int64) error {
	return p.dao.Stop(ctx, id)
}

func (p *PreemptCronJobRepository) UpdateNextUtime(ctx context.Context, id int64, next time.Time) error {
	return p.dao.UpdateNextUtime(ctx, id, next)
}

//func NewPreemptionCronJobRepository(dao dao.JobDAO) JobRepository {
//	return &PreemptCronJobRepository{dao: dao}
//}
