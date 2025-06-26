package utils

import (
	"testing"
)

func TestProd3d(t *testing.T) {
	a := Vec{2, 3, 4}
	b := Vec{5, 6, 7}

	result := a.Prod3d(b)
	expected := Vec{-3, 6, -3}

    if !expected.Equals(result){
		t.Errorf("Prod3d was incorrect, got: %v, want: %v.", result, expected)
	}
}
