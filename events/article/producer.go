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

type ReadEventV1 struct {
	Uids []int64
	Aids []int64
}

type Producer interface {
	ProducerReadEvent(ctx context.Context, event ReadEvent) error
	ProducerReadEventV1(ctx context.Context, info ReadEventV1)
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func (k *KafkaProducer) ProducerReadEventV1(ctx context.Context, info ReadEventV1) {
	// TODO: implement me
	panic("implement me")
}

func (k *KafkaProducer) ProducerReadEvent(ctx context.Context, event ReadEvent) error {
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
