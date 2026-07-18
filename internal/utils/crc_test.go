package utils

import (
	"errors"
	"hash/crc32"
	"strings"
	"testing"
)

type failingReader struct {
	err error
}

func (r failingReader) Read([]byte) (int, error) {
	return 0, r.err
}

func TestCalculateCRC(t *testing.T) {
	for _, content := range []string{"", "hello, world"} {
		got, err := CalculateCRC(strings.NewReader(content))
		if err != nil {
			t.Fatal(err)
		}
		want := crc32.ChecksumIEEE([]byte(content))
		if got != want {
			t.Fatalf("CRC(%q) = %08x, want %08x", content, got, want)
		}
	}
}

func TestCalculateCRCPropagatesReadError(t *testing.T) {
	want := errors.New("read failed")
	if _, err := CalculateCRC(failingReader{err: want}); !errors.Is(err, want) {
		t.Fatalf("CalculateCRC error = %v, want %v", err, want)
	}
}
