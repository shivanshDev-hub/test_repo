package worker

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/onioneye/crawler-worker/internal/clients"
	"github.com/onioneye/crawler-worker/internal/crawler"
	"github.com/onioneye/crawler-worker/internal/extract"
	"github.com/onioneye/crawler-worker/internal/frontier"
	"github.com/onioneye/crawler-worker/internal/job"
)

// activeJob is the job currently being processed by a worker goroutine,
// tracked so a panic can mark it failed. Protected by a mutex.
var (
	activeMu    sync.Mutex
	activeJobID string
)

func setActiveJob(id string) {
	activeMu.Lock()
	activeJobID = id
	activeMu.Unlock()
}

func getActiveJob() string {
	activeMu.Lock()
	defer activeMu.Unlock()
	return activeJobID
}

// updateStatus updates the status field of the task_status:<jobID> JSON
// document written by the ingestion-api taskbuilder (rdb.Set + JSON blob).
// A plain HSET would fail with WRONGTYPE because the key is a string, so we
// read-modify-write the JSON in place, preserving the other fields.
func (p *Pool) updateStatus(ctx context.Context, jobID, status string) {
	key := "task_status:" + jobID

	raw, err := p.rdb.Get(ctx, key).Result()
	if err != nil {
		// No document exists yet (e.g. offline enqueue) — create a minimal one
		// using the same JSON shape as the taskbuilder.
		doc := map[string]any{"status": status}
		b, _ := json.Marshal(doc)
		p.rdb.Set(ctx, key, b, 24*time.Hour)
		return
	}

	doc := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		doc = map[string]any{"status": status}
	}
	doc["status"] = status
	doc["last_updated"] = time.Now().UTC().Format(time.RFC3339)

	if b, err := json.Marshal(doc); err == nil {
		p.rdb.Set(ctx, key, b, 24*time.Hour)
	}
}

// failActiveJob marks the in-flight job failed in Redis on panic.
func (p *Pool) failActiveJob() {
	if id := getActiveJob(); id != "" {
		p.updateStatus(context.Background(), id, "failed")
	}
}

// enqueueManual pushes a blocked job to the manual-intervention stream with a
// 5-minute expiry slot, per the architecture's CAPTCHA handling.
func (p *Pool) enqueueManual(ctx context.Context, jb *job.Job, reason string) {
	p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: p.cfg.ManualStream,
		Values: map[string]interface{}{
			"job_id":     jb.JobID,
			"target":     jb.Target,
			"reason":     reason,
			"flagged_at": time.Now().UTC().Format(time.RFC3339),
			"expires_in": "300", // 5-minute TTL slot
		},
	})
}

func newFrontier(maxDepth int) *frontier.Frontier { return frontier.New(maxDepth) }

func hostOf(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return u.Hostname()
}

// alternateSchemeURL returns url with its http/https scheme flipped, or "" when
// the URL uses a non-http(s) scheme (there is no meaningful alternate to try).
// Used to recover seeds whose host only serves the other protocol.
func alternateSchemeURL(rawURL string) string {
	switch {
	case strings.HasPrefix(rawURL, "http://"):
		return "https://" + strings.TrimPrefix(rawURL, "http://")
	case strings.HasPrefix(rawURL, "https://"):
		return "http://" + strings.TrimPrefix(rawURL, "https://")
	default:
		return ""
	}
}

func (p *Pool) shieldClient() *clients.ShieldClient  { return p.shield }
func (p *Pool) storageClient() *clients.StorageClient { return p.storage }

func extractText(raw []byte) *clients.Extracted { return extract.Extract(raw) }

func scanPayloadFrom(jb *job.Job, f *crawler.Fetched, bodyRef string, ex *clients.Extracted) *clients.ScanPayload {
	return &clients.ScanPayload{
		JobID:        jb.JobID,
		SourceURL:    f.URL,
		Provenance:   provenanceTag(f.URL),
		Extracted:    *ex,
		BodyRef:      bodyRef,
		TorCircuitID: f.TorCircuitID,
		FetchedAt:    f.FetchedAt.Format(time.RFC3339),
		Terms:        jb.Terms,
	}
}

// provenanceTag is a rule-based classifier over the URL path (matches the
// shield's own assignment for the demo).
func provenanceTag(rawURL string) string {
	low := strings.ToLower(rawURL)
	switch {
	case strings.Contains(low, "/listing/"):
		return "marketplace_listing"
	case strings.Contains(low, "/thread/"):
		return "forum_post"
	case strings.Contains(low, "/chat/"):
		return "chat_message"
	default:
		return "generic_page"
	}
}
