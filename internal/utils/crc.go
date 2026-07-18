package utils

import (
	"bufio"
	"hash/crc32"
	"io"
)

var crcTable = crc32.MakeTable(crc32.IEEE)

func CalculateCRC(input io.Reader) (uint32, error) {
	reader := bufio.NewReader(input)
	crc := uint32(0)
	buffer := make([]byte, 8192)

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			crc = crc32.Update(crc, crcTable, buffer[:n])
		}
		if err == io.EOF {
			return crc, nil
		}
		if err != nil {
			return 0, err
		}
	}
}
