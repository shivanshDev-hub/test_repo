// Package session implements Phase 5: analyst session history via Redis.
package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/onioneye/ingestion-api/internal/model"
)

// Store tracks recent analyst queries for session-aware expansion.
type Store struct {
	rdb    *redis.Client
	window time.Duration
}

// New creates a session store.
func New(rdb *redis.Client, window time.Duration) *Store {
	return &Store{rdb: rdb, window: window}
}

// BuildContext merges the current query with recent session history (Phase 5).
func (s *Store) BuildContext(ctx context.Context, req *model.ValidatedRequest) (*model.SessionContext, error) {
	key := fmt.Sprintf("analyst_session:%s", req.AnalystID)
	cutoff := time.Now().Add(-s.window).Unix()

	// Fetch recent queries from the session list.
	raw, err := s.rdb.LRange(ctx, key, 0, 9).Result()
	if err != nil {
		// Redis unavailable: continue with current query only (no session history).
		raw = nil
	}

	recent := make([]string, 0, len(raw))
	for _, item := range raw {
		var q struct {
			Terms []string `json:"terms"`
		}
		if json.Unmarshal([]byte(item), &q) == nil {
			recent = append(recent, q.Terms...)
		}
	}

	sc := &model.SessionContext{
		AnalystID:          req.AnalystID,
		TargetType:         req.TargetType,
		TargetHost:         req.TargetHost,
		MaxDepth:           req.MaxDepth,
		CurrentQuery:       req.SeedTerms,
		RecentSessionQuery: recent,
		SessionWindowSecs:  int(s.window.Seconds()),
	}

	// Record this query into the session history.
	payload, _ := json.Marshal(map[string]any{
		"terms": req.SeedTerms,
		"ts":    time.Now().Unix(),
	})
	s.rdb.LPush(ctx, key, payload)
	s.rdb.LTrim(ctx, key, 0, 49)
	s.rdb.ZRemRangeByScore(ctx, key+":ts", "0", fmt.Sprintf("%d", cutoff))
	s.rdb.Expire(ctx, key, s.window+time.Minute)

	return sc, nil
}
