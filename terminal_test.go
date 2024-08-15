package utils

import (
	"testing"
)

func TestGetTermDims(t *testing.T) {
	_, _, err := GetTerminalSize()
	if err != nil {
		t.Error(err)
	}
}
