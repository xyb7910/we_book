package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"we_book/internal/repository"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
	"we_book/internal/service"
	"we_book/internal/service/sms/memory"
	"we_book/internal/web"
	"we_book/internal/web/middleware"
)

func main() {
	server := initServer()
	db := initDB()
	rdb := initRedis()
	//实现 user 相关路由的注册
	u := initUser(db, rdb)
	u.RegisterRoute(server)

	_ = server.Run(":8080")
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
	return redisClient
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

func initServer() *gin.Engine {
	server := gin.Default()
	// 实现跨域问题
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 新建一个 redis 客户端
	//redisClient := initRedis()
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	// 使用 cookie 实现 session
	//store := cookie.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	// 使用 内存 实现 session
	store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	// 使用 redis 实现 session
	//store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), redisClient)
	//if err != nil {
	//	panic(err)
	//}
	// 使用 gin 中间件实现 session
	server.Use(sessions.Sessions("my_session", store))
	// 使用中间件实现登录校验
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		Build())

	//login := middleware.NewLoginMiddlewareBuilder()
	//IgnorePaths := []string{"/users/signup", "/users/login"}
	//login.IgnorePaths(IgnorePaths)
	//server.Use(middleware.CheckLogin())
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
