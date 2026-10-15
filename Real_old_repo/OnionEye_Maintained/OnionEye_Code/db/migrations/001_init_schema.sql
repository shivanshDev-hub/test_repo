-- OnionEye — Initial Schema (Subsection 2 + Subsection 3 combined)
-- Run against PostgreSQL with pgvector extension enabled.

CREATE EXTENSION IF NOT EXISTS vector;

-- ── Actors ────────────────────────────────────────────────────────────────────
CREATE TABLE actors (
    actor_id       UUID PRIMARY KEY,
    primary_handle TEXT          NOT NULL,
    category       TEXT,
    last_scan_date TIMESTAMPTZ,
    created_at     TIMESTAMPTZ   DEFAULT now(),
    CONSTRAINT uq_actors_primary_handle UNIQUE (primary_handle)
);

-- ── Identifiers ───────────────────────────────────────────────────────────────
CREATE TABLE identifiers (
    identifier_id  UUID PRIMARY KEY,
    actor_id       UUID REFERENCES actors(actor_id),
    kind           TEXT,
    value          TEXT          NOT NULL,
    source         TEXT,
    provenance_tag TEXT,
    first_seen     TIMESTAMPTZ,
    last_seen      TIMESTAMPTZ,
    CONSTRAINT uq_identifiers_actor_kind_value UNIQUE (actor_id, kind, value)
);

-- ── Infrastructure indicators ─────────────────────────────────────────────────
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

-- ── Persona links (correlation edges) ────────────────────────────────────────
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

-- ── Scan records ──────────────────────────────────────────────────────────────
CREATE TABLE scans (
    scan_id       UUID NOT NULL,
    requested_by  TEXT,
    targets       TEXT[],
    status        TEXT,
    started_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at  TIMESTAMPTZ,
    PRIMARY KEY (scan_id, started_at)
) PARTITION BY RANGE (started_at);

-- ── Content flags ─────────────────────────────────────────────────────────────
CREATE TABLE content_flags (
    flag_id      UUID NOT NULL,
    scan_id      UUID,
    source_url   TEXT,
    category     TEXT,
    action_taken TEXT,
    flagged_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (flag_id, flagged_at)
) PARTITION BY RANGE (flagged_at);

-- ── Audit log ─────────────────────────────────────────────────────────────────
CREATE TABLE audit_log (
    log_id      UUID          NOT NULL,
    actor_type  TEXT,
    action      TEXT          NOT NULL,
    resource_id UUID,
    occurred_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    metadata    JSONB
) PARTITION BY RANGE (occurred_at);

-- ── Proxy leases ──────────────────────────────────────────────────────────────
CREATE TABLE proxy_leases (
    lease_id        UUID NOT NULL,
    job_id          UUID,
    circuit_id      TEXT,
    issued_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    released_at     TIMESTAMPTZ,
    released_reason TEXT,
    PRIMARY KEY (lease_id, issued_at)
) PARTITION BY RANGE (issued_at);

-- ── Dynamic confidence view ───────────────────────────────────────────────────
CREATE OR REPLACE VIEW actor_confidence AS
SELECT
    a.actor_id,
    a.primary_handle,
    COALESCE(AVG(pl.confidence), 0.0)::NUMERIC(3,2) AS attribution_confidence
FROM actors a
LEFT JOIN persona_links pl
    ON pl.actor_id_a = a.actor_id OR pl.actor_id_b = a.actor_id
GROUP BY a.actor_id, a.primary_handle;
