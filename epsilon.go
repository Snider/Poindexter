package poindexter

import "math"

// ApproxEqual returns true if the absolute difference between a and b
// is less than epsilon.
func ApproxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// ApproxZero returns true if the absolute value of v is less than epsilon.
func ApproxZero(v, epsilon float64) bool {
	return math.Abs(v) < epsilon
}
