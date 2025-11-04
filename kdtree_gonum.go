//go:build gonum

package poindexter

import (
	"math"
	"sort"
)

// Note: This file is compiled when built with the "gonum" tag. For now, we
// provide an internal KD-tree backend that performs balanced median-split
// construction and branch-and-bound queries. This gives sub-linear behavior on
// suitable datasets without introducing an external dependency. The public API
// and option names remain the same; a future change can swap this implementation
// to use gonum.org/v1/gonum/spatial/kdtree without altering callers.

// hasGonum reports whether the optimized backend is compiled in.
func hasGonum() bool { return true }

// kdNode represents a node in the median-split KD-tree.
type kdNode struct {
	axis  int
	idx   int // index into the original points slice
	val   float64
	left  *kdNode
	right *kdNode
}

// kdBackend holds the KD-tree root and metadata.
type kdBackend struct {
	root   *kdNode
	dim    int
	metric DistanceMetric
	// Access to original coords by index is done via a closure we capture at build
	coords func(i int) []float64
	len    int
}

// buildGonumBackend builds a balanced KD-tree using variance-based axis choice
// and median splits. It does not reorder the external points slice; it keeps
// indices and accesses the original data via closures, preserving caller order.
func buildGonumBackend[T any](points []KDPoint[T], metric DistanceMetric) (any, error) {
	// Only enable this backend for metrics where the axis-slab bound is valid
	// for pruning: L2/L1/L∞. For other metrics (e.g., Cosine), fall back.
	switch metric.(type) {
	case EuclideanDistance, ManhattanDistance, ChebyshevDistance:
		// supported
	default:
		return nil, ErrBackendUnavailable
	}
	if len(points) == 0 {
		return &kdBackend{root: nil, dim: 0, metric: metric, coords: func(int) []float64 { return nil }}, nil
	}
	dim := len(points[0].Coords)
	coords := func(i int) []float64 { return points[i].Coords }
	idxs := make([]int, len(points))
	for i := range idxs {
		idxs[i] = i
	}
	root := buildKDRecursive(idxs, coords, dim, 0)
	return &kdBackend{root: root, dim: dim, metric: metric, coords: coords, len: len(points)}, nil
}

// compute per-axis standard deviation (used for axis selection)
func axisStd(idxs []int, coords func(int) []float64, dim int) []float64 {
	vars := make([]float64, dim)
	means := make([]float64, dim)
	n := float64(len(idxs))
	if n == 0 {
		return vars
	}
	for _, i := range idxs {
		c := coords(i)
		for d := 0; d < dim; d++ {
			means[d] += c[d]
		}
	}
	for d := 0; d < dim; d++ {
		means[d] /= n
	}
	for _, i := range idxs {
		c := coords(i)
		for d := 0; d < dim; d++ {
			delta := c[d] - means[d]
			vars[d] += delta * delta
		}
	}
	for d := 0; d < dim; d++ {
		vars[d] = math.Sqrt(vars[d] / n)
	}
	return vars
}

func buildKDRecursive(idxs []int, coords func(int) []float64, dim int, depth int) *kdNode {
	if len(idxs) == 0 {
		return nil
	}
	// choose axis with max stddev
	stds := axisStd(idxs, coords, dim)
	axis := 0
	maxv := stds[0]
	for d := 1; d < dim; d++ {
		if stds[d] > maxv {
			maxv = stds[d]
			axis = d
		}
	}
	// nth-element (partial sort) by axis using sort.Slice for simplicity
	sort.Slice(idxs, func(i, j int) bool { return coords(idxs[i])[axis] < coords(idxs[j])[axis] })
	mid := len(idxs) / 2
	medianIdx := idxs[mid]
	n := &kdNode{axis: axis, idx: medianIdx, val: coords(medianIdx)[axis]}
	n.left = buildKDRecursive(append([]int(nil), idxs[:mid]...), coords, dim, depth+1)
	n.right = buildKDRecursive(append([]int(nil), idxs[mid+1:]...), coords, dim, depth+1)
	return n
}

// gonumNearest performs 1-NN search using the KD backend.
func gonumNearest[T any](backend any, query []float64) (int, float64, bool) {
	b, ok := backend.(*kdBackend)
	if !ok || b.root == nil || len(query) != b.dim {
		return -1, 0, false
	}
	bestIdx := -1
	bestDist := math.MaxFloat64
	var search func(*kdNode)
	search = func(n *kdNode) {
		if n == nil {
			return
		}
		c := b.coords(n.idx)
		d := b.metric.Distance(query, c)
		if d < bestDist {
			bestDist = d
			bestIdx = n.idx
		}
		axis := n.axis
		qv := query[axis]
		// choose side
		near, far := n.left, n.right
		if qv >= n.val {
			near, far = n.right, n.left
		}
		search(near)
		// prune if hyperslab distance is >= bestDist
		diff := qv - n.val
		if diff < 0 {
			diff = -diff
		}
		if diff <= bestDist {
			search(far)
		}
	}
	search(b.root)
	if bestIdx < 0 {
		return -1, 0, false
	}
	return bestIdx, bestDist, true
}

// small max-heap for (distance, index)
// We’ll use a slice maintaining the largest distance at [0] via container/heap-like ops.
type knnItem struct {
	idx  int
	dist float64
}

type knnHeap []knnItem

func (h knnHeap) Len() int           { return len(h) }
func (h knnHeap) less(i, j int) bool { return h[i].dist > h[j].dist } // max-heap by dist
func (h *knnHeap) push(x knnItem)    { *h = append(*h, x); h.up(len(*h) - 1) }
func (h *knnHeap) pop() knnItem {
	n := len(*h) - 1
	h.swap(0, n)
	v := (*h)[n]
	*h = (*h)[:n]
	h.down(0)
	return v
}
func (h *knnHeap) peek() knnItem { return (*h)[0] }
func (h knnHeap) swap(i, j int)  { h[i], h[j] = h[j], h[i] }
func (h *knnHeap) up(i int) {
	for i > 0 {
		p := (i - 1) / 2
		if !h.less(i, p) {
			break
		}
		h.swap(i, p)
		i = p
	}
}
func (h *knnHeap) down(i int) {
	for {
		l := 2*i + 1
		r := l + 1
		largest := i
		if l < h.Len() && h.less(l, largest) {
			largest = l
		}
		if r < h.Len() && h.less(r, largest) {
			largest = r
		}
		if largest == i {
			break
		}
		h.swap(i, largest)
		i = largest
	}
}

// gonumKNearest returns indices in ascending distance order.
func gonumKNearest[T any](backend any, query []float64, k int) ([]int, []float64) {
	b, ok := backend.(*kdBackend)
	if !ok || b.root == nil || len(query) != b.dim || k <= 0 {
		return nil, nil
	}
	var h knnHeap
	bestCap := k
	var search func(*kdNode)
	search = func(n *kdNode) {
		if n == nil {
			return
		}
		c := b.coords(n.idx)
		d := b.metric.Distance(query, c)
		if h.Len() < bestCap {
			h.push(knnItem{idx: n.idx, dist: d})
		} else if d < h.peek().dist {
			// replace max
			h[0] = knnItem{idx: n.idx, dist: d}
			h.down(0)
		}
		axis := n.axis
		qv := query[axis]
		near, far := n.left, n.right
		if qv >= n.val {
			near, far = n.right, n.left
		}
		search(near)
		// prune against current worst in heap if heap is full; otherwise use bestDist
		threshold := math.MaxFloat64
		if h.Len() == bestCap {
			threshold = h.peek().dist
		} else if h.Len() > 0 {
			// use best known (not strictly necessary)
			threshold = h.peek().dist
		}
		diff := qv - n.val
		if diff < 0 {
			diff = -diff
		}
		if diff <= threshold {
			search(far)
		}
	}
	search(b.root)
	// Extract to slices and sort ascending by distance
	res := make([]knnItem, len(h))
	copy(res, h)
	sort.Slice(res, func(i, j int) bool { return res[i].dist < res[j].dist })
	idxs := make([]int, len(res))
	dists := make([]float64, len(res))
	for i := range res {
		idxs[i] = res[i].idx
		dists[i] = res[i].dist
	}
	return idxs, dists
}

func gonumRadius[T any](backend any, query []float64, r float64) ([]int, []float64) {
	b, ok := backend.(*kdBackend)
	if !ok || b.root == nil || len(query) != b.dim || r < 0 {
		return nil, nil
	}
	var res []knnItem
	var search func(*kdNode)
	search = func(n *kdNode) {
		if n == nil {
			return
		}
		c := b.coords(n.idx)
		d := b.metric.Distance(query, c)
		if d <= r {
			res = append(res, knnItem{idx: n.idx, dist: d})
		}
		axis := n.axis
		qv := query[axis]
		near, far := n.left, n.right
		if qv >= n.val {
			near, far = n.right, n.left
		}
		search(near)
		diff := qv - n.val
		if diff < 0 {
			diff = -diff
		}
		if diff <= r {
			search(far)
		}
	}
	search(b.root)
	sort.Slice(res, func(i, j int) bool { return res[i].dist < res[j].dist })
	idxs := make([]int, len(res))
	dists := make([]float64, len(res))
	for i := range res {
		idxs[i] = res[i].idx
		dists[i] = res[i].dist
	}
	return idxs, dists
}
