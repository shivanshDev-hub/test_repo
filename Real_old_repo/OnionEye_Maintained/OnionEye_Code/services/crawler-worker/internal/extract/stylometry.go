package extract

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"time"
)

// VectorDim is the fixed width of the lexical style vector. It matches
// stylometric_profiles.feature_vector vector(384) so the fingerprint can be
// persisted and ANN-searched by pgvector without reshaping.
const VectorDim = 384

// MinStylometryWords mirrors the correlation-engine sample gate. Below this
// the profile is honestly downweighted rather than silently trusted.
const MinStylometryWords = 150

// AuthorBlock is one attributed writing sample on a page: the handle that
// wrote it, its normalized body text, and any post times rendered with it.
type AuthorBlock struct {
	Handle string
	Text   string
	Times  []time.Time
}

// Post-block anchors: forums/markets wrap each message in a card/post/message
// container carrying the author's <span class="username"> handle and a posted
// timestamp. The demo target mirrors this exact DOM shape.
var (
	rePostContainer = regexp.MustCompile(`(?is)<div[^>]*class="[^"]*(?:card|post|message|review)[^"]*"[^>]*>`)
	rePostUsername  = regexp.MustCompile(`(?is)class="username"[^>]*>([^<]+)<`)
)

// AttributionNoiseHandles are boilerplate authors that must never receive a
// writing profile (site chrome, banner text, stats rows).
var attributionNoiseHandles = map[string]bool{
	"admin": true, "administrator": true, "moderator": true, "mod": true,
	"system": true, "bot": true, "root": true, "anon": true, "anonymous": true,
}

// AttributeAuthorBlocks splits normalized (script/style-stripped) page text
// into per-author writing blocks. When the page provides no card/post
// containers, the whole body is attributed to the first real handle present
// (marketplace listing pages are single-author by construction). An empty
// result means there is no prose worth profiling.
func AttributeAuthorBlocks(normalized string) []AuthorBlock {
	idxs := rePostContainer.FindAllStringIndex(normalized, -1)
	if len(idxs) == 0 {
		// No thread structure: whole-page attribution to the primary handle.
		handle := firstRealHandle(normalized)
		if handle == "" {
			return nil
		}
		return []AuthorBlock{{Handle: handle, Text: strings.TrimSpace(normalized), Times: parseTimes(normalized)}}
	}

	var blocks []AuthorBlock
	seen := map[string]bool{}
	for i, span := range idxs {
		// The body is the text between this container's opening tag and the
		// next container's opening tag (or end of page) — the positional
		// equivalent of the original lookahead boundary, safe under RE2.
		bodyEnd := len(normalized)
		if i+1 < len(idxs) {
			bodyEnd = idxs[i+1][0]
		}
		body := normalized[span[1]:bodyEnd]
		hm := rePostUsername.FindStringSubmatch(body)
		if len(hm) < 2 {
			continue
		}
		handle := strings.TrimSpace(hm[1])
		if handle == "" || attributionNoiseHandles[strings.ToLower(handle)] {
			continue
		}
		if seen[strings.ToLower(handle)] {
			continue // keep the first block per author on this page
		}
		seen[strings.ToLower(handle)] = true
		blocks = append(blocks, AuthorBlock{
			Handle: handle,
			Text:   strings.TrimSpace(body),
			Times:  parseTimes(body),
		})
	}
	return blocks
}

// firstRealHandle returns the first non-noise DOM username handle on the page.
func firstRealHandle(text string) string {
	for _, m := range reHandleDOM.FindAllStringSubmatch(text, -1) {
		if len(m) > 1 {
			h := strings.TrimSpace(m[1])
			if h != "" && !attributionNoiseHandles[strings.ToLower(h)] && !isNoiseHandle(h) {
				return h
			}
		}
	}
	return ""
}

// parseTimes extracts ISO-8601 timestamps from an author block and parses them
// to UTC (this is the "posting time of day" source — real when the page renders
// per-post dates, honest fallback below when it does not).
func parseTimes(text string) []time.Time {
	var times []time.Time
	for _, raw := range reTs.FindAllString(text, -1) {
		// reTs allows a space separator; RFC3339 requires 'T'.
		raw = strings.Replace(raw, " ", "T", 1)
		t, err := time.Parse(time.RFC3339Nano, raw)
		if err != nil {
			continue
		}
		times = append(times, t.UTC())
	}
	return times
}

// CountWords is the sample-size used by the stylometry gate.
func CountWords(text string) int { return len(strings.Fields(text)) }

// HourHistogram buckets post times into 24 UTC hour bins (counts).
func HourHistogram(times []time.Time) []float64 {
	h := make([]float64, 24)
	for _, t := range times {
		h[t.Hour()]++
	}
	return h
}

// DayHistogram buckets post times into 7 UTC weekday bins (counts).
func DayHistogram(times []time.Time) []float64 {
	d := make([]float64, 7)
	for _, t := range times {
		d[int(t.Weekday())]++
	}
	return d
}

// MedianGapMinutes is the author's posting tempo (median inter-message gap in
// minutes). 0 when fewer than two times are present.
func MedianGapMinutes(times []time.Time) float64 {
	if len(times) < 2 {
		return 0
	}
	ts := make([]int64, len(times))
	for i, t := range times {
		ts[i] = t.Unix()
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i] < ts[j] })
	gaps := make([]float64, 0, len(ts)-1)
	for i := 1; i < len(ts); i++ {
		gaps = append(gaps, float64(ts[i]-ts[i-1])/60.0)
	}
	sort.Float64s(gaps)
	mid := len(gaps) / 2
	if len(gaps)%2 == 1 {
		return gaps[mid]
	}
	return (gaps[mid-1] + gaps[mid]) / 2.0
}

// LexicalVector hashes the author's prose deterministically into a fixed-width
// feature vector (signed count-min-sketch style bucketing). The token space is
// padded with fake-hash dimensions so the full 384 buckets are exercised:
//   - actual tokens (words / punctuation / structural markers) land in the
//     first half of the bucket space via a stable hash,
//   - a lexical "texture" profile (function-word density, avg word length,
//     punctuation rhythm, case/leet markers, emoji density) fills the second
//     half so short market-blurb samples still produce a spread fingerprint.
//
// The result is L2-normalized and deterministic: same text → same vector,
// never reversible back to the source prose.
func LexicalVector(text string) []float64 {
	v := make([]float64, VectorDim)

	words := strings.Fields(text)
	for _, w := range words {
		b := stableHash(strings.ToLower(w)) % (VectorDim / 2)
		v[b]++ // co-occurrence strength, not raw frequency
	}

	// Punctuation rhythm — punctuation is the single strongest style marker.
	var puncts float64
	for _, r := range []rune(text) {
		switch r {
		case '!':
			v[stableHash("!")%(VectorDim/2)] += 2
			puncts++
		case '?':
			v[stableHash("?")%(VectorDim/2)] += 2
			puncts++
		case '.', ',', ';', ':':
			v[stableHash(".")%(VectorDim/2)]++
			puncts++
		case '…':
			v[stableHash("...")%(VectorDim/2)] += 3
			puncts++
		}
	}

	// Lexical texture — second half of the bucket space.
	lower := strings.ToLower(text)
	textures := [][2]string{
		{"freq-the", "the_"}, {"freq-and", " and "}, {"freq-fuck", "fuck"},
		{"freq-bro", " bro"}, {"leet", "0"}, {"leet3", "3"},
	}
	for _, t := range textures {
		v[VectorDim/2+textureBucket(t[0])] = float64(strings.Count(lower, t[1]))
	}
	v[VectorDim/2+textureBucket("caps")] = float64(countUpper(text))
	v[VectorDim/2+textureBucket("emoji")] = float64(countEmoji(text))
	if len(words) > 0 {
		v[VectorDim/2+textureBucket("avg-wlen")] = float64(countLetters(text)) / float64(maxInt(len(words), 1))
	}
	if puncts > 0 {
		v[VectorDim/2+textureBucket("punct-density")] = puncts / float64(maxInt(len(words), 1))
	}

	return normalize(v)
}

// textureBucket maps a texture marker name into the second half of the vector.
func textureBucket(marker string) int { return int(stableHash(marker) % (VectorDim / 2)) }

// normalize L2-normalizes a vector in place and returns it.
func normalize(v []float64) []float64 {
	var n float64
	for _, x := range v {
		n += x * x
	}
	if n == 0 {
		return v
	}
	inv := 1.0 / math.Sqrt(n)
	for i := range v {
		v[i] *= inv
	}
	return v
}

// stableHash is a fixed (non-crypto) FNV-1a variant so vectors are reproducible
// across restarts and machines.
func stableHash(s string) uint64 {
	var h uint64 = 14695981039346656037
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}

func countUpper(s string) int {
	n := 0
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			n++
		}
	}
	return n
}

func countLetters(s string) int {
	n := 0
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			n++
		}
	}
	return n
}

func countEmoji(s string) int {
	n := 0
	for _, r := range s {
		if r > 0x1F000 {
			n++
		}
	}
	return n
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
