package job

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"time"
	"we_book/internal/domain"
	"we_book/internal/service"
	"we_book/pkg/logger"
)

type Executor interface {
	// Name Executor 名字
	Name() string
	Exec(ctx context.Context, job domain.Job) error
}

type Scheduler struct {
	execs   map[string]Executor
	svc     service.JobService
	l       logger.V1
	limiter *semaphore.Weighted
}

func NewScheduler(svc service.JobService, l logger.V1) *Scheduler {
	return &Scheduler{
		svc:     svc,
		l:       l,
		limiter: semaphore.NewWeighted(200),
		execs:   make(map[string]Executor),
	}
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.execs[exec.Name()] = exec
}

func (s *Scheduler) Scheduler(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		job, err := s.svc.Preempt(dbCtx)
		cancel()

		if err != nil {
			s.l.Error("Preempt job error", logger.Error(err))
		}

		exec, ok := s.execs[job.Executor]
		if !ok {
			s.l.Error("Executor not found", logger.String("executor", job.Executor))
			continue
		}

		go func() {
			defer func() {
				s.limiter.Release(1)
				err1 := job.CancelFunc()
				if err1 != nil {
					s.l.Error("release job error", logger.Error(err1), logger.Int64("job_id", job.Id))
				}
			}()
			err1 := exec.Exec(ctx, job)
			if err1 != nil {
				s.l.Error("exec job error", logger.Error(err1), logger.Int64("job_id", job.Id))
			}

			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err1 = s.svc.ResetNextTime(ctx, job)
			if err1 != nil {
				s.l.Error("reset job next time error", logger.Error(err1), logger.Int64("job_id", job.Id))
			}
		}()
	}
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.Job) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: make(map[string]func(ctx context.Context, j domain.Job) error),
	}
}

func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
	l.funcs[name] = fn
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, job domain.Job) error {
	fn, ok := l.funcs[job.Name]
	if !ok {
		return fmt.Errorf("unknown job name %s", job.Name)
	}
	return fn(ctx, job)
}
