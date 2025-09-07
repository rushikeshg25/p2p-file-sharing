package main

import (
	"fmt"
	"os"
	"p2p-file-sharing/internal/receiver"
	"p2p-file-sharing/internal/sender"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: p2p-share send|receive <file> <port>")
		os.Exit(1)
	}

	port := os.Args[3]
	file := os.Args[2]

	switch os.Args[1] {
	case "send":
		s := sender.NewSender(port, file)
		s.Send()
	case "receive":
		r := receiver.NewReceiver(port, file)
		r.Receive()
	default:
		fmt.Printf("Invalid command %s\n", os.Args[1])
		os.Exit(1)
	}
}
