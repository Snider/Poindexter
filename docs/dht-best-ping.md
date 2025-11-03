# Example: Find the best (lowest‑ping) peer in a DHT table

This example shows how to model a "made up" DHT routing table and use Poindexter's `KDTree` to quickly find:

- the single best peer by ping (nearest neighbor)
- the top N best peers by ping (k‑nearest neighbors)
- all peers under a ping threshold (radius search)

We keep it simple by mapping each peer to a 1‑dimensional coordinate: its ping in milliseconds. Using 1D means the KDTree's distance is just the absolute difference between pings.

> Tip: In a real system, you might expand to multiple dimensions (e.g., `[ping_ms, hop_count, geo_distance, score]`) and choose a metric (`L1`, `L2`, or `L∞`) that best matches your routing heuristic. See how to build normalized, weighted multi‑dimensional points with the public helpers `poindexter.Build2D/3D/4D` here: [Multi-Dimensional KDTree (DHT)](kdtree-multidimensional.md).

---

## Full example

```go
package main

import (
    "fmt"
    poindexter "github.com/Snider/Poindexter"
)

// Peer is our DHT peer entry (made up for this example).
type Peer struct {
    Addr string // multiaddr or host:port
    Ping int    // measured ping in milliseconds
}

func main() {
    // A toy DHT routing table with made-up ping values
    table := []Peer{
        {Addr: "peer1.example:4001", Ping: 74},
        {Addr: "peer2.example:4001", Ping: 52},
        {Addr: "peer3.example:4001", Ping: 110},
        {Addr: "peer4.example:4001", Ping: 35},
        {Addr: "peer5.example:4001", Ping: 60},
        {Addr: "peer6.example:4001", Ping: 44},
    }

    // Map peers to KD points in 1D where coordinate = ping (ms).
    // Use stable string IDs so we can delete/update later.
    pts := make([]poindexter.KDPoint[Peer], 0, len(table))
    for i, p := range table {
        pts = append(pts, poindexter.KDPoint[Peer]{
            ID:     fmt.Sprintf("peer-%d", i+1),
            Coords: []float64{float64(p.Ping)},
            Value:  p,
        })
    }

    // Build a KDTree. Euclidean metric is fine for 1D ping comparisons.
    kdt, err := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
    if err != nil {
        panic(err)
    }

    // 1) Find the best (lowest-ping) peer.
    //    Query is a 1D point representing desired ping target. Using 0 finds the min.
    best, d, ok := kdt.Nearest([]float64{0})
    if !ok {
        fmt.Println("no peers found")
        return
    }
    fmt.Printf("Best peer: %s (ping=%d ms), distance=%.0f\n", best.Value.Addr, best.Value.Ping, d)
    // Example output: Best peer: peer4.example:4001 (ping=35 ms), distance=35

    // 2) Top-N best peers by ping.
    top, dists := kdt.KNearest([]float64{0}, 3)
    fmt.Println("Top 3 peers by ping:")
    for i := range top {
        fmt.Printf("  #%d %s (ping=%d ms), distance=%.0f\n", i+1, top[i].Value.Addr, top[i].Value.Ping, dists[i])
    }

    // 3) All peers under a threshold (e.g., <= 50 ms): radius search.
    within, wd := kdt.Radius([]float64{0}, 50)
    fmt.Println("Peers with ping <= 50 ms:")
    for i := range within {
        fmt.Printf("  %s (ping=%d ms), distance=%.0f\n", within[i].Value.Addr, within[i].Value.Ping, wd[i])
    }

    // 4) Dynamic updates: if a peer improves ping, we can delete & re-insert with a new ID
    //    (or keep the same ID and just update the point if your application tracks indices).
    //    Here we simulate peer5 dropping from 60 ms to 30 ms.
    if kdt.DeleteByID("peer-5") {
        improved := poindexter.KDPoint[Peer]{
            ID:     "peer-5", // keep the same ID for simplicity
            Coords: []float64{30},
            Value:  Peer{Addr: "peer5.example:4001", Ping: 30},
        }
        _ = kdt.Insert(improved)
    }

    // Recompute the best after update
    best2, d2, _ := kdt.Nearest([]float64{0})
    fmt.Printf("After update, best peer: %s (ping=%d ms), distance=%.0f\n", best2.Value.Addr, best2.Value.Ping, d2)
}
```

### Why does querying with `[0]` work?
We use Euclidean distance in 1D, so `distance = |ping - target|`. With target `0`, minimizing the distance is equivalent to minimizing the ping itself.

### Extending the metric/space
- Multi-objective: encode more routing features (lower is better) as extra dimensions, e.g. `[ping_ms, hops, queue_delay_ms]`.
- Metric choice:
  - `EuclideanDistance` (L2): balances outliers smoothly.
  - `ManhattanDistance` (L1): linear penalty; robust for sparsity.
  - `ChebyshevDistance` (L∞): cares about the worst dimension.
- Normalization: when mixing units (ms, hops, km), normalize or weight dimensions so the metric reflects your priority.

### Notes
- This KDTree currently uses an internal linear scan for queries. The API is stable and designed so it can be swapped to use `gonum.org/v1/gonum/spatial/kdtree` under the hood later for sub-linear queries on large datasets.
- IDs are optional but recommended for O(1)-style deletes; keep them unique per tree.
