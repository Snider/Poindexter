package poindexter

import "math"

// Sum returns the sum of all values. Returns 0 for empty slices.
func Sum(data []float64) float64 {
	var s float64
	for _, v := range data {
		s += v
	}
	return s
}

// Mean returns the arithmetic mean. Returns 0 for empty slices.
func Mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	return Sum(data) / float64(len(data))
}

// Variance returns the population variance. Returns 0 for empty slices.
func Variance(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	m := Mean(data)
	var ss float64
	for _, v := range data {
		d := v - m
		ss += d * d
	}
	return ss / float64(len(data))
}

// StdDev returns the population standard deviation.
func StdDev(data []float64) float64 {
	return math.Sqrt(Variance(data))
}

// MinMax returns the minimum and maximum values.
// Returns (0, 0) for empty slices.
func MinMax(data []float64) (min, max float64) {
	if len(data) == 0 {
		return 0, 0
	}
	min, max = data[0], data[0]
	for _, v := range data[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

// IsUnderrepresented returns true if val is below threshold fraction of avg.
// For example, IsUnderrepresented(3, 10, 0.5) returns true because 3 < 10*0.5.
func IsUnderrepresented(val, avg, threshold float64) bool {
	return val < avg*threshold
}
