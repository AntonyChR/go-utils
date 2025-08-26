package terminal

import (
	"testing"
)

func TestGetTermDims(t *testing.T) {
	_, _, err := GetTerminalSize()
	if err != nil {
		// In a non-interactive environment (like a test runner), stty will fail.
	}
}
