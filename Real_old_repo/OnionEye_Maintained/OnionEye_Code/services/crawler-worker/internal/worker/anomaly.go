package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// workerOutcome is the result of a single worker turn. Workers are short-lived
// goroutines (one processOne per token), so anomaly counters live on the Pool
// and span turns — a broken Redis/circuit shows up as rapid fast turnover of
// the whole pool, which is exactly what we kill+backoff against.
type workerOutcome int

const (
	outcomeIdle   workerOutcome = iota // nothing claimed (block window elapsed)
	outcomeOK                          // job fully processed & acked
	outcomeFailed                      // job claimed but the turn died (network-dead)
	outcomePanic                       // goroutine panicked (job failed + immediate kill)
)

// anomalyKind is the machine-readable kind of a detected worker anomaly. It is
// sent to the dashboard as the `anomaly` field of a worker-anomaly event.
type anomalyKind string

const (
	anomalyPanic     anomalyKind = "panic"
	anomalyFailStorm anomalyKind = "rapid-failure-storm"
	anomalyIdleSpin  anomalyKind = "fast-idle-turnover"
)

// Defaults apply when config values are unset (0) — also what tests + bare
// Pool{} instances rely on. Two minutes of quiet resets anomaly streaks so a
// long idle pool never fires a stale kill on its first downhill turn.
const (
	defaultAnomalyFastIdleTurns = 3
	defaultAnomalyFastTurnDur   = 1 * time.Second
	defaultAnomalyMaxFailTurns  = 3
	defaultAnomalyBackoff       = 15 * time.Second
	defaultAnomalyMaxPerMinute  = 5
	anomalyResetAfter           = 2 * time.Minute
)

// anomalyThresholds are the tunables used by a single observe() evaluation.
type anomalyThresholds struct {
	fastIdleTurns int
	fastTurnDur   time.Duration
	maxFailTurns  int
	backoff       time.Duration
	maxPerMinute  int
}

func ifPosInt(v, def int) int {
	if v > 0 {
		return v
	}
	return def
}

func ifPosDur(v, def time.Duration) time.Duration {
	if v > 0 {
		return v
	}
	return def
}

// anomalyThresholds returns the config-backed thresholds with sane defaults so
// anomaly kill/remake works even on bare Pool{} instances (tests, misconfig).
func (p *Pool) anomalyThresholds() anomalyThresholds {
	if p.cfg == nil {
		return anomalyThresholds{
			fastIdleTurns: defaultAnomalyFastIdleTurns,
			fastTurnDur:   defaultAnomalyFastTurnDur,
			maxFailTurns:  defaultAnomalyMaxFailTurns,
			backoff:       defaultAnomalyBackoff,
			maxPerMinute:  defaultAnomalyMaxPerMinute,
		}
	}
	return anomalyThresholds{
		fastIdleTurns: ifPosInt(p.cfg.AnomalyFastIdleTurns, defaultAnomalyFastIdleTurns),
		fastTurnDur:   ifPosDur(p.cfg.AnomalyFastTurnDur, defaultAnomalyFastTurnDur),
		maxFailTurns:  ifPosInt(p.cfg.AnomalyMaxFailTurns, defaultAnomalyMaxFailTurns),
		backoff:       ifPosDur(p.cfg.AnomalyBackoff, defaultAnomalyBackoff),
		maxPerMinute:  ifPosInt(p.cfg.AnomalyMaxPerMinute, defaultAnomalyMaxPerMinute),
	}
}

// anomalyDetector accumulates turn outcomes across worker goroutines and
// decides when kill+remake is warranted. All counters are pool-global because a
// single anomalous condition (dead Redis, dead circuit) hurts every turn.
type anomalyDetector struct {
	mu sync.Mutex

	consecutiveFastIdle int       // consecutive no-job turns finishing under fastTurnDur
	consecutiveFails    int       // consecutive failed turns
	lastTurnAt          time.Time // freshness anchor for streak reset

	backoffUntil time.Time   // a killed slot must not be respawned before this
	kills        []time.Time // rolling 1-min kill log for the rate limiter
}

// observe records the outcome of one worker turn and returns the anomaly kind
// (and true) when a kill-worthy threshold is crossed. A real idle turn takes
// >= 5s (the XReadGroup block window), so ONLY fast no-job returns count toward
// the idle-spin detector — that is the Redis/circuit torn signal. outcomePanic
// is counted toward the failure streak but never returned here: the recover()
// path notifies + kills panics directly so each panic is a single event.
func (d *anomalyDetector) observe(outcome workerOutcome, turnDur time.Duration, t anomalyThresholds, now time.Time) (anomalyKind, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.lastTurnAt.IsZero() && now.Sub(d.lastTurnAt) > anomalyResetAfter {
		d.consecutiveFastIdle = 0
		d.consecutiveFails = 0
	}
	d.lastTurnAt = now

	switch outcome {
	case outcomeOK:
		d.consecutiveFastIdle = 0
		d.consecutiveFails = 0
	case outcomeIdle:
		d.consecutiveFails = 0
		if turnDur < t.fastTurnDur {
			d.consecutiveFastIdle++
			if d.consecutiveFastIdle >= t.fastIdleTurns {
				d.consecutiveFastIdle = 0
				return anomalyIdleSpin, true
			}
		} else {
			d.consecutiveFastIdle = 0
		}
	case outcomeFailed, outcomePanic:
		d.consecutiveFastIdle = 0
		d.consecutiveFails++
		if d.consecutiveFails >= t.maxFailTurns {
			d.consecutiveFails = 0
			return anomalyFailStorm, true
		}
	}
	return "", false
}

// recordKill logs a kill/remake and returns true when the per-minute rate limit
// has been exceeded (in which case the caller should NOT notify the dashboard).
func (d *anomalyDetector) recordKill(maxPerMinute int, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := now.Add(-time.Minute)
	kept := d.kills[:0]
	for _, k := range d.kills {
		if k.After(cutoff) {
			kept = append(kept, k)
		}
	}
	d.kills = kept

	if len(d.kills) >= maxPerMinute {
		return true // rate-limited: suppress this kill's dashboard event
	}
	d.kills = append(d.kills, now)
	return false
}

// beginBackoff holds a killed slot out of the pool until now+backoff, so the
// remade worker starts after the fresh circuit has settled.
func (d *anomalyDetector) beginBackoff(backoff time.Duration, now time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.backoffUntil.Before(now.Add(backoff)) {
		d.backoffUntil = now.Add(backoff)
	}
}

// backoffRemaining returns how long before a killed slot may be respawned.
func (d *anomalyDetector) backoffRemaining(now time.Time) time.Duration {
	d.mu.Lock()
	defer d.mu.Unlock()
	if now.After(d.backoffUntil) {
		return 0
	}
	return d.backoffUntil.Sub(now)
}

// postNotification is the shared best-effort dashboard POST. It never blocks
// or fails the caller: marshal/request/transport errors are logged + swallowed.
// The `at` stamp is always (re)written here so payloads carry a canonical time.
func (p *Pool) postNotification(ctx context.Context, payload map[string]any) {
	if p.cfg == nil {
		return
	}
	url := strings.TrimRight(p.cfg.DashboardAPIURL, "/") + p.cfg.DashboardNotifyPath
	payload["at"] = time.Now().UTC().Format(time.RFC3339Nano)

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("dashboard notify marshal: %v", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Printf("dashboard notify build request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	hc := &http.Client{Timeout: p.cfg.DashboardNotifyTimeout}
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("dashboard notify POST failed: %v", err)
		return
	}
	_ = resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("dashboard notify non-2xx: %s", resp.Status)
	}
}

// notifyAnomaly POSTs a worker-anomaly event so the dashboard can render the
// kill+remake popup.
func (p *Pool) notifyAnomaly(ctx context.Context, workerID string, kind anomalyKind, reason string) {
	p.postNotification(ctx, map[string]any{
		"event":      "worker-anomaly",
		"anomaly_id": fmt.Sprintf("anom-%d", time.Now().UnixNano()),
		"worker_id":  workerID,
		"anomaly":    string(kind),
		"reason":     reason,
	})
}

// handleWorkerAnomaly performs the kill+remake response to a detected anomaly:
//  1. begin the spawn backoff — the pool holds the killed slot out of the pool
//     until the backoff elapses, then respawns a brand-new worker goroutine.
//     The brake ALWAYS engages: the per-minute kill rate limiter only throttles
//     the dashboard event, never the backoff. Otherwise a sustained outage
//     would spin the pool unthrottled (each suppressed kill instantly
//     respawning yet another doomed worker).
//  2. rotate the Tor circuit (best-effort) so the remade worker starts fresh;
//  3. notify the dashboard (rate-limited to maxPerMinute).
//
// It is safe on Pool{} instances (nil cfg only logs) and never blocks.
func (p *Pool) handleWorkerAnomaly(ctx context.Context, workerID string, kind anomalyKind, reason string) {
	if p.cfg == nil {
		log.Printf("ANOMALY worker %s (%s): %s", workerID, kind, reason)
		return
	}
	t := p.anomalyThresholds()

	p.detector.beginBackoff(t.backoff, time.Now())

	if p.detector.recordKill(t.maxPerMinute, time.Now()) {
		log.Printf("ANOMALY [%s] suppressed: kill rate limit %d/min exceeded (worker %s %s)", workerID, t.maxPerMinute, kind, reason)
		return
	}

	log.Printf("ANOMALY KILL worker %s (%s): %s -> fresh circuit + remake in %s", workerID, kind, reason, t.backoff)

	// Fresh identity for the remade worker (best-effort, like rotation).
	p.rotateCircuit("anomaly-" + string(kind) + "-" + workerID)
	p.notifyAnomaly(ctx, workerID, kind, reason)
}

// observeWorkerTurn records a completed worker turn and, when the streak
// thresholds are crossed, initiates the kill+remake. Panics are handled by the
// recover() defer (they notify immediately), so they only count here toward the
// failure streak without generating a second event.
func (p *Pool) observeWorkerTurn(ctx context.Context, workerID string, outcome workerOutcome, turnDur time.Duration) {
	t := p.anomalyThresholds()

	kind, trigger := p.detector.observe(outcome, turnDur, t, time.Now())
	if !trigger {
		return
	}
	p.handleWorkerAnomaly(ctx, workerID, kind, fmt.Sprintf("pool-level %s: worker %s reported %d consecutive bad turns", kind, workerID, t.maxFailTurns))
}
