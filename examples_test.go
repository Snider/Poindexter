package poindexter_test

import (
	"fmt"

	poindexter "github.com/Snider/Poindexter"
)

func ExampleNewKDTree() {
	pts := []poindexter.KDPoint[string]{
		{ID: "A", Coords: []float64{0, 0}, Value: "alpha"},
		{ID: "B", Coords: []float64{1, 0}, Value: "bravo"},
	}
	tr, _ := poindexter.NewKDTree(pts)
	p, _, _ := tr.Nearest([]float64{0.2, 0})
	fmt.Println(p.ID)
	// Output: A
}

func ExampleBuild2D() {
	type rec struct{ ping, hops float64 }
	items := []rec{{ping: 20, hops: 3}, {ping: 30, hops: 2}, {ping: 15, hops: 4}}
	weights := [2]float64{1.0, 1.0}
	invert := [2]bool{false, false}
	pts, _ := poindexter.Build2D(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.ping },
		func(r rec) float64 { return r.hops },
		weights, invert,
	)
	tr, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.ManhattanDistance{}))
	_, _, _ = tr.Nearest([]float64{0, 0})
	// Querying the origin (0,0) in normalized space tends to favor minima on each axis.
	fmt.Printf("dim=%d len=%d", tr.Dim(), tr.Len())
	// Output: dim=2 len=3
}
