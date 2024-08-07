package main

import (
	"bytes"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	//server := initServer()
	//db := initDB()
	//rdb := initRedis()
	////实现 user 相关路由的注册
	//u := initUser(db, rdb)
	//u.RegisterRoute(server)
	//server := InitWebServer()
	//InitViperV2()
	//InitLogger()
	//_ = server.Run(":8080")

	app := InitWebServer()
	for _, c := range app.consumer {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.web
	server.Run(":8080")
}

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Info("init logger success")
}

func InitViper() {
	//读取的文件名字是 dev
	viper.SetConfigName("dev")
	//读取的文件类型为 yaml
	viper.SetConfigType("yaml")
	//当前工作目录下的 config 文件夹
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperV1() {
	// 设置 db.dsn 默认值
	viper.SetDefault("db.dsn", "root:123456@tcp(127.0.0.1:3306)/we_book?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperV2() {
	cfg := `
		db:
		  dsn: "root:123456789@tcp(127.0.0.1:3306)/we_book?charset=utf8mb4&parseTime=True&loc=Local"
		
		redis:
		  addr: "localhost:6379"
		`
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic("read config error")
	}
}
