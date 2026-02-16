package poindexter

import (
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		data []float64
		want float64
	}{
		{"empty", nil, 0},
		{"single", []float64{5}, 5},
		{"multiple", []float64{1, 2, 3, 4, 5}, 15},
		{"negative", []float64{-1, -2, 3}, 0},
		{"floats", []float64{0.1, 0.2, 0.3}, 0.6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.data)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Sum(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		name string
		data []float64
		want float64
	}{
		{"empty", nil, 0},
		{"single", []float64{5}, 5},
		{"multiple", []float64{2, 4, 6}, 4},
		{"floats", []float64{1.5, 2.5}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mean(tt.data)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Mean(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestVariance(t *testing.T) {
	tests := []struct {
		name string
		data []float64
		want float64
	}{
		{"empty", nil, 0},
		{"constant", []float64{5, 5, 5}, 0},
		{"simple", []float64{2, 4, 4, 4, 5, 5, 7, 9}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Variance(tt.data)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Variance(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestStdDev(t *testing.T) {
	got := StdDev([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	if math.Abs(got-2) > 1e-9 {
		t.Errorf("StdDev = %v, want 2", got)
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		name    string
		data    []float64
		wantMin float64
		wantMax float64
	}{
		{"empty", nil, 0, 0},
		{"single", []float64{3}, 3, 3},
		{"ordered", []float64{1, 2, 3}, 1, 3},
		{"reversed", []float64{3, 2, 1}, 1, 3},
		{"negative", []float64{-5, 0, 5}, -5, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := MinMax(tt.data)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("MinMax(%v) = (%v, %v), want (%v, %v)", tt.data, gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestIsUnderrepresented(t *testing.T) {
	tests := []struct {
		name      string
		val       float64
		avg       float64
		threshold float64
		want      bool
	}{
		{"below", 3, 10, 0.5, true},
		{"at", 5, 10, 0.5, false},
		{"above", 7, 10, 0.5, false},
		{"zero_avg", 0, 0, 0.5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUnderrepresented(tt.val, tt.avg, tt.threshold)
			if got != tt.want {
				t.Errorf("IsUnderrepresented(%v, %v, %v) = %v, want %v", tt.val, tt.avg, tt.threshold, got, tt.want)
			}
		})
	}
}
