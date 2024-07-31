package utils

import (
	"strconv"
	"testing"
)

func TestGetTermDims(t *testing.T) {
	height, width, err := GetTermDims()
	if err != nil {
		t.Error(err)
	}
	_, err = strconv.Atoi(height)
	if err != nil {
		t.Errorf("height value is not a number, got = \"%s\"", height)
	}
	_, err = strconv.Atoi(width)
	if err != nil {
		t.Errorf("width value is not a number, got = \"%s\"", width)
	}

}
