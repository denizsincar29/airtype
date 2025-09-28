package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/denizsincar29/airtype/internal/airtype"
)

func main() {
	// Define a flag for the input file
	filePath := flag.String("file", "text.txt", "Path to the text file to type.")
	flag.Parse()

	// Check if the file exists
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatalf("Error: File '%s' not found. Please create it or specify a different file with the -file flag.", *filePath)
	}

	// Read the content of the file
	content, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file '%s': %v", *filePath, err)
	}

	// Create and connect the AirType client
	client := airtype.NewClient()
	if err := client.Connect(); err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer client.Close()

	fmt.Printf("Typing content from '%s'...\n", *filePath)

	// Type out the content character by character
	for _, char := range content {
		if err := client.TypeChar(byte(char)); err != nil {
			log.Printf("Error sending character: %v", err)
		}
		// Add a small delay to simulate typing and prevent overwhelming the server
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("\nFinished typing.")
}