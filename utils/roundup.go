package utils

import "math"

func RoundUp(val float64, decimals int) float64 {
	factor := math.Pow(10, float64(decimals))
	return math.Ceil(val*factor) / factor
}
