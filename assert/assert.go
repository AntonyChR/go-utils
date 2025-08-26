// assert package provides helper functions for writing tests, such as Assert, AssertEq, and AssertNe.
package assert

import (
	"testing"
)

func assert(value bool, args []any){
    if !value {
        msg := ""
        var t *testing.T

        if len(args) > 0 {
            if m, ok := args[0].(string); ok{
                msg = m
            } else {
                panic("Assert: first argument must be a string")
            }
        } else {
            msg = "Assertion failed"
        }

        if len(args) > 1 {
            if test, ok := args[1].(*testing.T); ok {
                t = test
            } else {
                panic("Assert: second argument must be of type *testing.T")
            }
        }

        if t != nil {
            t.Error(msg)
        } else {
            panic(msg)
        }
    }
} 

// Assert checks if a boolean value is true. If the value is false,
// it generates an error message. Depending on the provided arguments,
// it can either panic or log an error in a testing.T instance.
//
// Parameters:
//  - value: The boolean value to check.
//  - args: Optional arguments. Can include:
//      - msg (string): A custom error message.
//      - t (*testing.T): A testing.T instance to log errors.
//
// Usage:
//  Assert(true)                          // Does nothing since the value is true.
//  Assert(false, "error message")        // Panics with the error message.
//  Assert(false, "error message", t)     // Logs an error in t and continues execution.
//
// Examples:
//  func TestAssert(t *testing.T) {
//      Assert(false, "the tested value is false", t)
//  }
//
// Notes:
//  - If a custom message is provided but no testing.T instance is given,
//    the function will panic with the message.
//  - If no custom message is provided, "Assertion failed" will be used as the default message.
//  - The function will panic if the argument types are not as expected.
func Assert(value bool, args ...any) {
    assert(value, args)
}

// AssertEq checks if two comparable values are equal. If the values are not equal,
// it generates an error message. Depending on the provided arguments, 
// it can either panic or log an error in a testing.T instance.
//
// Parameters:
//  - left_value: The first value to compare.
//  - right_value: The second value to compare.
//  - args: Optional arguments. Can include:
//      - msg (string): A custom error message.
//      - t (*testing.T): A testing.T instance to log errors.
//
// Usage:
//  AssertEq(1, 1)                          // Does nothing since the values are equal.
//  AssertEq(1, 2, "values are not equal")  // Panics with the error message.
//  AssertEq(1, 2, "values are not equal", t) // Logs an error in t and continues execution.
//
// Examples:
//  func TestAssertEq(t *testing.T) {
//      AssertEq(1, 2, "1 is not equal to 2", t)
//  }
//
// Notes:
//  - If a custom message is provided but no testing.T instance is given,
//    the function will panic with the message.
//  - If no custom message is provided, "Assertion failed" will be used as the default message.
//  - The function will panic if the argument types are not as expected.
func AssertEq[K comparable](left_value, right_value K, args ...any) {
    assert(left_value == right_value, args)
}

// AssertNe checks if two comparable values are not equal. If the values are equal,
// it generates an error message. Depending on the provided arguments, 
// it can either panic or log an error in a testing.T instance.
//
// Parameters:
//  - left_value: The first value to compare.
//  - right_value: The second value to compare.
//  - args: Optional arguments. Can include:
//      - msg (string): A custom error message.
//      - t (*testing.T): A testing.T instance to log errors.
//
// Usage:
//  AssertNe(1, 2)                          // Does nothing since the values are not equal.
//  AssertNe(1, 1, "values are equal")      // Panics with the error message.
//  AssertNe(1, 1, "values are equal", t)   // Logs an error in t and continues execution.
//
// Examples:
//  func TestAssertNe(t *testing.T) {
//      AssertNe(1, 1, "1 is equal to 1", t)
//  }
//
// Notes:
//  - If a custom message is provided but no testing.T instance is given,
//    the function will panic with the message.
//  - If no custom message is provided, "Assertion failed" will be used as the default message.
//  - The function will panic if the argument types are not as expected.
func AssertNe[K comparable](left_value, right_value K, args ...any) {
    assert(left_value != right_value, args)
}

// Must is a generic function that simplifies error handling by enforcing an automatic panic
// when an error is encountered. It accepts two parameters: a value `x` of any type `T` 
// and an `error` value `err`.
//
// If `err` is not `nil`, Must triggers a panic with the provided error.
// Otherwise, it simply returns `x`.
//
// This function is useful in cases where errors are unexpected or where explicit error handling
// would add unnecessary complexity to the code.
//
// Usage example:
//    file := Must(os.Open("file.txt"))
//
// In this example, if the `os.Open` function encounters an error opening the file,
// `Must` will panic with that error. Otherwise, it returns the opened file.
//
// Parameters:
// - x: a value of the generic type `T` to return if `err` is `nil`.
// - err: an `error` value that represents a potential error.
//
// Returns:
// - The value `x` of type `T`, if `err` is `nil`.
//
// Panic:
// - If `err` is not `nil`, Must panics with the provided error.
func Must[T any] (x T, err error) T {
    if err != nil {
        panic(err)
    }
    return x
}

// CheckErr is a helper function that simplifies error handling by panicking 
// if an error is encountered. It takes a single parameter, `err` of type `error`.
//
// If `err` is not `nil`, CheckErr triggers a panic with the provided error.
// This function is useful in situations where errors are considered fatal
// and handling them explicitly would add unnecessary complexity.
//
// Usage example:
//    file, err := os.Open("file.txt")
//    CheckErr(err)
//
// In this example, if there is an error while opening the file, `CheckErr` 
// will cause the program to panic with the given error. If no error is present, 
// the program continues as normal.
//
// Parameters:
// - err: an `error` value that represents a potential error.
//
// Panic:
// - If `err` is not `nil`, CheckErr panics with the provided error.
func CheckErr(err error){
    if err != nil {
        panic(err)
    }
}
