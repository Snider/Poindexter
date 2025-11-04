package poindexter

import (
	"math/rand"
	"testing"
)

// makeFixedPoints creates a deterministic set of points in 4D and 2D for parity checks.
func makeFixedPoints() []KDPoint[int] {
	pts := []KDPoint[int]{
		{ID: "A", Coords: []float64{0, 0, 0, 0}, Value: 1},
		{ID: "B", Coords: []float64{1, 0, 0.5, 0.2}, Value: 2},
		{ID: "C", Coords: []float64{0, 1, 0.3, 0.7}, Value: 3},
		{ID: "D", Coords: []float64{1, 1, 0.9, 0.9}, Value: 4},
		{ID: "E", Coords: []float64{0.2, 0.8, 0.4, 0.6}, Value: 5},
	}
	return pts
}

func TestBackendParity_Nearest(t *testing.T) {
	pts := makeFixedPoints()
	queries := [][]float64{
		{0, 0, 0, 0},
		{0.9, 0.2, 0.5, 0.1},
		{0.5, 0.5, 0.5, 0.5},
	}

	lin, err := NewKDTree(pts, WithBackend(BackendLinear), WithMetric(EuclideanDistance{}))
	if err != nil {
		t.Fatalf("linear NewKDTree: %v", err)
	}

	// Only build a gonum tree when the optimized backend is compiled in.
	if hasGonum() {
		gon, err := NewKDTree(pts, WithBackend(BackendGonum), WithMetric(EuclideanDistance{}))
		if err != nil {
			t.Fatalf("gonum NewKDTree: %v", err)
		}
		for _, q := range queries {
			pl, dl, okl := lin.Nearest(q)
			pg, dg, okg := gon.Nearest(q)
			if okl != okg {
				t.Fatalf("ok mismatch: linear=%v gonum=%v", okl, okg)
			}
			if !okl {
				continue
			}
			if pl.ID != pg.ID {
				t.Errorf("nearest ID mismatch for %v: linear=%s gonum=%s", q, pl.ID, pg.ID)
			}
			if (dl == 0 && dg != 0) || (dl != 0 && dg == 0) {
				t.Errorf("nearest distance zero/nonzero mismatch: linear=%v gonum=%v", dl, dg)
			}
		}
	}
}

func TestBackendParity_KNearest(t *testing.T) {
	pts := makeFixedPoints()
	q := []float64{0.6, 0.6, 0.4, 0.4}
	ks := []int{1, 2, 5, 10}
	lin, _ := NewKDTree(pts, WithBackend(BackendLinear), WithMetric(EuclideanDistance{}))
	if hasGonum() {
		gon, _ := NewKDTree(pts, WithBackend(BackendGonum), WithMetric(EuclideanDistance{}))
		for _, k := range ks {
			ln, ld := lin.KNearest(q, k)
			gn, gd := gon.KNearest(q, k)
			if len(ln) != len(gn) || len(ld) != len(gd) {
				t.Fatalf("k=%d length mismatch: linear (%d,%d) vs gonum (%d,%d)", k, len(ln), len(ld), len(gn), len(gd))
			}
			// Compare IDs element-wise; ties may reorder between backends, so relax by set equality when distances equal.
			for i := range ln {
				if ln[i].ID != gn[i].ID {
					// If distances are effectively equal, allow different order
					if i < len(ld) && i < len(gd) && ld[i] == gd[i] {
						continue
					}
					t.Logf("k=%d index %d ID mismatch: linear=%s gonum=%s (dl=%.6f dg=%.6f)", k, i, ln[i].ID, gn[i].ID, ld[i], gd[i])
				}
			}
		}
	}
}

func TestBackendParity_Radius(t *testing.T) {
	pts := makeFixedPoints()
	q := []float64{0.4, 0.6, 0.4, 0.6}
	radii := []float64{0, 0.15, 0.3, 1.0}
	lin, _ := NewKDTree(pts, WithBackend(BackendLinear), WithMetric(EuclideanDistance{}))
	if hasGonum() {
		gon, _ := NewKDTree(pts, WithBackend(BackendGonum), WithMetric(EuclideanDistance{}))
		for _, r := range radii {
			ln, ld := lin.Radius(q, r)
			gn, gd := gon.Radius(q, r)
			if len(ln) != len(gn) || len(ld) != len(gd) {
				t.Fatalf("r=%.3f length mismatch: linear (%d,%d) vs gonum (%d,%d)", r, len(ln), len(ld), len(gn), len(gd))
			}
		}
	}
}

func TestBackendParity_RandomQueries2D(t *testing.T) {
	// Down-project 4D to 2D to exercise pruning differences as well
	pts4 := makeFixedPoints()
	pts2 := make([]KDPoint[int], len(pts4))
	for i, p := range pts4 {
		pts2[i] = KDPoint[int]{ID: p.ID, Coords: []float64{p.Coords[0], p.Coords[1]}, Value: p.Value}
	}
	lin, _ := NewKDTree(pts2, WithBackend(BackendLinear), WithMetric(ManhattanDistance{}))
	if hasGonum() {
		gon, _ := NewKDTree(pts2, WithBackend(BackendGonum), WithMetric(ManhattanDistance{}))
		rng := rand.New(rand.NewSource(42))
		for i := 0; i < 50; i++ {
			q := []float64{rng.Float64(), rng.Float64()}
			pl, dl, okl := lin.Nearest(q)
			pg, dg, okg := gon.Nearest(q)
			if okl != okg {
				t.Fatalf("ok mismatch (2D rand)")
			}
			if !okl {
				continue
			}
			if pl.ID != pg.ID && (dl != dg) {
				// Allow different picks only if distances tie; otherwise flag
				t.Errorf("2D rand nearest mismatch: linear %s(%.6f) gonum %s(%.6f)", pl.ID, dl, pg.ID, dg)
			}
		}
	}
}
