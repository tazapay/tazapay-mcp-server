package math

import "math"

// Round2Decimal round off value to 2 decimal places
func Round2Decimal(x float64) float64 {
	return math.Round(x*Int100) / Int100
}

// Round3Decimal round off value to 3 decimal places
func Round3Decimal(x float64) float64 {
	return math.Round(x*Int1000) / Int1000
}

// Round6Decimal round off value to 6 decimal places
func Round6Decimal(x float64) float64 {
	return math.Round(x*Int1000000) / Int1000000
}
