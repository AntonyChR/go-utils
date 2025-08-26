package array

import "testing"

func TestFlat(t *testing.T) {

	nested := [][]int{{1, 2}, {3, 4, 5}, {6, 7, 8, 9}}
	flatten := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	res := Flatten(nested)

	if len(res) != len(flatten) {
		t.Errorf("incorrect array length, expected = %d, got = %d", len(flatten), len(res))
	}
}

func TestMap(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}

	res := Map(nums, func(n int) int {
		return n * 2
	})

	if len(res) != len(expected) {
		t.Errorf("incorrect slice length, expected = %d, got = %d", len(expected), len(res))
	}

	for i, el := range res {
		if el != expected[i] {
			t.Errorf("incorrect value, expected = %d, got = %d", expected[i], el)
		}
	}
}

func TestMutMap(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}

	MutMap(nums, func(n int) int {
		return n * 2
	})

	if len(nums) != len(expected) {
		t.Errorf("incorrect slice length, expected = %d, got = %d", len(expected), len(nums))
	}

	for i, el := range nums {
		if el != expected[i] {
			t.Errorf("incorrect value, expected = %d, got = %d", expected[i], el)
		}
	}

}

func TestMax(t *testing.T) {
	nums1 := []int{1, 2, 3, 4, 5}
	expected1 := 5

	res := Max(nums1)

	if res != expected1 {
		t.Errorf("incorrect max value, expected = %d, got = %d", expected1, res)
	}

	nums2 := []int{1, 2}
	expected2 := 2

	res = Max(nums2)

	if res != expected2 {
		t.Errorf("incorrect max value, expected = %d, got = %d", expected1, res)
	}
}

func TestFilter(t *testing.T) {
	input1 := []int{1,2,3,4,5,6,7,8,9,10}
	expect := []int{2,4,6,8,10}

	res := FilterPrealloc(input1,func(n int) bool {
		return n%2==0
	})

	if len(res) != len(expect) {
			t.Errorf("inconsistent slice length, got %v, expect %v", len(res), len(expect))
	}

	for i, n := range res {
		if n != expect[i] {
			t.Errorf("Error filtering elements, got %v, expect %v", res, expect)
		}
	}
}

func BenchmarkFilterPrealloc(b *testing.B) {
	size := 1_000_000
	testData := make([]int, size)
	for i := 0; i < size; i++ {
		testData[i] = i
	}

	// test callback 
	isEven := func(n int) bool {
		return n%2 == 0
	}

	// discard setup time 
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FilterPrealloc(testData, isEven)
	}
}

func BenchmarkFilterInPlace(b *testing.B) {
	size := 1_000_000
	testData := make([]int, size)
	for i := 0; i < size; i++ {
		testData[i] = i
	}

	isEven := func(n int) bool {
		return n%2 == 0
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FilterInPlace(testData, isEven)
	}
}
