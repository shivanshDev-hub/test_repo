# 📖 OnionEye — Complete Tool Usage Manual

**OnionEye** is a specialized dark-web reconnaissance and threat actor de-anonymization intelligence platform. It automates hidden service traversal, passive entity extraction, multi-layered content safety validation, persona graph correlation, and STIX 2.1 threat intelligence export.

---

## 📑 Table of Contents

1. [System Prerequisites](#1-system-prerequisites)
2. [Architecture & Service Topology](#2-architecture--service-topology)
3. [Quick Start with Docker Compose](#3-quick-start-with-docker-compose)
   - [Standard Mode (Lightweight / Offline AI Fallback)](#standard-mode-recommended-for-quick-start)
   - [AI-Enhanced Mode (With Ollama LLM Engine)](#ai-enhanced-mode-full-nlp-expansion)
4. [Local Development Setup](#4-local-development-setup)
   - [Running Frontend Dev Server](#running-the-frontend)
   - [Running Backend Services Independently](#running-backend-services-independently)
5. [Using the Tactical Web Console (UI)](#5-using-the-tactical-web-console-ui)
6. [API Reference & CLI Examples](#6-api-reference--cli-examples)
   - [Dashboard API — Live Data Endpoints](#6-dashboard-api--live-data-endpoints)
   - [Seeding Test Data](#7-seeding-test-data)
7. [Configuration & Target Allowlist](#8-configuration--target-allowlist)
8. [Fail-Open & Offline Fault Tolerance](#9-fail-open--offline-fault-tolerance)
9. [Troubleshooting & FAQs](#10-troubleshooting--faqs)

---

## 1. System Prerequisites

Before running OnionEye, ensure the following tools are installed on your host system:

- **Docker Engine:** `v24.0+` with `docker compose` plugin (`v2.20+`)
- **Go:** `1.22+` (if building/running Go microservices locally)
- **Node.js & npm:** `Node v18+` / `npm v9+` (for the frontend React UI)
- **Git:** for version control and workspace management
- *(Optional)* **Ollama:** if running local LLMs without Docker (model: `qwen2.5:3b`)

---

## 2. Architecture & Service Topology

The system comprises decoupled, zero-trust microservices:

| Service | Port | Description |
|---|---|---|
| **`nginx`** | `80` | Unified reverse proxy routing Dashboard UI, Ingestion, and Dashboard APIs |
| **`frontend`** | `80` (dev `5173`) | Tactical Threat Intel SOC UI (2D Canvas HUD Graph, STIX Exporter) |
| **`ingestion-api`** | `8080` | 10-phase input safety validation, sliding-window rate limit, and stream dispatch |
| **`crawler-worker`** | Internal | Tor-routed passive scraper with anti-obfuscation normalizer & link extractor |
| **`content-shield`** | `8082` | Stream-level byte inspection & instant memory zeroing (`phash.Zero`) |
| **`correlation-engine`**| `8083` | Multi-modal hard scoring, 384-d stylometry, time decay, & Negative Control dampener |
| **`storage-api`** | `8081` | Zero-trust PostgreSQL read/write boundary & partition manager |
| **`dashboard-api`** | `3000` | Telemetry, materialized views, actor stats, lead reviews, and STIX store |
| **`tor-sidecar`** | `9050` / `9051` | Hardened Tor SOCKS5 proxy sidecar with iptables killswitch & circuit renewal |
| **`redis`** | `6379` | In-memory message stream (`crawl:jobs`), rate-limit state, & session memory |
| **`postgres`** | `5432` | Relational audit store, HNSW stylometric vector index, & time-partitioned telemetry |
| **`ollama`** *(opt)* | `11434` | Local LLM inference engine for query expansion (`qwen2.5:3b`) & guarding (`llama-guard3:1b`) |

---

## 3. Quick Start with Docker Compose

All stack definitions reside under `OnionEye_Code/docker-compose.yml`.

### Standard Mode (Recommended for Quick Start)
Starts all core crawlers, safety shields, databases, ingestion pipeline, dashboard, and the web UI:

```bash
cd OnionEye_Code

# Build and bring up the entire stack (excluding heavy LLM weights)
docker compose up -d --build
```

- **Frontend SOC Console:** Access at `http://localhost:5173`
- **Ingestion API:** Reachable at `http://localhost:8080`
- **Dashboard API:** Reachable at `http://localhost:3000`

### AI-Enhanced Mode (Full NLP Expansion)
To include the Ollama local AI model container (`qwen2.5:3b`) for semantic seed generation:

```bash
cd OnionEye_Code

# Launch stack including the 'ai' profile container
docker compose --profile ai up -d --build
```

To stop all services:
```bash
docker compose down
```

---

## 4. Local Development Setup

### Running the Frontend
```bash
cd OnionEye_Code/frontend

# Install dependencies & start Vite live-reload dev server
npm install
npm run dev
```

### Running Backend Services Independently

1. **Ingestion API (`:8080`):**
   ```bash
   cd OnionEye_Code/services/ingestion-api
   go run ./cmd/server
   ```
2. **Storage API (`:8081`):**
   ```bash
   cd OnionEye_Code/services/storage-api
   go run ./cmd/server
   ```
3. **Content Shield (`:8082`):**
   ```bash
   cd OnionEye_Code/services/content-shield
   go run ./cmd/shield
   ```
4. **Correlation Engine (`:8083`):**
   ```bash
   cd OnionEye_Code/services/correlation-engine
   go run ./cmd/engine
   ```
5. **Dashboard API (`:3000`):**
   ```bash
   cd OnionEye_Code/services/dashboard-api
   go run ./cmd/server
   ```
6. **Crawler Worker:**
   ```bash
   cd OnionEye_Code/services/crawler-worker
   go run ./cmd/worker
   ```


---

## 5. Using the Tactical Web Console (UI)

### Dispatched Target Investigation
1. Open `http://localhost:5173` in your browser.
2. Click **"NEW INVESTIGATION"** (`[+]`) in the top navigation bar.
3. Select your **Target Type**:
   - `onion_address` — Tor v3 Hidden Service domain (e.g., `mock-market-1.onion`)
   - `vendor_handle` — Darknet vendor pseudonym (e.g., `shadow_vendor`)
   - `btc_wallet` — Cryptocurrency address (e.g., `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`)
   - `pgp_key` — PGP public key fingerprint or key ID
4. Fill in the **Target Value** and **Target Host** (e.g., `mock-market-1.onion`).
5. Enter comma-separated **Seed Terms** (e.g., `shadow_vendor, vault, escrow`).
6. Set the **Max Depth** (1 to 5) and click **"DISPATCH INVESTIGATION"**.

### Canvas Vector HUD Graph Navigation
- **Entity Glyphs:**
  - 👤 `Threat Actor` — Identified marketplace or forum persona
  - 🔑 `PGP Key` — Cryptographic signature / key fingerprint
  - ₿ `Crypto Wallet` — Bitcoin / Monero addresses
  - 🧅 `Hidden Service` — Tor `.onion` domain node
  - 🖥 `Clearnet Infra` — Discovered server IP / clearnet indicator
  - ✉ `Secure Drop` — Contact email / drop handle
- **Interactivity:**
  - **Click a Node:** Locks the targeting reticle `[  ]` on target, auto-centers the viewport, and opens the Entity Detail Drawer.
  - **Drag / Pan:** Click and drag canvas background to pan around the topology.
  - **Crimson Halos:** Pulsating crimson beacons highlight entities with content flags.

### Lead Attribution & Negative Control Inspection
Clicking on an edge or actor in the graph loads the attribution breakdown:
- **Confidence Rating:** Composite Bayesian/statistical score (0% to 100%).
- **Forensic Basis:** Justifications (`pgp_colocation`, `stylometry_hnsw`, `temporal_cooccurrence`).
- **Negative Control Status:**
  - `NC: PASS` — Clean correlation; passes false-positive suppression.
  - `NC: FLAGGED_REVIEW` — Penalized for co-locating on shared infrastructure (Tor exit relays, shared payment pools).

### Exporting Cases (STIX 2.1 / JSON / CSV)
1. In the top-right toolbar on the HUD Console:
2. Click **"EXPORT STIX 2.1"** to generate an OASIS-compliant CTI bundle (`threat-actor`, `identity`, `indicator`, `relationship` SDOs) for MISP or OpenCTI.
3. Click **"EXPORT RAW"** for full JSON/CSV forensic snapshot archives.

---

## 6. API Reference & CLI Examples

### 1. Dispatch an Investigation Scan (Ingestion API)
**Endpoint:** `POST http://localhost:8080/api/v1/investigations` (or via Nginx: `http://localhost/api/v1/investigations`)

```bash
curl -X POST http://localhost:8080/api/v1/investigations \
  -H "Content-Type: application/json" \
  -d '{
    "analyst_id": "analyst_042",
    "authorization_ref": "AUTH-NTRO-2026-8841",
    "target_type": "onion_address",
    "target_value": "mock-market-1.onion",
    "target_host": "mock-market-1.onion",
    "seed_terms": ["shadow_vendor", "0x71C3a9", "escrow"],
    "max_depth": 3
  }'
```

**Response (HTTP 202 Accepted):**
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "queued",
  "verified_terms": ["shadow_vendor", "0x71c3a9", "escrow"],
  "guard_verdict": "PASS",
  "llm_unavailable": false,
  "enqueued_at": "2026-09-06T18:40:00Z"
}
```

> **Term-fallback behavior:** If a crawl completes and **zero** extractable
> handles are found (e.g. a clean static portfolio site), the `seed_terms` are
> not lost — each one is persisted as a distinct **`TERM_FALLBACK`** actor with
> a `TERM:`-prefixed handle (e.g. `TERM:shadow_vendor`). This guarantees every
> investigation scan populates the knowledge graph, and the `TERM:` prefix
> keeps fallback actors visually distinct from real extracted handles. Real
> extracted handles always take priority over fallback actors.

### 2. Query Telemetry, Leads & STIX 2.1 (Dashboard API)
**Endpoints on Port `3000` (or via Nginx `http://localhost/api/v1/dashboard/*`):**

```bash
# Fetch aggregated platform stats
curl -s http://localhost:3000/api/v1/dashboard/stats | jq .

# Fetch interactive graph topology
curl -s http://localhost:3000/api/v1/dashboard/graph | jq .

# Fetch candidate leads for review
curl -s http://localhost:3000/api/v1/dashboard/leads | jq .

# Export OASIS STIX 2.1 Threat Intel bundle
curl -s http://localhost:3000/api/v1/dashboard/stix2 | jq .
```

### 3. Evaluate Persona Attribution & Stylometry (Correlation Engine)
**Endpoints on Port `8083`:**

```bash
# Evaluate actor composite attribution score with Negative Control dampening
curl -X POST http://localhost:8083/evaluate/actor \
  -H "Content-Type: application/json" \
  -d '{
    "actor_a": "shadow_vendor",
    "actor_b": "5h4d0w_v3nd0r",
    "signals": [
      { "type": "pgp_match", "strength": 0.95, "weight": 0.40, "age_days": 10 },
      { "type": "btc_wallet_reuse", "strength": 0.85, "weight": 0.35, "age_days": 5 }
    ],
    "nc_passed_count": 3
  }'

# Evaluate stylometric writing similarity between two 384-d dense vector embeddings
curl -X POST http://localhost:8083/evaluate/stylometry \
  -H "Content-Type: application/json" \
  -d '{
    "vector_a": [0.02, 0.87, 0.15, 0.42],
    "vector_b": [0.03, 0.85, 0.14, 0.40],
    "sample_words_a": 320,
    "sample_words_b": 280
  }'
```

### 4. Content Safety Stream Check (Content Shield)
**Endpoint on Port `8082`:**

```bash
# Send raw payload for zero-decode hash screening
curl -X POST http://localhost:8082/shield \
  -H "Content-Type: application/octet-stream" \
  -H "X-Job-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Source-URL: http://mock-market-1.onion/image.jpg" \
  --data-binary "CLEAN_PAYLOAD_BYTES"
```

### 5. Service Health Checks
```bash
curl -s http://localhost:8080/healthz | jq .
curl -s http://localhost:8081/health | jq .
curl -s http://localhost:8082/health | jq .
curl -s http://localhost:8083/healthz | jq .
curl -s http://localhost:3000/api/v1/dashboard/health | jq .
```

---

### 6. Dashboard API — Live Data Endpoints

All dashboard panels now fetch live data from PostgreSQL via `storage-api`. The `dashboard-api` proxies and enriches responses for the React frontend.

**Data flow:**
```
React Frontend → nginx → dashboard-api (enrichment/proxy) → storage-api → PostgreSQL
```

| Endpoint | Description |
|---|---|
| `GET /api/v1/dashboard/stats` | Aggregated actor count, link count, identifier stats, scan completion count, flag count |
| `GET /api/v1/dashboard/graph` | Graph nodes (all actors from `/actors`, including standalone `TERM_FALLBACK` nodes, plus identifiers) and edges (persona links) for RADAR canvas |
| `GET /api/v1/dashboard/leads` | Persona link leads enriched with `label_a`/`label_b` (actor `primary_handle` with fallbacks) |
| `GET /api/v1/dashboard/jobs` | Crawl scan history from the `scans` table (job_id, status, targets, timestamps) |
| `DELETE /api/v1/dashboard/delete/node` | Delete a single actor (and all child rows: identifiers, infra, persona_links) |
| `DELETE /api/v1/dashboard/delete/all` | Truncate all dashboard data (persona_links → identifiers → infra → actors) |

```bash
# Live stats
curl -s http://localhost:3000/api/v1/dashboard/stats | jq .

# Live graph (nodes + edges)
curl -s http://localhost:3000/api/v1/dashboard/graph | jq '.nodes | length'

# Live leads with enriched labels
curl -s http://localhost:3000/api/v1/dashboard/leads | jq '.[0]'

# Live crawl jobs
curl -s http://localhost:3000/api/v1/dashboard/jobs | jq '.[0]'
```

### 7. Seeding Test Data

To populate the dashboard with realistic test data:

```bash
docker exec -i onioneye_code-postgres-1 psql -U onioneye -d onioneye < db/seed.sql
```

This creates 9 actors, 12 identifiers, 5 infra indicators, 9 persona links, 5 scans, 3 content flags, and 2 audit entries — all idempotent (safe to re-run).

---

## 8. Configuration & Target Allowlist

Authorized target types and operational parameters are managed in configuration files:

- **Target Allowlist (`OnionEye_Code/config/target_allowlist.yaml`):**
  Defines permitted target categories for Phase 4 validation:
  ```yaml
  allowed_target_types:
    - vendor_handle
    - onion_address
    - btc_wallet
    - pgp_key
    - forum_user
    - email_address
    - jabber_id
    - infra_ip
  ```

- **Market Seeds (`OnionEye_Code/db/seed/markets.yaml`):**
  Maps predefined hidden services and simulated marketplace topologies for offline test runs.

---

## 9. Fail-Open & Offline Fault Tolerance

OnionEye is architected with graceful fallback pathways:

1. **Ollama / AI Offline:**
   - The Tier-1 regular-expression Content Guard remains fully active.
   - The Tier-2 LLM guard automatically outputs `PASS_NO_LLM`.
   - Seed expansion uses unexpanded original terms without throwing 500 errors.
2. **Redis Offline:**
   - Ingestion API generates deterministic SHA-256 task identifiers and returns HTTP 202 Accepted without stalling client workflows.
3. **Dashboard API Offline:**
   - The Vite frontend reverse-proxy traps `ECONNREFUSED` / `503` errors and delivers rich forensic mock telemetry to keep the UI interactive.

---

## 10. Troubleshooting & FAQs

### Q: Why do I see `(regex-only, LLM offline)` when creating a scan?
**A:** This indicates that the Ollama container or local LLM service is not running. The scan was safely validated using Tier-1 keyword filters and dispatched to the crawl queue. If you want AI term expansion, run Docker Compose with `--profile ai`.

### Q: Can OnionEye scan non-existent or offline `.onion` URLs?
**A:** Yes. The crawler routes all connections through Tor SOCKS5. If a host is down or nonexistent, Tor returns connection refused or timeout; the worker rotates circuits 3 times before safely marking the frontier exhausted and closing the job.

### Q: How do I export data to my SIEM?
**A:** Open the OnionEye console at `http://localhost:5173`, select your target investigation in the graph, and click **"EXPORT STIX 2.1"** in the top navigation bar.

---
*For in-depth architectural details, consult `../architecture/decided/`.*

