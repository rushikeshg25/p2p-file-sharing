package receiver

import (
	"hash/crc32"
	"net"
	"os"
	"p2p-file-sharing/internal/protocol"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func startTestSender(t *testing.T, handler func(net.Conn) error) (string, <-chan error) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	go func() {
		defer listener.Close()
		conn, err := listener.Accept()
		if err == nil {
			defer conn.Close()
			err = handler(conn)
		}
		errCh <- err
	}()

	return strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), errCh
}

func sendHeaderAndData(conn net.Conn, name string, size uint64, crc uint32, data []byte) error {
	header, err := (&protocol.FileHeader{Size: size, CRC: crc, Name: name}).Encode()
	if err != nil {
		return err
	}
	if _, err := conn.Write(header); err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

func TestReceivePublishesVerifiedFile(t *testing.T) {
	for _, data := range [][]byte{{}, {0, 1, 2, 3, 255}} {
		t.Run(strconv.Itoa(len(data))+"_bytes", func(t *testing.T) {
			port, senderErr := startTestSender(t, func(conn net.Conn) error {
				return sendHeaderAndData(conn, "source.bin", uint64(len(data)), crc32.ChecksumIEEE(data), data)
			})
			destination := filepath.Join(t.TempDir(), "received.bin")

			if err := NewReceiver("127.0.0.1", port, destination).Receive(); err != nil {
				t.Fatal(err)
			}
			if err := <-senderErr; err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(destination)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(data) {
				t.Fatalf("received data = %v, want %v", got, data)
			}
		})
	}
}

func TestReceiveRefusesExistingDestination(t *testing.T) {
	destination := filepath.Join(t.TempDir(), "existing.bin")
	if err := os.WriteFile(destination, []byte("keep me"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := NewReceiver("127.0.0.1", "3001", destination).Receive()
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Receive error = %v", err)
	}
	got, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "keep me" {
		t.Fatalf("existing destination changed to %q", got)
	}
}

func TestReceiveCleansUpAfterCRCFailure(t *testing.T) {
	data := []byte("corrupted")
	port, senderErr := startTestSender(t, func(conn net.Conn) error {
		return sendHeaderAndData(conn, "source.bin", uint64(len(data)), 0, data)
	})
	dir := t.TempDir()
	destination := filepath.Join(dir, "received.bin")

	err := NewReceiver("127.0.0.1", port, destination).Receive()
	if err == nil || !strings.Contains(err.Error(), "CRC mismatch") {
		t.Fatalf("Receive error = %v", err)
	}
	if err := <-senderErr; err != nil {
		t.Fatal(err)
	}
	assertNoOutputFiles(t, dir, destination)
}

func TestReceiveCleansUpAfterInterruptedTransfer(t *testing.T) {
	port, senderErr := startTestSender(t, func(conn net.Conn) error {
		return sendHeaderAndData(conn, "source.bin", 10, 0, []byte("short"))
	})
	dir := t.TempDir()
	destination := filepath.Join(dir, "received.bin")

	if err := NewReceiver("127.0.0.1", port, destination).Receive(); err == nil {
		t.Fatal("Receive succeeded after interrupted transfer")
	}
	if err := <-senderErr; err != nil {
		t.Fatal(err)
	}
	assertNoOutputFiles(t, dir, destination)
}

func TestReceiveHeaderTimeout(t *testing.T) {
	oldTimeout := headerTimeout
	headerTimeout = 20 * time.Millisecond
	defer func() { headerTimeout = oldTimeout }()

	port, senderErr := startTestSender(t, func(net.Conn) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})
	destination := filepath.Join(t.TempDir(), "received.bin")

	err := NewReceiver("127.0.0.1", port, destination).Receive()
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("Receive error = %v", err)
	}
	if err := <-senderErr; err != nil {
		t.Fatal(err)
	}
}

func assertNoOutputFiles(t *testing.T, dir string, destination string) {
	t.Helper()
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatalf("destination exists after failure: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(dir, ".received.bin.part-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary files remain: %v", matches)
	}
}
