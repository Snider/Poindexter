package main

import (
	"testing"

	poindexter "github.com/Snider/Poindexter"
)

type peer2 struct {
	ID     string
	PingMS float64
	Hops   float64
}

func TestExample2D(t *testing.T) {
	peers := []peer2{
		{ID: "A", PingMS: 22, Hops: 3},
		{ID: "B", PingMS: 34, Hops: 2},
		{ID: "C", PingMS: 15, Hops: 4},
		{ID: "D", PingMS: 55, Hops: 1},
		{ID: "E", PingMS: 18, Hops: 2},
	}
	weights := [2]float64{1.0, 1.0}
	invert := [2]bool{false, false}
	pts, err := poindexter.Build2D(
		peers,
		func(p peer2) string { return p.ID },
		func(p peer2) float64 { return p.PingMS },
		func(p peer2) float64 { return p.Hops },
		weights, invert,
	)
	if err != nil {
		t.Fatalf("Build2D err: %v", err)
	}
	tr, err := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.ManhattanDistance{}))
	if err != nil {
		t.Fatalf("NewKDTree err: %v", err)
	}
	best, d, ok := tr.Nearest([]float64{0, 0.3})
	if !ok {
		t.Fatalf("no nearest")
	}
	if best.ID == "" {
		t.Fatalf("unexpected empty ID")
	}
	if d < 0 {
		t.Fatalf("negative distance: %v", d)
	}
}
