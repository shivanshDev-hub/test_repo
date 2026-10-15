package engine

// engine.go implements Phase 9 Engine 2: Composite Confidence Score & Multi-Vector Attribution.
//
// Formula:
//    C_final = min( 1.0, sum(S_i * IDF_i * time_decay) / max(1.0, log2(1 + NC_passed)) )
//
// This evaluates multi-modal identifiers (Emails, Wallets, PGP Keys, Handles, Stylometry, CSAM Flags)
// while rapidly driving signal near zero for shared boilerplate infrastructure or high-volume noise.

import (
	"fmt"
	"math"
	"time"
)

// Evidence represents a single CSAM flag (either Hash hit or ML score >= 0.45).
type Evidence struct {
	ImageSrc        string    `json:"image_src"`
	ClassifierScore float64   `json:"classifier_score"` // [0.45, 1.0] for ML, or 1.0 for hashes
	IsHashMatch     bool      `json:"is_hash_match"`
	FlaggedAt       time.Time `json:"flagged_at"`
}

// ActorProfile represents the gathered forensic context and extracted vectors for an entity.
type ActorProfile struct {
	ActorID          string              `json:"actor_id"`
	NCPassed         int                 `json:"nc_passed"`
	NCFlagged        int                 `json:"nc_flagged"`
	Evidence         []Evidence          `json:"evidence,omitempty"`
	Artifacts        []ExtractedArtifact `json:"artifacts,omitempty"`
	StylometryVector []float64           `json:"stylometry_vector,omitempty"`
	SampleWordCount  int                 `json:"sample_word_count,omitempty"`
	// Temporal ("human") vector family — posting time-of-day rhythm.
	HourHistogram    []float64 `json:"hour_histogram,omitempty"` // 24 UTC bins
	DayHistogram     []float64 `json:"day_histogram,omitempty"`  // 7 UTC bins
	MedianGapMinutes float64   `json:"median_gap_minutes,omitempty"`
	// TempoHistogram is the inter-arrival cadence shape (12-bins: 8×5m + 4 tail)
	// derived from the last 200 evidence observations. Cosine on the normalized
	// shape measures rhythmic readiness, independent of raw posting volume.
	TempoHistogram []float64 `json:"tempo_histogram,omitempty"`
}

// Config holds the hyperparameters for Engine 2.
type Config struct {
	// HalfLifeDays determines how fast older flags/artifacts decay in relevance.
	HalfLifeDays float64 `json:"half_life_days"`
	// FloorThreshold is the minimum score to emit an alert (typically 0.85).
	FloorThreshold float64 `json:"floor_threshold"`
}

// DefaultConfig matches the production hyperparameters.
var DefaultConfig = Config{
	HalfLifeDays:   30.0,
	FloorThreshold: 0.85,
}

// CorrelationResult contains the full forensic evaluation and scored artifacts.
type CorrelationResult struct {
	ActorID         string           `json:"actor_id"`
	CompositeScore  float64          `json:"composite_score"`
	IsAlert         bool             `json:"is_alert"`
	Threshold       float64          `json:"threshold"`
	NCDampener      float64          `json:"nc_dampener"`
	ScoredArtifacts []ScoredArtifact `json:"scored_artifacts,omitempty"`
}

// calculateDecay computes the half-life exponential decay factor.
func calculateDecay(ageDays, halfLifeDays float64) float64 {
	if ageDays <= 0 {
		return 1.0
	}
	if halfLifeDays <= 0 {
		halfLifeDays = 30.0
	}
	return math.Pow(0.5, ageDays/halfLifeDays)
}

// ComputeComposite calculates the C_final score for an actor's profile.
func ComputeComposite(profile ActorProfile, cfg Config, now time.Time) float64 {
	res := EvaluateProfile(profile, cfg, now)
	return res.CompositeScore
}

// EvaluateProfile computes detailed multi-vector scores, applying time decay and negative control dampening.
func EvaluateProfile(profile ActorProfile, cfg Config, now time.Time) CorrelationResult {
	if cfg.HalfLifeDays <= 0 {
		cfg.HalfLifeDays = 30.0
	}
	if cfg.FloorThreshold <= 0 {
		cfg.FloorThreshold = 0.85
	}

	numerator := 0.0
	var scoredArtifacts []ScoredArtifact

	// 1. Process legacy/content-safety flags (Evidence)
	for _, ev := range profile.Evidence {
		s_i := ev.ClassifierScore
		if ev.IsHashMatch {
			s_i = 1.0
		}

		ageHours := now.Sub(ev.FlaggedAt).Hours()
		ageDays := ageHours / 24.0
		decay := calculateDecay(ageDays, cfg.HalfLifeDays)

		idf := 1.0 // Content safety flags have full weight
		weight := s_i * idf * decay
		numerator += weight
	}

	// 2. Process extracted intelligence artifacts (Emails, Wallets, PGP, Handles, Infra)
	for _, art := range profile.Artifacts {
		scored := ScoreArtifact(art, now, cfg.HalfLifeDays)
		scoredArtifacts = append(scoredArtifacts, scored)
		numerator += scored.FinalScore
	}

	// 2b. Stylometric fingerprint — the actor's own writing-style vector. This
	// is deliberately a modest supporting contribution (hard 0.62, IDF 0.80):
	// the Strong signal lives in the *pairwise* behavioral fusion that scores
	// candidate persona links, not in the per-actor composite.
	if len(profile.StylometryVector) > 0 && profile.SampleWordCount > 0 {
		fingerprint := ScoreArtifact(ExtractedArtifact{
			Type:        VectorStylometry,
			Value:       "stylometric_profile",
			SampleWords: profile.SampleWordCount,
			ObservedAt:  now,
		}, now, cfg.HalfLifeDays)
		fingerprint.Basis = fmt.Sprintf("Stylometric fingerprint (%d words)", profile.SampleWordCount)
		scoredArtifacts = append(scoredArtifacts, fingerprint)
		numerator += fingerprint.FinalScore
	}

	// 3. Denominator: logarithmic penalty for clean volume.
	// We use log2(1 + NC_passed). If NC_passed is 0, denominator is 0 -> floored at 1.0.
	denominator := math.Log2(1.0 + float64(profile.NCPassed))
	if denominator < 1.0 {
		denominator = 1.0
	}

	// 4. Final division & clamp to [0, 1.0]
	cFinal := numerator / denominator
	if cFinal > 1.0 {
		cFinal = 1.0
	}

	return CorrelationResult{
		ActorID:         profile.ActorID,
		CompositeScore:  cFinal,
		IsAlert:         cFinal >= cfg.FloorThreshold,
		Threshold:       cfg.FloorThreshold,
		NCDampener:      denominator,
		ScoredArtifacts: scoredArtifacts,
	}
}

// Evaluate determines if the actor crosses the notification threshold.
func Evaluate(profile ActorProfile, cfg Config, now time.Time) (float64, bool) {
	res := EvaluateProfile(profile, cfg, now)
	return res.CompositeScore, res.IsAlert
}
