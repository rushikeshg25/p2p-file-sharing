package main

import (
	"strings"
	"testing"
)

func TestReceiveCLIFormsRemainSupported(t *testing.T) {
	tests := [][]string{
		{"receive", "output.bin", "invalid"},
		{"receive", "192.168.1.10", "output.bin", "invalid"},
	}
	for _, args := range tests {
		err := run(args)
		if err == nil || !strings.Contains(err.Error(), "invalid port") {
			t.Fatalf("run(%q) error = %v", args, err)
		}
	}
}

func TestSendRejectsInvalidPortBeforeOpeningFile(t *testing.T) {
	err := run([]string{"send", "missing.bin", "invalid"})
	if err == nil || !strings.Contains(err.Error(), "invalid port") {
		t.Fatalf("run error = %v", err)
	}
}
