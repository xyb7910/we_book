package startup

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

//func InitRedis() redis.Cmdable {
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//	})
//	return redisClient
//}

func InitRedis() redis.Cmdable {
	// 我们需要一个 redis 的客户端
	cmd := redis.NewClient(&redis.Options{
		// 使用 GetString 方法获取配置文件中的 redis 地址
		Addr: viper.GetString("redis.addr"),
	})
	return cmd
}
