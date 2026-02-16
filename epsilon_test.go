package poindexter

import "testing"

func TestApproxEqual(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		epsilon float64
		want    bool
	}{
		{"equal", 1.0, 1.0, 0.01, true},
		{"close", 1.0, 1.005, 0.01, true},
		{"not_close", 1.0, 1.02, 0.01, false},
		{"negative", -1.0, -1.005, 0.01, true},
		{"zero", 0, 0.0001, 0.001, true},
		{"at_boundary", 1.0, 1.01, 0.01, false},
		{"large_epsilon", 100, 200, 150, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproxEqual(tt.a, tt.b, tt.epsilon)
			if got != tt.want {
				t.Errorf("ApproxEqual(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.epsilon, got, tt.want)
			}
		})
	}
}

func TestApproxZero(t *testing.T) {
	tests := []struct {
		name    string
		v       float64
		epsilon float64
		want    bool
	}{
		{"zero", 0, 0.01, true},
		{"small_pos", 0.005, 0.01, true},
		{"small_neg", -0.005, 0.01, true},
		{"not_zero", 0.02, 0.01, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproxZero(tt.v, tt.epsilon)
			if got != tt.want {
				t.Errorf("ApproxZero(%v, %v) = %v, want %v", tt.v, tt.epsilon, got, tt.want)
			}
		})
	}
}
