package poindexter

import (
	"errors"
	"math"
	"sort"
)

var (
	// ErrEmptyPoints indicates that no points were provided to build a KDTree.
	ErrEmptyPoints = errors.New("kdtree: no points provided")
	// ErrZeroDim indicates that points or tree dimension must be at least 1.
	ErrZeroDim = errors.New("kdtree: points must have at least one dimension")
	// ErrDimMismatch indicates inconsistent dimensionality among points.
	ErrDimMismatch = errors.New("kdtree: inconsistent dimensionality in points")
	// ErrDuplicateID indicates a duplicate point ID was encountered.
	ErrDuplicateID = errors.New("kdtree: duplicate point ID")
)

// KDPoint represents a point with coordinates and an attached payload/value.
// ID should be unique within a tree to enable O(1) deletes by ID.
// Coords must all have the same dimensionality within a given KDTree.
type KDPoint[T any] struct {
	ID     string
	Coords []float64
	Value  T
}

// DistanceMetric defines a metric over R^n.
type DistanceMetric interface {
	Distance(a, b []float64) float64
}

// EuclideanDistance implements the L2 metric.
type EuclideanDistance struct{}

func (EuclideanDistance) Distance(a, b []float64) float64 {
	var sum float64
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

// ManhattanDistance implements the L1 metric.
type ManhattanDistance struct{}

func (ManhattanDistance) Distance(a, b []float64) float64 {
	var sum float64
	for i := range a {
		d := a[i] - b[i]
		if d < 0 {
			d = -d
		}
		sum += d
	}
	return sum
}

// ChebyshevDistance implements the L-infinity (max) metric.
type ChebyshevDistance struct{}

func (ChebyshevDistance) Distance(a, b []float64) float64 {
	var max float64
	for i := range a {
		d := a[i] - b[i]
		if d < 0 {
			d = -d
		}
		if d > max {
			max = d
		}
	}
	return max
}

// KDOption configures KDTree construction (non-generic to allow inference).
type KDOption func(*kdOptions)

type kdOptions struct {
	metric DistanceMetric
}

// WithMetric sets the distance metric for the KDTree.
func WithMetric(m DistanceMetric) KDOption { return func(o *kdOptions) { o.metric = m } }

// KDTree is a lightweight wrapper providing nearest-neighbor operations.
// Note: This implementation currently uses linear scans for queries
// and is designed to be easily swappable with gonum.org/v1/gonum/spatial/kdtree
// in the future without breaking the public API.
type KDTree[T any] struct {
	points  []KDPoint[T]
	dim     int
	metric  DistanceMetric
	idIndex map[string]int
}

// NewKDTree builds a KDTree from the given points.
// All points must have the same dimensionality (>0).
func NewKDTree[T any](pts []KDPoint[T], opts ...KDOption) (*KDTree[T], error) {
	if len(pts) == 0 {
		return nil, ErrEmptyPoints
	}
	dim := len(pts[0].Coords)
	if dim == 0 {
		return nil, ErrZeroDim
	}
	idIndex := make(map[string]int, len(pts))
	for i, p := range pts {
		if len(p.Coords) != dim {
			return nil, ErrDimMismatch
		}
		if p.ID != "" {
			if _, exists := idIndex[p.ID]; exists {
				return nil, ErrDuplicateID
			}
			idIndex[p.ID] = i
		}
	}
	cfg := kdOptions{metric: EuclideanDistance{}}
	for _, o := range opts {
		o(&cfg)
	}
	t := &KDTree[T]{
		points:  append([]KDPoint[T](nil), pts...),
		dim:     dim,
		metric:  cfg.metric,
		idIndex: idIndex,
	}
	return t, nil
}

// NewKDTreeFromDim constructs an empty KDTree with the specified dimension.
// Call Insert to add points after construction.
func NewKDTreeFromDim[T any](dim int, opts ...KDOption) (*KDTree[T], error) {
	if dim <= 0 {
		return nil, ErrZeroDim
	}
	cfg := kdOptions{metric: EuclideanDistance{}}
	for _, o := range opts {
		o(&cfg)
	}
	return &KDTree[T]{
		points:  nil,
		dim:     dim,
		metric:  cfg.metric,
		idIndex: make(map[string]int),
	}, nil
}

// Dim returns the number of dimensions.
func (t *KDTree[T]) Dim() int { return t.dim }

// Len returns the number of points in the tree.
func (t *KDTree[T]) Len() int { return len(t.points) }

// Nearest returns the closest point to the query, along with its distance.
// ok is false if the tree is empty.
func (t *KDTree[T]) Nearest(query []float64) (KDPoint[T], float64, bool) {
	if len(query) != t.dim || t.Len() == 0 {
		return KDPoint[T]{}, 0, false
	}
	bestIdx := -1
	bestDist := math.MaxFloat64
	for i := range t.points {
		d := t.metric.Distance(query, t.points[i].Coords)
		if d < bestDist {
			bestDist = d
			bestIdx = i
		}
	}
	if bestIdx < 0 {
		return KDPoint[T]{}, 0, false
	}
	return t.points[bestIdx], bestDist, true
}

// KNearest returns up to k nearest neighbors to the query in ascending distance order.
func (t *KDTree[T]) KNearest(query []float64, k int) ([]KDPoint[T], []float64) {
	if k <= 0 || len(query) != t.dim || t.Len() == 0 {
		return nil, nil
	}
	tmp := make([]struct {
		idx  int
		dist float64
	}, len(t.points))
	for i := range t.points {
		tmp[i].idx = i
		tmp[i].dist = t.metric.Distance(query, t.points[i].Coords)
	}
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].dist < tmp[j].dist })
	if k > len(tmp) {
		k = len(tmp)
	}
	neighbors := make([]KDPoint[T], k)
	dists := make([]float64, k)
	for i := 0; i < k; i++ {
		neighbors[i] = t.points[tmp[i].idx]
		dists[i] = tmp[i].dist
	}
	return neighbors, dists
}

// Radius returns points within radius r (inclusive) from the query, sorted by distance.
func (t *KDTree[T]) Radius(query []float64, r float64) ([]KDPoint[T], []float64) {
	if r < 0 || len(query) != t.dim || t.Len() == 0 {
		return nil, nil
	}
	var sel []struct {
		idx  int
		dist float64
	}
	for i := range t.points {
		d := t.metric.Distance(query, t.points[i].Coords)
		if d <= r {
			sel = append(sel, struct {
				idx  int
				dist float64
			}{i, d})
		}
	}
	sort.Slice(sel, func(i, j int) bool { return sel[i].dist < sel[j].dist })
	neighbors := make([]KDPoint[T], len(sel))
	dists := make([]float64, len(sel))
	for i := range sel {
		neighbors[i] = t.points[sel[i].idx]
		dists[i] = sel[i].dist
	}
	return neighbors, dists
}

// Insert adds a point. Returns false if dimensionality mismatch or duplicate ID exists.
func (t *KDTree[T]) Insert(p KDPoint[T]) bool {
	if len(p.Coords) != t.dim {
		return false
	}
	if p.ID != "" {
		if _, exists := t.idIndex[p.ID]; exists {
			return false
		}
		// will set after append
	}
	t.points = append(t.points, p)
	if p.ID != "" {
		t.idIndex[p.ID] = len(t.points) - 1
	}
	return true
}

// DeleteByID removes a point by its ID. Returns false if not found or ID empty.
func (t *KDTree[T]) DeleteByID(id string) bool {
	if id == "" {
		return false
	}
	idx, ok := t.idIndex[id]
	if !ok {
		return false
	}
	last := len(t.points) - 1
	// swap delete
	t.points[idx] = t.points[last]
	if t.points[idx].ID != "" {
		t.idIndex[t.points[idx].ID] = idx
	}
	t.points = t.points[:last]
	delete(t.idIndex, id)
	return true
}
