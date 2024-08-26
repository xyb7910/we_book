package startup

import (
	"context"
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	"we_book/interactive/repository/dao"
)

var db *gorm.DB

func InitDB() *gorm.DB {
	if db == nil {
		dsn := "root:123456789@tcp(127.0.0.1:3306)/we_book_test?charset=utf8mb4&parseTime=True&loc=Local"
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err = sqlDB.PingContext(ctx)
			cancel()
			if err == nil {
				break
			}
			log.Println("db ping error, retrying...")
		}
		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			panic(err)
		}
		err = dao.InitTable(db)
		if err != nil {
			panic(err)
		}
		db = db.Debug()
	}
	return db
}
