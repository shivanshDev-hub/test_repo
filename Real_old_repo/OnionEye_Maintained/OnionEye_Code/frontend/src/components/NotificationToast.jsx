import React, { useEffect, useRef, useState } from 'react'

// NotificationToast renders a transient tactical popup (top-right, below the
// header) for each newly observed crawler-worker event. It auto-dismisses after
// 3s, matching the dashboard notification contract.
//
// Supported events (payload via `notification`):
//   worker-rotation → green accent, "Worker Circuit Rotation"
//   worker-anomaly  → red accent, "Worker Anomaly — Killed & Remade"
const EVENT_META = {
  'worker-rotation': {
    accent: 'var(--green)',
    glow: 'rgba(0, 255, 157, 0.15)',
    glyph: '🔄',
    title: 'Worker Circuit Rotation',
    idLabel: 'ROT',
    idKey: 'rotation_id',
  },
  'worker-anomaly': {
    accent: 'var(--red)',
    glow: 'rgba(255, 23, 68, 0.18)',
    glyph: '⚠️',
    title: 'Worker Anomaly — Killed & Remade',
    idLabel: 'WORKER',
    idKey: 'worker_id',
  },
}

// prettifyKind renders anomaly kinds ("rapid-failure-storm" → "Rapid Failure Storm")
function prettifyKind(kind) {
  return String(kind || '')
    .replace(/-/g, ' ')
    .replace(/\b\w/g, (c) => c.toUpperCase())
}

export default function NotificationToast({ notification, onDismiss }) {
  const [visible, setVisible] = useState(false)
  const timerRef = useRef(null)

  useEffect(() => {
    if (!notification) {
      setVisible(false)
      return
    }
    setVisible(true)
    // Auto-dismiss after 3s regardless of event/reason/at.
    if (timerRef.current) clearTimeout(timerRef.current)
    timerRef.current = setTimeout(() => {
      setVisible(false)
      if (onDismiss) onDismiss(notification)
    }, 3000)
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current)
    }
  }, [notification, onDismiss])

  if (!notification || !visible) return null

  const meta = EVENT_META[notification.event] || EVENT_META['worker-rotation']
  const reason = notification.reason || 'no reason reported'
  const idValue =
    notification[meta.idKey] || notification.anomaly_id || notification.rotation_id || '—'
  const kind = notification.event === 'worker-anomaly'
    ? (notification.anomaly ? prettifyKind(notification.anomaly) : 'Worker anomaly')
    : null
  const at = notification.at
  let timeLabel = ''
  if (at) {
    const d = new Date(at)
    if (!Number.isNaN(d.getTime())) {
      timeLabel = d.toTimeString().slice(0, 8) + ' UTC'
    }
  }

  return (
    <div style={{
      position: 'fixed',
      top: 62,
      right: 16,
      zIndex: 2000,
      minWidth: 260,
      maxWidth: 360,
      background: 'linear-gradient(180deg, #0a111c 0%, #05080e 100%)',
      border: '1px solid var(--border-bright)',
      borderLeft: `3px solid ${meta.accent}`,
      borderRadius: 6,
      padding: '10px 12px',
      boxShadow: `0 12px 40px rgba(0,0,0,0.9), 0 0 24px ${meta.glow}`,
      fontFamily: 'var(--font-mono)',
      animation: 'fade-in 0.25s ease-out',
    }}>
      {/* Top micro-line glow */}
      <div style={{
        position: 'absolute',
        top: 0, left: 0, right: 0,
        height: 2,
        background: meta.accent,
        opacity: 0.7,
      }} />

      <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 6 }}>
        <span style={{ fontSize: 12, filter: `drop-shadow(0 0 6px ${meta.accent})` }}>{meta.glyph}</span>
        <span style={{
          fontSize: 9,
          fontWeight: 800,
          letterSpacing: '0.08em',
          color: meta.accent,
          textTransform: 'uppercase',
        }}>
          {meta.title}
        </span>
        <span style={{ marginLeft: 'auto', fontSize: 8, color: 'var(--text-muted)' }}>{timeLabel}</span>
      </div>

      {kind && (
        <div style={{ fontSize: 9, color: meta.accent, marginBottom: 4, letterSpacing: '0.05em' }}>
          {kind}
        </div>
      )}
      <div style={{ fontSize: 10, color: 'var(--text-primary)', lineHeight: 1.5, wordBreak: 'break-word' }}>
        {reason}
      </div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 6, gap: 8 }}>
        <span style={{ fontSize: 8, color: 'var(--cyan)', letterSpacing: '0.04em' }}>
          {meta.idLabel}: {idValue}
        </span>
        <button
          onClick={() => {
            setVisible(false)
            if (onDismiss) onDismiss(notification)
          }}
          style={{
            background: 'none',
            border: 'none',
            color: 'var(--text-muted)',
            cursor: 'pointer',
            fontSize: 10,
            lineHeight: 1,
            padding: 0,
          }}
          title="Dismiss"
        >
          ✕
        </button>
      </div>
    </div>
  )
}