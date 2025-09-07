package receiver

import (
	"fmt"
	"log"
	"net"
	"os"
	"p2p-file-sharing/internal/protocol"
	"p2p-file-sharing/internal/utils"
)

type Receiver struct {
	Port     string
	FileName string
}

func NewReceiver(port string, FileName string) *Receiver {
	return &Receiver{
		Port:     port,
		FileName: FileName,
	}
}

const BUFFER_SIZE = 2048

func (r *Receiver) Receive() {
	conn, err := net.Dial("tcp", r.Port)
	if err != nil {
		log.Fatalf("Error connecting to the sender %v\n", err)
	}
	defer conn.Close()

	header, err := protocol.Decode(conn)
	if err != nil {
		log.Fatalf("Error reading header%v\n", err)
	}

	fmt.Println("Original filename", header.Name)

	file, err := os.Create(r.FileName)
	if err != nil {
		log.Fatalf("Error creating new file %v\n", err)
	}
	defer file.Close()

	r.receiveFileData(conn, file, int64(header.Size))

	fmt.Println("Checking CRC checksum")
	receivedFileCRC, err := utils.CalculateCRC(file)
	if err != nil {
		log.Fatalf("Error calculating file for received File %v\n", err)
	}

	if receivedFileCRC == header.CRC {
		fmt.Println("CRC verified")
		fmt.Println("p2p Transfer done")
	} else {
		fmt.Println("CRC mismatch")
	}

}

func (r *Receiver) receiveFileData(conn net.Conn, file *os.File, fileSize int64) {
	buffer := make([]byte, BUFFER_SIZE)
	var totalReceived int64

	for totalReceived < fileSize {
		remainingBytesSize := fileSize - totalReceived
		bufSize := BUFFER_SIZE
		if remainingBytesSize < int64(BUFFER_SIZE) {
			bufSize = int(remainingBytesSize)
		}

		n, err := conn.Read(buffer[:bufSize])
		if err != nil {
			log.Fatalf("Error reading from the stream %v\n", err)
		}

		if _, err := file.Write(buffer[:n]); err != nil {
			log.Fatalf("Error writing to file %v\n", err)
		}

		totalReceived += int64(n)
	}
}
