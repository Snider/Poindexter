package main

import (
	"fmt"
	po "github.com/Snider/Poindexter"
)

// BuildPingHop2D wraps poindexter.Build2D to construct 2D points from (ping_ms, hop_count).
func BuildPingHop2D[T any](
	items []T,
	id func(T) string,
	ping func(T) float64,
	hops func(T) float64,
	weights [2]float64,
	invert [2]bool,
) ([]po.KDPoint[T], error) {
	return po.Build2D(items, id, ping, hops, weights, invert)
}

// BuildPingHopGeo3D wraps poindexter.Build3D for (ping_ms, hop_count, geo_km).
func BuildPingHopGeo3D[T any](
	items []T,
	id func(T) string,
	ping func(T) float64,
	hops func(T) float64,
	geoKM func(T) float64,
	weights [3]float64,
	invert [3]bool,
) ([]po.KDPoint[T], error) {
	return po.Build3D(items, id, ping, hops, geoKM, weights, invert)
}

// BuildPingHopGeoScore4D wraps poindexter.Build4D for (ping_ms, hop_count, geo_km, score).
// Typical usage sets invert for score=true so higher score => lower cost.
func BuildPingHopGeoScore4D[T any](
	items []T,
	id func(T) string,
	ping func(T) float64,
	hops func(T) float64,
	geoKM func(T) float64,
	score func(T) float64,
	weights [4]float64,
	invert [4]bool,
) ([]po.KDPoint[T], error) {
	return po.Build4D(items, id, ping, hops, geoKM, score, weights, invert)
}

// Demo program that builds a small tree using the 2D helper and performs a query.
func main() {
	type Peer struct {
		ID           string
		PingMS, Hops float64
	}
	peers := []Peer{{"A", 20, 1}, {"B", 50, 2}, {"C", 10, 3}}

	pts, err := BuildPingHop2D(peers,
		func(p Peer) string { return p.ID },
		func(p Peer) float64 { return p.PingMS },
		func(p Peer) float64 { return p.Hops },
		[2]float64{1.0, 0.7},
		[2]bool{false, false},
	)
	if err != nil {
		panic(err)
	}
	kdt, _ := po.NewKDTree(pts, po.WithMetric(po.EuclideanDistance{}))
	best, dist, _ := kdt.Nearest([]float64{0, 0})
	fmt.Println(best.ID, dist)
}
