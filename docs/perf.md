# Performance: KDTree benchmarks and guidance

This page summarizes how to measure KDTree performance in this repository and how to compare the two internal backends (Linear vs Gonum) that you can select at build/runtime.

## How benchmarks are organized

- Micro-benchmarks live in `bench_kdtree_test.go` and `bench_kdtree_dual_test.go` and cover:
  - `Nearest` in 2D and 4D with N = 1k, 10k
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

## Sample results (from a recent local run)

Results vary by machine, Go version, and dataset seed. The following run was captured locally and is provided as a reference point.

- Machine: darwin/arm64, Apple M3 Ultra
- Package: `github.com/Snider/Poindexter`
- Command: `go test -bench . -benchmem -run=^$ ./... | tee bench.txt`

Full output:

```
goos: darwin
goarch: arm64
pkg: github.com/Snider/Poindexter
BenchmarkNearest_Linear_Uniform_1k_2D-32          	 409321	      3001 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Gonum_Uniform_1k_2D-32           	 413823	      2888 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Linear_Uniform_10k_2D-32         	  43053	     27809 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Gonum_Uniform_10k_2D-32          	  42996	     27936 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Linear_Uniform_1k_4D-32          	 326492	      3746 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Gonum_Uniform_1k_4D-32           	 338983	      3857 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Linear_Uniform_10k_4D-32         	  35661	     32985 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Gonum_Uniform_10k_4D-32          	  35678	     33388 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Linear_Clustered_1k_2D-32        	 425220	      2874 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Gonum_Clustered_1k_2D-32         	 420080	      2849 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Linear_Clustered_10k_2D-32       	  43242	     27776 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_Gonum_Clustered_10k_2D-32        	  42392	     27889 ns/op	       0 B/op	       0 allocs/op
BenchmarkKNN10_Linear_Uniform_10k_2D-32           	   1206	    977599 ns/op	  164492 B/op	       6 allocs/op
BenchmarkKNN10_Gonum_Uniform_10k_2D-32            	   1239	    972501 ns/op	  164488 B/op	       6 allocs/op
BenchmarkKNN10_Linear_Clustered_10k_2D-32         	   1219	    973242 ns/op	  164492 B/op	       6 allocs/op
BenchmarkKNN10_Gonum_Clustered_10k_2D-32          	   1214	    971017 ns/op	  164488 B/op	       6 allocs/op
BenchmarkRadiusMid_Linear_Uniform_10k_2D-32       	   1279	    917692 ns/op	  947529 B/op	      23 allocs/op
BenchmarkRadiusMid_Gonum_Uniform_10k_2D-32        	   1299	    918176 ns/op	  947529 B/op	      23 allocs/op
BenchmarkRadiusMid_Linear_Clustered_10k_2D-32     	   1059	   1123281 ns/op	 1217866 B/op	      24 allocs/op
BenchmarkRadiusMid_Gonum_Clustered_10k_2D-32      	   1063	   1149507 ns/op	 1217871 B/op	      24 allocs/op
BenchmarkNearest_1k_2D-32                         	 401595	      2964 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_10k_2D-32                        	  42129	     28229 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_1k_4D-32                         	 365626	      3642 ns/op	       0 B/op	       0 allocs/op
BenchmarkNearest_10k_4D-32                        	  36298	     33176 ns/op	       0 B/op	       0 allocs/op
BenchmarkKNearest10_1k_2D-32                      	  20348	     59568 ns/op	   17032 B/op	       6 allocs/op
BenchmarkKNearest10_10k_2D-32                     	   1224	    969093 ns/op	  164488 B/op	       6 allocs/op
BenchmarkRadiusMid_1k_2D-32                       	  21867	     53273 ns/op	   77512 B/op	      16 allocs/op
BenchmarkRadiusMid_10k_2D-32                      	   1302	    933791 ns/op	  955720 B/op	      23 allocs/op
PASS
ok  	github.com/Snider/Poindexter	40.102s
PASS
ok  	github.com/Snider/Poindexter/examples/dht_ping_1d	0.348s
PASS
ok  	github.com/Snider/Poindexter/examples/kdtree_2d_ping_hop	0.266s
PASS
ok  	github.com/Snider/Poindexter/examples/kdtree_3d_ping_hop_geo	0.272s
PASS
ok  	github.com/Snider/Poindexter/examples/kdtree_4d_ping_hop_geo_score	0.269s
```

Notes:
- The first block shows dual-backend benchmarks (Linear vs Gonum) for uniform and clustered datasets at 2D/4D with N=1k/10k.
- The final block includes the legacy single-backend benchmarks for additional sizes; both are useful for comparison.

To compare against the optimized KD backend explicitly, build with `-tags=gonum` and/or download `bench-gonum.txt` from CI artifacts.
