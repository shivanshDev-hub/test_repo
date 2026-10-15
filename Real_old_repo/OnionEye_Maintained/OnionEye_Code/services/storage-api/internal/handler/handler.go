// Package handler implements the storage-api HTTP endpoints.
//
// This is the ONLY service in the system that holds database credentials —
// every other service (crawler-worker, content-shield, correlation-engine,
// dashboard-api) persists and reads through this HTTP boundary, matching the
// zero-trust architecture.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func logErr(msg string, err error) {
	if err == nil {
		return
	}
	log.Printf("%s: %v", msg, err)
}

// Handler holds the pgx pool and exposes HTTP handlers.
type Handler struct {
	pool *pgxpool.Pool
}

// DeleteNodePayload requests deletion of a graph node.
// node_type is the UI node type (actor/pgp/wallet/onion/infra/email).
// node_id is the UI graph node id string.
type DeleteNodePayload struct {
	NodeID   string `json:"node_id"`
	NodeType string `json:"node_type"`
}

// DeleteAllResponse provides a summary of how many rows were affected.
type DeleteAllResponse struct {
	DeletedPersonaLinks int64 `json:"deleted_persona_links"`
	DeletedActors       int64 `json:"deleted_actors"`
	DeletedIdentifiers  int64 `json:"deleted_identifiers"`
	DeletedInfra        int64 `json:"deleted_infra_indicators"`
}

// New builds a Handler with the given connection pool.
func New(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

// ── Payload types (mirror the clients in crawler-worker / content-shield) ────

// ScanPayload is the extraction output persisted by the crawler-worker
// (POST /scans). Field names match clients.ScanPayload in crawler-worker.
type ScanPayload struct {
	JobID        string    `json:"job_id"`
	SourceURL    string    `json:"source_url"`
	Provenance   string    `json:"provenance_tag"`
	Extracted    Extracted `json:"extracted"`
	BodyRef      string    `json:"body_ref"`
	TorCircuitID string    `json:"tor_circuit_id"`
	FetchedAt    string    `json:"fetched_at"`
	Terms        []string  `json:"terms"`
}

// Extracted holds the regex-extracted fields from a crawled page.
type Extracted struct {
	Handles         []string         `json:"handles"`
	PGPKeys         []string         `json:"pgp_keys"`
	Wallets         []Wallet         `json:"wallets"`
	InfraIndicators []InfraIndicator `json:"infra_indicators"`
	Emails          []string         `json:"emails"`
	Links           []Link           `json:"links"`
	Images          []Image          `json:"images"`
	Timestamps      []string         `json:"timestamps"`
	Authors         []AuthorProfile  `json:"authors,omitempty"`
}

// AuthorProfile mirrors the crawler-worker behavioral fingerprint shape. The
// 384-d Vector persists into stylometric_profiles.feature_vector (pgvector);
// the histograms and post times drive the temporal ("human") vector family.
type AuthorProfile struct {
	Handle        string    `json:"handle"`
	Vector        []float64 `json:"vector,omitempty"`
	SampleWords   int       `json:"sample_words"`
	HourHistogram []float64 `json:"hour_histogram,omitempty"` // 24 UTC bins (counts)
	DayHistogram  []float64 `json:"day_histogram,omitempty"`  // 7 UTC bins (counts)
	PostTimes     []string  `json:"post_times,omitempty"`     // RFC3339
}

// Wallet is a crypto address extracted from a page.
type Wallet struct {
	Currency string `json:"currency"`
	Address  string `json:"address"`
}

// InfraIndicator is a hosting/infra signal extracted from a page.
type InfraIndicator struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Link is an outbound URL extracted from a page.
type Link struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// Image is an image reference extracted from a page.
type Image struct {
	Src  string `json:"src"`
	Type string `json:"type"`
}

// FlagPayload is the content_flags row shape posted by the content-shield
// (POST /content-flags). Field names match storage.FlagPayload.
type FlagPayload struct {
	JobID       string `json:"job_id"`
	SourceURL   string `json:"source_url"`
	Category    string `json:"category"`
	ActionTaken string `json:"action_taken"`
	HashMatch   bool   `json:"hash_match"`
	FlaggedAt   string `json:"flagged_at"`
}

// ActorPayload is an upsert for the actors table (POST /actors).
type ActorPayload struct {
	ActorID       string `json:"actor_id"`
	PrimaryHandle string `json:"primary_handle"`
	Category      string `json:"category"`
}

// IdentifierPayload is an upsert for the identifiers table (POST /identifiers).
type IdentifierPayload struct {
	ActorID       string `json:"actor_id"`
	Kind          string `json:"kind"`
	Value         string `json:"value"`
	Source        string `json:"source"`
	ProvenanceTag string `json:"provenance_tag"`
}

// PersonaLinkPayload is an upsert for the persona_links table (POST /persona-links).
type PersonaLinkPayload struct {
	LinkID                string          `json:"link_id"`
	ActorIDA              string          `json:"actor_id_a"`
	ActorIDB              string          `json:"actor_id_b"`
	Basis                 string          `json:"basis"`
	Confidence            float64         `json:"confidence"`
	NegativeControlPassed bool            `json:"negative_control_passed"`
	NCState               string          `json:"nc_state"`
	NCPenaltyScore        float64         `json:"nc_penalty_score"`
	NCAuditLog            json.RawMessage `json:"nc_audit_log"`
}

// ── Helpers ──────────────────────────────────────────────────────────────────

var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// validUUID reports whether s is a well-formed UUID string.
func validUUID(s string) bool { return uuidRe.MatchString(s) }

// parseTime parses an RFC3339 timestamp, falling back to now.
func parseTime(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Now().UTC()
}

// ── POST /scans ──────────────────────────────────────────────────────────────

// CreateScan persists a crawler scan plus its extracted actors/identifiers.
func (h *Handler) CreateScan(c *gin.Context) {
	var p ScanPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "begin tx: " + err.Error()})
		return
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op after commit

	// scan_id: use the job_id when it's a valid UUID, else let the DB generate.
	scanID := p.JobID
	if !validUUID(scanID) {
		scanID = "00000000-0000-0000-0000-000000000000" // placeholder; replaced below
	}

	fetchedAt := parseTime(p.FetchedAt)

	// Insert the scan row (partitioned by started_at). ON CONFLICT makes
	// retries idempotent — a re-delivered job won't double-insert. The
	// conflict target must include the partition key (started_at) because
	// scans has a composite PK (scan_id, started_at).
	if validUUID(p.JobID) {
		_, err = tx.Exec(ctx, `
			INSERT INTO scans (scan_id, requested_by, targets, status, started_at, completed_at)
			VALUES ($1, $2, $3, 'completed', $4, now())
			ON CONFLICT (scan_id, started_at) DO NOTHING`,
			scanID, "crawler-worker", []string{p.SourceURL}, fetchedAt)
	} else {
		err = tx.QueryRow(ctx, `
			INSERT INTO scans (scan_id, requested_by, targets, status, started_at, completed_at)
			VALUES (gen_random_uuid(), $1, $2, 'completed', $3, now())
			RETURNING scan_id`,
			"crawler-worker", []string{p.SourceURL}, fetchedAt).Scan(&scanID)
	}
	if err != nil {
		logErr("CreateScan insert scan failed", err)
		logErr("CreateScan insert scan failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert scan: " + err.Error()})
		return
	}

	// Upsert an actor per unique handle, then attach identifiers.
	firstActorID := ""
	allActorIDs := make([]string, 0)
	actorByHandle := map[string]string{}
	seen := map[string]bool{}
	for _, handle := range p.Extracted.Handles {
		if handle == "" || seen[handle] {
			continue
		}
		seen[handle] = true
		actorID, err := upsertActor(ctx, tx, handle, p.Provenance)
		if err != nil {
			logErr("CreateScan upsert actor failed", err)
			logErr("CreateScan upsert actor failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upsert actor: " + err.Error()})
			return
		}
		allActorIDs = append(allActorIDs, actorID)
		if firstActorID == "" {
			firstActorID = actorID
		}
		actorByHandle[strings.ToLower(handle)] = actorID
	}

	// Fallback mechanism: when the page yielded NO extractable handles at all,
	// promote the investigation terms into actor handles so the page still
	// produces a graph node. Fallback actors are visually distinct from real
	// handles (TERM: prefix + TERM_FALLBACK category) — they only exist to
	// anchor a scan when extraction found nothing, and are clearly labelled as
	// such in the UI/graph.
	if len(allActorIDs) == 0 {
		for _, term := range p.Terms {
			fbHandle := termFallbackHandle(term)
			if fbHandle == "" {
				continue
			}
			actorID, err := upsertActor(ctx, tx, fbHandle, "TERM_FALLBACK")
			if err != nil {
				logErr("CreateScan upsert fallback actor failed", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "upsert fallback actor: " + err.Error()})
				return
			}
			allActorIDs = append(allActorIDs, actorID)
			if firstActorID == "" {
				firstActorID = actorID
			}
		}
	}

	// Identifiers attach to every actor created for this page (co-occurrence
	// attribution). All identifiers extracted from a page belong to that
	// page's persona(s), and attaching them to each actor lets the correlation
	// engine create persona_links between any actors that share the same
	// identifier value (email, PGP key, wallet). NULL actor_id is still
	// allowed when the page produced no actors at all.
	attach := func(kind, value string) error {
		if value == "" {
			return nil
		}
		idList := allActorIDs
		if len(idList) == 0 {
			// No actors on the page: keep the identifier orphaned (NULL actor_id).
			idList = []string{firstActorID}
		}
		for _, aID := range idList {
			if err := upsertIdentifier(ctx, tx, aID, kind, value, p.Provenance, fetchedAt); err != nil {
				return err
			}
		}
		return nil
	}

	for _, key := range p.Extracted.PGPKeys {
		if err := attach("pgp_key", key); err != nil {
			logErr("CreateScan upsert pgp failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upsert pgp: " + err.Error()})
			return
		}
	}
	for _, w := range p.Extracted.Wallets {
		if err := attach("wallet:"+w.Currency, w.Address); err != nil {
			logErr("CreateScan upsert wallet failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upsert wallet: " + err.Error()})
			return
		}
	}
	for _, e := range p.Extracted.Emails {
		if err := attach("email", e); err != nil {
			logErr("CreateScan upsert email failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upsert email: " + err.Error()})
			return
		}
	}

	// Behavioral fingerprints (stylometry + temporal "human" vectors) for each
	// attributed handle on the page. Persisted in the same transaction as the
	// actor so a failed write never leaves dangling observations.
	for _, author := range p.Extracted.Authors {
		if author.Handle == "" {
			continue
		}
		actorID := actorByHandle[strings.ToLower(author.Handle)]
		if actorID == "" || !validUUID(actorID) {
			continue
		}
		if err := persistAuthorProfile(ctx, tx, actorID, author, fetchedAt, p.SourceURL, p.Provenance); err != nil {
			logErr("CreateScan persist author profile failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "persist behavioral profile: " + err.Error()})
			return
		}
	}

	// Actor-scoped inserts: infra/links/images.
	// actor_id is a UUID column; firstActorID may be empty if no handles were
	// extracted / upserted.
	attachable := validUUID(firstActorID)

	// Actor-scoped inserts: infra/links/images.
	// actor_id is a UUID column; firstActorID may be empty if no handles were
	// extracted / upserted.
	if attachable {
		for _, ind := range p.Extracted.InfraIndicators {
			// Hard guard: infra_indicators.actor_id is UUID; avoid passing "".
			if firstActorID == "" || !validUUID(firstActorID) {
				continue
			}
			if ind.Value == "" {
				continue
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO infra_indicators (indicator_id, actor_id, indicator_type, onion_host, evidence, detected_at)
				VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
				ON CONFLICT (actor_id, indicator_type, onion_host) DO NOTHING`,
				firstActorID, ind.Type, ind.Value, p.Provenance, fetchedAt)
			if err != nil {
				logErr("CreateScan insert infra failed", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "insert infra: " + err.Error()})
				return
			}
		}

		for _, l := range p.Extracted.Links {
			if l.URL == "" {
				continue
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO extracted_links (link_ref_id, actor_id, url, link_type, evidence, detected_at)
				VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`,
				firstActorID, l.URL, l.Type, p.Provenance, fetchedAt)
			if err != nil {
				logErr("CreateScan insert link failed", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "insert link: " + err.Error()})
				return
			}
		}

		for _, img := range p.Extracted.Images {
			if img.Src == "" {
				continue
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO image_hashes (image_ref_id, actor_id, perceptual_hash, evidence, detected_at)
				VALUES (gen_random_uuid(), $1, $2, $3, $4)`,
				firstActorID, img.Src, p.Provenance, fetchedAt)
			if err != nil {
				logErr("CreateScan insert image failed", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "insert image: " + err.Error()})
				return
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logErr("CreateScan commit failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit: " + err.Error()})
		return
	}

	// Auto-correlate: create persona_links for actors sharing the same
	// identifier value (wallet, email, PGP key, etc.).  Runs after commit
	// so we correlate against fully persisted data.
	for _, aID := range allActorIDs {
		if aID != "" && validUUID(aID) {
			h.correlateIdentifiers(ctx, aID)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"scan_id": scanID, "status": "completed"})
}

// correlateIdentifiers creates persona_links for actors that share the same
// identifier value (wallet address, email, PGP key, etc.).  It queries for any
// identifier under actorID whose kind+value also appears under a *different*
// actor, then upserts a persona_link between them with a confidence score
// appropriate to the identifier type.
func (h *Handler) correlateIdentifiers(ctx context.Context, actorID string) {
	// Use a generous timeout — this runs after the main transaction commits.
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	rows, err := h.pool.Query(ctx, `
		SELECT DISTINCT
			LEAST(i1.actor_id, i2.actor_id)  AS a_id,
			GREATEST(i1.actor_id, i2.actor_id) AS b_id,
			i1.kind,
			LEFT(i1.value, 40)               AS value_preview
		FROM identifiers i1
		JOIN identifiers i2
			ON  i1.kind  = i2.kind
			AND i1.value = i2.value
			AND i1.actor_id < i2.actor_id
		WHERE i1.actor_id = $1 OR i2.actor_id = $1
	`, actorID)
	if err != nil {
		logErr("correlateIdentifiers query", err)
		return
	}
	defer rows.Close()

	linksCreated := 0
	for rows.Next() {
		var aID, bID, kind, preview string
		if err := rows.Scan(&aID, &bID, &kind, &preview); err != nil {
			logErr("correlateIdentifiers scan", err)
			continue
		}

		// Confidence heuristic by identifier type.
		confidence := 0.55
		switch {
		case strings.Contains(kind, "pgp"):
			confidence = 0.90
		case strings.Contains(kind, "wallet"):
			confidence = 0.80
		case kind == "email":
			confidence = 0.65
		}

		basis := fmt.Sprintf("Shared %s: %s", kind, preview)

		_, err = h.pool.Exec(ctx, `
			INSERT INTO persona_links
				(link_id, actor_id_a, actor_id_b, basis, confidence,
				 negative_control_passed, nc_state, created_at)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, true, 'PASSED', now())
			ON CONFLICT (actor_id_a, actor_id_b, basis) DO NOTHING
		`, aID, bID, basis, confidence)
		if err != nil {
			logErr("correlateIdentifiers insert", err)
		} else {
			linksCreated++
		}
	}
	if linksCreated > 0 {
		log.Printf("[correlate] actor %s: created %d persona_link(s)", actorID, linksCreated)
	}
}

// ── POST /content-flags ──────────────────────────────────────────────────────

// CreateFlag persists a content-safety flag from the content-shield.
func (h *Handler) CreateFlag(c *gin.Context) {
	var p FlagPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	flaggedAt := parseTime(p.FlaggedAt)

	// Resolve scan_id from job_id when possible (nullable otherwise).
	var scanID *string
	if validUUID(p.JobID) {
		scanID = &p.JobID
	}

	_, err := h.pool.Exec(ctx, `
		INSERT INTO content_flags (flag_id, scan_id, source_url, category, action_taken, flagged_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`,
		scanID, p.SourceURL, p.Category, p.ActionTaken, flaggedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert flag: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "flagged"})
}

// ── POST /actors ─────────────────────────────────────────────────────────────

// UpsertActor creates or updates an actor row.
func (h *Handler) UpsertActor(c *gin.Context) {
	var p ActorPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}
	if p.PrimaryHandle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "primary_handle is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	actorID, err := upsertActor(ctx, h.pool, p.PrimaryHandle, p.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upsert actor: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"actor_id": actorID})
}

// ── POST /identifiers ────────────────────────────────────────────────────────

// UpsertIdentifier creates or updates an identifier row.
func (h *Handler) UpsertIdentifier(c *gin.Context) {
	var p IdentifierPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}
	if p.Value == "" || p.Kind == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kind and value are required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	now := time.Now().UTC()
	if err := upsertIdentifier(ctx, h.pool, p.ActorID, p.Kind, p.Value, p.ProvenanceTag, now); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upsert identifier: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "upserted"})
}

// ── POST /persona-links ──────────────────────────────────────────────────────

// UpsertPersonaLink creates or updates a correlation edge.
func (h *Handler) UpsertPersonaLink(c *gin.Context) {
	var p PersonaLinkPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}
	if p.ActorIDA == "" || p.ActorIDB == "" || p.Basis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "actor_id_a, actor_id_b, basis are required"})
		return
	}
	if p.Confidence < 0 || p.Confidence > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "confidence must be in [0,1]"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	ncState := p.NCState
	if ncState == "" {
		ncState = "INDETERMINATE"
	}
	audit := p.NCAuditLog
	if len(audit) == 0 {
		audit = json.RawMessage(`{}`)
	}

	_, err := h.pool.Exec(ctx, `
		INSERT INTO persona_links
			(link_id, actor_id_a, actor_id_b, basis, confidence, negative_control_passed,
			 nc_state, nc_penalty_score, nc_audit_log, created_at)
		VALUES
			(COALESCE($1::uuid, gen_random_uuid()), $2, $3, $4, $5, $6, $7, $8, $9, now())
		ON CONFLICT (actor_id_a, actor_id_b, basis) DO UPDATE SET
			confidence = EXCLUDED.confidence,
			negative_control_passed = EXCLUDED.negative_control_passed,
			nc_state = EXCLUDED.nc_state,
			nc_penalty_score = EXCLUDED.nc_penalty_score,
			nc_audit_log = EXCLUDED.nc_audit_log`,
		nullableUUID(p.LinkID), p.ActorIDA, p.ActorIDB, p.Basis,
		p.Confidence, p.NegativeControlPassed, ncState, p.NCPenaltyScore, audit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upsert persona link: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "upserted"})
}

// ── GET /actors/:id ──────────────────────────────────────────────────────────

// GetActor returns an actor profile with its identifiers and infra indicators.
func (h *Handler) DeleteNode(c *gin.Context) {
	var p DeleteNodePayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}
	p.NodeID = strings.TrimSpace(p.NodeID)
	p.NodeType = strings.TrimSpace(p.NodeType)
	if p.NodeID == "" || p.NodeType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "node_id and node_type are required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// ── Actor deletion (node_id is a UUID) ────────────────────────────────────
	if validUUID(p.NodeID) {
		actorID := p.NodeID
		tx, err := h.pool.Begin(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "begin tx: " + err.Error()})
			return
		}
		defer func() { _ = tx.Rollback(ctx) }()

		res, err := tx.Exec(ctx, `DELETE FROM persona_links WHERE actor_id_a = $1 OR actor_id_b = $1`, actorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete persona links: " + err.Error()})
			return
		}
		linksDeleted := res.RowsAffected()

		resI, _ := tx.Exec(ctx, `DELETE FROM identifiers WHERE actor_id = $1`, actorID)
		deletedIdentifiers := resI.RowsAffected()

		resInf, _ := tx.Exec(ctx, `DELETE FROM infra_indicators WHERE actor_id = $1`, actorID)
		deletedInfra := resInf.RowsAffected()

		resLinks, _ := tx.Exec(ctx, `DELETE FROM extracted_links WHERE actor_id = $1`, actorID)
		deletedExtractedLinks := resLinks.RowsAffected()

		resHashes, _ := tx.Exec(ctx, `DELETE FROM image_hashes WHERE actor_id = $1`, actorID)
		deletedImageHashes := resHashes.RowsAffected()

		resFlags, _ := tx.Exec(ctx, `DELETE FROM content_flags WHERE actor_id = $1`, actorID)
		deletedContentFlags := resFlags.RowsAffected()

		resA, err := tx.Exec(ctx, `DELETE FROM actors WHERE actor_id = $1`, actorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete actor: " + err.Error()})
			return
		}
		actorsDeleted := resA.RowsAffected()

		if err := tx.Commit(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "commit tx: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"deleted_persona_links":    linksDeleted,
			"deleted_actors":           actorsDeleted,
			"deleted_identifiers":      deletedIdentifiers,
			"deleted_infra_indicators": deletedInfra,
			"deleted_extracted_links":  deletedExtractedLinks,
			"deleted_image_hashes":     deletedImageHashes,
			"deleted_content_flags":    deletedContentFlags,
		})
		return
	}

	// ── Non-actor node deletion (identifier / infra) ──────────────────────────
	// Graph node IDs are built as:
	//   identifier: "id_" + kind + "_" + value   (e.g. id_wallet:BTC_1A1z...)
	//   infra:      "infra_" + indicatorType + "_" + onionHost
	// We reverse-match by concatenating DB columns and comparing against the
	// node_id suffix after stripping the known prefix.
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "begin tx: " + err.Error()})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	deletedIdentifiers := int64(0)
	deletedInfra := int64(0)

	switch {
	case strings.HasPrefix(p.NodeID, "id_"):
		rest := strings.TrimPrefix(p.NodeID, "id_")
		res, err := tx.Exec(ctx, `DELETE FROM identifiers WHERE (kind || '_' || value) = $1`, rest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete identifier: " + err.Error()})
			return
		}
		deletedIdentifiers = res.RowsAffected()

	case strings.HasPrefix(p.NodeID, "infra_"):
		rest := strings.TrimPrefix(p.NodeID, "infra_")
		res, err := tx.Exec(ctx, `DELETE FROM infra_indicators WHERE (indicator_type || '_' || onion_host) = $1`, rest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete infra: " + err.Error()})
			return
		}
		deletedInfra = res.RowsAffected()

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported node id prefix for deletion"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit tx: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted_identifiers":      deletedIdentifiers,
		"deleted_infra_indicators": deletedInfra,
	})
}

func (h *Handler) DeleteAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "begin tx: " + err.Error()})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Delete persona links first (no FK from persona_links).
	resLinks, err := tx.Exec(ctx, `DELETE FROM persona_links`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete persona_links: " + err.Error()})
		return
	}
	linksDeleted := resLinks.RowsAffected()

	// Then delete all FK children before actors.
	// Our schema doesn't declare ON DELETE CASCADE for most tables, so delete explicitly.
	resIdentifiers, err := tx.Exec(ctx, `DELETE FROM identifiers`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete identifiers: " + err.Error()})
		return
	}
	identDeleted := resIdentifiers.RowsAffected()

	resInfra, err := tx.Exec(ctx, `DELETE FROM infra_indicators`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete infra_indicators: " + err.Error()})
		return
	}
	infraDeleted := resInfra.RowsAffected()

	resExtractedLinks, err := tx.Exec(ctx, `DELETE FROM extracted_links`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete extracted_links: " + err.Error()})
		return
	}
	extractedLinksDeleted := resExtractedLinks.RowsAffected()

	resImageHashes, err := tx.Exec(ctx, `DELETE FROM image_hashes`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete image_hashes: " + err.Error()})
		return
	}
	imageHashesDeleted := resImageHashes.RowsAffected()

	resContentFlags, err := tx.Exec(ctx, `DELETE FROM content_flags`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete content_flags: " + err.Error()})
		return
	}
	contentFlagsDeleted := resContentFlags.RowsAffected()

	// Finally delete actors (stylometric_profiles has ON DELETE CASCADE).
	resActors, err := tx.Exec(ctx, `DELETE FROM actors`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete actors: " + err.Error()})
		return
	}
	actorsDeleted := resActors.RowsAffected()

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit tx: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted_persona_links":    linksDeleted,
		"deleted_actors":           actorsDeleted,
		"deleted_identifiers":      identDeleted,
		"deleted_infra_indicators": infraDeleted,
		"deleted_extracted_links":  extractedLinksDeleted,
		"deleted_image_hashes":     imageHashesDeleted,
		"deleted_content_flags":    contentFlagsDeleted,
	})
}

func (h *Handler) GetActor(c *gin.Context) {
	id := c.Param("id")
	if !validUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actor_id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var actor struct {
		ActorID       string     `json:"actor_id"`
		PrimaryHandle string     `json:"primary_handle"`
		Category      *string    `json:"category"`
		LastScanDate  *time.Time `json:"last_scan_date"`
		CreatedAt     time.Time  `json:"created_at"`
	}
	err := h.pool.QueryRow(ctx, `
		SELECT actor_id, primary_handle, category, last_scan_date, created_at
		FROM actors WHERE actor_id = $1`, id).
		Scan(&actor.ActorID, &actor.PrimaryHandle, &actor.Category, &actor.LastScanDate, &actor.CreatedAt)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "actor not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query actor: " + err.Error()})
		return
	}

	// Identifiers.
	rows, err := h.pool.Query(ctx, `
		SELECT kind, value, source, provenance_tag, first_seen, last_seen
		FROM identifiers WHERE actor_id = $1 ORDER BY last_seen DESC NULLS LAST`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query identifiers: " + err.Error()})
		return
	}
	defer rows.Close()

	type ident struct {
		Kind          string     `json:"kind"`
		Value         string     `json:"value"`
		Source        *string    `json:"source"`
		ProvenanceTag *string    `json:"provenance_tag"`
		FirstSeen     *time.Time `json:"first_seen"`
		LastSeen      *time.Time `json:"last_seen"`
	}
	idents := []ident{}
	for rows.Next() {
		var i ident
		if err := rows.Scan(&i.Kind, &i.Value, &i.Source, &i.ProvenanceTag, &i.FirstSeen, &i.LastSeen); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan identifier: " + err.Error()})
			return
		}
		idents = append(idents, i)
	}

	// Infra indicators.
	infraRows, err := h.pool.Query(ctx, `
		SELECT indicator_type, onion_host, matched_clearnet_host, evidence, detected_at
		FROM infra_indicators WHERE actor_id = $1 ORDER BY detected_at DESC`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query infra: " + err.Error()})
		return
	}
	defer infraRows.Close()

	type infra struct {
		IndicatorType       string     `json:"indicator_type"`
		OnionHost           string     `json:"onion_host"`
		MatchedClearnetHost *string    `json:"matched_clearnet_host"`
		Evidence            *string    `json:"evidence"`
		DetectedAt          *time.Time `json:"detected_at"`
	}
	infras := []infra{}
	for infraRows.Next() {
		var in infra
		if err := infraRows.Scan(&in.IndicatorType, &in.OnionHost, &in.MatchedClearnetHost, &in.Evidence, &in.DetectedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan infra: " + err.Error()})
			return
		}
		infras = append(infras, in)
	}

	c.JSON(http.StatusOK, gin.H{
		"actor":       actor,
		"identifiers": idents,
		"infra":       infras,
	})
}

// ── GET /actors ───────────────────────────────────────────────────────────────

// ListActors returns the full set of known actors (id + handle + category) so
// the dashboard graph can render standalone nodes even when an actor has no
// persona links yet (e.g. a TERM_FALLBACK actor from a page with no handles).
func (h *Handler) ListActors(c *gin.Context) {
	limit := 200
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.pool.Query(ctx, `
		SELECT actor_id, primary_handle, category, created_at
		FROM actors
		ORDER BY created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query actors: " + err.Error()})
		return
	}
	defer rows.Close()

	type actorRecord struct {
		ActorID       string    `json:"actor_id"`
		PrimaryHandle string    `json:"primary_handle"`
		Category      *string   `json:"category"`
		CreatedAt     time.Time `json:"created_at"`
	}

	actors := []actorRecord{}
	for rows.Next() {
		var a actorRecord
		if err := rows.Scan(&a.ActorID, &a.PrimaryHandle, &a.Category, &a.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan actor: " + err.Error()})
			return
		}
		actors = append(actors, a)
	}
	if actors == nil {
		actors = []actorRecord{}
	}

	c.JSON(http.StatusOK, gin.H{"actors": actors})
}

// ── GET /leads ───────────────────────────────────────────────────────────────

// ListLeads returns keyset-paginated leads from the active persona-links view.
//
// Keyset pagination (sorted on attribution_confidence DESC, link_id ASC) is
// used instead of OFFSET so live-updating correlation links can't silently
// skip or repeat rows mid-scroll. Cursor format: "confidence,link_id".
func (h *Handler) ListLeads(c *gin.Context) {
	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	query := `
		SELECT link_id, actor_id_a, actor_id_b, basis, attribution_confidence,
		       nc_state, created_at, enforced_policy_version
		FROM v_active_persona_links`
	args := []interface{}{}
	argN := 1

	if cursor := c.Query("cursor"); cursor != "" {
		parts := strings.SplitN(cursor, ",", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cursor must be 'confidence,link_id'"})
			return
		}
		conf, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cursor confidence"})
			return
		}
		linkID := parts[1]
		if !validUUID(linkID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cursor link_id"})
			return
		}
		// Keyset: strictly after (confidence, link_id) in DESC/ASC order.
		query += fmt.Sprintf(" WHERE (attribution_confidence, link_id) < ($%d::numeric, $%d::uuid)",
			argN, argN+1)
		args = append(args, conf, linkID)
		argN += 2
	}

	query += fmt.Sprintf(" ORDER BY attribution_confidence DESC, link_id ASC LIMIT $%d", argN)
	args = append(args, limit)

	rows, err := h.pool.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query leads: " + err.Error()})
		return
	}
	defer rows.Close()

	type lead struct {
		LinkID                string    `json:"link_id"`
		ActorIDA              string    `json:"actor_id_a"`
		ActorIDB              string    `json:"actor_id_b"`
		Basis                 string    `json:"basis"`
		AttributionConfidence float64   `json:"attribution_confidence"`
		NCState               string    `json:"nc_state"`
		CreatedAt             time.Time `json:"created_at"`
		PolicyVersion         int       `json:"enforced_policy_version"`
	}
	leads := []lead{}
	for rows.Next() {
		var l lead
		if err := rows.Scan(&l.LinkID, &l.ActorIDA, &l.ActorIDB, &l.Basis,
			&l.AttributionConfidence, &l.NCState, &l.CreatedAt, &l.PolicyVersion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan lead: " + err.Error()})
			return
		}
		leads = append(leads, l)
	}

	// Next cursor from the last row (if any).
	next := ""
	if len(leads) > 0 {
		last := leads[len(leads)-1]
		next = fmt.Sprintf("%.2f,%s", last.AttributionConfidence, last.LinkID)
	}

	c.JSON(http.StatusOK, gin.H{
		"leads": leads,
		"next":  next,
	})
}

// ── GET /images ───────────────────────────────────────────────────────────────

// ImageMeta is a single extracted image reference joined with its owning actor.
// The dashboard MEDIA tab renders these as thumbnails and deep-links each row
// back to the owning actor (click image = click node).
type ImageMeta struct {
	ImageRefID    string    `json:"image_ref_id"`
	ActorID       string    `json:"actor_id,omitempty"`
	Src           string    `json:"src"`
	Type          string    `json:"type,omitempty"`
	Evidence      string    `json:"evidence,omitempty"`
	DetectedAt    time.Time `json:"detected_at"`
	PrimaryHandle string    `json:"primary_handle,omitempty"`
}

// ListImages returns extracted image references (newest first) joined with the
// owning actor's primary handle. `?actor_id=` narrows to a single actor (used
// by the dashboard graph builder to attach evidence images to actor nodes).
func (h *Handler) ListImages(c *gin.Context) {
	limit := 200
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	query := `
		SELECT ih.image_ref_id, ih.actor_id, ih.perceptual_hash, ih.evidence, ih.detected_at,
		       COALESCE(a.primary_handle, '')
		FROM image_hashes ih
		LEFT JOIN actors a ON a.actor_id = ih.actor_id`
	args := []interface{}{}
	argN := 1

	if actorID := c.Query("actor_id"); actorID != "" {
		if !validUUID(actorID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actor_id"})
			return
		}
		query += fmt.Sprintf(" WHERE ih.actor_id = $%d", argN)
		args = append(args, actorID)
		argN++
	}

	query += fmt.Sprintf(" ORDER BY ih.detected_at DESC, ih.image_ref_id ASC LIMIT $%d", argN)
	args = append(args, limit)

	rows, err := h.pool.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query images: " + err.Error()})
		return
	}
	defer rows.Close()

	images := []ImageMeta{}
	for rows.Next() {
		var im ImageMeta
		if err := rows.Scan(&im.ImageRefID, &im.ActorID, &im.Src, &im.Evidence, &im.DetectedAt, &im.PrimaryHandle); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan image: " + err.Error()})
			return
		}
		images = append(images, im)
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

// ── GET /stats ────────────────────────────────────────────────────────────────

// GetStats returns aggregate platform counts for the dashboard stats bar.
// All counts are live database queries — zero mock data.
func (h *Handler) GetStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	type count struct{ N int64 }

	var actors, links, wallets, pgpKeys, flagged, jobsDone count

	// Unique threat actors.
	h.pool.QueryRow(ctx, `SELECT count(*) FROM actors`).Scan(&actors.N)

	// Active persona links (via policy-filtered view).
	h.pool.QueryRow(ctx, `SELECT count(*) FROM v_active_persona_links`).Scan(&links.N)

	// Crypto wallets (identifiers whose kind is wallet-like).
	h.pool.QueryRow(ctx, `SELECT count(DISTINCT value) FROM identifiers WHERE kind ILIKE '%wallet%' OR kind ILIKE '%btc%' OR kind ILIKE '%xmr%'`).Scan(&wallets.N)

	// PGP key fingerprints.
	h.pool.QueryRow(ctx, `SELECT count(DISTINCT value) FROM identifiers WHERE kind ILIKE '%pgp%'`).Scan(&pgpKeys.N)

	// Content flags (CSAM / illegal content detections).
	h.pool.QueryRow(ctx, `SELECT count(*) FROM content_flags`).Scan(&flagged.N)

	// Completed crawl scans (jobs done).
	h.pool.QueryRow(ctx, `SELECT count(*) FROM scans WHERE status = 'completed'`).Scan(&jobsDone.N)

	c.JSON(http.StatusOK, gin.H{
		"actors":    actors.N,
		"links":     links.N,
		"wallets":   wallets.N,
		"pgp_keys":  pgpKeys.N,
		"flagged":   flagged.N,
		"jobs_done": jobsDone.N,
	})
}

// ── GET /scans ───────────────────────────────────────────────────────────────

// ListScans returns crawl scan records for the jobs dashboard.
func (h *Handler) ListScans(c *gin.Context) {
	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.pool.Query(ctx, `
		SELECT scan_id, requested_by, targets, status, started_at, completed_at
		FROM scans
		ORDER BY started_at DESC
		LIMIT $1`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query scans: " + err.Error()})
		return
	}
	defer rows.Close()

	type scanRecord struct {
		ScanID      string     `json:"scan_id"`
		RequestedBy string     `json:"requested_by"`
		Targets     []string   `json:"targets"`
		Status      string     `json:"status"`
		StartedAt   time.Time  `json:"started_at"`
		CompletedAt *time.Time `json:"completed_at"`
	}

	var scans []scanRecord
	for rows.Next() {
		var s scanRecord
		if err := rows.Scan(&s.ScanID, &s.RequestedBy, &s.Targets, &s.Status, &s.StartedAt, &s.CompletedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan row: " + err.Error()})
			return
		}
		scans = append(scans, s)
	}
	if scans == nil {
		scans = []scanRecord{}
	}
	c.JSON(http.StatusOK, gin.H{"scans": scans})
}

// ── Internal helpers ─────────────────────────────────────────────────────────

// upsertActor inserts an actor by primary_handle, returning its ID.
// On conflict it updates last_scan_date and returns the existing ID.
func upsertActor(ctx context.Context, q interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}, handle, category string) (string, error) {
	var id string
	err := q.QueryRow(ctx, `
		INSERT INTO actors (actor_id, primary_handle, category, last_scan_date, created_at)
		VALUES (gen_random_uuid(), $1, $2, now(), now())
		ON CONFLICT (primary_handle) DO UPDATE SET
			last_scan_date = now(),
			category = COALESCE(EXCLUDED.category, actors.category)
		RETURNING actor_id`,
		handle, nullableString(category)).Scan(&id)
	return id, err
}

// upsertIdentifier inserts an identifier, updating last_seen on conflict.
func upsertIdentifier(ctx context.Context, q interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}, actorID, kind, value, provenance string, seenAt time.Time) error {
	_, err := q.Exec(ctx, `
		INSERT INTO identifiers (identifier_id, actor_id, kind, value, source, provenance_tag, first_seen, last_seen)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $6)
		ON CONFLICT (actor_id, kind, value) DO UPDATE SET
			last_seen = EXCLUDED.last_seen,
			provenance_tag = COALESCE(EXCLUDED.provenance_tag, identifiers.provenance_tag)`,
		nullableUUID(actorID), kind, value, "crawler", provenance, seenAt)
	return err
}

// nullableUUID returns nil for empty strings (maps to SQL NULL).
func nullableUUID(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// nullableString returns nil for empty strings (maps to SQL NULL).
func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// termFallbackHandle converts an investigation term into a stable, unique
// fallback actor handle. Terms may contain spaces, slashes, or URL text, so we
// sanitize to a clean handle-like token and prefix with "TERM:" to make
// fallback actors visually distinct from real extracted handles. Empty terms
// and terms that collapse to nothing are rejected.
func termFallbackHandle(term string) string {
	term = strings.TrimSpace(term)
	if term == "" {
		return ""
	}

	var b strings.Builder
	lastUnderscore := false
	for _, r := range term {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastUnderscore = false
			continue
		}
		// Collapse runs of non-alphanumeric characters into one underscore.
		if !lastUnderscore {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}
	out := strings.Trim(b.String(), "_")
	if out == "" {
		return ""
	}
	// Keep handles compact (max 40 chars body, matching web handle norms).
	if len(out) > 40 {
		out = out[:40]
	}
	return "TERM:" + out
}
