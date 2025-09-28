package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/denizsincar29/airtype/internal/airtype"
	"golang.org/x/term"
)

func main() {
	// Create and connect the AirType client
	client := airtype.NewClient()
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
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		ch := buf[0]
		// 3 = Ctrl+C, 27 = Esc
		if ch == 3 || ch == 27 {
			return
		}

		if err := client.TypeChar(ch); err != nil {
			log.Printf("Error sending character: %v", err)
			// If sending fails, we might have disconnected.
			// The reconnection logic will be added in the next step.
		}
	}
}