package engine

import (
	"fmt"
	"math"
)

const (
	// MinStylometryWords is the minimum sample threshold required for statistical significance.
	MinStylometryWords = 150
	// MinTempoGaps is the minimum number of recorded inter-arrival gaps before
	// the tempo cadence shape is trusted for fusion. Fewer gaps (a single
	// activity burst) are statistically meaningless rhythm evidence.
	MinTempoGaps = 5
	// HighStylometrySimThreshold is the cosine similarity cutoff for high confidence author attribution.
	HighStylometrySimThreshold = 0.88
	// ModerateStylometrySimThreshold is the cutoff for moderate similarity.
	ModerateStylometrySimThreshold = 0.75
)

// CosineSimilarity computes the cosine similarity between two feature vectors.
func CosineSimilarity(a, b []float64) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, fmt.Errorf("empty feature vector")
	}
	if len(a) != len(b) {
		return 0, fmt.Errorf("vector dimension mismatch: len(a)=%d, len(b)=%d", len(a), len(b))
	}

	dotProduct := 0.0
	normA := 0.0
	normB := 0.0

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0, nil
	}

	sim := dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
	// Clamp float precision errors to [-1.0, 1.0]
	if sim > 1.0 {
		sim = 1.0
	} else if sim < -1.0 {
		sim = -1.0
	}
	return sim, nil
}

// StylometryMatchResult contains the evaluation of two writing style profiles.
type StylometryMatchResult struct {
	Similarity float64 `json:"similarity"`
	HardScore  float64 `json:"hard_score"`
	IDFWeight  float64 `json:"idf_weight"`
	SampleGate bool    `json:"sample_gate_passed"`
	Basis      string  `json:"basis"`
}

// EvaluateStylometry computes cosine similarity and applies the statistical word-count sample gate.
func EvaluateStylometry(vectorA, vectorB []float64, sampleWordsA, sampleWordsB int) StylometryMatchResult {
	sim, err := CosineSimilarity(vectorA, vectorB)
	if err != nil {
		return StylometryMatchResult{
			Similarity: 0.0,
			HardScore:  0.0,
			IDFWeight:  0.0,
			SampleGate: false,
			Basis:      fmt.Sprintf("Stylometry evaluation error: %v", err),
		}
	}

	minSample := sampleWordsA
	if sampleWordsB < minSample && sampleWordsB > 0 {
		minSample = sampleWordsB
	}
	if minSample == 0 {
		minSample = sampleWordsA + sampleWordsB
	}

	gatePassed := minSample >= MinStylometryWords

	if !gatePassed {
		// Insufficient sample: heavily downweight attribution to prevent false stylometric clustering
		return StylometryMatchResult{
			Similarity: sim,
			HardScore:  sim * 0.25,
			IDFWeight:  0.30,
			SampleGate: false,
			Basis:      fmt.Sprintf("Insufficient word count (%d words < %d min required) — stylometry downweighted", minSample, MinStylometryWords),
		}
	}

	switch {
	case sim >= HighStylometrySimThreshold:
		return StylometryMatchResult{
			Similarity: sim,
			HardScore:  0.94,
			IDFWeight:  0.95,
			SampleGate: true,
			Basis:      fmt.Sprintf("High-confidence stylometric match (Cosine Similarity: %.3f >= %.2f)", sim, HighStylometrySimThreshold),
		}
	case sim >= ModerateStylometrySimThreshold:
		return StylometryMatchResult{
			Similarity: sim,
			HardScore:  0.72,
			IDFWeight:  0.75,
			SampleGate: true,
			Basis:      fmt.Sprintf("Moderate stylometric match (Cosine Similarity: %.3f >= %.2f)", sim, ModerateStylometrySimThreshold),
		}
	default:
		return StylometryMatchResult{
			Similarity: sim,
			HardScore:  0.15,
			IDFWeight:  0.20,
			SampleGate: true,
			Basis:      fmt.Sprintf("Low/divergent stylometric similarity (Cosine Similarity: %.3f < %.2f)", sim, ModerateStylometrySimThreshold),
		}
	}
}

// BehavioralFusionWeights are the documented fusion weights between the style
// and temporal ("human") modalities. Style leads; posting rhythm corroborates;
// tempo is a tie-breaker. Weights renormalize to whichever modalities are
// actually present so a missing histogram never zeroes the others out.
var BehavioralFusionWeights = struct {
	Style float64
	Hour  float64
	Day   float64
	Tempo float64
}{Style: 0.50, Hour: 0.20, Day: 0.10, Tempo: 0.20}

// BehavioralFusionResult is the fused same-author score for a candidate pair.
type BehavioralFusionResult struct {
	StylometrySimilarity float64 `json:"stylometry_similarity"`
	HourOverlap          float64 `json:"hour_overlap"`
	DayOverlap           float64 `json:"day_overlap"`
	TempoSimilarity      float64 `json:"tempo_similarity"`
	FusedScore           float64 `json:"fused_score"`
	SampleGatePassed     bool    `json:"sample_gate_passed"`
	Basis                string  `json:"basis"`
}

// EvaluateBehavioralFusion combines writing style (384-d vectors) with the
// temporal "human" vectors into one [0,1] same-author score.
//
// It reuses the engine's own CosineSimilarity for every modality.
// - hour/day histograms are treated as count vectors then normalized (shape).
// - tempo histogram is expected as a 12-d unit-ish vector (shape).
//
// When either profile is below the statistical sample gate the fused score is
// honestly downweighted (×0.35), mirroring EvaluateStylometry's skepticism —
// temporal rhythm is corroborative only.
func EvaluateBehavioralFusion(styleA, styleB []float64, hourA, hourB, dayA, dayB []float64, tempoA, tempoB []float64, gapA, gapB float64, wordsA, wordsB int) BehavioralFusionResult {
	var (
		styleSim, hourSim, daySim, tempoSim     float64
		haveStyle, haveHour, haveDay, haveTempo bool
	)

	if len(styleA) == len(styleB) && len(styleA) > 0 {
		if s, err := CosineSimilarity(styleA, styleB); err == nil {
			styleSim, haveStyle = s, true
		}
	}
	if len(hourA) == len(hourB) && len(hourA) > 0 {
		hnA := HistogramNormalize(hourA)
		hnB := HistogramNormalize(hourB)
		if s, err := CosineSimilarity(hnA, hnB); err == nil {
			hourSim, haveHour = s, true
		}
	}
	if len(dayA) == len(dayB) && len(dayA) > 0 {
		dnA := HistogramNormalize(dayA)
		dnB := HistogramNormalize(dayB)
		if s, err := CosineSimilarity(dnA, dnB); err == nil {
			daySim, haveDay = s, true
		}
	}
	if s, ok := TempoHistogramSim(tempoA, tempoB); ok {
		tempoSim, haveTempo = s, true
	} else if gapA > 0 && gapB > 0 {
		// Fallback tempo similarity: use median gap (scalar) when histogram is
		// missing or undersized — still a weak rhythm corroboration.
		mean := (gapA + gapB) / 2.0
		tempoSim, haveTempo = math.Exp(-math.Abs(gapA-gapB)/mean), true
	}

	// Renormalize available weights.
	var total float64
	if haveStyle {
		total += BehavioralFusionWeights.Style
	}
	if haveHour {
		total += BehavioralFusionWeights.Hour
	}
	if haveDay {
		total += BehavioralFusionWeights.Day
	}
	if haveTempo {
		total += BehavioralFusionWeights.Tempo
	}

	fused := 0.0
	if total > 0 {
		if haveStyle {
			fused += BehavioralFusionWeights.Style / total * styleSim
		}
		if haveHour {
			fused += BehavioralFusionWeights.Hour / total * hourSim
		}
		if haveDay {
			fused += BehavioralFusionWeights.Day / total * daySim
		}
		if haveTempo {
			fused += BehavioralFusionWeights.Tempo / total * tempoSim
		}
	}

	minSample := wordsA
	if wordsB < minSample && wordsB > 0 {
		minSample = wordsB
	}
	if minSample == 0 {
		minSample = wordsA + wordsB
	}
	gatePassed := minSample >= MinStylometryWords

	res := BehavioralFusionResult{
		StylometrySimilarity: styleSim,
		HourOverlap:          hourSim,
		DayOverlap:           daySim,
		TempoSimilarity:      tempoSim,
		FusedScore:           fused,
		SampleGatePassed:     gatePassed,
	}
	switch {
	case !gatePassed:
		res.FusedScore = fused * 0.35
		res.Basis = fmt.Sprintf("Behavioral fusion downweighted: %d words < %d minimum — style+rhythm corroborative only", minSample, MinStylometryWords)
	case fused >= HighStylometrySimThreshold:
		res.Basis = fmt.Sprintf("High-confidence behavioral match — style %.3f, hourly overlap %.3f, weekday %.3f, tempo %.3f", styleSim, hourSim, daySim, tempoSim)
	case fused >= ModerateStylometrySimThreshold:
		res.Basis = fmt.Sprintf("Moderate behavioral match — style %.3f, hourly overlap %.3f, weekday %.3f, tempo %.3f", styleSim, hourSim, daySim, tempoSim)
	default:
		res.Basis = fmt.Sprintf("Divergent behavioral profiles — style %.3f, hourly overlap %.3f (corroboration below %.2f)", styleSim, hourSim, ModerateStylometrySimThreshold)
	}
	return res
}

// tempoGapCount returns the cumulative number of recorded inter-arrival gaps in
// a tempo histogram — the sum of its count bins. Used as the sample-size gate
// for the cadence shape (see MinTempoGaps).
func tempoGapCount(h []float64) float64 {
	var n float64
	for _, v := range h {
		n += v
	}
	return n
}

// TempoHistogramSim is the normalized-cosine similarity between two tempo
// histograms. It implements the documented "rhythm shape" comparison: both
// profiles are first normalized to relative frequency (so a quiet author with
// 6 gaps is on equal footing with a prolific one), then compared with cosine.
// The MinTempoGaps gate on total gaps keeps a burst of 2-4 posts from being
// mistaken for a genuine cadence.
func TempoHistogramSim(a, b []float64) (float64, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	if tempoGapCount(a) < MinTempoGaps || tempoGapCount(b) < MinTempoGaps {
		return 0, false
	}
	na := HistogramNormalize(a)
	nb := HistogramNormalize(b)
	s, err := CosineSimilarity(na, nb)
	if err != nil {
		return 0, false
	}
	return s, true
}

// HistogramNormalize converts a count histogram into a unit vector so cosine
// similarity measures *shape* (relative activity), not raw volume.
func HistogramNormalize(h []float64) []float64 {
	norm := make([]float64, len(h))
	var sum float64
	for _, v := range h {
		sum += v
	}
	if sum == 0 {
		return norm
	}
	for i, v := range h {
		norm[i] = v / sum
	}
	return norm
}
