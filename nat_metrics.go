package poindexter

import (
	"math"
	"time"
)

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
