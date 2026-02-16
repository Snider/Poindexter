package poindexter

// Factor is a value–weight pair for composite scoring.
type Factor struct {
	Value  float64
	Weight float64
}

// WeightedScore computes the weighted sum of factors.
// Each factor contributes Value * Weight to the total.
// Returns 0 for empty slices.
func WeightedScore(factors []Factor) float64 {
	var total float64
	for _, f := range factors {
		total += f.Value * f.Weight
	}
	return total
}

// Ratio returns part/whole safely. Returns 0 if whole is 0.
func Ratio(part, whole float64) float64 {
	if whole == 0 {
		return 0
	}
	return part / whole
}

// Delta returns the difference new_ - old.
func Delta(old, new_ float64) float64 {
	return new_ - old
}

// DeltaPercent returns the percentage change from old to new_.
// Returns 0 if old is 0.
func DeltaPercent(old, new_ float64) float64 {
	if old == 0 {
		return 0
	}
	return (new_ - old) / old * 100
}
