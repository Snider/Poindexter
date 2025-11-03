package main

import (
	"fmt"
	poindexter "github.com/Snider/Poindexter"
)

type Peer4 struct {
	ID     string
	PingMS float64
	Hops   float64
	GeoKM  float64
	Score  float64
}

func main() {
	peers := []Peer4{
		{ID: "A", PingMS: 22, Hops: 3, GeoKM: 1200, Score: 0.86},
		{ID: "B", PingMS: 34, Hops: 2, GeoKM: 800, Score: 0.91},
		{ID: "C", PingMS: 15, Hops: 4, GeoKM: 4500, Score: 0.70},
		{ID: "D", PingMS: 55, Hops: 1, GeoKM: 300, Score: 0.95},
		{ID: "E", PingMS: 18, Hops: 2, GeoKM: 2200, Score: 0.80},
	}
	weights := [4]float64{1.0, 0.7, 0.2, 1.2}
	invert := [4]bool{false, false, false, true}
	pts, _ := poindexter.Build4D(
		peers,
		func(p Peer4) string { return p.ID },
		func(p Peer4) float64 { return p.PingMS },
		func(p Peer4) float64 { return p.Hops },
		func(p Peer4) float64 { return p.GeoKM },
		func(p Peer4) float64 { return p.Score },
		weights, invert,
	)
	tr, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
	best, _, _ := tr.Nearest([]float64{0, weights[1] * 0.2, weights[2] * 0.3, 0})
	fmt.Println("4D best:", best.ID)
}
