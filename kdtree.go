package poindexter

import (
	"errors"
	"math"
	"sort"
	"time"
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
	// ErrBackendUnavailable indicates that a requested backend cannot be used (e.g., not built/tagged).
	ErrBackendUnavailable = errors.New("kdtree: requested backend unavailable")
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

// CosineDistance implements 1 - cosine similarity.
//
// Distance is defined as 1 - (a·b)/(||a||*||b||). If both vectors are zero,
// distance is 0. If exactly one is zero, distance is 1. Numerical results are
// clamped to [0,2].
// Note: For typical normalized/weighted feature vectors with non-negative entries,
// the value will be in [0,1]. Opposite vectors in general spaces can yield up to 2.
type CosineDistance struct{}

func (CosineDistance) Distance(a, b []float64) float64 {
	var dot, na2, nb2 float64
	for i := range a {
		ai := a[i]
		bi := b[i]
		dot += ai * bi
		na2 += ai * ai
		nb2 += bi * bi
	}
	if na2 == 0 && nb2 == 0 {
		return 0
	}
	if na2 == 0 || nb2 == 0 {
		return 1
	}
	den := math.Sqrt(na2) * math.Sqrt(nb2)
	if den == 0 { // guard, though covered above
		return 1
	}
	cos := dot / den
	if cos > 1 {
		cos = 1
	} else if cos < -1 {
		cos = -1
	}
	d := 1 - cos
	if d < 0 {
		return 0
	}
	if d > 2 {
		return 2
	}
	return d
}

// WeightedCosineDistance implements 1 - weighted cosine similarity, where weights
// scale each axis in both the dot product and the norms.
// If Weights is nil or has zero length, this reduces to CosineDistance.
type WeightedCosineDistance struct{ Weights []float64 }

func (wcd WeightedCosineDistance) Distance(a, b []float64) float64 {
	w := wcd.Weights
	if len(w) == 0 || len(w) != len(a) || len(a) != len(b) {
		// Fallback to unweighted cosine when lengths mismatch or weights missing.
		return CosineDistance{}.Distance(a, b)
	}
	var dot, na2, nb2 float64
	for i := range a {
		wi := w[i]
		ai := a[i]
		bi := b[i]
		v := wi * ai
		dot += v * bi         // wi*ai*bi
		na2 += v * ai         // wi*ai*ai
		nb2 += (wi * bi) * bi // wi*bi*bi
	}
	if na2 == 0 && nb2 == 0 {
		return 0
	}
	if na2 == 0 || nb2 == 0 {
		return 1
	}
	den := math.Sqrt(na2) * math.Sqrt(nb2)
	if den == 0 {
		return 1
	}
	cos := dot / den
	if cos > 1 {
		cos = 1
	} else if cos < -1 {
		cos = -1
	}
	d := 1 - cos
	if d < 0 {
		return 0
	}
	if d > 2 {
		return 2
	}
	return d
}

// KDOption configures KDTree construction (non-generic to allow inference).
type KDOption func(*kdOptions)

type kdOptions struct {
	metric  DistanceMetric
	backend KDBackend
}

// defaultBackend returns the implicit backend depending on build tags.
// If built with the "gonum" tag, prefer the Gonum backend by default to keep
// code paths simple and performant; otherwise fall back to the linear backend.
func defaultBackend() KDBackend {
	if hasGonum() {
		return BackendGonum
	}
	return BackendLinear
}

// KDBackend selects the internal engine used by KDTree.
type KDBackend string

const (
	BackendLinear KDBackend = "linear"
	BackendGonum  KDBackend = "gonum"
)

// WithMetric sets the distance metric for the KDTree.
func WithMetric(m DistanceMetric) KDOption { return func(o *kdOptions) { o.metric = m } }

// WithBackend selects the internal KDTree backend ("linear" or "gonum").
// Default is linear. If the requested backend is unavailable (e.g., gonum build tag not enabled),
// the constructor will silently fall back to the linear backend.
func WithBackend(b KDBackend) KDOption { return func(o *kdOptions) { o.backend = b } }

// KDTree is a lightweight wrapper providing nearest-neighbor operations.
//
// Complexity: queries are O(n) linear scans in the current implementation.
// Inserts are O(1) amortized; deletes by ID are O(1) using swap-delete (order not preserved).
// Concurrency: KDTree is not safe for concurrent mutation. Guard with a mutex or
// share immutable snapshots for read-mostly workloads.
//
// This type is designed to be easily swappable with gonum.org/v1/gonum/spatial/kdtree
// in the future without breaking the public API.
type KDTree[T any] struct {
	points      []KDPoint[T]
	dim         int
	metric      DistanceMetric
	idIndex     map[string]int
	backend     KDBackend
	backendData any // opaque handle for backend-specific structures (e.g., gonum tree)

	// Analytics tracking (optional, enabled by default)
	analytics     *TreeAnalytics
	peerAnalytics *PeerAnalytics
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
	cfg := kdOptions{metric: EuclideanDistance{}, backend: defaultBackend()}
	for _, o := range opts {
		o(&cfg)
	}
	backend := cfg.backend
	var backendData any
	// Attempt to build gonum backend if requested and available.
	if backend == BackendGonum && hasGonum() {
		if bd, err := buildGonumBackend(pts, cfg.metric); err == nil {
			backendData = bd
		} else {
			backend = BackendLinear // fallback gracefully
		}
	} else if backend == BackendGonum && !hasGonum() {
		backend = BackendLinear // tag not enabled → fallback
	}
	t := &KDTree[T]{
		points:        append([]KDPoint[T](nil), pts...),
		dim:           dim,
		metric:        cfg.metric,
		idIndex:       idIndex,
		backend:       backend,
		backendData:   backendData,
		analytics:     NewTreeAnalytics(),
		peerAnalytics: NewPeerAnalytics(),
	}
	return t, nil
}

// NewKDTreeFromDim constructs an empty KDTree with the specified dimension.
// Call Insert to add points after construction.
func NewKDTreeFromDim[T any](dim int, opts ...KDOption) (*KDTree[T], error) {
	if dim <= 0 {
		return nil, ErrZeroDim
	}
	cfg := kdOptions{metric: EuclideanDistance{}, backend: defaultBackend()}
	for _, o := range opts {
		o(&cfg)
	}
	backend := cfg.backend
	if backend == BackendGonum && !hasGonum() {
		backend = BackendLinear
	}
	return &KDTree[T]{
		points:        nil,
		dim:           dim,
		metric:        cfg.metric,
		idIndex:       make(map[string]int),
		backend:       backend,
		backendData:   nil,
		analytics:     NewTreeAnalytics(),
		peerAnalytics: NewPeerAnalytics(),
	}, nil
}

// Dim returns the number of dimensions.
func (t *KDTree[T]) Dim() int { return t.dim }

// Len returns the number of points in the tree.
func (t *KDTree[T]) Len() int { return len(t.points) }

// Nearest returns the closest point to the query, along with its distance.
// ok is false if the tree is empty or the query dimensionality does not match Dim().
func (t *KDTree[T]) Nearest(query []float64) (KDPoint[T], float64, bool) {
	if len(query) != t.dim || t.Len() == 0 {
		return KDPoint[T]{}, 0, false
	}
	start := time.Now()
	defer func() {
		if t.analytics != nil {
			t.analytics.RecordQuery(time.Since(start).Nanoseconds())
		}
	}()

	// Gonum backend (if available and built)
	if t.backend == BackendGonum && t.backendData != nil {
		if idx, dist, ok := gonumNearest[T](t.backendData, query); ok && idx >= 0 && idx < len(t.points) {
			p := t.points[idx]
			if t.peerAnalytics != nil {
				t.peerAnalytics.RecordSelection(p.ID, dist)
			}
			return p, dist, true
		}
		// fall through to linear scan if backend didn't return a result
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
	p := t.points[bestIdx]
	if t.peerAnalytics != nil {
		t.peerAnalytics.RecordSelection(p.ID, bestDist)
	}
	return p, bestDist, true
}

// KNearest returns up to k nearest neighbors to the query in ascending distance order.
// If multiple points are at the same distance, tie ordering is arbitrary and not stable between calls.
func (t *KDTree[T]) KNearest(query []float64, k int) ([]KDPoint[T], []float64) {
	if k <= 0 || len(query) != t.dim || t.Len() == 0 {
		return nil, nil
	}
	start := time.Now()
	defer func() {
		if t.analytics != nil {
			t.analytics.RecordQuery(time.Since(start).Nanoseconds())
		}
	}()

	// Gonum backend path
	if t.backend == BackendGonum && t.backendData != nil {
		idxs, dists := gonumKNearest[T](t.backendData, query, k)
		if len(idxs) > 0 {
			neighbors := make([]KDPoint[T], len(idxs))
			for i := range idxs {
				neighbors[i] = t.points[idxs[i]]
				if t.peerAnalytics != nil {
					t.peerAnalytics.RecordSelection(neighbors[i].ID, dists[i])
				}
			}
			return neighbors, dists
		}
		// fall back on unexpected empty
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
		if t.peerAnalytics != nil {
			t.peerAnalytics.RecordSelection(neighbors[i].ID, dists[i])
		}
	}
	return neighbors, dists
}

// Radius returns points within radius r (inclusive) from the query, sorted by distance.
func (t *KDTree[T]) Radius(query []float64, r float64) ([]KDPoint[T], []float64) {
	if r < 0 || len(query) != t.dim || t.Len() == 0 {
		return nil, nil
	}
	start := time.Now()
	defer func() {
		if t.analytics != nil {
			t.analytics.RecordQuery(time.Since(start).Nanoseconds())
		}
	}()

	// Gonum backend path
	if t.backend == BackendGonum && t.backendData != nil {
		idxs, dists := gonumRadius[T](t.backendData, query, r)
		if len(idxs) > 0 {
			neighbors := make([]KDPoint[T], len(idxs))
			for i := range idxs {
				neighbors[i] = t.points[idxs[i]]
				if t.peerAnalytics != nil {
					t.peerAnalytics.RecordSelection(neighbors[i].ID, dists[i])
				}
			}
			return neighbors, dists
		}
		// fall back if no results
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
		if t.peerAnalytics != nil {
			t.peerAnalytics.RecordSelection(neighbors[i].ID, dists[i])
		}
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
	// Record insert in analytics
	if t.analytics != nil {
		t.analytics.RecordInsert()
	}
	// Rebuild backend if using Gonum
	if t.backend == BackendGonum && hasGonum() {
		if bd, err := buildGonumBackend(t.points, t.metric); err == nil {
			t.backendData = bd
			if t.analytics != nil {
				t.analytics.RecordRebuild()
			}
		} else {
			// fallback to linear if rebuild fails
			t.backend = BackendLinear
			t.backendData = nil
		}
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
	// Record delete in analytics
	if t.analytics != nil {
		t.analytics.RecordDelete()
	}
	// Rebuild backend if using Gonum
	if t.backend == BackendGonum && hasGonum() {
		if bd, err := buildGonumBackend(t.points, t.metric); err == nil {
			t.backendData = bd
			if t.analytics != nil {
				t.analytics.RecordRebuild()
			}
		} else {
			// fallback to linear if rebuild fails
			t.backend = BackendLinear
			t.backendData = nil
		}
	}
	return true
}

// Analytics returns the tree analytics tracker.
// Returns nil if analytics tracking is disabled.
func (t *KDTree[T]) Analytics() *TreeAnalytics {
	return t.analytics
}

// PeerAnalytics returns the peer analytics tracker.
// Returns nil if peer analytics tracking is disabled.
func (t *KDTree[T]) PeerAnalytics() *PeerAnalytics {
	return t.peerAnalytics
}

// GetAnalyticsSnapshot returns a point-in-time snapshot of tree analytics.
func (t *KDTree[T]) GetAnalyticsSnapshot() TreeAnalyticsSnapshot {
	if t.analytics == nil {
		return TreeAnalyticsSnapshot{}
	}
	return t.analytics.Snapshot()
}

// GetPeerStats returns per-peer selection statistics.
func (t *KDTree[T]) GetPeerStats() []PeerStats {
	if t.peerAnalytics == nil {
		return nil
	}
	return t.peerAnalytics.GetAllPeerStats()
}

// GetTopPeers returns the top N most frequently selected peers.
func (t *KDTree[T]) GetTopPeers(n int) []PeerStats {
	if t.peerAnalytics == nil {
		return nil
	}
	return t.peerAnalytics.GetTopPeers(n)
}

// ComputeDistanceDistribution analyzes the distribution of current point coordinates.
func (t *KDTree[T]) ComputeDistanceDistribution(axisNames []string) []AxisDistribution {
	return ComputeAxisDistributions(t.points, axisNames)
}

// ResetAnalytics clears all analytics data.
func (t *KDTree[T]) ResetAnalytics() {
	if t.analytics != nil {
		t.analytics.Reset()
	}
	if t.peerAnalytics != nil {
		t.peerAnalytics.Reset()
	}
}

// Points returns a copy of all points in the tree.
// This is useful for analytics and export operations.
func (t *KDTree[T]) Points() []KDPoint[T] {
	result := make([]KDPoint[T], len(t.points))
	copy(result, t.points)
	return result
}

// Backend returns the active backend type.
func (t *KDTree[T]) Backend() KDBackend {
	return t.backend
}
