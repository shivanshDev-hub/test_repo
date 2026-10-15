-- OnionEye — Initial partitions for time-partitioned tables
-- Depends on: 001_init_schema.sql, 002_add_nc_columns.sql
-- The architecture uses monthly range partitions; dropping a partition
-- expires data without row-level DELETE.

-- ── scans (partitioned by started_at) ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS scans_2026_09
    PARTITION OF scans
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');

CREATE TABLE IF NOT EXISTS scans_2026_10
    PARTITION OF scans
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');

CREATE TABLE IF NOT EXISTS scans_2026_11
    PARTITION OF scans
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');

-- ── content_flags (partitioned by flagged_at) ─────────────────────────────────
CREATE TABLE IF NOT EXISTS content_flags_2026_09
    PARTITION OF content_flags
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');

CREATE TABLE IF NOT EXISTS content_flags_2026_10
    PARTITION OF content_flags
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');

CREATE TABLE IF NOT EXISTS content_flags_2026_11
    PARTITION OF content_flags
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');

-- ── audit_log (partitioned by occurred_at) ────────────────────────────────────
CREATE TABLE IF NOT EXISTS audit_log_2026_09
    PARTITION OF audit_log
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');

CREATE TABLE IF NOT EXISTS audit_log_2026_10
    PARTITION OF audit_log
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');

CREATE TABLE IF NOT EXISTS audit_log_2026_11
    PARTITION OF audit_log
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');

-- ── proxy_leases (partitioned by issued_at) ───────────────────────────────────
CREATE TABLE IF NOT EXISTS proxy_leases_2026_09
    PARTITION OF proxy_leases
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');

CREATE TABLE IF NOT EXISTS proxy_leases_2026_10
    PARTITION OF proxy_leases
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');

CREATE TABLE IF NOT EXISTS proxy_leases_2026_11
    PARTITION OF proxy_leases
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');

-- ── Default fallback partitions (catches all out-of-range timestamps) ────────
CREATE TABLE IF NOT EXISTS scans_default PARTITION OF scans DEFAULT;
CREATE TABLE IF NOT EXISTS content_flags_default PARTITION OF content_flags DEFAULT;
CREATE TABLE IF NOT EXISTS audit_log_default PARTITION OF audit_log DEFAULT;
CREATE TABLE IF NOT EXISTS proxy_leases_default PARTITION OF proxy_leases DEFAULT;
