// NOTE: We intentionally removed fallback mock data for production delete work.
// If backend is unavailable, UI will show empty results instead.

// If the UI is served behind the same reverse-proxy as the Go services,
// requests should be sent to the proxy prefix.
// In docker/dev this is typically `/api/v1/dashboard`.
//
// NOTE: Keep these default paths aligned with nginx/reverse-proxy routing.
const DASHBOARD_BASE = import.meta.env.VITE_DASHBOARD_URL || '/api/v1/dashboard'
const INGESTION_BASE = import.meta.env.VITE_INGESTION_URL || '/api/v1/investigations'

// Debug helper: allow quick visibility in console without rebuilding.
const SCREENSHOT_ENDPOINTS_BASE = import.meta.env.VITE_SCREENSHOT_BASE || DASHBOARD_BASE

export async function fetchStats() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/stats`, { signal: AbortSignal.timeout(2000) })
    if (!res.ok) throw new Error('API Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

export async function fetchGraph() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/graph`, { signal: AbortSignal.timeout(2000) })
    if (!res.ok) throw new Error('API Error')
    const json = await res.json()
    return { data: json, live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

export async function fetchLeads() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/leads`, { signal: AbortSignal.timeout(2000) })
    if (!res.ok) throw new Error('API Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

// fetchMedia returns the extracted-evidence gallery: every image reference the
// crawler captured, joined with its owning actor. The sidebar MEDIA tab renders
// these as thumbnails; clicking one deep-links into the actor's dossier.
export async function fetchMedia() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/media`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error('Media API Error')
    const json = await res.json()
    return { data: Array.isArray(json?.images) ? json.images : [], live: true }
  } catch (err) {
    return { data: [], live: false, error: err?.message }
  }
}

export async function fetchJobs() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/jobs`, { signal: AbortSignal.timeout(2000) })
    if (!res.ok) throw new Error('API Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

// fetchLogs tails REAL docker logs for a compose service (or all services when
// service is null/omitted). The dashboard-api reads them straight from the
// Docker engine socket, so this is the same output as `docker compose logs`.
export async function fetchLogs(service = null, tail = 60) {
  try {
    const params = new URLSearchParams()
    if (service) params.set('service', service)
    params.set('tail', String(tail))
    const res = await fetch(`${DASHBOARD_BASE}/logs?${params}`, { signal: AbortSignal.timeout(4000) })
    if (!res.ok) throw new Error('Log API Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

// fetchAiStatus reports whether the stack is currently running in AI mode
// (i.e. the Ollama engine container is up).
export async function fetchAiStatus() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/ai/status`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error('AI status Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

// fetchComponents returns the live status of every OnionEye service so the
// right-rail TOOL MAP can render components glowing green/red by state.
export async function fetchComponents() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/components`, { signal: AbortSignal.timeout(4000) })
    if (!res.ok) throw new Error('Components API Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

export async function fetchTorStatus() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/tor/status`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error('Tor status Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

// restartStackWithAi restarts the entire docker-compose stack with the AI
// profile (enabled=true => Ollama) or without it (enabled=false). The backend
// drives `docker compose ... up -d --force-recreate` and the whole tool comes
// back up in ~15s, so callers should prompt the user before invoking this.
export async function restartStackWithAi(enabled) {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/ai/restart`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled: !!enabled }),
      signal: AbortSignal.timeout(150000),
    })
    const json = await res.json().catch(() => ({}))
    return { ok: res.ok, data: json }
  } catch (err) {
    return { ok: false, error: err?.message || 'Service unavailable' }
  }
}

export async function deleteNode(payload) {
  const DASHBOARD_BASE = import.meta.env.VITE_DASHBOARD_URL || '/api/v1/dashboard'
  const res = await fetch(`${DASHBOARD_BASE}/delete/node`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
    signal: AbortSignal.timeout(5000),
  })
  if (!res.ok) {
    const t = await res.text().catch(() => '')
    throw new Error(t || `Delete node failed (${res.status})`)
  }
  return await res.json().catch(() => ({}))
}

export async function deleteAll() {
  const DASHBOARD_BASE = import.meta.env.VITE_DASHBOARD_URL || '/api/v1/dashboard'
  const res = await fetch(`${DASHBOARD_BASE}/delete/all`, {
    method: 'DELETE',
    signal: AbortSignal.timeout(5000),
  })
  if (!res.ok) {
    const t = await res.text().catch(() => '')
    throw new Error(t || `Delete all failed (${res.status})`)
  }
  return await res.json().catch(() => ({}))
}

export async function submitInvestigation(payload) {
  try {
    const res = await fetch(INGESTION_BASE, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
      signal: AbortSignal.timeout(5000),
    })
    const json = await res.json()
    return { ok: res.ok, data: json }
  } catch (err) {
    return { ok: false, error: err.message || 'Service unavailable' }
  }
}

// fetchManualReviewQueue fetches pending manual review items from the
// crawl:manual Redis stream. Items appear here when the crawler hits CAPTCHA,
// bot blocks, or repeated scan failures requiring human intervention.
export async function fetchManualReviewQueue() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/review/pending`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error('Review queue Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

// fetchNotifications polls the dashboard's in-memory notification ring, where
// crawler-worker posts worker-rotation and worker-anomaly events. The UI shows
// a transient popup for each new event (deduplicated by rotation_id/anomaly_id).
export async function fetchNotifications() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/notifications`, { signal: AbortSignal.timeout(2000) })
    if (!res.ok) throw new Error('Notifications API Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

// dismissReviewItem removes a manual review item from the queue after the
// analyst has resolved it (either by manually visiting the site or deciding to
// skip it).
export async function dismissReviewItem(id) {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/review/dismiss`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id }),
      signal: AbortSignal.timeout(5000),
    })
    const json = await res.json().catch(() => ({}))
    return { ok: res.ok, data: json }
  } catch (err) {
    return { ok: false, error: err?.message || 'Service unavailable' }
  }
}

// fetchExistingScreenshot checks if a rendered PNG screenshot already exists in storage
export async function fetchExistingScreenshot(url) {
  try {
    const res = await fetch(`${SCREENSHOT_ENDPOINTS_BASE}/screenshot?url=${encodeURIComponent(url)}`, {
      signal: AbortSignal.timeout(4000),
    })
    if (!res.ok) return { ok: false }
    const blob = await res.blob()
    if (!blob || blob.size === 0) return { ok: false }
    return { ok: true, imageUrl: URL.createObjectURL(blob) }
  } catch {
    return { ok: false }
  }
}

// captureScreenshot triggers an on-demand, exploit-resistant Chromium snapshot
// through crawler-worker & render-sandbox with Content Shield verification.
export async function captureScreenshot(url) {
  try {
    const res = await fetch(`${SCREENSHOT_ENDPOINTS_BASE}/screenshot/capture`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url }),
      signal: AbortSignal.timeout(65000),
    })
    if (!res.ok) {
      const errJson = await res.json().catch(() => ({}))
      throw new Error(errJson.error || `Screenshot capture failed (${res.status})`)
    }
    const blob = await res.blob()
    return { ok: true, imageUrl: URL.createObjectURL(blob), live: true }
  } catch (err) {
    return { ok: false, error: err?.message || 'Screenshot service unavailable' }
  }
}
// ── Negative-control console API ──────────────────────────────────────────────
// Rules are (mode × artifact_type) damping policies; artifacts are concrete
// values with occurrence counts backing FREQUENCY-style dampening.

export async function fetchNegativeControl() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/negative-control`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error('Negative control API Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

export async function saveNegativeControlRule(payload) {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/negative-control`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
      signal: AbortSignal.timeout(5000),
    })
    const json = await res.json().catch(() => ({}))
    return { ok: res.ok, data: json }
  } catch (err) {
    return { ok: false, error: err?.message || 'Service unavailable' }
  }
}

export async function deleteNegativeControlRule(id) {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/negative-control/${id}`, {
      method: 'DELETE',
      signal: AbortSignal.timeout(5000),
    })
    const json = await res.json().catch(() => ({}))
    return { ok: res.ok, data: json }
  } catch (err) {
    return { ok: false, error: err?.message || 'Service unavailable' }
  }
}

export async function fetchArtifacts(type = '', mode = '') {
  try {
    const params = new URLSearchParams()
    if (type) params.set('type', type)
    if (mode) params.set('mode', mode)
    const qs = params.toString() ? `?${params.toString()}` : ''
    const res = await fetch(`${DASHBOARD_BASE}/artifacts${qs}`, { signal: AbortSignal.timeout(3000) })
    if (!res.ok) throw new Error('Artifacts API Error')
    return { data: await res.json(), live: true }
  } catch (err) {
    return { data: null, live: false, error: err?.message }
  }
}

export async function addNegativeControlArtifact(payload) {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/negative-control/artifacts`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
      signal: AbortSignal.timeout(5000),
    })
    const json = await res.json().catch(() => ({}))
    return { ok: res.ok, data: json }
  } catch (err) {
    return { ok: false, error: err?.message || 'Service unavailable' }
  }
}

// nukeAll wipes every docker volume and restarts the whole stack.
export async function nukeAll() {
  try {
    const res = await fetch(`${DASHBOARD_BASE}/system/nuke`, {
      method: 'POST',
      signal: AbortSignal.timeout(5000),
    })
    const json = await res.json().catch(() => ({}))
    return { ok: res.ok, data: json }
  } catch (err) {
    return { ok: false, error: err?.message || 'Service unavailable' }
  }
}
