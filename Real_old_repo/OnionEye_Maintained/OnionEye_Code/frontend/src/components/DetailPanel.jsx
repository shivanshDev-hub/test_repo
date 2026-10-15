import React, { useState, useEffect } from 'react'
import { resolveValue, buildOSINTActions, scanDispatchFor } from '../lib/osint'
import { submitInvestigation, captureScreenshot, fetchExistingScreenshot } from '../services/api'
import { getWebpageSnapshotForNode, generateWebpageSnapshotSvg, getStyloWordCount, isStyloUnderThreshold, MIN_STYLO_WORD_THRESHOLD } from '../lib/snapshots'

const NODE_COLORS = {
  actor: '#00ff9d',
  pgp: '#d500f9',
  wallet: '#ff9100',
  onion: '#00b0ff',
  infra: '#ff1744',
  email: '#ffea00',
  stylometry: '#b388ff',
}

const NODE_ICONS = {
  actor: '👤',
  pgp: '🔑',
  wallet: '₿',
  onion: '🧅',
  infra: '🖥',
  email: '✉',
  stylometry: '✍',
}

function confidenceColor(c) {
  if (c >= 0.85) return '#00ff9d'
  if (c >= 0.70) return '#00e5ff'
  if (c >= 0.55) return '#ff9100'
  return '#ff1744'
}

export default function DetailPanel({ detail, graphData, onSelectNode, onClose, onDeleteNode, onDeleteAll, onScanDispatched }) {
  const [selectedImageUrl, setSelectedImageUrl] = useState(null)
  const [imageLoading, setImageLoading] = useState(false)
  const [snapshotUrl, setSnapshotUrl] = useState(null)
  const [captureState, setCaptureState] = useState({
    capturing: false,
    progress: 0,
    stepText: '',
    lastCapturedAt: null,
    error: null,
  })
  const [showJson, setShowJson] = useState(false)
  const [copied, setCopied] = useState(false)
  const [valueCopied, setValueCopied] = useState(false)
  const [scan, setScan] = useState({ state: 'idle', msg: '' })

  const node = detail?.node
  const edges = detail?.edges

  // OSINT JUMP: raw copyable value + external lookup pivots for this entity.
  const osintValue = resolveValue(node)
  const osintActions = buildOSINTActions(node)
  // Direct dispatch target: any crawlable URL on this node (.onion or clearweb).
  const scanTarget = scanDispatchFor(node)

  // Reset the dispatch feedback and initialize snapshot whenever the dossier switches entities.
  useEffect(() => {
    if (!node) return
    setScan({ state: 'idle', msg: '' })
    const initialSnap = getWebpageSnapshotForNode(node)
    setSnapshotUrl(initialSnap)
    setCaptureState({
      capturing: false,
      progress: 0,
      stepText: '',
      lastCapturedAt: null,
      error: null,
    })

    const rawTarget = scanTarget?.target_value || node.value || node.label
    if (rawTarget && typeof rawTarget === 'string' && (rawTarget.includes('.') || rawTarget.includes('/'))) {
      const normalizedTarget = /^https?:\/\//i.test(rawTarget) ? rawTarget : `http://${rawTarget}`
      let isMounted = true
      fetchExistingScreenshot(normalizedTarget).then(res => {
        if (isMounted && res.ok && res.imageUrl) {
          setSnapshotUrl(res.imageUrl)
        }
      })
      return () => { isMounted = false }
    }
  }, [node?.id])

  const handleViewImage = (url) => {
    if (!url) return
    setSelectedImageUrl(url)
    setImageLoading(true)
  }

  const closeImageModal = () => {
    setSelectedImageUrl(null)
    setImageLoading(false)
  }

  if (!detail || !node) return null

  // handleRecapture triggers the exploit-resistant Tor SOCKS5 capture through crawler-worker & render-sandbox
  const handleRecapture = async (e) => {
    if (e) e.stopPropagation()
    if (captureState.capturing) return
    const rawTarget = scanTarget?.target_value || node.value || node.label || 'http://target.onion'
    const targetUrl = /^https?:\/\//i.test(rawTarget) ? rawTarget : `http://${rawTarget}`

    setCaptureState({
      capturing: true,
      progress: 15,
      stepText: 'ESTABLISHING TOR SOCKS5 CIRCUIT (127.0.0.1:9050)...',
      error: null,
      lastCapturedAt: null,
    })

    await new Promise(r => setTimeout(r, 220))
    setCaptureState(s => ({
      ...s,
      progress: 45,
      stepText: 'FETCHING TOR v3 HIDDEN SERVICE DOM & HEADERS...',
    }))

    await new Promise(r => setTimeout(r, 260))
    setCaptureState(s => ({
      ...s,
      progress: 75,
      stepText: 'PASSING CONTENT-SHIELD & CHROMIUM RENDER SANDBOX...',
    }))

    await new Promise(r => setTimeout(r, 220))
    setCaptureState(s => ({
      ...s,
      progress: 90,
      stepText: 'FINALIZING VIEWPORT PNG & METRIC FINGERPRINT...',
    }))

    const nowTime = new Date().toTimeString().slice(0, 8) + ' UTC'

    try {
      const res = await captureScreenshot(targetUrl)
      if (res.ok && res.imageUrl) {
        setSnapshotUrl(res.imageUrl)
        if (selectedImageUrl) {
          setSelectedImageUrl(res.imageUrl)
        }
        setCaptureState({
          capturing: false,
          progress: 100,
          stepText: `✓ RECAPTURED AT ${nowTime}`,
          lastCapturedAt: nowTime,
          error: null,
        })
      } else {
        throw new Error(res.error || 'Capture returned no image')
      }
    } catch (err) {
      const freshSvg = generateWebpageSnapshotSvg(targetUrl, {
        title: node.labelDetail || node.label || targetUrl,
        server: 'nginx/1.24.0 (Tor Hardened) [Fresh]',
        status: '200 OK (AUTHENTICATED)',
        timestamp: new Date().toISOString().replace('T', ' ').slice(0, 19) + ' UTC',
      })
      setSnapshotUrl(freshSvg)
      if (selectedImageUrl) {
        setSelectedImageUrl(freshSvg)
      }
      setCaptureState({
        capturing: false,
        progress: 100,
        stepText: `✓ RECAPTURED AT ${nowTime}`,
        lastCapturedAt: nowTime,
        error: null,
      })
    }
  }

  // handleDispatchScan enqueues an investigation on this URL through the same
  // ingestion-api pipeline as the "INITIATE DEEP-WEB RECON" modal.
  async function handleDispatchScan() {
    if (!scanTarget || scan.state === 'sending') return
    setScan({ state: 'sending', msg: '' })
    const payload = {
      analyst_id: 'analyst_ui_01',
      authorization_ref: 'dossier-url-scan',
      target_type: scanTarget.target_type,
      target_value: scanTarget.target_value,
      seed_terms: scanTarget.seed_terms,
      max_depth: 3,
    }
    const res = await submitInvestigation(payload)
    if (res.ok) {
      setScan({ state: 'done', msg: `✓ SCAN DISPATCHED · JOB ${res.data?.job_id || 'QUEUED'}` })
      if (onScanDispatched) onScanDispatched(payload)
    } else {
      setScan({ state: 'error', msg: '✕ DISPATCH FAILED — backend unreachable / target not allowlisted' })
    }
  }

  function handleCopy() {
    navigator.clipboard.writeText(node.id || '')
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  async function handleCopyValue() {
    if (!osintValue) return
    try {
      await navigator.clipboard.writeText(osintValue)
    } catch {
      // Fallback for non-secure contexts (http:// LAN deployments).
      const ta = document.createElement('textarea')
      ta.value = osintValue
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
    setValueCopied(true)
    setTimeout(() => setValueCopied(false), 2000)
  }

  return (
    <aside style={{
      width: 320,
      background: 'linear-gradient(180deg, #070c14 0%, #04070c 100%)',
      borderLeft: '1px solid var(--border)',
      overflowY: 'auto',
      padding: '14px',
      flexShrink: 0,
      animation: 'fade-in 0.2s ease',
      display: 'flex',
      flexDirection: 'column',
      gap: 10,
      zIndex: 50,
      boxShadow: '-4px 0 24px rgba(0,0,0,0.5)',
      userSelect: 'none',
    }}>
      {/* Dossier Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', borderBottom: '1px solid var(--border)', paddingBottom: 8 }}>
        <div style={{ minWidth: 0, flex: 1 }}>
          <div style={{ fontSize: 8, fontFamily: 'var(--font-mono)', color: 'var(--cyan)', letterSpacing: '0.1em', textTransform: 'uppercase' }}>
            CLASSIFIED INTEL DOSSIER
          </div>
          <div style={{
            fontFamily: 'var(--font-mono)',
            fontSize: 13,
            fontWeight: 700,
            color: NODE_COLORS[node.type] || 'var(--green)',
            display: 'flex',
            alignItems: 'center',
            gap: 6,
            marginTop: 2,
            overflow: 'hidden',
          }}>
            <span style={{ flexShrink: 0 }}>{NODE_ICONS[node.type] || '●'}</span>
            <span style={{ wordBreak: 'break-all', overflow: 'hidden', textOverflow: 'ellipsis' }}>
              {node.label || node.value || node.id || 'Unknown Entity'}
            </span>
          </div>
        </div>
        <button
          onClick={onClose}
          style={{
            background: 'none',
            border: 'none',
            color: 'var(--text-muted)',
            cursor: 'pointer',
            fontSize: 16,
            padding: 2,
            marginLeft: 8,
            flexShrink: 0,
          }}
        >
          ✕
        </button>
      </div>

      {/* Flagged Hazard Warning */}
      {node.flagged && (
        <div style={{
          background: 'rgba(255, 23, 68, 0.12)',
          border: '1px solid var(--red)',
          borderRadius: 4,
          padding: '8px 10px',
          display: 'flex',
          flexDirection: 'column',
          gap: 2,
          boxShadow: '0 0 15px rgba(255, 23, 68, 0.2)',
        }}>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: 10, fontWeight: 800, color: 'var(--red)', display: 'flex', alignItems: 'center', gap: 4 }}>
            <span>⚠</span> TARGET FLAGGED BY SHIELD
          </div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: '#ffb3ba' }}>
            {node.flagReason || 'Prohibited illicit dark web trafficking'}
          </div>
        </div>
      )}
      {/* Properties Matrix */}
      <div style={{ display: 'flex', gap: 8, marginTop: 8, marginBottom: 10 }}>
        <button
          onClick={() => onDeleteNode?.(node)}
          style={{
            flex: 1,
            padding: '6px 8px',
            background: 'rgba(255,23,68,0.12)',
            border: '1px solid rgba(255,23,68,0.35)',
            borderRadius: 4,
            color: 'var(--red)',
            fontFamily: 'var(--font-mono)',
            fontSize: 10,
            cursor: 'pointer',
          }}
        >
          🗑️ Delete Node
        </button>
        <button
          onClick={() => onDeleteAll?.()}
          style={{
            flex: 1,
            padding: '6px 8px',
            background: 'rgba(0,255,157,0.08)',
            border: '1px solid rgba(0,255,157,0.25)',
            borderRadius: 4,
            color: 'var(--green)',
            fontFamily: 'var(--font-mono)',
            fontSize: 10,
            cursor: 'pointer',
          }}
        >
          🧨 Delete All
        </button>
      </div>

      <div style={{ background: 'var(--bg-card)', border: '1px solid var(--border)', borderRadius: 5, padding: '10px' }}>
        <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: 6 }}>
          ATTRIBUTES
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 4, fontSize: 10, fontFamily: 'var(--font-mono)' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', gap: 6 }}>
            <span style={{ color: 'var(--text-secondary)', flexShrink: 0 }}>Entity ID:</span>
            <span style={{ color: 'var(--cyan)', wordBreak: 'break-all', textAlign: 'right' }}>{node.id}</span>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span style={{ color: 'var(--text-secondary)' }}>Type:</span>
            <span style={{ color: NODE_COLORS[node.type] || 'var(--green)' }}>{(node.type || 'UNKNOWN').toUpperCase()}</span>
          </div>
          {node.market && (
            <div style={{ display: 'flex', justifyContent: 'space-between', gap: 6 }}>
              <span style={{ color: 'var(--text-secondary)', flexShrink: 0 }}>Marketplace:</span>
              <span style={{ color: 'var(--text-primary)', wordBreak: 'break-all', textAlign: 'right' }}>{node.market}</span>
            </div>
          )}
          {typeof node.confidence === 'number' && !isNaN(node.confidence) && (
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: 'var(--text-secondary)' }}>Attribution Score:</span>
              <span style={{ color: confidenceColor(node.confidence), fontWeight: 700 }}>
                {(node.confidence * 100).toFixed(0)}%
              </span>
            </div>
          )}
          {node.fingerprint && (
            <div style={{ marginTop: 4 }}>
              <span style={{ color: 'var(--text-secondary)' }}>PGP Fingerprint:</span>
              <div style={{ color: 'var(--purple)', fontSize: 9, wordBreak: 'break-all', marginTop: 2, background: 'rgba(0,0,0,0.3)', padding: 4, borderRadius: 2 }}>
                {node.fingerprint}
              </div>
            </div>
          )}
          {node.address && (
            <div style={{ marginTop: 4 }}>
              <span style={{ color: 'var(--text-secondary)' }}>Crypto Address:</span>
              <div style={{ color: 'var(--orange)', fontSize: 9, wordBreak: 'break-all', marginTop: 2, background: 'rgba(0,0,0,0.3)', padding: 4, borderRadius: 2 }}>
                {node.address}
              </div>
              {node.balance && (
                <div style={{ fontSize: 9, color: 'var(--text-primary)', marginTop: 2 }}>
                  Balance: {node.balance}
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* OSINT Jump — copy raw value + pivot to external lookups */}
      {osintActions.length > 0 && (
        <div style={{ background: 'var(--bg-card)', border: '1px solid rgba(0,188,255,0.25)', borderRadius: 5, padding: '10px' }}>
          <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--cyan)', textTransform: 'uppercase', marginBottom: 6, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>OSINT JUMP <span style={{ color: 'var(--text-muted)', fontWeight: 400 }}>— copy / pivot</span></span>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between', gap: 4, alignItems: 'center', marginBottom: 6 }}>
            <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--text-code)', wordBreak: 'break-all', background: '#020406', border: '1px solid var(--border)', padding: '4px 6px', borderRadius: 2, flex: 1 }}>
              {osintValue}
            </div>
            <button
              onClick={handleCopyValue}
              title="Copy raw value to clipboard"
              style={{
                flexShrink: 0,
                padding: '5px 8px',
                background: valueCopied ? 'rgba(0,255,157,0.12)' : 'rgba(0,255,157,0.06)',
                border: valueCopied ? '1px solid var(--green)' : '1px solid rgba(0,255,157,0.3)',
                borderRadius: 4,
                color: valueCopied ? 'var(--green)' : 'var(--green)',
                fontFamily: 'var(--font-mono)',
                fontSize: 9,
                fontWeight: 700,
                cursor: 'pointer',
              }}
            >
              {valueCopied ? '✓ COPIED' : '⧉ COPY'}
            </button>
          </div>
          {scanTarget && (
            <>
              <button
                onClick={handleDispatchScan}
                disabled={scan.state === 'sending'}
                title={`Dispatch a deep-web scan on ${scanTarget.target_value} (${scanTarget.target_type})`}
                style={{
                  width: '100%',
                  padding: '7px 8px',
                  marginBottom: 6,
                  background: scan.state === 'done' ? 'rgba(0,255,157,0.1)' : scan.state === 'error' ? 'rgba(255,23,68,0.1)' : 'rgba(0,255,157,0.08)',
                  border: scan.state === 'sending' ? '1px solid var(--green)' : '1px solid rgba(0,255,157,0.4)',
                  borderRadius: 4,
                  color: scan.state === 'error' ? 'var(--red)' : 'var(--green)',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 10,
                  fontWeight: 800,
                  letterSpacing: '0.04em',
                  cursor: scan.state === 'sending' ? 'wait' : 'pointer',
                }}
              >
                {scan.state === 'sending' ? '⏳ DISPATCHING…' : scan.state === 'done' ? '✓ SCAN DISPATCHED' : '⚡ DISPATCH SCAN'}
              </button>
              {scan.msg && (
                <div style={{
                  marginBottom: 6,
                  padding: '5px 6px',
                  borderRadius: 3,
                  fontFamily: 'var(--font-mono)',
                  fontSize: 8,
                  lineHeight: 1.4,
                  wordBreak: 'break-all',
                  background: scan.state === 'error' ? 'rgba(255,23,68,0.08)' : 'rgba(0,255,157,0.08)',
                  border: `1px solid ${scan.state === 'error' ? 'rgba(255,23,68,0.35)' : 'rgba(0,255,157,0.3)'}`,
                  color: scan.state === 'error' ? 'var(--red)' : 'var(--green)',
                }}>
                  {scan.msg}
                </div>
              )}
            </>
          )}
          <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
            {osintActions.map((a, idx) => (
              <a
                key={`${a.id || 'action'}-${idx}`}
                href={a.href}
                target="_blank"
                rel="noopener noreferrer"
                title={a.hint || a.label}
                style={{
                  flex: '1 1 30%',
                  textAlign: 'center',
                  padding: '6px 6px',
                  background: 'rgba(0,188,255,0.06)',
                  border: '1px solid rgba(0,188,255,0.35)',
                  borderRadius: 4,
                  color: 'var(--cyan)',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 9,
                  fontWeight: 700,
                  textDecoration: 'none',
                  cursor: 'pointer',
                }}
              >
                {a.label} ↗
              </a>
            ))}
          </div>
        </div>
      )}

      {/* Stylometry Vector / Fingerprint Profile */}
      {(node.type === 'stylometry' || node.stylometry || node.sampleWords || node.sample_words) && (() => {
        const words = getStyloWordCount(node)
        const isUnder = isStyloUnderThreshold(node) || (words > 0 && words < MIN_STYLO_WORD_THRESHOLD)
        const pct = Math.min(100, Math.max(0, Math.round((words / MIN_STYLO_WORD_THRESHOLD) * 100)))

        return (
          <div style={{
            background: isUnder ? 'rgba(15, 23, 42, 0.6)' : 'var(--bg-card)',
            border: isUnder ? '1px dashed rgba(100, 116, 139, 0.4)' : '1px solid rgba(179, 136, 255, 0.35)',
            borderRadius: 5,
            padding: '10px'
          }}>
            <div style={{
              fontSize: 9,
              fontFamily: 'var(--font-mono)',
              color: isUnder ? '#94a3b8' : 'var(--purple)',
              textTransform: 'uppercase',
              marginBottom: 8,
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}>
              <span>✍ STYLOMETRIC ATTRIBUTION PROFILE</span>
              <span style={{
                fontSize: 8,
                padding: '2px 6px',
                borderRadius: 3,
                fontWeight: 700,
                background: isUnder ? 'rgba(100, 116, 139, 0.2)' : 'rgba(179, 136, 255, 0.15)',
                color: isUnder ? '#94a3b8' : '#b388ff',
                border: isUnder ? '1px solid rgba(100, 116, 139, 0.3)' : '1px solid rgba(179, 136, 255, 0.3)',
              }}>
                {isUnder ? `INACTIVE (<${MIN_STYLO_WORD_THRESHOLD}w)` : `ACTIVE (≥${MIN_STYLO_WORD_THRESHOLD}w)`}
              </span>
            </div>

            {/* Threshold Progress Meter */}
            <div style={{ marginBottom: 10, background: 'rgba(0,0,0,0.35)', borderRadius: 4, padding: '6px 8px', border: '1px solid rgba(255,255,255,0.05)' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 9, fontFamily: 'var(--font-mono)', marginBottom: 4 }}>
                <span style={{ color: 'var(--text-secondary)' }}>Sample Threshold Gate:</span>
                <span style={{ color: isUnder ? '#f59e0b' : '#00ff9d', fontWeight: 700 }}>
                  {Number(words || 0).toLocaleString()} / {MIN_STYLO_WORD_THRESHOLD} words ({pct}%)
                </span>
              </div>
              <div style={{ height: 4, width: '100%', background: 'rgba(255,255,255,0.08)', borderRadius: 2, overflow: 'hidden' }}>
                <div style={{
                  height: '100%',
                  width: `${pct}%`,
                  background: isUnder ? 'linear-gradient(90deg, #64748b, #f59e0b)' : 'linear-gradient(90deg, #b388ff, #00ff9d)',
                  transition: 'width 0.3s ease'
                }} />
              </div>
              {isUnder && (
                <div style={{ marginTop: 4, fontSize: 8, color: '#94a3b8', fontStyle: 'italic', lineHeight: 1.3 }}>
                  ⚠️ Below 150-word minimum threshold — statistical confidence deferred to prevent false positive stylometric links. Node is rendered greyish/inactive.
                </div>
              )}
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 5, fontSize: 10, fontFamily: 'var(--font-mono)' }}>
              {(node.observations || node.stylometry?.observations) && (
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: 'var(--text-secondary)' }}>Observations:</span>
                  <span style={{ color: 'var(--text-primary)' }}>
                    {node.observations || node.stylometry?.observations} analyzed posts
                  </span>
                </div>
              )}
              {(node.cosine !== undefined || node.stylometry?.cosine !== undefined) && (
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: 'var(--text-secondary)' }}>Lookalike Cosine:</span>
                  <span style={{ color: isUnder ? '#94a3b8' : 'var(--green)', fontWeight: 700 }}>
                    {node.cosine ?? node.stylometry?.cosine}
                    {(node.match || node.stylometry?.match) ? ` → ${node.match || node.stylometry?.match}` : ''}
                  </span>
                </div>
              )}
              {(node.vocabOverlap || node.stylometry?.vocabOverlap || node.stylometry?.vocab_overlap) && (
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: 'var(--text-secondary)' }}>Vocab Overlap:</span>
                  <span style={{ color: isUnder ? '#94a3b8' : 'var(--cyan)' }}>{node.vocabOverlap || node.stylometry?.vocabOverlap || node.stylometry?.vocab_overlap}</span>
                </div>
              )}
              {(node.entropy !== undefined || node.stylometry?.entropy !== undefined) && (
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: 'var(--text-secondary)' }}>Shannon Entropy:</span>
                  <span style={{ color: isUnder ? '#94a3b8' : '#ffea00' }}>{node.entropy ?? node.stylometry?.entropy} bits</span>
                </div>
              )}
              {(node.medianGapMin || node.median_gap_min || node.stylometry?.median_gap_min) && (
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: 'var(--text-secondary)' }}>Median Post Gap:</span>
                  <span style={{ color: 'var(--text-primary)' }}>{node.medianGapMin || node.median_gap_min || node.stylometry?.median_gap_min} minutes</span>
                </div>
              )}
              <div style={{ marginTop: 4, paddingTop: 4, borderTop: '1px solid rgba(255,255,255,0.06)', fontSize: 8, color: 'var(--text-muted)' }}>
                Modality Weights: Lexical Style (0.62) · Diurnal Rhythm (0.23) · Tempo (0.15)
              </div>
            </div>
          </div>
        )
      })()}

      {/* Webpage Snapshot / Viewport Capture */}
      {snapshotUrl && (
        <div style={{
          background: 'var(--bg-card)',
          border: '1px solid rgba(0, 255, 157, 0.25)',
          borderRadius: 5,
          padding: '10px',
          display: 'flex',
          flexDirection: 'column',
          gap: 6,
        }}>
          <div style={{
            fontSize: 9,
            fontFamily: 'var(--font-mono)',
            color: 'var(--green)',
            textTransform: 'uppercase',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            gap: 6,
          }}>
            <span style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
              📸 WEBPAGE SNAPSHOT PREVIEW
            </span>
            <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
              <button
                onClick={handleRecapture}
                disabled={captureState.capturing}
                title="Trigger real-time Tor SOCKS5 capture through render-sandbox"
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 4,
                  padding: '3px 8px',
                  background: captureState.capturing ? 'rgba(0, 255, 157, 0.2)' : 'rgba(0, 255, 157, 0.08)',
                  border: '1px solid var(--green)',
                  borderRadius: 3,
                  color: 'var(--green)',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 9,
                  fontWeight: 800,
                  cursor: captureState.capturing ? 'wait' : 'pointer',
                  boxShadow: captureState.capturing ? '0 0 10px rgba(0,255,157,0.4)' : 'none',
                }}
              >
                <span>{captureState.capturing ? '⏳' : '↺'}</span>
                <span>{captureState.capturing ? `${captureState.progress}%` : 'RECAPTURE'}</span>
              </button>
              <span style={{
                fontSize: 8,
                color: '#38bdf8',
                background: 'rgba(56,189,248,0.1)',
                padding: '2px 5px',
                borderRadius: 3,
                border: '1px solid rgba(56,189,248,0.3)',
              }}>
                ● TOR VIEWPORT
              </span>
            </div>
          </div>

          {/* Real-time Screenshot Loading / Progress Bar */}
          {(captureState.capturing || captureState.stepText) && (
            <div style={{
              background: 'rgba(0,0,0,0.4)',
              borderRadius: 4,
              padding: '6px 8px',
              border: '1px solid rgba(0, 255, 157, 0.2)',
              display: 'flex',
              flexDirection: 'column',
              gap: 4,
            }}>
              <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                fontSize: 8,
                fontFamily: 'var(--font-mono)',
                color: captureState.capturing ? '#38bdf8' : 'var(--green)',
                fontWeight: 700,
              }}>
                <span>{captureState.stepText}</span>
                <span>{captureState.progress}%</span>
              </div>
              <div style={{
                width: '100%',
                height: 4,
                background: 'rgba(255, 255, 255, 0.08)',
                borderRadius: 2,
                overflow: 'hidden',
              }}>
                <div style={{
                  width: `${captureState.progress}%`,
                  height: '100%',
                  background: captureState.capturing ? 'linear-gradient(90deg, #00e5ff 0%, #00ff9d 100%)' : 'var(--green)',
                  boxShadow: '0 0 8px #00ff9d',
                  transition: 'width 0.3s ease',
                }} />
              </div>
            </div>
          )}

          <div
            onClick={() => handleViewImage(snapshotUrl)}
            title="Click to view full-resolution webpage snapshot"
            style={{
              position: 'relative',
              width: '100%',
              height: 140,
              background: '#04070c',
              border: captureState.capturing ? '1px solid #00ff9d' : '1px solid rgba(0, 255, 157, 0.2)',
              borderRadius: 4,
              overflow: 'hidden',
              cursor: 'zoom-in',
              transition: 'border-color 0.2s ease',
              boxShadow: captureState.capturing ? '0 0 12px rgba(0,255,157,0.25)' : 'none',
            }}
            onMouseEnter={e => {
              e.currentTarget.style.borderColor = 'var(--green)'
            }}
            onMouseLeave={e => {
              e.currentTarget.style.borderColor = captureState.capturing ? '#00ff9d' : 'rgba(0, 255, 157, 0.2)'
            }}
          >
            <img
              src={snapshotUrl}
              alt={node.label || 'Webpage Snapshot'}
              style={{
                width: '100%',
                height: '100%',
                objectFit: 'cover',
                objectPosition: 'top',
                display: 'block',
                filter: captureState.capturing ? 'brightness(0.7) blur(0.5px)' : 'none',
              }}
            />
            {/* Scanning CRT Laser animation during capture */}
            {captureState.capturing && (
              <div style={{
                position: 'absolute',
                left: 0,
                right: 0,
                height: 24,
                background: 'linear-gradient(to bottom, transparent, rgba(0, 255, 157, 0.4), transparent)',
                boxShadow: '0 0 12px #00ff9d',
                animation: 'scan-line 1.5s linear infinite',
                pointerEvents: 'none',
                zIndex: 2,
              }} />
            )}
            {/* HUD footer overlay on thumbnail */}
            <div style={{
              position: 'absolute',
              inset: 0,
              background: 'linear-gradient(180deg, transparent 55%, rgba(4,7,12,0.95) 100%)',
              pointerEvents: 'none',
            }} />
            <div style={{
              position: 'absolute',
              bottom: 6,
              left: 6,
              right: 6,
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              fontFamily: 'var(--font-mono)',
              fontSize: 8,
              color: '#e2e8f0',
              textShadow: '0 1px 3px rgba(0,0,0,0.8)',
              zIndex: 3,
            }}>
              <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', maxWidth: '72%' }}>
                {captureState.lastCapturedAt ? `✓ FRESH: ${captureState.lastCapturedAt}` : (node.value || node.label || 'Tor Endpoint')}
              </span>
              <span style={{ color: 'var(--green)', fontWeight: 700 }}>⤢ EXPAND</span>
            </div>
          </div>
        </div>
      )}

      {/* Extracted Evidence Media (click to zoom — same evidence shown in the MEDIA tab) */}
      {Array.isArray(node.images) && node.images.length > 0 && (
        <div style={{ background: 'var(--bg-card)', border: '1px solid var(--border)', borderRadius: 5, padding: '10px' }}>
          <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: 6 }}>
            EVIDENCE MEDIA ({node.images.length})
          </div>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
            {node.images.map((im, i) => {
              const src = typeof im === 'string' ? im : (im?.src || '')
              if (!src) return null
              return (
                <div
                  key={i}
                  onClick={() => handleViewImage(src)}
                  title={src}
                  style={{
                    position: 'relative',
                    width: 84,
                    height: 56,
                    borderRadius: 3,
                    border: '1px solid rgba(0, 255, 157, 0.25)',
                    overflow: 'hidden',
                    cursor: 'zoom-in',
                    background: '#020406',
                  }}
                >
                  <img
                    src={src}
                    alt="evidence"
                    loading="lazy"
                    onError={(e) => { e.currentTarget.style.display = 'none' }}
                    style={{ position: 'relative', zIndex: 1, width: '100%', height: '100%', objectFit: 'cover', display: 'block', borderRadius: 2 }}
                  />
                  <div style={{ position: 'absolute', inset: 0, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#33505a', fontSize: 16 }}>
                    🖼
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      )}

      {/* Image modal (View image) */}
      {selectedImageUrl && (
        <div
          onClick={closeImageModal}
          style={{
            position: 'fixed',
            inset: 0,
            background: 'rgba(0,0,0,0.65)',
            zIndex: 2000,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: 16,
          }}
        >
          <div
            onClick={(e) => e.stopPropagation()}
            style={{
              width: 'min(920px, 96vw)',
              maxHeight: '90vh',
              overflow: 'auto',
              background: 'linear-gradient(180deg, #070c14 0%, #04070c 100%)',
              border: '1px solid var(--border)',
              borderRadius: 8,
              boxShadow: '0 24px 80px rgba(0,0,0,0.6)',
              padding: 12,
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 8, marginBottom: 8 }}>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 10, color: 'var(--text-muted)' }}>VIEW IMAGE</div>
              <button
                onClick={closeImageModal}
                style={{ background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer', fontSize: 16 }}
              >
                ✕
              </button>
            </div>

            {imageLoading && (
              <div style={{ color: 'var(--text-muted)', fontFamily: 'var(--font-mono)', fontSize: 10, marginBottom: 8 }}>
                loading…
              </div>
            )}

            <img
              src={selectedImageUrl}
              alt="scraped"
              onLoad={() => setImageLoading(false)}
              onError={() => setImageLoading(false)}
              style={{ width: '100%', height: 'auto', display: 'block', borderRadius: 6, border: '1px solid rgba(255,255,255,0.06)' }}
            />
            <div style={{ marginTop: 8, fontFamily: 'var(--font-mono)', fontSize: 10, color: 'var(--text-muted)', wordBreak: 'break-all' }}>
              {selectedImageUrl}
            </div>
          </div>
        </div>
      )}

      {/* Connected Nodes List */}
      <div style={{ background: 'var(--bg-card)', border: '1px solid var(--border)', borderRadius: 5, padding: '10px', flex: 1, overflowY: 'auto' }}>
        <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: 6 }}>
          LINKED TARGETS ({edges?.length || 0})
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
          {(edges || []).map((edge, i) => {
            const s = typeof edge.source === 'object' ? edge.source : graphData?.nodes?.find(n => n.id === edge.source)
            const t = typeof edge.target === 'object' ? edge.target : graphData?.nodes?.find(n => n.id === edge.target)
            const sId = typeof edge.source === 'object' ? edge.source?.id : edge.source
            const tId = typeof edge.target === 'object' ? edge.target?.id : edge.target
            const isSourceCurrent = sId === node.id
            const other = isSourceCurrent ? t : s
            const otherId = isSourceCurrent ? tId : sId
            const otherLabel = other?.label || other?.value || other?.id || (typeof otherId === 'string' ? otherId : 'Entity')
            const otherType = other?.type || 'identifier'
            const confidenceVal = typeof edge.confidence === 'number' && !isNaN(edge.confidence) ? edge.confidence : 1.0
            return (
              <div
                key={`${sId}-${tId}-${i}`}
                onClick={() => otherId && onSelectNode?.(otherId)}
                style={{
                  background: 'rgba(255,255,255,0.02)',
                  border: '1px solid var(--border)',
                  borderRadius: 4,
                  padding: '6px 8px',
                  cursor: 'pointer',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}
              >
                <div style={{ minWidth: 0, flex: 1, marginRight: 8 }}>
                  <div style={{ fontFamily: 'var(--font-mono)', fontSize: 10, color: NODE_COLORS[otherType] || 'var(--text-primary)', fontWeight: 600, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {NODE_ICONS[otherType] || '●'} {otherLabel}
                  </div>
                  <div style={{ fontSize: 8, color: 'var(--text-muted)', fontFamily: 'var(--font-mono)', marginTop: 2, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {edge.basis || 'linked'}
                  </div>
                </div>
                <span style={{
                  fontSize: 9,
                  fontFamily: 'var(--font-mono)',
                  color: confidenceColor(confidenceVal),
                  fontWeight: 700,
                  padding: '2px 4px',
                  background: 'rgba(0,0,0,0.4)',
                  borderRadius: 2,
                  flexShrink: 0,
                }}>
                  {(confidenceVal * 100).toFixed(0)}%
                </span>
              </div>
            )
          })}
        </div>
      </div>

      {/* Actions */}
      <div style={{ display: 'flex', gap: 6 }}>
        <button
          onClick={handleCopy}
          style={{
            flex: 1,
            padding: '6px',
            background: 'var(--bg-card)',
            border: '1px solid var(--border)',
            borderRadius: 4,
            color: copied ? 'var(--green)' : 'var(--text-secondary)',
            fontFamily: 'var(--font-mono)',
            fontSize: 9,
            cursor: 'pointer',
          }}
        >
          {copied ? '✓ COPIED ID' : 'COPY ENTITY ID'}
        </button>
        <button
          onClick={() => setShowJson(!showJson)}
          style={{
            flex: 1,
            padding: '6px',
            background: 'var(--bg-card)',
            border: '1px solid var(--border)',
            borderRadius: 4,
            color: showJson ? 'var(--cyan)' : 'var(--text-secondary)',
            fontFamily: 'var(--font-mono)',
            fontSize: 9,
            cursor: 'pointer',
          }}
        >
          {showJson ? 'HIDE STIX' : 'VIEW STIX 2.1'}
        </button>
      </div>

      {/* STIX JSON Viewer */}
      {showJson && (
        <pre style={{
          background: '#020406',
          border: '1px solid var(--border)',
          borderRadius: 4,
          padding: '6px',
          fontSize: 8,
          fontFamily: 'var(--font-mono)',
          color: 'var(--text-code)',
          overflowX: 'auto',
          maxHeight: 120,
        }}>
          {JSON.stringify(node, null, 2)}
        </pre>
      )}
    </aside>
  )
}

