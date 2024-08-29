package utils

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
