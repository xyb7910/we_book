package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type ReadEvent struct {
	Uid int64
	Aid int64
}

type Producer interface {
	ProducerReadEvent(ctx context.Context, event ReadEvent) error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func (k KafkaProducer) ProducerReadEvent(ctx context.Context, event ReadEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_article",
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func NewKafkaProducer(producer sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: producer,
	}
}
