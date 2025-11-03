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
	fmt.Printf("dim=%d len=%d", tr.Dim(), tr.Len())
	// Output: dim=2 len=3
}

func ExampleKDTree_Nearest() {
	pts := []poindexter.KDPoint[int]{
		{ID: "x", Coords: []float64{0, 0}, Value: 1},
		{ID: "y", Coords: []float64{2, 0}, Value: 2},
	}
	tr, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
	p, d, ok := tr.Nearest([]float64{1, 0})
	fmt.Printf("ok=%v id=%s d=%.1f", ok, p.ID, d)
	// Output: ok=true id=y d=1.0
}

func ExampleKDTree_KNearest() {
	pts := []poindexter.KDPoint[int]{
		{ID: "a", Coords: []float64{0}, Value: 0},
		{ID: "b", Coords: []float64{1}, Value: 0},
		{ID: "c", Coords: []float64{2}, Value: 0},
	}
	tr, _ := poindexter.NewKDTree(pts)
	ns, ds := tr.KNearest([]float64{0.6}, 2)
	fmt.Printf("%s %.1f | %s %.1f", ns[0].ID, ds[0], ns[1].ID, ds[1])
	// Output: b 0.4 | a 0.6
}

func ExampleKDTree_Radius() {
	pts := []poindexter.KDPoint[int]{
		{ID: "a", Coords: []float64{0}, Value: 0},
		{ID: "b", Coords: []float64{1}, Value: 0},
		{ID: "c", Coords: []float64{2}, Value: 0},
	}
	tr, _ := poindexter.NewKDTree(pts)
	within, _ := tr.Radius([]float64{0}, 1.0)
	fmt.Printf("%d %s %s", len(within), within[0].ID, within[1].ID)
	// Output: 2 a b
}

func ExampleKDTree_InsertDeleteByID() {
	pts := []poindexter.KDPoint[string]{
		{ID: "A", Coords: []float64{0}, Value: "a"},
	}
	tr, _ := poindexter.NewKDTree(pts)
	tr.Insert(poindexter.KDPoint[string]{ID: "Z", Coords: []float64{0.1}, Value: "z"})
	p, _, _ := tr.Nearest([]float64{0.09})
	fmt.Println(p.ID)
	tr.DeleteByID("Z")
	p2, _, _ := tr.Nearest([]float64{0.09})
	fmt.Println(p2.ID)
	// Output:
	// Z
	// A
}

func ExampleBuild3D() {
	type rec struct{ x, y, z float64 }
	items := []rec{{0, 0, 0}, {1, 1, 1}}
	weights := [3]float64{1, 1, 1}
	invert := [3]bool{false, false, false}
	pts, _ := poindexter.Build3D(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.x },
		func(r rec) float64 { return r.y },
		func(r rec) float64 { return r.z },
		weights, invert,
	)
	tr, _ := poindexter.NewKDTree(pts)
	fmt.Println(tr.Dim())
	// Output: 3
}

func ExampleBuild4D() {
	type rec struct{ a, b, c, d float64 }
	items := []rec{{0, 0, 0, 0}, {1, 1, 1, 1}}
	weights := [4]float64{1, 1, 1, 1}
	invert := [4]bool{false, false, false, false}
	pts, _ := poindexter.Build4D(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		func(r rec) float64 { return r.c },
		func(r rec) float64 { return r.d },
		weights, invert,
	)
	tr, _ := poindexter.NewKDTree(pts)
	fmt.Println(tr.Dim())
	// Output: 4
}

func ExampleBuild2DWithStats() {
	type rec struct{ ping, hops float64 }
	items := []rec{{20, 3}, {30, 2}, {15, 4}}
	weights := [2]float64{1.0, 1.0}
	invert := [2]bool{false, false}
	stats := poindexter.ComputeNormStats2D(items,
		func(r rec) float64 { return r.ping },
		func(r rec) float64 { return r.hops },
	)
	pts, _ := poindexter.Build2DWithStats(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.ping },
		func(r rec) float64 { return r.hops },
		weights, invert, stats,
	)
	tr, _ := poindexter.NewKDTree(pts)
	fmt.Printf("dim=%d len=%d", tr.Dim(), tr.Len())
	// Output: dim=2 len=3
}

func ExampleBuild4DWithStats() {
	type rec struct{ a, b, c, d float64 }
	items := []rec{{0, 0, 0, 0}, {1, 1, 1, 1}}
	weights := [4]float64{1, 1, 1, 1}
	invert := [4]bool{false, false, false, false}
	stats := poindexter.ComputeNormStats4D(items,
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		func(r rec) float64 { return r.c },
		func(r rec) float64 { return r.d },
	)
	pts, _ := poindexter.Build4DWithStats(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		func(r rec) float64 { return r.c },
		func(r rec) float64 { return r.d },
		weights, invert, stats,
	)
	tr, _ := poindexter.NewKDTree(pts)
	fmt.Println(tr.Dim())
	// Output: 4
}

func ExampleNewKDTreeFromDim_Insert() {
	// Construct an empty 2D tree, insert a point, then query.
	tr, _ := poindexter.NewKDTreeFromDim[string](2)
	tr.Insert(poindexter.KDPoint[string]{ID: "A", Coords: []float64{0.1, 0.2}, Value: "alpha"})
	p, _, ok := tr.Nearest([]float64{0, 0})
	fmt.Printf("ok=%v id=%s dim=%d len=%d", ok, p.ID, tr.Dim(), tr.Len())
	// Output: ok=true id=A dim=2 len=1
}

func ExampleKDTree_TiesBehavior() {
	// Two points equidistant from the query; tie ordering is arbitrary,
	// but distances are equal.
	pts := []poindexter.KDPoint[int]{
		{ID: "L", Coords: []float64{-1}},
		{ID: "R", Coords: []float64{+1}},
	}
	tr, _ := poindexter.NewKDTree(pts)
	ns, ds := tr.KNearest([]float64{0}, 2)
	_ = ns // neighbor order is unspecified
	fmt.Printf("equal=%.1f==%.1f? %v", ds[0], ds[1], ds[0] == ds[1])
	// Output: equal=1.0==1.0? true
}

func ExampleKDTree_Radius_none() {
	// Radius query that yields no matches.
	pts := []poindexter.KDPoint[int]{
		{ID: "a", Coords: []float64{10}},
		{ID: "b", Coords: []float64{20}},
	}
	tr, _ := poindexter.NewKDTree(pts)
	within, _ := tr.Radius([]float64{0}, 5)
	fmt.Println(len(within))
	// Output: 0
}
