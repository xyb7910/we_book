package main

import (
	"github.com/gin-gonic/gin"
	"we_book/internal/web"
)

func main() {
	server := gin.Default()
	// 实现跨域问题
	//server.Use(cors.New(cors.Config{
	//	AllowCredentials: true,
	//	AllowHeaders:     []string{"Content-Type"},
	//	AllowOriginFunc: func(origin string) bool {
	//		if strings.HasPrefix(origin, "http://localhost") {
	//			return true
	//		}
	//		return strings.Contains(origin, "your_company.com")
	//	},
	//	MaxAge: 12 * time.Hour,
	//}))
	// 实现 user 相关路由的注册
	user := web.NewUserHandler()
	user.RegisterRoute(server)

	_ = server.Run(":8080")
}
