package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/denizsincar29/airtype/airtype"
	"golang.org/x/term"
)

func main() {
	// Configure logging
	logFile, err := os.OpenFile("airtype.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Get IP and create client
	ip := getIPAddress()
	client := airtype.NewClient(ip)
	if err := client.Connect(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer client.Close()

	// Handle Ctrl+C gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nExiting...")
		client.Close()
		os.Exit(0)
	}()

	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to set raw mode: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Println("Terminal in raw mode. Press Ctrl+C or Esc to exit.")

	// Read from terminal and send to client
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Reading runes...")
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			log.Printf("Error reading rune: %v", err)
			return
		}

		// 3 = Ctrl+C, 27 = Esc
		if r == 3 || r == 27 {
			return
		}

		var data []byte
		switch r {
		case '\n', '\r': // Enter
			data = airtype.ENTER
		case 127, 8: // Backspace / DEL
			data = airtype.DELETE
		default:
			data = []byte(string(r))
		}

		if err := client.Write(data); err != nil {
			log.Printf("Error sending character: %v", err)
			// Reconnection logic is handled by the client's Write method
		}
	}
}

// getIPAddress reads the IP from a file or prompts the user for it.
func getIPAddress() string {
	if data, err := os.ReadFile("ip.txt"); err == nil {
		return strings.TrimSpace(string(data))
	}

	fmt.Print("Enter your iPhone's AirType IP address: ")
	reader := bufio.NewReader(os.Stdin)
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)
	if err := os.WriteFile("ip.txt", []byte(ip), 0644); err != nil {
		log.Printf("Warning: failed to save IP address to ip.txt: %v", err)
	}
	return ip
}
