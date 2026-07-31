package main

import (
	poindexter "github.com/Snider/Poindexter"
	"testing"
)

type peer4test struct {
	ID     string
	PingMS float64
	Hops   float64
	GeoKM  float64
	Score  float64
}

func TestExample4D(t *testing.T) {
	peers := []peer4test{
		{ID: "A", PingMS: 22, Hops: 3, GeoKM: 1200, Score: 0.86},
		{ID: "B", PingMS: 34, Hops: 2, GeoKM: 800, Score: 0.91},
		{ID: "C", PingMS: 15, Hops: 4, GeoKM: 4500, Score: 0.70},
		{ID: "D", PingMS: 55, Hops: 1, GeoKM: 300, Score: 0.95},
		{ID: "E", PingMS: 18, Hops: 2, GeoKM: 2200, Score: 0.80},
	}
	weights := [4]float64{1.0, 0.7, 0.2, 1.2}
	invert := [4]bool{false, false, false, true}
	pts, err := poindexter.Build4D(
		peers,
		func(p peer4test) string { return p.ID },
		func(p peer4test) float64 { return p.PingMS },
		func(p peer4test) float64 { return p.Hops },
		func(p peer4test) float64 { return p.GeoKM },
		func(p peer4test) float64 { return p.Score },
		weights, invert,
	)
	if err != nil {
		t.Fatalf("Build4D err: %v", err)
	}
	tr, err := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
	if err != nil {
		t.Fatalf("NewKDTree err: %v", err)
	}
	best, d, ok := tr.Nearest([]float64{0, weights[1] * 0.2, weights[2] * 0.3, 0})
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
