package engine_test

import (
	"testing"

	"github.com/onioneye/correlation-engine/internal/engine"
)

func TestCosineSimilarity(t *testing.T) {
	// 1. Identical vectors -> Cosine Similarity = 1.0
	v1 := []float64{0.2, 0.4, 0.6, 0.8}
	v2 := []float64{0.2, 0.4, 0.6, 0.8}
	sim, err := engine.CosineSimilarity(v1, v2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sim < 0.999 {
		t.Fatalf("expected similarity ~1.0 for identical vectors, got %f", sim)
	}

	// 2. Orthogonal vectors -> Cosine Similarity = 0.0
	v3 := []float64{1.0, 0.0}
	v4 := []float64{0.0, 1.0}
	simOrth, err := engine.CosineSimilarity(v3, v4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if simOrth > 0.001 {
		t.Fatalf("expected similarity 0.0 for orthogonal vectors, got %f", simOrth)
	}
}

func TestEvaluateStylometry_SampleWordCountGate(t *testing.T) {
	// Two vectors with high cosine similarity (~0.95)
	vA := []float64{0.35, 0.12, 0.88, 0.45, 0.22}
	vB := []float64{0.34, 0.11, 0.89, 0.44, 0.23}

	// Case 1: Insufficient sample size (< 150 words, e.g. 45 words from short chat)
	resSmall := engine.EvaluateStylometry(vA, vB, 45, 60)
	if resSmall.SampleGate {
		t.Fatalf("expected SampleGate=false for small sample word count (45 words)")
	}
	if resSmall.HardScore > 0.30 {
		t.Fatalf("expected severely penalized score for short text sample, got hard_score=%f", resSmall.HardScore)
	}

	// Case 2: Statistically significant sample size (>= 150 words, e.g. 350 words)
	resLarge := engine.EvaluateStylometry(vA, vB, 350, 420)
	if !resLarge.SampleGate {
		t.Fatalf("expected SampleGate=true for large sample word count (350 words)")
	}
	if resLarge.HardScore < 0.90 {
		t.Fatalf("expected high hard score >= 0.90 for high similarity + adequate sample, got %f", resLarge.HardScore)
	}
}

// TestTempoHistogramSim checks the normalized-cosine cadence comparison and its
// MinTempoGaps sample gate.
func TestTempoHistogramSim(t *testing.T) {
	twelve := make([]float64, 12)

	// Identical cadence shapes (both authors burst every ~5 minutes) -> 1.0.
	a := make([]float64, 12)
	a[0], a[1], a[2] = 8, 4, 2 // 14 gaps -> bursts of frequent posting
	b := make([]float64, 12)
	b[0], b[1], b[2] = 4, 2, 1 // same relative shape, smaller absolute volume

	s, ok := engine.TempoHistogramSim(a, b)
	if !ok || s < 0.999 {
		t.Fatalf("expected identical tempo shapes -> sim ~1.0, got %f ok=%v", s, ok)
	}

	// Divergent cadence (author c posts mostly 6-12h apart) -> clearly lower.
	c := twelve
	c[10] = 14
	if s2, ok2 := engine.TempoHistogramSim(a, c); !ok2 || s2 > 0.50 {
		t.Fatalf("expected divergent tempo shapes -> sim < 0.5, got %f ok=%v", s2, ok2)
	}

	// Sample gate: fewer than MinTempoGaps gaps is statistically meaningless.
	burst := make([]float64, 12)
	burst[0] = 3 // only 3 gaps
	if _, ok := engine.TempoHistogramSim(a, burst); ok {
		t.Fatalf("expected tempo sample gate to reject an undersized histogram (3 gaps < 5)")
	}
}

// TestEvaluateBehavioralFusion_TempoIntegration verifies the tempo histogram
// joins EvaluateBehavioralFusion's renormalized fusion and that an undersized
// tempo sample only yields the weak scalar fallback.
func TestEvaluateBehavioralFusion_TempoIntegration(t *testing.T) {
	styleA := []float64{0.35, 0.12, 0.88, 0.45, 0.22, 0.60}
	styleB := []float64{0.34, 0.11, 0.89, 0.44, 0.23, 0.61}
	hourA := make([]float64, 24)
	hourA[12] = 10
	hourB := make([]float64, 24)
	hourB[12] = 12
	dayA := make([]float64, 7)
	dayA[1] = 6
	dayB := make([]float64, 7)
	dayB[1] = 7

	// Tempo histograms describing the same "rapid-burst" cadence. 12 bins:
	// 8x5m then 40-60m, 1-6h, 6-12h, 12-24h.
	tempoA := make([]float64, 12)
	tempoA[0], tempoA[1], tempoA[2] = 9, 4, 3
	tempoB := make([]float64, 12)
	tempoB[0], tempoB[1], tempoB[2] = 5, 2, 2

	res := engine.EvaluateBehavioralFusion(styleA, styleB, hourA, hourB, dayA, dayB, tempoA, tempoB, 18, 20, 400, 450)
	if res.TempoSimilarity < 0.95 {
		t.Fatalf("expected high tempo similarity for matching cadence shapes, got %f", res.TempoSimilarity)
	}
	// With style ~1.0, hour ~1.0, day ~1.0 AND tempo ~1.0 the fused score should
	// be at the high-confidence tier.
	if res.FusedScore < 0.88 {
		t.Fatalf("expected high-confidence fused score with all four modalities aligned, got %f", res.FusedScore)
	}
	if !res.SampleGatePassed {
		t.Fatalf("expected sample gate to pass for 400+ word profiles")
	}

	// Undersized tempo (3 gaps) -> falls back to the scalar gap match (weak
	// corroboration). With widely divergent median gaps the fallback must NOT
	// outscore a genuine full-histogram cadence match.
	undersized := make([]float64, 12)
	undersized[0] = 3
	resFallback := engine.EvaluateBehavioralFusion(styleA, styleB, hourA, hourB, dayA, dayB, tempoA, undersized, 18, 720, 400, 450)
	if resFallback.FusedScore > res.FusedScore {
		t.Fatalf("expected undersized tempo not to beat a full tempo match: fallback=%f full=%f", resFallback.FusedScore, res.FusedScore)
	}
	if resFallback.TempoSimilarity > 0.30 {
		t.Fatalf("expected weak scalar fallback tempo for divergent median gaps, got %f", resFallback.TempoSimilarity)
	}
}
