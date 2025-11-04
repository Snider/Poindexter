package poindexter

import (
	"fmt"
	"math/rand"
	"testing"
)

// deterministicRand returns a rand.Rand with a fixed seed for reproducible datasets.
func deterministicRand() *rand.Rand { return rand.New(rand.NewSource(42)) }

func makeUniformPoints(n, dim int) []KDPoint[int] {
	r := deterministicRand()
	pts := make([]KDPoint[int], n)
	for i := 0; i < n; i++ {
		coords := make([]float64, dim)
		for d := 0; d < dim; d++ {
			coords[d] = r.Float64()
		}
		pts[i] = KDPoint[int]{ID: fmt.Sprint(i), Coords: coords, Value: i}
	}
	return pts
}

// makeClusteredPoints creates n points around c clusters with small variance.
func makeClusteredPoints(n, dim, c int) []KDPoint[int] {
	if c <= 0 {
		c = 1
	}
	r := deterministicRand()
	centers := make([][]float64, c)
	for i := 0; i < c; i++ {
		centers[i] = make([]float64, dim)
		for d := 0; d < dim; d++ {
			centers[i][d] = r.Float64()
		}
	}
	pts := make([]KDPoint[int], n)
	for i := 0; i < n; i++ {
		coords := make([]float64, dim)
		cent := centers[r.Intn(c)]
		for d := 0; d < dim; d++ {
			// small gaussian noise around center (Box-Muller)
			u1 := r.Float64()
			u2 := r.Float64()
			z := (rand.NormFloat64()) // uses global; fine for test speed
			_ = u1
			_ = u2
			coords[d] = cent[d] + 0.03*z
			if coords[d] < 0 {
				coords[d] = 0
			} else if coords[d] > 1 {
				coords[d] = 1
			}
		}
		pts[i] = KDPoint[int]{ID: fmt.Sprint(i), Coords: coords, Value: i}
	}
	return pts
}

func benchNearestBackend(b *testing.B, n, dim int, backend KDBackend, uniform bool, clusters int) {
	var pts []KDPoint[int]
	if uniform {
		pts = makeUniformPoints(n, dim)
	} else {
		pts = makeClusteredPoints(n, dim, clusters)
	}
	tr, _ := NewKDTree(pts, WithBackend(backend))
	q := make([]float64, dim)
	for i := range q {
		q[i] = 0.5
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = tr.Nearest(q)
	}
}

func benchKNNBackend(b *testing.B, n, dim, k int, backend KDBackend, uniform bool, clusters int) {
	var pts []KDPoint[int]
	if uniform {
		pts = makeUniformPoints(n, dim)
	} else {
		pts = makeClusteredPoints(n, dim, clusters)
	}
	tr, _ := NewKDTree(pts, WithBackend(backend))
	q := make([]float64, dim)
	for i := range q {
		q[i] = 0.5
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tr.KNearest(q, k)
	}
}

func benchRadiusBackend(b *testing.B, n, dim int, r float64, backend KDBackend, uniform bool, clusters int) {
	var pts []KDPoint[int]
	if uniform {
		pts = makeUniformPoints(n, dim)
	} else {
		pts = makeClusteredPoints(n, dim, clusters)
	}
	tr, _ := NewKDTree(pts, WithBackend(backend))
	q := make([]float64, dim)
	for i := range q {
		q[i] = 0.5
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tr.Radius(q, r)
	}
}

// Uniform 2D/4D, Linear vs Gonum (opt-in via build tag; falls back to linear if not available)
func BenchmarkNearest_Linear_Uniform_1k_2D(b *testing.B) {
	benchNearestBackend(b, 1_000, 2, BackendLinear, true, 0)
}
func BenchmarkNearest_Gonum_Uniform_1k_2D(b *testing.B) {
	benchNearestBackend(b, 1_000, 2, BackendGonum, true, 0)
}
func BenchmarkNearest_Linear_Uniform_10k_2D(b *testing.B) {
	benchNearestBackend(b, 10_000, 2, BackendLinear, true, 0)
}
func BenchmarkNearest_Gonum_Uniform_10k_2D(b *testing.B) {
	benchNearestBackend(b, 10_000, 2, BackendGonum, true, 0)
}

func BenchmarkNearest_Linear_Uniform_1k_4D(b *testing.B) {
	benchNearestBackend(b, 1_000, 4, BackendLinear, true, 0)
}
func BenchmarkNearest_Gonum_Uniform_1k_4D(b *testing.B) {
	benchNearestBackend(b, 1_000, 4, BackendGonum, true, 0)
}
func BenchmarkNearest_Linear_Uniform_10k_4D(b *testing.B) {
	benchNearestBackend(b, 10_000, 4, BackendLinear, true, 0)
}
func BenchmarkNearest_Gonum_Uniform_10k_4D(b *testing.B) {
	benchNearestBackend(b, 10_000, 4, BackendGonum, true, 0)
}

// Clustered 2D/4D (3 clusters)
func BenchmarkNearest_Linear_Clustered_1k_2D(b *testing.B) {
	benchNearestBackend(b, 1_000, 2, BackendLinear, false, 3)
}
func BenchmarkNearest_Gonum_Clustered_1k_2D(b *testing.B) {
	benchNearestBackend(b, 1_000, 2, BackendGonum, false, 3)
}
func BenchmarkNearest_Linear_Clustered_10k_2D(b *testing.B) {
	benchNearestBackend(b, 10_000, 2, BackendLinear, false, 3)
}
func BenchmarkNearest_Gonum_Clustered_10k_2D(b *testing.B) {
	benchNearestBackend(b, 10_000, 2, BackendGonum, false, 3)
}

func BenchmarkKNN10_Linear_Uniform_10k_2D(b *testing.B) {
	benchKNNBackend(b, 10_000, 2, 10, BackendLinear, true, 0)
}
func BenchmarkKNN10_Gonum_Uniform_10k_2D(b *testing.B) {
	benchKNNBackend(b, 10_000, 2, 10, BackendGonum, true, 0)
}
func BenchmarkKNN10_Linear_Clustered_10k_2D(b *testing.B) {
	benchKNNBackend(b, 10_000, 2, 10, BackendLinear, false, 3)
}
func BenchmarkKNN10_Gonum_Clustered_10k_2D(b *testing.B) {
	benchKNNBackend(b, 10_000, 2, 10, BackendGonum, false, 3)
}

func BenchmarkRadiusMid_Linear_Uniform_10k_2D(b *testing.B) {
	benchRadiusBackend(b, 10_000, 2, 0.5, BackendLinear, true, 0)
}
func BenchmarkRadiusMid_Gonum_Uniform_10k_2D(b *testing.B) {
	benchRadiusBackend(b, 10_000, 2, 0.5, BackendGonum, true, 0)
}
func BenchmarkRadiusMid_Linear_Clustered_10k_2D(b *testing.B) {
	benchRadiusBackend(b, 10_000, 2, 0.5, BackendLinear, false, 3)
}
func BenchmarkRadiusMid_Gonum_Clustered_10k_2D(b *testing.B) {
	benchRadiusBackend(b, 10_000, 2, 0.5, BackendGonum, false, 3)
}
