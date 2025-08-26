package utils

import "math"

// Round rounds a float64 value to a specified number of decimal places.
//
// Parameters:
//   - value: The float64 value to round.
//   - decimals: The number of decimal places to round to.
//
// Returns:
//   - The rounded float64 value.
func Round(value float64, decimals uint8) float64{
    return math.Round(value * 10.0 * float64(decimals)) / (10.0 * float64(decimals))
}
