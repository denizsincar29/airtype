package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
	"nhooyr.io/websocket"
)

func main() {
	// Catch Ctrl+C and exit gracefully
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nExiting (Ctrl+C pressed)")
		cancel()
		os.Exit(0)
	}()

	// Read or ask for IP address
	ip := ""
	if data, err := os.ReadFile("ip.txt"); err == nil {
		ip = strings.TrimSpace(string(data))
	} else {
		fmt.Print("Enter your iPhone's airtype IP address: ")
		reader := bufio.NewReader(os.Stdin)
		ip, _ = reader.ReadString('\n')
		ip = strings.TrimSpace(ip)
		os.WriteFile("ip.txt", []byte(ip), 0644)
	}

	// Connect to WebSocket
	url := fmt.Sprintf("ws://%s:8307/service", ip)
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer c.Close(websocket.StatusNormalClosure, "bye")

	fmt.Println("Connected to server:", url)

	// Read initial message
	_, msg, err := c.Read(ctx)
	if err != nil {
		log.Fatalf("Failed to read first message: %v", err)
	}
	fmt.Println("Received:", string(msg))

	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to set raw mode: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Println("Terminal in raw mode. Press Esc to exit.")

	var wordCount uint8 = 0
	buf := make([]byte, 1)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		ch := buf[0]

		switch ch {
		case 27: // Escape key
			fmt.Println("\nExiting (Esc pressed)")
			return
		case 13: // Enter
			err = c.Write(ctx, websocket.MessageText, []byte("\n"))
		case 127, 8: // Backspace / DEL
			err = c.Write(ctx, websocket.MessageText, []byte("#del$"))
		default:
			if ch == ' ' {
				wordCount++
			} else {
				wordCount = 0
			}
			err = c.Write(ctx, websocket.MessageText, []byte{ch})
		}

		if err != nil {
			log.Printf("Failed to send: %v", err)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
