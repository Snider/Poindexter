# Performance: KDTree benchmarks and guidance

This page summarizes how to measure KDTree performance in this repository and how to compare the two internal backends (Linear vs Gonum) that you can select at build/runtime.

## How benchmarks are organized

- Micro-benchmarks live in `bench_kdtree_test.go`, `bench_kdtree_dual_test.go`, and `bench_kdtree_dual_100k_test.go` and cover:
  - `Nearest` in 2D and 4D with N = 1k, 10k (both backends)
  - `Nearest` in 2D and 4D with N = 100k (gonum-tag job; linear also measured there)
  - `KNearest(k=10)` in 2D/4D with N = 1k, 10k
  - `Radius` (mid radius r≈0.5 after normalization) in 2D/4D with N = 1k, 10k
- Datasets: Uniform and 3-cluster synthetic generators in normalized [0,1] spaces.
- Backends: Linear (always available) and Gonum (enabled when built with `-tags=gonum`).

Run them locally:

```bash
# Linear backend (default)
go test -bench . -benchmem -run=^$ ./...

# Gonum backend (optimized KD; requires build tag)
go test -tags=gonum -bench . -benchmem -run=^$ ./...
```

GitHub Actions publishes benchmark artifacts on every push/PR:
- Linear job: artifact `bench-linear.txt`
- Gonum job: artifact `bench-gonum.txt`

## Backend selection and defaults

- Default backend is Linear.
- If you build with `-tags=gonum`, the default switches to the optimized KD backend.
- You can override at runtime:

```
// Force Linear
kdt, _ := poindexter.NewKDTree(pts, poindexter.WithBackend(poindexter.BackendLinear))
// Force Gonum (requires build tag)
kdt, _ := poindexter.NewKDTree(pts, poindexter.WithBackend(poindexter.BackendGonum))
```

Supported metrics in the optimized backend: L2 (Euclidean), L1 (Manhattan), L∞ (Chebyshev). Cosine/Weighted-Cosine currently use the Linear backend.

## What to expect (rule of thumb)

- Linear backend: O(n) per query; fast for small-to-medium datasets (≤10k), especially in low dims (≤4).
- Gonum backend: typically sub-linear for prunable datasets and dims ≤ ~8, with noticeable gains as N grows (≥10k–100k), especially on uniform or moderately clustered data and moderate radii.
- For large radii (many points within r) or highly correlated/pathological data, pruning may be less effective and behavior approaches O(n) even with KD-trees.

## Interpreting results

Benchmarks output something like:

```
BenchmarkNearest_10k_4D_Gonum_Uniform-8   50000  12,300 ns/op   0 B/op   0 allocs/op
```

- `ns/op`: lower is better (nanoseconds per operation)
- `B/op` and `allocs/op`: memory behavior; fewer is better
- `KNearest` incurs extra work due to sorting; `Radius` cost scales with the number of hits.

## Improving performance

- Normalize and weight features once; reuse across queries (see `Build*WithStats` helpers).
- Choose a metric aligned with your policy: L2 usually a solid default; L1 for per-axis penalties; L∞ for hard-threshold dominated objectives.
- Batch queries to benefit from CPU caches.
- Prefer the Gonum backend for larger N and dims ≤ ~8; stick to Linear for tiny datasets or when using Cosine metrics.

## Reproducing and tracking performance

- Local (Linear): `go test -bench . -benchmem -run=^$ ./...`
- Local (Gonum): `go test -tags=gonum -bench . -benchmem -run=^$ ./...`
- CI artifacts: download `bench-linear.txt` and `bench-gonum.txt` from the latest workflow run.
- Optional: add historical trend graphs via Benchstat or Codecov integration.
