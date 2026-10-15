```mermaid
flowchart TD
    subgraph T1["Team 1 - Ingestion and AI Guard - Jinesh and Rudra"]
        direction TB
        IN["Request Intake - rate limit, validate, authorize, session history"]
        GUARD["Safety Guard and Expansion - heuristic plus LLM screen, Qwen expansion, output guard"]
        QUEUE["Task Builder - idempotency lock, enqueue"]
        IN --> GUARD --> QUEUE
    end

    subgraph T2["Team 2 - Tor Crawler and Storage - Shivansh and Shubham"]
        direction TB
        CRAWL["Worker Pool - Tor circuits, zero trust, passive fetch only"]
        SHIELD{"Content Safety Shield - hash match discard"}
        PARSE["Parser and Extraction"]
        DB[("PostgreSQL")]
        CRAWL --> SHIELD
        SHIELD -->|"clean bytes"| PARSE --> DB
        SHIELD -->|"hash match, discard"| DB
    end

    subgraph T3["Team 3 - Correlation and Dashboard - Chanchal and Reshit"]
        direction TB
        ML{"ML Classifier - defense in depth"}
        CORR["Correlation Engine - negative control plus graph"]
        DASH[/"Dashboard, API, STIX Export"/]
    end

    QUEUE --> CRAWL
    DB -->|"images"| ML -->|"flags"| DB
    DB -->|"identifiers"| CORR -->|"persona links"| DB
    DB -->|"filtered views"| DASH
    ```

# Mentor Walkthrough Prep — Dark Web OSINT / Correlation Platform

**How to use this:** Read the 5-minute script out loud once or twice — it maps
1:1 onto the diagram, beat by beat. The Q&A bank underneath is longer than
you'll need; it's there so that whatever the mentor points at, you have a
ready answer instead of thinking on your feet.

---

## 30-Second Opener

> "We built a pipeline that takes an authorized analyst's search request,
> safely expands it into likely aliases, crawls dark web sources over Tor to
> collect evidence, and correlates that evidence into actor profiles — with
> safety and auditability treated as first-class requirements, not
> afterthoughts. Three sub-teams own three stages: **ingestion & AI
> guardrails**, **crawling & content safety**, and **correlation &
> dashboard**. Every stage that touches something risky — unsafe content,
> unverified attribution, or elevated access — has an explicit control around
> it, and every control is logged."

---

## The 5-Minute Script

### Beat 1 — Guarded Ingestion (~45s) · Team: Jinesh & Rudra
"A request comes in with an `authorization_ref` and seed terms. Before it
touches any model, it passes four gates: **rate limiting** (5 req/60s per
analyst, Redis sliding window), **sanitization** (strips zero-width/control
characters, regex-validates every term), **authorization** (checked against
a target-type allowlist), and a **15-minute session aggregator** that merges
this request with the analyst's recent history — that last one exists
specifically to catch someone splitting a malicious request across several
small ones."

### Beat 2 — Two-Tier Safety Guard Before the LLM Ever Runs (~30s)
"Only after those four gates does content reach AI. A sub-5ms heuristic trie
screens for known-bad terms first — cheap and fast. Anything ambiguous goes
to a second-tier LLM classifier that treats the query as **isolated data**,
wrapped in `<DATA>` tags inside a ChatML prompt, so the model can't confuse
the analyst's search terms with its own instructions. Output is
constrained to `PASS`/`BLOCK` JSON — nothing else."

### Beat 3 — Expansion, Then a Second Independent Check (~30s)
"A quantized Qwen2.5-3B model expands the term into typosquats and slang
using **grammar-constrained decoding** — it's physically incapable of
emitting anything except valid JSON matching our schema. Crypto addresses
and PGP keys pass through untouched; they're not semantic, so we don't let
the model 'creatively' mutate them. Then a second, independent guard
(Llama-Guard) re-checks every generated term and does a **partial drop** —
if 2 of 8 generated variants are bad, we keep the other 6 instead of
discarding the whole batch."

### Beat 4 — Crawling Under Tor (~45s) · Team: Shivansh & Shubham
"The job is idempotency-locked (SHA-256 hash, Redis `SET NX`) and pushed to
a Redis stream. A bounded worker pool picks it up. Every worker is
**zero-trust**: no DB credentials, read-only root filesystem, dropped Linux
capabilities, non-root, and destroyed after every job. Tor circuits are
leased by a separate manager with a **kill switch** — iptables rules applied
before the Tor daemon even starts, so a bug in our crawler physically cannot
leak traffic outside Tor. Fetching is passive only — GET/HEAD, never a
crafted payload — and CAPTCHA/blocking requires **three independent signals
to agree** before we call it a block, so we don't false-positive on a single
odd response."

### Beat 5 — Content Safety Shield, Before Anything Is Decoded (~45s)
"Raw bytes are perceptually hashed — PDQ for images, TMK+PDQF for video —
**before decoding, rendering, or storage.** On a hash match to known CSAM,
bytes are zeroed immediately and only a boolean flag continues downstream;
the pipeline stops there. This is the primary control. A second-layer ML
classifier (Team 3) catches material that's never been hashed before —
that's defense-in-depth, not a replacement."

### Beat 6 — Correlation With a Built-In Skeptic (~60s) · Team: Chanchal & Reshit
"Extracted identifiers — wallets, PGP keys, infra, stylometry — feed a
correlation engine that builds a property graph. But before any two actors
get linked, a **Negative Control engine** checks the match against known
confounders: shared hosting, disposable emails, Tor exit-node coincidences,
frequency dampening for identifiers reused across many unrelated actors. A
raw 0.85 wallet-reuse match can get downweighted to 0.42 once we see the
wallet is shared across three actors — that's the difference between 'same
person' and 'both used the same exchange.' Every accepted link carries a
**confidence score plus a human-readable rationale** — there's no code path
that emits a bare, unexplained number."

### Beat 7 — Dashboard, Audit, Export (~30s)
"Analysts triage leads as reviewed/dismissed/escalated. Viewing flagged
content requires elevated, **logged** access. Filter thresholds are
versioned so tightening a policy doesn't require reprocessing old data.
Exports are **frozen, hashed snapshots** — not live re-queries — so two
exports of the same case are provably identical or provably different.
We also export to **STIX 2.1**, the standard threat-intel format, so this
interoperates with other CTI tooling instead of being a dead end."

### Close (~15s)
"So: four independent gates before AI runs, two independent AI guards
around the one AI step that's user-facing, a hash-based discard before any
byte is decoded, a skeptic engine that actively argues against false
attribution, and an audit trail on every elevated read. Nothing in this
system trusts a single signal."

---

## Anticipated Questions

### Safety / CSAM handling
**Q: Why hash-match before decoding, instead of just running everything
through a classifier?**
Hashing operates on raw bytes and never requires the content to be rendered
or opened — it's cheaper, deterministic, and means the worst-category
content is never viewed by any part of our system, human or model. The
classifier is the second layer for content that has never been hashed
before, which hash-matching structurally cannot catch.

**Q: How does the ML classifier avoid needing a real CSAM dataset to train
on?**
It never targets the class directly. It's a composite of three
independently-trained, adults-only proxy vectors: an age model
(UTKFace/IMDB-WIKI/FairFace/AgeDB), a nudity model (NudeNet/open_nsfw
methodology), and a pose/behavior model (COCO Keypoints/MPII). A hard
age-gate (`p_minor < 0.40`) runs first and short-circuits most images to
"clean" before the more expensive models ever run — this also keeps
inference cheap since the overwhelming majority of images are adults.

**Q: What happens to a borderline score?**
Never auto-discarded. `[0.45, 0.65)` goes to a human review queue; only
`≥0.65` is treated as discard-and-flag, and `<0.65` never displays the
underlying image regardless — the raw bytes are gone either way once the
Phase 5 shield or 8A discard path fires.

**Q: Isn't storing `actor_id` on a flagged-content row risky?**
It's nullable and populated only via a best-effort lookup against handles
already extracted from the same page — if that's ambiguous (zero or
multiple candidates), it stays `NULL` rather than guessing. And the
underlying image is never stored regardless of whether the linkage
succeeds — only the fact "this actor is associated with a flagged event"
propagates.

### Correlation accuracy / false positives
**Q: How do you stop two unrelated actors who just share a hosting
provider from being scored as the same person?**
That's the Negative Control engine's entire job — every candidate link is
checked against static denylists (known shared infra, disposable domains,
tutorial/exchange wallets) and frequency-based dampening (an identifier
shared across many actors loses value as a unique signal). A raw 0.85 match
can end up at 0.42 after dampening.

**Q: What if a data source you need for negative control is down?**
Fail-open, not fail-closed: missing data defaults to zero penalty (not an
assumed confounder), the edge is tagged `INDETERMINATE` rather than
silently `PASSED`, and it's auto-re-queued once the source comes back
healthy. We explicitly didn't want "we couldn't check" to look identical to
"we checked and it's clean."

**Q: Are the weights (0.40 for hard identifiers, 0.15 for soft) arbitrary?**
Currently yes — they're an engineering starting point, not statistically
derived. The defensible next step (raised separately) is estimating them
empirically via a Fellegi-Sunter-style model once enough analyst
review outcomes (`lead_reviews`) accumulate, rather than hand-tuning them
forever.

**Q: Why avoid a single scan_id on a correlation edge?**
Because an edge represents accumulated evidence across time — a link built
from evidence in scan #100 and scan #500 can't honestly carry one scan ID.
Provenance is tracked in a separate junction table instead
(`persona_link_provenance`), one row per contributing scan.

**Q: How do you keep stylometry comparison from becoming O(n²) across every
actor pair?**
A three-gate blocking pipeline before any heavy comparison: minimum 500
aggregate words, a 365-day temporal-overlap window, then an HNSW
approximate-nearest-neighbor search (pgvector) capped at the top 50
candidates. Only survivors reach the full correlation engine.

### AI / prompt-injection safety
**Q: How do you stop the search terms themselves from being a prompt
injection against your expansion model?**
Input is wrapped in `<DATA>` tags inside a ChatML prompt specifically to
mark it as passive data the model must not treat as instructions, and two
independent guards run before generation (heuristic trie, then an LLM
classifier) and one after (Llama-Guard on the output). Output decoding is
grammar-constrained (GBNF/JSON schema), so even a successful injection
can't make the model emit anything outside the expected JSON shape.

**Q: Why grammar-constrained decoding specifically?**
It removes an entire failure class — markdown wrappers, preambles,
malformed JSON — by restricting token sampling to only valid schema tokens,
rather than post-hoc parsing and hoping the format is clean.

**Q: Why keep a 15-minute session history at all — isn't per-request
screening enough?**
Per-request screening can't see a request that's been deliberately split
into several innocuous-looking pieces ("salami slicing"). Merging the
current query with the last 15 minutes of that analyst's activity before
classification closes that gap.

### Crawling / operational security
**Q: How do you know the crawler can't leak traffic outside Tor?**
A container-level iptables kill switch (DOCKER-USER chain, IPv4 and IPv6)
is applied in the entrypoint script *before* the Tor daemon or the Go
worker binary ever starts, and the rules file lives outside Docker's own
management path so a Docker-level operation can't accidentally flush it.
Startup order is enforced: firewall → Tor → SOCKS health check → worker.

**Q: How do you avoid false-positive CAPTCHA/block detection triggering
unnecessary human intervention?**
Three independent signals have to agree: an explicit 403/408/503 status,
a response-body-length deviation from that target's historical Redis
baseline, and a second, independent Tor circuit confirming the same block.
A single odd response never trips it.

**Q: Why "passive fetch only"?**
It's an explicit ethical/legal boundary — plain GET/HEAD against whatever
the server already serves, never a crafted or exploitative payload. This
keeps the crawler firmly on the "collect what's public" side of the line.

**Q: Why zero-trust worker containers?**
Containers get destroyed after every job, hold no DB credentials (writes go
through a storage-service API only), run as non-root with all Linux
capabilities dropped, and mount the root filesystem read-only. A compromised
worker doesn't equal database access.

### Data governance / dashboard
**Q: Why keyset pagination instead of offset?**
The underlying result set updates live as new correlation links land.
Offset pagination can silently skip or repeat rows mid-scroll; keyset
pagination (sorted on confidence + a link_id tiebreaker) can't.

**Q: Why export snapshots instead of querying live at export time?**
Investigative/audit integrity — the same export run today vs. tomorrow
must not silently differ. We freeze a hashed (SHA-256) JSON payload at
export time and every renderer (CSV/JSON/PDF/STIX) reads only from that
frozen payload, so two analysts can verify two "same" exports are
byte-identical.

**Q: Why version the filter policy instead of just changing a threshold in
code?**
So a threshold change takes effect on the next read with zero backfill,
and every past export records which policy version was active — "why did
analyst X see this lead under the old rules" is always answerable by
re-running the view pinned to that version.

**Q: Why STIX 2.1 for export?**
It's the standard structured threat-intel format (OASIS), which means this
tool's output can plug into other CTI tooling instead of being a dead end
that only our own dashboard can read.

**Q: What happens to a scan that extracts zero actor handles — does the graph
stay empty?**
No. Investigation seed terms are carried through the crawl pipeline on the
job payload, and when a scan completes with zero extracted handles, the
storage layer creates one distinct actor per term using the `TERM_FALLBACK`
category and a `TERM:`-prefixed handle (e.g. `TERM:ShivanshSahajpal`). This
keeps the knowledge graph populated even for static/clean targets, and the
`TERM:` prefix guarantees fallback anchors can never be mistaken for real
extracted handles. Real handles always take priority — fallback fires only
when zero handles were found. The graph endpoint reads all actors (via
`GET /actors`) before adding persona-link edges, so these standalone nodes
render even with no links yet.

### System robustness
**Q: What happens if the ML classifier (8A) or the read-gate (8B) goes
down?**
Neither sits in the critical path. 8A only ever *adds* rows — disable it
and Phase 9/10 still function on whatever the hash shield already caught,
you just lose the second-layer signal. 8B only ever *filters a view* —
disable it and the dashboard reads the unfiltered view instead of failing.

**Q: How do you avoid double-processing the same request on a
retry/resubmit?**
An idempotency hash (`SHA256(analyst_id + target + host + minute_bucket)`)
is checked via Redis `SET NX EX 60` before a job is ever enqueued.

---

## If They Push on Limitations (answer honestly, don't oversell)

- **RBAC is a demo stand-in.** Elevated access is currently a hardcoded
  role flag toggled in the UI, not a real JWT/claims-based system —
  explicitly flagged in the docs as future work, not production-ready.
- **Row-level security is deferred.** Enforcement is currently
  application-layer, not DB-layer RLS.
- **Network micro-segmentation beyond Docker policy is a future feature.**
- **CSAM hash set is a synthetic stub for the demo** — production would
  point at a recognized hash-sharing program, not our own list.
- **Correlation weights are heuristic, not empirically fit yet** — the
  natural next step once enough analyst-reviewed outcomes exist.

Framing these as "known, explicitly scoped-out for the demo" rather than
"things we missed" is the honest and stronger answer.

---

## Quick-Reference Numbers (for rapid-fire recall)

| Control | Value |
|---|---|
| Rate limit | 5 requests / 60s / analyst |
| Session correlation window | 15 min (900s) |
| Fast heuristic screen latency | <5ms |
| Seed term count | 1–10 |
| max_depth | clamped 1–5, default 3 |
| Idempotency lock TTL | 60s |
| Politeness delay | 2–8s randomized, per host |
| Retry budget | 3 attempts max |
| CSAM discard threshold | classifier score ≥ 0.65 |
| Human-review band | 0.45 ≤ score < 0.65 |
| Age hard-gate short circuit | p_minor < 0.40 |
| Correlation accept threshold | attribution_confidence ≥ 0.50 |
| Hard identifier weight | ~0.40 |
| Soft/behavioral weight | ~0.15 |
| Stylometry min sample (extraction) | 150 words |
| Stylometry min sample (correlation) | 500 words |
| Stylometry temporal window | 365 days |
| Circuit idle rotation | 2 min |
| Overall request timeout | 90s |
| Max response body | 10MB |
