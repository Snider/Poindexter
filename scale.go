package poindexter

import "math"

// Lerp performs linear interpolation between a and b.
// t=0 returns a, t=1 returns b, t=0.5 returns the midpoint.
func Lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

// InverseLerp returns where v falls between a and b as a fraction [0,1].
// Returns 0 if a == b.
func InverseLerp(v, a, b float64) float64 {
	if a == b {
		return 0
	}
	return (v - a) / (b - a)
}

// Remap maps v from the range [inMin, inMax] to [outMin, outMax].
// Equivalent to Lerp(InverseLerp(v, inMin, inMax), outMin, outMax).
func Remap(v, inMin, inMax, outMin, outMax float64) float64 {
	return Lerp(InverseLerp(v, inMin, inMax), outMin, outMax)
}

// RoundToN rounds f to n decimal places.
func RoundToN(f float64, decimals int) float64 {
	mul := math.Pow(10, float64(decimals))
	return math.Round(f*mul) / mul
}

// Clamp restricts v to the range [min, max].
func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ClampInt restricts v to the range [min, max].
func ClampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// MinMaxScale normalizes v into [0,1] given its range [min, max].
// Returns 0 if min == max.
func MinMaxScale(v, min, max float64) float64 {
	if min == max {
		return 0
	}
	return (v - min) / (max - min)
}
