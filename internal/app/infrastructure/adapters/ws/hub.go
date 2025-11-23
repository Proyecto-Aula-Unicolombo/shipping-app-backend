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
	// Verificar si es una solicitud de upgrade de WebSocket
	if !websocket.IsWebSocketUpgrade(ctx) {
		log.Println("⚠️  No es una solicitud de upgrade de WebSocket")
		return fiber.NewError(fiber.StatusUpgradeRequired, "Se requiere WebSocket upgrade")
	}

	log.Println("🔌 Intentando establecer conexión WebSocket desde:", ctx.IP())

	return websocket.New(func(c *websocket.Conn) {
		log.Println("✅ Conexión WebSocket establecida con:", c.RemoteAddr())

		client := NewClient(c, hub)
		if addr := c.RemoteAddr(); addr != nil {
			client.id = addr.String()
		}

		hub.register <- client

		// Enviar mensaje de bienvenida
		welcomeMsg := WebSocketMessage{
			Type:    "connected",
			Payload: map[string]string{"message": "WebSocket connected successfully"},
		}
		if data, err := json.Marshal(welcomeMsg); err == nil {
			c.WriteMessage(websocket.TextMessage, data)
		}

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

// BroadcastToRole - Enviar mensaje solo a clientes con un rol específico
func (hub *Hub) BroadcastToRole(message interface{}, role string) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("marshal broadcast error:", err)
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for _, client := range hub.clients {
		if client.role == role {
			select {
			case client.send <- data:
			default:
				close(client.send)
			}
		}
	}
}

// BroadcastToOrder - Enviar mensaje a clientes siguiendo una orden específica
func (hub *Hub) BroadcastToOrder(message interface{}, orderID uint) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("marshal broadcast error:", err)
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for _, client := range hub.clients {
		// Enviar a admins y a clientes siguiendo esta orden
		if client.role == "admin" || hub.isFollowingOrder(client, orderID) {
			select {
			case client.send <- data:
			default:
				close(client.send)
			}
		}
	}
}

// isFollowingOrder - Verificar si un cliente está siguiendo una orden
func (hub *Hub) isFollowingOrder(client *Client, orderID uint) bool {
	for _, id := range client.orderIDs {
		if id == orderID {
			return true
		}
	}
	return false
}

// SendToClient - Enviar mensaje a un cliente específico por ID
func (hub *Hub) SendToClient(message interface{}, clientID string) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for _, client := range hub.clients {
		if client.id == clientID {
			select {
			case client.send <- data:
			default:
				close(client.send)
			}
			break
		}
	}
}
