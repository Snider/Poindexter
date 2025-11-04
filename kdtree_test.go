package poindexter

import (
	"testing"
)

func samplePoints() []KDPoint[string] {
	return []KDPoint[string]{
		{ID: "A", Coords: []float64{0, 0}, Value: "alpha"},
		{ID: "B", Coords: []float64{1, 0}, Value: "bravo"},
		{ID: "C", Coords: []float64{0, 1}, Value: "charlie"},
		{ID: "D", Coords: []float64{1, 1}, Value: "delta"},
		{ID: "E", Coords: []float64{2, 2}, Value: "echo"},
	}
}

func TestKDTree_Nearest(t *testing.T) {
	pts := samplePoints()
	tree, err := NewKDTree(pts, WithMetric(EuclideanDistance{}))
	if err != nil {
		t.Fatalf("NewKDTree error: %v", err)
	}

	p, dist, ok := tree.Nearest([]float64{0.9, 0.9})
	if !ok {
		t.Fatalf("expected a nearest neighbor")
	}
	if p.ID != "D" {
		t.Fatalf("expected D, got %s", p.ID)
	}
	if dist <= 0 {
		t.Fatalf("expected positive distance, got %v", dist)
	}
}

func TestKDTree_KNearest(t *testing.T) {
	pts := samplePoints()
	tree, err := NewKDTree(pts, WithMetric(ManhattanDistance{}))
	if err != nil {
		t.Fatalf("NewKDTree error: %v", err)
	}

	neighbors, dists := tree.KNearest([]float64{0.9, 0.9}, 3)
	if len(neighbors) != 3 || len(dists) != 3 {
		t.Fatalf("expected 3 neighbors, got %d", len(neighbors))
	}
	if neighbors[0].ID != "D" {
		t.Fatalf("expected first neighbor D, got %s", neighbors[0].ID)
	}
}

func TestKDTree_Radius(t *testing.T) {
	pts := samplePoints()
	tree, err := NewKDTree(pts, WithMetric(EuclideanDistance{}))
	if err != nil {
		t.Fatalf("NewKDTree error: %v", err)
	}

	neighbors, dists := tree.Radius([]float64{0, 0}, 1.01)
	if len(neighbors) < 2 {
		t.Fatalf("expected at least 2 neighbors within radius, got %d", len(neighbors))
	}
	// distances should be non-decreasing
	for i := 1; i < len(dists); i++ {
		if dists[i] < dists[i-1] {
			t.Fatalf("distances not sorted: %v", dists)
		}
	}
}

func TestKDTree_InsertDelete(t *testing.T) {
	pts := samplePoints()
	tree, err := NewKDTree(pts)
	if err != nil {
		t.Fatalf("NewKDTree error: %v", err)
	}
	// Insert a new close point near (0,0)
	ok := tree.Insert(KDPoint[string]{ID: "Z", Coords: []float64{0.05, 0.05}, Value: "zulu"})
	if !ok {
		t.Fatalf("insert failed")
	}
	p, _, found := tree.Nearest([]float64{0.04, 0.04})
	if !found || p.ID != "Z" {
		t.Fatalf("expected nearest to be Z after insert, got %+v", p)
	}

	// Delete and verify nearest changes back
	if !tree.DeleteByID("Z") {
		t.Fatalf("delete failed")
	}
	p, _, found = tree.Nearest([]float64{0.04, 0.04})
	if !found || p.ID != "A" {
		t.Fatalf("expected nearest to be A after delete, got %+v", p)
	}
}

func TestKDTree_DimAndLen(t *testing.T) {
	pts := samplePoints()
	tree, err := NewKDTree(pts)
	if err != nil {
		t.Fatalf("NewKDTree error: %v", err)
	}
	if tree.Len() != len(pts) {
		t.Fatalf("Len mismatch: %d vs %d", tree.Len(), len(pts))
	}
	if tree.Dim() != 2 {
		t.Fatalf("Dim mismatch: %d", tree.Dim())
	}
}
