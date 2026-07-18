package utils

import "testing"

func TestValidatePort(t *testing.T) {
	for _, port := range []string{"1", "3001", "65535"} {
		if err := ValidatePort(port); err != nil {
			t.Fatalf("ValidatePort(%q): %v", port, err)
		}
	}
	for _, port := range []string{"", "0", "65536", "http", "-1"} {
		if err := ValidatePort(port); err == nil {
			t.Fatalf("ValidatePort(%q) succeeded", port)
		}
	}
}
