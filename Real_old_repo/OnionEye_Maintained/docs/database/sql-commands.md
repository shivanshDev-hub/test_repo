-- ===============================================================================
-- OnionEye — FULL SCHEMA RESET + CREATE (Sub-architecture 2 & 3)
-- -------------------------------------------------------------------------------
-- Run against PostgreSQL with the pgvector extension installed. The full pipeline
-- uses the pgvector/pgvector:pg16 image in docker-compose, so `vector` is present.
--
-- USAGE:
--   cat db/sql_command.txt | psql -U $POSTGRES_USER -d $POSTGRES_DB
--   (or:  psql -U onioneye -d onioneye -f db/sql_command.txt)
--
-- NOTE: This file DROPS existing tables first so it is idempotent for demos.
-- The versioned migrations in db/migrations/ are the non-destructive equivalent.
-- ===============================================================================

-- ═══════════════════════════════════════════════════════════════════════════════
-- 0. PREP — extension + teardown (migrate down, demo reset)
-- ═══════════════════════════════════════════════════════════════════════════════
CREATE EXTENSION IF NOT EXISTS vector;

DROP TABLE IF EXISTS access_audit_logs CASCADE;
DROP TABLE IF EXISTS export_snapshots CASCADE;
DROP TABLE IF EXISTS analyst_exports CASCADE;
DROP TABLE IF EXISTS lead_reviews CASCADE;
DROP TABLE IF EXISTS persona_link_provenance CASCADE;
DROP TABLE IF EXISTS nc_static_sources CASCADE;
DROP TABLE IF EXISTS stylometric_profiles CASCADE;
DROP TABLE IF EXISTS image_hashes CASCADE;
DROP TABLE IF EXISTS extracted_links CASCADE;
DROP TABLE IF EXISTS filter_policies CASCADE;
DROP TABLE IF EXISTS proxy_leases CASCADE;
DROP TABLE IF EXISTS audit_log CASCADE;
DROP TABLE IF EXISTS content_flags CASCADE;
DROP TABLE IF EXISTS scans CASCADE;
DROP TABLE IF EXISTS persona_links CASCADE;
DROP TABLE IF EXISTS infra_indicators CASCADE;
DROP TABLE IF EXISTS identifiers CASCADE;
DROP TABLE IF EXISTS actors CASCADE;
DROP TABLE IF EXISTS audit_log_2026_09 CASCADE;

-- ═══════════════════════════════════════════════════════════════════════════════
-- 1. SUB-ARCHITECTURE 2 — Core crawler & storage schema
--    (crawler-worker / tor-sidecar / storage-api / content-shield)
-- ═══════════════════════════════════════════════════════════════════════════════

-- 1.1 Threat actors
CREATE TABLE actors (
    actor_id         UUID PRIMARY KEY,
    primary_handle   TEXT          NOT NULL,
    category         TEXT,
    last_scan_date   TIMESTAMPTZ,
    created_at       TIMESTAMPTZ   DEFAULT now()
);

-- 1.2 Identifiers found for a given threat actor
CREATE TABLE identifiers (
    identifier_id   UUID PRIMARY KEY,
    actor_id        UUID REFERENCES actors(actor_id),
    kind            TEXT,
    value           TEXT          NOT NULL,
    source          TEXT,
    provenance_tag  TEXT,
    first_seen      TIMESTAMPTZ,
    last_seen       TIMESTAMPTZ,
    CONSTRAINT uq_identifiers_actor_kind_value UNIQUE (actor_id, kind, value)
);

-- 1.3 Infrastructure indicators
CREATE TABLE infra_indicators (
    indicator_id          UUID PRIMARY KEY,
    actor_id              UUID REFERENCES actors(actor_id),
    indicator_type        TEXT,
    onion_host            TEXT,
    matched_clearnet_host TEXT,
    evidence              TEXT,
    detected_at           TIMESTAMPTZ,
    CONSTRAINT uq_infra_actor_type_host UNIQUE (actor_id, indicator_type, onion_host)
);

-- 1.4 Persona links between threat actors (correlation edges)
CREATE TABLE persona_links (
    link_id                 UUID PRIMARY KEY,
    actor_id_a              UUID REFERENCES actors(actor_id),
    actor_id_b              UUID REFERENCES actors(actor_id),
    basis                   TEXT,
    confidence              NUMERIC(3,2) NOT NULL,
    negative_control_passed BOOLEAN       NOT NULL DEFAULT false,
    created_at              TIMESTAMPTZ   DEFAULT now(),
    CONSTRAINT chk_persona_links_confidence CHECK (confidence >= 0.0 AND confidence <= 1.0),
    CONSTRAINT chk_persona_links_canonical_order CHECK (actor_id_a < actor_id_b),
    CONSTRAINT uq_persona_links_pair_basis  UNIQUE (actor_id_a, actor_id_b, basis)
);

-- 1.5 Scan records — RANGE PARTITIONED on started_at
CREATE TABLE scans (
    scan_id       UUID PRIMARY KEY,
    requested_by  TEXT,
    targets       TEXT[],
    status        TEXT,
    started_at    TIMESTAMPTZ,
    completed_at  TIMESTAMPTZ
) PARTITION BY RANGE (started_at);

-- 1.6 Content flags — RANGE PARTITIONED on flagged_at
CREATE TABLE content_flags (
    flag_id      UUID PRIMARY KEY,
    scan_id      UUID REFERENCES scans(scan_id),
    source_url   TEXT,
    category     TEXT,
    action_taken TEXT,
    flagged_at   TIMESTAMPTZ
) PARTITION BY RANGE (flagged_at);

-- 1.7 Audit log — RANGE PARTITIONED on occurred_at
CREATE TABLE audit_log (
    log_id      UUID          NOT NULL,
    actor_type  TEXT,
    action      TEXT          NOT NULL,
    resource_id UUID,
    occurred_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    metadata    JSONB
) PARTITION BY RANGE (occurred_at);

-- 1.8 Sample partition (expiry = DROP TABLE, no row-level DELETE)
CREATE TABLE audit_log_2026_09
    PARTITION OF audit_log
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');

-- 1.9 Proxy leases (Tor circuit tracking) — RANGE PARTITIONED on issued_at
CREATE TABLE proxy_leases (
    lease_id        UUID PRIMARY KEY,
    job_id          UUID,
    circuit_id      TEXT,
    issued_at       TIMESTAMPTZ,
    released_at     TIMESTAMPTZ,
    released_reason TEXT
) PARTITION BY RANGE (issued_at);

-- 1.10 Dynamic confidence view (attribution_confidence is NOT a stored column)
CREATE OR REPLACE VIEW actor_confidence AS
SELECT
    a.actor_id,
    a.primary_handle,
    COALESCE(AVG(pl.confidence), 0.0)::NUMERIC(3,2) AS attribution_confidence
FROM actors a
LEFT JOIN persona_links pl
    ON pl.actor_id_a = a.actor_id OR pl.actor_id_b = a.actor_id
GROUP BY a.actor_id, a.primary_handle;

-- ═══════════════════════════════════════════════════════════════════════════════
-- 2. SUB-ARCHITECTURE 3 — Filter, Correlation, Dashboard schema
--    (correlation-engine / dashboard-api)
-- ═══════════════════════════════════════════════════════════════════════════════

-- 2.1 Negative-control columns on persona_links (three-state evaluation)
ALTER TABLE persona_links
  ADD COLUMN nc_state VARCHAR(20) NOT NULL DEFAULT 'INDETERMINATE',
  ADD COLUMN nc_penalty_score NUMERIC(3,2) NOT NULL DEFAULT 0.00,
  ADD COLUMN nc_audit_log JSONB DEFAULT '{}'::jsonb;

CREATE INDEX idx_persona_links_nc_eval ON persona_links(nc_state) WHERE nc_state = 'INDETERMINATE';

-- 2.2 actor_id on content_flags (nullable — populated only by Phase 8A)
ALTER TABLE content_flags ADD COLUMN actor_id UUID REFERENCES actors(actor_id);
CREATE INDEX idx_content_flags_actor ON content_flags(actor_id) WHERE actor_id IS NOT NULL;

-- 2.3 Filter policies (versioned, query-time-enforced read gate)
CREATE TABLE filter_policies (
    policy_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version INT NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    min_confidence NUMERIC(3,2) NOT NULL DEFAULT 0.50,
    allowed_nc_states VARCHAR(20)[] NOT NULL DEFAULT '{"PASSED"}',
    category_thresholds JSONB NOT NULL DEFAULT '{
        "FINANCIAL": 0.60, "STYLOMETRY": 0.45, "INFRASTRUCTURE": 0.50
    }'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(128) NOT NULL
);

-- exactly one active policy version at a time
CREATE UNIQUE INDEX uq_active_policy ON filter_policies (is_active) WHERE is_active = TRUE;

-- 2.4 Active persona links view (gated through the active policy)
CREATE OR REPLACE VIEW v_active_persona_links AS
SELECT pl.link_id, pl.actor_id_a, pl.actor_id_b, pl.basis,
       pl.confidence AS attribution_confidence, pl.nc_state, pl.created_at,
       p.version AS enforced_policy_version
FROM persona_links pl
CROSS JOIN filter_policies p
WHERE p.is_active = TRUE
  AND pl.nc_state = ANY(p.allowed_nc_states)
  AND pl.confidence >= COALESCE(
      (p.category_thresholds ->> pl.basis)::numeric, p.min_confidence);

-- 2.5 Canonical actors view (recursive graph traversal: A↔B via PGP, B↔C via wallet)
CREATE OR REPLACE VIEW v_canonical_actors AS
WITH RECURSIVE link_graph AS (
    SELECT actor_id_a AS node_a, actor_id_b AS node_b
    FROM persona_links WHERE negative_control_passed = TRUE AND confidence >= 0.50
    UNION
    SELECT actor_id_b AS node_a, actor_id_a AS node_b
    FROM persona_links WHERE negative_control_passed = TRUE AND confidence >= 0.50
),
graph_traversal AS (
    SELECT actor_id AS original_actor_id, actor_id AS reached_node, ARRAY[actor_id] AS visited_path
    FROM actors
    UNION ALL
    SELECT gt.original_actor_id, lg.node_b, gt.visited_path || lg.node_b
    FROM graph_traversal gt
    JOIN link_graph lg ON gt.reached_node = lg.node_a
    WHERE NOT (lg.node_b = ANY(gt.visited_path))
),
cluster_roots AS (
    SELECT original_actor_id, MIN(reached_node) AS canonical_actor_id
    FROM graph_traversal GROUP BY original_actor_id
)
SELECT DISTINCT cr.canonical_actor_id, a.actor_id AS original_actor_id, a.primary_handle, a.created_at
FROM cluster_roots cr
JOIN actors a ON cr.original_actor_id = a.actor_id;

-- 2.6 Extracted links — table for frequency dampening (Engine 1)
CREATE TABLE extracted_links (
    link_ref_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id     UUID REFERENCES actors(actor_id),
    url          TEXT NOT NULL,
    link_type    TEXT,               -- 'internal' | 'external'
    evidence     TEXT DEFAULT 'clean',
    detected_at  TIMESTAMPTZ DEFAULT now()
);

-- 2.7 Image hashes — table for frequency dampening (reuses shield PDQ hash)
CREATE TABLE image_hashes (
    image_ref_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id       UUID REFERENCES actors(actor_id),
    perceptual_hash TEXT NOT NULL,
    evidence       TEXT DEFAULT 'clean',
    detected_at    TIMESTAMPTZ DEFAULT now()
);

-- 2.8 Stylometric profiles — pgvector HNSW for ANN candidate search
CREATE TABLE stylometric_profiles (
    actor_id UUID PRIMARY KEY REFERENCES actors(actor_id) ON DELETE CASCADE,
    sample_word_count INT NOT NULL,
    last_active_at TIMESTAMPTZ NOT NULL,
    feature_vector vector(384) NOT NULL
);

CREATE INDEX idx_stylometry_hnsw ON stylometric_profiles
USING hnsw (feature_vector vector_cosine_ops) WITH (m = 16, ef_construction = 64);

-- 2.9 Static negative-control source cache (Tor exits, VPN ranges, etc.)
CREATE TABLE nc_static_sources (
    source_id     TEXT PRIMARY KEY,
    source_url    TEXT NOT NULL,
    last_fetched  TIMESTAMPTZ,
    last_healthy  TIMESTAMPTZ,
    entry_count   INT,
    raw_payload   JSONB
);

-- 2.10 Persona link provenance (edge ↔ contributing scans junction)
CREATE TABLE persona_link_provenance (
    link_id UUID NOT NULL REFERENCES persona_links(link_id) ON DELETE CASCADE,
    scan_id UUID NOT NULL REFERENCES scans(scan_id) ON DELETE CASCADE,
    extracted_artifact_type VARCHAR(64) NOT NULL,
    associated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (link_id, scan_id, extracted_artifact_type)
);

-- 2.11 Lead reviews (analyst triage state, separate from raw evaluation)
CREATE TABLE lead_reviews (
    review_id SERIAL PRIMARY KEY,
    link_id UUID NOT NULL REFERENCES persona_links(link_id) ON DELETE CASCADE,
    analyst_username VARCHAR(64) NOT NULL DEFAULT 'analyst_demo',
    action VARCHAR(20) NOT NULL CHECK (action IN ('reviewed', 'dismissed', 'escalated')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_link_review UNIQUE (link_id)
);

-- 2.12 Analyst exports (which policy version was active at export time)
CREATE TABLE analyst_exports (
    export_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analyst_id UUID NOT NULL,
    policy_version_used INT NOT NULL REFERENCES filter_policies(version),
    exported_link_ids UUID[] NOT NULL,
    exported_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2.13 Export snapshots (frozen, SHA-256-hashed point-in-time payloads)
CREATE TABLE export_snapshots (
    export_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id VARCHAR(64) NOT NULL,
    analyst_username VARCHAR(64) NOT NULL,
    export_format VARCHAR(10) NOT NULL CHECK (export_format IN ('csv', 'json', 'pdf')),
    snapshot_hash VARCHAR(64) NOT NULL,
    snapshot_payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_export_case ON export_snapshots(case_id, created_at DESC);

-- 2.14 Access audit logs (every elevated read of flagged metadata)
CREATE TABLE access_audit_logs (
    log_id SERIAL PRIMARY KEY,
    analyst_username VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,          -- e.g. 'READ_FLAGGED_METADATA'
    target_actor_id UUID NOT NULL,
    accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ═══════════════════════════════════════════════════════════════════════════════
-- 3. SANITY CHECK — confirm every table exists
-- ═══════════════════════════════════════════════════════════════════════════════
SELECT tablename FROM pg_tables
WHERE schemaname = 'public'
ORDER BY tablename;