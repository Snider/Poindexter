package poindexter

import "testing"

func TestKNearest_EdgeCases(t *testing.T) {
	pts := []KDPoint[int]{
		{ID: "a", Coords: []float64{0}},
	}
	tr, _ := NewKDTree(pts)
	// k <= 0 → nil
	ns, ds := tr.KNearest([]float64{0}, 0)
	if ns != nil || ds != nil {
		t.Fatalf("expected nil for k<=0, got %v %v", ns, ds)
	}
	// query-dim mismatch → nil
	ns, ds = tr.KNearest([]float64{0, 1}, 1)
	if ns != nil || ds != nil {
		t.Fatalf("expected nil for dim mismatch, got %v %v", ns, ds)
	}
}

func TestRadius_QueryDimMismatch(t *testing.T) {
	pts := []KDPoint[int]{{ID: "p", Coords: []float64{0}}}
	tr, _ := NewKDTree(pts)
	ns, ds := tr.Radius([]float64{0, 0}, 1)
	if ns != nil || ds != nil {
		t.Fatalf("expected nil for dim mismatch, got %v %v", ns, ds)
	}
}

func TestInsert_DimMismatch(t *testing.T) {
	tr, _ := NewKDTreeFromDim[int](2)
	ok := tr.Insert(KDPoint[int]{ID: "bad", Coords: []float64{0}}) // wrong dim
	if ok {
		t.Fatalf("expected false on insert with dim mismatch")
	}
	// inserting with empty ID should succeed and not touch idIndex
	ok = tr.Insert(KDPoint[int]{ID: "", Coords: []float64{0, 0}})
	if !ok {
		t.Fatalf("expected true on insert with empty ID and matching dim")
	}
}
