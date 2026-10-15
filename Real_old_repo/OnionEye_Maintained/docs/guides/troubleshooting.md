# OnionEye Troubleshooting Guide

## How to use
When something fails, start from the **first failing component** in the request path, then verify the **stream/state transitions** and finally check **DB write constraints**.

Typical path:
`frontend → dashboard-api → ingestion-api → redis(crawl:jobs) → crawler-worker → content-shield → correlation-engine → storage-api → postgres → dashboard-api`

---

## 1) `POST /scans` returns HTTP 500 (UUID actor_id empty)

### Symptom
- `POST /scans` intermittent 500
- Log line contains: `invalid input syntax for type uuid: ""`

### Why it happens
- The crawler/extractor sometimes produces infra/links/images but no valid actor handle.
- The derived `firstActorID` is empty, and actor-scoped UUID inserts must not run.

### Fix (already implemented)
- `storage-api CreateScan()` skips actor-scoped inserts into:
  - `infra_indicators`
  - `extracted_links`
  - `image_hashes`
  when `firstActorID` is empty or invalid UUID.

### Verify
Run:
```bash
docker compose logs --tail=200 storage-api | egrep -n "POST /scans|invalid input syntax for type uuid|panic|error" || true
```
You should **not** see `invalid input syntax for type uuid: ""`.

---

## 2) Jobs not appearing in dashboard

### Checklist
1. Confirm ingestion dispatch:
   - `ingestion-api` logs show job enqueue.
2. Confirm Redis stream:
   - `redis` has entries in `crawl:jobs`.
3. Confirm crawler worker consumes:
   - crawler-worker logs show `XREADGROUP` activity.
4. Confirm storage persistence:
   - storage-api logs show successful scan persistence.

---

## 3) Graph shows empty / missing nodes

### Checklist
1. Confirm dashboard API endpoints respond 200:
   - `/api/v1/dashboard/graph`
   - `/api/v1/dashboard/stats`
2. Confirm correlation-engine persona upserts:
   - storage tables for actor/persona links populated.
3. Check Phase 8 read-time gate:
   - ensure the access gate policy isn’t filtering out all rows for the current role.

---

## 4) Delete node returns "missing node type" or "node_id must be a UUID" error

### Symptom
- `DELETE /api/v1/dashboard/delete/node` returns HTTP 400
- Error: `node_id and node_type are required` OR `missing node type for id=<identifier-id>` OR `node_id must be a UUID`

### Why it happens
- **Frontend**: `handleDeleteNode` only checked old prefixes (`actor_`, `email_`, etc.) for non-actor nodes.
- **Storage-API**: `DeleteNode` required UUID for all deletions; identifier/infra nodes have string IDs like `id_wallet:*`, `infra_*`

### Fix (already implemented)
- **Frontend**: `handleDeleteNode` now handles identifier/infra node ID patterns:
  - `id_email_*` → type: `email`
  - `id_wallet:*` → type: `wallet`
  - `id_pgp_key_*` → type: `pgp`
  - `id_onion_*` → type: `onion`
  - `infra_*` → type: `infra`
- **Storage-API**: `DeleteNode` now supports non-UUID node IDs via prefix matching (`id_` and `infra_`), converting to DB `kind || '_' || value` and `indicator_type || '_' || onion_host` columns.

### Verify
```bash
# Test delete identifier node
curl -X DELETE http://localhost:8081/delete/node \
  -H 'Content-Type: application/json' \
  -d '{"node_id":"id_email_support@darkvault.pro","node_type":"email"}'

# Test delete infra node
curl -X DELETE http://localhost:8081/delete/node \
  -H 'Content-Type: application/json' \
  -d '{"node_id":"infra_MARKET_ONION_http://redlanternabc.onion","node_type":"infra"}'
```

---

## 5) Content filter rejects too aggressively

### Checklist
- Verify classifier/weights presence in `content-shield` container.
- Ensure fail-open mode is enabled if desired.
- Inspect `content_flags` in postgres: are rows being written by Phase 8A?

---

## 6) Delete actor returns FK constraint violation

### Symptom
- `DELETE /api/v1/dashboard/delete/node` returns HTTP 500
- Error: `update or delete on table "actors" violates foreign key constraint`

### Why it happens
- Child tables (`identifiers`, `infra_indicators`, `extracted_links`, `image_hashes`, `content_flags`, etc.) have FK references to `actors.actor_id`
- The schema does not declare `ON DELETE CASCADE`, so child rows must be deleted before the parent actor

### Fix (already implemented)
- `storage-api DeleteNode()` deletes in this order:
  1. `persona_links` (references actor_id_a/b)
  2. `identifiers` (FK → actors)
  3. `infra_indicators` (FK → actors)
  4. `actors` (parent table)
- Additional FK children (`extracted_links`, `image_hashes`, `content_flags`, `stylometric_profiles`) are either rare in the seed data or have `ON DELETE CASCADE` already set

### Verify
```bash
curl -s -X DELETE http://localhost:3000/api/v1/dashboard/delete/node \
  -H 'Content-Type: application/json' \
  -d '{"node_id":"00000000-0000-0000-0000-000000000009","node_type":"actor"}' | jq .
```
Should return `{"deleted_actors":1, "deleted_identifiers":..., ...}` without error.

---

## 7) Tor-sidecar / egress failures

### Symptom
- network timeouts, repeated retry budget exhausted

### Checklist
1. Tor health endpoints from sidecar.
2. iptables killswitch rules load correctly.
3. Circuit renewal configuration matches expected ports.

---

## Need more help?
Provide:
- request id / job id
- time window
- last 200 log lines from `ingestion-api`, `crawler-worker`, `storage-api`
- whether AI profile (`--profile ai`) was used
