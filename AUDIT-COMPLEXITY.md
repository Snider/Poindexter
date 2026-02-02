# Code Complexity & Maintainability Audit

This document analyzes the code quality of the Poindexter library, identifies maintainability issues, and provides recommendations for improvement. The audit focuses on cyclomatic and cognitive complexity, code duplication, and other maintainability metrics.

## High-Impact Findings

This section summarizes the most critical issues that should be prioritized for refactoring.

## Detailed Findings

This section provides a detailed breakdown of all findings, categorized by file.

### `kdtree.go`

Overall, `kdtree.go` is well-structured and maintainable. The complexity is low, and the code is easy to understand. The following are minor points for consideration.

| Finding | Severity | Recommendation |
| --- | --- | --- |
| Mild Feature Envy | Low | The `KDTree` methods directly call analytics recording functions (`t.analytics.RecordQuery`, `t.peerAnalytics.RecordSelection`). This creates a tight coupling between the core tree logic and the analytics subsystem. |
| Minor Code Duplication | Low | The query methods (`Nearest`, `KNearest`, `Radius`) share some boilerplate code for dimension checks, analytics timing, and handling the `gonum` backend. |

#### Recommendations

**1. Decouple Analytics from Core Logic**

To address the feature envy, the analytics recording could be decoupled from the core tree operations. This would improve separation of concerns and make the `KDTree` struct more focused on its primary responsibility.

*   **Refactoring Approach**: Introduce an interface for a "query observer" or use an event-based system. The `KDTree` would accept an observer during construction and notify it of events like "query started," "query completed," and "point selected."
*   **Design Pattern**: Observer pattern or a simple event emitter.

**Example (Conceptual)**:

```go
// In kdtree.go
type QueryObserver[T any] interface {
    OnQueryStart(t *KDTree[T])
    OnQueryEnd(t *KDTree[T], duration time.Duration)
    OnPointSelected(p KDPoint[T], distance float64)
}

// KDTree would have:
// observer QueryObserver[T]

// In Nearest method:
// if t.observer != nil { t.observer.OnQueryStart(t) }
// ...
// if t.observer != nil { t.observer.OnPointSelected(p, bestDist) }
// ...
// if t.observer != nil { t.observer.OnQueryEnd(t, time.Since(start)) }

```

**2. Reduce Boilerplate in Query Methods**

The minor code duplication in the query methods could be reduced by extracting the common setup and teardown logic into a helper function. However, given that there are only three methods and the duplication is minimal, this is a very low-priority change.

*   **Refactoring Approach**: Create a private helper function that takes a query function as an argument and handles the boilerplate.

**Example (Conceptual)**:

```go
// In kdtree.go
func (t *KDTree[T]) executeQuery(query []float64, fn func() (any, any)) (any, any) {
    if len(query) != t.dim || t.Len() == 0 {
        return nil, nil
    }
    start := time.Now()
    defer func() {
        if t.analytics != nil {
            t.analytics.RecordQuery(time.Since(start).Nanoseconds())
        }
    }()
    return fn()
}

// KNearest would then be simplified:
// func (t *KDTree[T]) KNearest(query []float64, k int) ([]KDPoint[T], []float64) {
//     res, dists := t.executeQuery(query, func() (any, any) {
//         // ... core logic of KNearest ...
//         return neighbors, dists
//     })
//     return res.([]KDPoint[T]), dists.([]float64)
// }
```
This approach has tradeoffs with readability and type casting, so it should be carefully considered.

### `kdtree_helpers.go`

This file contains significant code duplication and could be simplified by consolidating redundant functions.

| Finding | Severity | Recommendation |
| --- | --- | --- |
| High Code Duplication | High | The functions `Build2D`, `Build3D`, `Build4D` and their `ComputeNormStats` counterparts are nearly identical. They should be removed in favor of the existing generic `BuildND` and `ComputeNormStatsND` functions. |
| Long Parameter Lists | Medium | Functions like `BuildND` accept a large number of parameters (`items`, `id`, `features`, `weights`, `invert`). This can be improved by introducing a configuration struct. |

#### Recommendations

**1. Consolidate Builder Functions (High Priority)**

The most impactful change would be to remove the duplicated builder functions. The generic `BuildND` function already provides the same functionality in a more flexible way.

*   **Refactoring Approach**:
    1.  Mark `Build2D`, `Build3D`, `Build4D`, `ComputeNormStats2D`, `ComputeNormStats3D`, and `ComputeNormStats4D` as deprecated.
    2.  In a future major version, remove the deprecated functions.
    3.  Update all internal call sites and examples to use `BuildND` and `ComputeNormStatsND`.

*   **Code Example of Improvement**:

Instead of calling the specialized function:
```go
// Before
pts, err := Build3D(
    items,
    func(it MyType) string { return it.ID },
    func(it MyType) float64 { return it.Feature1 },
    func(it MyType) float64 { return it.Feature2 },
    func(it MyType) float64 { return it.Feature3 },
    [3]float64{1.0, 2.0, 0.5},
    [3]bool{false, true, false},
)
```

The code would be refactored to use the generic version:
```go
// After
features := []func(MyType) float64{
    func(it MyType) float64 { return it.Feature1 },
    func(it MyType) float64 { return it.Feature2 },
    func(it MyType) float64 { return it.Feature3 },
}
weights := []float64{1.0, 2.0, 0.5}
invert := []bool{false, true, false}

pts, err := BuildND(
    items,
    func(it MyType) string { return it.ID },
    features,
    weights,
    invert,
)
```

**2. Introduce a Parameter Object for Configuration**

To make the function signatures cleaner and more extensible, a configuration struct could be used for the builder functions.

*   **Design Pattern**: Introduce Parameter Object.

*   **Code Example of Improvement**:

```go
// Define a configuration struct
type BuildConfig[T any] struct {
    IDFunc   func(T) string
    Features []func(T) float64
    Weights  []float64
    Invert   []bool
    Stats    *NormStats // Optional: for building with pre-computed stats
}

// Refactor BuildND to accept the config
func BuildND[T any](items []T, config BuildConfig[T]) ([]KDPoint[T], error) {
    // ... logic using config fields ...
}

// Example usage
config := BuildConfig[MyType]{
    IDFunc:   func(it MyType) string { return it.ID },
    Features: features,
    Weights:  weights,
    Invert:   invert,
}
pts, err := BuildND(items, config)
```
This makes the call site cleaner and adding new options in the future would not require changing the function signature.

### `kdtree_analytics.go`

This file has the most significant maintainability issues in the codebase. It appears to be a "God Class" that has accumulated multiple, loosely related responsibilities over time.

| Finding | Severity | Recommendation |
| --- | --- | --- |
| God Class / Low Cohesion | High | The file contains logic for: tree performance analytics, peer selection analytics, statistical distribution calculations, NAT routing metrics, and peer trust/reputation. These are all distinct concerns. |
| Speculative Generality | Medium | The file includes many structs and functions related to NAT routing and peer trust (`NATRoutingMetrics`, `TrustMetrics`, `PeerQualityScore`) that may not be used by all consumers of the library. This adds complexity for users who only need the core k-d tree functionality. |
| Magic Numbers | Low | The `PeerQualityScore` function contains several magic numbers for weighting and normalization (e.g., `metrics.AvgRTTMs/1000.0`, `metrics.BandwidthMbps/100.0`). These should be extracted into named constants. |

#### Recommendations

**1. Decompose the God Class (High Priority)**

The most important refactoring is to break `kdtree_analytics.go` into smaller, more focused files. This will improve cohesion and make the code easier to navigate and maintain.

*   **Refactoring Approach**:
    1.  **`tree_analytics.go`**: Keep `TreeAnalytics` and `TreeAnalyticsSnapshot`.
    2.  **`peer_analytics.go`**: Move `PeerAnalytics` and `PeerStats`.
    3.  **`stats.go`**: Move `DistributionStats`, `ComputeDistributionStats`, `percentile`, `AxisDistribution`, and `ComputeAxisDistributions`.
    4.  **`p2p_metrics.go`**: Create a new file for all the peer-to-peer and networking-specific logic. This would include `NATRoutingMetrics`, `QualityWeights`, `PeerQualityScore`, `TrustMetrics`, `ComputeTrustScore`, `StandardPeerFeatures`, etc. This makes it clear that this functionality is for a specific domain (P2P networking) and is not part of the core k-d tree library.

**2. Extract Magic Numbers as Constants**

The magic numbers in `PeerQualityScore` should be replaced with named constants to improve readability and make them easier to modify.

*   **Code Example of Improvement**:

```go
// Before
latencyScore := 1.0 - math.Min(metrics.AvgRTTMs/1000.0, 1.0)
bandwidthScore := math.Min(metrics.BandwidthMbps/100.0, 1.0)

// After
const (
    maxAcceptableRTTMs = 1000.0
    excellentBandwidthMbps = 100.0
)
latencyScore := 1.0 - math.Min(metrics.AvgRTTMs/maxAcceptableRTTMs, 1.0)
bandwidthScore := math.Min(metrics.BandwidthMbps/excellentBandwidthMbps, 1.0)
```

**3. Isolate Domain-Specific Logic**

By moving the P2P-specific logic into its own file (`p2p_metrics.go`), it becomes clearer that this is an optional, domain-specific extension to the core library. This reduces the cognitive load for developers who only need the generic k-d tree functionality. The use of build tags could even be considered to make this optional at compile time.
