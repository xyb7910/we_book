package job

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
	"we_book/pkg/logger"
)

type RankingJobAdapter struct {
	j Job
	l logger.V1
	p prometheus.Summary
}

func NewRankingJobAdapter(j Job, l logger.V1) *RankingJobAdapter {
	p := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "corn_job",
		ConstLabels: map[string]string{
			"name": j.Name(),
		},
	})
	prometheus.MustRegister(p)
	return &RankingJobAdapter{l: l, j: j}
}

func (r *RankingJobAdapter) Run() {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		r.p.Observe(float64(duration))
	}()
	err := r.j.Run()
	if err != nil {
		r.l.Error("job run error",
			logger.String("job", r.j.Name()))
	}
}
