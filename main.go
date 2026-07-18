package main

import (
	"fmt"
	"log"
	"os"
	"p2p-file-sharing/internal/receiver"
	"p2p-file-sharing/internal/sender"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		printUsage()
		return fmt.Errorf("missing command")
	}

	switch args[0] {
	case "send":
		if len(args) != 3 {
			printUsage()
			return fmt.Errorf("send requires a file and port")
		}
		s := sender.NewSender(args[2], args[1])
		return s.Send()
	case "receive":
		address := "localhost"
		fileArg := 1
		portArg := 2
		if len(args) == 4 {
			address = args[1]
			fileArg = 2
			portArg = 3
		} else if len(args) != 3 {
			printUsage()
			return fmt.Errorf("receive requires an output file and port, with an optional sender IP")
		}
		r := receiver.NewReceiver(address, args[portArg], args[fileArg])
		return r.Receive()
	default:
		printUsage()
		return fmt.Errorf("invalid command %q", args[0])
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  p2p-share send <file> <port>")
	fmt.Println("  p2p-share receive [sender-ip] <output-file> <port>")
}
