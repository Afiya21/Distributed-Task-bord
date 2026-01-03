package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

// Client represents a connected user
type Client struct {
	Hub    *Hub
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	Clients    map[string]*Client // Map UserID to Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	mu         sync.Mutex
}

type Message struct {
	UserID  string      `json:"userId"`
	Content interface{} `json:"content"`
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("User connected: %s", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Printf("User disconnected: %s", client.UserID)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.Lock()
			client, ok := h.Clients[message.UserID]
			if ok {
				select {
				case client.Send <- toJson(message.Content):
				default:
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
			h.mu.Unlock()
		}
	}
}

func toJson(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
