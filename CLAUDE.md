# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Poindexter is a Go library providing:
- **Sorting utilities** with custom comparators (ints, strings, floats, generic types)
- **KDTree** for nearest-neighbor search with multiple distance metrics (Euclidean, Manhattan, Chebyshev, Cosine)
- **Helper functions** for building normalized/weighted KD points (2D/3D/4D/ND)
- **DNS/RDAP tools** for domain, IP, and ASN lookups
- **Analytics** for tree operations and peer selection tracking
- **State index sidecar** (`stateindex/`) mapping index URIs to byte ranges inside
  State `.kv` container files

## Build Commands

```bash
# Run all tests
make test

# Run tests with race detector
make race

# Run tests with coverage (writes coverage.out)
make cover

# Run a single test
go test -run TestFunctionName ./...

# Build all packages
make build

# Build with gonum optimized backend
go build -tags=gonum ./...

# Run tests with gonum backend
go test -tags=gonum ./...

# Build WebAssembly module
make wasm-build

# Run benchmarks (writes bench.txt)
make bench

# Run benchmarks for specific backend
make bench-linear     # Linear backend
make bench-gonum      # Gonum backend (requires -tags=gonum)

# Fuzz testing (default 10s)
make fuzz
make fuzz FUZZTIME=30s

# Lint and static analysis
make lint   # requires golangci-lint
make vuln   # requires govulncheck

# CI-parity local run
make ci

# Format code
make fmt

# Tidy modules
make tidy
```

## Architecture

### Core Components

**kdtree.go** - Main KDTree implementation:
- `KDTree[T]` generic struct with payload type T
- `KDPoint[T]` represents points with ID, Coords, and Value
- Query methods: `Nearest()`, `KNearest()`, `Radius()`
- Mutation methods: `Insert()`, `DeleteByID()`
- Distance metrics implement `DistanceMetric` interface

**Dual Backend System:**
- `kdtree_gonum.go` (build tag `gonum`) - Optimized median-split KD-tree with branch-and-bound pruning
- `kdtree_gonum_stub.go` (default) - Stubs when gonum tag not set
- Linear backend is always available; gonum backend is default when tag is enabled
- Backend selected via `WithBackend()` option or defaults based on build tags

**kdtree_helpers.go** - Point construction utilities:
- `BuildND()`, `Build2D()`, `Build3D()`, `Build4D()` - Build normalized/weighted KD points
- `ComputeNormStatsND()` - Compute per-axis min/max for normalization
- Supports axis inversion and per-axis weights

**kdtree_analytics.go** - Operational metrics:
- `TreeAnalytics` - Query/insert/delete counts, timing stats
- `PeerAnalytics` - Per-peer selection tracking for DHT/NAT routing
- Distribution statistics and NAT routing metrics

**sort.go** - Sorting utilities:
- Type-specific sorts: `SortInts()`, `SortStrings()`, `SortFloat64s()`
- Generic sorts: `SortBy()`, `SortByKey()`
- Binary search: `BinarySearch()`, `BinarySearchStrings()`

**stateindex/** - Sidecar index over State `.kv` containers:
- `Index` / `Entry` — a JSON document listing segments, each naming a container
  file plus the payload's offset and length inside it
- `Append()`/`Set()`/`Lookup()`/`Remove()` keyed on `index_uri`; `ByPath()`/`Paths()`
  for the one-file and many-file cases
- Nothing here opens a container. Ranges are recorded at write time and taken at
  their word, so resolving an index URI never touches the binary State payload —
  entries may name files that do not exist

**dns_tools.go** - DNS and RDAP lookup utilities:
- DNS record types and lookup functions
- RDAP (modern WHOIS) for domains, IPs, ASNs
- External tool link generators

### Examples Directory

Runnable examples demonstrating KDTree usage patterns:
- `examples/dht_ping_1d/` - 1D DHT with ping latency
- `examples/kdtree_2d_ping_hop/` - 2D with ping + hop count
- `examples/kdtree_3d_ping_hop_geo/` - 3D with geographic distance
- `examples/kdtree_4d_ping_hop_geo_score/` - 4D with trust scores
- `examples/dht_helpers/` - Convenience wrappers for common DHT schemas
- `examples/wasm-browser-ts/` - TypeScript + Vite browser demo

### WebAssembly

WASM module in `wasm/main.go` exposes KDTree functionality to JavaScript. Build outputs to `dist/`.

## Key Patterns

### KDTree Construction
```go
pts := []poindexter.KDPoint[string]{
    {ID: "A", Coords: []float64{0, 0}, Value: "alpha"},
}
tree, err := poindexter.NewKDTree(pts,
    poindexter.WithMetric(poindexter.EuclideanDistance{}),
    poindexter.WithBackend(poindexter.BackendGonum))
```

### Building Normalized Points
```go
pts, err := poindexter.BuildND(items,
    func(p Peer) string { return p.ID },           // ID extractor
    []func(Peer) float64{getPing, getHops},        // Feature extractors
    []float64{1.0, 0.5},                           // Per-axis weights
    []bool{false, false})                          // Inversion flags
```

### Distance Metrics
- `EuclideanDistance{}` - L2 norm (default)
- `ManhattanDistance{}` - L1 norm
- `ChebyshevDistance{}` - L-infinity (max) norm
- `CosineDistance{}` - 1 - cosine similarity
- `WeightedCosineDistance{Weights: []float64{...}}` - Weighted cosine

Note: Cosine metrics use linear backend only; L2/L1/L-infinity work with gonum backend.

## Testing

Tests are organized by component:
- `kdtree_test.go`, `kdtree_nd_test.go` - Core KDTree tests
- `kdtree_gonum_test.go` - Gonum backend specific tests
- `kdtree_backend_parity_test.go` - Backend equivalence tests
- `fuzz_kdtree_test.go` - Fuzz targets
- `bench_*.go` - Benchmark suites

Coverage targets:
```bash
make cover           # Summary
make coverfunc       # Per-function breakdown
make cover-kdtree    # kdtree.go only
make coverhtml       # HTML report
```

## License

EUPL-1.2 (European Union Public Licence v1.2)