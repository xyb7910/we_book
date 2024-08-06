package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var addrs = []string{"localhost:9094"}

func TestSyncProduct(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Key:   sarama.StringEncoder("oid-123"),
		Value: sarama.StringEncoder("test-value"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "application/json",
	})
	assert.NoError(t, err)
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	require.NoError(t, err)
	msgCh := producer.Input()

	go func() {
		for {
			msg := &sarama.ProducerMessage{
				Topic: "test_topic",
				Key:   sarama.StringEncoder("oid-123"),
				Value: sarama.StringEncoder("test-value"),
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("trace_id"),
						Value: []byte("123456"),
					},
				},
				Metadata: "application/json",
			}
			select {
			case msgCh <- msg:
			}
		}
	}()

	errCh := producer.Errors()
	succCh := producer.Successes()

	for {
		select {
		case err := <-errCh:
			t.Logf("err: %v", err)
		case <-succCh:
			t.Logf("succ")
		}
	}
}

type JSONEncoder struct {
	Data any
}
