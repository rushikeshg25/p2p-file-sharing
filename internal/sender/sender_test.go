package sender

import (
	"bytes"
	"io"
	"net"
	"os"
	"p2p-file-sharing/internal/receiver"
	"p2p-file-sharing/internal/utils"
	"path/filepath"
	"strconv"
	"testing"
)

type shortWriteConn struct {
	net.Conn
}

func (c shortWriteConn) Write(data []byte) (int, error) {
	if len(data) > 2 {
		data = data[:2]
	}
	return c.Conn.Write(data)
}

func TestWriteAllHandlesShortWrites(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()
	want := []byte("complete payload")
	gotCh := make(chan []byte, 1)
	go func() {
		got := make([]byte, len(want))
		_, _ = io.ReadFull(server, got)
		gotCh <- got
	}()

	if err := writeAll(shortWriteConn{Conn: client}, want); err != nil {
		t.Fatal(err)
	}
	if got := <-gotCh; !bytes.Equal(got, want) {
		t.Fatalf("received %q, want %q", got, want)
	}
}

func TestSenderReceiverIntegration(t *testing.T) {
	data := bytes.Repeat([]byte{0, 1, 2, 3, 255}, 1000)
	dir := t.TempDir()
	source := filepath.Join(dir, "source.bin")
	if err := os.WriteFile(source, data, 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(source)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	crc, err := utils.CalculateCRC(file)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	s := NewSender(port, source)
	senderErr := make(chan error, 1)
	go func() {
		defer listener.Close()
		senderErr <- s.send(listener, file, info, filepath.Base(source), crc)
	}()

	destination := filepath.Join(dir, "received.bin")
	if err := receiver.NewReceiver("127.0.0.1", port, destination).Receive(); err != nil {
		t.Fatal(err)
	}
	if err := <-senderErr; err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("received file does not match source")
	}
}
