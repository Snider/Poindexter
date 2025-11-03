# Poindexter

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

## License

This project is licensed under the European Union Public Licence v1.2 (EUPL-1.2). See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.