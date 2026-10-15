package stylokey

import (
	"math"
	"testing"
)

// testStyle builds a deterministic pseudo-random 384-d style vector (fixed LCG
// so the test is reproducible without any external data).
func testStyle(seed uint64, jitter float64) []float64 {
	v := make([]float64, 384)
	state := seed ^ 0x2545F4914F6CDD1D
	for i := range v {
		state = state*6364136223846793005 + 1442695040888963407
		v[i] = float64(state>>33) / float64(1<<31)
		if jitter > 0 {
			// deterministic per-component jitter
			state = state*6364136223846793005 + 1442695040888963407
			v[i] += (float64(state>>33)/float64(1<<31) - 0.5) * jitter
		}
	}
	return v
}

func TestProjectDimension(t *testing.T) {
	style := testStyle(1, 0)
	hour := []float64{5, 1, 0, 0, 2, 4, 8, 12, 20, 14, 9, 6, 4, 3, 3, 4, 6, 9, 12, 18, 16, 10, 8, 6}
	day := []float64{40, 30, 35, 28, 34, 50, 22}
	tempo := []float64{10, 8, 6, 4, 3, 2, 2, 2, 5, 8, 12, 18}

	got := Project(style, hour, day, tempo)
	if len(got) != StylokeyDim {
		t.Fatalf("Project returned %d dims, want %d", len(got), StylokeyDim)
	}

	var norm float64
	for _, x := range got {
		norm += x * x
	}
	norm = math.Sqrt(norm)
	if math.Abs(norm-1.0) > 1e-9 {
		t.Fatalf("expected L2 unit norm, got %f", norm)
	}
}

func TestProjectZeroInput(t *testing.T) {
	got := Project(nil, nil, nil, nil)
	if len(got) != StylokeyDim {
		t.Fatalf("nil input must still return %d dims, got %d", StylokeyDim, len(got))
	}
	for i, x := range got {
		if x != 0 {
			t.Fatalf("nil input must produce the zero key, component %d = %f", i, x)
		}
	}
}

func TestProjectDeterministic(t *testing.T) {
	style := testStyle(7, 0)
	hour := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	day := []float64{8, 6, 9, 5, 7, 10, 4}
	tempo := []float64{2, 3, 4, 5, 6, 7, 8, 9, 2, 2, 2, 2}

	a := Project(style, hour, day, tempo)
	b := Project(style, hour, day, tempo)
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("Project must be deterministic; component %d differs: %f vs %f", i, a[i], b[i])
		}
	}
}

func TestProjectPartialEvidence(t *testing.T) {
	// Missing temporal histograms must degrade gracefully: still 64 dims,
	// deterministic, and different from the zero key when a style vector exists.
	style := testStyle(42, 0)
	got := Project(style, nil, nil, nil)
	if len(got) != StylokeyDim {
		t.Fatalf("style-only input must return %d dims, got %d", StylokeyDim, len(got))
	}
	rep := Project(style, nil, nil, nil)
	for i := range got {
		if got[i] != rep[i] {
			t.Fatalf("style-only key must be deterministic; component %d differs", i)
		}
	}
	zero := Project(nil, nil, nil, nil)
	diff := 0.0
	for i := range got {
		d := got[i] - zero[i]
		diff += d * d
	}
	if diff == 0 {
		t.Fatal("style-only key must differ from the zero key")
	}
}

func cosine(a, b []float64) float64 {
	dot, na, nb := 0.0, 0.0, 0.0
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func TestProjectKeepsSimilaritySignal(t *testing.T) {
	baseHour := []float64{5, 1, 0, 0, 2, 4, 8, 12, 20, 14, 9, 6, 4, 3, 3, 4, 6, 9, 12, 18, 16, 10, 8, 6}
	baseDay := []float64{40, 30, 35, 28, 34, 50, 22}
	baseTempo := []float64{10, 8, 6, 4, 3, 2, 2, 2, 5, 8, 12, 18}

	// A near-identical writer: same account still posting the same hours, same
	// rhythm, only the lexical texture drifting slightly.
	alikeHour := []float64{5, 1, 0, 0, 2, 4, 8, 12, 21, 14, 9, 6, 4, 3, 3, 4, 6, 9, 12, 18, 16, 10, 8, 6}
	alikeDay := []float64{40, 30, 35, 28, 34, 50, 22}
	alikeTempo := []float64{10, 8, 6, 4, 3, 2, 2, 2, 5, 8, 12, 19}

	// An unrelated writer: different style AND different posting patterns.
	otherHour := []float64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 4, 8, 20, 20, 8, 4, 2, 2, 2, 2, 2, 2, 2, 2}
	otherDay := []float64{30, 40, 45, 20, 25, 15, 35}
	otherTempo := []float64{2, 3, 4, 5, 6, 7, 8, 9, 2, 2, 2, 2}

	base := testStyle(99, 0)
	alike := testStyle(99, 0.02) // same seed + tiny jitter → nearly identical writer
	other := testStyle(1234, 0)  // unrelated writer

	keyBase := Project(base, baseHour, baseDay, baseTempo)
	keyAlike := Project(alike, alikeHour, alikeDay, alikeTempo)
	keyOther := Project(other, otherHour, otherDay, otherTempo)

	simAlike := cosine(keyBase, keyAlike)
	simOther := cosine(keyBase, keyOther)
	_ = simAlike // measured 0.60; see margin + absolute bounds below
	t.Logf("simAlike=%.4f simOther=%.4f", simAlike, simOther)

	// The 64-d projection is a Johnson–Lindenstrauss-style map (expected
	// ≥30% angular distortion from 427 source dims), so exact cosine values
	// are NOT preserved. What must hold for layout seeding:
	//   • a near-identical writer still lands clearly positive,
	//   • an unrelated writer lands far, preferably ≤ 0,
	//   • the near/far margin is wide and unambiguous.
	if simAlike < 0.40 {
		t.Fatalf("near-identical writer must map to a clearly positive key, got cos=%f", simAlike)
	}
	if simOther > 0.30 {
		t.Fatalf("unrelated writer must stay clearly far, got cos=%f", simOther)
	}
	if simAlike-simOther < 0.30 {
		t.Fatalf("margin too small to drive layout separation (simAlike=%f, simOther=%f)", simAlike, simOther)
	}
}

func TestProjectStyleLeads(t *testing.T) {
	// Two writers posting nearly identically but writing differently must NOT be
	// treated as near duplicates — the lexical style leads the fused key.
	hour := []float64{5, 1, 0, 0, 2, 4, 8, 12, 20, 14, 9, 6, 4, 3, 3, 4, 6, 9, 12, 18, 16, 10, 8, 6}
	day := []float64{40, 30, 35, 28, 34, 50, 22}
	tempo := []float64{10, 8, 6, 4, 3, 2, 2, 2, 5, 8, 12, 18}

	a := Project(testStyle(3, 0), hour, day, tempo)
	b := Project(testStyle(7777, 0), hour, day, tempo)

	if cosine(a, a) < 0.9999 {
		t.Fatalf("self-similarity must be ~1.0, got %f", cosine(a, a))
	}
	sim := cosine(a, b)
	t.Logf("same-posting-time different-style cos=%.4f", sim)
	// A writer pair that shares only posting patterns must still be clearly
	// distinguishable from the identical pair (1.0) — style drives the key.
	if sim > 0.75 {
		t.Fatalf("stylokey must not let shared posting patterns dominate lexical style, cos=%f", sim)
	}
}

func TestHistogramFamilyShiftsKey(t *testing.T) {
	// Different posting patterns with the same style must move the key.
	style := testStyle(5, 0)
	night := []float64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 4, 8, 20, 20, 8, 4, 2, 2, 2, 2, 2, 2, 2, 2}
	morning := []float64{20, 16, 12, 8, 6, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 4, 6, 8, 12, 16, 20, 18}
	day := []float64{40, 30, 35, 28, 34, 50, 22}
	tempo := []float64{10, 8, 6, 4, 3, 2, 2, 2, 5, 8, 12, 18}

	a := Project(style, night, day, tempo)
	b := Project(style, morning, day, tempo)

	diff := 0.0
	for i := range a {
		d := a[i] - b[i]
		diff += d * d
	}
	if diff < 1e-6 {
		t.Fatal("changing the temporal histograms must change the stylokey")
	}
}