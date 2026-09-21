package sender

import (
	"context"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

func TestCancelledAccept(t *testing.T) {
	f, e := os.CreateTemp(t.TempDir(), "source")
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	info, _ := f.Stat()
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer l.Close()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- (&Sender{}).sendContext(ctx, l, f, info, "empty", 0) }()
	cancel()
	select {
	case e := <-done:
		if e == nil {
			t.Fatal("cancellation succeeded")
		}
	case <-time.After(time.Second):
		t.Fatal("accept not cancelled")
	}
}
func TestSourceLengthBound(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "source")
	defer f.Close()
	f.WriteString("abcdef")
	f.Seek(0, 0)
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	done := make(chan error, 1)
	go func() { done <- (&Sender{}).SendFile(a, f, 3); a.Close() }()
	got, e := io.ReadAll(b)
	if e != nil || string(got) != "abc" {
		t.Fatal(string(got), e)
	}
	if e := <-done; e != nil {
		t.Fatal(e)
	}
}
func TestSourceTruncation(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "source")
	defer f.Close()
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	if (&Sender{}).SendFile(a, f, 1) == nil {
		t.Fatal("truncated source accepted")
	}
}
