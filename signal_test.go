package poindexter

import (
	"math"
	"testing"
)

func TestRampUp(t *testing.T) {
	tests := []struct {
		name             string
		elapsed, duration float64
		want             float64
	}{
		{"start", 0, 30, 0},
		{"mid", 15, 30, 0.5},
		{"end", 30, 30, 1},
		{"over", 60, 30, 1},
		{"negative", -5, 30, 0},
		{"zero_duration", 10, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RampUp(tt.elapsed, tt.duration)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("RampUp(%v, %v) = %v, want %v", tt.elapsed, tt.duration, got, tt.want)
			}
		})
	}
}

func TestSineWave(t *testing.T) {
	// At t=0, sin(0) = 0
	if got := SineWave(0, 10, 5); math.Abs(got) > 1e-9 {
		t.Errorf("SineWave(0, 10, 5) = %v, want 0", got)
	}
	// At t=period/4, sin(π/2) = 1, so result = amplitude
	got := SineWave(2.5, 10, 5)
	if math.Abs(got-5) > 1e-9 {
		t.Errorf("SineWave(2.5, 10, 5) = %v, want 5", got)
	}
	// Zero period returns 0
	if got := SineWave(5, 0, 5); got != 0 {
		t.Errorf("SineWave(5, 0, 5) = %v, want 0", got)
	}
}

func TestOscillate(t *testing.T) {
	// At t=0, sin(0)=0, so result = base * (1 + 0) = base
	got := Oscillate(100, 0, 10, 0.05)
	if math.Abs(got-100) > 1e-9 {
		t.Errorf("Oscillate(100, 0, 10, 0.05) = %v, want 100", got)
	}
	// At t=period/4, sin=1, so result = base * (1 + amplitude)
	got = Oscillate(100, 2.5, 10, 0.05)
	if math.Abs(got-105) > 1e-9 {
		t.Errorf("Oscillate(100, 2.5, 10, 0.05) = %v, want 105", got)
	}
	// Zero period returns base
	if got := Oscillate(100, 5, 0, 0.05); got != 100 {
		t.Errorf("Oscillate(100, 5, 0, 0.05) = %v, want 100", got)
	}
}

func TestNoise(t *testing.T) {
	n := NewNoise(42)

	// Float64 should be within [-variance, variance]
	for i := 0; i < 1000; i++ {
		v := n.Float64(0.1)
		if v < -0.1 || v > 0.1 {
			t.Fatalf("Float64(0.1) = %v, outside [-0.1, 0.1]", v)
		}
	}

	// Int should be within [0, max)
	n2 := NewNoise(42)
	for i := 0; i < 1000; i++ {
		v := n2.Int(10)
		if v < 0 || v >= 10 {
			t.Fatalf("Int(10) = %v, outside [0, 10)", v)
		}
	}

	// Int with zero max returns 0
	if got := n.Int(0); got != 0 {
		t.Errorf("Int(0) = %v, want 0", got)
	}
	if got := n.Int(-1); got != 0 {
		t.Errorf("Int(-1) = %v, want 0", got)
	}
}

func TestNoiseDeterministic(t *testing.T) {
	n1 := NewNoise(123)
	n2 := NewNoise(123)
	for i := 0; i < 100; i++ {
		a := n1.Float64(1.0)
		b := n2.Float64(1.0)
		if a != b {
			t.Fatalf("iteration %d: different values for same seed: %v != %v", i, a, b)
		}
	}
}
