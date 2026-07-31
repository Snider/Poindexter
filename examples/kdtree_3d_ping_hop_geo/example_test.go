package main

import (
	poindexter "github.com/Snider/Poindexter"
	"testing"
)

type peer3test struct {
	ID     string
	PingMS float64
	Hops   float64
	GeoKM  float64
}

func TestExample3D(t *testing.T) {
	peers := []peer3test{
		{ID: "A", PingMS: 22, Hops: 3, GeoKM: 1200},
		{ID: "B", PingMS: 34, Hops: 2, GeoKM: 800},
		{ID: "C", PingMS: 15, Hops: 4, GeoKM: 4500},
		{ID: "D", PingMS: 55, Hops: 1, GeoKM: 300},
		{ID: "E", PingMS: 18, Hops: 2, GeoKM: 2200},
	}
	weights := [3]float64{1.0, 0.7, 0.3}
	invert := [3]bool{false, false, false}
	pts, err := poindexter.Build3D(
		peers,
		func(p peer3test) string { return p.ID },
		func(p peer3test) float64 { return p.PingMS },
		func(p peer3test) float64 { return p.Hops },
		func(p peer3test) float64 { return p.GeoKM },
		weights, invert,
	)
	if err != nil {
		t.Fatalf("Build3D err: %v", err)
	}
	tr, err := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
	if err != nil {
		t.Fatalf("NewKDTree err: %v", err)
	}
	best, d, ok := tr.Nearest([]float64{0, weights[1] * 0.2, weights[2] * 0.4})
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
