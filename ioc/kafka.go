package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"we_book/events"
	"we_book/events/article"
)

func InitKafka() sarama.Client {
	type Config struct {
		Address []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Address, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

func NewConsumers(c1 *article.InteractiveReadEventBatchConsumer) []events.Consumer {
	return []events.Consumer{c1}
}
