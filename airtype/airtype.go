package airtype

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

var (
	DELETE = []byte("#del$")
	ENTER  = []byte("\n")
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
	logger         *slog.Logger
}

// NewClient creates a new AirType client.
func NewClient(ip string, logger *slog.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	// Create a logger with file if not provided
	if logger == nil {
		logFile, err := os.OpenFile("airtype.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		logger = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{AddSource: true}))
	}
	// Validate the IP address
	ip, err := IP(ip)
	if err != nil {
		logger.Error("Invalid IP address", slog.String("error", err.Error()))
	}

	return &Client{
		ctx:    ctx,
		cancel: cancel,
		ip:     ip,
		logger: logger,
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

	c.logger.Info("Attempting to reconnect to AirType server...")
	if c.conn != nil {
		c.conn.Close(websocket.StatusAbnormalClosure, "connection lost")
	}

	for {
		if c.ctx.Err() != nil {
			c.logger.Info("Reconnection aborted due to context cancellation.")
			return
		}

		conn, _, err := websocket.Dial(c.ctx, c.URL, nil)
		if err != nil {
			c.logger.Error("Reconnection attempt failed", slog.String("error", err.Error()))
			time.Sleep(5 * time.Second)
			continue
		}

		c.mu.Lock()
		c.conn = conn
		c.mu.Unlock()

		c.logger.Info("Reconnected to AirType server.")
		_, msg, err := c.conn.Read(c.ctx)
		if err != nil {
			c.logger.Error("Failed to read after reconnect", slog.String("error", err.Error()))
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

// Delete sends a delete command to the AirType server.
func (c *Client) Delete() error {
	return c.Write(DELETE)
}

// Enter sends an enter command to the AirType server.
func (c *Client) Enter() error {
	return c.Write(ENTER)
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

// GetLogger returns the client's logger.
func (c *Client) GetLogger() *slog.Logger {
	return c.logger
}

func IP(Ip string) (string, error) {
	// if the string ends with txt, read the file
	if strings.HasSuffix(Ip, ".txt") {
		data, err := os.ReadFile(Ip)
		if err != nil {
			return "", fmt.Errorf("failed to read IP from file: %w", err)
		}
		Ip = strings.TrimSpace(string(data))
	}
	// if the string has only 2 ip parts, add 192.168. at the beginning
	parts := strings.Split(Ip, ".")
	if len(parts) == 2 {
		Ip = "192.168." + Ip
	}
	// validate the ip address
	if net.ParseIP(Ip) == nil {
		return "", fmt.Errorf("invalid IP address: %s", Ip)
	}
	return Ip, nil
}
