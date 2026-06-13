package studio

import (
	"sync"

	"github.com/gorilla/websocket"
)

type client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub manages active WebSocket connections and fans out broadcast messages.
type Hub struct {
	mu         sync.Mutex
	clients    map[*client]struct{}
	broadcastC chan []byte
	registerC  chan *client
	unregisterC chan *client
}

func newHub() *Hub {
	return &Hub{
		clients:     make(map[*client]struct{}),
		broadcastC:  make(chan []byte, 16),
		registerC:   make(chan *client, 8),
		unregisterC: make(chan *client, 8),
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.registerC:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			h.mu.Unlock()

		case c := <-h.unregisterC:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcastC:
			h.mu.Lock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					delete(h.clients, c)
					close(c.send)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) broadcast(msg []byte) {
	h.broadcastC <- msg
}

func (c *client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
