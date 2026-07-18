package sender

import (
	"fmt"
	"io"
	"net"
	"os"
	"p2p-file-sharing/internal/protocol"
	"p2p-file-sharing/internal/utils"
	"path/filepath"
	"time"
)

type Sender struct {
	Port     string
	FileName string
}

const BUFFER_SIZE = 2048 //File chunk size 2KB
const transferIdleTimeout = 30 * time.Second

func NewSender(port string, FileName string) *Sender {
	return &Sender{
		Port:     port,
		FileName: FileName,
	}
}

func (s *Sender) Send() error {
	if err := utils.ValidatePort(s.Port); err != nil {
		return err
	}

	file, err := os.Open(s.FileName)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("get source file info: %w", err)
	}
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("source %q is not a regular file", s.FileName)
	}

	fileName := filepath.Base(s.FileName)
	if len(fileName) == 0 || len(fileName) > 255 {
		return fmt.Errorf("source filename must contain 1 to 255 bytes")
	}

	fmt.Printf("Calculating checksum for %s...\n", s.FileName)
	crcVal, err := utils.CalculateCRC(file)
	if err != nil {
		return fmt.Errorf("calculate source checksum: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("rewind source file: %w", err)
	}

	listener, err := net.Listen("tcp", net.JoinHostPort("", s.Port))
	if err != nil {
		return fmt.Errorf("listen on port %s: %w", s.Port, err)
	}
	defer listener.Close()

	return s.send(listener, file, fileInfo, fileName, crcVal)
}

func (s *Sender) send(listener net.Listener, file *os.File, fileInfo os.FileInfo, fileName string, crcVal uint32) error {
	fmt.Printf("Server listening on port %s\n", s.Port)
	fmt.Printf("File: %s\n", s.FileName)
	fmt.Printf("CRC32: %08x\n", crcVal)
	fmt.Println("Waiting for receiver...")

	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("accept receiver connection: %w", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", conn.RemoteAddr())

	header := protocol.FileHeader{
		Protocol: protocol.P2PF_PROTOCOL,
		Version:  protocol.VERSION,
		Size:     uint64(fileInfo.Size()),
		CRC:      crcVal,
		NameLen:  uint8(len(fileName)),
		Name:     fileName,
	}

	headerBytes, err := header.Encode()
	if err != nil {
		return fmt.Errorf("encode file header: %w", err)
	}

	if err := writeAll(conn, headerBytes); err != nil {
		return fmt.Errorf("send file header: %w", err)
	}

	fmt.Println("Starting to Send file Chunks")
	if err := s.SendFile(conn, file, fileInfo.Size()); err != nil {
		return err
	}
	return nil
}

func (s *Sender) SendFile(conn net.Conn, file *os.File, size int64) error {
	buffer := make([]byte, BUFFER_SIZE)

	progress := utils.NewProgressBar(size, "Sending")

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read source file: %w", err)
		}
		if err := writeAll(conn, buffer[:n]); err != nil {
			return fmt.Errorf("send file data: %w", err)
		}
		progress.Add(int64(n))
	}
	progress.Finish()
	fmt.Println("file sent")
	return nil
}

func writeAll(conn net.Conn, data []byte) error {
	for len(data) > 0 {
		if err := conn.SetWriteDeadline(time.Now().Add(transferIdleTimeout)); err != nil {
			return err
		}
		n, err := conn.Write(data)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		data = data[n:]
	}
	return nil
}
