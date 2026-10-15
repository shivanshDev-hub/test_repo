// Package ratelimit implements Phase 2: Redis sliding-window rate limiting.
package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter enforces per-analyst quotas via a Redis sorted set.
type Limiter struct {
	rdb    *redis.Client
	max    int
	window time.Duration
}

// New builds a sliding-window limiter.
func New(rdb *redis.Client, max int, window time.Duration) *Limiter {
	return &Limiter{rdb: rdb, max: max, window: window}
}

// Allow checks and records a request.
// Returns (allowed, countInWindow, remainingQuota, retryAfterSecs).
func (l *Limiter) Allow(ctx context.Context, analystID string) (bool, int, int, int) {
	key := fmt.Sprintf("analyst_rate:%s", analystID)
	now := float64(time.Now().UnixNano()) / 1e9
	cutoff := now - l.window.Seconds()

	pipe := l.rdb.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%f", cutoff))
	pipe.ZCard(ctx, key)
	pipe.ZAdd(ctx, key, redis.Z{Score: now, Member: fmt.Sprintf("%f", now)})
	pipe.Expire(ctx, key, l.window+time.Second)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		// Redis unavailable: fail-open so local/dev stacks without Redis still ingest.
		return true, 0, l.max, 0
	}

	count := int(cmds[1].(*redis.IntCmd).Val())
	if count >= l.max {
		l.rdb.ZRem(ctx, key, fmt.Sprintf("%f", now))
		retry := int(l.window.Seconds())
		if retry < 1 {
			retry = 1
		}
		return false, count, 0, retry
	}

	return true, count + 1, l.max - count - 1, 0
}
