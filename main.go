package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: p2p-share send|receive <file> <port>")
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Invalid port %s", os.Args[3])
		os.Exit(1)
	}

	file := os.Args[2]

	switch os.Args[1] {
	case "send":
		fmt.Printf("Sending file %s to port %d\n", file, port)
	case "receive":
		fmt.Printf("Receiving file %s from port %d\n", file, port)
	default:
		fmt.Printf("Invalid command %s\n", os.Args[1])
		os.Exit(1)
	}
}
