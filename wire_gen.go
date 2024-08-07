// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	article3 "we_book/events/article"
	"we_book/internal/repository"
	article2 "we_book/internal/repository/article"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
	"we_book/internal/repository/dao/article"
	"we_book/internal/service"
	"we_book/internal/web"
	"we_book/internal/web/jwt"
	"we_book/ioc"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	v1 := ioc.InitLogger()
	v := ioc.InitMiddlewares(cmdable, handler, v1)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	articleDAO := article.NewGORMArticleDAO(db)
	articleRepository := article2.NewArticleRepository(articleDAO)
	client := ioc.InitKafka()
	syncProducer := ioc.NewSyncProducer(client)
	producer := article3.NewKafkaProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, v1, producer)
	articleHandler := web.NewArticleHandler(articleService, v1)
	wechatService := ioc.InitWechatService(v1)
	wechatHandlerConfig := ioc.NewWechatHandlerConfig()
	oAuth2WeChatHandler := web.NewOAuth2WeChatHandler(wechatService, userService, handler, wechatHandlerConfig)
	engine := ioc.InitWebServer(v, userHandler, articleHandler, oAuth2WeChatHandler)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveRepository := repository.NewCacheReadCntRepository(interactiveCache, interactiveDAO, v1)
	interactiveReadEventConsumer := article3.NewInteractiveReadEventConsumer(client, interactiveRepository, v1)
	v2 := ioc.NewConsumers(interactiveReadEventConsumer)
	app := &App{
		web:      engine,
		consumer: v2,
	}
	return app
}
