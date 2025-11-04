//go:build gonum

package poindexter

import "testing"

// 100k-size benchmarks run only in the gonum-tag job to keep CI time reasonable.

func BenchmarkNearest_Linear_Uniform_100k_2D(b *testing.B) {
	benchNearestBackend(b, 100_000, 2, BackendLinear, true, 0)
}
func BenchmarkNearest_Gonum_Uniform_100k_2D(b *testing.B) {
	benchNearestBackend(b, 100_000, 2, BackendGonum, true, 0)
}

func BenchmarkNearest_Linear_Uniform_100k_4D(b *testing.B) {
	benchNearestBackend(b, 100_000, 4, BackendLinear, true, 0)
}
func BenchmarkNearest_Gonum_Uniform_100k_4D(b *testing.B) {
	benchNearestBackend(b, 100_000, 4, BackendGonum, true, 0)
}

func BenchmarkNearest_Linear_Clustered_100k_2D(b *testing.B) {
	benchNearestBackend(b, 100_000, 2, BackendLinear, false, 3)
}
func BenchmarkNearest_Gonum_Clustered_100k_2D(b *testing.B) {
	benchNearestBackend(b, 100_000, 2, BackendGonum, false, 3)
}

func BenchmarkNearest_Linear_Clustered_100k_4D(b *testing.B) {
	benchNearestBackend(b, 100_000, 4, BackendLinear, false, 3)
}
func BenchmarkNearest_Gonum_Clustered_100k_4D(b *testing.B) {
	benchNearestBackend(b, 100_000, 4, BackendGonum, false, 3)
}
