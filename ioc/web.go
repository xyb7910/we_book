package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"we_book/internal/web"
	ijwt "we_book/internal/web/jwt"
	"we_book/internal/web/middleware"
)

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	handler *web.OAuth2WeChatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	handler.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	return []gin.HandlerFunc{
		corsMiddleware(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/refresh_token").
			Build(),
		//ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
		sessions.Sessions("my_session", store),
	}
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
