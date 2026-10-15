package worker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onioneye/crawler-worker/internal/config"
)

func TestPoolNotifyAnomalyPayload(t *testing.T) {
	srv, req, body := captureNotification(t, http.StatusOK)
	defer srv.Close()

	p := &Pool{
		cfg: &config.Config{
			DashboardAPIURL:        srv.URL,
			DashboardNotifyPath:    "/api/v1/dashboard/notify",
			DashboardNotifyTimeout: 2 * time.Second,
		},
	}

	p.notifyAnomaly(context.Background(), "w-7", anomalyPanic, "recovered panic: runtime error")

	if req() == nil {
		t.Fatal("expected HTTP POST request to the mock dashboard, got none")
	}
	if req().Method != http.MethodPost {
		t.Fatalf("expected POST, got %s", req().Method)
	}

	b := body()
	if b["event"] != "worker-anomaly" {
		t.Fatalf("event = %v, want worker-anomaly", b["event"])
	}
	if b["worker_id"] != "w-7" {
		t.Fatalf("worker_id = %v, want w-7", b["worker_id"])
	}
	if b["anomaly"] != "panic" {
		t.Fatalf("anomaly = %v, want panic", b["anomaly"])
	}
	if b["reason"] != "recovered panic: runtime error" {
		t.Fatalf("reason = %v, want the panic reason", b["reason"])
	}
	if _, ok := b["anomaly_id"].(string); !ok {
		t.Fatalf("anomaly_id missing or not a string: %T", b["anomaly_id"])
	}
	if _, ok := b["at"].(string); !ok {
		t.Fatalf("at missing or not a string: %T", b["at"])
	}
}

func TestAnomalyDetectorFastIdleSpin(t *testing.T) {
	d := &anomalyDetector{}
	thr := anomalyThresholds{fastIdleTurns: 3, fastTurnDur: time.Second, maxFailTurns: 3}
	now := time.Now()

	for i := 0; i < 2; i++ {
		kind, ok := d.observe(outcomeIdle, 100*time.Millisecond, thr, now.Add(time.Duration(i)*time.Millisecond))
		if ok {
			t.Fatalf("turn %d: unexpected anomaly %s", i+1, kind)
		}
	}
	kind, ok := d.observe(outcomeIdle, 50*time.Millisecond, thr, now.Add(3*time.Millisecond))
	if !ok {
		t.Fatal("expected fast-idle-turnover anomaly on the 3rd fast no-job turn")
	}
	if kind != anomalyIdleSpin {
		t.Fatalf("kind = %s, want %s", kind, anomalyIdleSpin)
	}
}

func TestAnomalyDetectorSlowIdleDoesNotTrigger(t *testing.T) {
	d := &anomalyDetector{}
	thr := anomalyThresholds{fastIdleTurns: 3, fastTurnDur: time.Second, maxFailTurns: 3}
	now := time.Now()

	// A genuine idle read blocks >= 5s — never an anomaly, no matter how many.
	for i := 0; i < 25; i++ {
		_, ok := d.observe(outcomeIdle, 5*time.Second, thr, now.Add(time.Duration(i)*time.Second))
		if ok {
			t.Fatalf("turn %d: genuine idle incorrectly flagged as anomaly", i+1)
		}
	}
}

func TestAnomalyDetectorFailStorm(t *testing.T) {
	d := &anomalyDetector{}
	thr := anomalyThresholds{fastIdleTurns: 3, fastTurnDur: time.Second, maxFailTurns: 3}
	now := time.Now()

	for i := 0; i < 2; i++ {
		kind, ok := d.observe(outcomeFailed, 10*time.Second, thr, now.Add(time.Duration(i)*time.Second))
		if ok {
			t.Fatalf("turn %d: unexpected anomaly %s", i+1, kind)
		}
	}
	kind, ok := d.observe(outcomeFailed, 10*time.Second, thr, now.Add(3*time.Second))
	if !ok {
		t.Fatal("expected rapid-failure-storm anomaly on the 3rd consecutive failed turn")
	}
	if kind != anomalyFailStorm {
		t.Fatalf("kind = %s, want %s", kind, anomalyFailStorm)
	}
}
func TestAnomalyDetectorOKResetsStreak(t *testing.T) {
	d := &anomalyDetector{}
	thr := anomalyThresholds{fastIdleTurns: 3, fastTurnDur: time.Second, maxFailTurns: 3}
	now := time.Now()

	// fail, fail, OK (resets), then a further fail must NOT trigger yet.
	d.observe(outcomeFailed, 10*time.Second, thr, now)
	d.observe(outcomeFailed, 10*time.Second, thr, now.Add(time.Second))
	d.observe(outcomeOK, 10*time.Second, thr, now.Add(2*time.Second))
	if kind, ok := d.observe(outcomeFailed, 10*time.Second, thr, now.Add(3*time.Second)); ok {
		t.Fatalf("failure streak should have been reset by the OK turn, got anomaly %s", kind)
	}
}

func TestAnomalyDetectorStaleStreakResets(t *testing.T) {
	d := &anomalyDetector{}
	thr := anomalyThresholds{fastIdleTurns: 3, fastTurnDur: time.Second, maxFailTurns: 3}
	base := time.Now()

	// Two fast idle turns, then a long quiet gap: the streak is stale and must
	// not contribute to a later episode.
	d.observe(outcomeIdle, 50*time.Millisecond, thr, base)
	d.observe(outcomeIdle, 50*time.Millisecond, thr, base.Add(time.Second))

	gap := base.Add(10 * time.Minute)
	if kind, ok := d.observe(outcomeIdle, 50*time.Millisecond, thr, gap.Add(1*time.Second)); ok {
		t.Fatalf("stale streak leaked into a new episode: %s", kind)
	}
	if kind, ok := d.observe(outcomeIdle, 50*time.Millisecond, thr, gap.Add(2*time.Second)); ok {
		t.Fatalf("expected no anomaly after only 2 fresh fast idle turns, got %s", kind)
	}
}

func TestAnomalyDetectorRateLimit(t *testing.T) {
	d := &anomalyDetector{}
	t0 := time.Now()

	if d.recordKill(2, t0) {
		t.Fatal("first kill should be allowed")
	}
	if d.recordKill(2, t0.Add(time.Second)) {
		t.Fatal("second kill should be allowed")
	}
	if !d.recordKill(2, t0.Add(2*time.Second)) {
		t.Fatal("third kill inside the window must be rate-limited")
	}
	// Window rolls over.
	if d.recordKill(2, t0.Add(61*time.Second)) {
		t.Fatal("kill outside the rolling window should be allowed again")
	}
}

func TestPoolHandleAnomalyNilConfigNoop(t *testing.T) {
	p := &Pool{}                                                             // cfg == nil
	p.handleWorkerAnomaly(context.Background(), "w-1", anomalyPanic, "boom") // must not panic
}

func TestPoolHandleAnomalyRateLimitSuppressesPOST(t *testing.T) {
	var posts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	p := &Pool{
		cfg: &config.Config{
			DashboardAPIURL:        srv.URL,
			DashboardNotifyPath:    "/notify",
			DashboardNotifyTimeout: 2 * time.Second,
			AnomalyMaxPerMinute:    2,
			AnomalyBackoff:         time.Second,
		},
	}

	p.handleWorkerAnomaly(context.Background(), "w-1", anomalyPanic, "panic 1")
	p.handleWorkerAnomaly(context.Background(), "w-2", anomalyPanic, "panic 2")
	p.handleWorkerAnomaly(context.Background(), "w-3", anomalyPanic, "panic 3") // limited

	if got := posts.Load(); got != 2 {
		t.Fatalf("dashboard POSTs = %d, want 2 (3rd kill must be rate-limited)", got)
	}
	// The backoff brake must engage even for a suppressed kill — the rate limit
	// only throttles the dashboard event, never the local kill+remake. Without
	// this, a sustained outage respawns doomed workers unthrottled (fast-churn).
	if rem := p.detector.backoffRemaining(time.Now()); rem <= 0 {
		t.Fatalf("backoff brake not engaged after rate-limited kill (backoffRemaining=%v)", rem)
	}
}
