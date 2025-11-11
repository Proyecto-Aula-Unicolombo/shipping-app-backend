package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/saveblush/gofiber3-contrib/websocket"
)

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mu         *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mu:         &sync.Mutex{},
	}
}

func (hub *Hub) HandleWebSocketConnection(ctx fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		log.Println("Conexión WebSocket establecida con:", c.RemoteAddr())

		client := NewClient(c, hub)
		if addr := c.RemoteAddr(); addr != nil {
			client.id = addr.String()
		}

		hub.register <- client

		go client.WritePump()
		client.ReadPump()
	}, websocket.Config{
		Origins: []string{"*"},
	})(ctx)
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	addrStr := client.id
	if addrStr == "" && client.conn != nil {
		if ra := client.conn.RemoteAddr(); ra != nil {
			addrStr = ra.String()
		}
	}

	log.Println("Client Connected", addrStr)

	hub.mu.Lock()
	defer hub.mu.Unlock()
	client.id = addrStr
	hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client Disconnected", client.conn.RemoteAddr())
	hub.mu.Lock()
	defer hub.mu.Unlock()
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
		}
	}
	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]
}

// helper para enviar un mensaje a todos
func (hub *Hub) BroadcastJSON(message interface{}, ignore *Client) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("marshal broadcast error:", err)
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for _, client := range hub.clients {
		if client != ignore {
			select {
			case client.send <- data:
			default:
				close(client.send)
			}
		}
	}
}
