package utils

import "testing"

func TestFlat(t *testing.T){

    nested := [][]int{{1,2},{3,4,5},{6,7,8,9}}
    flatten := []int{1,2,3,4,5,6,7,8,9}

    res := Flatten(nested)

    if len(res) != len(flatten) {
        t.Errorf("error length, expected = %d, got = %d",  len(flatten), len(res))
    }
}

