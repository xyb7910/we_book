//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"we_book/internal/repository"
	"we_book/internal/repository/article"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
	article2 "we_book/internal/repository/dao/article"
	"we_book/internal/service"
	"we_book/internal/web"
	ijwt "we_book/internal/web/jwt"
	"we_book/ioc"
)

var thirdProviderSet = wire.NewSet(InitRedis, InitLogger, InitDB)
var userSvcProviderSet = wire.NewSet(dao.NewUserDAO, cache.NewUserCache, repository.NewUserRepository, service.NewUserService, web.NewUserHandler)
var codeSvcProviderSet = wire.NewSet(cache.NewRedisCodeCache, repository.NewCodeRepository, service.NewCodeService)
var articleSvcProviderSet = wire.NewSet(article2.NewGORMArticleDAO, article.NewArticleRepository, service.NewArticleService, web.NewArticleHandler)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 首先引入最基本的第三方依赖
		thirdProviderSet,
		// 用户服务
		userSvcProviderSet,
		// 验证码服务
		codeSvcProviderSet,
		// 文章服务
		articleSvcProviderSet,

		InitSMSService,
		InitPhoneToWechatService,
		NewWechatHandlerConfig,
		web.NewOAuth2WeChatHandler,
		ijwt.NewRedisJWTHandler,

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	return new(gin.Engine)
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdProviderSet,
		articleSvcProviderSet,
	)
	return new(web.ArticleHandler)
}
