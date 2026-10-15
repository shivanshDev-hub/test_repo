// Package handler wires the 10-phase ingestion pipeline into an HTTP handler.
package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/onioneye/ingestion-api/internal/config"
	"github.com/onioneye/ingestion-api/internal/expander"
	"github.com/onioneye/ingestion-api/internal/guard"
	"github.com/onioneye/ingestion-api/internal/model"
	"github.com/onioneye/ingestion-api/internal/ratelimit"
	"github.com/onioneye/ingestion-api/internal/session"
	"github.com/onioneye/ingestion-api/internal/stream"
	"github.com/onioneye/ingestion-api/internal/taskbuilder"
	"github.com/onioneye/ingestion-api/internal/validator"
)

// Handler holds all pipeline dependencies.
type Handler struct {
	cfg        *config.Config
	rdb        *redis.Client
	limiter    *ratelimit.Limiter
	allowlist  *validator.Allowlist
	sessions   *session.Store
	inputGuard *guard.Guard
	expander   *expander.Expander
	outGuard   *expander.OutputGuard
	builder    *taskbuilder.Builder
	emitter    *stream.Emitter
}

// New constructs the handler with all phases wired.
func New(cfg *config.Config, rdb *redis.Client) (*Handler, error) {
	allowlistData, err := os.ReadFile(cfg.AllowlistPath)
	if err != nil {
		return nil, err
	}
	allowlist, err := validator.LoadAllowlist(allowlistData)
	if err != nil {
		return nil, err
	}

	return &Handler{
		cfg:        cfg,
		rdb:        rdb,
		limiter:    ratelimit.New(rdb, cfg.RateLimitMax, time.Duration(cfg.RateLimitWindow)*time.Second),
		allowlist:  allowlist,
		sessions:   session.New(rdb, time.Duration(cfg.SessionWindow)*time.Second),
		inputGuard: guard.New(cfg.OllamaURL, cfg.GuardModel),
		expander:   expander.New(cfg.OllamaURL, cfg.ExpanderModel),
		outGuard:   expander.NewOutputGuard(),
		builder:    taskbuilder.New(rdb),
		emitter:    stream.New(rdb),
	}, nil
}

// RegisterRoutes mounts the ingestion endpoint.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/investigations", h.Ingest)
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Ingest runs the full 10-phase pipeline.
func (h *Handler) Ingest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
	defer cancel()

	// Phase 1: Ingest raw payload.
	var req model.InvestigationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "detail": err.Error()})
		return
	}
	env := &model.Envelope{
		ClientIP:   c.ClientIP(),
		AnalystID:  req.AnalystID,
		RawPayload: req,
		ReceivedAt: time.Now().UTC(),
	}

	// Phase 2: Rate limit.
	allowed, count, remaining, retry := h.limiter.Allow(ctx, env.AnalystID)
	if !allowed {
		c.Header("Retry-After", fmt.Sprintf("%d", retry))
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":              "rate limit exceeded",
			"requests_in_window": count,
			"retry_after":        retry,
		})
		return
	}

	// Phase 3: Sanitize + validate.
	validated, _ := validator.ValidateRequest(&req, h.cfg.MaxSeedTerms, h.cfg.MaxDepthMax)
	if validated.AnalystID == "" || validated.TargetType == "" || validated.TargetValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	// Phase 4: Authorization (allowlist + target value charset).
	auth := h.allowlist.Check(validated.TargetType)
	if !auth.IsAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized target type", "reason": auth.Reason})
		return
	}
	if !validator.ValidateTargetValue(validated.TargetValue) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target value format"})
		return
	}

	// Phase 5: Session context.
	sessCtx, err := h.sessions.BuildContext(ctx, validated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session error"})
		return
	}

	// Phase 6: Input guard.
	guardDecision := h.inputGuard.Check(ctx, sessCtx)
	if guardDecision.Decision == "BLOCK" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":           "input blocked by safety guard",
			"reason_code":     guardDecision.ReasonCode,
			"violating_terms": guardDecision.ViolatingTerms,
		})
		return
	}

	// Phase 7: Term expansion (per seed term).
	var allVerified []string
	var allStripped []string
	verdict := "PASS"
	llmUnavailable := false
	if guardDecision.ReasonCode == "TIER2_LLM_UNAVAILABLE" {
		llmUnavailable = true
		verdict = "PASS_NO_LLM"
	}
	for _, term := range validated.SeedTerms {
		exp, err := h.expander.Expand(ctx, term, sessCtx)
		if err != nil {
			// Fallback: use the original term unexpanded when Ollama is offline.
			allVerified = append(allVerified, term)
			llmUnavailable = true
			if verdict == "PASS" {
				verdict = "PASS_NO_LLM"
			}
			continue
		}
		// Phase 8: Output guard.
		verified := h.outGuard.Verify(exp)
		allVerified = append(allVerified, verified.VerifiedTerms...)
		allStripped = append(allStripped, verified.StrippedTerms...)
		if verified.GuardVerdict == "PARTIAL_DROP" {
			verdict = "PARTIAL_DROP"
		}
	}

	// Phase 9: Task builder (idempotency).
	verifiedExp := &model.VerifiedExpansion{
		ExpansionID:   "pipeline",
		VerifiedTerms: allVerified,
		StrippedTerms: allStripped,
		GuardVerdict:  verdict,
		TargetHost:    validated.TargetHost,
		MaxDepth:      validated.MaxDepth,
		AnalystID:     validated.AnalystID,
	}
	msg, isDup, jobID, err := h.builder.Build(ctx, validated, verifiedExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "task build error"})
		return
	}

	// Phase 10: Stream emission.
	if !isDup {
		if _, err := h.emitter.Publish(ctx, msg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "stream publish error"})
			return
		}
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":          jobID,
		"status":          "queued",
		"duplicate":       isDup,
		"rate_limit":      gin.H{"allowed": allowed, "count": count, "remaining": remaining},
		"guard_verdict":   verdict,
		"stripped_terms":  allStripped,
		"expanded_terms":  allVerified,
		"sanitized":       validated.Sanitized,
		"received_at":     env.ReceivedAt,
		"llm_unavailable": llmUnavailable,
	})
}
