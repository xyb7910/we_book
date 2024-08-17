//go:build wireinject

package main

import (
	"github.com/google/wire"
	article "we_book/events/article"
	"we_book/internal/repository"
	article2 "we_book/internal/repository/article"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
	article3 "we_book/internal/repository/dao/article"
	"we_book/internal/service"
	"we_book/internal/web"
	ijwt "we_book/internal/web/jwt"
	"we_book/ioc"
)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCacheInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

var rankingServerProvider = wire.NewSet(
	repository.NewRankingRepository,
	cache.NewRankingRedisCache,
	service.NewBatchRankingService,
)

func InitWebServer() *App {
	wire.Build(
		// 首先引入最基本的第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,

		ioc.InitKafka,

		interactiveSvcProvider,
		rankingServerProvider,
		ioc.InitRankingJob,
		ioc.InitJobs,
		ioc.InitRLockClient,

		ioc.NewSyncProducer,
		ioc.NewConsumers,

		// consumer
		article.NewInteractiveReadEventBatchConsumer,
		article.NewKafkaProducer,

		// 初始化 DAO
		dao.NewUserDAO,
		article3.NewGORMArticleDAO,

		cache.NewUserCache,
		cache.NewRedisCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		article2.NewArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// 基于内存实现存储
		ioc.InitSMSService,
		ioc.InitMiddlewares,

		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WeChatHandler,

		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
