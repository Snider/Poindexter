package poindexter

import (
	"errors"
	"testing"
)

func TestNewKDTree_Errors(t *testing.T) {
	// empty points
	if _, err := NewKDTree[string](nil); !errors.Is(err, ErrEmptyPoints) {
		t.Fatalf("want ErrEmptyPoints, got %v", err)
	}
	// zero-dim
	pts0 := []KDPoint[string]{{ID: "A", Coords: nil}}
	if _, err := NewKDTree(pts0); !errors.Is(err, ErrZeroDim) {
		t.Fatalf("want ErrZeroDim, got %v", err)
	}
	// dim mismatch
	ptsDim := []KDPoint[string]{
		{ID: "A", Coords: []float64{0}},
		{ID: "B", Coords: []float64{0, 1}},
	}
	if _, err := NewKDTree(ptsDim); !errors.Is(err, ErrDimMismatch) {
		t.Fatalf("want ErrDimMismatch, got %v", err)
	}
	// duplicate IDs
	ptsDup := []KDPoint[string]{
		{ID: "X", Coords: []float64{0}},
		{ID: "X", Coords: []float64{1}},
	}
	if _, err := NewKDTree(ptsDup); !errors.Is(err, ErrDuplicateID) {
		t.Fatalf("want ErrDuplicateID, got %v", err)
	}
}

func TestDeleteByID_NotFound(t *testing.T) {
	pts := []KDPoint[int]{
		{ID: "A", Coords: []float64{0}, Value: 1},
	}
	tr, err := NewKDTree(pts)
	if err != nil {
		t.Fatalf("NewKDTree err: %v", err)
	}
	if tr.DeleteByID("NOPE") {
		t.Fatalf("expected false for missing ID")
	}
}

func TestKNearest_KGreaterThanN(t *testing.T) {
	pts := []KDPoint[int]{
		{ID: "a", Coords: []float64{0}},
		{ID: "b", Coords: []float64{2}},
	}
	tr, _ := NewKDTree(pts)
	ns, ds := tr.KNearest([]float64{1}, 5)
	if len(ns) != 2 || len(ds) != 2 {
		t.Fatalf("want 2 neighbors, got %d", len(ns))
	}
	if !(ds[0] <= ds[1]) {
		t.Fatalf("distances not sorted: %v", ds)
	}
}

func TestRadius_BoundaryAndZero(t *testing.T) {
	pts := []KDPoint[int]{
		{ID: "o", Coords: []float64{0}},
		{ID: "one", Coords: []float64{1}},
	}
	tr, _ := NewKDTree(pts, WithMetric(EuclideanDistance{}))
	// radius exactly includes point at distance 1
	within, _ := tr.Radius([]float64{0}, 1)
	foundOne := false
	for _, p := range within {
		if p.ID == "one" {
			foundOne = true
		}
	}
	if !foundOne {
		t.Fatalf("expected to include point at exact radius")
	}
	// radius zero should include exact match only
	within0, _ := tr.Radius([]float64{0}, 0)
	if len(within0) == 0 || within0[0].ID != "o" {
		t.Fatalf("expected only origin at r=0, got %v", within0)
	}
}

func TestNewKDTreeFromDim_WithMetric_InsertQuery(t *testing.T) {
	tr, err := NewKDTreeFromDim[string](2, WithMetric(ManhattanDistance{}))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ok := tr.Insert(KDPoint[string]{ID: "A", Coords: []float64{0, 0}, Value: "a"})
	if !ok {
		t.Fatalf("insert failed")
	}
	tr.Insert(KDPoint[string]{ID: "B", Coords: []float64{2, 2}, Value: "b"})
	p, d, ok := tr.Nearest([]float64{1, 0})
	if !ok || p.ID != "A" {
		t.Fatalf("expected A nearest, got %v", p)
	}
	if d != 1 { // ManhattanDistance from (1,0) to (0,0) is 1
		t.Fatalf("expected manhattan distance 1, got %v", d)
	}
}

func TestNearest_QueryDimMismatch(t *testing.T) {
	pts := []KDPoint[int]{
		{ID: "a", Coords: []float64{0, 0}},
	}
	tr, _ := NewKDTree(pts)
	_, _, ok := tr.Nearest([]float64{0})
	if ok {
		t.Fatalf("expected ok=false for query dim mismatch")
	}
}
