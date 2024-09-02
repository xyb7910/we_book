package link_pool

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}
	pool := &ConnectionPool{
		connections: make(map[*websocket.Conn]bool),
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "无法升级到 WebSocket", http.StatusInternalServerError)
		return
	}

	// 添加新的连接到连接池
	pool.Add(conn)
	defer pool.Remove(conn)

	// 读取消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取消息错误: %v", err)
			break
		}
		log.Printf("接收到消息: %s", message)

		// 广播消息到所有连接
		pool.Broadcast(messageType, message)
	}
}
