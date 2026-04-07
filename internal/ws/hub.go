package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// writeDeadline is the maximum time allowed to write a message to the client.
	writeDeadline = 54 * time.Second

	// readDeadline is the maximum time allowed to read the next message from the client.
	// Resets on each ping/pong exchange.
	readDeadline = 60 * time.Second

	// pingInterval is how often to send a ping frame to keep the connection alive.
	pingInterval = 45 * time.Second

	// sendBufferSize is the number of outbound messages to buffer per client before
	// the client is considered too slow and is dropped.
	sendBufferSize = 256
)

var upgrader = websocket.Upgrader{
	// Allow connections from any origin — the hub lives on the local network
	// and we don't want to fight CORS when phones and tablets connect.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// BroadcastEvent is the envelope sent over WebSocket to every connected client.
//
// Type identifies what changed. Defined types:
//
//	"kasa:state"         — TP-Link Kasa device state update
//	                       Payload: { device_id, alias, relay_state, brightness? }
//	"hue:state"          — Philips Hue light state update
//	                       Payload: { device_id, name, on, brightness?, color_temp? }
//	"device:discovered"  — a new device was found on the network
//	                       Payload: { device_id, name, type, ip }
//	"device:offline"     — a device stopped responding
//	                       Payload: { device_id }
//	"layout:updated"     — the generated UI layout changed
//	                       Payload: { layout JSON }
//	"command:result"     — result of a command sent to a device
//	                       Payload: { device_id, success, error? }
type BroadcastEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Client represents a single active WebSocket connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub manages all connected WebSocket clients and routes broadcasts to them.
// Create one with NewHub and start it with go hub.Run().
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

// NewHub returns an initialised Hub ready to be started with Run.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

// Run starts the Hub's main event loop. Must be called in a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("[ws] client connected (%d total)", h.clientCount())

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("[ws] client disconnected (%d remaining)", len(h.clients))
			}
			h.mu.Unlock()

		case data := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					// Client's send buffer is full — it's too slow to keep up.
					// Drop it now rather than blocking the whole broadcast.
					delete(h.clients, client)
					close(client.send)
					log.Printf("[ws] dropped slow client")
				}
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastGlobal marshals event to JSON and sends it to every connected client.
// It is safe to call from multiple goroutines (e.g. the Kasa and Hue pollers).
func (h *Hub) BroadcastGlobal(event BroadcastEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("[ws] marshal broadcast: %v", err)
		return
	}

	// Send to the broadcast channel rather than iterating clients directly —
	// the Run loop owns the map and handles slow-client cleanup safely.
	select {
	case h.broadcast <- data:
	default:
		log.Printf("[ws] broadcast channel full, event dropped: %s", event.Type)
	}
}

// HandleWebSocket upgrades an HTTP connection to WebSocket, registers the new
// client with the Hub, and starts its read/write pumps.
//
// Mount this at your WebSocket endpoint, e.g. GET /api/ws
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade: %v", err)
		return
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, sendBufferSize),
	}

	h.register <- client

	// Each pump runs in its own goroutine. writePump owns the write side of the
	// connection; readPump owns the read side (and detects disconnection).
	go client.writePump()
	go client.readPump()
}

// readPump drains incoming frames from the client connection. Its main job is
// to detect disconnection (and respond to pong frames) so the client can be
// cleanly unregistered from the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(readDeadline))
	c.conn.SetPongHandler(func(string) error {
		// Each pong resets the read deadline, keeping the connection alive.
		c.conn.SetReadDeadline(time.Now().Add(readDeadline))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			// Normal close or network error — either way, we're done.
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("[ws] unexpected close: %v", err)
			}
			break
		}
	}
}

// writePump drains the client's send channel and writes frames to the WebSocket
// connection. It also sends periodic ping frames to keep the connection alive
// and detect dead clients before the read deadline expires.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if !ok {
				// Hub closed the channel — send a clean close frame and exit.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[ws] write: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[ws] ping: %v", err)
				return
			}
		}
	}
}

// clientCount returns the number of currently registered clients.
// Acquires a read lock — safe to call from any goroutine.
func (h *Hub) clientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
