# OnionEye

> **📋 Team Structure:** See [TEAM_STRUCTURE.md](./TEAM_STRUCTURE.md) for current team assignments and responsibilities.

<p align="center">
  <img src="docs/architecture/assets/official_logo.svg" alt="OnionEye Logo" width="420"/>
</p>

<h1 align="center">🕵️♂️ OnionEye</h1>
<p align="center">
  <b>Hardened Dark-Web Reconnaissance, Entity Correlation & Adversary De-Anonymization Platform</b>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Security-Zero--Trust%20Storage-00ff9d?style=for-the-badge&logo=shield" alt="Zero Trust">
  <img src="https://img.shields.io/badge/Egress-Hardened%20Tor%20SOCKS5-7d4698?style=for-the-badge&logo=torproject" alt="Tor Network">
  <img src="https://img.shields.io/badge/Safety-CSAM%20Zero--Decode%20Shield-ff1744?style=for-the-badge" alt="Content Shield">
  <img src="https://img.shields.io/badge/Standard-OASIS%20STIX%202.1-00b0ff?style=for-the-badge" alt="STIX 2.1">
  <img src="https://img.shields.io/badge/Engine-Mathematical%20Dampening-ff9100?style=for-the-badge" alt="Dampening Engine">
</p>

---

## 🎯 Executive Overview & Problem Statement

Modern dark-web threat intelligence is broken by three systemic vulnerabilities:

1. **The False-Attribution Trap:** Threat actors routinely change aliases across forums and marketplaces. Traditional link-analysis engines naively group entities based on co-occurring artifacts, generating catastrophic false-positive clusters by attributing shared infrastructure (Tor exit nodes, coin tumblers, forum signature templates) to individual threat actors.
2. **The Contamination & Legal Liability:** Passive scraping of unindexed `.onion` nodes constantly risks ingesting CSAM and illicit media into government or enterprise databases, creating severe criminal liability and analyst trauma.
3. **The Egress Leak Threat:** Traditional web scrapers operating over proxies frequently suffer from clearnet DNS leaks or kernel-level TCP bypasses, de-anonymizing intelligence operations to hostile server administrators.

**OnionEye** is a hardened, zero-trust dark-web intelligence gathering and persona-correlation platform. It provides cyber defense teams, intelligence units, and financial investigators with **provably safe, mathematically dampened, and standards-compliant threat intelligence**.

---

## ⚡ Core Innovations & Breakthrough Features

### 1. 🛡️ Zero-Decode Content Safety Shield (Stream-Level Protection)
- **Zero-Decode Guarantee:** Raw byte streams are hashed and screened *before* memory allocation, rendering, or parsing. Illicit media is never decoded into memory.
- **Immediate RAM Sanitization:** If a perceptual/exact match hits a known-bad hash set, the in-memory buffer is instantly zeroed out (`phash.Zero(b)`), and a cryptographic discard audit record is persisted without ever touching disk.
- **Secondary ML Fail-Open Pipeline:** Unseen media streams pass through a localized classifier container for computer-vision safety verification before any textual parsing is permitted.

### 2. 🧮 Two-Stage Negative Control & Composite Attribution Engine
Unlike heuristic tools that inflate confidence on sheer link count, OnionEye applies rigorous mathematical noise dampening:
- **Stage 1 (Negative Control State Machine):** Evaluates whether extracted artifacts originate from shared/benign darknet infrastructure. Tracks state transitions (`PASSED`, `FLAGGED`, `INDETERMINATE`) across baseline health feeds.
- **Stage 2 (Time-Decayed Composite Confidence):** Computes confidence scores $C_{\text{final}} \in [0.0, 1.0]$ using TF-IDF inverse frequency and an exponential half-life decay:

$$\text{Confidence} = \min\left(1.0, \; \frac{\sum_{i=1}^{k} \left( S_i \cdot \text{IDF}_i \cdot 0.5^{(\Delta t_i / 30\text{d})} \right)}{\max\left(1.0, \; \log_2(1 + \text{NC}_{\text{passed}})\right)} \right)$$

- **30-Day Half-Life Decay:** Digital footprints (PGP keys, wallets, emails) naturally decay in attribution weight if not re-observed.
- **Logarithmic Frequency Dampening:** Shared artifacts observed across dozens of clean sites automatically have their attribution score dampened toward zero.

### 3. 🔒 Zero-Leak Hardened Tor Egress Sidecar
- **Kernel-Level Egress Enforcement:** All crawler containers run inside an isolated network sandbox where `iptables` rules drop all non-loopback outbound traffic except through the dedicated Tor SOCKS5 proxy (`tor-sidecar:9050`).
- **Unprivileged Execution:** Tor daemons execute under unprivileged `debian-tor` (UID 100) users with automated circuit renewal upon timeout or connection reset.

### 4. 🗄️ Zero-Trust Storage & Vector Analytics (PostgreSQL + pgvector)
- **Isolated Storage Boundary:** `storage-api` is the **only** microservice in the entire architecture that holds database credentials. All crawler, engine, and dashboard interactions flow through authenticated REST contracts.
- **Time-Partitioned Scans:** High-throughput telemetry and flag logs are partitioned into monthly physical tables (`PARTITION BY RANGE (started_at)`).
- **Stylometric HNSW Vector Matching:** Author writing fingerprints are embedded into 384-dimensional dense vectors and queried using pgvector HNSW cosine similarity.
- **Recursive SQL Graph Clustering:** The `v_canonical_actors` view executes recursive Common Table Expressions (CTEs) across active `filter_policies` to resolve fragmented alias graphs into unified canonical actor roots.

### 5. 🌐 OASIS STIX 2.1 Standardized Threat Intel Export
- Every identified entity, wallet, alias link, and persona cluster can be exported as standard **OASIS STIX 2.1 JSON bundles**.
- Directly consumable by enterprise TIP/SIEM platforms (OpenCTI, MISP, Splunk, Microsoft Sentinel).

### 6. 📊 Real-Time 2D Canvas Intelligence HUD
- GPU-accelerated 2D HTML5 Canvas graph visualizer supporting dynamic force-directed node physics, actor clustering, and negative-control state inspection.
- Filter policies can be dynamically modified at runtime to set confidence cutoffs and allowlist/denylist rules without re-indexing historical raw data.

---

## 🏛️ System Architecture

```
                                 [ TARGET DARKNET ]
                                          │ (.onion)
                                          ▼
                         ┌─────────────────────────────────┐
                         │   tor-sidecar (SOCKS5 Proxy)    │
                         │   • iptables Zero-Leak Egress   │
                         │   • Automated Circuit Rotation  │
                         └────────────────┬────────────────┘
                                          │
                                          ▼
                         ┌─────────────────────────────────┐
                         │   crawler-worker (Go Engine)    │
                         │   • Multi-worker Concurrency    │
                         │   • Regex Artifact Extraction   │
                         └────────────────┬────────────────┘
                                          │ (Raw Stream)
                                          ▼
                         ┌─────────────────────────────────┐
                         │    content-shield (Go / ML)     │
                         │   • Zero-Decode Raw Hashing     │
                         │   • RAM Sanitization (phash.Zero)
                         │   • Python ML Classification    │
                         └────────────────┬────────────────┘
                                          │ (Clean Scans)
                                          ▼
┌──────────────────────────────────────────────────────────────────────────────────┐
│                             storage-api (Go / pgx)                               │
│       * Sole Credential Holder — Enforces Zero-Trust Database Isolation *        │
└───────────────┬──────────────────────────────────────────────────┬───────────────┘
                │                                                  │
                ▼                                                  ▼
┌────────────────────────────────┐                ┌────────────────────────────────┐
│   PostgreSQL 16 + pgvector     │                │   correlation-engine (Go)      │
│   • Time-Partitioned Scans     │◄───────────────┤   • Negative Control FSM       │
│   • HNSW Stylometry Index      │                │   • Half-Life Time Decay       │
│   • Recursive Actor Views (CTE)│                │   • Logarithmic Dampening      │
└────────────────────────────────┘                └────────────────┬───────────────┘
                                                                   │
                                                                   ▼
                                                  ┌────────────────────────────────┐
                                                  │ dashboard-api & React 2D HUD   │
                                                  │ • OASIS STIX 2.1 Exporter      │
                                                  │ • Interactive Force Graph      │
                                                  └────────────────────────────────┘
```

---

## 🚦 Microservice Architecture & Port Allocations

| Service | Technology | Port | Primary Responsibility |
|---|---|:---:|---|
| **`nginx`** | Nginx Reverse Proxy | `80` | Single-entrypoint gateway for Dashboard UI, Ingestion API, and Dashboard API. |
| **`frontend`** | React / Vite / 2D Canvas HUD | `80` (proxied) / `5173` (dev) | Interactive analyst HUD, force-directed graph console, and STIX 2.1 exporter. |
| **`ingestion-api`** | Go (Gin / Redis) | `8080` | 10-phase input safety validation, sliding-window rate limiter, and task dispatcher. |
| **`storage-api`** | Go (pgx/v5) | `8081` | Zero-trust CRUD gateway; sole PostgreSQL database credential holder. |
| **`content-shield`** | Go / SHA-256 / FastScan | `8082` | Stream-level byte inspection & instant memory zeroing (`phash.Zero`). |
| **`correlation-engine`**| Go (Math & Vector Graph) | `8083` | Multi-modal hard scoring, 384-d stylometry, time decay, and negative-control dampening. |
| **`dashboard-api`** | Go (Gin / pgx) | `3000` | REST API for graph telemetry, metrics, lead reviews, and forensic snapshots. |
| **`crawler-worker`** | Go / Tor SOCKS5 | Internal | Autonomous hidden-service discovery, anti-obfuscation normalizer & artifact extractor. |
| **`tor-sidecar`** | Alpine / Tor / iptables | `9050`, `9051` | Hardened SOCKS5 egress & Tor control port with automated circuit renewal. |
| **`postgres`** | PostgreSQL 16 + pgvector | `5432` | Relational audit store, HNSW stylometric vector index, and time-partitioned telemetry. |
| **`redis`** | Redis 7 Alpine | `6379` | In-memory message streams (`crawl:jobs`), rate limits, and session context cache. |
| **`ollama`** *(Optional)* | Ollama LLM Runtime | `11434` | Local inference engine for query expansion (`qwen2.5:3b`) & safety (`llama-guard3:1b`). |

---

## 🚀 Quickstart & Deployment

### Prerequisites
- **Docker Engine 24.0+** & **Docker Compose v2**
- Minimum 4 GB RAM (8 GB recommended for `--profile ai`)

### 1. Launch OnionEye in Standard Mode
```bash
# Clone the platform repository
git clone https://github.com/sudo-problems/OnionEye.git
cd OnionEye/OnionEye_Code

# Build and start all core zero-trust services in detached mode
docker compose up -d --build
```

### 2. Launch with Local AI Expansion & Safety Guard (Optional)
```bash
# Launches the Ollama local AI model container alongside core services
docker compose --profile ai up -d --build
```

### 3. Verify System Health & UI
```bash
# Inspect overall reverse proxy and service health
curl -s http://localhost/health | jq .
curl -s http://localhost:8080/healthz | jq .
curl -s http://localhost:8081/health | jq .
curl -s http://localhost:8082/health | jq .
curl -s http://localhost:8083/healthz | jq .
curl -s http://localhost:3000/api/v1/dashboard/health | jq .

# Access the analyst dashboard in your browser
open http://localhost/
```

---

## 🔬 Mathematical Verification & Unit Tests

All attribution formulas, exponential decay schedules, stylometric cosine similarity, and negative control state transitions are rigorously validated with automated test suites:

```bash
# Ingestion API validation, rate limiting, and Tier-1 safety guard tests
cd OnionEye_Code/services/ingestion-api && go test -v ./...

# Crawler anti-obfuscation normalization & extraction tests
cd ../crawler-worker && go test -v ./internal/extract

# Multi-modal vector scoring, stylometry cosine similarity & negative control tests
cd ../correlation-engine && go test -v ./internal/engine ./internal/negcontrol ./internal/handler
```

**Test Coverage Highlights:**
- `TestNegativeControl_Transitions`: Verifies deterministic state machine movement between `PASSED`, `FLAGGED`, and `INDETERMINATE`.
- `TestCompositeScore_TimeDecay`: Validates that a 60-day-old artifact has its confidence decayed by exactly $0.5^{(60/30)} = 0.25$.
- `TestCompositeScore_CleanDampening`: Verifies that high benign hit counts ($NC_{\text{passed}} = 7$) logarithmically suppress confidence to prevent false linkages.
- `TestTermFallbackHandle` (storage-api): Verifies investigation terms with no extractable handles are sanitized into distinct `TERM:`-prefixed fallback actor handles — 11 cases covering whitespace, URL fragments, unicode, >40-char truncation, and pure-noise rejection.

---

## 📊 Live Dashboard & Seed Data

The dashboard now displays **live data from PostgreSQL** — zero mock/placeholder data remains. All dashboard panels are backed by real database queries:

| Panel | Data Source | Backend Path |
|---|---|---|
| **Stats Bar** | `actors`, `v_active_persona_links`, `identifiers`, `content_flags`, `scans` | `storage-api /stats` → `dashboard-api /stats` |
| **Graph (RADAR)** | `persona_links` + `actors` + `identifiers` + `infra_indicators` | `storage-api /actors` + `/leads` → `dashboard-api /graph` (all actors render, including standalone `TERM_FALLBACK` nodes) |
| **Leads** | `v_active_persona_links` (policy-filtered) | `storage-api /leads` → `dashboard-api /leads` (enriched with actor handles) |
| **Jobs (Crawl)** | `scans` table | `storage-api /scans` → `dashboard-api /jobs` |

#### Term-Fallback Actors

When a scan's seed **terms** produce **zero extractable handles** (e.g. a clean
portfolio or static page), the investigation terms themselves become distinct
actor nodes in the knowledge graph instead of the scan rendering nothing:

- Terms are sanitized by `termFallbackHandle()` (`storage-api`) — non-alphanumeric
  runs collapse to `_`, handles are truncated past 40 chars, and empty/noise terms
  are rejected.
- Fallback actors use the `TERM_FALLBACK` category and a `TERM:` handle prefix
  (e.g. `TERM:example_handle`), so they can never be confused with real
  extracted handles.
- Real extracted handles always take priority — fallback only triggers when
  **zero** handles were extracted for the whole scan.
- The graph endpoint now seeds nodes from `GET /actors` (all known actors) and
  only then adds `persona_link` edges, so standalone fallback nodes render even
  with no links.

### Seeding Test Data

To populate the dashboard with realistic test data for development/demo:

```bash
# Apply seed data (idempotent — safe to re-run)
docker exec -i onioneye_code-postgres-1 psql -U onioneye -d onioneye < OnionEye_Code/db/seed.sql
```

The seed file (`OnionEye_Code/db/seed.sql`) creates:
- 9 threat actors with realistic handles (cipher7, phantom_null, ghostwire, etc.)
- 12 identifiers (6 crypto wallets, 2 PGP keys, 4 handles)
- 5 infrastructure indicators (bulletproof hosting, hidden services)
- 9 persona links (8 PASSED negative control, 1 FAILED for filtering)
- 5 scan records (3 completed → `jobs_done` count)
- 3 content safety flags (CSAM detections)
- 2 audit log entries

All inserts use `ON CONFLICT DO NOTHING` so the seed is additive and safe to re-apply.

### Verifying Live Endpoints

```bash
# Stats (live DB counts)
curl -s http://localhost:3000/api/v1/dashboard/stats | jq .

# Graph nodes + edges
curl -s http://localhost:3000/api/v1/dashboard/graph | jq '.nodes | length'

# Leads with enriched labels
curl -s http://localhost:3000/api/v1/dashboard/leads | jq '.[0]'

# Crawl jobs from scans table
curl -s http://localhost:3000/api/v1/dashboard/jobs | jq '.[0]'
```

---

## 🛡️ Operational Security & Forensic Integrity

- **Zero Clearnet Leakage:** Docker container egress is strictly restricted at the packet filter level (`iptables` / `ip6tables`) in the Tor sidecar before worker containers spin up.
- **Zero Raw Abusive Storage:** No illicit media files or raw buffers are written to the database or cached on disk.
- **Auditable Provenance:** Every persona link records its cryptographic basis, exact extraction timestamp, and negative control audit log for chain-of-custody forensic admissibility.

---

## 🚧 Operational Boundaries & Scope Notes

- **Passive OSINT Design (No Active Exploitation):** OnionEye is engineered for passive OSINT reconnaissance and threat-actor correlation. The crawler never executes active penetration testing payloads or brute-force requests.
- **JavaScript Proof-of-Work (PoW) & DOM Engines:** Client-side JavaScript execution engines (e.g., embedded V8/Goja runtimes or headless browser automation like Playwright/Chromedp) for automated JS-based Proof-of-Work challenges are **not implemented yet**. 
- **CAPTCHA / JS-Gate Handling:** Targets requiring complex client-side JavaScript execution, anti-DDoS PoW puzzles, or interactive visual CAPTCHAs are detected by the crawler's heuristic block detector (`detectBlock()`) and gracefully routed to the **Redis `crawl:manual` intervention stream** for operator review.

---

## 📄 License & Responsible Disclosure

OnionEye is engineered strictly for authorized cybersecurity research, national defense, and law enforcement intelligence operations. Users are responsible for ensuring compliance with all local, federal, and international laws regarding dark-web reconnaissance.
