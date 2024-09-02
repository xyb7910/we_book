package link_pool

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
)

func Test_wsHandler(t *testing.T) {
	http.HandleFunc("/ws", wsHandler)

	go func() {
		server := gin.Default()
		server.GET("/", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "Hello, WebSocket!")
		})
		server.Run(":8082")
	}()

	log.Fatal(http.ListenAndServe(":8081", nil))
}
