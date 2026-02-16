package poindexter

import (
	"math"
	"testing"
)

func TestLerp(t *testing.T) {
	tests := []struct {
		name    string
		t_, a, b float64
		want    float64
	}{
		{"start", 0, 10, 20, 10},
		{"end", 1, 10, 20, 20},
		{"mid", 0.5, 10, 20, 15},
		{"quarter", 0.25, 0, 100, 25},
		{"extrapolate", 2, 0, 10, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Lerp(tt.t_, tt.a, tt.b)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Lerp(%v, %v, %v) = %v, want %v", tt.t_, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestInverseLerp(t *testing.T) {
	tests := []struct {
		name    string
		v, a, b float64
		want    float64
	}{
		{"start", 10, 10, 20, 0},
		{"end", 20, 10, 20, 1},
		{"mid", 15, 10, 20, 0.5},
		{"equal_range", 5, 5, 5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InverseLerp(tt.v, tt.a, tt.b)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("InverseLerp(%v, %v, %v) = %v, want %v", tt.v, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestRemap(t *testing.T) {
	tests := []struct {
		name                         string
		v, inMin, inMax, outMin, outMax float64
		want                         float64
	}{
		{"identity", 5, 0, 10, 0, 10, 5},
		{"scale_up", 5, 0, 10, 0, 100, 50},
		{"reverse", 3, 0, 10, 10, 0, 7},
		{"offset", 0, 0, 1, 100, 200, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Remap(tt.v, tt.inMin, tt.inMax, tt.outMin, tt.outMax)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Remap = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundToN(t *testing.T) {
	tests := []struct {
		name     string
		f        float64
		decimals int
		want     float64
	}{
		{"zero_dec", 3.456, 0, 3},
		{"one_dec", 3.456, 1, 3.5},
		{"two_dec", 3.456, 2, 3.46},
		{"three_dec", 3.4564, 3, 3.456},
		{"negative", -2.555, 2, -2.56},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoundToN(tt.f, tt.decimals)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("RoundToN(%v, %v) = %v, want %v", tt.f, tt.decimals, got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name       string
		v, min, max float64
		want       float64
	}{
		{"within", 5, 0, 10, 5},
		{"below", -5, 0, 10, 0},
		{"above", 15, 0, 10, 10},
		{"at_min", 0, 0, 10, 0},
		{"at_max", 10, 0, 10, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clamp(tt.v, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("Clamp(%v, %v, %v) = %v, want %v", tt.v, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestClampInt(t *testing.T) {
	if got := ClampInt(5, 0, 10); got != 5 {
		t.Errorf("ClampInt(5, 0, 10) = %v, want 5", got)
	}
	if got := ClampInt(-1, 0, 10); got != 0 {
		t.Errorf("ClampInt(-1, 0, 10) = %v, want 0", got)
	}
	if got := ClampInt(15, 0, 10); got != 10 {
		t.Errorf("ClampInt(15, 0, 10) = %v, want 10", got)
	}
}

func TestMinMaxScale(t *testing.T) {
	tests := []struct {
		name       string
		v, min, max float64
		want       float64
	}{
		{"mid", 5, 0, 10, 0.5},
		{"at_min", 0, 0, 10, 0},
		{"at_max", 10, 0, 10, 1},
		{"equal_range", 5, 5, 5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinMaxScale(tt.v, tt.min, tt.max)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("MinMaxScale(%v, %v, %v) = %v, want %v", tt.v, tt.min, tt.max, got, tt.want)
			}
		})
	}
}
