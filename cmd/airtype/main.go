package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/denizsincar29/airtype/airtype"
	"github.com/denizsincar29/goerror"
	"golang.design/x/clipboard"
	"golang.org/x/term"
)

func main() {
	// Define flags
	is_clipboard := flag.Bool("c", false, "Use clipboard input")
	filename := flag.String("file", "", "Path to the text file")
	ip := flag.String("ip", "ip.txt", "IP address")
	flag.Parse()
	if *ip == "" {
		fmt.Println("IP address is required. Use -ip flag to specify it.")
		return
	}

	client := airtype.NewClient(*ip, nil) // nil logger means default logger
	logger := client.GetLogger()
	e := goerror.NewError(logger)
	err := client.Connect()
	e.Must(err, "Failed to connect to AirType device at %s", *ip)
	defer client.Close()

	// initialize clipboard if needed
	err = clipboard.Init()
	e.Must(err, "Failed to initialize clipboard")

	// If a file is provided, read and send its content
	if *filename != "" || *is_clipboard {
		var content []byte
		if *is_clipboard {
			content = clipboard.Read(clipboard.FmtText)
		} else {
			content, err = os.ReadFile(*filename)
			e.Must(err, "Failed to read file: %s", *filename)
		}
		err = client.Write(content)
		e.Must(err, "Failed to send file content to AirType device")
		fmt.Println("File or clipboard content sent successfully.")
		return
	}

	// If no file is provided, enter interactive mode
	fmt.Println("Entering interactive mode. Type your input:")
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
	e.Must(err, "Failed to set terminal to raw mode")
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Println("Terminal in raw mode. Press Ctrl+C or Esc to exit.")

	// Read from terminal and send to client
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Reading runes...")
	for {
		r, _, err := reader.ReadRune()
		e.Check(err, "Failed to read rune from terminal")
		if e.IsError() { // saved the last error
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
		err = client.Write(data)
		e.Check(err, "Failed to send character to AirType device")
	}
}
