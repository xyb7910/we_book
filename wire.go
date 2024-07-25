//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"we_book/internal/repository"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
	"we_book/internal/service"
	"we_book/internal/web"
	ijwt "we_book/internal/web/jwt"
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

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,

		web.NewUserHandler,
		web.NewOAuth2WeChatHandler,

		// 基于内存实现存储
		ioc.InitSMSService,
		ioc.InitMiddlewares,

		ioc.InitWebServer,

		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,
	)
	return new(gin.Engine)
}
