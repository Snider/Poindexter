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

We’ll use a small, made‑up set of DHT peers in each runnable example below. Each example declares its own `Peer` type and dataset so you can copy‑paste and run independently.

## Normalization and weights

To make heterogeneous units comparable (ms, hops, km, score), use the library helpers which:
- Min‑max normalize each axis to [0,1] over your provided dataset
- Optionally invert axes where “higher is better” so they become “lower cost”
- Apply per‑axis weights so you can emphasize what matters

Build 4‑D points and query them with helpers (full program):

```go
package main

import (
    "fmt"
    poindexter "github.com/Snider/Poindexter"
)

type Peer struct {
    ID        string
    PingMS    float64
    Hops      float64
    GeoKM     float64
    Score     float64
}

var peers = []Peer{
    {ID: "A", PingMS: 22, Hops: 3, GeoKM: 1200, Score: 0.86},
    {ID: "B", PingMS: 34, Hops: 2, GeoKM: 800,  Score: 0.91},
    {ID: "C", PingMS: 15, Hops: 4, GeoKM: 4500, Score: 0.70},
    {ID: "D", PingMS: 55, Hops: 1, GeoKM: 300,  Score: 0.95},
    {ID: "E", PingMS: 18, Hops: 2, GeoKM: 2200, Score: 0.80},
}

func main() {
    // Build 4‑D KDTree using Euclidean (L2)
    weights4 := [4]float64{1.0, 0.7, 0.2, 1.2}
    invert4 := [4]bool{false, false, false, true} // invert score (higher is better)
    pts, err := poindexter.Build4D(
        peers,
        func(p Peer) string { return p.ID },
        func(p Peer) float64 { return p.PingMS },
        func(p Peer) float64 { return p.Hops },
        func(p Peer) float64 { return p.GeoKM },
        func(p Peer) float64 { return p.Score },
        weights4, invert4,
    )
    if err != nil { panic(err) }
    tree, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))

    // Query target preferences (construct a query in normalized/weighted space)
    // Example: seek very low ping, low hops, moderate geo, high score (low score_cost)
    query := []float64{weights4[0]*0.0, weights4[1]*0.2, weights4[2]*0.3, weights4[3]*0.0}

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

Sometimes you want a strict trade‑off between just latency and path length. Build 2‑D points using helpers:

```go
package main

import (
    "fmt"
    poindexter "github.com/Snider/Poindexter"
)

type Peer struct {
    ID     string
    PingMS float64
    Hops   float64
}

var peers = []Peer{
    {ID: "A", PingMS: 22, Hops: 3},
    {ID: "B", PingMS: 34, Hops: 2},
    {ID: "C", PingMS: 15, Hops: 4},
    {ID: "D", PingMS: 55, Hops: 1},
    {ID: "E", PingMS: 18, Hops: 2},
}

func main() {
    weights2 := [2]float64{1.0, 1.0}
    invert2  := [2]bool{false, false}

    pts2, err := poindexter.Build2D(
        peers,
        func(p Peer) string { return p.ID },     // id
        func(p Peer) float64 { return p.PingMS },// f1: ping
        func(p Peer) float64 { return p.Hops },  // f2: hops
        weights2, invert2,
    )
    if err != nil { panic(err) }

    tree2, _ := poindexter.NewKDTree(pts2, poindexter.WithMetric(poindexter.ManhattanDistance{})) // L1 favors axis‑aligned tradeoffs
    // Prefer very low ping, modest hops
    query2 := []float64{weights2[0]*0.0, weights2[1]*0.3}
    best2, _, _ := tree2.Nearest(query2)
    fmt.Println("2D best (ping+hop):", best2.ID)
}
```

## 3‑D: Ping + Hop + Geo

Add geography to discourage far peers when latency is similar. Use the 3‑D helper:

```go
package main

import (
    "fmt"
    poindexter "github.com/Snider/Poindexter"
)

type Peer struct {
    ID     string
    PingMS float64
    Hops   float64
    GeoKM  float64
}

var peers = []Peer{
    {ID: "A", PingMS: 22, Hops: 3, GeoKM: 1200},
    {ID: "B", PingMS: 34, Hops: 2, GeoKM: 800},
    {ID: "C", PingMS: 15, Hops: 4, GeoKM: 4500},
    {ID: "D", PingMS: 55, Hops: 1, GeoKM: 300},
    {ID: "E", PingMS: 18, Hops: 2, GeoKM: 2200},
}

func main() {
    weights3 := [3]float64{1.0, 0.7, 0.3}
    invert3  := [3]bool{false, false, false}

    pts3, err := poindexter.Build3D(
        peers,
        func(p Peer) string { return p.ID },
        func(p Peer) float64 { return p.PingMS },
        func(p Peer) float64 { return p.Hops },
        func(p Peer) float64 { return p.GeoKM },
        weights3, invert3,
    )
    if err != nil { panic(err) }

    tree3, _ := poindexter.NewKDTree(pts3, poindexter.WithMetric(poindexter.EuclideanDistance{}))
    // Prefer low ping/hop, modest geo
    query3 := []float64{weights3[0]*0.0, weights3[1]*0.2, weights3[2]*0.4}
    top3, _, _ := tree3.Nearest(query3)
    fmt.Println("3D best (ping+hop+geo):", top3.ID)
}
```

## Dynamic updates

Your routing table changes constantly. Insert/remove peers. For consistent normalization, compute and reuse your min/max stats (preferred) or rebuild points when the candidate set changes.

Tip: Use the WithStats helpers to reuse normalization across updates:

```go
// Compute once over your baseline
stats := poindexter.ComputeNormStats2D(peers,
    func(p Peer) float64 { return p.PingMS },
    func(p Peer) float64 { return p.Hops },
)

// Build now or later using the same stats
ts, _ := poindexter.Build2DWithStats(
    peers,
    func(p Peer) string { return p.ID },
    func(p Peer) float64 { return p.PingMS },
    func(p Peer) float64 { return p.Hops },
    [2]float64{1,1}, [2]bool{false,false}, stats,
)
```

```go
package main

import (
    "fmt"
    poindexter "github.com/Snider/Poindexter"
)

type Peer struct {
    ID     string
    PingMS float64
    Hops   float64
}

var peers = []Peer{
    {ID: "A", PingMS: 22, Hops: 3},
    {ID: "B", PingMS: 34, Hops: 2},
    {ID: "C", PingMS: 15, Hops: 4},
}

func main() {
    // Initial 2‑D build (ping + hops)
    weights2 := [2]float64{1.0, 1.0}
    invert2  := [2]bool{false, false}
    pts, _ := poindexter.Build2D(
        peers,
        func(p Peer) string { return p.ID },
        func(p Peer) float64 { return p.PingMS },
        func(p Peer) float64 { return p.Hops },
        weights2, invert2,
    )
    tree, _ := poindexter.NewKDTree(pts)

    // Insert a new peer: rebuild its point using the same helper.
    newPeer := Peer{ID: "Z", PingMS: 12, Hops: 2}
    addPts, _ := poindexter.Build2D(
        []Peer{newPeer},
        func(p Peer) string { return p.ID },
        func(p Peer) float64 { return p.PingMS },
        func(p Peer) float64 { return p.Hops },
        weights2, invert2,
    )
    _ = tree.Insert(addPts[0])

    // Verify nearest now prefers Z for low ping target
    best, _, _ := tree.Nearest([]float64{0, 0})
    fmt.Println("Best after insert:", best.ID)

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
