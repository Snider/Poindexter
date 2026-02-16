package poindexter

import (
	"math"
	"math/rand"
)

// RampUp returns a linear ramp from 0 to 1 over the given duration.
// The result is clamped to [0, 1].
func RampUp(elapsed, duration float64) float64 {
	if duration <= 0 {
		return 1
	}
	return Clamp(elapsed/duration, 0, 1)
}

// SineWave returns a sine value with the given period and amplitude.
// Output range is [-amplitude, amplitude].
func SineWave(t, period, amplitude float64) float64 {
	if period == 0 {
		return 0
	}
	return math.Sin(t/period*2*math.Pi) * amplitude
}

// Oscillate modulates a base value with a sine wave.
// Returns base * (1 + sin(t/period*2π) * amplitude).
func Oscillate(base, t, period, amplitude float64) float64 {
	if period == 0 {
		return base
	}
	return base * (1 + math.Sin(t/period*2*math.Pi)*amplitude)
}

// Noise generates seeded pseudo-random values.
type Noise struct {
	rng *rand.Rand
}

// NewNoise creates a seeded noise generator.
func NewNoise(seed int64) *Noise {
	return &Noise{rng: rand.New(rand.NewSource(seed))}
}

// Float64 returns a random value in [-variance, variance].
func (n *Noise) Float64(variance float64) float64 {
	return (n.rng.Float64() - 0.5) * 2 * variance
}

// Int returns a random integer in [0, max).
// Returns 0 if max <= 0.
func (n *Noise) Int(max int) int {
	if max <= 0 {
		return 0
	}
	return n.rng.Intn(max)
}
