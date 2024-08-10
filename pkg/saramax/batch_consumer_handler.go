package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
	"we_book/pkg/logger"
)

type BatchConsumerHandler[T any] struct {
	l             logger.V1
	fn            func(msg []*sarama.ConsumerMessage, t []T) error
	bathSize      int
	batchDuration time.Duration
}

func NewBatchConsumerHandler[T any](l logger.V1, fn func(msg []*sarama.ConsumerMessage, t []T) error, bathSize int, batchDuration time.Duration) *BatchConsumerHandler[T] {
	return &BatchConsumerHandler[T]{
		l:             l,
		fn:            fn,
		bathSize:      bathSize,
		batchDuration: batchDuration,
	}
}

func (b *BatchConsumerHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchConsumerHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchConsumerHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()
	batchSize := b.bathSize
	for {
		// 凑一批处理的时候，要注意超时
		ctx, cancel := context.WithTimeout(context.Background(), b.batchDuration)
		done := false
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgsCh:
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化失败",
						logger.Error(err),
						logger.String("topic", msg.Topic),
						logger.Int64("partition", int64(msg.Partition)),
						logger.Int64("offset", msg.Offset))
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		if err := b.fn(msgs, ts); err != nil {
			b.l.Error("处理消息失败",
				logger.Error(err))
		}
		for _, msg := range msgs {
			session.MarkMessage(msg, "")
		}
	}
}
