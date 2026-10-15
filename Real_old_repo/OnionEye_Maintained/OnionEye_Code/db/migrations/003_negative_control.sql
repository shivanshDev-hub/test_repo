-- OnionEye — Negative Control Console (LeftControlRail NEGATIVE CONTROL tab)
-- Depends on: 001_init_schema.sql, 002_add_nc_columns.sql
--
-- Two tables back the analyst-facing negative-control panel:
--   * negative_control_rules      — per (mode × artifact_type) damping policy.
--   * negative_control_artifacts  — concrete artifact values + occurrence counts.
--
-- Modes:
--   FREQUENCY — artifacts whose occurrence_count exceeds frequency_cap are
--               dampened toward zero (logarithmic frequency dampening).
--   STATIC    — a fixed denylist (frequency_cap ignored); values are always
--               suppressed regardless of how many actors reference them.

-- ── Negative control rules ────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS negative_control_rules (
    rule_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mode          VARCHAR(20) NOT NULL DEFAULT 'FREQUENCY'
                  CHECK (mode IN ('FREQUENCY', 'STATIC')),
    artifact_type VARCHAR(40) NOT NULL,
    frequency_cap INT NOT NULL DEFAULT 0,          -- 0 => dampening disabled
    enabled       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- One active policy per (mode, artifact_type) — upserts replace in place.
CREATE UNIQUE INDEX IF NOT EXISTS uq_nc_rule
    ON negative_control_rules (mode, artifact_type);

-- ── Negative control artifacts ────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS negative_control_artifacts (
    artifact_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id         UUID REFERENCES negative_control_rules(rule_id)
                    ON DELETE CASCADE,
    artifact_type   VARCHAR(40) NOT NULL,
    value           TEXT NOT NULL,
    occurrence_count INT NOT NULL DEFAULT 1,
    first_seen_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_nc_artifact UNIQUE (rule_id, artifact_type, value)
);

CREATE INDEX IF NOT EXISTS idx_nca_type
    ON negative_control_artifacts (artifact_type);

CREATE INDEX IF NOT EXISTS idx_nca_type_occ
    ON negative_control_artifacts (artifact_type, occurrence_count DESC);

-- ── Seed baseline rules so the panel shows real state on first boot ───────────
INSERT INTO negative_control_rules (mode, artifact_type, frequency_cap, enabled)
VALUES
    ('FREQUENCY', 'onion_link', 3, TRUE),
    ('FREQUENCY', 'email',      5, TRUE),
    ('FREQUENCY', 'btc_address', 2, TRUE),
    ('STATIC',    'onion_link', 0, TRUE)
ON CONFLICT (mode, artifact_type) DO NOTHING;