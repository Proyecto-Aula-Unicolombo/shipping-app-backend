package ws

import (
	"log"
	"time"

	"github.com/saveblush/gofiber3-contrib/websocket"
)

const (
	writeWait      = 30 * time.Second    // max time to wait
	pongWait       = 60 * time.Second    // max time to inactive
	pingPeriod     = (pongWait * 9) / 10 // frecuence to send a ping
	maxMessageSize = 512
)

type Client struct {
	hub  *Hub
	id   string
	conn *websocket.Conn
	send chan []byte
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  hub,
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod) // tempo to send ping
	defer ticker.Stop()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // limit time to write
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage) // get writer to write message
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // send ping, and virify if the client is disconeted
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))                                                           // limit time of inactive
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil }) // when the client response this function keep the conection

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) { //  if the client is disconnect this end the for
				log.Printf("error de lectura: %v", err)
			}
			break
		}
	}
}
