package utils

import (
	"testing"
)

func TestGetTermDims(t *testing.T) {
	_, _, err := GetTermDims()
	if err != nil {
		t.Error()
	}
}
