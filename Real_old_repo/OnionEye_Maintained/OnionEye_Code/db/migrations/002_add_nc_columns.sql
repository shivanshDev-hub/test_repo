-- OnionEye — Phase 8/9/10 columns and tables (Subsection 3 additions)
-- Depends on: 001_init_schema.sql

-- ── Negative control columns on persona_links ─────────────────────────────────
ALTER TABLE persona_links
  ADD COLUMN nc_state VARCHAR(20) NOT NULL DEFAULT 'INDETERMINATE',
  ADD COLUMN nc_penalty_score NUMERIC(3,2) NOT NULL DEFAULT 0.00,
  ADD COLUMN nc_audit_log JSONB DEFAULT '{}'::jsonb;

CREATE INDEX idx_persona_links_nc_eval ON persona_links(nc_state) WHERE nc_state = 'INDETERMINATE';

-- ── actor_id on content_flags (nullable, populated by Phase 8A) ───────────────
-- NOTE: scan_id FK was dropped because scans is partitioned by started_at and
-- PostgreSQL cannot FK a partitioned table to another partitioned table unless
-- the referencing columns include the partition key. Referential integrity is
-- enforced at the application layer (storage-api).
ALTER TABLE content_flags ADD COLUMN actor_id UUID REFERENCES actors(actor_id);
CREATE INDEX idx_content_flags_actor ON content_flags(actor_id) WHERE actor_id IS NOT NULL;

-- ── Filter policies (versioned read-time thresholds) ──────────────────────────
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

CREATE UNIQUE INDEX uq_active_policy ON filter_policies (is_active) WHERE is_active = TRUE;

INSERT INTO filter_policies (version, is_active, min_confidence, allowed_nc_states, category_thresholds, created_by)
VALUES (1, TRUE, 0.50, '{"PASSED"}', '{"FINANCIAL": 0.60, "STYLOMETRY": 0.45, "INFRASTRUCTURE": 0.50}'::jsonb, 'system')
ON CONFLICT (version) DO NOTHING;


-- ── Active persona links view (gated by filter policy) ────────────────────────
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

-- ── Canonical actors view (recursive graph traversal) ─────────────────────────
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
    SELECT original_actor_id, MIN(reached_node::text)::uuid AS canonical_actor_id
    FROM graph_traversal GROUP BY original_actor_id
)
SELECT DISTINCT cr.canonical_actor_id, a.actor_id AS original_actor_id, a.primary_handle, a.created_at
FROM cluster_roots cr
JOIN actors a ON cr.original_actor_id = a.actor_id;

-- ── Extracted links (for frequency dampening) ─────────────────────────────────
CREATE TABLE extracted_links (
    link_ref_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id     UUID REFERENCES actors(actor_id),
    url          TEXT NOT NULL,
    link_type    TEXT,
    evidence     TEXT DEFAULT 'clean',
    detected_at  TIMESTAMPTZ DEFAULT now()
);

-- ── Image hashes (for frequency dampening) ────────────────────────────────────
CREATE TABLE image_hashes (
    image_ref_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id       UUID REFERENCES actors(actor_id),
    perceptual_hash TEXT NOT NULL,
    evidence       TEXT DEFAULT 'clean',
    detected_at    TIMESTAMPTZ DEFAULT now()
);

-- ── Stylometric profiles (pgvector HNSW) ──────────────────────────────────────
CREATE TABLE stylometric_profiles (
    actor_id UUID PRIMARY KEY REFERENCES actors(actor_id) ON DELETE CASCADE,
    sample_word_count INT NOT NULL,
    last_active_at TIMESTAMPTZ NOT NULL,
    feature_vector vector(384) NOT NULL
);

CREATE INDEX idx_stylometry_hnsw ON stylometric_profiles
USING hnsw (feature_vector vector_cosine_ops) WITH (m = 16, ef_construction = 64);

-- ── Static source cache (Tor exits, VPN ranges, etc.) ─────────────────────────
CREATE TABLE nc_static_sources (
    source_id     TEXT PRIMARY KEY,
    source_url    TEXT NOT NULL,
    last_fetched  TIMESTAMPTZ,
    last_healthy  TIMESTAMPTZ,
    entry_count   INT,
    raw_payload   JSONB
);

-- ── Persona link provenance ───────────────────────────────────────────────────
-- NOTE: scan_id FK dropped for the same partitioned-PK reason as content_flags.
CREATE TABLE persona_link_provenance (
    link_id UUID NOT NULL REFERENCES persona_links(link_id) ON DELETE CASCADE,
    scan_id UUID NOT NULL,
    extracted_artifact_type VARCHAR(64) NOT NULL,
    associated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (link_id, scan_id, extracted_artifact_type)
);

-- ── Lead reviews ──────────────────────────────────────────────────────────────
CREATE TABLE lead_reviews (
    review_id SERIAL PRIMARY KEY,
    link_id UUID NOT NULL REFERENCES persona_links(link_id) ON DELETE CASCADE,
    analyst_username VARCHAR(64) NOT NULL DEFAULT 'analyst_demo',
    action VARCHAR(20) NOT NULL CHECK (action IN ('reviewed', 'dismissed', 'escalated')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_link_review UNIQUE (link_id)
);

-- ── Analyst exports ───────────────────────────────────────────────────────────
CREATE TABLE analyst_exports (
    export_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analyst_id UUID NOT NULL,
    policy_version_used INT NOT NULL REFERENCES filter_policies(version),
    exported_link_ids UUID[] NOT NULL,
    exported_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ── Export snapshots ──────────────────────────────────────────────────────────
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

-- ── Access audit logs ─────────────────────────────────────────────────────────
CREATE TABLE access_audit_logs (
    log_id SERIAL PRIMARY KEY,
    analyst_username VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    target_actor_id UUID NOT NULL,
    accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
