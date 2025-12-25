package poindexter

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// TreeAnalytics tracks operational statistics for a KDTree.
// All counters are safe for concurrent reads; use the Reset() method for atomic reset.
type TreeAnalytics struct {
	QueryCount  atomic.Int64 // Total nearest/kNearest/radius queries
	InsertCount atomic.Int64 // Total successful inserts
	DeleteCount atomic.Int64 // Total successful deletes

	// Timing statistics (nanoseconds)
	TotalQueryTimeNs  atomic.Int64
	LastQueryTimeNs   atomic.Int64
	MinQueryTimeNs    atomic.Int64
	MaxQueryTimeNs    atomic.Int64
	LastQueryAt       atomic.Int64 // Unix nanoseconds
	CreatedAt         time.Time
	LastRebuiltAt     atomic.Int64 // Unix nanoseconds (for gonum backend rebuilds)
	BackendRebuildCnt atomic.Int64 // Number of backend rebuilds
}

// NewTreeAnalytics creates a new analytics tracker.
func NewTreeAnalytics() *TreeAnalytics {
	a := &TreeAnalytics{
		CreatedAt: time.Now(),
	}
	a.MinQueryTimeNs.Store(math.MaxInt64)
	return a
}

// RecordQuery records a query operation with timing.
func (a *TreeAnalytics) RecordQuery(durationNs int64) {
	a.QueryCount.Add(1)
	a.TotalQueryTimeNs.Add(durationNs)
	a.LastQueryTimeNs.Store(durationNs)
	a.LastQueryAt.Store(time.Now().UnixNano())

	// Update min/max (best-effort, not strictly atomic)
	for {
		cur := a.MinQueryTimeNs.Load()
		if durationNs >= cur || a.MinQueryTimeNs.CompareAndSwap(cur, durationNs) {
			break
		}
	}
	for {
		cur := a.MaxQueryTimeNs.Load()
		if durationNs <= cur || a.MaxQueryTimeNs.CompareAndSwap(cur, durationNs) {
			break
		}
	}
}

// RecordInsert records a successful insert.
func (a *TreeAnalytics) RecordInsert() {
	a.InsertCount.Add(1)
}

// RecordDelete records a successful delete.
func (a *TreeAnalytics) RecordDelete() {
	a.DeleteCount.Add(1)
}

// RecordRebuild records a backend rebuild.
func (a *TreeAnalytics) RecordRebuild() {
	a.BackendRebuildCnt.Add(1)
	a.LastRebuiltAt.Store(time.Now().UnixNano())
}

// Snapshot returns a point-in-time view of the analytics.
func (a *TreeAnalytics) Snapshot() TreeAnalyticsSnapshot {
	avgNs := int64(0)
	qc := a.QueryCount.Load()
	if qc > 0 {
		avgNs = a.TotalQueryTimeNs.Load() / qc
	}
	minNs := a.MinQueryTimeNs.Load()
	if minNs == math.MaxInt64 {
		minNs = 0
	}
	return TreeAnalyticsSnapshot{
		QueryCount:        qc,
		InsertCount:       a.InsertCount.Load(),
		DeleteCount:       a.DeleteCount.Load(),
		AvgQueryTimeNs:    avgNs,
		MinQueryTimeNs:    minNs,
		MaxQueryTimeNs:    a.MaxQueryTimeNs.Load(),
		LastQueryTimeNs:   a.LastQueryTimeNs.Load(),
		LastQueryAt:       time.Unix(0, a.LastQueryAt.Load()),
		CreatedAt:         a.CreatedAt,
		BackendRebuildCnt: a.BackendRebuildCnt.Load(),
		LastRebuiltAt:     time.Unix(0, a.LastRebuiltAt.Load()),
	}
}

// Reset atomically resets all counters.
func (a *TreeAnalytics) Reset() {
	a.QueryCount.Store(0)
	a.InsertCount.Store(0)
	a.DeleteCount.Store(0)
	a.TotalQueryTimeNs.Store(0)
	a.LastQueryTimeNs.Store(0)
	a.MinQueryTimeNs.Store(math.MaxInt64)
	a.MaxQueryTimeNs.Store(0)
	a.LastQueryAt.Store(0)
	a.BackendRebuildCnt.Store(0)
	a.LastRebuiltAt.Store(0)
}

// TreeAnalyticsSnapshot is an immutable snapshot for JSON serialization.
type TreeAnalyticsSnapshot struct {
	QueryCount        int64     `json:"queryCount"`
	InsertCount       int64     `json:"insertCount"`
	DeleteCount       int64     `json:"deleteCount"`
	AvgQueryTimeNs    int64     `json:"avgQueryTimeNs"`
	MinQueryTimeNs    int64     `json:"minQueryTimeNs"`
	MaxQueryTimeNs    int64     `json:"maxQueryTimeNs"`
	LastQueryTimeNs   int64     `json:"lastQueryTimeNs"`
	LastQueryAt       time.Time `json:"lastQueryAt"`
	CreatedAt         time.Time `json:"createdAt"`
	BackendRebuildCnt int64     `json:"backendRebuildCount"`
	LastRebuiltAt     time.Time `json:"lastRebuiltAt"`
}

// PeerAnalytics tracks per-peer selection statistics for NAT routing optimization.
type PeerAnalytics struct {
	mu sync.RWMutex

	// Per-peer hit counters (peer ID -> selection count)
	hitCounts map[string]*atomic.Int64
	// Per-peer cumulative distance sums for average calculation
	distanceSums map[string]*atomic.Uint64 // stored as bits of float64
	// Last selection time per peer
	lastSelected map[string]*atomic.Int64 // Unix nano
}

// NewPeerAnalytics creates a new peer analytics tracker.
func NewPeerAnalytics() *PeerAnalytics {
	return &PeerAnalytics{
		hitCounts:    make(map[string]*atomic.Int64),
		distanceSums: make(map[string]*atomic.Uint64),
		lastSelected: make(map[string]*atomic.Int64),
	}
}

// RecordSelection records that a peer was selected/returned in a query result.
func (p *PeerAnalytics) RecordSelection(peerID string, distance float64) {
	if peerID == "" {
		return
	}

	p.mu.RLock()
	hc, hok := p.hitCounts[peerID]
	ds, dok := p.distanceSums[peerID]
	ls, lok := p.lastSelected[peerID]
	p.mu.RUnlock()

	if !hok || !dok || !lok {
		p.mu.Lock()
		if _, ok := p.hitCounts[peerID]; !ok {
			p.hitCounts[peerID] = &atomic.Int64{}
		}
		if _, ok := p.distanceSums[peerID]; !ok {
			p.distanceSums[peerID] = &atomic.Uint64{}
		}
		if _, ok := p.lastSelected[peerID]; !ok {
			p.lastSelected[peerID] = &atomic.Int64{}
		}
		hc = p.hitCounts[peerID]
		ds = p.distanceSums[peerID]
		ls = p.lastSelected[peerID]
		p.mu.Unlock()
	}

	hc.Add(1)
	// Atomic float add via CAS
	for {
		old := ds.Load()
		oldF := math.Float64frombits(old)
		newF := oldF + distance
		if ds.CompareAndSwap(old, math.Float64bits(newF)) {
			break
		}
	}
	ls.Store(time.Now().UnixNano())
}

// GetPeerStats returns statistics for a specific peer.
func (p *PeerAnalytics) GetPeerStats(peerID string) PeerStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	hc, hok := p.hitCounts[peerID]
	ds, dok := p.distanceSums[peerID]
	ls, lok := p.lastSelected[peerID]

	stats := PeerStats{PeerID: peerID}
	if hok {
		stats.SelectionCount = hc.Load()
	}
	if dok && stats.SelectionCount > 0 {
		stats.AvgDistance = math.Float64frombits(ds.Load()) / float64(stats.SelectionCount)
	}
	if lok {
		stats.LastSelectedAt = time.Unix(0, ls.Load())
	}
	return stats
}

// GetAllPeerStats returns statistics for all tracked peers.
func (p *PeerAnalytics) GetAllPeerStats() []PeerStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]PeerStats, 0, len(p.hitCounts))
	for id := range p.hitCounts {
		stats := PeerStats{PeerID: id}
		if hc := p.hitCounts[id]; hc != nil {
			stats.SelectionCount = hc.Load()
		}
		if ds := p.distanceSums[id]; ds != nil && stats.SelectionCount > 0 {
			stats.AvgDistance = math.Float64frombits(ds.Load()) / float64(stats.SelectionCount)
		}
		if ls := p.lastSelected[id]; ls != nil {
			stats.LastSelectedAt = time.Unix(0, ls.Load())
		}
		result = append(result, stats)
	}

	// Sort by selection count descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].SelectionCount > result[j].SelectionCount
	})
	return result
}

// GetTopPeers returns the top N most frequently selected peers.
func (p *PeerAnalytics) GetTopPeers(n int) []PeerStats {
	all := p.GetAllPeerStats()
	if n > len(all) {
		n = len(all)
	}
	return all[:n]
}

// Reset clears all peer analytics data.
func (p *PeerAnalytics) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.hitCounts = make(map[string]*atomic.Int64)
	p.distanceSums = make(map[string]*atomic.Uint64)
	p.lastSelected = make(map[string]*atomic.Int64)
}

// PeerStats holds statistics for a single peer.
type PeerStats struct {
	PeerID         string    `json:"peerId"`
	SelectionCount int64     `json:"selectionCount"`
	AvgDistance    float64   `json:"avgDistance"`
	LastSelectedAt time.Time `json:"lastSelectedAt"`
}

// DistributionStats provides statistical analysis of distances in query results.
type DistributionStats struct {
	Count      int       `json:"count"`
	Min        float64   `json:"min"`
	Max        float64   `json:"max"`
	Mean       float64   `json:"mean"`
	Median     float64   `json:"median"`
	StdDev     float64   `json:"stdDev"`
	P25        float64   `json:"p25"`  // 25th percentile
	P75        float64   `json:"p75"`  // 75th percentile
	P90        float64   `json:"p90"`  // 90th percentile
	P99        float64   `json:"p99"`  // 99th percentile
	Variance   float64   `json:"variance"`
	Skewness   float64   `json:"skewness"`
	SampleSize int       `json:"sampleSize"`
	ComputedAt time.Time `json:"computedAt"`
}

// ComputeDistributionStats calculates distribution statistics from a slice of distances.
func ComputeDistributionStats(distances []float64) DistributionStats {
	n := len(distances)
	if n == 0 {
		return DistributionStats{ComputedAt: time.Now()}
	}

	// Sort for percentile calculations
	sorted := make([]float64, n)
	copy(sorted, distances)
	sort.Float64s(sorted)

	// Basic stats
	min, max := sorted[0], sorted[n-1]
	sum := 0.0
	for _, d := range sorted {
		sum += d
	}
	mean := sum / float64(n)

	// Variance and standard deviation
	sumSqDiff := 0.0
	for _, d := range sorted {
		diff := d - mean
		sumSqDiff += diff * diff
	}
	variance := sumSqDiff / float64(n)
	stdDev := math.Sqrt(variance)

	// Skewness
	skewness := 0.0
	if stdDev > 0 {
		sumCubeDiff := 0.0
		for _, d := range sorted {
			diff := (d - mean) / stdDev
			sumCubeDiff += diff * diff * diff
		}
		skewness = sumCubeDiff / float64(n)
	}

	return DistributionStats{
		Count:      n,
		Min:        min,
		Max:        max,
		Mean:       mean,
		Median:     percentile(sorted, 0.5),
		StdDev:     stdDev,
		P25:        percentile(sorted, 0.25),
		P75:        percentile(sorted, 0.75),
		P90:        percentile(sorted, 0.90),
		P99:        percentile(sorted, 0.99),
		Variance:   variance,
		Skewness:   skewness,
		SampleSize: n,
		ComputedAt: time.Now(),
	}
}

// percentile returns the p-th percentile from a sorted slice.
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}
	idx := p * float64(len(sorted)-1)
	lower := int(idx)
	upper := lower + 1
	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	frac := idx - float64(lower)
	return sorted[lower]*(1-frac) + sorted[upper]*frac
}

// AxisDistribution provides per-axis (feature) distribution analysis.
type AxisDistribution struct {
	Axis  int               `json:"axis"`
	Name  string            `json:"name,omitempty"`
	Stats DistributionStats `json:"stats"`
}

// ComputeAxisDistributions analyzes the distribution of values along each axis.
func ComputeAxisDistributions[T any](points []KDPoint[T], axisNames []string) []AxisDistribution {
	if len(points) == 0 {
		return nil
	}
	dim := len(points[0].Coords)
	result := make([]AxisDistribution, dim)

	for axis := 0; axis < dim; axis++ {
		values := make([]float64, len(points))
		for i, p := range points {
			if axis < len(p.Coords) {
				values[i] = p.Coords[axis]
			}
		}
		name := ""
		if axis < len(axisNames) {
			name = axisNames[axis]
		}
		result[axis] = AxisDistribution{
			Axis:  axis,
			Name:  name,
			Stats: ComputeDistributionStats(values),
		}
	}
	return result
}

// NATRoutingMetrics provides metrics specifically for NAT traversal routing decisions.
type NATRoutingMetrics struct {
	// Connectivity score (0-1): higher means better reachability
	ConnectivityScore float64 `json:"connectivityScore"`
	// Symmetry score (0-1): higher means more symmetric NAT (easier to traverse)
	SymmetryScore float64 `json:"symmetryScore"`
	// Relay requirement probability (0-1): likelihood peer needs relay
	RelayProbability float64 `json:"relayProbability"`
	// Direct connection success rate (historical)
	DirectSuccessRate float64 `json:"directSuccessRate"`
	// Average RTT in milliseconds
	AvgRTTMs float64 `json:"avgRttMs"`
	// Jitter (RTT variance) in milliseconds
	JitterMs float64 `json:"jitterMs"`
	// Packet loss rate (0-1)
	PacketLossRate float64 `json:"packetLossRate"`
	// Bandwidth estimate in Mbps
	BandwidthMbps float64 `json:"bandwidthMbps"`
	// NAT type classification
	NATType string `json:"natType"`
	// Last probe timestamp
	LastProbeAt time.Time `json:"lastProbeAt"`
}

// NATTypeClassification enumerates common NAT types for routing decisions.
type NATTypeClassification string

const (
	NATTypeOpen            NATTypeClassification = "open"             // No NAT / Public IP
	NATTypeFullCone        NATTypeClassification = "full_cone"        // Easy to traverse
	NATTypeRestrictedCone  NATTypeClassification = "restricted_cone"  // Moderate difficulty
	NATTypePortRestricted  NATTypeClassification = "port_restricted"  // Harder to traverse
	NATTypeSymmetric       NATTypeClassification = "symmetric"        // Hardest to traverse
	NATTypeSymmetricUDP    NATTypeClassification = "symmetric_udp"    // UDP-only symmetric
	NATTypeUnknown         NATTypeClassification = "unknown"          // Not yet classified
	NATTypeBehindCGNAT     NATTypeClassification = "cgnat"            // Carrier-grade NAT
	NATTypeFirewalled      NATTypeClassification = "firewalled"       // Blocked by firewall
	NATTypeRelayRequired   NATTypeClassification = "relay_required"   // Must use relay
)

// PeerQualityScore computes a composite quality score for peer selection.
// Higher scores indicate better peers for routing.
// Weights can be customized; default weights emphasize latency and reliability.
func PeerQualityScore(metrics NATRoutingMetrics, weights *QualityWeights) float64 {
	w := DefaultQualityWeights()
	if weights != nil {
		w = *weights
	}

	// Normalize metrics to 0-1 scale (higher is better)
	latencyScore := 1.0 - math.Min(metrics.AvgRTTMs/1000.0, 1.0)         // <1000ms is acceptable
	jitterScore := 1.0 - math.Min(metrics.JitterMs/100.0, 1.0)           // <100ms jitter
	lossScore := 1.0 - metrics.PacketLossRate                            // 0 loss is best
	bandwidthScore := math.Min(metrics.BandwidthMbps/100.0, 1.0)         // 100Mbps is excellent
	connectivityScore := metrics.ConnectivityScore                        // Already 0-1
	symmetryScore := metrics.SymmetryScore                                // Already 0-1
	directScore := metrics.DirectSuccessRate                              // Already 0-1
	relayPenalty := 1.0 - metrics.RelayProbability                       // Prefer non-relay

	// NAT type bonus/penalty
	natScore := natTypeScore(metrics.NATType)

	// Weighted combination
	score := (w.Latency*latencyScore +
		w.Jitter*jitterScore +
		w.PacketLoss*lossScore +
		w.Bandwidth*bandwidthScore +
		w.Connectivity*connectivityScore +
		w.Symmetry*symmetryScore +
		w.DirectSuccess*directScore +
		w.RelayPenalty*relayPenalty +
		w.NATType*natScore) / w.Total()

	return math.Max(0, math.Min(1, score))
}

// QualityWeights configures the importance of each metric in peer selection.
type QualityWeights struct {
	Latency       float64 `json:"latency"`
	Jitter        float64 `json:"jitter"`
	PacketLoss    float64 `json:"packetLoss"`
	Bandwidth     float64 `json:"bandwidth"`
	Connectivity  float64 `json:"connectivity"`
	Symmetry      float64 `json:"symmetry"`
	DirectSuccess float64 `json:"directSuccess"`
	RelayPenalty  float64 `json:"relayPenalty"`
	NATType       float64 `json:"natType"`
}

// Total returns the sum of all weights for normalization.
func (w QualityWeights) Total() float64 {
	return w.Latency + w.Jitter + w.PacketLoss + w.Bandwidth +
		w.Connectivity + w.Symmetry + w.DirectSuccess + w.RelayPenalty + w.NATType
}

// DefaultQualityWeights returns sensible defaults for peer selection.
func DefaultQualityWeights() QualityWeights {
	return QualityWeights{
		Latency:       3.0, // Most important
		Jitter:        1.5,
		PacketLoss:    2.0,
		Bandwidth:     1.0,
		Connectivity:  2.0,
		Symmetry:      1.0,
		DirectSuccess: 2.0,
		RelayPenalty:  1.5,
		NATType:       1.0,
	}
}

// natTypeScore returns a 0-1 score based on NAT type (higher is better for routing).
func natTypeScore(natType string) float64 {
	switch NATTypeClassification(natType) {
	case NATTypeOpen:
		return 1.0
	case NATTypeFullCone:
		return 0.9
	case NATTypeRestrictedCone:
		return 0.7
	case NATTypePortRestricted:
		return 0.5
	case NATTypeSymmetric:
		return 0.3
	case NATTypeSymmetricUDP:
		return 0.25
	case NATTypeBehindCGNAT:
		return 0.2
	case NATTypeFirewalled:
		return 0.1
	case NATTypeRelayRequired:
		return 0.05
	default:
		return 0.4 // Unknown gets middle score
	}
}

// TrustMetrics tracks trust and reputation for peer selection.
type TrustMetrics struct {
	// ReputationScore (0-1): aggregated trust score
	ReputationScore float64 `json:"reputationScore"`
	// SuccessfulTransactions: count of successful exchanges
	SuccessfulTransactions int64 `json:"successfulTransactions"`
	// FailedTransactions: count of failed/aborted exchanges
	FailedTransactions int64 `json:"failedTransactions"`
	// AgeSeconds: how long this peer has been known
	AgeSeconds int64 `json:"ageSeconds"`
	// LastSuccessAt: last successful interaction
	LastSuccessAt time.Time `json:"lastSuccessAt"`
	// LastFailureAt: last failed interaction
	LastFailureAt time.Time `json:"lastFailureAt"`
	// VouchCount: number of other peers vouching for this peer
	VouchCount int `json:"vouchCount"`
	// FlagCount: number of reports against this peer
	FlagCount int `json:"flagCount"`
	// ProofOfWork: computational proof of stake/work
	ProofOfWork float64 `json:"proofOfWork"`
}

// ComputeTrustScore calculates a composite trust score from trust metrics.
func ComputeTrustScore(t TrustMetrics) float64 {
	total := t.SuccessfulTransactions + t.FailedTransactions
	if total == 0 {
		// New peer with no history: moderate trust with age bonus
		ageBonus := math.Min(float64(t.AgeSeconds)/(86400*30), 0.2) // Up to 0.2 for 30 days
		return 0.5 + ageBonus
	}

	// Base score from success rate
	successRate := float64(t.SuccessfulTransactions) / float64(total)

	// Volume confidence (more transactions = more confident)
	volumeConfidence := 1 - 1/(1+float64(total)/10)

	// Vouch/flag adjustment
	vouchBonus := math.Min(float64(t.VouchCount)*0.02, 0.15)
	flagPenalty := math.Min(float64(t.FlagCount)*0.05, 0.3)

	// Recency bonus (recent success = better)
	recencyBonus := 0.0
	if !t.LastSuccessAt.IsZero() {
		hoursSince := time.Since(t.LastSuccessAt).Hours()
		recencyBonus = 0.1 * math.Exp(-hoursSince/168) // Decays over ~1 week
	}

	// Proof of work bonus
	powBonus := math.Min(t.ProofOfWork*0.1, 0.1)

	score := successRate*volumeConfidence + vouchBonus - flagPenalty + recencyBonus + powBonus
	return math.Max(0, math.Min(1, score))
}

// NetworkHealthSummary aggregates overall network health metrics.
type NetworkHealthSummary struct {
	TotalPeers        int       `json:"totalPeers"`
	ActivePeers       int       `json:"activePeers"`       // Peers queried recently
	HealthyPeers      int       `json:"healthyPeers"`      // Peers with good metrics
	DegradedPeers     int       `json:"degradedPeers"`     // Peers with some issues
	UnhealthyPeers    int       `json:"unhealthyPeers"`    // Peers with poor metrics
	AvgLatencyMs      float64   `json:"avgLatencyMs"`
	MedianLatencyMs   float64   `json:"medianLatencyMs"`
	AvgTrustScore     float64   `json:"avgTrustScore"`
	AvgQualityScore   float64   `json:"avgQualityScore"`
	DirectConnectRate float64   `json:"directConnectRate"` // % of peers directly reachable
	RelayDependency   float64   `json:"relayDependency"`   // % of peers needing relay
	ComputedAt        time.Time `json:"computedAt"`
}

// FeatureVector represents a normalized feature vector for a peer.
// This is the core structure for KD-Tree based peer selection.
type FeatureVector struct {
	PeerID   string    `json:"peerId"`
	Features []float64 `json:"features"`
	Labels   []string  `json:"labels,omitempty"` // Optional feature names
}

// StandardPeerFeatures defines the standard feature set for peer selection.
// These map to dimensions in the KD-Tree.
type StandardPeerFeatures struct {
	LatencyMs       float64 `json:"latencyMs"`       // Lower is better
	HopCount        int     `json:"hopCount"`        // Lower is better
	GeoDistanceKm   float64 `json:"geoDistanceKm"`   // Lower is better
	TrustScore      float64 `json:"trustScore"`      // Higher is better (invert)
	BandwidthMbps   float64 `json:"bandwidthMbps"`   // Higher is better (invert)
	PacketLossRate  float64 `json:"packetLossRate"`  // Lower is better
	ConnectivityPct float64 `json:"connectivityPct"` // Higher is better (invert)
	NATScore        float64 `json:"natScore"`        // Higher is better (invert)
}

// ToFeatureSlice converts structured features to a slice for KD-Tree operations.
// Inversion is handled so that lower distance = better peer.
func (f StandardPeerFeatures) ToFeatureSlice() []float64 {
	return []float64{
		f.LatencyMs,
		float64(f.HopCount),
		f.GeoDistanceKm,
		1 - f.TrustScore,       // Invert: higher trust = lower value
		100 - f.BandwidthMbps,  // Invert: higher bandwidth = lower value (capped at 100)
		f.PacketLossRate,
		100 - f.ConnectivityPct, // Invert: higher connectivity = lower value
		1 - f.NATScore,          // Invert: higher NAT score = lower value
	}
}

// StandardFeatureLabels returns the labels for standard peer features.
func StandardFeatureLabels() []string {
	return []string{
		"latency_ms",
		"hop_count",
		"geo_distance_km",
		"trust_score_inv",
		"bandwidth_inv",
		"packet_loss_rate",
		"connectivity_inv",
		"nat_score_inv",
	}
}

// FeatureRanges defines min/max ranges for feature normalization.
type FeatureRanges struct {
	Ranges []AxisStats `json:"ranges"`
}

// DefaultPeerFeatureRanges returns sensible default ranges for peer features.
func DefaultPeerFeatureRanges() FeatureRanges {
	return FeatureRanges{
		Ranges: []AxisStats{
			{Min: 0, Max: 1000},   // Latency: 0-1000ms
			{Min: 0, Max: 20},     // Hops: 0-20
			{Min: 0, Max: 20000},  // Geo distance: 0-20000km (half Earth circumference)
			{Min: 0, Max: 1},      // Trust score (inverted): 0-1
			{Min: 0, Max: 100},    // Bandwidth (inverted): 0-100Mbps
			{Min: 0, Max: 1},      // Packet loss: 0-100%
			{Min: 0, Max: 100},    // Connectivity (inverted): 0-100%
			{Min: 0, Max: 1},      // NAT score (inverted): 0-1
		},
	}
}

// NormalizePeerFeatures normalizes peer features to [0,1] using provided ranges.
func NormalizePeerFeatures(features []float64, ranges FeatureRanges) []float64 {
	result := make([]float64, len(features))
	for i, v := range features {
		if i < len(ranges.Ranges) {
			result[i] = scale01(v, ranges.Ranges[i].Min, ranges.Ranges[i].Max)
		} else {
			result[i] = v
		}
	}
	return result
}

// WeightedPeerFeatures applies per-feature weights after normalization.
func WeightedPeerFeatures(normalized []float64, weights []float64) []float64 {
	result := make([]float64, len(normalized))
	for i, v := range normalized {
		w := 1.0
		if i < len(weights) {
			w = weights[i]
		}
		result[i] = v * w
	}
	return result
}
