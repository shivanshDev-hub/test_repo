// Package stylokey builds a fixed-width, deterministic behavioral key for a
// single actor from the raw lexical style vector (384-d) fused with the
// temporal "human" histogram family (hour/day/tempo).
//
// The key is:
//
//	1. per-modality normalization — every block is mean-centered and
//	   L2-normalized (count histograms are first converted to shape), so the
//	   projection captures writing/rhythm *shape*, not raw volume and not the
//	   shared non-negativity bias of the source fingerprints,
//	2. concatenation,
//	3. a deterministic ±1 feature-hash (Rademacher) projection down to
//	   StylokeyDim fixed dimensions — a Johnson–Lindenstrauss-style random
//	   projection whose coefficients come from a fixed-seed splitmix64 hash,
//	   so no training artifact or stored matrix is ever required,
//	4. final L2 normalization.
//
// Everything is pure integer/float arithmetic over a fixed seed: identical
// inputs always yield identical keys on any platform and any restart, which is
// what makes the dashboard graph layout reproducible across polls and deploys.
package stylokey

import "math"

// StylokeyDim is the fixed width of the projected vector. It must never change
// once a deployment goes live — graph layouts and any future ANN index depend
// on the embedding's dimensionality staying stable.
const StylokeyDim = 64

// projectionSeed is the fixed hash seed for the deterministic projection. It
// is arbitrary but immutable: changing it silently reshuffles every layout.
const projectionSeed uint64 = 0x5EED5EED

// splitmix64 is a public-domain deterministic 64-bit hash. It generates the
// ±1 projection coefficients so no random number generator state or training
// artifact is needed.
func splitmix64(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	z := x
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// normalizeL2 scales v to unit Euclidean norm. A zero or near-zero vector is
// left untouched.
func normalizeL2(v []float64) {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	if sum <= 0 {
		return
	}
	n := math.Sqrt(sum)
	for i := range v {
		v[i] /= n
	}
}

// histogramShape normalizes a count histogram to a relative-frequency shape, so
// a quiet author (6 gaps) and a prolific one (600 gaps) compare on rhythm, not
// volume. A zero-sum histogram yields an all-zero block.
func histogramShape(h []float64) []float64 {
	out := make([]float64, len(h))
	var total float64
	for _, v := range h {
		total += v
	}
	if total <= 0 {
		return out
	}
	for i, v := range h {
		out[i] = v / total
	}
	return out
}

// meanCenter subtracts the block mean so non-negative modalities do not all
// collapse onto the shared "positive-orthant" direction. This is what lets the
// projection separate writers by *shape* rather than by raw positivity.
func meanCenter(v []float64) {
	var sum float64
	for _, x := range v {
		sum += x
	}
	mean := sum / float64(len(v))
	for i := range v {
		v[i] -= mean
	}
}

// centeredUnit normalizes a modality block for projection: it mean-centers the
// block (removing the shared positivity bias) and then L2-normalizes it so
// every block contributes on an equal footing.
func centeredUnit(v []float64) []float64 {
	out := append([]float64(nil), v...)
	meanCenter(out)
	normalizeL2(out)
	return out
}

// Modality weights mirror the correlation engine's BehavioralFusionWeights so
// the fused key agrees with the documented scoring philosophy: the lexical
// style leads (0.50), posting rhythm corroborates (0.20 hour / 0.10 day /
// 0.20 tempo). Each normalized block is scaled by sqrt(weight) so its energy
// in the projected key equals its share of the modality budget.
const (
	weightStyle = 0.50
	weightHour  = 0.20
	weightDay   = 0.10
	weightTempo = 0.20
)

// modalityBlock is one normalized-to-shape input block with its fusion weight.
type modalityBlock struct {
	data   []float64
	weight float64
}

// anyNonZero reports whether at least one component of v is non-zero.
func anyNonZero(v []float64) bool {
	for _, x := range v {
		if x != 0 {
			return true
		}
	}
	return false
}

// Project computes the StylokeyDim-wide stylokey for one actor.
//
//   - style:   raw lexical vector (typically 384-d, but any width works).
//   - hour:    posting-time-of-day count histogram (24 UTC bins in practice).
//   - day:     posting weekday count histogram (7 bins in practice).
//   - tempo:   inter-arrival cadence histogram (12 bins in practice).
//
// Every modality is normalized to a centered unit vector (count histograms
// first become shape), then scaled by its fusion weight — so *shape*, never
// raw volume or the shared non-negativity bias of the source fingerprints,
// drives the key, and the temporal family can never swamp the lexical style.
// The modality weights renormalize over whatever blocks are actually present
// (mirroring EvaluateBehavioralFusion), and each block may be empty/nil. An
// entirely empty input yields the zero vector of length StylokeyDim (the
// caller should treat that as "no behavioral evidence available").
func Project(style, hour, day, tempo []float64) []float64 {
	out := make([]float64, StylokeyDim)

	// 1. Build the fused source vector, one weighted normalized block per
	//    modality that actually carries evidence.
	var src []float64
	modalSnap := []modalityBlock{
		{data: style, weight: weightStyle},
		{data: hour, weight: weightHour},
		{data: day, weight: weightDay},
		{data: tempo, weight: weightTempo},
	}

	present := 0.0
	blocks := make([]modalityBlock, 0, len(modalSnap))
	for _, b := range modalSnap {
		if len(b.data) == 0 {
			continue
		}
		cooked := b.data
		if b.weight != weightStyle { // lexical vector is already a fingerprint; histograms are counts
			cooked = histogramShape(b.data)
		}
		cooked = centeredUnit(cooked)
		if anyNonZero(cooked) {
			present += b.weight
			blocks = append(blocks, modalityBlock{data: cooked, weight: b.weight})
		}
	}
	if present == 0 {
		return out // zero key → no behavioral evidence
	}
	for _, b := range blocks {
		scale := math.Sqrt(b.weight / present)
		for _, x := range b.data {
			src = append(src, x*scale)
		}
	}

	// 2. Deterministic ±1 random projection (feature hashing) down to
	//    StylokeyDim. sign(hash(i,j)) ∈ {-1,+1}; the 1/sqrt(k) scale keeps the
	//    projected vector's variance comparable to the source's.
	k := float64(len(src))
	for i := range src {
		h := projectionSeed ^ uint64(i+1)
		for j := range out {
			h = splitmix64(h)
			sign := 1.0
			if h&1 == 1 {
				sign = -1.0
			}
			out[j] += sign * src[i] / math.Sqrt(k)
		}
	}

	// 3. Final L2 normalization so cosine similarity on stylokeys is a true
	//    [−1, 1] comparable metric (identical to the engine's vector usage).
	normalizeL2(out)
	return out
}