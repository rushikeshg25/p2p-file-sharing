package utils

import (
	"fmt"
	"strconv"
)

func ValidatePort(port string) error {
	value, err := strconv.Atoi(port)
	if err != nil || value < 1 || value > 65535 {
		return fmt.Errorf("invalid port %q: must be a number from 1 to 65535", port)
	}
	return nil
}
