# Frequently Asked Questions (FAQ)

This page answers common questions about the Poindexter library.

## General

### What is Poindexter?

Poindexter is a Go library that provides a collection of utility functions, including sorting algorithms and a K-D tree implementation for nearest neighbor searches.

### What is the license?

Poindexter is licensed under the European Union Public Licence v1.2 (EUPL-1.2). See the [LICENSE](LICENSE) file for more details.

## K-D Tree

### What is a K-D tree?

A K-D tree is a data structure used for organizing points in a k-dimensional space. It is particularly useful for nearest neighbor searches.

### How do I choose a distance metric?

The choice of distance metric depends on your specific use case. Here are some general guidelines:

-   **Euclidean (L2):** A good default for most cases.
-   **Manhattan (L1):** Useful when movement is restricted to a grid.
-   **Chebyshev (L∞):** Useful for cases where the maximum difference between coordinates is the most important factor.

### Is the K-D tree thread-safe?

The K-D tree implementation is not safe for concurrent mutations. If you need to use it in a concurrent environment, you should protect it with a mutex or share immutable snapshots for read-mostly workloads.

## WASM

### Can I use Poindexter in the browser?

Yes, Poindexter provides a WebAssembly (WASM) build that allows you to use the K-D tree in browser environments. See the [WASM documentation](wasm.md) for more details.

### What is the performance of the WASM build?

The performance of the WASM build is generally good, but it will be slower than the native Go implementation. For performance-critical applications, it is recommended to benchmark your specific use case.
