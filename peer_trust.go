package poindexter

import (
	"math"
	"time"
)

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
