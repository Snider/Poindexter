package poindexter

import (
	"math"
	"testing"
)

func TestWeightedScore(t *testing.T) {
	tests := []struct {
		name    string
		factors []Factor
		want    float64
	}{
		{"empty", nil, 0},
		{"single", []Factor{{Value: 5, Weight: 2}}, 10},
		{"multiple", []Factor{
			{Value: 3, Weight: 2},  // 6
			{Value: 1, Weight: -5}, // -5
		}, 1},
		{"lek_heuristic", []Factor{
			{Value: 2, Weight: 2},   // engagement × 2
			{Value: 1, Weight: 3},   // creative × 3
			{Value: 1, Weight: 1.5}, // first person × 1.5
			{Value: 3, Weight: -5},  // compliance × -5
		}, -6.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeightedScore(tt.factors)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("WeightedScore = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRatio(t *testing.T) {
	tests := []struct {
		name        string
		part, whole float64
		want        float64
	}{
		{"half", 5, 10, 0.5},
		{"full", 10, 10, 1},
		{"zero_whole", 5, 0, 0},
		{"zero_part", 0, 10, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ratio(tt.part, tt.whole)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Ratio(%v, %v) = %v, want %v", tt.part, tt.whole, got, tt.want)
			}
		})
	}
}

func TestDelta(t *testing.T) {
	if got := Delta(10, 15); got != 5 {
		t.Errorf("Delta(10, 15) = %v, want 5", got)
	}
	if got := Delta(15, 10); got != -5 {
		t.Errorf("Delta(15, 10) = %v, want -5", got)
	}
}

func TestDeltaPercent(t *testing.T) {
	tests := []struct {
		name      string
		old, new_ float64
		want      float64
	}{
		{"increase", 100, 150, 50},
		{"decrease", 100, 75, -25},
		{"zero_old", 0, 10, 0},
		{"no_change", 50, 50, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeltaPercent(tt.old, tt.new_)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("DeltaPercent(%v, %v) = %v, want %v", tt.old, tt.new_, got, tt.want)
			}
		})
	}
}
