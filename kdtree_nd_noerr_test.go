package poindexter

import (
	"fmt"
	"math"
	"testing"
)

// TestBuildNDNoErr_Parity checks that BuildNDNoErr matches BuildND on valid inputs.
func TestBuildNDNoErr_Parity(t *testing.T) {
	type rec struct {
		A, B, C float64
		ID      string
	}
	items := []rec{
		{A: 10, B: 100, C: 1, ID: "x"},
		{A: 20, B: 200, C: 2, ID: "y"},
		{A: 30, B: 300, C: 3, ID: "z"},
	}
	features := []func(rec) float64{
		func(r rec) float64 { return r.A },
		func(r rec) float64 { return r.B },
		func(r rec) float64 { return r.C },
	}
	weights := []float64{1, 0.5, 2}
	invert := []bool{false, true, false}
	idfn := func(r rec) string { return r.ID }

	ptsStrict, err := BuildND(items, idfn, features, weights, invert)
	if err != nil {
		t.Fatalf("BuildND returned error: %v", err)
	}
	ptsLoose := BuildNDNoErr(items, idfn, features, weights, invert)
	if len(ptsStrict) != len(ptsLoose) {
		t.Fatalf("length mismatch: strict %d loose %d", len(ptsStrict), len(ptsLoose))
	}
	for i := range ptsStrict {
		if ptsStrict[i].ID != ptsLoose[i].ID {
			t.Fatalf("ID mismatch at %d: %s vs %s", i, ptsStrict[i].ID, ptsLoose[i].ID)
		}
		if len(ptsStrict[i].Coords) != len(ptsLoose[i].Coords) {
			t.Fatalf("dim mismatch at %d: %d vs %d", i, len(ptsStrict[i].Coords), len(ptsLoose[i].Coords))
		}
		for d := range ptsStrict[i].Coords {
			if math.Abs(ptsStrict[i].Coords[d]-ptsLoose[i].Coords[d]) > 1e-12 {
				t.Fatalf("coord mismatch at %d dim %d: %v vs %v", i, d, ptsStrict[i].Coords[d], ptsLoose[i].Coords[d])
			}
		}
	}
}

// TestBuildNDNoErr_Lenient ensures the no-error builder is lenient and returns nil on bad lengths.
func TestBuildNDNoErr_Lenient(t *testing.T) {
	type rec struct{ A float64 }
	items := []rec{{A: 1}, {A: 2}}
	features := []func(rec) float64{func(r rec) float64 { return r.A }}
	weightsBad := []float64{} // wrong length
	invert := []bool{false}
	pts := BuildNDNoErr(items, func(r rec) string { return fmt.Sprint(r.A) }, features, weightsBad, invert)
	if pts != nil {
		t.Fatalf("expected nil result on bad weights length, got %v", pts)
	}
}
