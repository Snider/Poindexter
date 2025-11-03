package poindexter

import (
	"fmt"
	"testing"
)

func TestInsert_DuplicateID(t *testing.T) {
	tr, err := NewKDTreeFromDim[string](1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ok := tr.Insert(KDPoint[string]{ID: "X", Coords: []float64{0}})
	if !ok {
		t.Fatalf("first insert should succeed")
	}
	// duplicate ID should fail
	if tr.Insert(KDPoint[string]{ID: "X", Coords: []float64{1}}) {
		t.Fatalf("expected insert duplicate ID to return false")
	}
}

func TestDeleteByID_SwapDelete(t *testing.T) {
	// Arrange 3 points so that deleting the middle triggers swap-delete path
	pts := []KDPoint[int]{
		{ID: "A", Coords: []float64{0}},
		{ID: "B", Coords: []float64{1}},
		{ID: "C", Coords: []float64{2}},
	}
	tr, err := NewKDTree(pts)
	if err != nil {
		t.Fatalf("NewKDTree err: %v", err)
	}
	if !tr.DeleteByID("B") {
		t.Fatalf("delete B failed")
	}
	if tr.Len() != 2 {
		t.Fatalf("expected len 2, got %d", tr.Len())
	}
	// Ensure B is gone and A/C remain reachable
	ids := make(map[string]bool)
	for _, q := range [][]float64{{0}, {2}} {
		p, _, ok := tr.Nearest(q)
		if ok {
			ids[p.ID] = true
		}
	}
	if ids["B"] {
		t.Fatalf("B should not be present after delete")
	}
	if !(ids["A"] || ids["C"]) {
		t.Fatalf("expected either A or C to be nearest for respective queries: %v", ids)
	}
}

func TestRadius_NegativeReturnsNil(t *testing.T) {
	pts := []KDPoint[int]{{ID: "z", Coords: []float64{0}}}
	tr, _ := NewKDTree(pts)
	ns, ds := tr.Radius([]float64{0}, -1)
	if ns != nil || ds != nil {
		// Both should be nil on invalid radius
		t.Fatalf("expected nil slices on negative radius, got %v %v", ns, ds)
	}
}

func TestNearest_EmptyTree(t *testing.T) {
	tr, _ := NewKDTreeFromDim[int](2)
	_, _, ok := tr.Nearest([]float64{0, 0})
	if ok {
		t.Fatalf("expected ok=false for empty tree")
	}
}

func TestWeightedCosineMetric_ViaKDTree(t *testing.T) {
	// Two points oriented differently around the query; ensure call path exercised
	type rec struct{ a, b float64 }
	items := []rec{{1, 0}, {0, 1}}
	weights := []float64{1, 2}
	invert := []bool{false, false}
	features := []func(rec) float64{
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
	}
	pts, err := BuildND(items, func(r rec) string { return fmt.Sprintf("%v", r) }, features, weights, invert)
	if err != nil {
		t.Fatalf("buildND err: %v", err)
	}
	tr, err := NewKDTree(pts, WithMetric(WeightedCosineDistance{Weights: weights}))
	if err != nil {
		t.Fatalf("kdt err: %v", err)
	}
	q := []float64{0.5 * weights[0], 0.5 * weights[1]} // mid direction
	_, d, ok := tr.Nearest(q)
	if !ok {
		t.Fatalf("no nearest")
	}
	if d < 0 || d > 2 {
		t.Fatalf("cosine distance out of bounds: %v", d)
	}
}
