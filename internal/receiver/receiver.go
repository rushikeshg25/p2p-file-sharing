package receiver

import (
	"fmt"
	"log"
	"net"
	"os"
	"p2p-file-sharing/internal/protocol"
	"p2p-file-sharing/internal/utils"
	"time"
)

type Receiver struct {
	Address  string
	Port     string
	FileName string
}

func NewReceiver(address string, port string, FileName string) *Receiver {
	return &Receiver{
		Address:  address,
		Port:     port,
		FileName: FileName,
	}
}

const BUFFER_SIZE = 2048

func (r *Receiver) Receive() {
	serverAddress := net.JoinHostPort(r.Address, r.Port)
	fmt.Printf("Connecting to sender at %s...\n", serverAddress)
	conn, err := net.DialTimeout("tcp", serverAddress, 10*time.Second)
	if err != nil {
		log.Fatalf("Error connecting to sender at %s: %v\n", serverAddress, err)
	}
	defer conn.Close()
	fmt.Printf("Connected to %s; waiting for file header...\n", conn.RemoteAddr())

	if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		log.Fatalf("Error setting header timeout: %v\n", err)
	}
	header, err := protocol.Decode(conn)
	if err != nil {
		log.Fatalf("Error reading file header from sender: %v\n", err)
	}
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		log.Fatalf("Error clearing header timeout: %v\n", err)
	}

	fmt.Println("Original filename", header.Name)

	file, err := os.Create(r.FileName)
	if err != nil {
		log.Fatalf("Error creating new file %v\n", err)
	}
	defer file.Close()

	r.receiveFileData(conn, file, int64(header.Size))

	fmt.Println("Checking CRC checksum")

	// Reset file pointer to beginning before CRC calculation
	if _, err := file.Seek(0, 0); err != nil {
		log.Fatalf("Error seeking to beginning of file: %v", err)
	}
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

	progress := utils.NewProgressBar(fileSize, "Receiving")

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
		progress.Add(int64(n))
	}

	progress.Finish()
}
