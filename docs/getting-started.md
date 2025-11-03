# Getting Started

This guide will help you get started with the Poindexter library.

## Installation

To install Poindexter, use `go get`:

```bash
go get github.com/Snider/Poindexter
```

## Basic Usage

### Importing the Library

```go
import "github.com/Snider/Poindexter"
```

### Using the Hello Function

The `Hello` function returns a greeting message:

```go
package main

import (
    "fmt"
    "github.com/Snider/Poindexter"
)

func main() {
    // Say hello to the world
    fmt.Println(poindexter.Hello(""))
    // Output: Hello, World!

    // Say hello to someone specific
    fmt.Println(poindexter.Hello("Poindexter"))
    // Output: Hello, Poindexter!
}
```

### Getting the Version

You can check the library version:

```go
package main

import (
    "fmt"
    "github.com/Snider/Poindexter"
)

func main() {
    version := poindexter.Version()
    fmt.Println("Library version:", version)
}
```

## Sorting Data

Poindexter includes comprehensive sorting utilities:

### Basic Sorting

```go
package main

import (
    "fmt"
    "github.com/Snider/Poindexter"
)

func main() {
    // Sort integers
    numbers := []int{3, 1, 4, 1, 5, 9}
    poindexter.SortInts(numbers)
    fmt.Println(numbers) // [1 1 3 4 5 9]

    // Sort strings
    words := []string{"banana", "apple", "cherry"}
    poindexter.SortStrings(words)
    fmt.Println(words) // [apple banana cherry]
}
```

### Advanced Sorting with Custom Keys

```go
package main

import (
    "fmt"
    "github.com/Snider/Poindexter"
)

type Product struct {
    Name  string
    Price float64
}

func main() {
    products := []Product{
        {"Apple", 1.50},
        {"Banana", 0.75},
        {"Cherry", 3.00},
    }

    // Sort by price using SortByKey
    poindexter.SortByKey(products, func(p Product) float64 {
        return p.Price
    })

    for _, p := range products {
        fmt.Printf("%s: $%.2f\n", p.Name, p.Price)
    }
}
```

## Next Steps

- Check out the [API Reference](api.md) for detailed documentation
- Read about the [License](license.md)
