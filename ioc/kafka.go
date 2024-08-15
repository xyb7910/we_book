package ioc

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"we_book/events"
	"we_book/events/article"
)

type Config struct {
	Address []string `yaml:"addrs"`
}

func InitKafka() sarama.Client {
	// 初始化 Viper
	viper.SetConfigName("dev")  // 配置文件名 (不带扩展名)
	viper.SetConfigType("yaml") // 配置文件类型
	viper.AddConfigPath(".")    // 配置文件路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// 反序列化配置
	var cfg Config
	if err := viper.UnmarshalKey("kafka", &cfg); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %v", err))
	}

	fmt.Println(cfg) // 确保配置已被正确读取

	// 设置 Sarama 配置
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true

	// 创建 Sarama 客户端
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
