//go:build !gonum

package poindexter

// hasGonum reports whether the gonum backend is compiled in (build tag 'gonum').
func hasGonum() bool { return false }

// buildGonumBackend is unavailable without the 'gonum' build tag.
func buildGonumBackend[T any](pts []KDPoint[T], metric DistanceMetric) (any, error) {
	return nil, ErrEmptyPoints // sentinel non-nil error to force fallback
}

func gonumNearest[T any](backend any, query []float64) (int, float64, bool) {
	return -1, 0, false
}

func gonumKNearest[T any](backend any, query []float64, k int) ([]int, []float64) {
	return nil, nil
}

func gonumRadius[T any](backend any, query []float64, r float64) ([]int, []float64) {
	return nil, nil
}
