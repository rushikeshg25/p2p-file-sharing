package main

import (
	"fmt"
	"os"
	"p2p-file-sharing/internal/receiver"
	"p2p-file-sharing/internal/sender"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "send":
		if len(os.Args) != 4 {
			printUsage()
			os.Exit(1)
		}
		s := sender.NewSender(os.Args[3], os.Args[2])
		s.Send()
	case "receive":
		address := "localhost"
		fileArg := 2
		portArg := 3
		if len(os.Args) == 5 {
			address = os.Args[2]
			fileArg = 3
			portArg = 4
		} else if len(os.Args) != 4 {
			printUsage()
			os.Exit(1)
		}
		r := receiver.NewReceiver(address, os.Args[portArg], os.Args[fileArg])
		r.Receive()
	default:
		fmt.Printf("Invalid command %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  p2p-share send <file> <port>")
	fmt.Println("  p2p-share receive [sender-ip] <output-file> <port>")
}
