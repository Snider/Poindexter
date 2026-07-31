package main

import (
	"fmt"
	poindexter "github.com/Snider/Poindexter"
)

type Peer3 struct {
	ID     string
	PingMS float64
	Hops   float64
	GeoKM  float64
}

func main() {
	peers := []Peer3{
		{ID: "A", PingMS: 22, Hops: 3, GeoKM: 1200},
		{ID: "B", PingMS: 34, Hops: 2, GeoKM: 800},
		{ID: "C", PingMS: 15, Hops: 4, GeoKM: 4500},
		{ID: "D", PingMS: 55, Hops: 1, GeoKM: 300},
		{ID: "E", PingMS: 18, Hops: 2, GeoKM: 2200},
	}
	weights := [3]float64{1.0, 0.7, 0.3}
	invert := [3]bool{false, false, false}
	pts, _ := poindexter.Build3D(
		peers,
		func(p Peer3) string { return p.ID },
		func(p Peer3) float64 { return p.PingMS },
		func(p Peer3) float64 { return p.Hops },
		func(p Peer3) float64 { return p.GeoKM },
		weights, invert,
	)
	tr, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
	best, _, _ := tr.Nearest([]float64{0, weights[1] * 0.2, weights[2] * 0.4})
	fmt.Println("3D best:", best.ID)
}
