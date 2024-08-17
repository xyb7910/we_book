package job

import (
	"github.com/prometheus/client_golang/prometheus"
	corn "github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
	"strconv"
	"time"
	"we_book/pkg/logger"
)

type CronJobBuilder struct {
	l      logger.V1
	p      *prometheus.SummaryVec
	tracer trace.Tracer
}

type cronJobFuncAdaptor func() error

func (c cronJobFuncAdaptor) Run() {
	_ = c()
}

func NewCronJobBuilder(l logger.V1) *CronJobBuilder {
	p := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "we_book",
		Subsystem: "job",
		Help:      "job summary",
		Name:      "cron_job",
	}, []string{"name", "success"})
	prometheus.MustRegister(p)
	return &CronJobBuilder{
		l:      l,
		p:      p,
		tracer: otel.GetTracerProvider().Tracer("we_book/internal/job"),
	}
}

func (b *CronJobBuilder) Build(job Job) corn.Job {
	name := job.Name()
	return cronJobFuncAdaptor(func() error {
		_, span := b.tracer.Start(context.Background(), name)
		defer span.End()
		start := time.Now()
		b.l.Info("job start", logger.String("name", name))
		var success bool
		defer func() {
			b.l.Info("job end", logger.String("name", name))
			duration := time.Since(start).Milliseconds()
			b.p.WithLabelValues(name, strconv.FormatBool(success)).Observe(float64(duration))
		}()
		err := job.Run()
		success = err == nil
		if err != nil {
			span.RecordError(err)
			b.l.Error("job error", logger.String("name", name), logger.Error(err))
		}
		return err
	})
}
