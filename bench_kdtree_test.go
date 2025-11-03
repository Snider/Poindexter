package poindexter

import (
	"fmt"
	"math/rand"
	"testing"
)

func makePoints(n, dim int) []KDPoint[int] {
	pts := make([]KDPoint[int], n)
	for i := 0; i < n; i++ {
		coords := make([]float64, dim)
		for d := 0; d < dim; d++ {
			coords[d] = rand.Float64()
		}
		pts[i] = KDPoint[int]{ID: fmt.Sprint(i), Coords: coords, Value: i}
	}
	return pts
}

func benchNearest(b *testing.B, n, dim int) {
	pts := makePoints(n, dim)
	tr, _ := NewKDTree(pts)
	q := make([]float64, dim)
	for i := range q {
		q[i] = 0.5
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = tr.Nearest(q)
	}
}

func benchKNearest(b *testing.B, n, dim, k int) {
	pts := makePoints(n, dim)
	tr, _ := NewKDTree(pts)
	q := make([]float64, dim)
	for i := range q {
		q[i] = 0.5
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tr.KNearest(q, k)
	}
}

func benchRadius(b *testing.B, n, dim int, r float64) {
	pts := makePoints(n, dim)
	tr, _ := NewKDTree(pts)
	q := make([]float64, dim)
	for i := range q {
		q[i] = 0.5
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tr.Radius(q, r)
	}
}

func BenchmarkNearest_1k_2D(b *testing.B)  { benchNearest(b, 1_000, 2) }
func BenchmarkNearest_10k_2D(b *testing.B) { benchNearest(b, 10_000, 2) }
func BenchmarkNearest_1k_4D(b *testing.B)  { benchNearest(b, 1_000, 4) }
func BenchmarkNearest_10k_4D(b *testing.B) { benchNearest(b, 10_000, 4) }

func BenchmarkKNearest10_1k_2D(b *testing.B)  { benchKNearest(b, 1_000, 2, 10) }
func BenchmarkKNearest10_10k_2D(b *testing.B) { benchKNearest(b, 10_000, 2, 10) }

func BenchmarkRadiusMid_1k_2D(b *testing.B)  { benchRadius(b, 1_000, 2, 0.5) }
func BenchmarkRadiusMid_10k_2D(b *testing.B) { benchRadius(b, 10_000, 2, 0.5) }
