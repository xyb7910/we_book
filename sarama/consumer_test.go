package sarama

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs, "test_group", cfg)
	assert.NoError(t, err)
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	time.AfterFunc(time.Second*11, func() {
		cancel()
	})
	err = consumer.Consume(ctx, []string{"test_topic"}, &testConsumerGroupHandler{})
	t.Log(err, time.Since(start).String())
}

type testConsumerGroupHandler struct {
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("setup")
	return nil
}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("clean up")
	return nil
}

func (t testConsumerGroupHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

func (t testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var eg errgroup.Group
		var last *sarama.ConsumerMessage
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				last = msg
				eg.Go(func() error {
					time.Sleep(time.Second)
					log.Println(string(msg.Value))
					return nil
				})
			}
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			continue
		}
		if last != nil {
			session.MarkMessage(last, "")
		}
	}

}
