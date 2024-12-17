package commands

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func CreateTCPServer() error {
	// Start listening on the specified address
	port := "8989"
	args := os.Args
	n := len(args)

	if n == 2 {
		port = args[1]
	}

	if n > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}

	listener, err := net.Listen("tcp", "localhost: "+port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", port, err)
	}
	defer listener.Close()
	fmt.Printf("Listening on the port :%s\n", port)

	connCount := 0
	for connCount < 10 {
		// Accept a client connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		connCount++

		fmt.Println("New client connected:", conn.RemoteAddr())

		// Handle the client in a separate goroutine
		go handleClient(conn)
	}

	return nil
}

func getCientName(conn net.Conn) string {
	defer conn.Close()

	logo, err := os.ReadFile("logo.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = conn.Write(logo)
	if err != nil {
		fmt.Println("Error writing writing logo:", err)
		return ""
	}

	_, err = conn.Write([]byte("[ENTER YOUR NAME]: "))
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return ""
	}

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		return name
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading name from client:", err)
	}
	return ""
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	clientName := getCientName(conn)

	// Read data from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// Echo back the received data
		message := scanner.Text()
		if message == "" {
			continue
		}

		fmt.Printf("%s: %s\n", clientName + ":", message)

		_, err := conn.Write([]byte("Echo: " + message + "\n"))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}

	// Handle any scanner error
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from client:", err)
	}
	fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
}
