package utils

import (
	"bufio"
	"hash/crc32"
	"os"
)

var crcTable = crc32.MakeTable(crc32.IEEE)

func CalculateCRC(fd *os.File) (uint32, error) {
	reader := bufio.NewReader(fd)
	crc := uint32(0)
	buffer := make([]byte, 8192)

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			crc = crc32.Update(crc, crcTable, buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return crc, nil
}

func CalculateCRC32Stream(data []byte) uint32 {
	return crc32.Checksum(data, crcTable)
}
