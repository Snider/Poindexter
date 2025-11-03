package poindexter

import (
	"math/rand"
	"testing"
)

// FuzzKDTreeNearest_NoPanic ensures Nearest never panics and distances are non-negative.
func FuzzKDTreeNearest_NoPanic(f *testing.F) {
	// Seed with small cases
	f.Add(3, 2)
	f.Add(5, 4)
	f.Fuzz(func(t *testing.T, n int, dim int) {
		if n <= 0 {
			n = 1
		}
		if n > 64 {
			n = 64
		}
		if dim <= 0 {
			dim = 1
		}
		if dim > 8 {
			dim = 8
		}

		pts := make([]KDPoint[int], n)
		for i := 0; i < n; i++ {
			coords := make([]float64, dim)
			for d := 0; d < dim; d++ {
				coords[d] = rand.Float64()*100 - 50
			}
			pts[i] = KDPoint[int]{ID: "", Coords: coords, Value: i}
		}
		tr, err := NewKDTree(pts)
		if err != nil {
			t.Skip()
		}
		q := make([]float64, dim)
		for d := range q {
			q[d] = rand.Float64()*100 - 50
		}
		_, dist, _ := tr.Nearest(q)
		if dist < 0 {
			t.Fatalf("negative distance: %v", dist)
		}
	})
}

// FuzzMetrics_NoNegative checks Manhattan, Euclidean, Chebyshev don't return negatives for random inputs.
func FuzzMetrics_NoNegative(f *testing.F) {
	f.Add(2)
	f.Add(4)
	f.Fuzz(func(t *testing.T, dim int) {
		if dim <= 0 {
			dim = 1
		}
		if dim > 8 {
			dim = 8
		}
		a := make([]float64, dim)
		b := make([]float64, dim)
		for i := 0; i < dim; i++ {
			a[i] = rand.Float64()*10 - 5
			b[i] = rand.Float64()*10 - 5
		}
		m1 := EuclideanDistance{}.Distance(a, b)
		m2 := ManhattanDistance{}.Distance(a, b)
		m3 := ChebyshevDistance{}.Distance(a, b)
		m4 := CosineDistance{}.Distance(a, b)
		w := make([]float64, dim)
		for i := range w {
			w[i] = 1
		}
		m5 := WeightedCosineDistance{Weights: w}.Distance(a, b)
		if m1 < 0 || m2 < 0 || m3 < 0 || m4 < 0 || m5 < 0 {
			t.Fatalf("negative metric: %v %v %v %v %v", m1, m2, m3, m4, m5)
		}
		if m4 > 2 || m5 > 2 {
			t.Fatalf("cosine distance out of bounds: %v %v", m4, m5)
		}
	})
}

// FuzzDimensionMismatch_NoPanic ensures queries with wrong dims return ok=false and not panic.
func FuzzDimensionMismatch_NoPanic(f *testing.F) {
	f.Add(3, 2, 1)
	f.Fuzz(func(t *testing.T, n, dim, qdim int) {
		if n <= 0 {
			n = 1
		}
		if n > 32 {
			n = 32
		}
		if dim <= 0 {
			dim = 1
		}
		if dim > 6 {
			dim = 6
		}
		if qdim < 0 {
			qdim = 0
		}
		if qdim > 6 {
			qdim = 6
		}
		pts := make([]KDPoint[int], n)
		for i := 0; i < n; i++ {
			coords := make([]float64, dim)
			pts[i] = KDPoint[int]{Coords: coords}
		}
		tr, err := NewKDTree(pts)
		if err != nil {
			t.Skip()
		}
		q := make([]float64, qdim)
		_, _, ok := tr.Nearest(q)
		if qdim != dim && ok {
			t.Fatalf("expected ok=false for dim mismatch; dim=%d qdim=%d", dim, qdim)
		}
	})
}
