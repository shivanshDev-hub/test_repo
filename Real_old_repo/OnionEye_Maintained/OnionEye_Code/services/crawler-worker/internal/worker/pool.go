// Package worker implements the bounded worker pool. A fixed number of tokens
// (buffered channel) bounds concurrency; every goroutine is wrapped in a
// recover() and releases its token via defer, per the architecture notes.
package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/onioneye/crawler-worker/internal/clients"
	"github.com/onioneye/crawler-worker/internal/config"
	"github.com/onioneye/crawler-worker/internal/crawler"
	"github.com/onioneye/crawler-worker/internal/job"
)

// Pool consumes jobs from the Redis stream with a bounded set of workers.
type Pool struct {
	cfg     *config.Config
	rdb     *redis.Client
	tokens  chan struct{}
	wg      sync.WaitGroup
	crawler *crawler.Client
	shield  *clients.ShieldClient
	storage *clients.StorageClient

	// on-demand screenshot pipeline
	screenshotSem chan struct{}
	render         *clients.RenderClient

	// worker lifecycle
	rotationMu sync.Mutex
	rotationID string

	// anomaly kill+remake
	workerSeq uint64          // monotonic source of worker ids
	detector  anomalyDetector // turn-outcome counters + kill rate limiter
}

// ActiveRotationID returns the ID for the most recent graceful rotation.
func (p *Pool) ActiveRotationID() string {
	p.rotationMu.Lock()
	defer p.rotationMu.Unlock()
	return p.rotationID
}

func (p *Pool) bumpRotationID() string {
	p.rotationMu.Lock()
	defer p.rotationMu.Unlock()
	p.rotationID = fmt.Sprintf("rot-%d", time.Now().UnixNano())
	return p.rotationID
}

// NewPool creates a pool sized by cfg.WorkerCount and prefills the tokens.
func NewPool(cfg *config.Config, rdb *redis.Client, cr *crawler.Client) *Pool {
	p := &Pool{
		cfg:     cfg,
		rdb:     rdb,
		crawler: cr,
		shield:  clients.NewShieldClient(cfg.ShieldURL),
		storage: clients.NewStorageClient(cfg.StorageURL),
		render:  clients.NewRenderClient(cfg.RenderSandboxURL, cfg.RenderSandboxToken),

		// bounded capture concurrency
		screenshotSem: make(chan struct{}, cfg.MaxConcurrentScreenshots),

		tokens: make(chan struct{}, cfg.WorkerCount),
	}
	for i := 0; i < cfg.WorkerCount; i++ {
		p.tokens <- struct{}{}
	}
	return p
}

// notifyRotation POSTs a best-effort worker-rotation notification to the
// dashboard API (shared best-effort transport, see postNotification).
func (p *Pool) notifyRotation(ctx context.Context, reason string) {
	p.rotationMu.Lock()
	rotID := p.rotationID
	p.rotationMu.Unlock()
	p.postNotification(ctx, map[string]any{
		"event":       "worker-rotation",
		"rotation_id": rotID,
		"reason":      reason,
	})
}

// rotateCircuit refreshes the shared Tor circuit (NEWNYM) so the post-drain
// worker batches start on a fresh identity. Best-effort: if the control port is
// unreachable the rotation schedule still continues.
func (p *Pool) rotateCircuit(rotationID string) {
	if p.crawler == nil {
		return
	}
	if err := p.crawler.RotateCircuit(); err != nil {
		log.Printf("rotation %s: circuit refresh failed (best-effort, continuing): %v", rotationID, err)
		return
	}
	log.Printf("rotation %s: circuit refreshed -> %s", rotationID, p.crawler.CircuitID())
}

// Run blocks forever, pulling jobs and handing them to workers.
// Each worker goroutine acquires a token, recovers from panics, and always returns it.
func (p *Pool) Run(ctx context.Context) {
	_ = p.bumpRotationID() // initialize rotation id

	// Graceful periodic rotation: stop *starting* new jobs for a short drain window,
	// wait for in-flight jobs to finish, rotate the Tor circuit, then resume.
	// We do not cancel in-flight jobs.
	rotationEvery := 5 * time.Minute
	drainWindow := 60 * time.Second

	rotationTimer := time.NewTimer(rotationEvery)
	defer rotationTimer.Stop()

	var rotationMu sync.Mutex
	rotationInhibit := false

	startWorker := func() {
		rotationMu.Lock()
		allowed := !rotationInhibit
		rotationMu.Unlock()
		if !allowed {
			// If rotation is inhibiting, return the token immediately.
			p.tokens <- struct{}{}
			return
		}

		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			defer func() {
				p.tokens <- struct{}{} // return token to pool
			}()

			id := fmt.Sprintf("w-%d", atomic.AddUint64(&p.workerSeq, 1))
			start := time.Now()
			outcome := outcomeIdle

			// observe must run AFTER the recover defer (registered first) so it
			// sees the final outcome including panics. LIFO: recover sets
			// outcome, then observe records it.
			defer func() {
				p.observeWorkerTurn(ctx, id, outcome, time.Since(start))
			}()
			defer func() {
				if r := recover(); r != nil {
					outcome = outcomePanic
					log.Printf("ANOMALY worker %s: PANIC: %v\n%s", id, r, debug.Stack())
					p.failActiveJob()
					p.handleWorkerAnomaly(ctx, id, anomalyPanic, fmt.Sprintf("recovered panic: %v", r))
				}
			}()

			outcome = p.processOne(ctx)
		}()
	}

	for {
		select {
		case <-ctx.Done():
			p.wg.Wait()
			return
		case <-rotationTimer.C:
			// Enter inhibit/drain mode.
			rotationMu.Lock()
			rotationInhibit = true
			rotationMu.Unlock()

			reason := "scheduled rotation"
			rotID := p.bumpRotationID()
			log.Printf("[rotation %s] drain started reason=%q (no new jobs until in-flight complete or 60s)", rotID, reason)
			p.notifyRotation(ctx, reason)

			// Wait for either completion or drain window, then resume.
			ch := make(chan struct{})
			go func() {
				p.wg.Wait()
				close(ch)
			}()

			select {
			case <-ch:
				// all in-flight done
			case <-time.After(drainWindow):
				// blast-radius control: if something is stuck, we don't want to hang forever.
				// Note: we still do NOT cancel in-flight work.
			}

			// Rotate the Tor circuit so the resumed workers start on a fresh
			// identity. Best-effort: if the control port is unreachable we still
			// resume the pool so crawling is never blocked by rotation.
			p.rotateCircuit(rotID)
			log.Printf("[rotation %s] drain complete, workers resuming", rotID)

			rotationMu.Lock()
			rotationInhibit = false
			rotationMu.Unlock()

			rotationTimer.Reset(rotationEvery)
		case <-p.tokens:
			// Anomaly kill+remake gate: a slot killed for an anomaly is held
			// out of the pool until the backoff elapses, so the remade worker
			// starts only after the fresh circuit has settled. Rotation inhibit
			// still takes precedence inside startWorker.
			if rem := p.detector.backoffRemaining(time.Now()); rem > 0 {
				go func() {
					select {
					case <-time.After(rem):
						p.tokens <- struct{}{}
					case <-ctx.Done():
					}
				}()
				continue
			}
			startWorker()
		}
	}
}

// processOne reads a single message from the stream, processes the job, and
// acknowledges it. On unprocessable input it is acknowledged (skipped) so the
// stream does not stall; on transient failures it is left pending for redelivery.
//
// The workerOutcome return feeds the anomaly detector: outcomeFailed is
// returned when the whole turn died through Tor (the anomaly consumer decides
// whether the streak warrants a kill+remake).
func (p *Pool) processOne(ctx context.Context) workerOutcome {
	msgs, err := p.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    p.cfg.Group,
		Consumer: p.cfg.Consumer,
		Streams:  []string{p.cfg.Stream, ">"},
		Count:    1,
		Block:    5 * time.Second,
	}).Result()
	if err != nil || len(msgs) == 0 || len(msgs[0].Messages) == 0 {
		// A genuine idle read takes >= 5s (Block window). An instant return
		// means Redis/circuit trouble and becomes an idle-spin signal.
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Printf("worker stream read error: %v", err)
		}
		return outcomeIdle
	}

	msg := msgs[0].Messages[0]
	msgID := msg.ID

	// Decode fields into a job
	fields := make(map[string]string, len(msg.Values))
	for k, v := range msg.Values {
		fields[k] = fmt.Sprintf("%v", v)
	}
	jb, err := job.Parse(fields)
	if err != nil {
		log.Printf("invalid job %s: %v", msgID, err)
		p.rdb.XAck(ctx, p.cfg.Stream, p.cfg.Group, msgID)
		return outcomeOK // queue pollution, not a worker anomaly
	}

	setActiveJob(jb.JobID)
	log.Printf("[job %s] target=%s depth=%d terms=%d", jb.JobID, jb.Target, jb.DepthRemaining, len(jb.Terms))
	if err := p.processJob(ctx, jb); err != nil {
		// Whole turn died through Tor (see processJob). Ack so the job does not
		// loop forever in the pending queue, and mark it failed.
		log.Printf("[job %s] abnormal net outcome: %v", jb.JobID, err)
		p.rdb.XAck(ctx, p.cfg.Stream, p.cfg.Group, msgID)
		p.updateStatus(ctx, jb.JobID, "failed")
		return outcomeFailed
	}
	p.rdb.XAck(ctx, p.cfg.Stream, p.cfg.Group, msgID)
	p.updateStatus(ctx, jb.JobID, "completed")
	log.Printf("[job %s] status=completed written -> task_status:%s (acked %s)", jb.JobID, jb.JobID, msgID)
	return outcomeOK
}

// processJob crawls the job's frontier through Tor, runs every page through
// the content shield, extracts, and persists via storage-api.
//
// Every stage is logged to stdout so the SYS LOGS panel follows the job's
// exact journey: seed -> fetch attempt -> shield verdict -> extraction counts
// -> storage persist, plus a per-job DONE summary.
//
// It returns an error ONLY when the entire turn died through Tor (every
// reached URL failed to fetch) — the signal the anomaly detector uses to kill
// + remake workers on a fresh circuit when the streak persists.
func (p *Pool) processJob(ctx context.Context, jb *job.Job) error {
	// Stage counters for the per-job DONE summary line.
	fetched, blocked, fetchErrors, shieldErrors, discarded, persisted := 0, 0, 0, 0, 0, 0
	// Seed the frontier with the target host.
	front := newFrontier(p.cfg.MaxDepth)
	done := func() {
		// front.Len() after exhaustion is the total deduped URLs seen, i.e. the
		// breadth reached by the BFS landing within max_depth.
		log.Printf("[job %s] DONE urls_reached=%d fetched=%d blocked=%d fetch_errors=%d shield_errors=%d discarded=%d persisted=%d",
			jb.JobID, front.Len(), fetched, blocked, fetchErrors, shieldErrors, discarded, persisted)
	}
	// Normalize the target: preserve any scheme the analyst explicitly supplied
	// (a bare http:// vs https:// choice is meaningful to the target — see the
	// arsenal6 market below, which is HTTP-only and REFUSES https connections).
	// Only when no scheme was given do we pick a default (http:// for .onion:
	// the majority of hidden services, including this market class, serve plain
	// HTTP), and keep the other scheme as a fallback if that seed refuses.
	target := jb.Target
	seedScheme := ""
	if i := strings.Index(target, "://"); i >= 0 {
		seedScheme = target[:i+3] // includes "://"; lowercased below
		target = target[i+3:]
	}
	seedScheme = strings.ToLower(seedScheme)
	if seedScheme == "" {
		seedScheme = "http://"
	}
	seed := seedScheme + target
	front.Add(seed, 0)
	// One alternative-scheme retry per seed host (http→https / https→http) so a
	// scheme-less .onion that insists on the other protocol still gets crawled.
	altSeedsTried := map[string]bool{}
	log.Printf("[job %s] seed=%s max_depth=%d", jb.JobID, seed, p.cfg.MaxDepth)

	for {
		entry, ok := front.Next()
		if !ok {
			done()
			// Network-dead turn: every URL this worker reached failed through
			// Tor (e.g. the SOCKS path or the whole circuit is down). Returns
			// an error so the anomaly detector can kill + remake on a fresh
			// circuit when the streak persists.
			if fetched == 0 && fetchErrors > 0 {
				return fmt.Errorf("network-dead turn: 0/%d URLs fetched, %d fetch errors", front.Len(), fetchErrors)
			}
			return nil
		}

		t0 := time.Now()
		f, err := p.crawler.Fetch(ctx, entry.URL)
		if err != nil {
			fetchErrors++
			log.Printf("[job %s] fetch error %s: %v", jb.JobID, entry.URL, err)
			// Scheme fallback for seed hosts: an HTTP-only hidden service
			// refuses https (and vice-versa). If the alternate scheme is
			// untried for this exact URL, requeue it before giving up.
			if alt := alternateSchemeURL(entry.URL); alt != "" && !altSeedsTried[alt] {
				altSeedsTried[alt] = true
				log.Printf("[job %s] %s unreachable -> trying alternate scheme %s", jb.JobID, entry.URL, alt)
				front.Add(alt, entry.Depth)
			}
			continue
		}
		fetched++
		log.Printf("[job %s] fetch OK %s -> HTTP %d, %d bytes, %s, circ=%s",
			jb.JobID, entry.URL, f.HTTPStatus, len(f.RawBody), time.Since(t0).Round(time.Millisecond), f.TorCircuitID)

		// Politeness delay after each fetch, keyed on host.
		p.crawler.PoliteDelay(hostOf(entry.URL))

		// CAPTCHA / block → manual intervention queue.
		if f.Blocked {
			blocked++
			log.Printf("[job %s] BLOCK on %s (%s) -> manual queue", jb.JobID, entry.URL, f.BlockReason)
			p.enqueueManual(ctx, jb, f.BlockReason)
			done()
			return nil
		}

		// Content shield runs on raw bytes BEFORE parsing.
		f.JobID = jb.JobID
		verdict, err := p.shieldClient().Check(ctx, f)
		if err != nil {
			shieldErrors++
			log.Printf("[job %s] shield error %s: %v", jb.JobID, entry.URL, err)
			continue
		}
		if !verdict.ShieldPassed {
			discarded++
			log.Printf("[job %s] content discarded by shield: %s (body_ref=%s)", jb.JobID, verdict.ContentFlags[0].Category, verdict.BodyRef)
			continue // pipeline terminates here for this page
		}
		log.Printf("[job %s] shield PASS %s (body_ref=%s)", jb.JobID, entry.URL, verdict.BodyRef)

		// Parse + extract on clean bytes.
		extracted := extractText(f.RawBody)
		log.Printf("[job %s] extracted %s: links=%d handles=%d wallets=%d pgp=%d emails=%d infra=%d",
			jb.JobID, entry.URL, len(extracted.Links), len(extracted.Handles), len(extracted.Wallets),
			len(extracted.PGPKeys), len(extracted.Emails), len(extracted.InfraIndicators))

		// Enqueue discovered links for deeper crawling (BFS).
		for _, l := range extracted.Links {
			front.Add(crawler.NormalizeURL(l.URL), entry.Depth+1)
		}
		log.Printf("[job %s] frontier=%d URLs queued for depth<=%d", jb.JobID, front.Len(), entry.Depth+1)

		// Persist via storage-api (never direct DB access).
		payload := scanPayloadFrom(jb, f, verdict.BodyRef, extracted)
		if err := p.storageClient().PersistScan(ctx, payload); err != nil {
			log.Printf("[job %s] storage error: %v", jb.JobID, err)
		} else {
			persisted++
			log.Printf("[job %s] persisted %s -> storage-api", jb.JobID, entry.URL)
		}
	}
}
