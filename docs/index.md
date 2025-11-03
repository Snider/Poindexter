# Poindexter

Welcome to the Poindexter Go library documentation!

## Overview

Poindexter is a Go library package licensed under EUPL-1.2.

## Features

- Simple and easy to use
- Comprehensive sorting utilities with custom comparators
- Generic sorting functions with type safety
- Binary search capabilities
- Well-documented API
- Comprehensive test coverage
- Cross-platform support

## Quick Start

Install the library:

```bash
go get github.com/Snider/Poindexter
```

Use it in your code:

```go
package main

import (
    "fmt"
    "github.com/Snider/Poindexter"
)

func main() {
    fmt.Println(poindexter.Hello("World"))
    fmt.Println("Version:", poindexter.Version())
}
```

## License

This project is licensed under the European Union Public Licence v1.2 (EUPL-1.2).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.


## Examples

- Find the best (lowest‑ping) DHT peer using KDTree: [Best Ping Peer (DHT)](dht-best-ping.md)
