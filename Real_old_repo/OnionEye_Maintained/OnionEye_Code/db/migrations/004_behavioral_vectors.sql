-- OnionEye — Behavioral / Human-Vector Enrichment (Phase 9, §2.7 extension)
--
-- Extends the existing stylometric profile table with the temporal ("human")
-- vector family and adds a partition-friendly raw observation log so posting
-- time-of-day and response tempo can be modeled per actor.
--
-- The 384-d feature_vector stays the canonical writing-style fingerprint; the
-- new histograms are the activity signature. Both feed the correlation engine's
-- behavioral fusion and the negative-control temporal overlap check.

-- ── Extend stylometric_profiles (in place — graceful reuse) ────────────────
ALTER TABLE stylometric_profiles
    ADD COLUMN IF NOT EXISTS hour_histogram     real[24]    DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS day_histogram      real[7]     DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS median_gap_minutes int         DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS observations_count int         DEFAULT 0,
    ADD COLUMN IF NOT EXISTS provenance         TEXT        DEFAULT NULL;

-- ── Raw posting-time observations ──────────────────────────────────────────
-- Every post timestamp an author publishes on a crawled page, kept as evidence
-- so hour-of-day/weekday histograms can be recomputed or audited without
-- trusting the aggregate. source_url + provenance preserve the per-market
-- context (timezone/rhythm fusion must stay per-provenance to avoid smearing).
CREATE TABLE IF NOT EXISTS actor_observations (
    observation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id       UUID        NOT NULL REFERENCES actors(actor_id) ON DELETE CASCADE,
    observed_at    TIMESTAMPTZ NOT NULL,
    source_url     TEXT,
    provenance_tag TEXT,
    created_at     TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_actor_observations_actor_time
    ON actor_observations (actor_id, observed_at DESC);
CREATE INDEX IF NOT EXISTS idx_actor_observations_time
    ON actor_observations (observed_at);

-- ── Inter-arrival (tempo / response-time proxy) histogram ───────────────
-- We approximate response latency as inter-arrival gaps between an actor's
-- own consecutive posts observed in this crawler stream.
--
-- Hybrid binning (Option 3):
--   - 8 linear bins of width 5 minutes for 0-40 minutes
--   - 4 long bins (tail catch): 40-60m, 1-6h, 6-12h, 12-24h
-- Total = 12 bins.
ALTER TABLE stylometric_profiles
    ADD COLUMN IF NOT EXISTS tempo_histogram real[12] DEFAULT NULL;