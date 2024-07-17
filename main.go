package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"we_book/internal/repository"
	"we_book/internal/repository/dao"
	"we_book/internal/service"
	"we_book/internal/web"
)

func main() {
	server := initServer()
	db := initDB()

	// 实现 user 相关路由的注册
	u := initUser(db)
	u.RegisterRoute(server)

	_ = server.Run(":8080")
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initServer() *gin.Engine {
	server := gin.Default()
	// 实现跨域问题
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	return server
}

func initDB() *gorm.DB {
	dsn := "root:123456789@tcp(127.0.0.1:3306)/we_book?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 只会在初始化是 panic， 一旦初始化过程出错，应用无法启动
		panic(err)
	}
	err = dao.InitTable(DB)
	if err != nil {
		panic(err)
	}
	return DB
}
