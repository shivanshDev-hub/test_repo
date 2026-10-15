package negcontrol

// temporal.go — Engine 1's built-in skeptic, now temporal. Before any
// candidate persona link is accepted on hard identifiers, this check asks:
// "do these two handles actually *behave* like the same human?".
//
// Design stance (from Architecture §3.1):
//   - Completely disjoint posting windows over meaningful sample sizes are a
//     WEAK-NEGATIVE signal — different timezones explain most disjointness —
//     so they apply a penalty, never a hard veto.
//   - High hourly overlap is corroborative only: rhythm entrainment supports a
//     link that hard shared identifiers already proposed; it never creates one
//     on its own.

import "math"

const (
	// MinTemporalSamples is the observation floor before any temporal claim.
	MinTemporalSamples = 4
	// DisjointThreshold is the hourly-overlap level that triggers the
	// coincidental-timing penalty (mirroring the design doc's
	// temporal_overlap_window check).
	DisjointThreshold = 0.12
	// EntrainedThreshold is the hourly-overlap level treated as corroboration.
	EntrainedThreshold = 0.60
	// TemporalPenalty is the confidence deduction for a disjoint window pair.
	TemporalPenalty = 0.13
)

// TemporalCheckResult is the output of the temporal overlap guard.
type TemporalCheckResult struct {
	Overlap    float64 `json:"overlap"`    // cosine on normalized 24-h histograms
	Penalty    float64 `json:"penalty"`    // subtracted from candidate link confidence
	Supportive bool    `json:"supportive"` // entrained rhythm corroboration
	Basis      string  `json:"basis"`
}

// EvaluateTemporalOverlap compares two actors' 24-bin hourly activity profiles
// (count histograms, any timezone-normalized offset) and returns the guard
// verdict for a candidate persona link between them.
func EvaluateTemporalOverlap(hourA, hourB []float64, nA, nB int) TemporalCheckResult {
	if len(hourA) != 24 || len(hourB) != 24 {
		return TemporalCheckResult{Basis: "Temporal overlap skipped — no hourly activity profile"}
	}
	if nA < MinTemporalSamples || nB < MinTemporalSamples {
		return TemporalCheckResult{
			Basis: "Temporal overlap skipped — insufficient observations (n<4) to claim rhythm",
		}
	}

	// Normalize to unit vectors: cosine over histograms measures SHAPE (the
	// rhythm) rather than raw posting volume.
	na, nb := HistogramNorm(hourA), HistogramNorm(hourB)
	var dot, na2, nb2 float64
	for i := range na {
		dot += na[i] * nb[i]
		na2 += na[i] * na[i]
		nb2 += nb[i] * nb[i]
	}
	overlap := 0.0
	if na2 > 0 && nb2 > 0 {
		overlap = dot / (math.Sqrt(na2) * math.Sqrt(nb2))
		if overlap > 1 {
			overlap = 1
		}
	}

	switch {
	case overlap <= DisjointThreshold:
		return TemporalCheckResult{
			Overlap: overlap,
			Penalty: TemporalPenalty,
			Basis:   "Coincidental-timing guard: disjoint posting windows (hourly overlap <= 0.12) — different timezone or different operator",
		}
	case overlap >= EntrainedThreshold:
		return TemporalCheckResult{
			Overlap:    overlap,
			Supportive: true,
			Basis:      "Shared posting rhythm (hourly overlap >= 0.60) corroborates hard-identifier link",
		}
	default:
		return TemporalCheckResult{
			Overlap: overlap,
			Basis:   "Neutral temporal overlap — no penalty, no corroboration",
		}
	}
}

// HistogramNorm L1-normalizes a count histogram to a [0,1] shape vector.
func HistogramNorm(h []float64) []float64 {
	out := make([]float64, len(h))
	var sum float64
	for _, v := range h {
		sum += v
	}
	if sum == 0 {
		return out
	}
	for i, v := range h {
		out[i] = v / sum
	}
	return out
}
