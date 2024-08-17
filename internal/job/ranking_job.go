package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"sync"
	"time"
	"we_book/internal/service"
	"we_book/pkg/logger"
)

type RankingJob struct {
	svc       service.RankingService
	timeout   time.Duration
	client    *rlock.Client
	key       string
	l         logger.V1
	lock      *rlock.Lock
	localLock *sync.Mutex
}

func NewRankingJob(svc service.RankingService,
	client *rlock.Client,
	l logger.V1,
	timeout time.Duration) *RankingJob {
	return &RankingJob{
		svc:       svc,
		timeout:   timeout,
		client:    client,
		l:         l,
		key:       "ranking_job",
		localLock: &sync.Mutex{},
	}
}

func (r *RankingJob) Run() error {
	r.localLock.Lock()
	defer r.localLock.Unlock()
	if r.lock == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      0,
		}, time.Second)
		if err != nil {
			return nil
		}
		r.lock = lock

		go func() {
			r.localLock.Lock()
			defer r.localLock.Unlock()
			err1 := lock.AutoRefresh(r.timeout/2, time.Second)

			if err1 != nil {
				r.l.Error("ranking job auto refresh error")
			}
			r.lock = nil
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.lock = nil
	r.localLock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}

func (r *RankingJob) Name() string {
	return "ranking_job"
}
