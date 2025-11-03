package poindexter

// Helper builders for KDTree points with min-max normalization, optional inversion per-axis,
// and per-axis weights. These are convenience utilities to make it easy to map domain
// records into KD space for 2D/3D/4D use-cases.

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

// Build2D constructs normalized-and-weighted KD points from items using two feature extractors.
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

// Build3D constructs normalized-and-weighted KD points using three feature extractors.
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

// Build4D constructs normalized-and-weighted KD points using four feature extractors.
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
