# KDTree: Multi‑Dimensional Search (DHT peers)

This example extends the single‑dimension "best ping" demo to a realistic multi‑dimensional selection:

- ping_ms (lower is better)
- hop_count (lower is better)
- geo_distance_km (lower is better)
- score (higher is better — e.g., capacity/reputation)

We will:
- Build 4‑D points over these features
- Run `Nearest`, `KNearest`, and `Radius` queries
- Show subsets: ping+hop (2‑D) and ping+hop+geo (3‑D)
- Demonstrate weighting/normalization to balance disparate units

> Tip: KDTree distances are geometric. Mixing units (ms, hops, km, arbitrary score) requires scaling so that each axis contributes proportionally to your decision policy.

## Dataset

```go
package main

import (
    "fmt"
    poindexter "github.com/Snider/Poindexter"
)

type Peer struct {
    ID        string
    PingMS    float64 // milliseconds
    Hops      float64 // hop count
    GeoKM     float64 // crow‑flight distance in kilometers
    Score     float64 // [0..1] trust/rep/capacity score (higher is better)
}

var peers = []Peer{
    {ID: "A", PingMS: 22, Hops: 3, GeoKM: 1200, Score: 0.86},
    {ID: "B", PingMS: 34, Hops: 2, GeoKM: 800,  Score: 0.91},
    {ID: "C", PingMS: 15, Hops: 4, GeoKM: 4500, Score: 0.70},
    {ID: "D", PingMS: 55, Hops: 1, GeoKM: 300,  Score: 0.95},
    {ID: "E", PingMS: 18, Hops: 2, GeoKM: 2200, Score: 0.80},
}
```

## Normalization and weights

We scale raw features to comparable magnitudes and flip `Score` so lower is better. For demo simplicity we will:
- Min‑max normalize each axis to [0,1] over the current candidate set
- Convert `Score` to a cost: `score_cost = 1 - score`
- Apply weights to emphasize certain axes

Helper functions:

```go
// minMax returns (min, max) of a slice.
func minMax(xs []float64) (float64, float64) {
    if len(xs) == 0 { return 0, 0 }
    mn, mx := xs[0], xs[0]
    for _, v := range xs[1:] {
        if v < mn { mn = v }
        if v > mx { mx = v }
    }
    return mn, mx
}

// scale01 maps v from [min,max] to [0,1]. If min==max, returns 0.
func scale01(v, min, max float64) float64 {
    if max == min { return 0 }
    return (v - min) / (max - min)
}
```

Build 4‑D points:

```go
// Weights to balance axes (tune to taste)
var wPing, wHop, wGeo, wScore = 1.0, 0.7, 0.2, 1.2

func build4D(peers []Peer) ([]poindexter.KDPoint[Peer], error) {
    pings := make([]float64, len(peers))
    hops  := make([]float64, len(peers))
    geos  := make([]float64, len(peers))
    scores:= make([]float64, len(peers))
    for i, p := range peers {
        pings[i], hops[i], geos[i], scores[i] = p.PingMS, p.Hops, p.GeoKM, p.Score
    }
    pMin, pMax := minMax(pings)
    hMin, hMax := minMax(hops)
    gMin, gMax := minMax(geos)
    sMin, sMax := minMax(scores)

    pts := make([]poindexter.KDPoint[Peer], len(peers))
    for i, p := range peers {
        pingN  := scale01(p.PingMS, pMin, pMax)
        hopN   := scale01(p.Hops,   hMin, hMax)
        geoN   := scale01(p.GeoKM,  gMin, gMax)
        scoreC := 1 - scale01(p.Score, sMin, sMax) // lower is better

        pts[i] = poindexter.KDPoint[Peer]{
            ID:    p.ID,
            Value: p,
            Coords: []float64{
                wPing*pingN,
                wHop*hopN,
                wGeo*geoN,
                wScore*scoreC,
            },
        }
    }
    return pts, nil
}
```

## 4‑D KDTree: Nearest, k‑NN, Radius

```go
func main() {
    // Build 4‑D KDTree using Euclidean (L2)
    pts, _ := build4D(peers)
    tree, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))

    // Query target preferences (you may construct a query in normalized/weighted space)
    // Example: seek very low ping, low hops, moderate geo, high score (low score_cost)
    query := []float64{wPing*0.0, wHop*0.2, wGeo*0.3, wScore*0.0}

    // 1‑NN
    best, dist, ok := tree.Nearest(query)
    if ok {
        fmt.Printf("Best peer: %s (dist=%.4f)\n", best.ID, dist)
    }

    // k‑NN (top 3)
    neigh, dists := tree.KNearest(query, 3)
    for i := range neigh {
        fmt.Printf("%d) %s dist=%.4f\n", i+1, neigh[i].ID, dists[i])
    }

    // Radius query
    within, wd := tree.Radius(query, 0.35)
    fmt.Printf("Within radius 0.35: ")
    for i := range within {
        fmt.Printf("%s(%.3f) ", within[i].ID, wd[i])
    }
    fmt.Println()
}
```

## 2‑D: Ping + Hop

Sometimes you want a strict trade‑off between just latency and path length. Build 2‑D points (reuse normalization):

```go
var wPing2, wHop2 = 1.0, 1.0

func build2D_pingHop(peers []Peer) []poindexter.KDPoint[Peer] {
    pings := make([]float64, len(peers))
    hops  := make([]float64, len(peers))
    for i, p := range peers { pings[i], hops[i] = p.PingMS, p.Hops }
    pMin, pMax := minMax(pings)
    hMin, hMax := minMax(hops)

    pts := make([]poindexter.KDPoint[Peer], len(peers))
    for i, p := range peers {
        pingN := scale01(p.PingMS, pMin, pMax)
        hopN  := scale01(p.Hops,   hMin, hMax)
        pts[i] = poindexter.KDPoint[Peer]{
            ID:    p.ID,
            Value: p,
            Coords: []float64{ wPing2*pingN, wHop2*hopN },
        }
    }
    return pts
}

func demo2D() {
    pts := build2D_pingHop(peers)
    tree, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.ManhattanDistance{})) // L1 favors axis‑aligned tradeoffs
    // Prefer very low ping, modest hops
    query := []float64{wPing2*0.0, wHop2*0.3}
    best, _, _ := tree.Nearest(query)
    fmt.Println("2D best (ping+hop):", best.ID)
}
```

## 3‑D: Ping + Hop + Geo

Add geography to discourage far peers when latency is similar:

```go
var wPing3, wHop3, wGeo3 = 1.0, 0.7, 0.3

func build3D_pingHopGeo(peers []Peer) []poindexter.KDPoint[Peer] {
    pings := make([]float64, len(peers))
    hops  := make([]float64, len(peers))
    geos  := make([]float64, len(peers))
    for i, p := range peers { pings[i], hops[i], geos[i] = p.PingMS, p.Hops, p.GeoKM }
    pMin, pMax := minMax(pings)
    hMin, hMax := minMax(hops)
    gMin, gMax := minMax(geos)

    pts := make([]poindexter.KDPoint[Peer], len(peers))
    for i, p := range peers {
        pingN := scale01(p.PingMS, pMin, pMax)
        hopN  := scale01(p.Hops,   hMin, hMax)
        geoN  := scale01(p.GeoKM,  gMin, gMax)
        pts[i] = poindexter.KDPoint[Peer]{
            ID:    p.ID,
            Value: p,
            Coords: []float64{ wPing3*pingN, wHop3*hopN, wGeo3*geoN },
        }
    }
    return pts
}

func demo3D() {
    pts := build3D_pingHopGeo(peers)
    tree, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
    // Prefer low ping/hop, modest geo
    query := []float64{wPing3*0.0, wHop3*0.2, wGeo3*0.4}
    top, _, _ := tree.Nearest(query)
    fmt.Println("3D best (ping+hop+geo):", top.ID)
}
```

## Dynamic updates

Your routing table changes constantly. Insert/remove peers without rebuilding:

```go
func updatesExample() {
    pts := build2D_pingHop(peers)
    tree, _ := poindexter.NewKDTree(pts)

    // Insert a new peer
    newPeer := Peer{ID: "Z", PingMS: 12, Hops: 2, GeoKM: 900, Score: 0.88}
    // Build consistent 2D point for the new peer. In a real system retain normalization mins/maxes.
    ptsZ := build2D_pingHop([]Peer{newPeer})
    _ = tree.Insert(ptsZ[0])

    // Delete by ID when peer goes offline
    _ = tree.DeleteByID("Z")
}
```

## Choosing a metric

- Euclidean (L2): smooth trade‑offs across axes; good default for blended preferences
- Manhattan (L1): emphasizes per‑axis absolute differences; useful when each unit of ping/hop matters equally
- Chebyshev (L∞): min‑max style; dominated by the worst axis (e.g., reject any peer with too many hops regardless of ping)

## Notes on production use

- Keep and reuse normalization parameters (min/max or mean/std) rather than recomputing per query to avoid drift.
- Consider capping outliers (e.g., clamp geo distances > 5000 km).
- For large N (≥ 1e5) and low dims (≤ 8), consider swapping the internal engine to `gonum.org/v1/gonum/spatial/kdtree` behind the same API for faster queries.
