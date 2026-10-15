import React from 'react'

function StatCard({ label, value, icon, accent = 'var(--green)', subText, isAlert }) {
  return (
    <div style={{
      background: 'linear-gradient(180deg, #0d1522 0%, #080e16 100%)',
      border: isAlert ? '1px solid rgba(255, 23, 68, 0.4)' : '1px solid var(--border)',
      borderRadius: 4,
      padding: '8px 12px',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      position: 'relative',
      overflow: 'hidden',
      boxShadow: isAlert ? '0 0 15px rgba(255, 23, 68, 0.15)' : 'none',
    }}>
      {/* Top micro-line glow */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        height: 2,
        background: accent,
        opacity: isAlert ? 1 : 0.6,
      }} />

      <div>
        <div style={{
          display: 'flex',
          alignItems: 'baseline',
          gap: 6,
        }}>
          <span style={{
            fontFamily: 'var(--font-mono)',
            fontSize: 20,
            fontWeight: 800,
            color: accent,
            letterSpacing: '-0.02em',
          }}>
            {value}
          </span>
          {subText && (
            <span style={{ fontSize: 9, color: 'var(--text-muted)', fontFamily: 'var(--font-mono)' }}>
              {subText}
            </span>
          )}
        </div>
        <div style={{
          fontSize: 9,
          fontFamily: 'var(--font-mono)',
          color: 'var(--text-secondary)',
          textTransform: 'uppercase',
          letterSpacing: '0.08em',
          marginTop: 2,
        }}>
          {label}
        </div>
      </div>

      <div style={{
        fontSize: 18,
        opacity: 0.8,
        filter: `drop-shadow(0 0 8px ${accent})`,
      }}>
        {icon}
      </div>
    </div>
  )
}

export default function StatsBar({ stats }) {
  const s = stats || {
    actors: 0,
    links: 0,
    jobs_done: 0,
    flagged: 0,
    wallets: 0,
    pgp_keys: 0,
  }

  return (
    <div style={{
      display: 'grid',
      gridTemplateColumns: 'repeat(6, 1fr)',
      gap: 10,
      padding: '8px 16px',
      borderBottom: '1px solid var(--border)',
      background: 'var(--bg-deep)',
      zIndex: 10,
    }}>
      <StatCard label="Threat Actors" value={s.actors || 0} icon="💀" accent="var(--green)" subText="TRACKED" />
      <StatCard label="Persona Links" value={s.links || 0} icon="🕸" accent="var(--cyan)" subText="CORRELATED" />
      <StatCard label="Flagged CSAM" value={s.flagged || 0} icon="⚠" accent="var(--red)" subText="CRITICAL" isAlert={true} />
      <StatCard label="Crypto Wallets" value={s.wallets || 0} icon="₿" accent="var(--orange)" subText="MONITORED" />
      <StatCard label="PGP Keyrings" value={s.pgp_keys || 0} icon="🔑" accent="var(--purple)" subText="FINGERPRINTS" />
      <StatCard label="Crawl Depth" value={s.jobs_done || 0} icon="⚡" accent="var(--yellow)" subText="TASKS DONE" />
    </div>
  )
}

