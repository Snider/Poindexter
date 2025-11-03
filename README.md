# Poindexter

[![Go Reference](https://pkg.go.dev/badge/github.com/Snider/Poindexter.svg)](https://pkg.go.dev/github.com/Snider/Poindexter)
[![CI](https://github.com/Snider/Poindexter/actions/workflows/ci.yml/badge.svg)](https://github.com/Snider/Poindexter/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Snider/Poindexter)](https://goreportcard.com/report/github.com/Snider/Poindexter)
[![Vulncheck](https://img.shields.io/badge/govulncheck-enabled-brightgreen.svg)](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
[![codecov](https://codecov.io/gh/Snider/Poindexter/branch/main/graph/badge.svg)](https://codecov.io/gh/Snider/Poindexter)

A Go library package providing utility functions including sorting algorithms with custom comparators.

## Features

- 🔢 **Sorting Utilities**: Sort integers, strings, and floats in ascending or descending order
- 🎯 **Custom Sorting**: Sort any type with custom comparison functions or key extractors
- 🔍 **Binary Search**: Fast search on sorted data
- 🧭 **KDTree (NN Search)**: Build a KDTree over points with generic payloads; nearest, k-NN, and radius queries with Euclidean or Manhattan metrics
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

### KDTree performance and notes
- Current KDTree queries are O(n) linear scans, which are great for small-to-medium datasets or low-latency prototyping. For 1e5+ points and low/medium dimensions, consider swapping the internal engine to `gonum.org/v1/gonum/spatial/kdtree` (the API here is compatible by design).
- Insert is O(1) amortized; delete by ID is O(1) via swap-delete; order is not preserved.
- Concurrency: the KDTree type is not safe for concurrent mutation. Protect with a mutex or share immutable snapshots for read-mostly workloads.
- See multi-dimensional examples (ping/hops/geo/score) in docs and `examples/`.

## License

This project is licensed under the European Union Public Licence v1.2 (EUPL-1.2). See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.