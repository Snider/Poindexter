package poindexter

import (
	"reflect"
	"testing"
)

func TestSortInts(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"empty slice", []int{}, []int{}},
		{"single element", []int{5}, []int{5}},
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"unsorted", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"with negatives", []int{-3, 5, -1, 0, 2}, []int{-3, -1, 0, 2, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]int, len(tt.input))
			copy(data, tt.input)
			SortInts(data)
			if !reflect.DeepEqual(data, tt.expected) {
				t.Errorf("SortInts(%v) = %v, want %v", tt.input, data, tt.expected)
			}
		})
	}
}

func TestSortIntsDescending(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"empty slice", []int{}, []int{}},
		{"single element", []int{5}, []int{5}},
		{"ascending order", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"unsorted", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{9, 6, 5, 4, 3, 2, 1, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]int, len(tt.input))
			copy(data, tt.input)
			SortIntsDescending(data)
			if !reflect.DeepEqual(data, tt.expected) {
				t.Errorf("SortIntsDescending(%v) = %v, want %v", tt.input, data, tt.expected)
			}
		})
	}
}

func TestSortStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"empty slice", []string{}, []string{}},
		{"single element", []string{"hello"}, []string{"hello"}},
		{"already sorted", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"reverse sorted", []string{"z", "y", "x"}, []string{"x", "y", "z"}},
		{"unsorted", []string{"banana", "apple", "cherry", "date"}, []string{"apple", "banana", "cherry", "date"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]string, len(tt.input))
			copy(data, tt.input)
			SortStrings(data)
			if !reflect.DeepEqual(data, tt.expected) {
				t.Errorf("SortStrings(%v) = %v, want %v", tt.input, data, tt.expected)
			}
		})
	}
}

func TestSortStringsDescending(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"empty slice", []string{}, []string{}},
		{"ascending order", []string{"a", "b", "c"}, []string{"c", "b", "a"}},
		{"unsorted", []string{"banana", "apple", "cherry"}, []string{"cherry", "banana", "apple"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]string, len(tt.input))
			copy(data, tt.input)
			SortStringsDescending(data)
			if !reflect.DeepEqual(data, tt.expected) {
				t.Errorf("SortStringsDescending(%v) = %v, want %v", tt.input, data, tt.expected)
			}
		})
	}
}

func TestSortFloat64s(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected []float64
	}{
		{"empty slice", []float64{}, []float64{}},
		{"single element", []float64{3.14}, []float64{3.14}},
		{"unsorted", []float64{3.14, 1.41, 2.71, 1.73}, []float64{1.41, 1.73, 2.71, 3.14}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]float64, len(tt.input))
			copy(data, tt.input)
			SortFloat64s(data)
			if !reflect.DeepEqual(data, tt.expected) {
				t.Errorf("SortFloat64s(%v) = %v, want %v", tt.input, data, tt.expected)
			}
		})
	}
}

func TestSortFloat64sDescending(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected []float64
	}{
		{"empty slice", []float64{}, []float64{}},
		{"unsorted", []float64{3.14, 1.41, 2.71}, []float64{3.14, 2.71, 1.41}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]float64, len(tt.input))
			copy(data, tt.input)
			SortFloat64sDescending(data)
			if !reflect.DeepEqual(data, tt.expected) {
				t.Errorf("SortFloat64sDescending(%v) = %v, want %v", tt.input, data, tt.expected)
			}
		})
	}
}

func TestSortBy(t *testing.T) {
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
	SortBy(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})

	expected := []Person{
		{"Bob", 25},
		{"Alice", 30},
		{"Charlie", 35},
	}

	if !reflect.DeepEqual(people, expected) {
		t.Errorf("SortBy (by age) = %v, want %v", people, expected)
	}
}

func TestSortByKey(t *testing.T) {
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
	SortByKey(products, func(p Product) float64 {
		return p.Price
	})

	expected := []Product{
		{"Banana", 0.75},
		{"Apple", 1.50},
		{"Cherry", 3.00},
	}

	if !reflect.DeepEqual(products, expected) {
		t.Errorf("SortByKey (by price) = %v, want %v", products, expected)
	}
}

func TestSortByKeyDescending(t *testing.T) {
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
	SortByKeyDescending(students, func(s Student) int {
		return s.Score
	})

	expected := []Student{
		{"Bob", 92},
		{"Alice", 85},
		{"Charlie", 78},
	}

	if !reflect.DeepEqual(students, expected) {
		t.Errorf("SortByKeyDescending (by score) = %v, want %v", students, expected)
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{"empty slice", []int{}, true},
		{"single element", []int{5}, true},
		{"sorted ascending", []int{1, 2, 3, 4, 5}, true},
		{"not sorted", []int{1, 3, 2, 4}, false},
		{"sorted with duplicates", []int{1, 1, 2, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSorted(tt.input)
			if result != tt.expected {
				t.Errorf("IsSorted(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSortedStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected bool
	}{
		{"sorted", []string{"a", "b", "c"}, true},
		{"not sorted", []string{"b", "a", "c"}, false},
		{"empty", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSortedStrings(tt.input)
			if result != tt.expected {
				t.Errorf("IsSortedStrings(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSortedFloat64s(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected bool
	}{
		{"sorted", []float64{1.1, 2.2, 3.3}, true},
		{"not sorted", []float64{2.2, 1.1, 3.3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSortedFloat64s(tt.input)
			if result != tt.expected {
				t.Errorf("IsSortedFloat64s(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBinarySearch(t *testing.T) {
	data := []int{1, 3, 5, 7, 9, 11}

	tests := []struct {
		name     string
		target   int
		expected int
	}{
		{"found at start", 1, 0},
		{"found in middle", 5, 2},
		{"found at end", 11, 5},
		{"not found", 4, -1},
		{"not found - too small", 0, -1},
		{"not found - too large", 12, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BinarySearch(data, tt.target)
			if result != tt.expected {
				t.Errorf("BinarySearch(%v, %d) = %d, want %d", data, tt.target, result, tt.expected)
			}
		})
	}
}

func TestBinarySearchStrings(t *testing.T) {
	data := []string{"apple", "banana", "cherry", "date"}

	tests := []struct {
		name     string
		target   string
		expected int
	}{
		{"found at start", "apple", 0},
		{"found in middle", "banana", 1},
		{"found at end", "date", 3},
		{"not found", "grape", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BinarySearchStrings(data, tt.target)
			if result != tt.expected {
				t.Errorf("BinarySearchStrings(%v, %q) = %d, want %d", data, tt.target, result, tt.expected)
			}
		})
	}
}
