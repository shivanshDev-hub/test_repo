import React, { useState, useEffect, useMemo, useRef, useCallback } from 'react'
import AiRestartModal from './AiRestartModal'
import { fetchComponents } from '../services/api'

/* ── PanelShell ────────────────────────────────────────────────────────────── */
function PanelShell({ title, children }) {
  return (
    <div style={{
      flex: 1, overflow: 'hidden',
      background: 'linear-gradient(180deg, #070c14 0%, #04070c 100%)',
      border: '1px solid var(--border)', borderRadius: 6,
      display: 'flex', flexDirection: 'column', minHeight: 0,
    }}>
      <div style={{
        padding: '10px 10px', borderBottom: '1px solid var(--border)',
        background: 'rgba(5, 8, 13, 0.95)',
        display: 'flex', alignItems: 'center',
      }}>
        <span style={{ fontFamily: 'var(--font-mono)', fontSize: 10, color: 'var(--green)', fontWeight: 800 }}>
          {title}
        </span>
      </div>
      <div style={{ flex: 1, overflowY: 'auto', padding: 10 }}>{children}</div>
    </div>
  )
}

/* ── TinyRow (label/value) ─────────────────────────────────────────────────── */
function TinyRow({ k, v, color }) {
  return (
    <div style={{ display: 'flex', justifyContent: 'space-between', gap: 10, marginBottom: 6 }}>
      <span style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>{k}</span>
      <span style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 800, color: color || 'var(--text-primary)', textAlign: 'right', wordBreak: 'break-word' }}>{v}</span>
    </div>
  )
}

/* ── ToolNode (glowing component box in TOOL MAP) ──────────────────────────── */
function ToolNode({ comp }) {
  const up = comp.status === 'up'
  const degraded = comp.status === 'degraded'
  const color = up ? 'var(--green)' : degraded ? 'var(--yellow)' : 'var(--red)'
  const anim = up ? 'pulse-green 1.6s infinite' : 'pulse-red 1.6s infinite'
  return (
    <div style={{
      borderRadius: 8,
      border: `1px solid ${up ? 'rgba(0,255,157,0.55)' : degraded ? 'rgba(255,209,102,0.45)' : 'rgba(255,23,68,0.45)'}`,
      background: up ? 'rgba(0,255,157,0.06)' : degraded ? 'rgba(255,209,102,0.06)' : 'rgba(255,23,68,0.06)',
      padding: '8px 10px',
      display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 10,
    }}>
      <div style={{ minWidth: 0 }}>
        <div style={{ fontFamily: 'var(--font-mono)', fontSize: 10, fontWeight: 800, color, letterSpacing: '0.02em', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{comp.name}</div>
        <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)', marginTop: 2, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{comp.role}</div>
      </div>
      <span style={{
        width: 10, height: 10, borderRadius: '50%', flexShrink: 0,
        background: color, boxShadow: `0 0 10px ${color}40`,
        animation: anim,
      }} />
    </div>
  )
}

/* ═══ Main RightRail ════════════════════════════════════════════════════════ */
export default function RightRail({
  searchQuery, setSearchQuery, graphData, onSelectNode,
  torStatus, torLive,
  aiMode, aiRestarting, onRestartStack,
  onMergeGraph,
}) {
  const [active, setActive] = useState('search')
  const [searchResults, setSearchResults] = useState([])
  const [isSearchOpen, setIsSearchOpen] = useState(false)
  const [aiModal, setAiModal] = useState(null) // null | true | false
  const [components, setComponents] = useState([])
  const [componentSummary, setComponentSummary] = useState(null)
  const [componentsLoading, setComponentsLoading] = useState(false)
  const [importMsg, setImportMsg] = useState(null)
  const [exportMsg, setExportMsg] = useState(null)
  const [importedFileName, setImportedFileName] = useState(null)
  const [importedGraph, setImportedGraph] = useState(null) // { nodes, links }
  const [saveRecordMsg, setSaveRecordMsg] = useState(null)
  const [savedRecords, setSavedRecords] = useState([])
  const [savedRecordsLoading, setSavedRecordsLoading] = useState(false)
  const [loadMsg, setLoadMsg] = useState(null)
  const fileInputRef = useRef(null)
  const width = 290
  const [collapsed, setCollapsed] = useState(false)

  /* ── Item strip (mirrors LeftControlRail) ──────────────────────────────── */
  const items = useMemo(() => [
    { id: 'search',     label: 'SEARCH',         icon: '🔎' },
    { id: 'tor',        label: 'TOR CONNECTION',  icon: '🧅' },
    { id: 'ai',         label: 'AI ENGINE',      icon: '🧠' },
    { id: 'toolmap',    label: 'TOOL MAP',        icon: '🗺' },
    { id: 'export',     label: 'EXPORT',         icon: '📦' },
    { id: 'openrecord', label: 'OPEN RECORD',    icon: '📂' },
  ], [])
/* ── Search effect (migrated from Header) ──────────────────────────────── */
  useEffect(() => {
    if (!searchQuery?.trim() || !graphData?.nodes) {
      setSearchResults([])
      setIsSearchOpen(false)
      return
    }
    const q = searchQuery.toLowerCase()
    const matches = graphData.nodes.filter(n =>
      (n.label && n.label.toLowerCase().includes(q)) ||
      (n.id && n.id.toLowerCase().includes(q)) ||
      (n.type && n.type.toLowerCase().includes(q)) ||
      (n.market && n.market.toLowerCase().includes(q))
    ).slice(0, 8)
    setSearchResults(matches)
    setIsSearchOpen(true)
  }, [searchQuery, graphData])

  /* ── Components polling (TOOL MAP) ─────────────────────────────────────── */
  const loadComponents = useCallback(async () => {
    setComponentsLoading(true)
    const res = await fetchComponents()
    if (res.live && res.data?.components) {
      setComponents(res.data.components)
      setComponentSummary(res.data.summary || null)
    }
    setComponentsLoading(false)
  }, [])

  useEffect(() => {
    if (active === 'toolmap') {
      loadComponents()
      const iv = setInterval(loadComponents, 5000)
      return () => clearInterval(iv)
    }
  }, [active, loadComponents])

  /* ── Saved records polling (OPEN RECORD) ───────────────────────────────── */
  const loadSavedRecords = useCallback(async () => {
    setSavedRecordsLoading(true)
    try {
      const base = import.meta.env.VITE_DASHBOARD_URL || '/api/v1/dashboard'
      const res = await fetch(`${base}/records`, { signal: AbortSignal.timeout(3000) })
      if (res.ok) {
        const j = await res.json()
        setSavedRecords(j.records || [])
      }
    } catch { /* offline */ }
    setSavedRecordsLoading(false)
  }, [])

  useEffect(() => {
    if (active === 'openrecord') {
      loadSavedRecords()
      const iv = setInterval(loadSavedRecords, 8000)
      return () => clearInterval(iv)
    }
  }, [active, loadSavedRecords])

  /* ── Download helper ───────────────────────────────────────────────────── */
  function downloadBlob(blob, filename) {
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
  }

  function normalizedGraph() {
    return {
      nodes: graphData.nodes.map(n => ({ ...n })),
      links: (graphData.links || []).map(l => ({
        source: typeof l.source === 'object' ? l.source.id : l.source,
        target: typeof l.target === 'object' ? l.target.id : l.target,
        type: l.type, basis: l.basis, confidence: l.confidence,
      })),
    }
  }

  /* ── EXPORT: JSON Record ───────────────────────────────────────────────── */
  function handleExportJSON() {
    if (!graphData?.nodes) return
    const record = {
      schema: 'onioneye-record', version: 1, name: 'graph-export',
      exported_at: new Date().toISOString(),
      graph: normalizedGraph(),
    }
    downloadBlob(
      new Blob([JSON.stringify(record, null, 2)], { type: 'application/json' }),
      `onioneye-graph-${Date.now()}.json`
    )
    setExportMsg('✓ JSON Record exported')
    setTimeout(() => setExportMsg(null), 3000)
  }

  /* ── EXPORT: STIX 2.1 bundle ───────────────────────────────────────────── */
  function handleExportSTIX() {
    if (!graphData?.nodes) return
    const bundle = {
      type: 'bundle', id: `bundle--${Date.now()}`, spec_version: '2.1',
      objects: graphData.nodes.map(n => ({
        type: n.type === 'actor' ? 'threat-actor' : 'observed-data',
        id: `${n.type}--${n.id}`,
        created: new Date().toISOString(), modified: new Date().toISOString(),
        name: n.label,
        labels: [n.type, n.flagged ? 'flagged-prohibited' : 'monitored'],
        confidence: Math.round((n.confidence || 0.8) * 100),
        description: n.flagReason || `OnionEye dark web intelligence on ${n.label}`,
      })),
    }
    downloadBlob(
      new Blob([JSON.stringify(bundle, null, 2)], { type: 'application/json' }),
      `onioneye-stix2.1-${Date.now()}.json`
    )
    setExportMsg('✓ STIX 2.1 Bundle exported')
    setTimeout(() => setExportMsg(null), 3000)
  }

  /* ── EXPORT: CSV ───────────────────────────────────────────────────────── */
  function handleExportCSV() {
    if (!graphData?.nodes?.length) return
    const header = 'id,label,type,market,confidence,flagged'
    const rows = graphData.nodes.map(n =>
      `${n.id},"${(n.label||'').replace(/"/g,'""')}",${n.type||''},${n.market||''},${n.confidence||''},${n.flagged||false}`
    )
    downloadBlob(
      new Blob([header + '\n' + rows.join('\n')], { type: 'text/csv' }),
      `onioneye-nodes-${Date.now()}.csv`
    )
    setExportMsg('✓ CSV exported')
    setTimeout(() => setExportMsg(null), 3000)
  }
/* ── OPEN RECORD: file import ──────────────────────────────────────────── */
  function handleFileImport(e) {
    const file = e.target.files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = () => {
      try {
        const json = JSON.parse(reader.result)
        const g = json.graph
        if (!g?.nodes || !Array.isArray(g.nodes)) {
          setImportMsg('✕ Invalid record format — no nodes found')
          return
        }
        setImportedGraph(g)
        setImportedFileName(file.name)
        setImportMsg(`✓ Loaded ${g.nodes.length} nodes, ${(g.links||[]).length} links`)
      } catch {
        setImportMsg('✕ Failed to parse JSON file')
      }
    }
    reader.readAsText(file)
    e.target.value = ''
  }

  function handleApplyImport() {
    if (!importedGraph || !onMergeGraph) return
    onMergeGraph(importedGraph)
    setImportMsg('✓ Merged into current graph')
    setImportedGraph(null)
    setImportedFileName(null)
    setTimeout(() => setImportMsg(null), 3000)
  }

  /* ── Save current graph to backend records ─────────────────────────────── */
  async function handleSaveRecord() {
    if (!graphData?.nodes?.length) return
    const name = prompt('Record name:', `case-${Date.now()}`)
    if (!name) return
    setSaveRecordMsg('Saving…')
    try {
      const base = import.meta.env.VITE_DASHBOARD_URL || '/api/v1/dashboard'
      const res = await fetch(`${base}/records`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, graph: normalizedGraph() }),
        signal: AbortSignal.timeout(5000),
      })
      setSaveRecordMsg(res.ok ? '✓ Record saved' : '✕ Save failed')
      if (res.ok) loadSavedRecords()
    } catch {
      setSaveRecordMsg('✕ Backend unreachable')
    }
    setTimeout(() => setSaveRecordMsg(null), 3000)
  }

  /* ── Load saved record (backend) into graph ────────────────────────────── */
  async function handleLoadSavedRecord(recordId) {
    setLoadMsg('Loading…')
    try {
      const base = import.meta.env.VITE_DASHBOARD_URL || '/api/v1/dashboard'
      const res = await fetch(`${base}/records/${recordId}`, { signal: AbortSignal.timeout(3000) })
      if (!res.ok) { setLoadMsg('✕ Load failed'); return }
      const json = await res.json()
      if (json.graph?.nodes && onMergeGraph) {
        onMergeGraph(json.graph)
        setLoadMsg('✓ Loaded into graph')
      } else {
        setLoadMsg('✕ No graph data in record')
      }
    } catch {
      setLoadMsg('✕ Backend unreachable')
    }
    setTimeout(() => setLoadMsg(null), 3000)
  }

  const btnBase = {
    fontFamily: 'var(--font-mono)', fontSize: 10, fontWeight: 900,
    padding: '8px 10px', borderRadius: 8, cursor: 'pointer',
    border: '1px solid var(--border)', background: 'rgba(0,0,0,0.25)',
    color: 'var(--text-secondary)',
  }
if (collapsed) {
    return (
      <>
        <AiRestartModal
          withAi={aiModal}
          aiRestarting={aiRestarting}
          onConfirm={(w) => { setAiModal(null); onRestartStack?.(w) }}
          onClose={() => setAiModal(null)}
        />
        <aside style={{
          width: 28,
          borderLeft: '1px solid var(--border)',
          display: 'flex', flexDirection: 'column', alignItems: 'center',
          flexShrink: 0, zIndex: 50, background: '#04070c',
        }}>
          <button
            onClick={() => setCollapsed(false)}
            title="Expand TOOLS & RECORDS"
            style={{
              marginTop: 8, width: 24, height: 24, borderRadius: 4,
              border: '1px solid var(--border)', background: 'rgba(5, 8, 13, 0.95)',
              color: 'var(--green)', cursor: 'pointer', fontFamily: 'var(--font-mono)',
              fontSize: 13, fontWeight: 800, lineHeight: 1,
              display: 'flex', alignItems: 'center', justifyContent: 'center',
            }}
          >
            &lt;
          </button>
        </aside>
      </>
    )
  }

  return (
    <>
      <AiRestartModal
        withAi={aiModal}
        aiRestarting={aiRestarting}
        onConfirm={(w) => { setAiModal(null); onRestartStack?.(w) }}
        onClose={() => setAiModal(null)}
      />
      <aside style={{
        width,
        borderLeft: '1px solid var(--border)',
        display: 'flex', flexDirection: 'column',
        flexShrink: 0, zIndex: 50, background: '#04070c',
      }}>
        {/* ── Collapse toggle ─────────────────────────────────────────────── */}
        <div style={{
          display: 'flex', alignItems: 'center', justifyContent: 'flex-end',
          padding: '4px 6px', borderBottom: '1px solid var(--border)',
          background: 'rgba(5, 8, 13, 0.95)',
        }}>
          <button
            onClick={() => setCollapsed(true)}
            title="Collapse TOOLS & RECORDS"
            style={{
              width: 22, height: 22, borderRadius: 4,
              border: '1px solid var(--border)', background: 'transparent',
              color: 'var(--text-secondary)', cursor: 'pointer',
              fontFamily: 'var(--font-mono)', fontSize: 12, fontWeight: 800, lineHeight: 1,
              display: 'flex', alignItems: 'center', justifyContent: 'center',
            }}
          >
            &gt;
          </button>
        </div>

        {/* ── Button strip ───────────────────────────────────────────────── */}
        <div style={{
          display: 'flex', flexDirection: 'column',
          padding: 8, gap: 6,
          borderBottom: '1px solid var(--border)',
          background: 'rgba(5, 8, 13, 0.95)',
          overflowY: 'auto', maxHeight: 420,
        }}>
          {items.map(it => {
            const isActive = active === it.id
            const dotColor =
              it.id === 'search' ? (searchQuery?.trim() ? 'var(--green)' : null) :
              it.id === 'tor' ? (torLive ? 'var(--green)' : 'var(--red)') :
              it.id === 'ai' ? (aiMode ? 'var(--purple)' : null) : null
            return (
              <button
                key={it.id}
                onClick={() => setActive(it.id)}
                style={{
                  width: '100%', borderRadius: 10,
                  border: `1px solid ${isActive ? 'rgba(0,255,157,0.55)' : 'rgba(255,255,255,0.06)'}`,
                  background: isActive ? 'rgba(0,255,157,0.08)' : 'rgba(0,0,0,0.25)',
                  color: isActive ? 'var(--green)' : 'var(--text-secondary)',
                  padding: '10px 8px', cursor: 'pointer',
                  fontFamily: 'var(--font-mono)', fontSize: 10, fontWeight: 900,
                  display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 8,
                }}
                title={it.label}
              >
                <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <span>{it.icon}</span>
                  <span style={{ letterSpacing: '0.03em' }}>{it.label}</span>
                </span>
                {dotColor && (
                  <span style={{
                    width: 6, height: 6, borderRadius: '50%',
                    background: dotColor, boxShadow: `0 0 6px ${dotColor}`,
                    animation: 'pulse-green 1.6s infinite',
                  }} />
                )}
              </button>
            )
          })}
        </div>

        {/* ── Panel area ─────────────────────────────────────────────────── */}
        <div style={{ flex: 1, overflow: 'hidden', padding: 10, display: 'flex', flexDirection: 'column', gap: 10 }}>
{/* ═══ SEARCH PANEL ═════════════════════════════════════════════ */}
          {active === 'search' && (
            <PanelShell title="SEARCH TARGETS">
              <div style={{ position: 'relative' }}>
                <div style={{
                  display: 'flex', alignItems: 'center',
                  background: 'rgba(12, 20, 32, 0.9)', border: '1px solid var(--border-bright)',
                  borderRadius: 4, padding: '4px 10px', gap: 8,
                }}>
                  <span style={{ color: 'var(--green)', fontSize: 11 }}>🔍</span>
                  <input
                    type="text"
                    placeholder="Search Target, PGP, Wallet…"
                    value={searchQuery || ''}
                    onChange={e => setSearchQuery?.(e.target.value)}
                    onFocus={() => searchQuery && setIsSearchOpen(true)}
                    style={{
                      background: 'transparent', border: 'none', outline: 'none',
                      color: 'var(--text-primary)', fontFamily: 'var(--font-mono)',
                      fontSize: 11, width: '100%',
                    }}
                  />
                  {searchQuery && (
                    <button onClick={() => { setSearchQuery?.(''); setIsSearchOpen(false) }}
                      style={{ background: 'transparent', border: 'none', color: 'var(--text-muted)', cursor: 'pointer', fontSize: 11 }}>
                      ✕
                    </button>
                  )}
                </div>

                {isSearchOpen && searchResults.length > 0 && (
                  <div style={{
                    position: 'relative', marginTop: 4,
                    background: 'var(--bg-card)', border: '1px solid var(--border-bright)',
                    borderRadius: 4, boxShadow: '0 8px 30px rgba(0,0,0,0.8)',
                    maxHeight: 260, overflowY: 'auto', zIndex: 10,
                  }}>
                    <div style={{ padding: '4px 8px', fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-muted)', borderBottom: '1px solid var(--border)' }}>
                      MATCHED ({searchResults.length})
                    </div>
                    {searchResults.map(n => (
                      <div
                        key={n.id}
                        onClick={() => { onSelectNode?.(n.id); setIsSearchOpen(false) }}
                        style={{
                          padding: '6px 10px', display: 'flex', alignItems: 'center', justifyContent: 'space-between',
                          borderBottom: '1px solid rgba(255,255,255,0.03)', cursor: 'pointer',
                          fontSize: 11, fontFamily: 'var(--font-mono)',
                        }}
                      >
                        <div>
                          <div style={{ color: 'var(--text-primary)', fontWeight: 600 }}>{n.label}</div>
                          <div style={{ fontSize: 9, color: 'var(--text-muted)' }}>ID: {n.id} · {n.type?.toUpperCase()}</div>
                        </div>
                        {n.flagged && (
                          <span style={{ fontSize: 9, color: 'var(--red)', background: 'rgba(255,23,68,0.1)', padding: '2px 4px', borderRadius: 2 }}>FLAGGED</span>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginTop: 10 }}>
                {graphData?.nodes?.length || 0} nodes · {searchResults.length} matches
              </div>
            </PanelShell>
          )}
{/* ═══ TOR PANEL ═════════════════════════════════════════════════ */}
          {active === 'tor' && (
            <PanelShell title="TOR CONNECTION">
              <div style={{ borderRadius: 8, padding: 10, marginBottom: 10, border: '1px solid var(--border)', background: 'rgba(12,20,32,0.35)' }}>
                <TinyRow k="STATUS" v={torStatus?.status?.toUpperCase() || '…'} color={
                  torStatus?.status === 'connected' ? 'var(--green)' :
                  torStatus?.status === 'connecting' ? 'var(--yellow)' : 'var(--red)'
                } />
                <TinyRow k="LIVE" v={torLive ? 'YES' : 'NO'} color={torLive ? 'var(--green)' : 'var(--red)'} />
                <TinyRow k="SOCKS OK" v={torStatus?.socks_ok ? 'YES' : 'NO'} color={torStatus?.socks_ok ? 'var(--green)' : 'var(--red)'} />
                <TinyRow k="CONTROL OK" v={torStatus?.control_ok ? 'YES' : 'NO'} color={torStatus?.control_ok ? 'var(--green)' : 'var(--red)'} />
              </div>
              {torStatus?.note && (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-secondary)', wordBreak: 'break-word', padding: 8, borderRadius: 6, border: '1px solid var(--border)', background: 'rgba(12,20,32,0.35)' }}>
                  {torStatus.note}
                </div>
              )}
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginTop: 10 }}>
                Tor sidecar · SOCKS5 9050 · Control 9051
              </div>
            </PanelShell>
          )}

          {/* ═══ AI PANEL ══════════════════════════════════════════════════ */}
          {active === 'ai' && (
            <PanelShell title="AI ENGINE (OLLAMA)">
              <div style={{ borderRadius: 8, padding: 10, marginBottom: 10, border: `1px solid ${aiMode ? 'rgba(213,0,249,0.4)' : 'var(--border)'}`, background: aiMode ? 'rgba(213,0,249,0.06)' : 'rgba(12,20,32,0.35)' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                  <span style={{ width: 10, height: 10, borderRadius: '50%', background: aiMode ? 'var(--purple)' : 'var(--text-muted)', boxShadow: aiMode ? '0 0 8px #d500f9' : 'none' }} />
                  <span style={{ fontFamily: 'var(--font-mono)', fontSize: 11, fontWeight: 800, color: aiMode ? 'var(--purple)' : 'var(--text-muted)' }}>
                    {aiMode ? 'AI MODE: ACTIVE' : 'AI MODE: OFFLINE'}
                  </span>
                </div>
                <TinyRow k="MODEL" v={aiMode ? 'qwen2.5:3b + guard3:1b' : '—'} color={aiMode ? 'var(--purple)' : 'var(--text-muted)'} />
                <TinyRow k="ENGINE" v={aiMode ? 'Ollama (running)' : 'Ollama (stopped)'} />
              </div>
              <button onClick={() => setAiModal(true)} disabled={aiRestarting} style={{ ...btnBase, width: '100%', borderColor: 'var(--purple)', color: 'var(--purple)', marginBottom: 6 }}>
                🧠 {aiRestarting ? 'RESTARTING…' : 'RESTART WITH AI'}
              </button>
              <button onClick={() => setAiModal(false)} disabled={aiRestarting} style={{ ...btnBase, width: '100%', borderColor: 'var(--cyan)', color: 'var(--cyan)' }}>
                🧅 {aiRestarting ? 'RESTARTING…' : 'RESTART NO AI'}
              </button>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginTop: 10 }}>
                Ollama inference · Qwen 3B + Llama Guard 3
              </div>
            </PanelShell>
          )}
{/* ═══ TOOL MAP PANEL ════════════════════════════════════════════ */}
          {active === 'toolmap' && (
            <PanelShell title="TOOL ARCHITECTURE MAP">
              {componentsLoading && components.length === 0 && (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>Probing services…</div>
              )}
              {componentSummary && (
                <div style={{ display: 'flex', gap: 8, marginBottom: 10, padding: '6px 8px', borderRadius: 6, border: '1px solid var(--border)', background: 'rgba(12,20,32,0.35)' }}>
                  <span style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 800, color: 'var(--green)' }}>● {componentSummary.up} UP</span>
                  <span style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 800, color: 'var(--red)' }}>● {componentSummary.down} DOWN</span>
                </div>
              )}
              {['ui','api','engine','worker','shield','infra','network','ai'].map(cat => {
                const group = components.filter(c => c.category === cat)
                if (!group.length) return null
                const labels = { ui:'UI LAYER', api:'API LAYER', engine:'ENGINE', worker:'WORKER', shield:'SHIELD', infra:'INFRASTRUCTURE', network:'NETWORK', ai:'AI' }
                return (
                  <div key={cat} style={{ marginBottom: 10 }}>
                    <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, fontWeight: 800, color: 'var(--text-muted)', letterSpacing: '0.1em', marginBottom: 6, paddingBottom: 4, borderBottom: '1px solid rgba(255,255,255,0.03)' }}>
                      {labels[cat]}
                    </div>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 6 }}>
                      {group.map(c => <ToolNode key={c.id} comp={c} />)}
                    </div>
                  </div>
                )
              })}
              {components.length > 0 && (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)', marginTop: 4 }}>
                  Polling every 5s · Docker + TCP probes
                </div>
              )}
            </PanelShell>
          )}
{/* ═══ EXPORT PANEL ═══════════════════════════════════════════════ */}
          {active === 'export' && (
            <PanelShell title="EXPORT GRAPH">
              {exportMsg && (
                <div style={{ border: '1px solid var(--border)', borderRadius: 6, padding: 8, marginBottom: 10, background: 'rgba(12,20,32,0.4)', fontFamily: 'var(--font-mono)', fontSize: 9, color: exportMsg.startsWith('✓') ? 'var(--green)' : 'var(--red)' }}>
                  {exportMsg}
                </div>
              )}
              <button onClick={handleExportJSON} style={{ ...btnBase, width: '100%', marginBottom: 6, borderColor: 'var(--green)', color: 'var(--green)' }}>
                📋 JSON Record
              </button>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)', marginBottom: 10 }}>
                Full graph snapshot · Nodes + Links · Re-importable via Open Record
              </div>
              <button onClick={handleExportSTIX} style={{ ...btnBase, width: '100%', marginBottom: 6, borderColor: 'var(--purple)', color: 'var(--purple)' }}>
                📦 STIX 2.1 Bundle
              </button>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)', marginBottom: 10 }}>
                STIX standard · Compatible with threat-intel platforms
              </div>
              <button onClick={handleExportCSV} style={{ ...btnBase, width: '100%', marginBottom: 6, borderColor: 'var(--cyan)', color: 'var(--cyan)' }}>
                📊 CSV Nodes
              </button>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)', marginBottom: 10 }}>
                Tabular format · For spreadsheet / SIEM import
              </div>
              <button onClick={handleSaveRecord} style={{ ...btnBase, width: '100%', borderColor: 'var(--orange)', color: 'var(--orange)' }}>
                💾 Save to Dashboard
              </button>
              {saveRecordMsg && (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: saveRecordMsg.startsWith('✓') ? 'var(--green)' : 'var(--red)', marginTop: 6 }}>
                  {saveRecordMsg}
                </div>
              )}
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginTop: 10 }}>
                {graphData?.nodes?.length || 0} nodes · {(graphData?.links || []).length} links
              </div>
            </PanelShell>
          )}
{/* ═══ OPEN RECORD PANEL ═════════════════════════════════════════ */}
          {active === 'openrecord' && (
            <PanelShell title="OPEN RECORD">
              <input ref={fileInputRef} type="file" accept=".json" onChange={handleFileImport} style={{ display: 'none' }} />
              {importMsg && (
                <div style={{ border: '1px solid var(--border)', borderRadius: 6, padding: 8, marginBottom: 10, background: 'rgba(12,20,32,0.4)', fontFamily: 'var(--font-mono)', fontSize: 9, color: importMsg.startsWith('✓') ? 'var(--green)' : 'var(--red)', wordBreak: 'break-word' }}>
                  {importMsg}
                </div>
              )}
              <button onClick={() => fileInputRef.current?.click()} style={{ ...btnBase, width: '100%', marginBottom: 6, borderColor: 'var(--cyan)', color: 'var(--cyan)' }}>
                📁 Import JSON File
              </button>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)', marginBottom: 10 }}>
                Open an exported OnionEye JSON Record from disk
              </div>
              {importedGraph && (
                <div style={{ border: '1px solid rgba(0,255,157,0.35)', borderRadius: 6, padding: 8, marginBottom: 10, background: 'rgba(0,255,157,0.04)' }}>
                  <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 800, color: 'var(--green)', marginBottom: 4 }}>STAGED: {importedFileName}</div>
                  <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-secondary)', marginBottom: 8 }}>
                    {importedGraph.nodes.length} nodes · {(importedGraph.links || []).length} links ready to merge
                  </div>
                  <button onClick={handleApplyImport} style={{ ...btnBase, width: '100%', borderColor: 'var(--green)', color: 'var(--green)' }}>
                    ⚡ APPLY TO GRAPH
                  </button>
                </div>
              )}
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, fontWeight: 800, color: 'var(--text-muted)', letterSpacing: '0.1em', marginBottom: 6, marginTop: 4 }}>
                SAVED RECORDS (DASHBOARD)
              </div>
              {savedRecordsLoading && <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>Loading…</div>}
              {!savedRecordsLoading && savedRecords.length === 0 && <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>No saved records yet</div>}
              {savedRecords.map(r => (
                <div key={r.id} style={{ borderRadius: 6, padding: '6px 8px', marginBottom: 6, border: '1px solid var(--border)', background: 'rgba(12,20,32,0.35)', display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 6 }}>
                  <div style={{ minWidth: 0 }}>
                    <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 800, color: 'var(--text-primary)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{r.name}</div>
                    <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)' }}>{r.node_count || '?'} nodes · {r.link_count || '?'} links</div>
                  </div>
                  <button onClick={() => handleLoadSavedRecord(r.id)} style={{ ...btnBase, padding: '4px 8px', fontSize: 9, borderColor: 'var(--green)', color: 'var(--green)', flexShrink: 0 }}>LOAD</button>
                </div>
              ))}
              {loadMsg && (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: loadMsg.startsWith('✓') ? 'var(--green)' : 'var(--red)', marginTop: 4 }}>
                  {loadMsg}
                </div>
              )}
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginTop: 10 }}>
                Imported nodes merge into the existing graph (deduplicated by ID)
              </div>
            </PanelShell>
          )}

        </div>
      </aside>
    </>
  )
}
