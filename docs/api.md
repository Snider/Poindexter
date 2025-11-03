# API Reference

Complete API documentation for the Poindexter library.

## Core Functions

### Version

```go
func Version() string
```

Returns the current version of the library.

**Returns:**
- `string`: The version string (e.g., "0.2.1")

**Example:**

```go
version := poindexter.Version()
fmt.Println(version) // Output: 0.2.1
```

---

### Hello

```go
func Hello(name string) string
```

Returns a greeting message.

**Parameters:**
- `name` (string): The name to greet. If empty, defaults to "World"

**Returns:**
- `string`: A greeting message

**Examples:**

```go
// Greet the world
message := poindexter.Hello("")
fmt.Println(message) // Output: Hello, World!

// Greet a specific person
message = poindexter.Hello("Alice")
fmt.Println(message) // Output: Hello, Alice!
```

---

## Sorting Functions

### Basic Sorting

#### SortInts

```go
func SortInts(data []int)
```

Sorts a slice of integers in ascending order in place.

**Example:**

```go
numbers := []int{3, 1, 4, 1, 5, 9}
poindexter.SortInts(numbers)
fmt.Println(numbers) // Output: [1 1 3 4 5 9]
```

---

#### SortIntsDescending

```go
func SortIntsDescending(data []int)
```

Sorts a slice of integers in descending order in place.

**Example:**

```go
numbers := []int{3, 1, 4, 1, 5, 9}
poindexter.SortIntsDescending(numbers)
fmt.Println(numbers) // Output: [9 5 4 3 1 1]
```

---

#### SortStrings

```go
func SortStrings(data []string)
```

Sorts a slice of strings in ascending order in place.

**Example:**

```go
words := []string{"banana", "apple", "cherry"}
poindexter.SortStrings(words)
fmt.Println(words) // Output: [apple banana cherry]
```

---

#### SortStringsDescending

```go
func SortStringsDescending(data []string)
```

Sorts a slice of strings in descending order in place.

---

#### SortFloat64s

```go
func SortFloat64s(data []float64)
```

Sorts a slice of float64 values in ascending order in place.

---

#### SortFloat64sDescending

```go
func SortFloat64sDescending(data []float64)
```

Sorts a slice of float64 values in descending order in place.

---

### Advanced Sorting

#### SortBy

```go
func SortBy[T any](data []T, less func(i, j int) bool)
```

Sorts a slice using a custom comparison function.

**Parameters:**
- `data`: The slice to sort
- `less`: A function that returns true if data[i] should come before data[j]

**Example:**

```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
}

// Sort by age
poindexter.SortBy(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
// Result: [Bob(25) Alice(30) Charlie(35)]
```

---

#### SortByKey

```go
func SortByKey[T any, K int | float64 | string](data []T, key func(T) K)
```

Sorts a slice by extracting a comparable key from each element in ascending order.

**Parameters:**
- `data`: The slice to sort
- `key`: A function that extracts a sortable key from each element

**Example:**

```go
type Product struct {
    Name  string
    Price float64
}

products := []Product{
    {"Apple", 1.50},
    {"Banana", 0.75},
    {"Cherry", 3.00},
}

// Sort by price
poindexter.SortByKey(products, func(p Product) float64 {
    return p.Price
})
// Result: [Banana(0.75) Apple(1.50) Cherry(3.00)]
```

---

#### SortByKeyDescending

```go
func SortByKeyDescending[T any, K int | float64 | string](data []T, key func(T) K)
```

Sorts a slice by extracting a comparable key from each element in descending order.

**Example:**

```go
type Student struct {
    Name  string
    Score int
}

students := []Student{
    {"Alice", 85},
    {"Bob", 92},
    {"Charlie", 78},
}

// Sort by score descending
poindexter.SortByKeyDescending(students, func(s Student) int {
    return s.Score
})
// Result: [Bob(92) Alice(85) Charlie(78)]
```

---

### Checking if Sorted

#### IsSorted

```go
func IsSorted(data []int) bool
```

Checks if a slice of integers is sorted in ascending order.

---

#### IsSortedStrings

```go
func IsSortedStrings(data []string) bool
```

Checks if a slice of strings is sorted in ascending order.

---

#### IsSortedFloat64s

```go
func IsSortedFloat64s(data []float64) bool
```

Checks if a slice of float64 values is sorted in ascending order.

---

### Binary Search

#### BinarySearch

```go
func BinarySearch(data []int, target int) int
```

Performs a binary search on a sorted slice of integers.

**Parameters:**
- `data`: A sorted slice of integers
- `target`: The value to search for

**Returns:**
- `int`: The index where target is found, or -1 if not found

**Example:**

```go
numbers := []int{1, 3, 5, 7, 9, 11}
index := poindexter.BinarySearch(numbers, 7)
fmt.Println(index) // Output: 3
```

---

#### BinarySearchStrings

```go
func BinarySearchStrings(data []string, target string) int
```

Performs a binary search on a sorted slice of strings.

**Parameters:**
- `data`: A sorted slice of strings
- `target`: The value to search for

**Returns:**
- `int`: The index where target is found, or -1 if not found


## KDTree Helpers

Poindexter provides helpers to build normalized, weighted KD points from your own records. These functions min–max normalize each axis over your dataset, optionally invert axes where higher is better (to turn them into “lower cost”), and apply per‑axis weights.

```go
func Build2D[T any](
    items []T,
    id func(T) string,
    f1, f2 func(T) float64,
    weights [2]float64,
    invert [2]bool,
) ([]KDPoint[T], error)

func Build3D[T any](
    items []T,
    id func(T) string,
    f1, f2, f3 func(T) float64,
    weights [3]float64,
    invert [3]bool,
) ([]KDPoint[T], error)

func Build4D[T any](
    items []T,
    id func(T) string,
    f1, f2, f3, f4 func(T) float64,
    weights [4]float64,
    invert [4]bool,
) ([]KDPoint[T], error)
```

Example (4D over ping, hops, geo, score):

```go
// weights and inversion: flip score so higher is better → lower cost
weights := [4]float64{1.0, 0.7, 0.2, 1.2}
invert  := [4]bool{false, false, false, true}

pts, err := poindexter.Build4D(
    peers,
    func(p Peer) string { return p.ID },
    func(p Peer) float64 { return p.PingMS },
    func(p Peer) float64 { return p.Hops },
    func(p Peer) float64 { return p.GeoKM },
    func(p Peer) float64 { return p.Score },
    weights, invert,
)
if err != nil { panic(err) }

kdt, _ := poindexter.NewKDTree(pts, poindexter.WithMetric(poindexter.EuclideanDistance{}))
best, dist, _ := kdt.Nearest([]float64{0, 0, 0, 0})
```

Notes:
- Keep and reuse your normalization parameters (min/max) if you need consistency across updates; otherwise rebuild points when the candidate set changes.
- Use `invert` to turn “higher is better” features (like scores) into lower costs for distance calculations.


---

## KDTree Constructors and Errors

### NewKDTree

```go
func NewKDTree[T any](pts []KDPoint[T], opts ...KDOption) (*KDTree[T], error)
```

Build a KDTree from the provided points. All points must have the same dimensionality (> 0) and IDs (if provided) must be unique.

Possible errors:
- `ErrEmptyPoints`: no points provided
- `ErrZeroDim`: dimension must be at least 1
- `ErrDimMismatch`: inconsistent dimensionality among points
- `ErrDuplicateID`: duplicate point ID encountered

### NewKDTreeFromDim

```go
func NewKDTreeFromDim[T any](dim int, opts ...KDOption) (*KDTree[T], error)
```

Construct an empty KDTree with the given dimension, then populate later via `Insert`.

---

## KDTree Notes: Complexity, Ties, Concurrency

- Complexity: current implementation uses O(n) linear scans for queries (`Nearest`, `KNearest`, `Radius`). Inserts are O(1) amortized. Deletes by ID are O(1) using swap-delete (order not preserved).
- Tie ordering: when multiple neighbors have the same distance, ordering of ties is arbitrary and not stable between calls.
- Concurrency: KDTree is not safe for concurrent mutation. Wrap with a mutex or share immutable snapshots for read-mostly workloads.

See runnable examples in the repository `examples/` and the docs pages for 1D DHT and multi-dimensional KDTree usage.


## KDTree Normalization Stats (reuse across updates)

To keep normalization consistent across dynamic updates, compute per‑axis min/max once and reuse it to build points later. This avoids drift when the candidate set changes.

### Types

```go
// AxisStats holds the min/max observed for a single axis.
type AxisStats struct {
    Min float64
    Max float64
}

// NormStats holds per‑axis normalisation stats; for D dims, Stats has length D.
type NormStats struct {
    Stats []AxisStats
}
```

### Compute normalization stats

```go
func ComputeNormStats2D[T any](items []T, f1, f2 func(T) float64) NormStats
func ComputeNormStats3D[T any](items []T, f1, f2, f3 func(T) float64) NormStats
func ComputeNormStats4D[T any](items []T, f1, f2, f3, f4 func(T) float64) NormStats
```

### Build with precomputed stats

```go
func Build2DWithStats[T any](
    items []T,
    id func(T) string,
    f1, f2 func(T) float64,
    weights [2]float64,
    invert [2]bool,
    stats NormStats,
) ([]KDPoint[T], error)

func Build3DWithStats[T any](
    items []T,
    id func(T) string,
    f1, f2, f3 func(T) float64,
    weights [3]float64,
    invert [3]bool,
    stats NormStats,
) ([]KDPoint[T], error)

func Build4DWithStats[T any](
    items []T,
    id func(T) string,
    f1, f2, f3, f4 func(T) float64,
    weights [4]float64,
    invert [4]bool,
    stats NormStats,
) ([]KDPoint[T], error)
```

#### Example (2D)
```go
// Compute stats once over your baseline set
stats := poindexter.ComputeNormStats2D(peers,
    func(p Peer) float64 { return p.PingMS },
    func(p Peer) float64 { return p.Hops },
)

// Build points using those stats (now or later)
pts, _ := poindexter.Build2DWithStats(
    peers,
    func(p Peer) string { return p.ID },
    func(p Peer) float64 { return p.PingMS },
    func(p Peer) float64 { return p.Hops },
    [2]float64{1,1}, [2]bool{false,false}, stats,
)
```

Notes:
- If `min==max` for an axis, normalized value is `0` for that axis.
- `invert[i]` flips the normalized axis as `1 - n` before applying `weights[i]`.
- These helpers mirror `Build2D/3D/4D`, but use your provided `NormStats` instead of recomputing from the items slice.
