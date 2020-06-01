package util

import "math"

// compareFloatPrecision is used to ignore a specified amount of imprecision
// when comparing two floating point numbers.
const compareFloatPrecision = 1e-9

// CompareFloat compares two floating point numbers, ignoring small differences
// caused by rounding errors due to imprecision; returns -1 if float a is less
// than float b, 0 if float a is equal to float b, and 1 if float a is greater
// than float b.
func CompareFloat(a, b float64) int {

	if math.Abs(a-b) < compareFloatPrecision {
		return 0
	}

	if a < b {
		return -1
	}

	return 1

}
