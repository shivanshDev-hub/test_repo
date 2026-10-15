# Dark Web Threat Actor De-anonymization Platform — System Architecture Specification

**Version:** 1.0 (Engineering Design Draft)
**Target runtime:** Go 1.22+

## Overview

This document specifies the software architecture behind the platform's Collection → Processing → Analytics pipeline: automatically crawling dark web, deep web, and clearnet sources; extracting structured identifiers; and correlating those identifiers into candidate real-world actor profiles with evidence-backed confidence scores. It expands the high-level architecture into an implementation-level spec — repository layout, component contracts, data flow, I/O contracts, and failure handling.

**Scope.** This document specifies architecture, module boundaries, and interface contracts. It stops at the design/pseudocode level rather than a production-ready codebase, and treats de-anonymization as a human-reviewed, evidence-weighted output (candidate leads with confidence scores) rather than an automatic identity oracle.

**Non-goals.** Anti-bot/CAPTCHA evasion, exploitation of discovered misconfigurations, and legal/authorization workflows are out of scope for this spec — they belong to separate governance and reconnaissance documents referenced by, but not defined in, this one.

## Table of Contents
1. [Repository Structure](#1-repository-structure)
2. [Core Working Mechanism & Component Pipeline](#2-core-working-mechanism--component-pipeline)
3. [Input & Output Specifications](#3-input--output-specifications)
4. [Resilience & Error-Handling Strategy](#4-resilience--error-handling-strategy)

---

## 1. Repository Structure

```
deanon-platform/
├── cmd/
│   ├── server/                  # REST API + dashboard backend
│   │   └── main.go
│   ├── worker/                  # Standalone crawl/extraction/correlation worker
│   │   └── main.go
│   └── cli/                     # Operator CLI (scan launcher, admin tasks)
│       └── main.go
│
├── internal/                    # Application code, not importable outside this module
│   ├── domain/                  # Layer 0 — core types + ports (interfaces). No internal deps.
│   │   ├── actor.go             #   Actor, ActorLink, Identifier
│   │   ├── artifact.go          #   Artifact, ArtifactType
│   │   ├── scan.go              #   Scan, ScanStatus, ScanConfig
│   │   └── ports.go             #   ArtifactStore, GraphStore, RawObjectStore, ProxyDialer interfaces
│   │
│   ├── config/                  # Layer 1 — config loading/validation
│   │   ├── loader.go
│   │   └── schema.go
│   │
│   ├── proxy/                   # Layer 1 — Tor/SOCKS5 proxy + circuit management
│   │   ├── tor_controller.go
│   │   ├── circuit_pool.go
│   │   └── health_check.go
│   │
│   ├── storage/                 # Layer 1 — storage adapters (implement domain ports)
│   │   ├── objectstore/         #   Raw page/object storage (S3-compatible)
│   │   ├── postgres/            #   Structured artifact + scan metadata storage
│   │   └── graphdb/              #   Neo4j-backed actor/identifier graph
│   │
│   ├── crawler/                 # Layer 2 — crawl engine (depends on domain + proxy port)
│   │   ├── frontier.go
│   │   ├── worker_pool.go
│   │   ├── fetcher.go
│   │   └── ratelimit.go
│   │
│   ├── parser/                  # Layer 2 — extraction pipeline (depends on domain)
│   │   ├── extractor.go
│   │   ├── patterns.go
│   │   ├── validators.go
│   │   └── normalizer.go
│   │
│   ├── correlation/             # Layer 3 — de-anonymization/correlation engine
│   │   ├── resolver.go
│   │   ├── stylometry.go
│   │   ├── infra_link.go
│   │   └── confidence.go
│   │
│   └── api/                     # Layer 4 — HTTP handlers, middleware, DTOs
│       ├── handlers/
│       ├── middleware/
│       └── dto/
│
├── pkg/                          # Reusable, externally-importable libraries
│   ├── retry/                    # Generic backoff/retry helper
│   ├── ratelimiter/               # Token-bucket limiter
│   └── logctx/                    # Structured logging w/ correlation IDs
│
├── configs/
│   ├── config.example.yaml
│   └── torrc.template
│
├── scripts/
│   ├── setup_tor.sh
│   ├── migrate_db.sh
│   └── seed_test_data.sh
│
├── migrations/                     # SQL schema migrations (goose/golang-migrate)
├── deployments/
│   ├── docker-compose.yaml
│   └── k8s/
├── go.mod
└── go.sum
```

### Dependency direction

Packages depend on **interfaces defined in `domain`**, never on each other's concrete types. `correlation`, for example, takes a `domain.ArtifactStore` interface in its constructor; at wiring time in `cmd/server/main.go`, a concrete `storage/postgres.Store` is injected. This is what keeps `storage` from ever needing to import `correlation`, and is the standard fix any time two internal packages want to import each other.

| Layer | Packages | May import | Must not import |
|---|---|---|---|
| 0 | `internal/domain` | stdlib only | any other `internal/*` package |
| 1 | `config`, `proxy`, `storage/*` | `domain` | `crawler`, `parser`, `correlation`, `api` |
| 2 | `crawler`, `parser` | `domain`, proxy/storage **interfaces only** | each other's concrete internals; `correlation`, `api` |
| 3 | `correlation` | `domain`, storage ports, `parser` output types | `api` |
| 4 | `api`, `cmd/*` | everything | — (nothing imports layer 4) |

---

## 2. Core Working Mechanism & Component Pipeline

### 2.1 End-to-End Data Flow

```
CLI / REST API
   → Scan Queue (job dispatch)
   → URL Frontier (dedup + depth tracking)
   → Crawler Worker Pool (fetch via Tor/proxy)
   → Raw Object Store (immutable, content-addressed)
   → Extraction Pipeline (parser)
   → Structured Storage (Postgres + Graph DB)
   → Correlation Engine (entity resolution + confidence scoring)
   → Analyst Review (dashboard/API)
   → Reports, Alerts, Exports
```

1. **Ingress.** CLI or REST validates the request against `ScanConfig` (target, depth, proxy params, scope) and persists a `Scan{status: queued}` row.
2. **Dispatch.** The scan is pushed onto the job queue (RabbitMQ/Kafka topic `scans.queued`); a `worker` process claims it.
3. **Frontier seeding.** The target is normalized (scheme, trailing slash, `.onion` canonicalization) and inserted into the `URLFrontier` at depth 0.
4. **Crawl.** The worker pool drains the frontier through the proxy layer. Each successful fetch is written to the raw object store; discovered same-scope links are pushed back onto the frontier up to `max_depth`/`max_pages`.
5. **Extraction.** As raw objects land, the parser pipeline runs registered extractors over normalized text, producing `Artifact` records with provenance (source URL, scan ID, timestamp, snippet).
6. **Structured persistence.** Artifacts are written to Postgres (queryable metadata) and upserted as nodes/edges in the graph store.
7. **Correlation.** The resolver clusters artifacts into candidate `ActorLink`s, attaches evidence, computes a 0–100 confidence score, and writes `status=candidate` links to the graph.
8. **Review & output.** Analysts confirm/reject candidates via the dashboard or API; confirmed links, reports, and alerts are exported. Nothing is auto-published as a confirmed identity — see §2.5.

### 2.2 Proxy & Tor Management

**Responsibilities:** route all outbound crawl traffic through Tor (and, for clearnet sources, an optional forward-proxy pool); distribute load across multiple circuits; detect and replace unhealthy circuits.

**Components**
- `TorController` — owns the connection to Tor's control port (default `9051`), authenticates, and issues control-protocol commands (`SIGNAL NEWNYM` to request new circuits for future streams — see Tor's `control-spec` for the full command reference).
- `CircuitPool` — a set of `IsolatedDialer`s, each using a distinct SOCKS5 username so Tor's stream isolation (`IsolateSOCKSAuth`) forces each identity onto its own circuit instead of reusing one circuit for every request.
- `HealthChecker` — periodically dials a known-reachable check target through each dialer; a dialer failing N consecutive checks is marked unhealthy and quarantined pending renewal.

Minimal `torrc` this assumes:
```
SOCKSPort 9050
ControlPort 9051
CookieAuthentication 1
```

```go
// internal/proxy/tor_controller.go
func (c *TorController) RenewCircuit(ctx context.Context, identity string) error {
    conn, err := net.DialTimeout("tcp", c.controlAddr, 5*time.Second)
    if err != nil {
        return fmt.Errorf("control port dial: %w", err)
    }
    defer conn.Close()

    if err := c.authenticate(conn); err != nil { // cookie or hashed-password auth, per torrc
        return fmt.Errorf("control port auth: %w", err)
    }
    if _, err := fmt.Fprintf(conn, "SIGNAL NEWNYM\r\n"); err != nil {
        return err
    }
    return c.expectOK(conn)
}
```

```go
// internal/proxy/circuit_pool.go
type CircuitPool struct {
    mu      sync.Mutex
    dialers []*IsolatedDialer
    next    int
}

type IsolatedDialer struct {
    Dial         proxy.Dialer // SOCKS5 dialer bound to a unique auth identity
    RequestCount int
    LastRotated  time.Time
    Healthy      bool
}

func (p *CircuitPool) Next() *IsolatedDialer {
    p.mu.Lock()
    defer p.mu.Unlock()
    d := p.dialers[p.next%len(p.dialers)]
    p.next++
    return d
}
```

**Rotation policy** (configurable): renew a dialer's circuit after `N` requests, after a max age, or immediately after `k` consecutive failures attributable to that circuit (as opposed to the target genuinely being down).

### 2.3 Crawler Engine

**`URLFrontier`** — a durable, prioritized queue seeded from scan targets. Responsibilities: dedup (hash of the normalized URL, backed by a unique DB constraint or an in-memory Bloom filter fast-path), depth tracking, per-host politeness state.

**`WorkerPool`** — fixed-size goroutine pool bounded by a semaphore, with `sync.WaitGroup` tracking in-flight fetches so the crawl phase of a scan knows when it's actually done (as opposed to just "queue empty," which can race with in-flight workers still discovering new URLs).

```go
// internal/crawler/worker_pool.go
type WorkerPool struct {
    sem  chan struct{}
    wg   sync.WaitGroup
    jobs chan Job
}

func NewWorkerPool(size, queueDepth int) *WorkerPool {
    return &WorkerPool{
        sem:  make(chan struct{}, size),
        jobs: make(chan Job, queueDepth),
    }
}

func (p *WorkerPool) Start(ctx context.Context, handle func(context.Context, Job)) {
    for {
        select {
        case job, ok := <-p.jobs:
            if !ok {
                return
            }
            p.sem <- struct{}{}
            p.wg.Add(1)
            go func(j Job) {
                defer func() {
                    <-p.sem
                    p.wg.Done()
                    if r := recover(); r != nil {
                        logctx.From(ctx).Error("worker panic recovered", "job", j.ID, "panic", r)
                    }
                }()
                handle(ctx, j)
            }(job)
        case <-ctx.Done():
            return
        }
    }
}

func (p *WorkerPool) Submit(j Job) { p.jobs <- j }
func (p *WorkerPool) Wait()        { p.wg.Wait() }
```

**`Fetcher`** — performs the HTTP(S) request through a `CircuitPool`-provided dialer, bounded by a per-request `context.WithTimeout`, and streams the response body into the object store rather than buffering large pages fully in memory.

**`ratelimit.go`** — per-host token bucket; backs off on `429`/`503`; enforces a configured minimum inter-request delay per host, independent of overall worker concurrency.

### 2.4 Parsing & Extraction Pipeline

Each raw object is normalized to text (HTML → text extraction, charset normalization) and run through a chain of `Extractor`s:

```go
// internal/parser/extractor.go
type Extractor interface {
    Type() domain.ArtifactType
    Extract(text string) []domain.Artifact
}
```

Registered extractors cover the identifier categories relevant to actor linking: handles/usernames, emails, cryptocurrency addresses, PGP key material, and `.onion`/domain references. Pattern matching is the first pass; every extractor is paired with a validator that filters false positives before an artifact is persisted.

```go
// internal/parser/patterns.go — illustrative, not exhaustive
var patterns = map[domain.ArtifactType]*regexp.Regexp{
    ArtifactEmail:    regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`),
    ArtifactBTC:      regexp.MustCompile(`\b(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}\b`),
    ArtifactXMR:      regexp.MustCompile(`\b4[0-9AB][1-9A-HJ-NP-Za-km-z]{93}\b`),
    ArtifactPGPBlock: regexp.MustCompile(`-----BEGIN PGP PUBLIC KEY BLOCK-----[\s\S]+?-----END PGP PUBLIC KEY BLOCK-----`),
    ArtifactOnion:    regexp.MustCompile(`\b[a-z2-7]{16,56}\.onion\b`),
}
```

```go
// internal/parser/validators.go
func ValidateBTC(candidate string) bool { return base58check.Verify(candidate) }
func ValidatePGPBlock(b string) (fingerprint string, ok bool) { /* parse via RFC 4880 armor, don't just delimiter-match */ }
```

Two accuracy notes worth designing around explicitly:

- **Regex over-generates.** A bare pattern match is a *candidate*, not a fact — Base58Check validation cuts most accidental BTC-shaped strings; PGP blocks should be parsed (not just delimiter-matched) to pull a real fingerprint. Free-text PII beyond structured identifiers (names, physical addresses) generally needs an NER model rather than regex, and NER output is probabilistic — route it through the confidence/review path in §2.5 rather than treating it as ground truth.
- **Monero addresses are not transaction graphs.** Extracting a Monero address only captures the fact that someone *posted* that address as text (e.g., in a donation request). It does not de-anonymize on-chain Monero activity — Monero's stealth-address design specifically prevents linking a posted public address to individual transactions from the ledger alone. Don't let an "address extracted" stat imply more than that.

Normalization attaches provenance to every surviving artifact: `source_url`, `scan_id`, `first_seen`, and a short `context_snippet` for analyst review.

### 2.5 De-anonymization & Correlation Engine

**Goal:** cluster artifacts into candidate `ActorLink`s and score how much evidence supports treating them as the same actor — as *leads for a human analyst*, not as automatic conclusions.

```go
// internal/correlation/resolver.go
type Identifier struct {
    Type  domain.ArtifactType
    Value string
}

type LinkEvidence struct {
    Method string  // "exact_match" | "stylometry" | "infra_overlap"
    Detail string
    Weight float64
}

type ActorLink struct {
    ActorID     string
    Identifiers []Identifier
    Evidence    []LinkEvidence
    Confidence  float64 // 0-100
    Status      string  // "candidate" | "analyst_confirmed" | "rejected"
}

type Resolver interface {
    Resolve(ctx context.Context, artifacts []domain.Artifact) ([]ActorLink, error)
}
```

**Resolution runs in two passes:**

1. **Deterministic clustering** on exact-match keys: a shared PGP fingerprint, a reused wallet address, an identical handle across two sources. These are high-confidence, explainable signals.
2. **Probabilistic corroboration**, which can *raise* confidence on an existing cluster but — by design — never creates a link on its own:
   - **Stylometry** (`stylometry.go`) — extracts a feature vector per authored-text sample (function-word frequencies, average sentence length, type-token ratio, punctuation/emoji habits — the standard "Writeprints"-style feature set from the authorship-attribution literature) and compares vectors across handles with cosine similarity.
   - **Infrastructure correlation** (`infra_link.go`) — shared hosting/ASN, SSL certificate fingerprint reuse, matching mail-server headers, where such technical artifacts are available.

```go
// internal/correlation/stylometry.go
type StyleVector struct {
    AvgSentenceLen float64
    TypeTokenRatio float64
    FunctionWords  map[string]float64
    Punctuation    map[string]float64
}

func ExtractStyleVector(text string) StyleVector { /* feature extraction */ }
func Similarity(a, b StyleVector) float64        { /* cosine similarity */ }
```

**Confidence scoring** is a weighted sum over evidence, clamped to `[0, 100]`:

```go
// internal/correlation/confidence.go
func Score(evidence []LinkEvidence) float64 {
    var total float64
    for _, e := range evidence {
        total += e.Weight
    }
    return math.Min(total, 100)
}
```

Illustrative starting weights (tune against labeled evaluation data before relying on the absolute numbers):

| Evidence type | Example weight | Why |
|---|---|---|
| Exact PGP fingerprint match | 40 | Cryptographic identity material, very low false-positive rate |
| Reused wallet address | 35 | Reuse is a deliberate or careless act, strong signal |
| Shared unique infra artifact (cert fingerprint) | 25 | Rare enough to be meaningful, but can reflect shared hosting |
| Stylometric similarity above threshold | 10–15 | Corroborating only; too weak alone to found a link |
| Shared generic infra (same hosting provider) | 3–5 | Weak; many unrelated actors share large providers |

**Human-in-the-loop gate.** No `ActorLink` reaches `analyst_confirmed` — and nothing is exported to reports/alerts as an identification — without explicit analyst sign-off. This is an architectural invariant, not just a policy note: `api/handlers` should not expose any path that flips `Status` to `analyst_confirmed` except one requiring an authenticated analyst action, and that action should be audit-logged (§4.6).

---

## 3. Input & Output Specifications

### 3.1 CLI

```
deanon-cli scan start --target <url|.onion> [flags]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--target` | string | *(required)* | Target URL or `.onion` address |
| `--depth` | int | `2` | Max link-following depth |
| `--max-pages` | int | `500` | Hard cap on pages per scan |
| `--proxy-port` | int | `9050` | Tor SOCKS5 port |
| `--control-port` | int | `9051` | Tor control port |
| `--concurrency` | int | `10` | Worker pool size |
| `--timeout` | duration | `30s` | Per-request timeout |
| `--scope` | enum | `domain` | `domain` \| `subdomain` \| `explicit` |
| `--output` | path | `-` | Report destination (`-` = stdout) |
| `--format` | enum | `json` | `json` \| `ndjson` |

| Subcommand | Purpose |
|---|---|
| `scan start` | Launch a scan |
| `scan status <id>` | Poll status |
| `scan report <id>` | Fetch the completed report |
| `scan cancel <id>` | Request graceful cancellation |
| `config validate` | Validate a config file before deploy |
| `proxy test` | Verify Tor connectivity + circuit renewal |

### 3.2 REST API

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/v1/scans` | Create/enqueue a scan |
| `GET` | `/api/v1/scans/{id}` | Get scan status |
| `GET` | `/api/v1/scans/{id}/report` | Get full report |
| `POST` | `/api/v1/scans/{id}/cancel` | Cancel a running scan |
| `POST` | `/api/v1/actors/{id}/confirm` | Analyst confirms a candidate link (audited) |
| `GET` | `/healthz` | Liveness probe |

`POST /api/v1/scans` request:
```json
{
  "target": "http://exampleonionaddress1234567890abcdef.onion",
  "depth": 2,
  "max_pages": 500,
  "proxy": {
    "socks5_port": 9050,
    "control_port": 9051,
    "circuit_rotation_interval_requests": 25
  },
  "concurrency": 10,
  "timeout_seconds": 30,
  "scope": "domain",
  "callback_url": "https://internal.example.org/hooks/scan-complete"
}
```

`202 Accepted`:
```json
{
  "scan_id": "b6f1c9de-4a2e-4b7a-9c1a-6e2f8d3a9b10",
  "status": "queued",
  "created_at": "2026-08-29T10:15:00Z"
}
```

### 3.3 JSON Report Format

```json
{
  "scan_id": "b6f1c9de-4a2e-4b7a-9c1a-6e2f8d3a9b10",
  "target": "http://exampleonionaddress1234567890abcdef.onion",
  "status": "completed",
  "created_at": "2026-08-29T10:15:00Z",
  "started_at": "2026-08-29T10:15:04Z",
  "completed_at": "2026-08-29T10:22:41Z",
  "stats": {
    "pages_crawled": 187,
    "pages_failed": 6,
    "artifacts_extracted": 42,
    "actor_links_found": 3,
    "circuit_rotations": 9
  },
  "artifacts": [
    {
      "artifact_id": "a-0001",
      "type": "pgp_fingerprint",
      "value": "EXAMPLE-FPR-0000-0000-0000-0000-0000-0000",
      "source_url": "http://exampleonionaddress.../contact",
      "first_seen": "2026-08-29T10:16:12Z",
      "context_snippet": "PGP fingerprint listed on contact page"
    },
    {
      "artifact_id": "a-0002",
      "type": "btc_address",
      "value": "1ExampleAddressNotReal00000000000",
      "source_url": "http://exampleonionaddress.../donate",
      "first_seen": "2026-08-29T10:17:03Z",
      "context_snippet": "Donation address on landing page"
    }
  ],
  "actor_links": [
    {
      "actor_id": "actor-001",
      "identifiers": ["handle:example_handle", "pgp_fingerprint:EXAMPLE-FPR-..."],
      "confidence": 78,
      "status": "candidate",
      "evidence": [
        {"method": "exact_match", "detail": "shared PGP fingerprint across 2 sources", "weight": 40},
        {"method": "stylometry", "detail": "similarity 0.86 vs handle corpus", "weight": 15}
      ]
    }
  ],
  "errors": [
    {"url": "http://exampleonionaddress.../forum?p=44", "stage": "fetch", "message": "context deadline exceeded", "retryable": true}
  ]
}
```

| Field | Type | Notes |
|---|---|---|
| `status` | enum | `queued` \| `running` \| `completed` \| `partial` \| `failed` |
| `stats.circuit_rotations` | int | Count of forced Tor circuit renewals during the scan |
| `artifacts[].type` | enum | `email` \| `btc_address` \| `xmr_address` \| `pgp_fingerprint` \| `handle` \| `onion_link` |
| `actor_links[].status` | enum | `candidate` \| `analyst_confirmed` \| `rejected` — only ever set to `analyst_confirmed` via the audited confirm endpoint |
| `errors[].retryable` | bool | Whether the pipeline itself already retried this before giving up |

---

## 4. Resilience & Error-Handling Strategy

### 4.1 Failure Classification & Retry

Every network-facing call classifies its error before deciding whether to retry: transient (timeout, connection reset, `5xx`, proxy handshake failure) vs. permanent (`404`, DNS failure for a non-existent `.onion`, malformed target). Only transient errors are retried.

```go
// pkg/retry/retry.go
type Config struct {
    MaxAttempts int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Jitter      float64 // fraction of delay, 0.0–1.0
}

func Do(ctx context.Context, cfg Config, fn func(ctx context.Context) error) error {
    var lastErr error
    for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
        err := fn(ctx)
        if err == nil {
            return nil
        }
        if !isRetryable(err) {
            return err
        }
        lastErr = err

        delay := time.Duration(math.Min(
            float64(cfg.BaseDelay)*math.Pow(2, float64(attempt)),
            float64(cfg.MaxDelay),
        ))
        delay += time.Duration(rand.Float64() * cfg.Jitter * float64(delay))

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(delay):
        }
    }
    return fmt.Errorf("giving up after %d attempts: %w", cfg.MaxAttempts, lastErr)
}
```

### 4.2 Tor Circuit & Proxy Handshake Failures

Circuit-level failures are handled one layer below the generic retry above, because the right response to "this circuit is bad" is different from "retry the same request":

1. A fetch through a given `IsolatedDialer` fails with a proxy/handshake-level error (as opposed to the *target* returning an HTTP error).
2. The dialer's failure counter increments; on reaching a threshold (default 3 consecutive), the dialer is marked unhealthy and pulled from `CircuitPool.Next()` rotation.
3. `TorController.RenewCircuit` is called for that identity (`SIGNAL NEWNYM`); the dialer is re-enabled once a subsequent health check succeeds.
4. The original job is re-queued (via the standard retry path in §4.1) onto a *different* healthy dialer rather than blindly retried on the same one — retrying a bad circuit on itself just reproduces the failure.

### 4.3 Context Timeout Handling

Timeouts are layered so a stuck request can't hang a worker indefinitely, and a stuck worker can't hang a whole scan:

```go
// per-request timeout
func (f *Fetcher) Fetch(ctx context.Context, url string) (*Response, error) {
    ctx, cancel := context.WithTimeout(ctx, f.cfg.RequestTimeout)
    defer cancel()
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    return f.do(req)
}
```

```go
// scan-level deadline, propagated to every worker via ctx
scanCtx, cancel := context.WithTimeout(parentCtx, scan.Config.MaxDuration)
defer cancel()
pool.Start(scanCtx, handleJob)
```

Cancellation propagates downward: canceling `scanCtx` unblocks every in-flight `select` in the worker pool and fetcher without needing to track goroutines individually.

### 4.4 Panic Recovery & Bounded Concurrency

Every worker goroutine and every HTTP handler recovers from panics locally, so one bad page (malformed HTML, an extractor hitting an edge case) can't crash the process or silently wedge the `WaitGroup`:

```go
// internal/api/middleware/recover.go
func Recover(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                logctx.From(r.Context()).Error("handler panic", "err", rec, "stack", string(debug.Stack()))
                http.Error(w, `{"error":"internal_error"}`, http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

Concurrency is always bounded, never spawned per-URL unboundedly: the semaphore in `WorkerPool` (§2.3) caps in-flight fetches, and the frontier's channel is buffered rather than unbounded, so a burst of discovered links applies backpressure to the crawl instead of ballooning goroutine/memory usage. The same pattern bounds the extraction and correlation stages.

### 4.5 Circuit Breakers for Downstream Dependencies

Postgres, the graph store, and the object store each sit behind a simple breaker: after a failure-rate threshold trips, calls fail fast instead of piling up waiting on a dependency that's down; a half-open probe periodically checks whether the dependency has recovered before fully closing the breaker again. This keeps one degraded storage dependency from backing up the entire worker pool — a scan can still make crawl progress and buffer results briefly while, say, the graph store is unavailable.

### 4.6 Observability

Every log line and metric carries a `scan_id` (and, within a scan, a `job_id`) correlation ID end-to-end, so a single failed page can be traced from frontier → fetch → extraction → storage. Key metrics to export (Prometheus-style): `crawler_fetch_duration_seconds`, `crawler_fetch_errors_total{reason}`, `proxy_circuit_rotations_total`, `parser_artifacts_extracted_total{type}`, `correlation_links_created_total`, `worker_pool_inflight`, `worker_pool_queue_depth`. Analyst-facing actions (confirming/rejecting an `ActorLink`, exporting a report) are written to an append-only audit log distinct from operational logs, since that trail is itself part of the platform's accountability model.
