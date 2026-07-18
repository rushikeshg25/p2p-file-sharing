package protocol

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

func TestHeaderRoundTrip(t *testing.T) {
	want := &FileHeader{Size: 1234, CRC: 0x12345678, NameLen: 1, Name: "photo.jpg"}
	encoded, err := want.Encode()
	if err != nil {
		t.Fatal(err)
	}

	got, err := Decode(bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}
	if got.Protocol != P2PF_PROTOCOL || got.Version != VERSION || got.Size != want.Size || got.CRC != want.CRC || got.Name != want.Name {
		t.Fatalf("decoded header = %+v", got)
	}
	if got.NameLen != uint8(len(want.Name)) {
		t.Fatalf("name length = %d, want %d", got.NameLen, len(want.Name))
	}
}

func TestHeaderEncodeRejectsInvalidName(t *testing.T) {
	for _, name := range []string{"", strings.Repeat("a", 256)} {
		if _, err := (&FileHeader{Name: name}).Encode(); err == nil {
			t.Fatalf("Encode accepted filename of length %d", len(name))
		}
	}
}

func TestDecodeRejectsInvalidProtocolAndVersion(t *testing.T) {
	tests := []struct {
		name     string
		protocol uint32
		version  uint32
	}{
		{name: "protocol", protocol: 0, version: VERSION},
		{name: "version", protocol: P2PF_PROTOCOL, version: VERSION + 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := make([]byte, HEADER_SIZE)
			binary.BigEndian.PutUint32(header[0:4], tt.protocol)
			binary.BigEndian.PutUint32(header[4:8], tt.version)
			if _, err := Decode(bytes.NewReader(header)); err == nil {
				t.Fatal("Decode accepted invalid header")
			}
		})
	}
}

func TestDecodeRejectsTruncatedHeaderAndName(t *testing.T) {
	if _, err := Decode(bytes.NewReader(make([]byte, HEADER_SIZE-1))); err == nil {
		t.Fatal("Decode accepted truncated fixed header")
	}

	header := make([]byte, HEADER_SIZE)
	binary.BigEndian.PutUint32(header[0:4], P2PF_PROTOCOL)
	binary.BigEndian.PutUint32(header[4:8], VERSION)
	header[20] = 3
	if _, err := Decode(bytes.NewReader(append(header, 'a'))); err == nil {
		t.Fatal("Decode accepted truncated filename")
	}
}
