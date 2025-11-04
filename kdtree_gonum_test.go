//go:build gonum

package poindexter

import (
	"math"
	"testing"
)

func equalish(a, b []float64, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > tol {
			return false
		}
	}
	return true
}

func TestGonumKnnHeap(t *testing.T) {
	h := knnHeap{}

	h.push(knnItem{idx: 1, dist: 1.0})
	h.push(knnItem{idx: 2, dist: 2.0})
	h.push(knnItem{idx: 3, dist: 0.5})

	if h.Len() != 3 {
		t.Errorf("expected heap length 3, got %d", h.Len())
	}

	item := h.pop()
	if item.idx != 2 || item.dist != 2.0 {
		t.Errorf("expected item with index 2 and dist 2.0, got idx %d dist %f", item.idx, item.dist)
	}

	item = h.pop()
	if item.idx != 1 || item.dist != 1.0 {
		t.Errorf("expected item with index 1 and dist 1.0, got idx %d dist %f", item.idx, item.dist)
	}

	item = h.pop()
	if item.idx != 3 || item.dist != 0.5 {
		t.Errorf("expected item with index 3 and dist 0.5, got idx %d dist %f", item.idx, item.dist)
	}

	if h.Len() != 0 {
		t.Errorf("expected heap length 0, got %d", h.Len())
	}
}

func TestGonumNearest(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{2, 2}},
		{ID: "3", Coords: []float64{3, 3}},
	}

	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}

	p, dist, ok := tree.Nearest([]float64{1.1, 1.1})
	if !ok || p.ID != "1" || math.Abs(dist-0.1414213562373095) > 1e-9 {
		t.Errorf("expected point 1 with dist ~0.14, got point %s with dist %f", p.ID, dist)
	}
}

func TestGonumKNearest(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{2, 2}},
		{ID: "3", Coords: []float64{3, 3}},
	}

	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}

	ps, dists := tree.KNearest([]float64{1.1, 1.1}, 2)
	if len(ps) != 2 || ps[0].ID != "1" || ps[1].ID != "2" {
		t.Errorf("expected points 1 and 2, got %v", ps)
	}

	expectedDists := []float64{0.1414213562373095, 1.2727922061357854}
	if !equalish(dists, expectedDists, 1e-9) {
		t.Errorf("expected dists %v, got %v", expectedDists, dists)
	}
}

func TestGonumRadius(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{2, 2}},
		{ID: "3", Coords: []float64{3, 3}},
	}

	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}

	ps, dists := tree.Radius([]float64{1.1, 1.1}, 1.5)
	if len(ps) != 2 || ps[0].ID != "1" || ps[1].ID != "2" {
		t.Errorf("expected points 1 and 2, got %v", ps)
	}

	expectedDists := []float64{0.1414213562373095, 1.2727922061357854}
	if !equalish(dists, expectedDists, 1e-9) {
		t.Errorf("expected dists %v, got %v", expectedDists, dists)
	}
}

func TestBuildGonumBackendWithNonSupportedMetric(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum), WithMetric(CosineDistance{}))
	if err != nil {
		t.Fatal(err)
	}
	if tree.backend != BackendLinear {
		t.Errorf("expected fallback to linear backend, but got %s", tree.backend)
	}
}

func TestGonumNearestWithEmptyTree(t *testing.T) {
	_, err := NewKDTree([]KDPoint[int]{}, WithBackend(BackendGonum))
	if err != ErrEmptyPoints {
		t.Fatalf("expected ErrEmptyPoints, got %v", err)
	}
}

func TestAxisStdWithNoPoints(t *testing.T) {
	stds := axisStd(nil, nil, 2)
	if len(stds) != 2 || stds[0] != 0 || stds[1] != 0 {
		t.Errorf("expected [0, 0], got %v", stds)
	}
}

func TestGonumNearestWithNilRoot(t *testing.T) {
	backend := &kdBackend{root: nil, dim: 2}
	_, _, ok := gonumNearest[int](backend, []float64{1, 1})
	if ok {
		t.Error("expected no point found, but got one")
	}
}

func TestGonumNearestWithMismatchedDimensions(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}

	_, _, ok := tree.Nearest([]float64{1, 1, 1})
	if ok {
		t.Error("expected no point found, but got one")
	}
}

func TestGonumKNearestWithEmptyTree(t *testing.T) {
	_, err := NewKDTree([]KDPoint[int]{}, WithBackend(BackendGonum))
	if err != ErrEmptyPoints {
		t.Fatalf("expected ErrEmptyPoints, got %v", err)
	}
}

func TestGonumRadiusWithEmptyTree(t *testing.T) {
	_, err := NewKDTree([]KDPoint[int]{}, WithBackend(BackendGonum))
	if err != ErrEmptyPoints {
		t.Fatalf("expected ErrEmptyPoints, got %v", err)
	}
}

func TestGonumKNearestWithZeroK(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.KNearest([]float64{1, 1}, 0)
	if len(ps) != 0 {
		t.Error("expected 0 points, got some")
	}
}

func TestGonumRadiusWithNegativeRadius(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.Radius([]float64{1, 1}, -1)
	if len(ps) != 0 {
		t.Error("expected 0 points, got some")
	}
}

func TestGonumNearestWithSinglePoint(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, ok := tree.Nearest([]float64{1, 1})
	if !ok || p.ID != "1" {
		t.Errorf("expected point 1, got %v", p)
	}
}
func TestGonumKnnHeapPop(t *testing.T) {
	h := knnHeap{}
	h.push(knnItem{idx: 1, dist: 1.0})
	h.push(knnItem{idx: 2, dist: 2.0})
	h.push(knnItem{idx: 3, dist: 0.5})

	if h.Len() != 3 {
		t.Errorf("expected heap length 3, got %d", h.Len())
	}

	item := h.pop()
	if item.idx != 2 || item.dist != 2.0 {
		t.Errorf("expected item with index 2 and dist 2.0, got idx %d dist %f", item.idx, item.dist)
	}
}

func TestGonumKNearestWithSmallK(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{2, 2}},
		{ID: "3", Coords: []float64{3, 3}},
	}

	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}

	ps, _ := tree.KNearest([]float64{1.1, 1.1}, 1)
	if len(ps) != 1 || ps[0].ID != "1" {
		t.Errorf("expected point 1, got %v", ps)
	}
}
func TestGonumNearestReturnsFalseForNoPoints(t *testing.T) {
	tree, err := NewKDTreeFromDim[int](2, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	_, _, ok := tree.Nearest([]float64{0, 0})
	if ok {
		t.Errorf("expected ok to be false, but it was true")
	}
}
func TestBuildKDRecursiveWithSinglePoint(t *testing.T) {
	idxs := []int{0}
	coords := func(i int) []float64 { return []float64{1, 1} }
	node := buildKDRecursive(idxs, coords, 2, 0)
	if node == nil {
		t.Fatal("expected a node, got nil")
	}
	if node.idx != 0 {
		t.Errorf("expected index 0, got %d", node.idx)
	}
}

func TestGonumKNearestWithLargeK(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{2, 2}},
		{ID: "3", Coords: []float64{3, 3}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.KNearest([]float64{1.1, 1.1}, 5)
	if len(ps) != 3 {
		t.Errorf("expected 3 points, got %d", len(ps))
	}
}
func TestGonumNearestWithIdenticalPoints(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, _ := tree.Nearest([]float64{1, 1})
	if p.ID != "1" && p.ID != "2" {
		t.Errorf("expected point 1 or 2, got %v", p)
	}
}
func TestGonumKNearestPrefersCloserPoints(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{1.1, 1.1}},
		{ID: "3", Coords: []float64{1.2, 1.2}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.KNearest([]float64{0.9, 0.9}, 2)
	if len(ps) != 2 || ps[0].ID != "1" || ps[1].ID != "2" {
		t.Errorf("expected points 1 and 2, got %v", ps)
	}
}
func TestGonumKNearestWithFewerPointsThanK(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.KNearest([]float64{1, 1}, 2)
	if len(ps) != 1 {
		t.Errorf("expected 1 point, got %d", len(ps))
	}
}
func TestGonumKNearestReturnsCorrectOrder(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{3, 3}},
		{ID: "2", Coords: []float64{1, 1}},
		{ID: "3", Coords: []float64{2, 2}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.KNearest([]float64{0, 0}, 3)
	if ps[0].ID != "2" || ps[1].ID != "3" || ps[2].ID != "1" {
		t.Errorf("expected points in order 2, 3, 1, got %v", ps)
	}
}

func TestGonumRadiusReturnsAllWithinRadius(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{2, 2}},
		{ID: "3", Coords: []float64{10, 10}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.Radius([]float64{0, 0}, 3)
	if len(ps) != 2 {
		t.Errorf("expected 2 points, got %d", len(ps))
	}
}

func TestGonumRadiusReturnsEmptyForLargeRadiusWithNoPoints(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{10, 10}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.Radius([]float64{0, 0}, 1)
	if len(ps) != 0 {
		t.Errorf("expected 0 points, got %d", len(ps))
	}
}
func TestGonumNearestWithNegativeCoords(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{-1, -1}},
		{ID: "2", Coords: []float64{-2, -2}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, _ := tree.Nearest([]float64{-1.1, -1.1})
	if p.ID != "1" {
		t.Errorf("expected point 1, got %v", p)
	}
}
func TestGonumRadiusWithOverlappingPoints(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.Radius([]float64{1, 1}, 0.1)
	if len(ps) != 2 {
		t.Errorf("expected 2 points, got %d", len(ps))
	}
}

func TestGonumNearestWithFurtherPoints(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{10, 10}},
		{ID: "2", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, _ := tree.Nearest([]float64{0, 0})
	if p.ID != "2" {
		t.Errorf("expected point 2, got %v", p)
	}
}

func TestGonumNearestWithZeroDistance(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	_, dist, _ := tree.Nearest([]float64{1, 1})
	if dist != 0 {
		t.Errorf("expected distance 0, got %f", dist)
	}
}
func TestGonumNearestWithRightChild(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{5, 5}},
		{ID: "2", Coords: []float64{1, 1}},
		{ID: "3", Coords: []float64{8, 8}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, _ := tree.Nearest([]float64{9, 9})
	if p.ID != "3" {
		t.Errorf("expected point 3, got %v", p)
	}
}

func TestGonumKNearestHeapBehavior(t *testing.T) {
	h := knnHeap{}
	h.push(knnItem{idx: 1, dist: 1.0})
	h.push(knnItem{idx: 2, dist: 3.0})
	h.push(knnItem{idx: 3, dist: 2.0})

	if h.peek().dist != 3.0 {
		t.Errorf("expected max dist to be 3.0, got %f", h.peek().dist)
	}

	h.push(knnItem{idx: 4, dist: 0.5})
	if h.peek().dist != 3.0 {
		t.Errorf("expected max dist to be 3.0, got %f", h.peek().dist)
	}
}

func TestGonumNearestWithUnsortedPoints(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{10, 0}},
		{ID: "2", Coords: []float64{0, 10}},
		{ID: "3", Coords: []float64{5, 5}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, _ := tree.Nearest([]float64{4, 4})
	if p.ID != "3" {
		t.Errorf("expected point 3, got %v", p)
	}
}
func TestGonumKNearestWithThreshold(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{10, 10}},
		{ID: "3", Coords: []float64{2, 2}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.KNearest([]float64{0, 0}, 2)
	if len(ps) != 2 {
		t.Fatalf("expected 2 points, got %d", len(ps))
	}
	if ps[0].ID != "1" || ps[1].ID != "3" {
		t.Errorf("expected points 1 and 3, got %v and %v", ps[0].ID, ps[1].ID)
	}
}

func TestGonumRadiusSearch(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{1.5, 1.5}},
		{ID: "3", Coords: []float64{3, 3}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.Radius([]float64{1.2, 1.2}, 0.5)
	if len(ps) != 2 {
		t.Errorf("expected 2 points, got %d", len(ps))
	}
}
func TestGonumNearestWithFloatMinMax(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1e150, 1e150}},
		{ID: "2", Coords: []float64{-1e150, -1e150}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, ok := tree.Nearest([]float64{0, 0})
	if !ok {
		t.Fatal("expected to find a point, but didn't")
	}
	if p.ID != "1" && p.ID != "2" {
		t.Errorf("expected point 1 or 2, got %v", p)
	}
}
func TestGonumKnnHeapWithDuplicateDistances(t *testing.T) {
	h := knnHeap{}
	h.push(knnItem{idx: 1, dist: 1.0})
	h.push(knnItem{idx: 2, dist: 1.0})
	if h.Len() != 2 {
		t.Errorf("expected heap length 2, got %d", h.Len())
	}
}
func TestGonumNearestToPointOnAxis(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{0, 10}},
		{ID: "2", Coords: []float64{0, -10}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, _ := tree.Nearest([]float64{0, 0})
	if p.ID != "1" && p.ID != "2" {
		t.Errorf("expected point 1 or 2, got %v", p)
	}
}
func TestGonumNearestReturnsCorrectlyWhenPointsAreCollinear(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 1}},
		{ID: "2", Coords: []float64{2, 2}},
		{ID: "3", Coords: []float64{3, 3}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, _ := tree.Nearest([]float64{1.9, 1.9})
	if p.ID != "2" {
		t.Errorf("expected point 2, got %v", p)
	}
}
func TestGonumKNearestWithMorePoints(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{0, 0}},
		{ID: "2", Coords: []float64{1, 1}},
		{ID: "3", Coords: []float64{2, 2}},
		{ID: "4", Coords: []float64{3, 3}},
		{ID: "5", Coords: []float64{4, 4}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.KNearest([]float64{0.5, 0.5}, 3)
	if len(ps) != 3 {
		t.Fatalf("expected 3 points, got %d", len(ps))
	}
	if !((ps[0].ID == "1" && ps[1].ID == "2") || (ps[0].ID == "2" && ps[1].ID == "1")) {
		t.Errorf("expected first two points to be 1 and 2, got %s and %s", ps[0].ID, ps[1].ID)
	}
	if ps[2].ID != "3" {
		t.Errorf("expected third point to be 3, got %s", ps[2].ID)
	}
}
func TestGonumRadiusReturnsSorted(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1.2, 1.2}},
		{ID: "2", Coords: []float64{1.1, 1.1}},
		{ID: "3", Coords: []float64{1.0, 1.0}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	ps, _ := tree.Radius([]float64{0, 0}, 2)
	if len(ps) != 3 {
		t.Errorf("expected 3 points, got %d", len(ps))
	}
	if ps[0].ID != "3" || ps[1].ID != "2" || ps[2].ID != "1" {
		t.Errorf("expected order 3, 2, 1, got %v, %v, %v", ps[0].ID, ps[1].ID, ps[2].ID)
	}
}
func TestGonumNearestWithMultipleDimensions(t *testing.T) {
	points := []KDPoint[int]{
		{ID: "1", Coords: []float64{1, 2, 3, 4}},
		{ID: "2", Coords: []float64{5, 6, 7, 8}},
	}
	tree, err := NewKDTree(points, WithBackend(BackendGonum))
	if err != nil {
		t.Fatal(err)
	}
	p, _, _ := tree.Nearest([]float64{1.1, 2.1, 3.1, 4.1})
	if p.ID != "1" {
		t.Errorf("expected point 1, got %v", p)
	}
}
