import React, { useState, useEffect, useRef, useCallback } from 'react'
import { fetchLogs } from '../services/api'

// Default service order from docker-compose.yml. The API also returns whichever
// services are actually present at runtime, so extras get appended dynamically.
// "pipeline" is a virtual feed (no compose container) — the backend synthesises
// Redis queue + task_status lines for it.
const DEFAULT_SERVICES = [
  'pipeline',
  'redis',
  'postgres',
  'tor-sidecar',
  'crawler-worker',
  'content-shield',
  'storage-api',
  'ingestion-api',
  'correlation-engine',
  'dashboard-api',
  'frontend',
  'nginx',
]

// Virtual AI log feeds — only surfaced (and selectable) when the tool runs in
// AI mode (Ollama engine up). `ai` is the combined view; `ai-seed` is the
// initial-seed / term-expansion stage (qwen2.5:3b) and `ai-guard` is the
// CASM / safety-guard block stage (llama-guard3:1b).
const AI_SERVICES = ['ai', 'ai-seed', 'ai-guard']

const SERVICE_COLORS = {
  pipeline: '#7c4dff',
  redis: '#ff9100',
  postgres: '#00e5ff',
  'tor-sidecar': '#d500f9',
  'crawler-worker': '#00ff9d',
  'content-shield': '#ff1744',
  'storage-api': '#00b0ff',
  'ingestion-api': '#ffea00',
  'correlation-engine': '#e040fb',
  'dashboard-api': '#ff6d00',
  frontend: '#76ff03',
  nginx: '#18ffff',
  ai: '#d500f9',
  'ai-seed': '#aa00ff',
  'ai-guard': '#ff1744',
}

const COLOR_FALLBACKS = ['#00ff9d', '#00e5ff', '#d500f9', '#ff9100', '#ffea00', '#18ffff', '#76ff03', '#e040fb', '#ff6d00']

function colorFor(service) {
  return SERVICE_COLORS[service] || COLOR_FALLBACKS[service.charCodeAt(0) % COLOR_FALLBACKS.length]
}

export default function LiveLogs({ aiMode }) {
  const [service, setService] = useState('all')
  const [services, setServices] = useState(
    aiMode ? [...DEFAULT_SERVICES, ...AI_SERVICES] : DEFAULT_SERVICES
  )
  const [lines, setLines] = useState([])
  const [follow, setFollow] = useState(true)
  const [isLive, setIsLive] = useState(true)
  const logRef = useRef()
  const followRef = useRef(true)

  const refresh = useCallback(async () => {
    const res = await fetchLogs(service === 'all' ? null : service, 60)
    if (!res.live || !res.data) {
      setIsLive(false)
      return
    }
    setIsLive(true)

    // In non-AI mode the AI log feeds are disabled: never collect their lines
    // or register them as selectable services, even if the backend returns them.
    const isVisible = (svc) => aiMode || !AI_SERVICES.includes(svc)

    const flat = []
    const seenServices = new Set()
    if (service === 'all') {
      const map = res.data.services || {}
      Object.entries(map).forEach(([svc, logs]) => {
        if (!isVisible(svc)) return
        seenServices.add(svc)
        if (!Array.isArray(logs)) return
        logs.forEach(txt => flat.push({ svc, text: txt }))
      })
    } else {
      if (!isVisible(service)) {
        setLines([])
        setServices(DEFAULT_SERVICES)
        setIsLive(false)
        return
      }
      seenServices.add(res.data.service || service)
      ;(res.data.logs || []).forEach(txt => flat.push({ svc: res.data.service || service, text: txt }))
    }

    if (flat.length) {
      setLines(prev => {
        // Dedupe against EVERY line already rendered (not just the last one).
        // The API returns the same `--tail 60` window each poll, so without
        // this the whole window gets re-appended as duplicates every 2s and
        // real new output is drowned out — making the panel look frozen.
        const seen = new Set(prev.map(p => p.key))
        const merged = [...prev]
        for (const line of flat) {
          const key = `${line.svc}|${line.text}`
          if (seen.has(key)) continue
          seen.add(key)
          merged.push({ ...line, key })
        }
        return merged.slice(-250)
      })
    }

    setServices(prev => {
      const existing = new Set(prev)
      const combined = [...prev]
      seenServices.forEach(s => {
        if (!existing.has(s)) combined.push(s)
      })
      return combined
    })
  }, [service, aiMode])

  useEffect(() => {
    refresh()
    const timer = setInterval(refresh, 2000)
    return () => clearInterval(timer)
  }, [refresh])

  // Keep "follow" in sync to avoid stale closure in the scroll handler.
  useEffect(() => {
    followRef.current = follow
  }, [follow])

  // Sync available services + selected feed when AI mode toggles on/off so the
  // AI log feeds only appear/remain when the tool actually runs in AI mode.
  useEffect(() => {
    setServices(prev => {
      const hasAI = prev.some(s => AI_SERVICES.includes(s))
      if (aiMode && !hasAI) return [...DEFAULT_SERVICES, ...AI_SERVICES]
      if (!aiMode && hasAI) {
        // Drop AI feeds; if one was selected, fall back to "all".
        if (AI_SERVICES.includes(service)) setService('all')
        return DEFAULT_SERVICES
      }
      return prev
    })
  }, [aiMode, service])

  // Auto-scroll to the newest line while follow mode is on.
  useEffect(() => {
    const el = logRef.current
    if (el && followRef.current) {
      el.scrollTop = el.scrollHeight
    }
  }, [lines])

  function handleScroll() {
    const el = logRef.current
    if (!el) return
    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 40
    if (atBottom) {
      followRef.current = true
      setFollow(true)
    } else {
      followRef.current = false
      setFollow(false)
    }
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: 0 }}>
      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 4 }}>
        <div style={{ fontSize: 8, fontFamily: 'var(--font-mono)', color: 'var(--green)', textTransform: 'uppercase', letterSpacing: '0.04em' }}>
          <span style={{
            display: 'inline-block', width: 7, height: 7, borderRadius: '50%',
            background: isLive ? '#00ff9d' : '#ff1744',
            boxShadow: isLive ? '0 0 6px #00ff9d' : '0 0 6px #ff1744',
            marginRight: 5,
          }} />
          LIVE PIPELINE LOGS {isLive ? '' : '(OFFLINE)'}
        </div>
        <button
          onClick={() => { setFollow(true); followRef.current = true; if (logRef.current) logRef.current.scrollTop = logRef.current.scrollHeight }}
          style={{
            background: 'none', border: 'none', cursor: 'pointer',
            color: follow ? 'var(--green)' : 'var(--text-muted)',
            fontSize: 8, fontFamily: 'var(--font-mono)',
          }}
          title="Toggle follow-mode"
        >
          {follow ? '● FOLLOW:ON' : '○ FOLLOW:OFF'}
        </button>
      </div>

      {/* Component chips */}
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: 3, marginBottom: 5 }}>
        <button
          onClick={() => setService('all')}
          style={{
            border: service === 'all' ? '1px solid var(--green)' : '1px solid var(--border)',
            background: service === 'all' ? 'rgba(0,255,157,0.1)' : 'rgba(0,0,0,0.3)',
            color: service === 'all' ? 'var(--green)' : 'var(--text-secondary)',
            fontSize: 8, fontFamily: 'var(--font-mono)', cursor: 'pointer',
            borderRadius: 3, padding: '1px 5px',
          }}
        >
          ALL
        </button>
        {services.map(svc => (
          <button
            key={svc}
            onClick={() => setService(svc)}
            style={{
              border: service === svc ? '1px solid var(--cyan)' : '1px solid var(--border)',
              background: service === svc ? 'rgba(0,229,255,0.1)' : 'rgba(0,0,0,0.3)',
              color: service === svc ? 'var(--cyan)' : colorFor(svc),
              fontSize: 8, fontFamily: 'var(--font-mono)', cursor: 'pointer',
              borderRadius: 3, padding: '1px 5px',
            }}
          >
            {svc}
          </button>
        ))}
      </div>

      {/* Log area */}
      <div
        ref={logRef}
        onScroll={handleScroll}
        style={{
          background: '#01030a',
          border: '1px solid var(--border)',
          borderRadius: 4,
          flex: 1,
          minHeight: 0,
          overflowY: 'auto',
          padding: '6px 8px',
          fontFamily: 'var(--font-mono)',
          fontSize: 8,
          lineHeight: 1.5,
          whiteSpace: 'pre-wrap',
          wordBreak: 'break-all',
        }}
      >
        {lines.length === 0 && (
          <div style={{ color: 'var(--text-muted)' }}>waiting for container activity — start a scan or watch the pipeline...</div>
        )}
        {lines.map((l, i) => (
          <div key={l.key + '-' + i} style={{ display: 'flex', gap: 5 }}>
            <span style={{ color: colorFor(l.svc), flexShrink: 0, fontWeight: 700 }}>
              {l.svc}:
            </span>
            <span style={{ color: /fatal|error|panic/i.test(l.text) ? '#ff8080' : '#88a8c4' }}>
              {l.text}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}