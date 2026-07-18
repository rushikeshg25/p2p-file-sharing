package receiver

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net"
	"os"
	"p2p-file-sharing/internal/protocol"
	"p2p-file-sharing/internal/utils"
	"path/filepath"
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

var headerTimeout = 30 * time.Second
var transferIdleTimeout = 30 * time.Second

func (r *Receiver) Receive() error {
	if err := utils.ValidatePort(r.Port); err != nil {
		return err
	}
	if r.FileName == "" {
		return fmt.Errorf("output filename cannot be empty")
	}
	if _, err := os.Lstat(r.FileName); err == nil {
		return fmt.Errorf("output file %q already exists", r.FileName)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("check output file: %w", err)
	}

	serverAddress := net.JoinHostPort(r.Address, r.Port)
	fmt.Printf("Connecting to sender at %s...\n", serverAddress)
	conn, err := net.DialTimeout("tcp", serverAddress, 10*time.Second)
	if err != nil {
		return fmt.Errorf("connect to sender at %s: %w", serverAddress, err)
	}
	defer conn.Close()
	fmt.Printf("Connected to %s; waiting for file header...\n", conn.RemoteAddr())

	if err := conn.SetReadDeadline(time.Now().Add(headerTimeout)); err != nil {
		return fmt.Errorf("set file header timeout: %w", err)
	}
	header, err := protocol.Decode(conn)
	if err != nil {
		return fmt.Errorf("read file header from sender: %w", err)
	}
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		return fmt.Errorf("clear file header timeout: %w", err)
	}
	if header.Size > math.MaxInt64 {
		return fmt.Errorf("file is too large: %d bytes", header.Size)
	}
	if header.Name == "" || filepath.Base(header.Name) != header.Name {
		return fmt.Errorf("sender provided invalid filename %q", header.Name)
	}

	fmt.Println("Original filename", header.Name)

	destinationDir := filepath.Dir(r.FileName)
	tempPattern := "." + filepath.Base(r.FileName) + ".part-*"
	file, err := os.CreateTemp(destinationDir, tempPattern)
	if err != nil {
		return fmt.Errorf("create temporary output file: %w", err)
	}
	tempName := file.Name()
	tempPublished := false
	defer func() {
		file.Close()
		if !tempPublished {
			os.Remove(tempName)
		}
	}()

	if err := r.receiveFileData(conn, file, int64(header.Size)); err != nil {
		return err
	}

	fmt.Println("Checking CRC checksum")

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("rewind received file: %w", err)
	}
	receivedFileCRC, err := utils.CalculateCRC(file)
	if err != nil {
		return fmt.Errorf("calculate received file checksum: %w", err)
	}

	if receivedFileCRC != header.CRC {
		return fmt.Errorf("CRC mismatch: expected %08x, got %08x", header.CRC, receivedFileCRC)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync received file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close received file: %w", err)
	}
	if err := os.Link(tempName, r.FileName); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("output file %q already exists", r.FileName)
		}
		return fmt.Errorf("publish received file: %w", err)
	}
	if err := os.Remove(tempName); err != nil {
		return fmt.Errorf("remove temporary output file: %w", err)
	}
	tempPublished = true

	fmt.Println("CRC verified")
	fmt.Println("p2p Transfer done")
	return nil
}

func (r *Receiver) receiveFileData(conn net.Conn, file *os.File, fileSize int64) error {
	buffer := make([]byte, BUFFER_SIZE)
	var totalReceived int64

	progress := utils.NewProgressBar(fileSize, "Receiving")

	for totalReceived < fileSize {
		remainingBytesSize := fileSize - totalReceived
		bufSize := BUFFER_SIZE
		if remainingBytesSize < int64(BUFFER_SIZE) {
			bufSize = int(remainingBytesSize)
		}

		if err := conn.SetReadDeadline(time.Now().Add(transferIdleTimeout)); err != nil {
			return fmt.Errorf("set transfer timeout: %w", err)
		}
		n, err := io.ReadFull(conn, buffer[:bufSize])
		if err != nil {
			return fmt.Errorf("receive file data after %d of %d bytes: %w", totalReceived, fileSize, err)
		}

		if _, err := file.Write(buffer[:n]); err != nil {
			return fmt.Errorf("write temporary output file: %w", err)
		}

		totalReceived += int64(n)
		progress.Add(int64(n))
	}

	progress.Finish()
	return nil
}
