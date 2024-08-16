package web

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type ObservabilityHandler struct {
}

func (h *ObservabilityHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/observability")
	g.GET("/metric", func(ctx *gin.Context) {
		sleep := rand.Int31n(1000)
		time.Sleep(time.Millisecond * time.Duration(sleep))
		ctx.String(http.StatusOK, "Ok")
	})
}
