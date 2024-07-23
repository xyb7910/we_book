package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"we_book/internal/repository/dao"
)

func InitDB() *gorm.DB {
	dsn := "root:123456789@tcp(127.0.0.1:3306)/we_book?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
