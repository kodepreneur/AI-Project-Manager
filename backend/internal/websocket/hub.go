package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type EventMessage struct {
	Type      string      `json:"type"`      // "log", "antigravity_output", "process_status", "system_metrics", "notification"
	ProjectID *uint       `json:"project_id,omitempty"`
	Payload   interface{} `json:"payload"`
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan EventMessage
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.RWMutex
}

var GlobalHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan EventMessage, 256),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			data, err := json.Marshal(message)
			if err != nil {
				continue
			}
			h.mutex.RLock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					go func(c *websocket.Conn) {
						h.unregister <- c
					}(client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(msgType string, projectID *uint, payload interface{}) {
	h.broadcast <- EventMessage{
		Type:      msgType,
		ProjectID: projectID,
		Payload:   payload,
	}
}

func (h *Hub) Register(c *websocket.Conn) {
	h.register <- c
}

func (h *Hub) Unregister(c *websocket.Conn) {
	h.unregister <- c
}
