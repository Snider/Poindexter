package main

import (
	"fmt"
	poindexter "github.com/Snider/Poindexter"
)

type Peer struct {
	Addr string
	Ping int
}

func main() {
	// Toy DHT routing table
	table := []Peer{
		{Addr: "peer1.example:4001", Ping: 74},
		{Addr: "peer2.example:4001", Ping: 52},
		{Addr: "peer3.example:4001", Ping: 110},
		{Addr: "peer4.example:4001", Ping: 35},
		{Addr: "peer5.example:4001", Ping: 60},
		{Addr: "peer6.example:4001", Ping: 44},
	}
	pts := make([]poindexter.KDPoint[Peer], 0, len(table))
	for i, p := range table {
		pts = append(pts, poindexter.KDPoint[Peer]{
			ID:     fmt.Sprintf("peer-%d", i+1),
			Coords: []float64{float64(p.Ping)},
			Value:  p,
		})
	}
	kdt, err := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
	if err != nil {
		panic(err)
	}
	best, d, ok := kdt.Nearest([]float64{0})
	if !ok {
		fmt.Println("no peers found")
		return
	}
	fmt.Printf("Best peer: %s (ping=%d ms), distance=%.0f\n", best.Value.Addr, best.Value.Ping, d)
}
