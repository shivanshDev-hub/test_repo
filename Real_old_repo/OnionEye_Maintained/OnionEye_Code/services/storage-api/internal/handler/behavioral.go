package handler

import (
	"context"
	"database/sql"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetActorBehavioral returns an actor's behavioral profile (384-d stylometry
// through pgvector) plus its k nearest "writes-like" lookalikes over the HNSW
// index. Consumed by dashboard-api for the dossier STYLOMETRY + TEMPORAL
// sections.
func (h *Handler) GetActorBehavioral(c *gin.Context) {
	actorID := strings.TrimSpace(c.Param("id"))
	k := 8
	if v := c.Query("k"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 25 {
			k = n
		}
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var (
		words      int
		obs        int
		lastActive time.Time
		gap        sql.NullInt64
		hour       []float32
		day        []float32
		tempo      []float32
		vecTxt     string
	)
	err := h.pool.QueryRow(ctx, `
		SELECT sp.sample_word_count, sp.observations_count, sp.last_active_at,
		       sp.median_gap_minutes,
		       COALESCE(sp.hour_histogram, ARRAY[]::real[]),
		       COALESCE(sp.day_histogram, ARRAY[]::real[]),
		       COALESCE(sp.tempo_histogram, ARRAY[]::real[]),
		       sp.feature_vector::text
		FROM stylometric_profiles sp
		WHERE sp.actor_id = $1`, actorID).
		Scan(&words, &obs, &lastActive, &gap, &hour, &day, &tempo, &vecTxt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no behavioral profile for actor: " + err.Error()})
		return
	}

	vec := parseVectorText(vecTxt)
	lookalikes := make([]gin.H, 0)
	if len(vec) == 384 {
		rows, qerr := h.pool.Query(ctx, `
			SELECT a2.primary_handle, s2.sample_word_count, s2.observations_count,
			       1 - (s2.feature_vector <=> $1::vector) AS cosim
			FROM stylometric_profiles s2
			JOIN actors a2 ON a2.actor_id = s2.actor_id
			WHERE s2.actor_id <> $2
			  AND s2.feature_vector <=> $1::vector <= 0.60
			ORDER BY s2.feature_vector <=> $1::vector
			LIMIT $3`, vectorText(vec), actorID, k)
		if qerr == nil {
			defer rows.Close()
			for rows.Next() {
				var handle string
				var w, o int
				var cosim float64
				if rows.Scan(&handle, &w, &o, &cosim) == nil {
					lookalikes = append(lookalikes, gin.H{
						"handle":       handle,
						"sample_words": w,
						"observations": o,
						"cosine":       cosim,
					})
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"sample_words":    words,
			"observations":    obs,
			"last_active_at":  lastActive.Format(time.RFC3339),
			"median_gap_min":  gapOrZero(gap),
			"hour_histogram":  hour,
			"day_histogram":   day,
			"tempo_histogram": tempo,
			// Raw 384-d lexical fingerprint. dashboard-api projects it (with the
			// histograms above) into the fixed-width determinisic stylokey that
			// drives the graph layout — no other service can read the DB, so this
			// is the single passthrough point for the vector.
			"feature_vector": vec,
			"vector_present": len(vec) == 384,
			"sample_gate":    words >= MinStylometrySampleWords,
		},
		"lookalikes": lookalikes,
	})
}

// MinStylometrySampleWords mirrors the correlation engine's statistical gate.
const MinStylometrySampleWords = 150

func gapOrZero(g sql.NullInt64) int {
	if g.Valid {
		return int(g.Int64)
	}
	return 0
}

// persistAuthorProfile writes an author's behavioral fingerprint inside the
// CreateScan transaction: raw post-time observations first (as audit evidence),
// then a weighted-mean merge into stylometric_profiles. The weighted mean is
// computed in Go so the 384-d vector stays exact without pgvector arithmetic
// or a growing running-sum column. Zero-trust is preserved — this is still
// storage-api's own DB transaction.
func persistAuthorProfile(ctx context.Context, tx pgx.Tx, actorID string, a AuthorProfile, fetchedAt time.Time, sourceURL, provenance string) error {
	// 1. Raw observations — the "message timing" evidence this profile is built
	//    from. Kept per-source so hour histograms can be recomputed per market.
	for _, raw := range a.PostTimes {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			continue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO actor_observations (actor_id, observed_at, source_url, provenance_tag)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (actor_id, observed_at, source_url) DO NOTHING`,
			actorID, t.UTC(), sourceURL, provenance); err != nil {
			return err
		}
	}

	// 2. Current profile state (if any).
	var (
		curWords  int
		curVecTxt sql.NullString
		curHour   []float32
		curDay    []float32
		curObs    int
	)
	err := tx.QueryRow(ctx, `
		SELECT sample_word_count,
		       feature_vector::text,
		       COALESCE(hour_histogram, ARRAY[]::real[]),
		       COALESCE(day_histogram, ARRAY[]::real[]),
		       observations_count
		FROM stylometric_profiles
		WHERE actor_id = $1`, actorID).
		Scan(&curWords, &curVecTxt, &curHour, &curDay, &curObs)
	exists := err == nil
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	// 3. Weighted-mean vector merge (style is weighted by sample size so a
	//    single long page never drowns accumulated short posts, and vice versa).
	validVec := len(a.Vector) == 384 && a.SampleWords > 0
	if !validVec && !exists {
		// No style sample AND no pre-existing row: without a 384-d vector the
		// NOT NULL feature_vector cannot be satisfied — temporal evidence was
		// already recorded as observations above, so nothing else to do.
		return nil
	}

	var mergeVec []float64
	if validVec && exists {
		curVec := parseVectorText(curVecTxt.String)
		if len(curVec) == 384 {
			total := curWords + a.SampleWords
			mergeVec = make([]float64, 384)
			for i := range a.Vector {
				mergeVec[i] = (curVec[i]*float64(curWords) + a.Vector[i]*float64(a.SampleWords)) / float64(total)
			}
		} else {
			mergeVec = a.Vector
		}
	}

	// 4. Temporal merge — element-wise histogram accumulation.
	hour := mergeHistograms(curHour, a.HourHistogram)
	day := mergeHistograms(curDay, a.DayHistogram)

	// Inter-arrival tempo histogram is recomputed from evidence (not streamed)
	// to avoid compounding bin drift across uneven crawls.
	tempo := interArrivalTempoHistogramMinutes(ctx, tx, actorID)

	// 5. Recompute median inter-post gap from stored observations (true median
	//    over evidence, not a streaming approximation).
	var gap *int
	if mg := medianGapMinutes(ctx, tx, actorID); mg != nil {
		g := *mg
		gap = &g
	}

	newObs := curObs + len(a.PostTimes)

	if !exists {
		_, err = tx.Exec(ctx, `
			INSERT INTO stylometric_profiles
				(actor_id, sample_word_count, last_active_at, hour_histogram,
				 day_histogram, tempo_histogram, median_gap_minutes, observations_count,
				 provenance, feature_vector)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::vector)`,
			actorID, a.SampleWords, fetchedAt, hour, day, tempo, gap, newObs, provenance, vectorText(a.Vector))
		return err
	}

	totalWords := curWords + a.SampleWords
	if !validVec {
		// Profile exists but this page contributed no style sample — update
		// only the temporal family so the 384-d fingerprint is never clobbered.
		_, err = tx.Exec(ctx, `
			UPDATE stylometric_profiles SET
				sample_word_count   = $2,
				last_active_at      = GREATEST(last_active_at, $3),
				hour_histogram      = $4,
				day_histogram       = $5,
				tempo_histogram     = $6,
				median_gap_minutes  = $7,
				observations_count  = $8,
				provenance          = COALESCE($9, provenance)
			WHERE actor_id = $1`,
			actorID, totalWords, fetchedAt, hour, day, tempo, gap, newObs, provenance)
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE stylometric_profiles SET
			sample_word_count   = $2,
			last_active_at      = GREATEST(last_active_at, $3),
			hour_histogram      = $4,
			day_histogram       = $5,
			tempo_histogram     = $6,
			median_gap_minutes  = $7,
			observations_count  = $8,
			provenance          = COALESCE($9, provenance),
			feature_vector      = $10::vector
		WHERE actor_id = $1`,
		actorID, totalWords, fetchedAt, hour, day, tempo, gap, newObs, provenance, vectorText(mergeVec))
	return err
}

// mergeHistograms adds a new count histogram onto the stored one (same length).
func mergeHistograms(stored []float32, added []float64) []float32 {
	out := make([]float32, len(stored))
	copy(out, stored)
	for i, v := range added {
		if i < len(out) {
			out[i] += float32(v)
		}
	}
	return out
}

// medianGapMinutes returns the median inter-post gap (minutes) over the most
// recent 200 observations, or nil when there are fewer than 2.
func interArrivalTempoHistogramMinutes(ctx context.Context, tx pgx.Tx, actorID string) []float64 {
	const bins = 12
	// Hybrid binning (Option 3): 8 linear bins of 5 minutes from 0-40,
	// then 4 long bins: 40-60m, 1-6h, 6-12h, 12-24h.
	// We store counts, then caller normalizes in correlation-engine.
	h := make([]float64, bins)

	rows, err := tx.Query(ctx, `
		SELECT observed_at FROM actor_observations
		WHERE actor_id = $1 ORDER BY observed_at DESC LIMIT 200`, actorID)
	if err != nil {
		return h
	}
	defer rows.Close()

	var times []time.Time
	for rows.Next() {
		var t time.Time
		if rows.Scan(&t) == nil {
			times = append(times, t)
		}
	}
	if len(times) < 2 {
		return h
	}

	// Sort ascending to compute forward gaps.
	sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })

	// Inter-arrival gaps in minutes.
	for i := 1; i < len(times); i++ {
		gapMin := times[i].Sub(times[i-1]).Minutes()
		if gapMin < 0 {
			continue
		}
		if gapMin < 0.01 {
			gapMin = 0
		}

		switch {
		case gapMin < 40:
			// 0-40 split into 8 bins of width 5.
			idx := int(gapMin / 5.0) // 0..7
			if idx < 0 {
				idx = 0
			}
			if idx > 7 {
				idx = 7
			}
			h[idx] += 1
		case gapMin < 60:
			h[8] += 1
		case gapMin < 360: // 1-6h
			h[9] += 1
		case gapMin < 720: // 6-12h
			h[10] += 1
		case gapMin < 1440: // 12-24h
			h[11] += 1
		default:
			// >24h: ignore (doesn't help cosine similarity shape).
		}
	}

	return h
}

// medianGapMinutes returns the median inter-post gap (minutes) over the most
// recent 200 observations, or nil when there are fewer than 2.
func medianGapMinutes(ctx context.Context, tx pgx.Tx, actorID string) *int {
	rows, err := tx.Query(ctx, `
		SELECT observed_at FROM actor_observations
		WHERE actor_id = $1 ORDER BY observed_at DESC LIMIT 200`, actorID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var times []time.Time
	for rows.Next() {
		var t time.Time
		if rows.Scan(&t) == nil {
			times = append(times, t)
		}
	}
	if len(times) < 2 {
		return nil
	}
	sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
	gaps := make([]int, 0, len(times)-1)
	for i := 1; i < len(times); i++ {
		gaps = append(gaps, int(times[i].Sub(times[i-1]).Minutes()))
	}
	sort.Ints(gaps)
	mid := len(gaps) / 2
	if len(gaps)%2 == 1 {
		v := gaps[mid]
		return &v
	}
	v := (gaps[mid-1] + gaps[mid]) / 2
	return &v
}

// parseVectorText parses pgvector's text form ("[0.1,0.2,...]") into floats.
func parseVectorText(s string) []float64 {
	s = strings.TrimSpace(strings.Trim(s, "[]"))
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]float64, len(parts))
	for i, p := range parts {
		f, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return nil
		}
		out[i] = f
	}
	return out
}

// vectorText renders a float vector into pgvector's text form.
func vectorText(v []float64) string {
	if len(v) == 0 {
		return "[]"
	}
	parts := make([]string, len(v))
	for i, f := range v {
		parts[i] = strconv.FormatFloat(f, 'g', -1, 64)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
