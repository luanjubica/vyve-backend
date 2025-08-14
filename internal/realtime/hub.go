package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/vyve/vyve-backend/pkg/cache"
	"bufio"
)

// Hub maintains active client connections and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// User ID to clients mapping
	userClients map[uuid.UUID][]*Client

	// Inbound messages from clients
	broadcast chan Message

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Redis client for pub/sub
	cache cache.Cache

	// Mutex for concurrent access
	mu sync.RWMutex
}

// Client represents a connected user
type Client struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

// Message represents a real-time message
type Message struct {
	Type      string                 `json:"type"`
	UserID    uuid.UUID              `json:"user_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventType constants
const (
	EventTypeNudge            = "nudge"
	EventTypeInteraction      = "interaction"
	EventTypePersonUpdate     = "person_update"
	EventTypeReflection       = "reflection"
	EventTypeStreakUpdate     = "streak_update"
	EventTypeHealthScore      = "health_score"
	EventTypeNotification     = "notification"
	EventTypePing             = "ping"
	EventTypePong             = "pong"
)

// NewHub creates a new real-time hub
func NewHub(cache cache.Cache) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[uuid.UUID][]*Client),
		broadcast:   make(chan Message, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		cache:       cache,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	// Subscribe to Redis pub/sub for cross-server communication
	go h.subscribeToRedis()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	
	// Add to user clients mapping
	if _, ok := h.userClients[client.UserID]; !ok {
		h.userClients[client.UserID] = []*Client{}
	}
	h.userClients[client.UserID] = append(h.userClients[client.UserID], client)

	log.Printf("Client registered: User %s, Total clients: %d", client.UserID, len(h.clients))
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		// Remove from user clients mapping
		if clients, ok := h.userClients[client.UserID]; ok {
			for i, c := range clients {
				if c == client {
					h.userClients[client.UserID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			if len(h.userClients[client.UserID]) == 0 {
				delete(h.userClients, client.UserID)
			}
		}

		log.Printf("Client unregistered: User %s, Total clients: %d", client.UserID, len(h.clients))
	}
}

// broadcastMessage broadcasts a message to relevant clients
func (h *Hub) broadcastMessage(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	// If message is for a specific user
	if message.UserID != uuid.Nil {
		if clients, ok := h.userClients[message.UserID]; ok {
			for _, client := range clients {
				select {
				case client.Send <- data:
				default:
					// Client's send channel is full, close it
					h.unregisterClient(client)
				}
			}
		}
	} else {
		// Broadcast to all clients
		for client := range h.clients {
			select {
			case client.Send <- data:
			default:
				// Client's send channel is full, close it
				h.unregisterClient(client)
			}
		}
	}

	// Also publish to Redis for cross-server communication
	h.publishToRedis(message)
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID uuid.UUID, messageType string, data map[string]interface{}) {
	message := Message{
		Type:      messageType,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now(),
	}
	h.broadcast <- message
}

// SendToAll sends a message to all connected clients
func (h *Hub) SendToAll(messageType string, data map[string]interface{}) {
	message := Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
	}
	h.broadcast <- message
}

// subscribeToRedis subscribes to Redis pub/sub
func (h *Hub) subscribeToRedis() {
	ctx := context.Background()
	channel, err := h.cache.Subscribe(ctx, "vyve:realtime")
	if err != nil {
		log.Printf("Error subscribing to Redis: %v", err)
		return
	}

	for msg := range channel {
		var message Message
		if err := json.Unmarshal([]byte(msg), &message); err != nil {
			log.Printf("Error unmarshaling Redis message: %v", err)
			continue
		}
		
		// Don't rebroadcast, just send to local clients
		h.sendToLocalClients(message)
	}
}

// publishToRedis publishes a message to Redis
func (h *Hub) publishToRedis(message Message) {
	ctx := context.Background()
	if err := h.cache.Publish(ctx, "vyve:realtime", message); err != nil {
		log.Printf("Error publishing to Redis: %v", err)
	}
}

// sendToLocalClients sends a message to local clients only
func (h *Hub) sendToLocalClients(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	if message.UserID != uuid.Nil {
		if clients, ok := h.userClients[message.UserID]; ok {
			for _, client := range clients {
				select {
				case client.Send <- data:
				default:
					// Client's send channel is full
				}
			}
		}
	}
}

// HandleWebSocket handles WebSocket connections
func (h *Hub) HandleWebSocket(c *fiber.Ctx) error {
	// Upgrade to WebSocket
	return websocket.New(func(conn *websocket.Conn) {
		// Get user ID from context
		userID, ok := c.Locals("user_id").(uuid.UUID)
		if !ok {
			conn.WriteMessage(websocket.CloseMessage, []byte("Unauthorized"))
			conn.Close()
			return
		}

		// Create client
		client := &Client{
			ID:     uuid.New(),
			UserID: userID,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			Hub:    h,
		}

		// Register client
		h.register <- client

		// Start goroutines for reading and writing
		go client.writePump()
		go client.readPump()
	})(c)
}

// HandleSSE handles Server-Sent Events connections
func (h *Hub) HandleSSE(c *fiber.Ctx) error {
	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Get user ID from context
	userVal := c.Locals("user_id")
	userID, ok := userVal.(uuid.UUID)
	if !ok {
		// Sometimes middleware may store as string; attempt to parse
		if s, ok2 := userVal.(string); ok2 {
			if parsed, err := uuid.Parse(s); err == nil {
				userID = parsed
				ok = true
			}
		}
	}
	if !ok {
		return fiber.ErrUnauthorized
	}

	// userID validated; keep reference to avoid unused var warning
	_ = userID

	// Create SSE stream
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Send initial connection message
		fmt.Fprintf(w, "event: connected\ndata: {\"message\": \"Connected to SSE\"}\n\n")
		w.Flush()

		// Create a channel to receive messages
		messages := make(chan Message, 10)
		
		// Subscribe to user's messages
		h.mu.Lock()
		// Note: This is a simplified version. In production, you'd need proper SSE client management
		h.mu.Unlock()

		// Keep connection alive with periodic pings
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg := <-messages:
				data, _ := json.Marshal(msg)
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Type, data)
				w.Flush()

			case <-ticker.C:
				fmt.Fprintf(w, "event: ping\ndata: {\"time\": %d}\n\n", time.Now().Unix())
				w.Flush()

			case <-c.Context().Done():
				return
			}
		}
	})

	return nil
}

// Client methods

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Handle ping
		if msg.Type == EventTypePing {
			pong := Message{
				Type:      EventTypePong,
				Timestamp: time.Now(),
			}
			data, _ := json.Marshal(pong)
			c.Send <- data
			continue
		}

		// Process other message types
		// TODO: Add message processing logic
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}