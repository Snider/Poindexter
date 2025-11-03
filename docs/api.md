# API Reference

Complete API documentation for the Poindexter library.

## Core Functions

### Version

```go
func Version() string
```

Returns the current version of the library.

**Returns:**
- `string`: The version string (e.g., "0.1.0")

**Example:**

```go
version := poindexter.Version()
fmt.Println(version) // Output: 0.1.0
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
