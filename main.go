package main

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"we_book/internal/repository"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
	"we_book/internal/service"
	"we_book/internal/service/sms/memory"
	"we_book/internal/web"
)

func main() {
	//server := initServer()
	//db := initDB()
	//rdb := initRedis()
	////实现 user 相关路由的注册
	//u := initUser(db, rdb)
	//u.RegisterRoute(server)
	server := InitWebServer()

	_ = server.Run(":8080")
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}
