# Troubleshooting Guide

This guide provides solutions to common problems you may encounter when using the Poindexter library.

## General

### `go get` command fails

If you are having trouble installing the library with `go get`, make sure you have a working Go environment and that you have network connectivity. You can check your Go installation by running `go version`.

## K-D Tree

### `ErrDimMismatch` error when building a K-D tree

This error occurs when the points you are using to build the K-D tree have inconsistent dimensions. Make sure that all points have the same number of coordinates.

### `ErrDuplicateID` error when building a K-D tree

This error occurs when you are using points with duplicate IDs. If you are using IDs, they must be unique for each point.

### Nearest neighbor search is slow

The K-D tree is designed for efficient nearest neighbor searches, but its performance can be affected by the distribution of your data. If you are experiencing slow performance, consider the following:

-   **Data Distribution:** The K-D tree performs best with uniformly distributed data. If your data is highly clustered, the performance may be degraded.
-   **Dimensionality:** The performance of the K-D tree can degrade as the number of dimensions increases. For very high-dimensional data, other methods may be more appropriate.

## WASM

### WASM module fails to load

If the WASM module is not loading, make sure that the `.wasm` file is being served correctly and that the path to the file is correct in your HTML. You can check the browser's developer console for any errors related to loading the module.

### Performance is slow in the browser

The WASM build will generally be slower than the native Go implementation. If you are experiencing performance issues, consider the following:

-   **Data Size:** The amount of data you are processing can affect performance. Try to minimize the amount of data you are sending to the WASM module.
-   **Benchmarking:** Benchmark your specific use case to identify any performance bottlenecks.
