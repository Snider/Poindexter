package poindexter

import (
	"math"
	"testing"
	"time"
)

// ============================================================================
// TreeAnalytics Tests
// ============================================================================

func TestNewTreeAnalytics(t *testing.T) {
	a := NewTreeAnalytics()
	if a == nil {
		t.Fatal("NewTreeAnalytics returned nil")
	}
	if a.QueryCount.Load() != 0 {
		t.Errorf("expected QueryCount=0, got %d", a.QueryCount.Load())
	}
	if a.InsertCount.Load() != 0 {
		t.Errorf("expected InsertCount=0, got %d", a.InsertCount.Load())
	}
	if a.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestTreeAnalyticsRecordQuery(t *testing.T) {
	a := NewTreeAnalytics()

	a.RecordQuery(1000) // 1μs
	a.RecordQuery(2000) // 2μs
	a.RecordQuery(500)  // 0.5μs

	if a.QueryCount.Load() != 3 {
		t.Errorf("expected QueryCount=3, got %d", a.QueryCount.Load())
	}
	if a.TotalQueryTimeNs.Load() != 3500 {
		t.Errorf("expected TotalQueryTimeNs=3500, got %d", a.TotalQueryTimeNs.Load())
	}
	if a.MinQueryTimeNs.Load() != 500 {
		t.Errorf("expected MinQueryTimeNs=500, got %d", a.MinQueryTimeNs.Load())
	}
	if a.MaxQueryTimeNs.Load() != 2000 {
		t.Errorf("expected MaxQueryTimeNs=2000, got %d", a.MaxQueryTimeNs.Load())
	}
	if a.LastQueryTimeNs.Load() != 500 {
		t.Errorf("expected LastQueryTimeNs=500, got %d", a.LastQueryTimeNs.Load())
	}
}

func TestTreeAnalyticsSnapshot(t *testing.T) {
	a := NewTreeAnalytics()

	a.RecordQuery(1000)
	a.RecordQuery(3000)
	a.RecordInsert()
	a.RecordInsert()
	a.RecordDelete()
	a.RecordRebuild()

	snap := a.Snapshot()

	if snap.QueryCount != 2 {
		t.Errorf("expected QueryCount=2, got %d", snap.QueryCount)
	}
	if snap.InsertCount != 2 {
		t.Errorf("expected InsertCount=2, got %d", snap.InsertCount)
	}
	if snap.DeleteCount != 1 {
		t.Errorf("expected DeleteCount=1, got %d", snap.DeleteCount)
	}
	if snap.AvgQueryTimeNs != 2000 {
		t.Errorf("expected AvgQueryTimeNs=2000, got %d", snap.AvgQueryTimeNs)
	}
	if snap.MinQueryTimeNs != 1000 {
		t.Errorf("expected MinQueryTimeNs=1000, got %d", snap.MinQueryTimeNs)
	}
	if snap.MaxQueryTimeNs != 3000 {
		t.Errorf("expected MaxQueryTimeNs=3000, got %d", snap.MaxQueryTimeNs)
	}
	if snap.BackendRebuildCnt != 1 {
		t.Errorf("expected BackendRebuildCnt=1, got %d", snap.BackendRebuildCnt)
	}
}

func TestTreeAnalyticsReset(t *testing.T) {
	a := NewTreeAnalytics()

	a.RecordQuery(1000)
	a.RecordInsert()
	a.RecordDelete()

	a.Reset()

	if a.QueryCount.Load() != 0 {
		t.Errorf("expected QueryCount=0 after reset, got %d", a.QueryCount.Load())
	}
	if a.InsertCount.Load() != 0 {
		t.Errorf("expected InsertCount=0 after reset, got %d", a.InsertCount.Load())
	}
	if a.DeleteCount.Load() != 0 {
		t.Errorf("expected DeleteCount=0 after reset, got %d", a.DeleteCount.Load())
	}
}

// ============================================================================
// PeerAnalytics Tests
// ============================================================================

func TestNewPeerAnalytics(t *testing.T) {
	p := NewPeerAnalytics()
	if p == nil {
		t.Fatal("NewPeerAnalytics returned nil")
	}
}

func TestPeerAnalyticsRecordSelection(t *testing.T) {
	p := NewPeerAnalytics()

	p.RecordSelection("peer1", 0.5)
	p.RecordSelection("peer1", 0.3)
	p.RecordSelection("peer2", 1.0)

	stats := p.GetPeerStats("peer1")
	if stats.SelectionCount != 2 {
		t.Errorf("expected peer1 SelectionCount=2, got %d", stats.SelectionCount)
	}
	if math.Abs(stats.AvgDistance-0.4) > 0.001 {
		t.Errorf("expected peer1 AvgDistance~0.4, got %f", stats.AvgDistance)
	}

	stats2 := p.GetPeerStats("peer2")
	if stats2.SelectionCount != 1 {
		t.Errorf("expected peer2 SelectionCount=1, got %d", stats2.SelectionCount)
	}
}

func TestPeerAnalyticsGetAllPeerStats(t *testing.T) {
	p := NewPeerAnalytics()

	p.RecordSelection("peer1", 0.5)
	p.RecordSelection("peer1", 0.5)
	p.RecordSelection("peer2", 1.0)
	p.RecordSelection("peer3", 0.8)
	p.RecordSelection("peer3", 0.8)
	p.RecordSelection("peer3", 0.8)

	all := p.GetAllPeerStats()
	if len(all) != 3 {
		t.Errorf("expected 3 peers, got %d", len(all))
	}

	// Should be sorted by selection count descending
	if all[0].PeerID != "peer3" || all[0].SelectionCount != 3 {
		t.Errorf("expected first peer to be peer3 with count=3, got %s with count=%d",
			all[0].PeerID, all[0].SelectionCount)
	}
}

func TestPeerAnalyticsGetTopPeers(t *testing.T) {
	p := NewPeerAnalytics()

	for i := 0; i < 5; i++ {
		p.RecordSelection("peer1", 0.5)
	}
	for i := 0; i < 3; i++ {
		p.RecordSelection("peer2", 0.3)
	}
	p.RecordSelection("peer3", 0.1)

	top := p.GetTopPeers(2)
	if len(top) != 2 {
		t.Errorf("expected 2 top peers, got %d", len(top))
	}
	if top[0].PeerID != "peer1" {
		t.Errorf("expected top peer to be peer1, got %s", top[0].PeerID)
	}
	if top[1].PeerID != "peer2" {
		t.Errorf("expected second peer to be peer2, got %s", top[1].PeerID)
	}
}

func TestPeerAnalyticsReset(t *testing.T) {
	p := NewPeerAnalytics()

	p.RecordSelection("peer1", 0.5)
	p.Reset()

	stats := p.GetAllPeerStats()
	if len(stats) != 0 {
		t.Errorf("expected 0 peers after reset, got %d", len(stats))
	}
}

// ============================================================================
// DistributionStats Tests
// ============================================================================

func TestComputeDistributionStatsEmpty(t *testing.T) {
	stats := ComputeDistributionStats(nil)
	if stats.Count != 0 {
		t.Errorf("expected Count=0 for empty input, got %d", stats.Count)
	}
}

func TestComputeDistributionStatsSingle(t *testing.T) {
	stats := ComputeDistributionStats([]float64{5.0})
	if stats.Count != 1 {
		t.Errorf("expected Count=1, got %d", stats.Count)
	}
	if stats.Min != 5.0 || stats.Max != 5.0 {
		t.Errorf("expected Min=Max=5.0, got Min=%f, Max=%f", stats.Min, stats.Max)
	}
	if stats.Mean != 5.0 {
		t.Errorf("expected Mean=5.0, got %f", stats.Mean)
	}
	if stats.Median != 5.0 {
		t.Errorf("expected Median=5.0, got %f", stats.Median)
	}
}

func TestComputeDistributionStatsMultiple(t *testing.T) {
	// Values: 1, 2, 3, 4, 5 - mean=3, median=3
	stats := ComputeDistributionStats([]float64{1, 2, 3, 4, 5})

	if stats.Count != 5 {
		t.Errorf("expected Count=5, got %d", stats.Count)
	}
	if stats.Min != 1.0 {
		t.Errorf("expected Min=1.0, got %f", stats.Min)
	}
	if stats.Max != 5.0 {
		t.Errorf("expected Max=5.0, got %f", stats.Max)
	}
	if stats.Mean != 3.0 {
		t.Errorf("expected Mean=3.0, got %f", stats.Mean)
	}
	if stats.Median != 3.0 {
		t.Errorf("expected Median=3.0, got %f", stats.Median)
	}
	// Variance = 2.0 for this dataset
	if math.Abs(stats.Variance-2.0) > 0.001 {
		t.Errorf("expected Variance~2.0, got %f", stats.Variance)
	}
}

func TestComputeDistributionStatsPercentiles(t *testing.T) {
	// 100 values from 0 to 99
	values := make([]float64, 100)
	for i := 0; i < 100; i++ {
		values[i] = float64(i)
	}
	stats := ComputeDistributionStats(values)

	// P25 should be around 24.75, P75 around 74.25
	if math.Abs(stats.P25-24.75) > 0.1 {
		t.Errorf("expected P25~24.75, got %f", stats.P25)
	}
	if math.Abs(stats.P75-74.25) > 0.1 {
		t.Errorf("expected P75~74.25, got %f", stats.P75)
	}
	if math.Abs(stats.P90-89.1) > 0.1 {
		t.Errorf("expected P90~89.1, got %f", stats.P90)
	}
}

// ============================================================================
// AxisDistribution Tests
// ============================================================================

func TestComputeAxisDistributions(t *testing.T) {
	points := []KDPoint[string]{
		{ID: "a", Coords: []float64{1.0, 10.0}},
		{ID: "b", Coords: []float64{2.0, 20.0}},
		{ID: "c", Coords: []float64{3.0, 30.0}},
	}

	dists := ComputeAxisDistributions(points, []string{"x", "y"})

	if len(dists) != 2 {
		t.Errorf("expected 2 axis distributions, got %d", len(dists))
	}

	if dists[0].Axis != 0 || dists[0].Name != "x" {
		t.Errorf("expected first axis=0, name=x, got axis=%d, name=%s", dists[0].Axis, dists[0].Name)
	}
	if dists[0].Stats.Mean != 2.0 {
		t.Errorf("expected axis 0 mean=2.0, got %f", dists[0].Stats.Mean)
	}

	if dists[1].Axis != 1 || dists[1].Name != "y" {
		t.Errorf("expected second axis=1, name=y, got axis=%d, name=%s", dists[1].Axis, dists[1].Name)
	}
	if dists[1].Stats.Mean != 20.0 {
		t.Errorf("expected axis 1 mean=20.0, got %f", dists[1].Stats.Mean)
	}
}

// ============================================================================
// NAT Routing Tests
// ============================================================================

func TestPeerQualityScoreDefaults(t *testing.T) {
	// Perfect peer
	perfect := NATRoutingMetrics{
		ConnectivityScore: 1.0,
		SymmetryScore:     1.0,
		RelayProbability:  0.0,
		DirectSuccessRate: 1.0,
		AvgRTTMs:          10,
		JitterMs:          5,
		PacketLossRate:    0.0,
		BandwidthMbps:     100,
		NATType:           string(NATTypeOpen),
	}
	score := PeerQualityScore(perfect, nil)
	if score < 0.9 {
		t.Errorf("expected perfect peer score > 0.9, got %f", score)
	}

	// Poor peer
	poor := NATRoutingMetrics{
		ConnectivityScore: 0.2,
		SymmetryScore:     0.1,
		RelayProbability:  0.9,
		DirectSuccessRate: 0.1,
		AvgRTTMs:          500,
		JitterMs:          100,
		PacketLossRate:    0.5,
		BandwidthMbps:     1,
		NATType:           string(NATTypeSymmetric),
	}
	poorScore := PeerQualityScore(poor, nil)
	if poorScore > 0.5 {
		t.Errorf("expected poor peer score < 0.5, got %f", poorScore)
	}
	if poorScore >= score {
		t.Error("poor peer should have lower score than perfect peer")
	}
}

func TestPeerQualityScoreCustomWeights(t *testing.T) {
	metrics := NATRoutingMetrics{
		ConnectivityScore: 1.0,
		SymmetryScore:     0.5,
		RelayProbability:  0.0,
		DirectSuccessRate: 1.0,
		AvgRTTMs:          100,
		JitterMs:          10,
		PacketLossRate:    0.01,
		BandwidthMbps:     50,
		NATType:           string(NATTypeFullCone),
	}

	// Weight latency heavily
	latencyWeights := QualityWeights{
		Latency:       10.0,
		Jitter:        1.0,
		PacketLoss:    1.0,
		Bandwidth:     1.0,
		Connectivity:  1.0,
		Symmetry:      1.0,
		DirectSuccess: 1.0,
		RelayPenalty:  1.0,
		NATType:       1.0,
	}
	scoreLatency := PeerQualityScore(metrics, &latencyWeights)

	// Weight bandwidth heavily
	bandwidthWeights := QualityWeights{
		Latency:       1.0,
		Jitter:        1.0,
		PacketLoss:    1.0,
		Bandwidth:     10.0,
		Connectivity:  1.0,
		Symmetry:      1.0,
		DirectSuccess: 1.0,
		RelayPenalty:  1.0,
		NATType:       1.0,
	}
	scoreBandwidth := PeerQualityScore(metrics, &bandwidthWeights)

	// Scores should differ based on weights
	if scoreLatency == scoreBandwidth {
		t.Error("different weights should produce different scores")
	}
}

func TestDefaultQualityWeights(t *testing.T) {
	w := DefaultQualityWeights()
	if w.Latency <= 0 {
		t.Error("Latency weight should be positive")
	}
	if w.Total() <= 0 {
		t.Error("Total weights should be positive")
	}
}

func TestNatTypeScore(t *testing.T) {
	tests := []struct {
		natType  string
		minScore float64
		maxScore float64
	}{
		{string(NATTypeOpen), 0.9, 1.0},
		{string(NATTypeFullCone), 0.8, 1.0},
		{string(NATTypeSymmetric), 0.2, 0.4},
		{string(NATTypeRelayRequired), 0.0, 0.1},
		{"unknown", 0.3, 0.5},
	}

	for _, tc := range tests {
		score := natTypeScore(tc.natType)
		if score < tc.minScore || score > tc.maxScore {
			t.Errorf("natType %s: expected score in [%f, %f], got %f",
				tc.natType, tc.minScore, tc.maxScore, score)
		}
	}
}

// ============================================================================
// Trust Score Tests
// ============================================================================

func TestComputeTrustScoreNewPeer(t *testing.T) {
	// New peer with no history
	metrics := TrustMetrics{
		SuccessfulTransactions: 0,
		FailedTransactions:     0,
		AgeSeconds:             86400, // 1 day old
	}
	score := ComputeTrustScore(metrics)
	// New peer should get moderate trust
	if score < 0.4 || score > 0.7 {
		t.Errorf("expected new peer score in [0.4, 0.7], got %f", score)
	}
}

func TestComputeTrustScoreGoodPeer(t *testing.T) {
	metrics := TrustMetrics{
		SuccessfulTransactions: 100,
		FailedTransactions:     2,
		AgeSeconds:             86400 * 30, // 30 days
		VouchCount:             5,
		FlagCount:              0,
		LastSuccessAt:          time.Now(),
	}
	score := ComputeTrustScore(metrics)
	if score < 0.8 {
		t.Errorf("expected good peer score > 0.8, got %f", score)
	}
}

func TestComputeTrustScoreBadPeer(t *testing.T) {
	metrics := TrustMetrics{
		SuccessfulTransactions: 5,
		FailedTransactions:     20,
		AgeSeconds:             86400,
		VouchCount:             0,
		FlagCount:              10,
	}
	score := ComputeTrustScore(metrics)
	if score > 0.3 {
		t.Errorf("expected bad peer score < 0.3, got %f", score)
	}
}

// ============================================================================
// Feature Normalization Tests
// ============================================================================

func TestStandardPeerFeaturesToSlice(t *testing.T) {
	features := StandardPeerFeatures{
		LatencyMs:       100,
		HopCount:        5,
		GeoDistanceKm:   1000,
		TrustScore:      0.9,
		BandwidthMbps:   50,
		PacketLossRate:  0.01,
		ConnectivityPct: 95,
		NATScore:        0.8,
	}

	slice := features.ToFeatureSlice()
	if len(slice) != 8 {
		t.Errorf("expected 8 features, got %d", len(slice))
	}

	// TrustScore should be inverted (0.9 -> 0.1)
	if math.Abs(slice[3]-0.1) > 0.001 {
		t.Errorf("expected inverted trust score ~0.1, got %f", slice[3])
	}
}

func TestNormalizePeerFeatures(t *testing.T) {
	features := []float64{100, 5, 1000, 0.5, 50, 0.01, 50, 0.5}
	ranges := DefaultPeerFeatureRanges()

	normalized := NormalizePeerFeatures(features, ranges)

	for i, v := range normalized {
		if v < 0 || v > 1 {
			t.Errorf("normalized feature %d out of range [0,1]: %f", i, v)
		}
	}
}

func TestWeightedPeerFeatures(t *testing.T) {
	normalized := []float64{0.5, 0.5, 0.5, 0.5}
	weights := []float64{1.0, 2.0, 0.5, 1.5}

	weighted := WeightedPeerFeatures(normalized, weights)

	expected := []float64{0.5, 1.0, 0.25, 0.75}
	for i, v := range weighted {
		if math.Abs(v-expected[i]) > 0.001 {
			t.Errorf("weighted feature %d: expected %f, got %f", i, expected[i], v)
		}
	}
}

func TestStandardFeatureLabels(t *testing.T) {
	labels := StandardFeatureLabels()
	if len(labels) != 8 {
		t.Errorf("expected 8 feature labels, got %d", len(labels))
	}
}

// ============================================================================
// KDTree Analytics Integration Tests
// ============================================================================

func TestKDTreeAnalyticsIntegration(t *testing.T) {
	points := []KDPoint[string]{
		{ID: "a", Coords: []float64{0, 0}, Value: "A"},
		{ID: "b", Coords: []float64{1, 1}, Value: "B"},
		{ID: "c", Coords: []float64{2, 2}, Value: "C"},
	}
	tree, err := NewKDTree(points)
	if err != nil {
		t.Fatal(err)
	}

	// Check initial analytics
	if tree.Analytics() == nil {
		t.Fatal("Analytics should not be nil")
	}
	if tree.PeerAnalytics() == nil {
		t.Fatal("PeerAnalytics should not be nil")
	}

	// Perform queries
	tree.Nearest([]float64{0.1, 0.1})
	tree.Nearest([]float64{0.9, 0.9})
	tree.KNearest([]float64{0.5, 0.5}, 2)

	snap := tree.GetAnalyticsSnapshot()
	if snap.QueryCount != 3 {
		t.Errorf("expected QueryCount=3, got %d", snap.QueryCount)
	}
	if snap.InsertCount != 0 {
		t.Errorf("expected InsertCount=0, got %d", snap.InsertCount)
	}

	// Check peer stats
	peerStats := tree.GetPeerStats()
	if len(peerStats) == 0 {
		t.Error("expected some peer stats after queries")
	}

	// Peer 'a' should have been selected for query [0.1, 0.1]
	var foundA bool
	for _, ps := range peerStats {
		if ps.PeerID == "a" && ps.SelectionCount > 0 {
			foundA = true
			break
		}
	}
	if !foundA {
		t.Error("expected peer 'a' to be recorded in analytics")
	}

	// Test top peers
	topPeers := tree.GetTopPeers(1)
	if len(topPeers) != 1 {
		t.Errorf("expected 1 top peer, got %d", len(topPeers))
	}

	// Test insert analytics
	tree.Insert(KDPoint[string]{ID: "d", Coords: []float64{3, 3}, Value: "D"})
	snap = tree.GetAnalyticsSnapshot()
	if snap.InsertCount != 1 {
		t.Errorf("expected InsertCount=1, got %d", snap.InsertCount)
	}

	// Test delete analytics
	tree.DeleteByID("d")
	snap = tree.GetAnalyticsSnapshot()
	if snap.DeleteCount != 1 {
		t.Errorf("expected DeleteCount=1, got %d", snap.DeleteCount)
	}

	// Test reset
	tree.ResetAnalytics()
	snap = tree.GetAnalyticsSnapshot()
	if snap.QueryCount != 0 || snap.InsertCount != 0 || snap.DeleteCount != 0 {
		t.Error("expected all counts to be 0 after reset")
	}
}

func TestKDTreeDistanceDistribution(t *testing.T) {
	points := []KDPoint[string]{
		{ID: "a", Coords: []float64{0, 10}, Value: "A"},
		{ID: "b", Coords: []float64{1, 20}, Value: "B"},
		{ID: "c", Coords: []float64{2, 30}, Value: "C"},
	}
	tree, _ := NewKDTree(points)

	dists := tree.ComputeDistanceDistribution([]string{"x", "y"})
	if len(dists) != 2 {
		t.Errorf("expected 2 axis distributions, got %d", len(dists))
	}

	if dists[0].Name != "x" || dists[0].Stats.Mean != 1.0 {
		t.Errorf("unexpected axis 0 distribution: name=%s, mean=%f",
			dists[0].Name, dists[0].Stats.Mean)
	}
	if dists[1].Name != "y" || dists[1].Stats.Mean != 20.0 {
		t.Errorf("unexpected axis 1 distribution: name=%s, mean=%f",
			dists[1].Name, dists[1].Stats.Mean)
	}
}

func TestKDTreePointsExport(t *testing.T) {
	points := []KDPoint[string]{
		{ID: "a", Coords: []float64{0, 0}, Value: "A"},
		{ID: "b", Coords: []float64{1, 1}, Value: "B"},
	}
	tree, _ := NewKDTree(points)

	exported := tree.Points()
	if len(exported) != 2 {
		t.Errorf("expected 2 points, got %d", len(exported))
	}

	// Verify it's a copy, not a reference
	exported[0].ID = "modified"
	original := tree.Points()
	if original[0].ID == "modified" {
		t.Error("Points() should return a copy, not a reference")
	}
}

func TestKDTreeBackend(t *testing.T) {
	tree, _ := NewKDTreeFromDim[string](2)
	backend := tree.Backend()
	if backend != BackendLinear && backend != BackendGonum {
		t.Errorf("unexpected backend: %s", backend)
	}
}
