package poindexter

import "errors"

// Helper builders for KDTree points with min-max normalisation, optional inversion per-axis,
// and per-axis weights. These are convenience utilities to make it easy to map domain
// records into KD space for 2D/3D/4D use-cases.

// Errors for helper builders.
var (
	// ErrInvalidFeatures indicates that no features were provided or nil feature encountered.
	ErrInvalidFeatures = errors.New("kdtree: invalid features: provide at least one feature and ensure none are nil")
	// ErrInvalidWeights indicates weights length doesn't match features length.
	ErrInvalidWeights = errors.New("kdtree: invalid weights length; must match number of features")
	// ErrInvalidInvert indicates invert flags length doesn't match features length.
	ErrInvalidInvert = errors.New("kdtree: invalid invert length; must match number of features")
	// ErrStatsDimMismatch indicates NormStats dimensions do not match features length.
	ErrStatsDimMismatch = errors.New("kdtree: stats dimensionality mismatch")
)

// AxisStats holds the min/max observed for a single axis.
type AxisStats struct {
	Min float64
	Max float64
}

// NormStats holds per-axis normalisation statistics.
// For D dimensions, Stats has length D.
type NormStats struct {
	Stats []AxisStats
}

// ComputeNormStatsND computes per-axis min/max for an arbitrary number of features.
func ComputeNormStatsND[T any](items []T, features []func(T) float64) (NormStats, error) {
	if len(features) == 0 {
		return NormStats{}, ErrInvalidFeatures
	}
	// Initialise mins/maxes on first item where possible
	stats := make([]AxisStats, len(features))
	if len(items) == 0 {
		// empty items → zero stats slice of correct dim
		return NormStats{Stats: stats}, nil
	}
	// Seed with first item values
	first := items[0]
	for i, f := range features {
		if f == nil {
			return NormStats{}, ErrInvalidFeatures
		}
		v := f(first)
		stats[i] = AxisStats{Min: v, Max: v}
	}
	// Process remaining items
	for _, it := range items[1:] {
		for i, f := range features {
			v := f(it)
			if v < stats[i].Min {
				stats[i].Min = v
			}
			if v > stats[i].Max {
				stats[i].Max = v
			}
		}
	}
	return NormStats{Stats: stats}, nil
}

// BuildND constructs normalised-and-weighted KD points from arbitrary amount features.
// Features are min-max normalised per axis over the provided items, optionally inverted,
// then multiplied by per-axis weights.
func BuildND[T any](items []T, id func(T) string, features []func(T) float64, weights []float64, invert []bool) ([]KDPoint[T], error) {
	if len(items) == 0 {
		return nil, nil
	}
	if len(features) == 0 {
		return nil, ErrInvalidFeatures
	}
	if len(weights) != len(features) {
		return nil, ErrInvalidWeights
	}
	if len(invert) != len(features) {
		return nil, ErrInvalidInvert
	}
	stats, err := ComputeNormStatsND(items, features)
	if err != nil {
		return nil, err
	}
	return BuildNDWithStats(items, id, features, weights, invert, stats)
}

// BuildNDWithStats builds points using provided normalisation stats.
func BuildNDWithStats[T any](items []T, id func(T) string, features []func(T) float64, weights []float64, invert []bool, stats NormStats) ([]KDPoint[T], error) {
	if len(items) == 0 {
		return nil, nil
	}
	if len(features) == 0 {
		return nil, ErrInvalidFeatures
	}
	if len(weights) != len(features) {
		return nil, ErrInvalidWeights
	}
	if len(invert) != len(features) {
		return nil, ErrInvalidInvert
	}
	if len(stats.Stats) != len(features) {
		return nil, ErrStatsDimMismatch
	}
	pts := make([]KDPoint[T], len(items))
	for i, it := range items {
		coords := make([]float64, len(features))
		for d, f := range features {
			if f == nil {
				return nil, ErrInvalidFeatures
			}
			n := scale01(f(it), stats.Stats[d].Min, stats.Stats[d].Max)
			if invert[d] {
				n = 1 - n
			}
			coords[d] = weights[d] * n
		}
		var pid string
		if id != nil {
			pid = id(it)
		}
		pts[i] = KDPoint[T]{ID: pid, Value: it, Coords: coords}
	}
	return pts, nil
}

// minMax returns (min,max) of a slice.
func minMax(xs []float64) (float64, float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	mn, mx := xs[0], xs[0]
	for _, v := range xs[1:] {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return mn, mx
}

// scale01 maps v from [min,max] to [0,1]. If min==max, returns 0.
func scale01(v, min, max float64) float64 {
	if max == min {
		return 0
	}
	return (v - min) / (max - min)
}

// ComputeNormStats2D computes per-axis min/max for two features.
func ComputeNormStats2D[T any](items []T, f1, f2 func(T) float64) NormStats {
	vals1 := make([]float64, len(items))
	vals2 := make([]float64, len(items))
	for i, it := range items {
		vals1[i] = f1(it)
		vals2[i] = f2(it)
	}
	mn1, mx1 := minMax(vals1)
	mn2, mx2 := minMax(vals2)
	return NormStats{Stats: []AxisStats{{mn1, mx1}, {mn2, mx2}}}
}

// ComputeNormStats3D computes per-axis min/max for three features.
func ComputeNormStats3D[T any](items []T, f1, f2, f3 func(T) float64) NormStats {
	vals1 := make([]float64, len(items))
	vals2 := make([]float64, len(items))
	vals3 := make([]float64, len(items))
	for i, it := range items {
		vals1[i] = f1(it)
		vals2[i] = f2(it)
		vals3[i] = f3(it)
	}
	mn1, mx1 := minMax(vals1)
	mn2, mx2 := minMax(vals2)
	mn3, mx3 := minMax(vals3)
	return NormStats{Stats: []AxisStats{{mn1, mx1}, {mn2, mx2}, {mn3, mx3}}}
}

// ComputeNormStats4D computes per-axis min/max for four features.
func ComputeNormStats4D[T any](items []T, f1, f2, f3, f4 func(T) float64) NormStats {
	vals1 := make([]float64, len(items))
	vals2 := make([]float64, len(items))
	vals3 := make([]float64, len(items))
	vals4 := make([]float64, len(items))
	for i, it := range items {
		vals1[i] = f1(it)
		vals2[i] = f2(it)
		vals3[i] = f3(it)
		vals4[i] = f4(it)
	}
	mn1, mx1 := minMax(vals1)
	mn2, mx2 := minMax(vals2)
	mn3, mx3 := minMax(vals3)
	mn4, mx4 := minMax(vals4)
	return NormStats{Stats: []AxisStats{{mn1, mx1}, {mn2, mx2}, {mn3, mx3}, {mn4, mx4}}}
}

// Build2D constructs normalised-and-weighted KD points from items using two feature extractors.
// - id: function to provide a stable string ID (can return "" if you don't need DeleteByID)
// - f1,f2: feature extractors (raw values)
// - weights: per-axis weights applied after normalization
// - invert: per-axis flags; if true, the axis is inverted (1-norm) so that higher raw values become lower cost
func Build2D[T any](items []T, id func(T) string, f1, f2 func(T) float64, weights [2]float64, invert [2]bool) ([]KDPoint[T], error) {
	if len(items) == 0 {
		return nil, nil
	}
	vals1 := make([]float64, len(items))
	vals2 := make([]float64, len(items))
	for i, it := range items {
		vals1[i] = f1(it)
		vals2[i] = f2(it)
	}
	mn1, mx1 := minMax(vals1)
	mn2, mx2 := minMax(vals2)

	pts := make([]KDPoint[T], len(items))
	for i, it := range items {
		n1 := scale01(vals1[i], mn1, mx1)
		n2 := scale01(vals2[i], mn2, mx2)
		if invert[0] {
			n1 = 1 - n1
		}
		if invert[1] {
			n2 = 1 - n2
		}
		pts[i] = KDPoint[T]{
			ID:    id(it),
			Value: it,
			Coords: []float64{
				weights[0] * n1,
				weights[1] * n2,
			},
		}
	}
	return pts, nil
}

// Build2DWithStats builds points using provided normalisation stats.
func Build2DWithStats[T any](items []T, id func(T) string, f1, f2 func(T) float64, weights [2]float64, invert [2]bool, stats NormStats) ([]KDPoint[T], error) {
	if len(items) == 0 {
		return nil, nil
	}
	if len(stats.Stats) != 2 {
		return nil, ErrStatsDimMismatch
	}
	pts := make([]KDPoint[T], len(items))
	for i, it := range items {
		n1 := scale01(f1(it), stats.Stats[0].Min, stats.Stats[0].Max)
		n2 := scale01(f2(it), stats.Stats[1].Min, stats.Stats[1].Max)
		if invert[0] {
			n1 = 1 - n1
		}
		if invert[1] {
			n2 = 1 - n2
		}
		pts[i] = KDPoint[T]{
			ID:     id(it),
			Value:  it,
			Coords: []float64{weights[0] * n1, weights[1] * n2},
		}
	}
	return pts, nil
}

// Build3D constructs normalised-and-weighted KD points using three feature extractors.
func Build3D[T any](items []T, id func(T) string, f1, f2, f3 func(T) float64, weights [3]float64, invert [3]bool) ([]KDPoint[T], error) {
	if len(items) == 0 {
		return nil, nil
	}
	vals1 := make([]float64, len(items))
	vals2 := make([]float64, len(items))
	vals3 := make([]float64, len(items))
	for i, it := range items {
		vals1[i] = f1(it)
		vals2[i] = f2(it)
		vals3[i] = f3(it)
	}
	mn1, mx1 := minMax(vals1)
	mn2, mx2 := minMax(vals2)
	mn3, mx3 := minMax(vals3)

	pts := make([]KDPoint[T], len(items))
	for i, it := range items {
		n1 := scale01(vals1[i], mn1, mx1)
		n2 := scale01(vals2[i], mn2, mx2)
		n3 := scale01(vals3[i], mn3, mx3)
		if invert[0] {
			n1 = 1 - n1
		}
		if invert[1] {
			n2 = 1 - n2
		}
		if invert[2] {
			n3 = 1 - n3
		}
		pts[i] = KDPoint[T]{
			ID:    id(it),
			Value: it,
			Coords: []float64{
				weights[0] * n1,
				weights[1] * n2,
				weights[2] * n3,
			},
		}
	}
	return pts, nil
}

// Build3DWithStats builds points using provided normalisation stats.
func Build3DWithStats[T any](items []T, id func(T) string, f1, f2, f3 func(T) float64, weights [3]float64, invert [3]bool, stats NormStats) ([]KDPoint[T], error) {
	if len(items) == 0 {
		return nil, nil
	}
	if len(stats.Stats) != 3 {
		return nil, ErrStatsDimMismatch
	}
	pts := make([]KDPoint[T], len(items))
	for i, it := range items {
		n1 := scale01(f1(it), stats.Stats[0].Min, stats.Stats[0].Max)
		n2 := scale01(f2(it), stats.Stats[1].Min, stats.Stats[1].Max)
		n3 := scale01(f3(it), stats.Stats[2].Min, stats.Stats[2].Max)
		if invert[0] {
			n1 = 1 - n1
		}
		if invert[1] {
			n2 = 1 - n2
		}
		if invert[2] {
			n3 = 1 - n3
		}
		pts[i] = KDPoint[T]{
			ID:     id(it),
			Value:  it,
			Coords: []float64{weights[0] * n1, weights[1] * n2, weights[2] * n3},
		}
	}
	return pts, nil
}

// Build4D constructs normalised-and-weighted KD points using four feature extractors.
func Build4D[T any](items []T, id func(T) string, f1, f2, f3, f4 func(T) float64, weights [4]float64, invert [4]bool) ([]KDPoint[T], error) {
	if len(items) == 0 {
		return nil, nil
	}
	vals1 := make([]float64, len(items))
	vals2 := make([]float64, len(items))
	vals3 := make([]float64, len(items))
	vals4 := make([]float64, len(items))
	for i, it := range items {
		vals1[i] = f1(it)
		vals2[i] = f2(it)
		vals3[i] = f3(it)
		vals4[i] = f4(it)
	}
	mn1, mx1 := minMax(vals1)
	mn2, mx2 := minMax(vals2)
	mn3, mx3 := minMax(vals3)
	mn4, mx4 := minMax(vals4)

	pts := make([]KDPoint[T], len(items))
	for i, it := range items {
		n1 := scale01(vals1[i], mn1, mx1)
		n2 := scale01(vals2[i], mn2, mx2)
		n3 := scale01(vals3[i], mn3, mx3)
		n4 := scale01(vals4[i], mn4, mx4)
		if invert[0] {
			n1 = 1 - n1
		}
		if invert[1] {
			n2 = 1 - n2
		}
		if invert[2] {
			n3 = 1 - n3
		}
		if invert[3] {
			n4 = 1 - n4
		}
		pts[i] = KDPoint[T]{
			ID:    id(it),
			Value: it,
			Coords: []float64{
				weights[0] * n1,
				weights[1] * n2,
				weights[2] * n3,
				weights[3] * n4,
			},
		}
	}
	return pts, nil
}

// Build4DWithStats builds points using provided normalisation stats.
func Build4DWithStats[T any](items []T, id func(T) string, f1, f2, f3, f4 func(T) float64, weights [4]float64, invert [4]bool, stats NormStats) ([]KDPoint[T], error) {
	if len(items) == 0 {
		return nil, nil
	}
	if len(stats.Stats) != 4 {
		return nil, ErrStatsDimMismatch
	}
	pts := make([]KDPoint[T], len(items))
	for i, it := range items {
		n1 := scale01(f1(it), stats.Stats[0].Min, stats.Stats[0].Max)
		n2 := scale01(f2(it), stats.Stats[1].Min, stats.Stats[1].Max)
		n3 := scale01(f3(it), stats.Stats[2].Min, stats.Stats[2].Max)
		n4 := scale01(f4(it), stats.Stats[3].Min, stats.Stats[3].Max)
		if invert[0] {
			n1 = 1 - n1
		}
		if invert[1] {
			n2 = 1 - n2
		}
		if invert[2] {
			n3 = 1 - n3
		}
		if invert[3] {
			n4 = 1 - n4
		}
		pts[i] = KDPoint[T]{
			ID:     id(it),
			Value:  it,
			Coords: []float64{weights[0] * n1, weights[1] * n2, weights[2] * n3, weights[3] * n4},
		}
	}
	return pts, nil
}
