package article

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"we_book/internal/repository"
	"we_book/pkg/logger"
	"we_book/pkg/saramax"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.V1
}

func NewInteractiveReadEventBatchConsumer(client sarama.Client, repo repository.InteractiveRepository, l logger.V1) *InteractiveReadEventBatchConsumer {
	return &InteractiveReadEventBatchConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (r *InteractiveReadEventBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewBatchConsumerHandler[ReadEvent](r.l, r.Consume, 10, time.Second))
		if err != nil {
			r.l.Error("循环消费消息失败", logger.Error(err))
		}
	}()
	return err
}

func (r *InteractiveReadEventBatchConsumer) Consume(msgs []*sarama.ConsumerMessage, ts []ReadEvent) error {
	ids := make([]int64, 0, len(ts))
	bizs := make([]string, 0, len(ts))
	for _, evt := range ts {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := r.repo.BatchIncrReadCnt(ctx, ids, bizs)
	if err != nil {
		r.l.Error("批量更新阅读数失败",
			logger.Field{Key: "ids", Value: ids},
			logger.Error(err))
	}
	return nil
}
