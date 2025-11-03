package poindexter

import (
	"fmt"
	"testing"
)

func TestCosineDistance_Basics(t *testing.T) {
	// identical vectors → distance 0
	a := []float64{1, 0, 0}
	b := []float64{1, 0, 0}
	d := CosineDistance{}.Distance(a, b)
	if d != 0 {
		t.Fatalf("expected 0, got %v", d)
	}
	// orthogonal → distance 1
	b = []float64{0, 1, 0}
	d = CosineDistance{}.Distance(a, b)
	if d < 0.999 || d > 1.001 {
		t.Fatalf("expected ~1, got %v", d)
	}
	// opposite → distance 2
	a = []float64{1, 0}
	b = []float64{-1, 0}
	d = CosineDistance{}.Distance(a, b)
	if d < 1.999 || d > 2.001 {
		t.Fatalf("expected ~2, got %v", d)
	}
	// zero vectors
	a = []float64{0, 0}
	b = []float64{0, 0}
	d = CosineDistance{}.Distance(a, b)
	if d != 0 {
		t.Fatalf("both zero → 0, got %v", d)
	}
	// one zero
	a = []float64{0, 0}
	b = []float64{1, 2}
	d = CosineDistance{}.Distance(a, b)
	if d != 1 {
		t.Fatalf("one zero → 1, got %v", d)
	}
}

func TestWeightedCosineDistance_Basics(t *testing.T) {
	w := WeightedCosineDistance{Weights: []float64{2, 0.5}}
	a := []float64{1, 0}
	b := []float64{1, 0}
	d := w.Distance(a, b)
	if d != 0 {
		t.Fatalf("expected 0, got %v", d)
	}
	// orthogonal remains ~1 regardless of weights for these axes
	b = []float64{0, 3}
	d = w.Distance(a, b)
	if d < 0.999 || d > 1.001 {
		t.Fatalf("expected ~1, got %v", d)
	}
}

func TestBuildND_ParityWithBuild4D(t *testing.T) {
	type rec struct{ a, b, c, d float64 }
	items := []rec{{0, 10, 100, 1}, {10, 20, 200, 2}, {5, 15, 150, 1.5}}
	weights4 := [4]float64{1.0, 0.5, 2.0, 1.0}
	invert4 := [4]bool{false, true, false, true}
	pts4, err := Build4D(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		func(r rec) float64 { return r.c },
		func(r rec) float64 { return r.d },
		weights4, invert4,
	)
	if err != nil {
		t.Fatalf("build4d err: %v", err)
	}

	features := []func(rec) float64{
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		func(r rec) float64 { return r.c },
		func(r rec) float64 { return r.d },
	}
	wts := []float64{weights4[0], weights4[1], weights4[2], weights4[3]}
	inv := []bool{invert4[0], invert4[1], invert4[2], invert4[3]}
	ptsN, err := BuildND(items, func(r rec) string { return "" }, features, wts, inv)
	if err != nil {
		t.Fatalf("buildND err: %v", err)
	}
	if len(ptsN) != len(pts4) {
		t.Fatalf("len mismatch")
	}
	for i := range ptsN {
		if len(ptsN[i].Coords) != 4 {
			t.Fatalf("dim != 4")
		}
		for d := 0; d < 4; d++ {
			if fmt.Sprintf("%.6f", ptsN[i].Coords[d]) != fmt.Sprintf("%.6f", pts4[i].Coords[d]) {
				t.Fatalf("coords mismatch at i=%d d=%d: %v vs %v", i, d, ptsN[i].Coords, pts4[i].Coords)
			}
		}
	}
}

func TestBuildNDWithStats_Errors(t *testing.T) {
	type rec struct{ x float64 }
	items := []rec{{1}, {2}}
	features := []func(rec) float64{func(r rec) float64 { return r.x }}
	wts := []float64{1}
	inv := []bool{false}
	// stats dim mismatch
	stats := NormStats{Stats: []AxisStats{{Min: 0, Max: 1}, {Min: 0, Max: 1}}}
	_, err := BuildNDWithStats(items, func(r rec) string { return "" }, features, wts, inv, stats)
	if err == nil {
		t.Fatalf("expected error for stats dim mismatch")
	}
}
