package poindexter

import (
	"fmt"
	"testing"
)

func TestBuild2D_NormalizationAndInversion(t *testing.T) {
	type rec struct{ a, b float64 }
	items := []rec{{a: 0, b: 100}, {a: 10, b: 300}}
	// f1 over [0,10], f2 over [100,300]
	pts, err := Build2D(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		[2]float64{2.0, 0.5},
		[2]bool{true, false}, // invert first axis, not second
	)
	if err != nil {
		t.Fatalf("Build2D err: %v", err)
	}
	if len(pts) != 2 {
		t.Fatalf("expected 2 points, got %d", len(pts))
	}
	// item0: a=0 -> n1=0 -> invert -> 1 -> *2 = 2; b=100 -> n2=0 -> *0.5 = 0
	if got := fmt.Sprintf("%.1f,%.1f", pts[0].Coords[0], pts[0].Coords[1]); got != "2.0,0.0" {
		t.Fatalf("coords[0] = %s, want 2.0,0.0", got)
	}
	// item1: a=10 -> n1=1 -> invert -> 0 -> *2 = 0; b=300 -> n2=1 -> *0.5=0.5
	if got := fmt.Sprintf("%.1f,%.1f", pts[1].Coords[0], pts[1].Coords[1]); got != "0.0,0.5" {
		t.Fatalf("coords[1] = %s, want 0.0,0.5", got)
	}
}

func TestBuild3D_AllEqualSafe(t *testing.T) {
	type rec struct{ x, y, z float64 }
	items := []rec{{1, 1, 1}, {1, 1, 1}}
	pts, err := Build3D(items,
		func(r rec) string { return "id" },
		func(r rec) float64 { return r.x },
		func(r rec) float64 { return r.y },
		func(r rec) float64 { return r.z },
		[3]float64{1, 1, 1},
		[3]bool{false, false, false},
	)
	if err != nil {
		t.Fatalf("Build3D err: %v", err)
	}
	if len(pts) != 2 {
		t.Fatalf("len = %d", len(pts))
	}
	for i := range pts {
		if len(pts[i].Coords) != 3 {
			t.Fatalf("dim = %d", len(pts[i].Coords))
		}
		for _, c := range pts[i].Coords {
			if c != 0 {
				t.Fatalf("expected 0 when min==max, got %v", c)
			}
		}
	}
}

// Example-style end-to-end sanity on 4D using the documented Peer data
func TestBuild4D_EndToEnd_Example(t *testing.T) {
	type Peer struct {
		ID     string
		PingMS float64
		Hops   float64
		GeoKM  float64
		Score  float64
	}
	peers := []Peer{
		{ID: "A", PingMS: 22, Hops: 3, GeoKM: 1200, Score: 0.86},
		{ID: "B", PingMS: 34, Hops: 2, GeoKM: 800, Score: 0.91},
		{ID: "C", PingMS: 15, Hops: 4, GeoKM: 4500, Score: 0.70},
		{ID: "D", PingMS: 55, Hops: 1, GeoKM: 300, Score: 0.95},
		{ID: "E", PingMS: 18, Hops: 2, GeoKM: 2200, Score: 0.80},
	}
	weights := [4]float64{1.0, 0.7, 0.2, 1.2}
	invert := [4]bool{false, false, false, true} // flip score so higher score -> lower cost
	pts, err := Build4D(peers,
		func(p Peer) string { return p.ID },
		func(p Peer) float64 { return p.PingMS },
		func(p Peer) float64 { return p.Hops },
		func(p Peer) float64 { return p.GeoKM },
		func(p Peer) float64 { return p.Score },
		weights, invert,
	)
	if err != nil {
		t.Fatalf("Build4D err: %v", err)
	}
	if len(pts) != len(peers) {
		t.Fatalf("len pts=%d", len(pts))
	}
	// Build KDTree and query near origin in normalized/weighted space (prefer minima on all axes)
	tree, err := NewKDTree(pts, WithMetric(EuclideanDistance{}))
	if err != nil {
		t.Fatalf("NewKDTree err: %v", err)
	}
	if tree.Dim() != 4 {
		t.Fatalf("dim=%d", tree.Dim())
	}
	best, _, ok := tree.Nearest([]float64{0, 0, 0, 0})
	if !ok {
		t.Fatalf("no nearest")
	}
	// With these weights and inversions, peer B emerges as closest in this setup.
	if best.ID != "B" {
		t.Fatalf("expected best B, got %s", best.ID)
	}
}

func TestComputeNormStatsAndWithStats_Parity2D(t *testing.T) {
	type rec struct{ a, b float64 }
	items := []rec{{0, 10}, {5, 20}, {10, 30}}
	weights := [2]float64{1, 2}
	invert := [2]bool{false, true}
	// Build using automatic stats
	autoPts, err := Build2D(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		weights, invert,
	)
	if err != nil {
		t.Fatalf("auto build err: %v", err)
	}
	// Compute stats and build with stats
	stats := ComputeNormStats2D(items,
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
	)
	withPts, err := Build2DWithStats(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		weights, invert, stats,
	)
	if err != nil {
		t.Fatalf("with-stats build err: %v", err)
	}
	if len(withPts) != len(autoPts) {
		t.Fatalf("len mismatch")
	}
	for i := range withPts {
		if len(withPts[i].Coords) != 2 {
			t.Fatalf("dim mismatch")
		}
		if withPts[i].Coords[0] != autoPts[i].Coords[0] || withPts[i].Coords[1] != autoPts[i].Coords[1] {
			t.Fatalf("coords mismatch at %d: %v vs %v", i, withPts[i].Coords, autoPts[i].Coords)
		}
	}
}

func TestBuild3DWithStats_MinEqualsMax_Safe(t *testing.T) {
	type rec struct{ x, y, z float64 }
	items := []rec{{1, 2, 3}, {1, 5, 3}, {1, 9, 3}}
	weights := [3]float64{1, 1, 1}
	invert := [3]bool{false, false, false}
	// x and z min==max across items for x=1, z=3
	stats := NormStats{Stats: []AxisStats{{Min: 1, Max: 1}, {Min: 2, Max: 9}, {Min: 3, Max: 3}}}
	pts, err := Build3DWithStats(items,
		func(r rec) string { return "" },
		func(r rec) float64 { return r.x },
		func(r rec) float64 { return r.y },
		func(r rec) float64 { return r.z },
		weights, invert, stats,
	)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, p := range pts {
		if p.Coords[0] != 0 || p.Coords[2] != 0 {
			t.Fatalf("expected zero for min==max axes, got %v", p.Coords)
		}
	}
}

func TestBuild4DWithStats_DynamicUpdateExample(t *testing.T) {
	type Peer struct {
		ID                     string
		Ping, Hops, Geo, Score float64
	}
	base := []Peer{{"A", 20, 3, 1000, 0.8}, {"B", 30, 2, 800, 0.9}}
	weights := [4]float64{1, 1, 0.2, 1.2}
	invert := [4]bool{false, false, false, true}
	stats := ComputeNormStats4D(base,
		func(p Peer) float64 { return p.Ping },
		func(p Peer) float64 { return p.Hops },
		func(p Peer) float64 { return p.Geo },
		func(p Peer) float64 { return p.Score },
	)
	pts, err := Build4DWithStats(base,
		func(p Peer) string { return p.ID },
		func(p Peer) float64 { return p.Ping },
		func(p Peer) float64 { return p.Hops },
		func(p Peer) float64 { return p.Geo },
		func(p Peer) float64 { return p.Score },
		weights, invert, stats,
	)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	tr, err := NewKDTree(pts)
	if err != nil {
		t.Fatalf("kdt err: %v", err)
	}
	// add a new peer using same stats
	newPeer := Peer{"Z", 15, 2, 1200, 0.85}
	newPts, _ := Build4DWithStats([]Peer{newPeer},
		func(p Peer) string { return p.ID },
		func(p Peer) float64 { return p.Ping },
		func(p Peer) float64 { return p.Hops },
		func(p Peer) float64 { return p.Geo },
		func(p Peer) float64 { return p.Score },
		weights, invert, stats,
	)
	if !tr.Insert(newPts[0]) {
		t.Fatalf("insert failed")
	}
	if tr.Dim() != 4 {
		t.Fatalf("dim != 4")
	}
}

func TestComputeNormStats3D(t *testing.T) {
	type rec struct{ x, y, z float64 }
	items := []rec{{1, 10, 100}, {2, 20, 200}, {3, 30, 300}}
	stats := ComputeNormStats3D(items,
		func(r rec) float64 { return r.x },
		func(r rec) float64 { return r.y },
		func(r rec) float64 { return r.z },
	)
	expected := NormStats{
		Stats: []AxisStats{
			{Min: 1, Max: 3},
			{Min: 10, Max: 30},
			{Min: 100, Max: 300},
		},
	}
	if stats.Stats[0] != expected.Stats[0] || stats.Stats[1] != expected.Stats[1] || stats.Stats[2] != expected.Stats[2] {
		t.Fatalf("expected %v, got %v", expected, stats)
	}
}

func TestBuildND(t *testing.T) {
	type rec struct{ a, b, c float64 }
	items := []rec{{1, 2, 3}, {4, 5, 6}}
	extractors := []func(rec) float64{
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		func(r rec) float64 { return r.c },
	}
	weights := []float64{1, 1, 1}
	invert := []bool{false, false, false}
	pts, err := BuildND(items, func(r rec) string { return "" }, extractors, weights, invert)
	if err != nil {
		t.Fatal(err)
	}
	if len(pts) != 2 {
		t.Fatalf("expected 2 points, got %d", len(pts))
	}
	if len(pts[0].Coords) != 3 {
		t.Fatalf("expected 3 dimensions, got %d", len(pts[0].Coords))
	}
}

func TestBuildNDError(t *testing.T) {
	type rec struct{ a, b, c float64 }
	items := []rec{{1, 2, 3}, {4, 5, 6}}
	extractors := []func(rec) float64{
		func(r rec) float64 { return r.a },
		func(r rec) float64 { return r.b },
		func(r rec) float64 { return r.c },
	}
	weights := []float64{1, 1} // Mismatched length
	invert := []bool{false, false, false}
	_, err := BuildND(items, func(r rec) string { return "" }, extractors, weights, invert)
	if err == nil {
		t.Fatal("expected an error, but got nil")
	}
}
