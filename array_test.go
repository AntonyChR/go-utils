package utils

import "testing"

func TestFlat(t *testing.T){

    nested := [][]int{{1,2},{3,4,5},{6,7,8,9}}
    flatten := []int{1,2,3,4,5,6,7,8,9}

    res := Flatten(nested)

    if len(res) != len(flatten) {
        t.Errorf("incorrect array length, expected = %d, got = %d",  len(flatten), len(res))
    }
}

func TestMap(t *testing.T){
    nums := []int{1,2,3,4,5}
    expected := []int{2,4,6,8,10}

    res := Map(nums, func(n int) int {
        return n * 2
    })

    if len(res) != len(expected){
        t.Errorf("incorrect slice length, expected = %d, got = %d", len(expected), len(res))
    }

    for i, el := range res {
        if el != expected[i]{
            t.Errorf("incorrect value, expected = %d, got = %d", expected[i], el)
        }
    }
}


func TestMutMap(t *testing.T){
    nums := []int{1,2,3,4,5}
    expected := []int{2,4,6,8,10}

    MutMap(nums, func(n int) int {
        return n * 2
    })

    if len(nums) != len(expected){
        t.Errorf("incorrect slice length, expected = %d, got = %d", len(expected), len(nums))
    }

    for i, el := range nums {
        if el != expected[i]{
            t.Errorf("incorrect value, expected = %d, got = %d", expected[i], el)
        }
    }

}
