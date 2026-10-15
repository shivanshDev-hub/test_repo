-- ===============================================================================
-- OnionEye — Demo Seed Data (live dashboard test dataset)
-- -------------------------------------------------------------------------------
-- Populates a multi-actor investigation graph so the dashboard shows realistic
-- live numbers in the stats bar, graph, leads panel, and jobs panel.
--
-- USAGE (containerized):
--   docker compose exec -T postgres psql -U onioneye -d onioneye < db/seed.sql
--
-- Safe to run repeatedly: every INSERT is idempotent via ON CONFLICT DO NOTHING.
-- ===============================================================================

-- 1. ACTORS (9 persona nodes)
INSERT INTO actors (actor_id, primary_handle, category, last_scan_date, created_at) VALUES
    ('00000000-0000-0000-0000-000000000001', 'phantom_matrix',        'COORDINATOR',    '2026-09-01T03:12:00Z', '2026-08-20T10:00:00Z'),
    ('00000000-0000-0000-0000-000000000002', 'cipher_syndicate',      'MARKETPLACE',    '2026-09-02T07:44:00Z', '2026-08-21T14:30:00Z'),
    ('00000000-0000-0000-0000-000000000003', 'darkmint_exchange',     'FINANCIAL',      '2026-09-03T05:23:00Z', '2026-08-22T09:15:00Z'),
    ('00000000-0000-0000-0000-000000000004', 'northstar_pgp',         'SERVICES',       '2026-09-04T11:01:00Z', '2026-08-23T16:45:00Z'),
    ('00000000-0000-0000-0000-000000000005', 'redlantern_market',     'MARKETPLACE',    '2026-09-05T02:18:00Z', '2026-08-24T08:05:00Z'),
    ('00000000-0000-0000-0000-000000000006', 'torrent_tunnel_infra',  'INFRASTRUCTURE', '2026-09-06T14:52:00Z', '2026-08-25T12:20:00Z'),
    ('00000000-0000-0000-0000-000000000007', 'silkroute_supplier',    'DISTRIBUTOR',    '2026-09-06T18:37:00Z', '2026-08-26T19:55:00Z'),
    ('00000000-0000-0000-0000-000000000008', 'quantum_ledger',        'FINANCIAL',      '2026-09-07T09:26:00Z', '2026-08-27T07:40:00Z'),
    ('00000000-0000-0000-0000-000000000009', 'midnight_escrow',       'FINANCIAL',      '2026-09-08T16:09:00Z', '2026-08-28T11:10:00Z')
ON CONFLICT (primary_handle) DO NOTHING;

-- 2. IDENTIFIERS
-- Wallet kinds matched by: kind ILIKE '%wallet%' OR '%btc%' OR '%xmr%'.
-- PGP kinds matched by:     kind ILIKE '%pgp%'.
INSERT INTO identifiers (identifier_id, actor_id, kind, value, source, provenance_tag, first_seen, last_seen) VALUES
    ('01000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'btc_wallet', 'bc1qphantommatrix0x9f', 'crawler', 'market_parsed', '2026-08-20T10:00:00Z', '2026-09-01T03:12:00Z'),
    ('01000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', 'btc_wallet', 'bc1qcipher7s7ndc4a2',  'crawler', 'listing_page',  '2026-08-21T14:30:00Z', '2026-09-02T07:44:00Z'),
    ('01000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000003', 'xmr_wallet', '8BdVhS2darkmint3',     'crawler', 'escrow_page',   '2026-08-22T09:15:00Z', '2026-09-03T05:23:00Z'),
    ('01000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000004', 'btc_wallet', 'bc1qnorthstar4f1e9a',  'crawler', 'forum_sig',     '2026-08-23T16:45:00Z', '2026-09-04T11:01:00Z'),
    ('01000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000005', 'xmr_wallet', '4AD6rKpredlantern',    'crawler', 'listing_page',  '2026-08-24T08:05:00Z', '2026-09-05T02:18:00Z'),
    ('01000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000009', 'btc_wallet', 'bc1qzerotrust8c3b7f',  'crawler', 'escrow_page',   '2026-08-28T11:10:00Z', '2026-09-08T16:09:00Z'),
    ('01000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000004', 'pgp_fingerprint', '3B4FC9A217DE88B10E459C31AA62D8F01E7702CD', 'crawler', 'trusted_key', '2026-08-23T16:45:00Z', '2026-09-04T11:01:00Z'),
    ('01000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000007', 'pgp_fingerprint', '9D28E1AB44C066F23B0751E8F7A0C3DD906B14F2', 'crawler', 'market_sig',  '2026-08-26T19:55:00Z', '2026-09-06T18:37:00Z'),
    ('01000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000001', 'telegram',    '@phantom_matrix',      'osint',   'manual_entry',  '2026-08-20T10:00:00Z', '2026-09-01T03:12:00Z'),
    ('01000000-0000-0000-0000-00000000000a', '00000000-0000-0000-0000-000000000002', 'telegram',    '@cipher_syndicate',    'osint',   'manual_entry',  '2026-08-21T14:30:00Z', '2026-09-02T07:44:00Z'),
    ('01000000-0000-0000-0000-00000000000b', '00000000-0000-0000-0000-000000000003', 'email',       'ops@darkmint.example', 'osint',   'leak_db',       '2026-08-22T09:15:00Z', '2026-09-03T05:23:00Z'),
    ('01000000-0000-0000-0000-00000000000c', '00000000-0000-0000-0000-000000000005', 'onion_alias', 'redlantern_market',    'crawler', 'site_title',    '2026-08-24T08:05:00Z', '2026-09-05T02:18:00Z')
ON CONFLICT (actor_id, kind, value) DO NOTHING;

-- 3. INFRASTRUCTURE INDICATORS
INSERT INTO infra_indicators (indicator_id, actor_id, indicator_type, onion_host, matched_clearnet_host, evidence, detected_at) VALUES
    ('02000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'MARKET_ONION', 'http://cipher7snd4ca2.onion',  'cipher-market.example',    'linked via PGP sig on about page', '2026-08-21T14:30:00Z'),
    ('02000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000006', 'MIRROR_HOST',  'http://torrenttunnel1x.onion', 'torrent-proxy.example',    'DNS + TLS cert reuse',             '2026-08-25T12:20:00Z'),
    ('02000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000005', 'MARKET_ONION', 'http://redlanternabc.onion',   'redlantern-shop.example',  'shared styling fingerprint',       '2026-08-24T08:05:00Z'),
    ('02000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000007', 'BLOG_ONION',   'http://silkroute0f2a.onion',   NULL,                       'referenced in vendor thread',      '2026-08-26T19:55:00Z'),
    ('02000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000009', 'ESCROW_ONION', 'http://midnightescrow1.onion', NULL,                       'escrow payment instr page',        '2026-08-28T11:10:00Z')
ON CONFLICT (actor_id, indicator_type, onion_host) DO NOTHING;

-- 4. PERSONA LINKS
-- actor_id_a must be < actor_id_b (schema constraint).
-- nc_state = PASSED + confidence >= category threshold => visible in v_active_persona_links.
INSERT INTO persona_links (link_id, actor_id_a, actor_id_b, basis, confidence, negative_control_passed,
                           nc_state, nc_penalty_score, nc_audit_log, created_at) VALUES
    ('03000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'STYLOMETRY',     0.88, TRUE,  'PASSED', 0.00, '{"round":2,"static_checks":["style_ngram"]}'::jsonb, '2026-08-30T08:00:00Z'),
    ('03000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000003', 'FINANCIAL',      0.79, TRUE,  'PASSED', 0.00, '{"round":1,"static_checks":["wallet_graph"]}'::jsonb, '2026-08-30T09:00:00Z'),
    ('03000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000004', 'STYLOMETRY',     0.83, TRUE,  'PASSED', 0.00, '{"round":2,"static_checks":["style_ngram"]}'::jsonb, '2026-08-30T10:00:00Z'),
    ('03000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000005', 'INFRASTRUCTURE', 0.72, TRUE,  'PASSED', 0.00, '{"round":3,"static_checks":["tls_cert"]}'::jsonb, '2026-08-30T11:00:00Z'),
    ('03000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000007', 'STYLOMETRY',     0.91, TRUE,  'PASSED', 0.00, '{"round":1,"static_checks":["style_ngram"]}'::jsonb, '2026-08-30T12:00:00Z'),
    ('03000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000006', 'INFRASTRUCTURE', 0.69, TRUE,  'PASSED', 0.00, '{"round":2,"static_checks":["dns_reuse"]}'::jsonb, '2026-08-30T13:00:00Z'),
    ('03000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000008', 'FINANCIAL',      0.62, TRUE,  'PASSED', 0.00, '{"round":1,"static_checks":["wallet_graph"]}'::jsonb, '2026-08-30T14:00:00Z'),
    ('03000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000004', 'STYLOMETRY',     0.55, TRUE,  'PASSED', 0.00, '{"round":3,"static_checks":["style_ngram"]}'::jsonb, '2026-08-30T15:00:00Z'),
    ('03000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000005', 'FINANCIAL',      0.42, FALSE, 'FAILED', 0.55, '{"round":1,"failure":"confidence_below_policy"}'::jsonb, '2026-08-30T16:00:00Z')
ON CONFLICT (actor_id_a, actor_id_b, basis) DO NOTHING;

-- 5. SCAN RECORDS (jobs panel). status = completed counted in stats jobs_done.
INSERT INTO scans (scan_id, requested_by, targets, status, started_at, completed_at) VALUES
    ('04000000-0000-0000-0000-000000000001', 'analyst_demo',  ARRAY['http://cipher7snd4ca2.onion'], 'completed', '2026-09-01T03:00:00Z', '2026-09-01T03:12:00Z'),
    ('04000000-0000-0000-0000-000000000002', 'analyst_demo',  ARRAY['http://redlanternabc.onion'],   'completed', '2026-09-02T07:30:00Z', '2026-09-02T07:44:00Z'),
    ('04000000-0000-0000-0000-000000000003', 'auto_pipeline', ARRAY['http://midnightescrow1.onion'], 'completed', '2026-09-03T05:10:00Z', '2026-09-03T05:23:00Z'),
    ('04000000-0000-0000-0000-000000000004', 'analyst_demo',  ARRAY['http://torrenttunnel1x.onion'], 'running',   '2026-09-08T16:00:00Z', NULL),
    ('04000000-0000-0000-0000-000000000005', 'auto_pipeline', ARRAY['http://silkroute0f2a.onion'],   'queued',    '2026-09-08T16:05:00Z', NULL)
ON CONFLICT (scan_id, started_at) DO NOTHING;

-- 6. CONTENT FLAGS (stats flagged)
INSERT INTO content_flags (flag_id, scan_id, source_url, category, action_taken, actor_id, flagged_at) VALUES
    ('05000000-0000-0000-0000-000000000001', '04000000-0000-0000-0000-000000000001', 'http://cipher7snd4ca2.onion/listing/8832', 'ILLICIT_GOODS',   'BLOCKED',   '00000000-0000-0000-0000-000000000002', '2026-09-01T03:12:00Z'),
    ('05000000-0000-0000-0000-000000000002', '04000000-0000-0000-0000-000000000002', 'http://redlanternabc.onion/vendor',       'FRAUD_SERVICE',   'BLOCKED',   '00000000-0000-0000-0000-000000000005', '2026-09-02T07:44:00Z'),
    ('05000000-0000-0000-0000-000000000003', '04000000-0000-0000-0000-000000000003', 'http://midnightescrow1.onion/pay',        'FINANCIAL_FRAUD', 'ESCALATED', '00000000-0000-0000-0000-000000000009', '2026-09-03T05:23:00Z')
ON CONFLICT (flag_id, flagged_at) DO NOTHING;

-- 7. AUDIT LOG (no unique constraint on the partitioned parent, so skip ON CONFLICT)
INSERT INTO audit_log (log_id, actor_type, action, resource_id, occurred_at, metadata)
SELECT * FROM (VALUES
    ('06000000-0000-0000-0000-000000000001'::uuid, 'analyst', 'graph_export',     '00000000-0000-0000-0000-000000000001'::uuid, '2026-09-05T10:00:00Z'::timestamptz, '{"format":"json","scope":"network"}'::jsonb),
    ('06000000-0000-0000-0000-000000000002'::uuid, 'system',  'filter_policy_ev', '00000000-0000-0000-0000-000000000007'::uuid, '2026-09-06T09:00:00Z'::timestamptz, '{"nc_round":2}'::jsonb)
) AS v(log_id, actor_type, action, resource_id, occurred_at, metadata)
WHERE NOT EXISTS (SELECT 1 FROM audit_log a WHERE a.log_id = v.log_id);
