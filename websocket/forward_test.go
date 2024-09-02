package websocket

import (
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
)

type Hub struct {
	// 封装的 map key 为房间号，value 为房间内的所有连接
	conns *syncx.Map[string, *websocket.Conn]
}

func (h *Hub) AddConn(name string, conn *websocket.Conn) {
	h.conns.Store(name, conn)
	go func() {
		// 接收消息
		typ, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		switch typ {
		case websocket.CloseMessage:
			h.conns.Delete(name)
			conn.Close()
			return
		default:
			// 广播消息
			log.Println("from client:", typ, string(msg), name)
			h.conns.Range(func(key string, value *websocket.Conn) bool {
				if key == name {
					// 不发送给自己
					return true
				}
				log.Println("to client:", key)
				err := value.WriteMessage(typ, msg)
				if err != nil {
					log.Println(err)
				}
				return true
			})
		}
	}()
}

func TestForward(t *testing.T) {
	upgrader := websocket.Upgrader{}
	hub := &Hub{
		conns: &syncx.Map[string, *websocket.Conn]{},
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Write([]byte("upgrade error"))
			return
		}
		name := r.URL.Query().Get("name")
		hub.AddConn(name, c)
	})
	// ws://localhost:8081/ws?name=ypb
	http.ListenAndServe(":8081", nil)
}
