package main

import (
	"fmt"
	poindexter "github.com/Snider/Poindexter"
)

type Peer2 struct {
	ID     string
	PingMS float64
	Hops   float64
}

func main() {
	peers := []Peer2{
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
		func(p Peer2) string { return p.ID },
		func(p Peer2) float64 { return p.PingMS },
		func(p Peer2) float64 { return p.Hops },
		weights, invert,
	)
	if err != nil {
		panic(fmt.Sprintf("Build2D failed: %v", err))
	}
	tr, err := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.ManhattanDistance{}))
	if err != nil {
		panic(fmt.Sprintf("NewKDTree failed: %v", err))
	}
	best, _, ok := tr.Nearest([]float64{0, 0.3})
	if !ok {
		panic("no nearest neighbour found")
	}
	fmt.Println("2D best:", best.ID)
}
