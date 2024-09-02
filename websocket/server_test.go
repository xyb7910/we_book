package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
	"time"
)

type WsServer struct {
	*websocket.Conn
}

func TestServer(t *testing.T) {
	upgrader := &websocket.Upgrader{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Write([]byte("upgrade error"))
			return
		}
		conn := &WsServer{Conn: c}

		go func() {
			for {
				typ, msg, err := conn.ReadMessage()
				if err != nil {
					return
				}
				switch typ {
				case websocket.CloseMessage:
					conn.Close()
					return
				default:
					t.Logf("msg:%s", msg)
				}
			}
		}()

		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for now := range ticker.C {
				err := conn.WriteMessage(websocket.TextMessage, []byte("Hello"+now.String()))
				if err != nil {
					return
				}
			}
		}()
	})
	http.ListenAndServe(":8081", nil)
}
