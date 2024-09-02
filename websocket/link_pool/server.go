package link_pool

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// Ws websocket
type Ws struct {
	Conn *websocket.Conn
}

// ConnectionPool connection pool
type ConnectionPool struct {
	connections map[*websocket.Conn]bool
	lock        sync.Mutex
}

// Add new connection
func (p *ConnectionPool) Add(conn *websocket.Conn) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.connections[conn] = true
}

// Remove connection
func (p *ConnectionPool) Remove(conn *websocket.Conn) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, ok := p.connections[conn]; ok; {
		delete(p.connections, conn)
		conn.Close()
	}
}

// Broadcast message
func (p *ConnectionPool) Broadcast(messageType int, message []byte) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for conn := range p.connections {
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("发送消息错误: %v", err)
			conn.Close()
			delete(p.connections, conn)
		}
	}
}
