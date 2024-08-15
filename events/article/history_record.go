package article

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"we_book/internal/repository"
	"we_book/pkg/logger"
	"we_book/pkg/saramax"
)

type HistoryReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.V1
}

func NewHistoryReadEventConsumer(client sarama.Client, repo repository.InteractiveRepository, l logger.V1) *HistoryReadEventConsumer {
	return &HistoryReadEventConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (r *HistoryReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(), []string{"history_read"}, saramax.NewHandler[ReadEvent](r.l, r.Consumer))
		if err != nil {
			r.l.Error("consumer error")
		}
	}()
	return err
}

func (r *HistoryReadEventConsumer) Consumer(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.repo.AddRecord(ctx, t.Aid, t.Uid)
}
