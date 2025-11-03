# Poindexter

A Go library package providing utility functions including sorting algorithms with custom comparators.

## Features

- 🔢 **Sorting Utilities**: Sort integers, strings, and floats in ascending or descending order
- 🎯 **Custom Sorting**: Sort any type with custom comparison functions or key extractors
- 🔍 **Binary Search**: Fast search on sorted data
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
    "github.com/Snider/Poindexter"
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
    
    products := []Product{
        {"Apple", 1.50},
        {"Banana", 0.75},
        {"Cherry", 3.00},
    }
    
    poindexter.SortByKey(products, func(p Product) float64 {
        return p.Price
    })
}
```

## Documentation

Full documentation is available at [https://snider.github.io/Poindexter/](https://snider.github.io/Poindexter/)

## License

This project is licensed under the European Union Public Licence v1.2 (EUPL-1.2). See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.