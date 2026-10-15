import React, { useMemo, useState, useEffect, useRef } from 'react'
import LiveLogs from './LiveLogs'
import {
  fetchNotifications,
  fetchNegativeControl,
  saveNegativeControlRule,
  deleteNegativeControlRule,
  fetchArtifacts,
  addNegativeControlArtifact,
  fetchManualReviewQueue,
  dismissReviewItem,
  nukeAll,
  deleteAll,
} from '../services/api'

function PanelShell({ title, children }) {
  return (
    <div
      style={{
        flex: 1,
        overflow: 'hidden',
        background: 'linear-gradient(180deg, #070c14 0%, #04070c 100%)',
        border: '1px solid var(--border)',
        borderRadius: 6,
        display: 'flex',
        flexDirection: 'column',
        minHeight: 0,
      }}
    >
      <div
        style={{
          padding: '10px 10px',
          borderBottom: '1px solid var(--border)',
          background: 'rgba(5, 8, 13, 0.95)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
        }}
      >
        <div style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <span style={{ fontFamily: 'var(--font-mono)', fontSize: 10, color: 'var(--green)', fontWeight: 800 }}>
            {title}
          </span>
          <span style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)' }}>
            SLIDING WINDOW
          </span>
        </div>
      </div>
      <div style={{ flex: 1, overflowY: 'auto', padding: 10 }}>{children}</div>
    </div>
  )
}

function TinyRow({ k, v, color }) {
  return (
    <div style={{ display: 'flex', justifyContent: 'space-between', gap: 10, marginBottom: 6 }}>
      <span style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>{k}</span>
      <span style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 800, color: color || 'var(--text-primary)', textAlign: 'right', wordBreak: 'break-word' }}>
        {v}
      </span>
    </div>
  )
}

const inputStyle = {
  width: '100%',
  background: 'rgba(3, 6, 11, 0.85)',
  border: '1px solid var(--border)',
  borderRadius: 5,
  color: 'var(--text-primary)',
  fontFamily: 'var(--font-mono)',
  fontSize: 10,
  padding: '6px 8px',
  outline: 'none',
}

const btnBase = {
  padding: '8px 10px',
  borderRadius: 6,
  border: '1px solid var(--border)',
  background: 'rgba(7, 12, 20, 0.8)',
  color: 'var(--text-secondary)',
  cursor: 'pointer',
  fontFamily: 'var(--font-mono)',
  fontSize: 10,
  fontWeight: 800,
}

function SliderBlock({ title, count, color, children, defaultOpen = true }) {
  const [open, setOpen] = useState(defaultOpen)
  return (
    <div style={{ border: '1px solid var(--border)', borderRadius: 6, marginBottom: 8, background: 'rgba(12, 20, 32, 0.35)', overflow: 'hidden' }}>
      <button
        onClick={() => setOpen(o => !o)}
        style={{
          width: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: 8,
          padding: '7px 9px',
          background: 'rgba(5, 8, 13, 0.9)',
          border: 'none',
          borderBottom: open ? '1px solid var(--border)' : 'none',
          cursor: 'pointer',
          fontFamily: 'var(--font-mono)',
        }}
      >
        <span style={{ fontSize: 9, fontWeight: 900, color: color || 'var(--cyan)', letterSpacing: '0.04em' }}>{title}</span>
        <span style={{ fontSize: 9, color: 'var(--text-muted)', fontWeight: 700 }}>{count} {open ? '▾' : '▸'}</span>
      </button>
      {open && <div style={{ padding: 8 }}>{children}</div>}
    </div>
  )
}

export default function LeftControlRail({ onOpenScanModal, graphData, media, onDeleteAll }) {
  const [active, setActive] = useState('dispatch')
  const [collapsed, setCollapsed] = useState(false)
  const width = 290

  const [notifItems, setNotifItems] = useState([])
  const [unreadCount, setUnreadCount] = useState(0)
  const seenNotifRef = useRef(new Set())

  // Always poll notifications in the background (independent of active tab) so
  // the badge count and panel data are fresh the moment the user opens it.
  useEffect(() => {
    let alive = true
    async function tick() {
      try {
        const res = await fetchNotifications()
        if (!alive) return
        if (res?.live && Array.isArray(res.data?.items)) {
          const items = res.data.items.slice(0, 30)
          setNotifItems(items)
          const fresh = items.filter(n => !seenNotifRef.current.has(n.anomaly_id || n.rotation_id))
          if (fresh.length > 0) {
            fresh.forEach(n => seenNotifRef.current.add(n.anomaly_id || n.rotation_id))
            setUnreadCount(c => (active === 'notifications' ? 0 : c + fresh.length))
          } else if (active === 'notifications') {
            setUnreadCount(0)
          }
        }
      } catch {
        // ignore
      }
    }

    tick()
    const interval = setInterval(tick, 3000)
    return () => {
      alive = false
      clearInterval(interval)
    }
  }, [active])

  // ── DATA FETCHED: categorize graph nodes + media into slider buckets ────────
  const fetched = useMemo(() => {
    const nodes = graphData?.nodes || []
    const ipRe = /^(\d{1,3}\.){3}\d{1,3}$/
    const actors = []
    const ips = []
    const onionHosts = []
    const endpoints = []
    const walletsPgp = []
    const unknowns = []

    nodes.forEach(n => {
      const t = n.type || ''
      const label = String(n.label || n.id || '').trim()
      const hostPart = label.split('/')[0].trim()
      if (t === 'actor') actors.push(n)
      else if (t === 'wallet' || t === 'pgp') walletsPgp.push(n)
      else if (t === 'onion' || label.includes('.onion')) onionHosts.push(n)
      else if (t === 'infra' && ipRe.test(hostPart)) ips.push(n)
      else if (t === 'infra' || label.includes('http') || label.includes('/') || (n.id && String(n.id).startsWith('id_'))) endpoints.push(n)
      else unknowns.push(n)
    })

    return {
      actors,
      ips,
      onionHosts,
      endpoints,
      walletsPgp,
      unknowns,
      media: Array.isArray(media) ? media : [],
      totalNodes: nodes.length,
    }
  }, [graphData, media])

  // ── NEGATIVE CONTROL ─────────────────────────────────────────────────────────
  const [ncMode, setNcMode] = useState('FREQUENCY')
  const [ncType, setNcType] = useState('onion_link')
  const [ncRules, setNcRules] = useState([])
  const [ncArtifacts, setNcArtifacts] = useState([])
  const [ncCap, setNcCap] = useState(3)
  const [ncNewValue, setNcNewValue] = useState('')
  const [ncBusy, setNcBusy] = useState(false)

  const ncTypeOptions = [
    { value: 'email', label: 'EMAIL' },
    { value: 'btc_address', label: 'BTC ADDRESS' },
    { value: 'xmr_address', label: 'XMR ADDRESS' },
    { value: 'pgp_fingerprint', label: 'PGP FINGERPRINT' },
    { value: 'handle', label: 'HANDLE' },
    { value: 'onion_link', label: 'ONION LINK' },
  ]

  useEffect(() => {
    if (active !== 'negcontrol') return
    let alive = true
    async function load() {
      const [rulesRes, artRes] = await Promise.all([
        fetchNegativeControl(),
        fetchArtifacts(ncType, ncMode),
      ])
      if (!alive) return
      if (rulesRes.live && rulesRes.data?.rules) setNcRules(rulesRes.data.rules)
      if (artRes.live && artRes.data?.artifacts) setNcArtifacts(artRes.data.artifacts)
    }
    load()
    const interval = setInterval(load, 8000)
    return () => {
      alive = false
      clearInterval(interval)
    }
  }, [active]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (active !== 'negcontrol') return
    let alive = true
    async function loadArtifacts() {
      const artRes = await fetchArtifacts(ncType, ncMode)
      if (alive && artRes.live && artRes.data?.artifacts) setNcArtifacts(artRes.data.artifacts)
    }
    loadArtifacts()
    return () => {
      alive = false
    }
  }, [active, ncType, ncMode]) // eslint-disable-line react-hooks/exhaustive-deps

  async function handleSaveNcRule(cap) {
    setNcBusy(true)
    const res = await saveNegativeControlRule({
      mode: ncMode,
      artifact_type: ncType,
      frequency_cap: cap,
      enabled: true,
    })
    setNcBusy(false)
    if (res.ok) {
      const rulesRes = await fetchNegativeControl()
      if (rulesRes.live && rulesRes.data?.rules) setNcRules(rulesRes.data.rules)
    } else {
      alert(`Failed to save rule: ${res.error || 'unknown error'}`)
    }
  }

  async function handleDeleteRule(id) {
    await deleteNegativeControlRule(id)
    const rulesRes = await fetchNegativeControl()
    if (rulesRes.live && rulesRes.data?.rules) setNcRules(rulesRes.data.rules)
  }

  async function handleAddArtifact() {
    const value = ncNewValue.trim()
    if (!value) return
    const res = await addNegativeControlArtifact({ artifact_type: ncType, value })
    if (res.ok) {
      setNcNewValue('')
      const artRes = await fetchArtifacts(ncType, ncMode)
      if (artRes.live && artRes.data?.artifacts) setNcArtifacts(artRes.data.artifacts)
    } else {
      alert(`Failed to add artifact: ${res.error || 'unknown error'}`)
    }
  }

  // ── MANUAL REVIEW (re-homed from Header modal) ───────────────────────────────
  const [reviewItems, setReviewItems] = useState([])
  const [reviewBusy, setReviewBusy] = useState(false)

  useEffect(() => {
    if (active !== 'manual') return
    let alive = true
    async function load() {
      const res = await fetchManualReviewQueue()
      if (!alive) return
      if (res.live && res.data?.items) setReviewItems(res.data.items.filter(it => !it.expired))
    }
    load()
    const interval = setInterval(load, 10000)
    return () => {
      alive = false
      clearInterval(interval)
    }
  }, [active])

  async function handleDismissReview(id) {
    setReviewBusy(true)
    const res = await dismissReviewItem(id)
    setReviewBusy(false)
    if (res.ok) {
      setReviewItems(prev => prev.filter(it => it.id !== id))
    } else {
      alert(`Failed to dismiss: ${res.error || 'Unknown error'}`)
    }
  }

  // ── DELETE ALL ────────────────────────────────────────────────────────────────
  const [deleteArmed, setDeleteArmed] = useState(false)
  const [deleteBusy, setDeleteBusy] = useState(false)
  const [deleteMsg, setDeleteMsg] = useState('')

  async function handleDeleteAll() {
    setDeleteBusy(true)
    setDeleteMsg('Deleting all correlation data…')
    try {
      await deleteAll() // throws on HTTP error
      setDeleteMsg('✓ All correlation data deleted.')
    } catch (e) {
      setDeleteMsg(`✕ Delete failed: ${e?.message || 'unknown error'}`)
      setDeleteBusy(false)
      return
    }
    setDeleteBusy(false)
    // Refresh the live graph after wipe.
    if (onDeleteAll) onDeleteAll()
  }

  // ── NUKE ALL ──────────────────────────────────────────────────────────────────
  const [nukeArmed, setNukeArmed] = useState(false)
  const [nukeBusy, setNukeBusy] = useState(false)
  const [nukeMsg, setNukeMsg] = useState('')

  async function handleNukeAll() {
    setNukeBusy(true)
    setDeleteBusy(true)
    setNukeMsg('Step 1/2 — deleting persisting DB data…')
    try {
      await deleteAll()
    } catch (e) {
      setNukeMsg(`✕ Data wipe failed: ${e?.message || 'unknown error'} — aborting.`)
      setNukeBusy(false)
      setDeleteBusy(false)
      setNukeArmed(false)
      return
    }
    setNukeMsg('Step 2/2 — docker compose down -v && up -d…')
    const res = await nukeAll()
    setNukeBusy(false)
    setDeleteBusy(false)
    if (res.ok) {
      setNukeMsg(res.data?.message || '✓ NUKE dispatched — stack down, volumes wiped, restarting.')
    } else {
      setNukeMsg(`✕ NUKE failed: ${res.error || 'unknown error'}`)
    }
    if (onDeleteAll) onDeleteAll()
  }


  const items = useMemo(
    () => [
      { id: 'dispatch', label: 'DISPATCH', icon: '⚡' },
      { id: 'data', label: 'DATA FETCHED', icon: '📥' },
      { id: 'logs', label: 'LOG TABLE', icon: '🖥' },
      { id: 'negcontrol', label: 'NEGATIVE CONTROL', icon: '🧯' },
      { id: 'manual', label: 'MANUAL REVIEW', icon: '🔍' },
      { id: 'notifications', label: 'NOTIFICATIONS', icon: '🔔' },
      { id: 'guide', label: 'GUIDE', icon: '📚' },
      { id: 'deleteAll', label: 'DELETE ALL', icon: '🗑️' },
      { id: 'nukeAll', label: 'NUKE ALL', icon: '💥' },
    ],
    []
  )

  if (collapsed) {
    return (
      <aside
        style={{
          width: 28,
          borderRight: '1px solid var(--border)',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          flexShrink: 0,
          zIndex: 50,
          background: '#04070c',
        }}
      >
        <button
          onClick={() => setCollapsed(false)}
          title="Expand CONTROL RAIL"
          style={{
            marginTop: 8,
            width: 24,
            height: 24,
            borderRadius: 4,
            border: '1px solid var(--border)',
            background: 'rgba(5, 8, 13, 0.95)',
            color: 'var(--green)',
            cursor: 'pointer',
            fontFamily: 'var(--font-mono)',
            fontSize: 13,
            fontWeight: 800,
            lineHeight: 1,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          &gt;
        </button>
      </aside>
    )
  }

  return (
    <aside
      style={{
        width,
        borderRight: '1px solid var(--border)',
        display: 'flex',
        flexDirection: 'column',
        flexShrink: 0,
        zIndex: 50,
        background: '#04070c',
      }}
    >
      {/* ── Collapse toggle ─────────────────────────────────────────────────── */}
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'flex-end',
          padding: '4px 6px',
          borderBottom: '1px solid var(--border)',
          background: 'rgba(5, 8, 13, 0.95)',
        }}
      >
        <button
          onClick={() => setCollapsed(true)}
          title="Collapse CONTROL RAIL"
          style={{
            width: 22,
            height: 22,
            borderRadius: 4,
            border: '1px solid var(--border)',
            background: 'transparent',
            color: 'var(--text-secondary)',
            cursor: 'pointer',
            fontFamily: 'var(--font-mono)',
            fontSize: 12,
            fontWeight: 800,
            lineHeight: 1,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          &lt;
        </button>
      </div>

      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          padding: 8,
          gap: 6,
          borderBottom: '1px solid var(--border)',
          background: 'rgba(5, 8, 13, 0.95)',
          overflowY: 'auto',
          maxHeight: 420,
        }}
      >
        {items.map(it => {
          const isActive = active === it.id
          return (
            <button
              key={it.id}
              onClick={() => setActive(it.id)}
              style={{
                width: '100%',
                borderRadius: 10,
                border: `1px solid ${isActive ? 'rgba(0,255,157,0.55)' : 'rgba(255,255,255,0.06)'}`,
                background: isActive ? 'rgba(0,255,157,0.08)' : 'rgba(0,0,0,0.25)',
                color: isActive ? 'var(--green)' : 'var(--text-secondary)',
                padding: '10px 8px',
                cursor: 'pointer',
                fontFamily: 'var(--font-mono)',
                fontSize: 10,
                fontWeight: 900,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                gap: 8,
              }}
              title={it.label}
            >
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <span>{it.icon}</span>
                <span style={{ letterSpacing: '0.03em' }}>{it.label}</span>
              </span>
              {it.id === 'notifications' && unreadCount > 0 && (
                <span
                  style={{
                    background: 'var(--red)',
                    color: '#fff',
                    fontSize: 8,
                    fontWeight: 900,
                    borderRadius: 10,
                    padding: '1px 6px',
                    minWidth: 16,
                    textAlign: 'center',
                  }}
                >
                  {unreadCount}
                </span>
              )}
            </button>
          )
        })}
      </div>

      <div style={{ flex: 1, overflow: 'hidden', padding: 10, display: 'flex', flexDirection: 'column', gap: 10 }}>
        {active === 'dispatch' && (
          <PanelShell title="DISPATCH INVESTIGATION">
            <button
              onClick={onOpenScanModal}
              style={{
                padding: '10px 12px',
                borderRadius: 8,
                border: '1px dashed var(--green)',
                background: 'rgba(0,255,157,0.08)',
                color: 'var(--green)',
                cursor: 'pointer',
                fontFamily: 'var(--font-mono)',
                fontWeight: 900,
                fontSize: 11,
              }}
            >
              ⚡ START SCAN
            </button>
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginTop: 10 }}>
              Uses existing scan modal.
            </div>
          </PanelShell>
        )}

        {active === 'data' && (
          <PanelShell title={`DATA FETCHED · ${fetched.totalNodes} NODES`}>
            <SliderBlock title="ACTORS" count={fetched.actors.length} color="var(--green)">
              {fetched.actors.length === 0 ? (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>no actors</div>
              ) : (
                fetched.actors.map(a => (
                  <div key={a.id} style={{ padding: '4px 0', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                    <TinyRow k="ID" v={String(a.id).slice(0, 8) + '…'} color="var(--green)" />
                    <TinyRow k="HANDLE" v={String(a.label || a.id)} color="var(--text-primary)" />
                  </div>
                ))
              )}
            </SliderBlock>

            <SliderBlock title="IPS" count={fetched.ips.length} color="var(--red)">
              {fetched.ips.length === 0 ? (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>no IP indicators</div>
              ) : (
                fetched.ips.map(n => (
                  <div key={n.id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)', padding: '4px 0' }}>
                    <TinyRow k="IP" v={String(n.label || n.id)} color="var(--red)" />
                    {n.actor_id && <TinyRow k="ACTOR" v={n.actor_id.slice(0, 8) + '…'} color="var(--text-muted)" />}
                  </div>
                ))
              )}
            </SliderBlock>

            <SliderBlock title="ENDPOINTS" count={fetched.endpoints.length} color="var(--purple)">
              {fetched.endpoints.length === 0 ? (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>no endpoint paths</div>
              ) : (
                fetched.endpoints.map(n => (
                  <div key={n.id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)', padding: '4px 0' }}>
                    <TinyRow k="PATH" v={String(n.label || n.id)} color="var(--magenta)" />
                  </div>
                ))
              )}
            </SliderBlock>

            <SliderBlock title="WALLETS / PGP" count={fetched.walletsPgp.length} color="var(--orange)">
              {fetched.walletsPgp.length === 0 ? (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>no wallet/PGP keys</div>
              ) : (
                fetched.walletsPgp.map(n => (
                  <div key={n.id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)', padding: '4px 0' }}>
                    <TinyRow k={n.type === 'wallet' ? 'WALLET' : 'PGP'} v={String(n.label || n.id)} color="var(--orange)" />
                  </div>
                ))
              )}
            </SliderBlock>

            <SliderBlock title="ONION HOSTS" count={fetched.onionHosts.length} color="var(--cyan)">
              {fetched.onionHosts.length === 0 ? (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>no onion hosts</div>
              ) : (
                fetched.onionHosts.map(n => (
                  <div key={n.id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)', padding: '4px 0' }}>
                    <TinyRow k="ONION" v={String(n.label || n.id)} color="var(--cyan)" />
                  </div>
                ))
              )}
            </SliderBlock>

            <SliderBlock title="IMAGES" count={fetched.media.length} color="var(--yellow)">
              {fetched.media.length === 0 ? (
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>no extracted media</div>
              ) : (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 6 }}>
                  {fetched.media.slice(0, 30).map((m, i) => (
                    <img
                      key={m.src || i}
                      src={m.src}
                      alt="evidence"
                      style={{
                        width: '100%',
                        aspectRatio: '1',
                        objectFit: 'cover',
                        borderRadius: 4,
                        border: '1px solid var(--border)',
                        background: '#000',
                      }}
                    />
                  ))}
                </div>
              )}
            </SliderBlock>

            {fetched.unknowns.length > 0 && (
              <SliderBlock title="UNCATEGORIZED" count={fetched.unknowns.length} color="var(--text-muted)">
                {fetched.unknowns.map(n => (
                  <TinyRow key={n.id} k="?" v={String(n.label || n.id)} color="var(--text-muted)" />
                ))}
              </SliderBlock>
            )}
          </PanelShell>
        )}

        {active === 'logs' && (
          <PanelShell title="LOG TABLE">
            <LiveLogs />
          </PanelShell>
        )}

        {active === 'negcontrol' && (
          <PanelShell title={`NEGATIVE CONTROL · ${ncMode}`}>
            {/* Mode toggle: FREQUENCY (occurrence dampening) vs STATIC (denylist) */}
            <div style={{ display: 'flex', gap: 6, marginBottom: 8 }}>
              <button
                onClick={() => setNcMode('FREQUENCY')}
                style={{
                  ...btnBase,
                  flex: 1,
                  borderColor: ncMode === 'FREQUENCY' ? 'var(--green)' : 'var(--border)',
                  color: ncMode === 'FREQUENCY' ? 'var(--green)' : 'var(--text-muted)',
                  background: ncMode === 'FREQUENCY' ? 'rgba(0,255,157,0.08)' : 'transparent',
                }}
              >
                FREQUENCY
              </button>
              <button
                onClick={() => setNcMode('STATIC')}
                style={{
                  ...btnBase,
                  flex: 1,
                  borderColor: ncMode === 'STATIC' ? 'var(--red)' : 'var(--border)',
                  color: ncMode === 'STATIC' ? 'var(--red)' : 'var(--text-muted)',
                  background: ncMode === 'STATIC' ? 'rgba(255,23,68,0.08)' : 'transparent',
                }}
              >
                STATIC
              </button>
            </div>

            {/* Artifact type selector */}
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginBottom: 4 }}>
              ARTIFACT TYPE
            </div>
            <select
              value={ncType}
              onChange={e => setNcType(e.target.value)}
              style={{ ...inputStyle, marginBottom: 10 }}
            >
              {ncTypeOptions.map(o => (
                <option key={o.value} value={o.value}>{o.label}</option>
              ))}
            </select>

            {/* Rule controls */}
            <div style={{ border: '1px solid var(--border)', borderRadius: 6, padding: 8, marginBottom: 10, background: 'rgba(12,20,32,0.4)' }}>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginBottom: 6 }}>{ncMode === 'FREQUENCY' ? 'FREQUENCY CAP' : 'STATIC MODE (cap ignored)'}</div>
              {ncMode === 'FREQUENCY' ? (
                <div style={{ display: 'flex', gap: 6, marginBottom: 6 }}>
                  <input
                    type="number"
                    min={0}
                    value={ncCap}
                    onChange={e => setNcCap(Number(e.target.value))}
                    style={{ ...inputStyle, flex: 1 }}
                  />
                  <button onClick={() => handleSaveNcRule(ncCap)} style={{ ...btnBase, borderColor: 'var(--green)', color: 'var(--green)', whiteSpace: 'nowrap' }} disabled={ncBusy}>
                    {ncBusy ? 'SAVING…' : 'SAVE RULE'}
                  </button>
                </div>
              ) : (
                <button onClick={() => handleSaveNcRule(0)} style={{ ...btnBase, borderColor: 'var(--red)', color: 'var(--red)', width: '100%' }} disabled={ncBusy}>
                  {ncBusy ? 'SAVING…' : 'ADD STATIC RULE'}
                </button>
              )}
            </div>

            {/* Existing rules */}
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginBottom: 4 }}>
              ACTIVE RULES ({ncRules.length})
            </div>
            {ncRules.length === 0 ? (
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginBottom: 10 }}>no rules yet</div>
            ) : (
              ncRules.slice(0, 8).map(r => (
                <div key={r.rule_id} style={{ border: '1px solid var(--border)', borderRadius: 6, padding: 8, marginBottom: 8, background: 'rgba(12,20,32,0.35)' }}>
                  <TinyRow k="MODE" v={r.mode} color={r.mode === 'FREQUENCY' ? 'var(--green)' : 'var(--red)'} />
                  <TinyRow k="TYPE" v={r.artifact_type} color="var(--cyan)" />
                  <TinyRow k="CAP" v={String(r.frequency_cap)} color="var(--text-primary)" />
                  <TinyRow k="ARTIFACTS" v={String(r.artifact_count || 0)} color="var(--text-muted)" />
                  <button
                    onClick={() => handleDeleteRule(r.rule_id)}
                    style={{ ...btnBase, borderColor: 'var(--red)', color: 'var(--red)', width: '100%', marginTop: 4, padding: '5px 8px' }}
                  >
                    REMOVE RULE
                  </button>
                </div>
              ))
            )}

            {/* Artifact list for selected type + mode */}
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginBottom: 4, marginTop: 4 }}>
              ARTIFACTS · {ncType.toUpperCase()} ({ncArtifacts.length})
            </div>
            {ncArtifacts.length === 0 ? (
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginBottom: 10 }}>no artifacts for this type/mode</div>
            ) : (
              ncArtifacts.slice(0, 12).map(a => (
                <div key={a.artifact_id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)', padding: '5px 0' }}>
                  <TinyRow k="VALUE" v={a.value} color="var(--text-primary)" />
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <TinyRow k="OCC" v={String(a.occurrence_count)} color="var(--cyan)" />
                    <TinyRow k="SEEN" v={String(a.last_seen_at || '').slice(0, 16).replace('T', ' ')} color="var(--text-muted)" />
                  </div>
                </div>
              ))
            )}

            {/* Add artifact */}
            <div style={{ display: 'flex', gap: 6, marginTop: 10 }}>
              <input
                value={ncNewValue}
                onChange={e => setNcNewValue(e.target.value)}
                placeholder={`add ${ncType} value…`}
                style={{ ...inputStyle, flex: 1 }}
                onKeyDown={e => e.key === 'Enter' && handleAddArtifact()}
              />
              <button onClick={handleAddArtifact} style={{ ...btnBase, borderColor: 'var(--green)', color: 'var(--green)' }}>
                ADD
              </button>
            </div>
          </PanelShell>
        )}

        {active === 'notifications' && (
          <PanelShell title="NOTIFICATIONS">
            {notifItems.length === 0 ? (
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)' }}>
                Waiting for events…
              </div>
            ) : (
              notifItems.map((n, idx) => {
                const kind = n.event === 'worker-anomaly' ? 'var(--red)' : 'var(--green)'
                const id = n.anomaly_id || n.rotation_id || '—'
                return (
                  <div
                    key={(n.anomaly_id || n.rotation_id || idx) + String(n.event || 'evt')}
                    style={{
                      background: 'rgba(12, 20, 32, 0.45)',
                      border: '1px solid var(--border)',
                      borderRadius: 6,
                      padding: 10,
                      marginBottom: 10,
                    }}
                  >
                    <TinyRow k="TYPE" v={String(n.event || '—').replace(/-/g, ' ')} color={kind} />
                    <TinyRow k="ID" v={String(id)} color={'var(--cyan)'} />
                    <TinyRow k="REASON" v={String(n.reason || 'no reason')} color={'var(--text-primary)'} />
                    <TinyRow
                      k="AT"
                      v={n.at ? String(n.at).slice(0, 19).replace('T', ' ') : '—'}
                      color={'var(--text-muted)'}
                    />
                  </div>
                )
              })
            )}
          </PanelShell>
        )}

        {active === 'manual' && (
          <PanelShell title={`MANUAL REVIEW · ${reviewItems.length} PENDING`}>
            {reviewItems.length === 0 ? (
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', lineHeight: 1.6 }}>
                Queue clear — no CAPTCHA / bot-block items awaiting review.
              </div>
            ) : (
              reviewItems.map(it => (
                <div key={it.id} style={{ border: '1px solid var(--border)', borderRadius: 6, padding: 10, marginBottom: 10, background: 'rgba(12,20,32,0.4)' }}>
                  <TinyRow k="JOB" v={String(it.job_id || it.id)} color="var(--cyan)" />
                  <TinyRow k="TARGET" v={String(it.target || '—')} color="var(--text-primary)" />
                  <TinyRow k="REASON" v={String(it.reason || 'no reason')} color="var(--yellow)" />
                  <TinyRow k="FLAGGED" v={String(it.flagged_at || '').slice(0, 16).replace('T', ' ')} color="var(--text-muted)" />
                  <button
                    onClick={() => handleDismissReview(it.id)}
                    disabled={reviewBusy}
                    style={{ ...btnBase, borderColor: 'var(--green)', color: 'var(--green)', width: '100%', marginTop: 6 }}
                  >
                    {reviewBusy ? 'DISMISSING…' : '✓ RESOLVED — REMOVE'}
                  </button>
                </div>
              ))
            )}
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', marginTop: 6 }}>
              Source: crawl:manual Redis stream · 10s refresh
            </div>
          </PanelShell>
        )}

        {active === 'guide' && (
          <PanelShell title="GUIDE · GITHUB">
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', lineHeight: 1.7, marginBottom: 8 }}>
              OnionEye — operations reference. All artifacts are live DB-backed; there is no mock data.
            </div>
            <div style={{ border: '1px solid var(--border)', borderRadius: 6, padding: 8, marginBottom: 10, background: 'rgba(12,20,32,0.35)' }}>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 900, color: 'var(--green)', marginBottom: 4 }}>
                📂 REPOSITORY
              </div>
              <a
                href="https://github.com/Sudo-Problems/OnionEye-"
                target="_blank"
                rel="noopener noreferrer"
                style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--cyan)', wordBreak: 'break-all', textDecoration: 'none' }}
              >
                github.com/Sudo-Problems/OnionEye-
              </a>
            </div>
            <div style={{ border: '1px solid var(--border)', borderRadius: 6, padding: 8, marginBottom: 10, background: 'rgba(12,20,32,0.35)' }}>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 900, color: 'var(--green)', marginBottom: 4 }}>
                📖 GUIDES / WIKI
              </div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                {[
                  { label: 'README — platform architecture & setup', href: 'https://github.com/Sudo-Problems/OnionEye-#readme' },
                  { label: '📚 Docs — index (all guides & architecture)', href: 'https://github.com/Sudo-Problems/OnionEye-/blob/main/docs/README.md' },
                  { label: 'GitHub Wiki — operations & hardening', href: 'https://github.com/Sudo-Problems/OnionEye-/wiki' },
                  { label: 'Issues — report a bug', href: 'https://github.com/Sudo-Problems/OnionEye-/issues' },
                  { label: 'PRs — contribute', href: 'https://github.com/Sudo-Problems/OnionEye-/pulls' },
                ].map(l => (
                  <a
                    key={l.href}
                    href={l.href}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--cyan)', textDecoration: 'none', borderBottom: '1px solid rgba(0,255,157,0.25)', paddingBottom: 2 }}
                  >
                    ↗ {l.label}
                  </a>
                ))}
              </div>
            </div>
            {[
              { label: '01 · DISPATCH', text: 'START SCAN opens the dispatch modal for a new crawl job.', doc: 'guides/tool-usage.md' },
              { label: '02 · DATA FETCHED', text: 'Every node in the correlation graph bucketed by type: actors, IPs, endpoints, wallets/PGP, onion hosts, images.', doc: 'architecture/decided/subsection-3.md' },
              { label: '03 · LOG TABLE', text: 'Live docker logs of compose services (1s refresh).', doc: 'guides/troubleshooting.md' },
              { label: '04 · NEGATIVE CONTROL', text: 'FREQUENCY mode damps recurring artifacts above the cap; STATIC mode is a hard denylist. Pick artifact type, set a cap, add values.', doc: 'guides/tool-usage.md' },
              { label: '05 · MANUAL REVIEW', text: 'CAPTCHA / bot-block triage items streamed from crawl:manual.', doc: 'guides/tool-usage.md' },
              { label: '06 · NOTIFICATIONS', text: 'Worker-rotation and worker-anomaly events (3s poll).', doc: 'architecture/decided/subsection-3.md' },
              { label: '07 · DELETE ALL', text: 'Wipes correlation graph data (persona links + dependent rows) via storage-api.', doc: 'guides/tool-usage.md' },
              { label: '08 · NUKE ALL', text: 'Wipes DB data, then docker compose down -v && up -d — destroys every volume and restarts the whole stack.', doc: 'guides/tool-usage.md' },
            ].map(sec => (
              <div key={sec.label} style={{ border: '1px solid var(--border)', borderRadius: 6, padding: 8, marginBottom: 8, background: 'rgba(12,20,32,0.35)' }}>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, fontWeight: 900, color: 'var(--green)', marginBottom: 4 }}>{sec.label}</div>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--text-muted)', lineHeight: 1.6 }}>{sec.text}</div>
                <a
                  href={`https://github.com/Sudo-Problems/OnionEye-/blob/main/docs/${sec.doc}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  style={{ fontFamily: 'var(--font-mono)', fontSize: 8, color: 'var(--cyan)', textDecoration: 'none', borderBottom: '1px solid rgba(0,255,157,0.25)', paddingBottom: 2, display: 'inline-flex', gap: 4, marginTop: 6, alignItems: 'center' }}
                >
                  <span>↗</span>
                  <span>doc: {sec.doc}</span>
                </a>
              </div>
            ))}
          </PanelShell>
        )}

        {active === 'deleteAll' && (
          <PanelShell title={`DELETE ALL${deleteBusy ? ' · WORKING…' : ''}`}>
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', lineHeight: 1.6, marginBottom: 10 }}>
              Wipes all persona links, identifier rows, and dependent graph data. Actor rows are preserved.
            </div>
            {deleteMsg && (
              <div style={{
                border: '1px solid var(--border)',
                borderRadius: 6,
                padding: 10,
                marginBottom: 10,
                background: 'rgba(12,20,32,0.4)',
                fontFamily: 'var(--font-mono)',
                fontSize: 9,
                color: deleteMsg.startsWith('✓') ? 'var(--green)' : deleteMsg.startsWith('✕') ? 'var(--red)' : 'var(--yellow)',
                lineHeight: 1.6,
              }}>
                {deleteMsg}
              </div>
            )}
            {!deleteArmed ? (
              <button
                onClick={() => setDeleteArmed(true)}
                disabled={deleteBusy}
                style={{ ...btnBase, borderColor: 'var(--orange)', color: 'var(--orange)', width: '100%' }}
              >
                {deleteBusy ? 'DELETING…' : '⚠ ARM DELETE'}
              </button>
            ) : (
              <div>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--orange)', marginBottom: 8, fontWeight: 800 }}>
                  CONFIRM: delete all correlation data?
                </div>
                <div style={{ display: 'flex', gap: 6 }}>
                  <button
                    onClick={() => { setDeleteArmed(false); handleDeleteAll() }}
                    disabled={deleteBusy}
                    style={{ ...btnBase, borderColor: 'var(--red)', color: 'var(--red)', flex: 1 }}
                  >
                    {deleteBusy ? 'DELETING…' : '☠ CONFIRM'}
                  </button>
                  <button
                    onClick={() => setDeleteArmed(false)}
                    disabled={deleteBusy}
                    style={{ ...btnBase, flex: 1 }}
                  >
                    CANCEL
                  </button>
                </div>
              </div>
            )}
          </PanelShell>
        )}

        {active === 'nukeAll' && (
          <PanelShell title={`NUKE ALL${nukeBusy || deleteBusy ? ' · WORKING…' : ''}`}>
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--red)', lineHeight: 1.6, marginBottom: 10, fontWeight: 800 }}>
              ⚠ DESTRUCTIVE — 1) deletes all persisting data, 2) docker compose down -v, 3) restarts the whole stack. Data lost permanently.
            </div>
            {nukeMsg && (
              <div style={{
                border: '1px solid var(--red)',
                borderRadius: 6,
                padding: 10,
                marginBottom: 10,
                background: 'rgba(255,23,68,0.08)',
                fontFamily: 'var(--font-mono)',
                fontSize: 9,
                color: nukeMsg.startsWith('✓') ? 'var(--green)' : nukeMsg.startsWith('✕') ? 'var(--red)' : 'var(--yellow)',
                lineHeight: 1.6,
              }}>
                {nukeMsg}
              </div>
            )}
            {!nukeArmed ? (
              <button
                onClick={() => setNukeArmed(true)}
                disabled={nukeBusy || deleteBusy}
                style={{ ...btnBase, borderColor: 'var(--orange)', color: 'var(--orange)', width: '100%' }}
              >
                {nukeBusy || deleteBusy ? 'WORKING…' : '☢ ARM NUKE'}
              </button>
            ) : (
              <div>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--red)', marginBottom: 8, fontWeight: 800 }}>
                  CONFIRM: wipe ALL data + docker volumes?
                </div>
                <div style={{ display: 'flex', gap: 6 }}>
                  <button
                    onClick={() => { setNukeArmed(false); handleNukeAll() }}
                    disabled={nukeBusy || deleteBusy}
                    style={{ ...btnBase, borderColor: 'var(--red)', color: 'var(--red)', flex: 1 }}
                  >
                    {nukeBusy || deleteBusy ? 'NUKING…' : '💥 CONFIRM NUKE'}
                  </button>
                  <button
                    onClick={() => setNukeArmed(false)}
                    disabled={nukeBusy || deleteBusy}
                    style={{ ...btnBase, flex: 1 }}
                  >
                    CANCEL
                  </button>
                </div>
              </div>
            )}
          </PanelShell>
        )}

        {active !== 'dispatch' && active !== 'data' && active !== 'logs' && active !== 'negcontrol' && active !== 'notifications' && active !== 'manual' && active !== 'guide' && active !== 'deleteAll' && active !== 'nukeAll' && (
          <PanelShell title={items.find(i => i.id === active)?.label || active.toUpperCase()}>
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-muted)', lineHeight: 1.6 }}>
              Placeholder: next step wire to existing/implemented tabs.
            </div>
          </PanelShell>
        )}
      </div>
    </aside>
  )
}
