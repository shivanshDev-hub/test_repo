// Negative-control handlers (storage-api) backing the LeftControlRail
// "NEGATIVE CONTROL" console.  These are the ONLY tables that persist the
// analyst's frequency/static damping policy — every other service reads
// through this HTTP boundary (zero-trust architecture).
package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ModeFrequency = "FREQUENCY"
	ModeStatic    = "STATIC"
)

// NegativeControlRulePayload accepted by POST /negative-control.
type NegativeControlRulePayload struct {
	Mode         string `json:"mode"`
	ArtifactType string `json:"artifact_type"`
	FrequencyCap int    `json:"frequency_cap"`
	Enabled      *bool  `json:"enabled"`
}

func (p NegativeControlRulePayload) valid() bool {
	mode := strings.ToUpper(p.Mode)
	return (mode == ModeFrequency || mode == ModeStatic) && strings.TrimSpace(p.ArtifactType) != ""
}

// NegativeControlArtifactPayload accepted by POST /negative-control/artifacts.
type NegativeControlArtifactPayload struct {
	ArtifactType string `json:"artifact_type"`
	Value        string `json:"value"`
	RuleID       string `json:"rule_id,omitempty"`
}

// ── GET /negative-control ─────────────────────────────────────────────────────
// Returns every policy rule (with live artifact counts) newest first.
func (h *Handler) ListNegativeControl(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
	defer cancel()

	rows, err := h.pool.Query(ctx, `
		SELECT r.rule_id, r.mode, r.artifact_type, r.frequency_cap, r.enabled,
		       r.created_at, r.updated_at,
		       COALESCE(counts.ac, 0)::int,
		       COALESCE(counts.mx, 0)::int
		FROM negative_control_rules r
		LEFT JOIN (
			SELECT rule_id, count(*)::int AS ac,
			       COALESCE(max(occurrence_count), 0)::int AS mx
			FROM negative_control_artifacts GROUP BY rule_id
		) counts ON counts.rule_id = r.rule_id
		ORDER BY r.created_at DESC, r.artifact_type ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query rules: " + err.Error()})
		return
	}
	defer rows.Close()

	type rr struct {
		RuleID        string `json:"rule_id"`
		Mode          string `json:"mode"`
		ArtifactType  string `json:"artifact_type"`
		FrequencyCap  int    `json:"frequency_cap"`
		Enabled       bool   `json:"enabled"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
		ArtifactCount int    `json:"artifact_count"`
		MaxOccurrence int    `json:"max_occurrences"`
	}
	out := []rr{}
	for rows.Next() {
		var r rr
		if err := rows.Scan(&r.RuleID, &r.Mode, &r.ArtifactType, &r.FrequencyCap,
			&r.Enabled, &r.CreatedAt, &r.UpdatedAt, &r.ArtifactCount, &r.MaxOccurrence); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan rule: " + err.Error()})
			return
		}
		out = append(out, r)
	}
	c.JSON(http.StatusOK, gin.H{"rules": out})
}

// ── POST /negative-control ────────────────────────────────────────────────────
func (h *Handler) CreateNegativeControlRule(c *gin.Context) {
	var p NegativeControlRulePayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !p.valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mode FREQUENCY|STATIC and artifact_type required"})
		return
	}
	mode := strings.ToUpper(p.Mode)
	enabled := true
	if p.Enabled != nil {
		enabled = *p.Enabled
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
	defer cancel()

	var ruleID, cm, at, ca, ua string
	var fc int
	var en bool
	if err := h.pool.QueryRow(ctx, `
		INSERT INTO negative_control_rules (mode,artifact_type,frequency_cap,enabled)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (mode,artifact_type) DO UPDATE
			SET frequency_cap=EXCLUDED.frequency_cap,enabled=EXCLUDED.enabled,updated_at=NOW()
		RETURNING rule_id,mode,artifact_type,frequency_cap,enabled,created_at,updated_at`,
		mode, strings.TrimSpace(p.ArtifactType), p.FrequencyCap, enabled,
	).Scan(&ruleID, &cm, &at, &fc, &en, &ca, &ua); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"rule_id": ruleID, "mode": cm, "artifact_type": at,
		"frequency_cap": fc, "enabled": en, "created_at": ca, "updated_at": ua,
	})
}

// ── DELETE /negative-control/:id ──────────────────────────────────────────────
func (h *Handler) DeleteNegativeControlRule(c *gin.Context) {
	id := c.Param("id")
	if !validUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule_id"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
	defer cancel()
	if _, err := h.pool.Exec(ctx,
		`DELETE FROM negative_control_rules WHERE rule_id=$1`, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted", "rule_id": id})
}

// ── POST /negative-control/artifacts ──────────────────────────────────────────
func (h *Handler) AddNegativeControlArtifact(c *gin.Context) {
	var p NegativeControlArtifactPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(p.ArtifactType) == "" || strings.TrimSpace(p.Value) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artifact_type and value required"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
	defer cancel()

	// Auto-link to an active rule for this type if caller omits rule_id.
	var ruleID *string
	if p.RuleID != "" {
		_ = h.pool.QueryRow(ctx,
			`SELECT rule_id FROM negative_control_rules WHERE rule_id=$1`, p.RuleID).Scan(&ruleID)
	} else {
		_ = h.pool.QueryRow(ctx,
			`SELECT rule_id FROM negative_control_rules
			 WHERE artifact_type=$1 AND enabled=TRUE LIMIT 1`,
			p.ArtifactType).Scan(&ruleID)
	}

	var aid, ca string
	var cnt int
	if err := h.pool.QueryRow(ctx, `
		INSERT INTO negative_control_artifacts (rule_id,artifact_type,value,occurrence_count)
		VALUES ($1,$2,$3,1)
		ON CONFLICT (rule_id,artifact_type,value) DO UPDATE
			SET occurrence_count=negative_control_artifacts.occurrence_count+1, last_seen_at=NOW()
		RETURNING artifact_id, occurrence_count, first_seen_at`,
		ruleID, strings.TrimSpace(p.ArtifactType), strings.TrimSpace(p.Value),
	).Scan(&aid, &cnt, &ca); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"artifact_id": aid, "rule_id": ruleID, "artifact_type": strings.TrimSpace(p.ArtifactType),
		"value": strings.TrimSpace(p.Value), "occurrence_count": cnt, "first_seen_at": ca,
	})
}

// ── GET /artifacts?type=&mode=&limit= ────────────────────────────────────────
func (h *Handler) ListNegativeControlArtifacts(c *gin.Context) {
	limit := 200
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
	defer cancel()

	q := `SELECT a.artifact_id, COALESCE(a.rule_id::text,''), a.artifact_type,
	       a.value, a.occurrence_count, a.first_seen_at, a.last_seen_at,
	       COALESCE(r.mode,''), COALESCE(r.enabled,FALSE)
	      FROM negative_control_artifacts a
	      LEFT JOIN negative_control_rules r ON r.rule_id=a.rule_id`
	args := []interface{}{}
	argN := 1
	wh := []string{}

	if t := c.Query("type"); t != "" {
		wh = append(wh, "a.artifact_type=$"+strconv.Itoa(argN))
		args = append(args, t)
		argN++
	}
	if m := strings.ToUpper(c.Query("mode")); m == ModeFrequency || m == ModeStatic {
		wh = append(wh, "r.mode=$"+strconv.Itoa(argN))
		args = append(args, m)
		argN++
	}
	if len(wh) > 0 {
		q += " WHERE " + strings.Join(wh, " AND ")
	}
	q += " ORDER BY a.occurrence_count DESC, a.last_seen_at DESC LIMIT $" + strconv.Itoa(argN)
	args = append(args, limit)

	rows, err := h.pool.Query(ctx, q, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type ar struct {
		ArtifactID      string `json:"artifact_id"`
		RuleID          string `json:"rule_id"`
		ArtifactType    string `json:"artifact_type"`
		Value           string `json:"value"`
		OccurrenceCount int    `json:"occurrence_count"`
		FirstSeenAt     string `json:"first_seen_at"`
		LastSeenAt      string `json:"last_seen_at"`
		Mode            string `json:"mode"`
		Enabled         bool   `json:"enabled"`
	}
	out := []ar{}
	for rows.Next() {
		var a ar
		if err := rows.Scan(&a.ArtifactID, &a.RuleID, &a.ArtifactType, &a.Value,
			&a.OccurrenceCount, &a.FirstSeenAt, &a.LastSeenAt, &a.Mode, &a.Enabled); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		out = append(out, a)
	}
	c.JSON(http.StatusOK, gin.H{"artifacts": out})
}