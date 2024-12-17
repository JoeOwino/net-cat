package main

import (
	"fmt"
	"os"

	"net-cat/commands"
)

func main() {
	err := commands.CreateTCPServer()
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
