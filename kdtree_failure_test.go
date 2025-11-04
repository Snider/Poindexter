package poindexter

import (
	"testing"
)

func TestNewKDTree_Empty(t *testing.T) {
	_, err := NewKDTree[any](nil)
	if err != ErrEmptyPoints {
		t.Errorf("expected ErrEmptyPoints, got %v", err)
	}
}

func TestWasmSmokeEquivalent(t *testing.T) {
	tree, err := NewKDTreeFromDim[string](2)
	if err != nil {
		t.Fatalf("NewKDTreeFromDim failed: %v", err)
	}
	tree.Insert(KDPoint[string]{ID: "a", Coords: []float64{0, 0}, Value: "A"})
	tree.Insert(KDPoint[string]{ID: "b", Coords: []float64{1, 0}, Value: "B"})
	_, _, found := tree.Nearest([]float64{0.9, 0.1})
	if !found {
		t.Error("expected to find a nearest point")
	}
}

func TestNewKDTree_ZeroDim(t *testing.T) {
	pts := []KDPoint[any]{{Coords: []float64{}}}
	_, err := NewKDTree[any](pts)
	if err != ErrZeroDim {
		t.Errorf("expected ErrZeroDim, got %v", err)
	}
}

func TestNewKDTree_DimMismatch(t *testing.T) {
	pts := []KDPoint[any]{
		{Coords: []float64{1, 2}},
		{Coords: []float64{3}},
	}
	_, err := NewKDTree[any](pts)
	if err != ErrDimMismatch {
		t.Errorf("expected ErrDimMismatch, got %v", err)
	}
}

func TestNewKDTree_DuplicateID(t *testing.T) {
	pts := []KDPoint[any]{
		{ID: "a", Coords: []float64{1, 2}},
		{ID: "a", Coords: []float64{3, 4}},
	}
	_, err := NewKDTree[any](pts)
	if err != ErrDuplicateID {
		t.Errorf("expected ErrDuplicateID, got %v", err)
	}
}

func TestBuildND_InvalidFeatures(t *testing.T) {
	_, err := BuildND[struct{}]([]struct{}{{}}, func(a struct{}) string { return "" }, nil, nil, nil)
	if err != ErrInvalidFeatures {
		t.Errorf("expected ErrInvalidFeatures, got %v", err)
	}
}

func TestBuildND_InvalidWeights(t *testing.T) {
	_, err := BuildND[struct{}]([]struct{}{{}}, func(a struct{}) string { return "" }, []func(struct{}) float64{func(a struct{}) float64 { return 0 }}, nil, nil)
	if err != ErrInvalidWeights {
		t.Errorf("expected ErrInvalidWeights, got %v", err)
	}
}

func TestBuildND_InvalidInvert(t *testing.T) {
	_, err := BuildND[struct{}]([]struct{}{{}}, func(a struct{}) string { return "" }, []func(struct{}) float64{func(a struct{}) float64 { return 0 }}, []float64{0}, nil)
	if err != ErrInvalidInvert {
		t.Errorf("expected ErrInvalidInvert, got %v", err)
	}
}

func TestBuildNDWithStats_DimMismatch(t *testing.T) {
	_, err := BuildNDWithStats[struct{}]([]struct{}{{}}, func(a struct{}) string { return "" }, []func(struct{}) float64{func(a struct{}) float64 { return 0 }}, []float64{0}, []bool{false}, NormStats{})
	if err != ErrStatsDimMismatch {
		t.Errorf("expected ErrStatsDimMismatch, got %v", err)
	}
}

func TestGonumStub_BuildBackend(t *testing.T) {
	_, err := buildGonumBackend[any](nil, nil)
	if err != ErrEmptyPoints {
		t.Errorf("expected ErrEmptyPoints, got %v", err)
	}
}
