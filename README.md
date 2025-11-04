# Poindexter

[![Go Reference](https://pkg.go.dev/badge/github.com/Snider/Poindexter.svg)](https://pkg.go.dev/github.com/Snider/Poindexter)
[![CI](https://github.com/Snider/Poindexter/actions/workflows/ci.yml/badge.svg)](https://github.com/Snider/Poindexter/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Snider/Poindexter)](https://goreportcard.com/report/github.com/Snider/Poindexter)
[![Vulncheck](https://img.shields.io/badge/govulncheck-enabled-brightgreen.svg)](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
[![codecov](https://codecov.io/gh/Snider/Poindexter/branch/main/graph/badge.svg)](https://codecov.io/gh/Snider/Poindexter)
[![Release](https://img.shields.io/github/v/release/Snider/Poindexter?display_name=tag)](https://github.com/Snider/Poindexter/releases)

A Go library package providing utility functions including sorting algorithms with custom comparators.

## Features

- 🔢 **Sorting Utilities**: Sort integers, strings, and floats in ascending or descending order
- 🎯 **Custom Sorting**: Sort any type with custom comparison functions or key extractors
- 🔍 **Binary Search**: Fast search on sorted data
- 🧭 **KDTree (NN Search)**: Build a KDTree over points with generic payloads; nearest, k-NN, and radius queries with Euclidean, Manhattan, Chebyshev, and Cosine metrics
- 📦 **Generic Functions**: Type-safe operations using Go generics
- ✅ **Well-Tested**: Comprehensive test coverage
- 📖 **Documentation**: Full documentation available at GitHub Pages

## Installation

```bash
go get github.com/Snider/Poindexter
```

## Quick Start

```go
package main

import (
    "fmt"
    poindexter "github.com/Snider/Poindexter"
)

func main() {
    // Basic sorting
    numbers := []int{3, 1, 4, 1, 5, 9}
    poindexter.SortInts(numbers)
    fmt.Println(numbers) // [1 1 3 4 5 9]

    // Custom sorting with key function
    type Product struct {
        Name  string
        Price float64
    }

    products := []Product{{"Apple", 1.50}, {"Banana", 0.75}, {"Cherry", 3.00}}
    poindexter.SortByKey(products, func(p Product) float64 { return p.Price })

    // KDTree quick demo
    pts := []poindexter.KDPoint[string]{
        {ID: "A", Coords: []float64{0, 0}, Value: "alpha"},
        {ID: "B", Coords: []float64{1, 0}, Value: "bravo"},
        {ID: "C", Coords: []float64{0, 1}, Value: "charlie"},
    }
    tree, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
    nearest, dist, _ := tree.Nearest([]float64{0.9, 0.1})
    fmt.Println(nearest.ID, nearest.Value, dist) // B bravo ~0.141...
}
```

## Documentation

Full documentation is available at [https://snider.github.io/Poindexter/](https://snider.github.io/Poindexter/)

Explore runnable examples in the repository:
- examples/dht_ping_1d
- examples/kdtree_2d_ping_hop
- examples/kdtree_3d_ping_hop_geo
- examples/kdtree_4d_ping_hop_geo_score
- examples/wasm-browser (browser demo using the ESM loader)

### KDTree performance and notes
- Dual backend support: Linear (always available) and an optimized KD backend enabled when building with `-tags=gonum`. Linear is the default; with the `gonum` tag, the optimized backend becomes the default.
- Complexity: Linear backend is O(n) per query. Optimized KD backend is typically sub-linear on prunable datasets and dims ≤ ~8, especially as N grows (≥10k–100k).
- Insert is O(1) amortized; delete by ID is O(1) via swap-delete; order is not preserved.
- Concurrency: the KDTree type is not safe for concurrent mutation. Protect with a mutex or share immutable snapshots for read-mostly workloads.
- See multi-dimensional examples (ping/hops/geo/score) in docs and `examples/`.
- Performance guide: see docs/Performance for benchmark guidance and tips: [docs/perf.md](docs/perf.md) • Hosted: https://snider.github.io/Poindexter/perf/

### Backend selection
- Default backend is Linear. If you build with `-tags=gonum`, the default becomes the optimized KD backend.
- You can override per tree at construction:

```go
// Force Linear (always available)
kdt1, _ := poindexter.NewKDTree(pts, poindexter.WithBackend(poindexter.BackendLinear))

// Force Gonum (requires build tag)
kdt2, _ := poindexter.NewKDTree(pts, poindexter.WithBackend(poindexter.BackendGonum))
```

- Supported metrics in the optimized backend: Euclidean (L2), Manhattan (L1), Chebyshev (L∞).
- Cosine and Weighted-Cosine currently run on the Linear backend.
- See the Performance guide for measured comparisons and when to choose which backend.

#### Choosing a metric (quick tips)
- Euclidean (L2): smooth trade-offs across axes; solid default for blended preferences.
- Manhattan (L1): emphasizes per-axis absolute differences; good when each unit of ping/hop matters equally.
- Chebyshev (L∞): dominated by the worst axis; useful for strict thresholds (e.g., reject high hop count regardless of ping).
- Cosine: angle-based for vector similarity; pair it with normalized/weighted features when direction matters more than magnitude.

See the multi-dimensional KDTree docs for end-to-end examples and weighting/normalization helpers: [Multi-Dimensional KDTree (DHT)](docs/kdtree-multidimensional.md).

## Maintainer Makefile

The repository includes a maintainer-friendly `Makefile` that mirrors CI tasks and speeds up local workflows.

- help — list available targets
- tidy / tidy-check — run `go mod tidy`, optionally verify no diffs
- fmt — format code (`go fmt ./...`)
- vet — `go vet ./...`
- build — `go build ./...`
- examples — build all programs under `examples/` (if present)
- test — run unit tests
- race — run tests with the race detector
- cover — run tests with race + coverage (writes `coverage.out` and prints summary)
- coverhtml — render HTML coverage report to `coverage.html`
- coverfunc — print per-function coverage (from `coverage.out`)
- cover-kdtree — print coverage details filtered to `kdtree.go`
- fuzz — run Go fuzzing for a configurable time (default 10s) matching CI
- bench — run benchmarks with `-benchmem` (writes `bench.txt`)
- lint — run `golangci-lint` (if installed)
- vuln — run `govulncheck` (if installed)
- ci — CI-parity aggregate: tidy-check, build, vet, cover, examples, bench, lint, vuln
- release — run GoReleaser with the canonical `.goreleaser.yaml` (for tagged releases)
- snapshot — GoReleaser snapshot (no publish)
- docs-serve — serve MkDocs locally on 127.0.0.1:8000
- docs-build — build MkDocs site into `site/`

Quick usage:

- See all targets:

```bash
make help
```

- Fast local cycle:

```bash
make fmt
make vet
make test
```

- CI-parity run (what GitHub Actions does, locally):

```bash
make ci
```

- Coverage summary:

```bash
make cover
```

- Generate HTML coverage report (writes coverage.html):

```bash
make coverhtml
```

- Fuzz for 10 seconds (default):

```bash
make fuzz
```

- Fuzz with a custom time (e.g., 30s):

```bash
make fuzz FUZZTIME=30s
```

- Run benchmarks (writes bench.txt):

```bash
make bench
```

- Build examples (if any under ./examples):

```bash
make examples
```

- Serve docs locally (requires mkdocs-material):

```bash
make docs-serve
```

Configurable variables:

- `FUZZTIME` (default `10s`) — e.g. `make fuzz FUZZTIME=30s`
- `BENCHOUT` (default `bench.txt`), `COVEROUT` (default `coverage.out`), `COVERHTML` (default `coverage.html`)
- Tool commands are overridable via env: `GO`, `GOLANGCI_LINT`, `GORELEASER`, `MKDOCS`

Requirements for optional targets:

- `golangci-lint` for `make lint`
- `golang.org/x/vuln/cmd/govulncheck` for `make vuln`
- `goreleaser` for `make release` / `make snapshot`
- `mkdocs` + `mkdocs-material` for `make docs-serve` / `make docs-build`

See the full Makefile at the repo root for authoritative target definitions.

## License

This project is licensed under the European Union Public Licence v1.2 (EUPL-1.2). See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.