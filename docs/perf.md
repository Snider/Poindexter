# Performance: KDTree benchmarks and guidance

This page summarizes how to measure KDTree performance in this repository and when to consider switching the internal engine to `gonum.org/v1/gonum/spatial/kdtree` for large datasets.

## How benchmarks are organized

- Micro-benchmarks live in `bench_kdtree_test.go` and cover:
  - `Nearest` in 2D and 4D with N = 1k, 10k
  - `KNearest(k=10)` in 2D with N = 1k, 10k
  - `Radius` (mid radius) in 2D with N = 1k, 10k
- All benchmarks operate in normalized [0,1] spaces and use the current linear-scan implementation.

Run them locally:

```bash
go test -bench . -benchmem -run=^$ ./...
```

GitHub Actions publishes benchmark artifacts for Go 1.23 on every push/PR. Look for artifacts named `bench-<go-version>.txt` in the CI run.

## What to expect (rule of thumb)

- Time complexity is O(n) per query in the current implementation.
- For small-to-medium datasets (up to ~10k points), linear scans are often fast enough, especially for low dimensionality (≤4) and if queries are batched efficiently.
- For larger datasets (≥100k) and low/medium dimensions (≤8), a true KD-tree (like Gonum’s) often yields sub-linear queries and significantly lower latency.

## Interpreting results

Benchmarks output something like:

```
BenchmarkNearest_10k_4D-8      50000         23,000 ns/op      0 B/op      0 allocs/op
```

- `ns/op`: lower is better (nanoseconds per operation)
- `B/op` and `allocs/op`: memory behavior; fewer is better

Because `KNearest` sorts by distance, you should expect additional cost over `Nearest`. `Radius` cost depends on how many points fall within the radius; tighter radii usually run faster.

## Improving performance

- Prefer Euclidean (L2) over metrics that require extra branching for CPU pipelines, unless your policy prefers otherwise.
- Normalize and weight features once; reuse coordinates across queries.
- Batch queries to amortize overhead of data locality and caches.
- Consider a backend swap to Gonum’s KD-tree for large N (we plan to add a `WithBackend("gonum")` option).

## Reproducing and tracking performance

- Local: run `go test -bench . -benchmem -run=^$ ./...`
- CI: download `bench-*.txt` artifacts from the latest workflow run
- Optional: we can add historical trend graphs via Codecov or Benchstat integration if desired.
