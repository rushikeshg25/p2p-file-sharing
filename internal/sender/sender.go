package sender

import (
	"log"
	"os"
)

type Sender struct {
	Port     int
	FileName string
}

const BUFFER_SIZE = 8192 //File chunk size 8KB

func NewSender(port int, FileName string) *Sender {
	return &Sender{
		Port:     port,
		FileName: FileName,
	}
}

func (s *Sender) SendFile() {
	file, err := os.Open(s.FileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

}
