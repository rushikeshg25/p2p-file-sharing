package sender

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"p2p-file-sharing/internal/protocol"
	"p2p-file-sharing/internal/utils"
)

type Sender struct {
	Port     string
	FileName string
}

const BUFFER_SIZE = 2048 //File chunk size 2KB

func NewSender(port string, FileName string) *Sender {
	return &Sender{
		Port:     port,
		FileName: FileName,
	}
}

func (s *Sender) Send() {
	file, err := os.Open(s.FileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}
	crcVal, err := utils.CalculateCRC(file)
	if err != nil {
		log.Fatalf("Error calculating CRC: %v", err)
	}

	// Reset file pointer to beginning after CRC calculation
	if _, err := file.Seek(0, 0); err != nil {
		log.Fatalf("Error seeking to beginning of file: %v", err)
	}

	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Fatalf("Couldnt start tcp sender")
	}

	fmt.Printf("Server listening on port %s\n", s.Port)
	fmt.Printf("File: %s\n", s.FileName)
	fmt.Printf("CRC32: %08x\n", crcVal)
	fmt.Println("Waiting for receiver...")

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", conn.RemoteAddr())

	header := protocol.FileHeader{
		Protocol: protocol.P2PF_PROTOCOL,
		Version:  protocol.VERSION,
		Size:     uint64(fileInfo.Size()),
		CRC:      crcVal,
		NameLen:  uint8(len(s.FileName)),
		Name:     s.FileName,
	}

	headerBytes, err := header.Encode()
	if err != nil {
		log.Fatalf("Error encoding header: %v", err)
	}

	if _, err := conn.Write(headerBytes); err != nil {
		log.Fatalf("Couldn't send header Bytes")
	}

	fmt.Println("Starting to Send file Chunks")
	s.SendFile(conn, file, fileInfo.Size())

}

func (s *Sender) SendFile(conn net.Conn, file *os.File, size int64) {
	buffer := make([]byte, BUFFER_SIZE)
	var totalSent int64

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading file in bytes %v\n", err)
		}
		if _, err := conn.Write(buffer[:n]); err != nil {
			log.Fatalf("Error seding file bytes to the sender %v\n", err)
		}
		totalSent += int64(n)
	}
	fmt.Println("file sent")
}
