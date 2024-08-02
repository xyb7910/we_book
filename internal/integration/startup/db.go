package startup

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"we_book/internal/repository/dao"
)

//func InitDB() *gorm.DB {
//	dsn := "root:123456789@tcp(127.0.0.1:3306)/we_book?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		panic(err)
//	}
//	err = dao.InitTable(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}

	// 在结构体中使用默认的 dsn
	c := Config{
		DSN: "root:123456789@tcp(127.0.0.1:3306)/we_book?charset=utf8mb4&parseTime=True&loc=Local",
	}

	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic("db config error")
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		panic("db connect error")
	}
	err = dao.InitTable(db)
	if err != nil {
		panic("db init table error")
	}
	return db
}

var mongoDB *mongo.Database

func InitMongoDB() *mongo.Database {
	if mongoDB == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context,
				startedEvent *event.CommandStartedEvent) {
				fmt.Println(startedEvent.Command)
			},
		}
		opts := options.Client().
			ApplyURI("mongodb://root:example@localhost:27017/").
			SetMonitor(monitor)
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			panic(err)
		}
		mongoDB = client.Database("webook")
	}
	return mongoDB
}
