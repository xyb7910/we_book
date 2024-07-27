package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	return db
}
