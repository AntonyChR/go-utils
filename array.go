package utils

import (
	"cmp"
)

// Map applies a transformation function to each element of a slice,
// returning a new slice with the results. The original slice is not modified.
//
// Parameters:
//   - arr: The input slice containing elements of any type.
//   - cb: A function that takes an element of type T and returns an
//         element of the same type T.
//
// Returns:
//   - A new slice of type T with the transformed elements.
func Map[T any](arr []T, cb func(T)T)[]T{
    c :=  make([]T, len(arr))
    for i,e := range arr {
        c[i] = cb(e)
    }
    return c
}

// MutMap applies a transformation function to each element of a slice,
// modifying the elements of the original slice in place rather than creating a new one.
//
// Parameters:
//   - arr: The input slice containing elements of any type. This slice
//          will be modified directly.
//   - cb: A function that takes an element of type T and returns an
//         element of the same type T.
//
// Returns:
//   - The original slice with the transformed elements.
func MutMap[T any](arr []T, cb func(T) T) []T {
    for i := range arr {
        arr[i] = cb(arr[i])
    }
    return arr
}

func Flatten[T any](arr [][]T) []T{
    var narr []T
    for _,el := range arr {
        narr = append(narr, el...)
    }
    return narr
}

// Contains checks if a given element exists in a slice.
func Contains[T comparable](arr []T, el T) bool {
    for _, e := range arr {
        if e == el {
            return true
        }
    }
    return false
}


func Max[T cmp.Ordered](arr []T) T {
	maxVal := arr[0]
	for _, num := range arr[1:] {
		if num > maxVal {
			maxVal = num
		}
	}
	return maxVal
}

// FilterPrealloc creates a new slice containing only the elements that satisfy the given condition.
//
// Parameters:
//   - arr: The input slice containing elements of any type.
//   - cb: A function that takes an element of type T and returns a boolean value.
//
// Returns:
//   - A new slice of type T containing only the elements that satisfy the condition.
func FilterPrealloc[T any](arr []T, cb func(T)bool) []T {
    l := make([]T, 0, len(arr))
    for _,b := range arr {
        if cb(b) {
            l = append(l, b)
        }
    }
    return l
}

// FilterInPlace modifies the original slice in place, removing elements that do not satisfy the given condition.
//
// Parameters:
//   - arr: The input slice containing elements of any type. This slice
//          will be modified directly.
//   - cb: A function that takes an element of type T and returns a boolean value.
//
// Returns:
//   - The original slice with only the elements that satisfy the condition.
func FilterInPlace[T any](arr []T, cb func(T)bool) []T {
    writeIndex := 0
    for _, b := range arr {
        if cb(b) {
            arr[writeIndex] = b
            writeIndex++
        }
    }
    return arr[:writeIndex]
}




func DeepCopyArrMap[K comparable, V any](original []map[K]V) []map[K]V {
	if original == nil {
		return nil
	}
    deepCopy := make([]map[K]V, len(original))

    for i, originalMap := range original {
        newMap := make(map[K]V, len(originalMap))

        for key, value := range originalMap {
            newMap[key] = value
        }

        deepCopy[i] = newMap
    }

    return deepCopy
}
