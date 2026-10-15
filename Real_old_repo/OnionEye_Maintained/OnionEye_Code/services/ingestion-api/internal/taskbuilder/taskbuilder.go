// Package taskbuilder implements Phase 9: idempotency + task status hash.
package taskbuilder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/onioneye/ingestion-api/internal/model"
)

// Builder creates idempotent task messages and records status hashes.
type Builder struct {
	rdb *redis.Client
}

// New creates a task builder.
func New(rdb *redis.Client) *Builder {
	return &Builder{rdb: rdb}
}

// Build computes the idempotency key and checks for duplicates (Phase 9).
// Returns the task message, whether it's a duplicate, and the job ID.
func (b *Builder) Build(ctx context.Context, req *model.ValidatedRequest, verified *model.VerifiedExpansion) (*model.TaskMessage, bool, string, error) {
	// Idempotency key: hash of analyst + target + terms + depth.
	termsJSON, _ := json.Marshal(verified.VerifiedTerms)
	keyInput := fmt.Sprintf("%s|%s|%s|%d", req.AnalystID, req.TargetValue, termsJSON, req.MaxDepth)
	sum := sha256.Sum256([]byte(keyInput))
	jobID := hex.EncodeToString(sum[:16])

	// Check for existing job.
	exists, err := b.rdb.Exists(ctx, "task_status:"+jobID).Result()
	if err != nil {
		exists = 0
	}
	if exists > 0 {
		return nil, true, jobID, nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	msg := &model.TaskMessage{
		JobID:          jobID,
		Target:         req.TargetValue,
		Terms:          verified.VerifiedTerms,
		DepthRemaining: req.MaxDepth,
		EnqueuedAt:     now,
	}

	// Record status hash.
	status := &model.TaskStatus{
		Status:       "queued",
		AnalystID:    req.AnalystID,
		Target:       req.TargetValue,
		TermsCount:   fmt.Sprintf("%d", len(verified.VerifiedTerms)),
		GuardVerdict: verified.GuardVerdict,
		CreatedAt:    now,
		LastUpdated:  now,
	}
	statusJSON, _ := json.Marshal(status)
	b.rdb.Set(ctx, "task_status:"+jobID, statusJSON, 24*time.Hour)

	return msg, false, jobID, nil
}
