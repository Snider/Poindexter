package main

import (
	"fmt"
	poindexter "github.com/Snider/Poindexter"
	"testing"
)

type peer struct {
	Addr string
	Ping int
}

// TestExample1D ensures the 1D example logic runs and exercises KDTree paths.
func TestExample1D(t *testing.T) {
	// Same toy table as the example
	table := []peer{
		{Addr: "peer1.example:4001", Ping: 74},
		{Addr: "peer2.example:4001", Ping: 52},
		{Addr: "peer3.example:4001", Ping: 110},
		{Addr: "peer4.example:4001", Ping: 35},
		{Addr: "peer5.example:4001", Ping: 60},
		{Addr: "peer6.example:4001", Ping: 44},
	}
	pts := make([]poindexter.KDPoint[peer], 0, len(table))
	for i, p := range table {
		pts = append(pts, poindexter.KDPoint[peer]{
			ID:     fmt.Sprintf("peer-%d", i+1),
			Coords: []float64{float64(p.Ping)},
			Value:  p,
		})
	}
	kdt, err := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
	if err != nil {
		t.Fatalf("NewKDTree err: %v", err)
	}
	best, d, ok := kdt.Nearest([]float64{0})
	if !ok {
		t.Fatalf("no nearest")
	}
	// Expect the minimum ping (35ms)
	if best.Value.Ping != 35 {
		t.Fatalf("expected best ping 35ms, got %d", best.Value.Ping)
	}
	// Distance from [0] to [35] should be 35
	if d != 35 {
		t.Fatalf("expected distance 35, got %v", d)
	}
}
