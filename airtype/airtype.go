package airtype

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

// Client handles the connection and communication with the AirType server.
type Client struct {
	conn           *websocket.Conn
	ctx            context.Context
	cancel         context.CancelFunc
	URL            string
	ip             string
	mu             sync.Mutex
	isReconnecting bool
}

// NewClient creates a new AirType client.
func NewClient(ip string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ctx:    ctx,
		cancel: cancel,
		ip:     ip,
	}
}

// Connect establishes a WebSocket connection to the AirType server.
func (c *Client) Connect() error {
	c.URL = fmt.Sprintf("ws://%s:8307/service", c.ip)

	var err error
	c.conn, _, err = websocket.Dial(c.ctx, c.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", c.URL, err)
	}

	fmt.Println("Connected to server:", c.URL)

	// Read initial message
	_, msg, err := c.conn.Read(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to read first message: %w", err)
	}
	fmt.Println("Received:", string(msg))
	return nil
}

func (c *Client) reconnect() {
	c.mu.Lock()
	if c.isReconnecting {
		c.mu.Unlock()
		return
	}
	c.isReconnecting = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.isReconnecting = false
		c.mu.Unlock()
	}()

	log.Println("Connection lost. Attempting to reconnect...")
	if c.conn != nil {
		c.conn.Close(websocket.StatusAbnormalClosure, "connection lost")
	}

	for {
		if c.ctx.Err() != nil {
			log.Println("Context cancelled, stopping reconnection.")
			return
		}

		conn, _, err := websocket.Dial(c.ctx, c.URL, nil)
		if err != nil {
			log.Printf("Failed to reconnect: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		c.mu.Lock()
		c.conn = conn
		c.mu.Unlock()

		log.Println("Reconnected successfully!")

		_, msg, err := c.conn.Read(c.ctx)
		if err != nil {
			log.Printf("Failed to read first message after reconnect: %v", err)
			c.conn.Close(websocket.StatusAbnormalClosure, "post-reconnect read failed")
			continue
		}
		fmt.Println("Received (after reconnect):", string(msg))
		return
	}
}

// Write sends a message to the AirType server and handles reconnection.
func (c *Client) Write(data []byte) error {
	c.mu.Lock()
	conn := c.conn
	isReconnecting := c.isReconnecting
	c.mu.Unlock()

	if conn == nil || isReconnecting {
		return fmt.Errorf("not connected or reconnecting")
	}

	err := conn.Write(c.ctx, websocket.MessageText, data)
	if err != nil {
		go c.reconnect()
		return fmt.Errorf("write failed, triggering reconnect: %w", err)
	}
	return nil
}

// Close gracefully closes the connection.
func (c *Client) Close() {
	c.cancel() // This will stop reconnection attempts
	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close(websocket.StatusNormalClosure, "bye")
	}
	c.mu.Unlock()
}

// TypeChar sends a single character or a special command.
func (c *Client) TypeChar(ch byte) error {
	var data []byte
	switch ch {
	case '\n', '\r': // Enter
		data = []byte("\n")
	case 127, 8: // Backspace / DEL
		data = []byte("#del$")
	default:
		data = []byte{ch}
	}
	return c.Write(data)
}