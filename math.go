package  utils

import "math"

func Round(value float64, decimals uint8) float64{
    return math.Round(value * 10.0 * float64(decimals)) / (10.0 * float64(decimals))
}
