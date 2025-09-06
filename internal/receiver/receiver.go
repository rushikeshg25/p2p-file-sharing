package receiver

type Receiver struct {
	Port     int
	FilePath string
}

func NewSender(port int, filePath string) *Receiver {
	return &Receiver{
		Port:     port,
		FilePath: filePath,
	}
}
