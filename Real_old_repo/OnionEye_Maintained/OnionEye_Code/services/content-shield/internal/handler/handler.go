// Package handler implements POST /shield — the raw-byte safety gate.
package handler

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/onioneye/content-shield/internal/hashset"
	"github.com/onioneye/content-shield/internal/phash"
	"github.com/onioneye/content-shield/internal/storage"
)

// ShieldRequest is the worker's inbound payload (raw bytes travel base64 for
// the demo; a blob store would replace this with a raw_bytes_ref).
type ShieldRequest struct {
	JobID        string `json:"job_id"`
	SourceURL    string `json:"source_url"`
	HTTPStatus   int    `json:"http_status"`
	RawBytesB64  string `json:"raw_bytes"`
	TorCircuitID string `json:"tor_circuit_id"`
}

// ContentFlag mirrors a content_flags row.
type ContentFlag struct {
	Category    string `json:"category"`
	ActionTaken string `json:"action_taken"`
	HashMatch   bool   `json:"hash_match"`
	HashAlgo    string `json:"hash_algorithm,omitempty"`
	FlaggedAt   string `json:"flagged_at"`
}

// ShieldResponse is the verdict handed back to the worker.
type ShieldResponse struct {
	ShieldPassed bool          `json:"shield_passed"`
	BodyRef      string        `json:"body_ref"`      // minted ONLY on clean pass
	ContentFlags []ContentFlag `json:"content_flags"` // present only on discard
}

// Shield handles the POST /shield endpoint.
type Shield struct {
	set      *hashset.HashSet
	ph       phash.Hasher
	store    *storage.Client
	hashAlgo string
}

// New builds a Shield handler.
func New(set *hashset.HashSet, ph phash.Hasher, store *storage.Client) *Shield {
	return &Shield{set: set, ph: ph, store: store, hashAlgo: "sha256"}
}

// Handle is the Gin handler. Pipeline order: hashing + hash check happen before
// any decoding or rendering of the bytes.
func (s *Shield) Handle(c *gin.Context) {
	var req ShieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	raw, err := base64.StdEncoding.DecodeString(req.RawBytesB64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid raw_bytes encoding"})
		return
	}

	// 1) Hash the RAW byte stream (never decoded).
	digest := s.ph.HashOf(raw)

	// 2) Membership check.
	if s.set.Contains(digest) {
		// Zero the buffer IMMEDIATELY before any other statement.
		phash.Zero(raw)

		flag := ContentFlag{
			Category:    "csam",
			ActionTaken: "content_discarded",
			HashMatch:   true,
			HashAlgo:    s.hashAlgo,
			FlaggedAt:   time.Now().UTC().Format(time.RFC3339),
		}

		// The pipeline terminates here: no body_ref is minted, the parser never
		// sees this job's data. Only the flag continues, via storage-api.
		s.writeFlag(c, req, flag)

		c.JSON(http.StatusOK, ShieldResponse{
			ShieldPassed: false,
			BodyRef:      "",
			ContentFlags: []ContentFlag{flag},
		})
		return
	}

	// 3) Clean pass — mint the durable body_ref now. This is the ONLY point the
	// handle is created (the downloader's raw bytes are consumed/released here).
	c.JSON(http.StatusOK, ShieldResponse{
		ShieldPassed: true,
		BodyRef:      uuid.NewString(),
		ContentFlags: nil,
	})
}

// writeFlag best-effort persists the flag to PostgreSQL via storage-api.
func (s *Shield) writeFlag(c *gin.Context, req ShieldRequest, flag ContentFlag) {
	payload := storage.FlagPayload{
		JobID:      req.JobID,
		SourceURL:  req.SourceURL,
		Category:   flag.Category,
		ActionTaken: flag.ActionTaken,
		HashMatch:  flag.HashMatch,
		FlaggedAt:  flag.FlaggedAt,
	}
	if err := s.store.PersistFlag(c.Request.Context(), &payload); err != nil {
		// Non-fatal: the verdict is still returned to the worker; the flag row
		// is retried by storage retention. Log for audit.
		log.Printf("flag persistence failed: %v", err)
	}
}