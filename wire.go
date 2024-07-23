package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
	"we_book/internal/service"
	"we_book/internal/web"
	"we_book/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 首先引入最基本的第三方依赖
		ioc.InitDB,
		ioc.InitRedis,

		// 初始化 DAO
		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewRedisCodeCache,

		service.NewUserService,
		service.NewCodeService,

		// 基于内存实现存储
		ioc.InitSMSService,
		web.NewUserHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
