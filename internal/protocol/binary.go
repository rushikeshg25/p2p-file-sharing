package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	P2PF_PROTOCOL = 0x50325046 // "P2PF" in hex
	VERSION       = 1
	HEADER_SIZE   = 4 + 4 + 8 + 4 + 1 // Protocol + version + size + crc + name_len
)

type Header struct {
	Protocol uint32
	Version  uint32
	Size     uint64
	CRC      uint32
	NameLen  uint8
	Name     string
}

func (h *Header) Encode() ([]byte, error) {
	if len(h.Name) > 255 {
		return nil, errors.New("name too long")
	}
	headerSize := HEADER_SIZE + len(h.Name)

	headerBuf := make([]byte, headerSize)

	binary.BigEndian.PutUint32(headerBuf[0:4], P2PF_PROTOCOL)
	binary.BigEndian.PutUint32(headerBuf[4:8], VERSION)
	binary.BigEndian.PutUint64(headerBuf[8:16], h.Size)
	binary.BigEndian.PutUint32(headerBuf[16:20], h.CRC)
	headerBuf[20] = h.NameLen
	copy(headerBuf[21:], h.Name)

	return headerBuf, nil
}

func Decodeheader(r io.Reader) (*Header, error) {
	header := make([]byte, HEADER_SIZE)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	h := &Header{
		Protocol: binary.BigEndian.Uint32(header[0:4]),
		Version:  binary.BigEndian.Uint32(header[4:8]),
		Size:     binary.BigEndian.Uint64(header[8:16]),
		CRC:      binary.BigEndian.Uint32(header[16:20]),
		NameLen:  header[20],
	}

	if h.Protocol != P2PF_PROTOCOL {
		return nil, errors.New("invalid protocol")
	}

	if h.Version != VERSION {
		return nil, errors.New("invalid version")
	}

	if h.NameLen > 0 {
		nameBytes := make([]byte, h.NameLen)
		if _, err := io.ReadFull(r, nameBytes); err != nil {
			return nil, err
		}
		h.Name = string(nameBytes)
	}

	return h, nil
}
