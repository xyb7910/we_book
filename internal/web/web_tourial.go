package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Tourial struct{}

func (t *Tourial) RegisterRoutes(server *gin.Engine) {
	server.GET("/ping", func(c *gin.Context) {
		// 输出请求路径
		fmt.Println(c.Request.URL.Path)
		c.String(http.StatusOK, "Hello World")
	})

	// 静态路由匹配
	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "这是 hello 请求的路由返回结果")
	})

	// 参数路由
	server.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "这是你传过来的名字: %s", name)
	})

	// 通配符路由
	server.GET("/views/*.html", func(c *gin.Context) {
		path := c.Param(".html")
		c.String(http.StatusOK, "这是你传过来的路径: %s", path)
	})

	// 查询参数
	server.GET("/order", func(c *gin.Context) {
		// 获取查询参数
		id := c.Query("id")
		c.String(http.StatusOK, "这是你传过来的订单id: %s", id)
	})
}
