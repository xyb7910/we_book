package article

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"we_book/internal/repository"
	"we_book/pkg/logger"
	"we_book/pkg/saramax"
)

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.V1
}

func NewInteractiveReadEventConsumer(
	client sarama.Client,
	repo repository.InteractiveRepository,
	l logger.V1) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (r *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive-read-event-consumer", r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewHandler[ReadEvent](r.l, r.Consume),
		)
		if err != nil {
			r.l.Error("consumer error", logger.Error(err))
		}
	}()
	return err
}

func (r *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.repo.IncrReadCnt(ctx, "article", t.Aid)
}
