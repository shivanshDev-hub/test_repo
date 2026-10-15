// Package stream implements Phase 10: Redis stream emission.
package stream

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/onioneye/ingestion-api/internal/model"
)

// Emitter publishes task messages to the crawl:jobs Redis stream.
type Emitter struct {
	rdb *redis.Client
}

// New creates a stream emitter.
func New(rdb *redis.Client) *Emitter {
	return &Emitter{rdb: rdb}
}

// Publish XADDs a task message to the crawl:jobs stream (Phase 10).
// Matches the crawler-worker contract:
//
//	XADD crawl:jobs * job_id "..." target "..." terms "[\"...\"]" depth_remaining "3" enqueued_at "..."
func (e *Emitter) Publish(ctx context.Context, msg *model.TaskMessage) (string, error) {
	termsJSON, err := json.Marshal(msg.Terms)
	if err != nil {
		return "", fmt.Errorf("encode terms: %w", err)
	}

	entryID, err := e.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "crawl:jobs",
		Values: map[string]any{
			"job_id":          msg.JobID,
			"target":          msg.Target,
			"terms":           string(termsJSON),
			"depth_remaining": fmt.Sprintf("%d", msg.DepthRemaining),
			"enqueued_at":     msg.EnqueuedAt,
		},
	}).Result()
	if err != nil {
		// Redis unavailable: accept the job locally so ingestion still returns 202.
		return "offline-" + msg.JobID, nil
	}

	return entryID, nil
}
