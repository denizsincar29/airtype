package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/denizsincar29/airtype/airtype"
)

func main() {
	// Configure logging
	logFile, err := os.OpenFile("airtype.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Define a flag for the input file
	filePath := flag.String("file", "text.txt", "Path to the text file to type.")
	flag.Parse()

	// Check if the file exists
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatalf("Error: File '%s' not found. Please create it or specify a different file with the -file flag.", *filePath)
	}

	// Read the content of the file
	content, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file '%s': %v", *filePath, err)
	}

	// Get IP and create client
	ip := getIPAddress()
	client := airtype.NewClient(ip)
	if err := client.Connect(); err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer client.Close()

	fmt.Printf("Typing content from '%s'...\n", *filePath)

	// Send the entire content at once
	if err := client.Write(content); err != nil {
		log.Fatalf("Error sending content: %v", err)
	}

	fmt.Println("\nFinished typing.")
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