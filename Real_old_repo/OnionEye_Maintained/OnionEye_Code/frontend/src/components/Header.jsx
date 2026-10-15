import React, { useState, useEffect } from 'react'
import logoUrl from '../assets/official_logo.svg'

export default function Header({
  isLiveApi,
}) {
  const [timeStr, setTimeStr] = useState('')

  useEffect(() => {
    function updateTime() {
      const now = new Date()
      setTimeStr(now.toUTCString().slice(17, 25) + ' UTC')
    }
    updateTime()
    const timer = setInterval(updateTime, 1000)
    return () => clearInterval(timer)
  }, [])

  return (
    <header style={{
      height: 64,
      background: 'linear-gradient(180deg, #09101a 0%, #05080e 100%)',
      borderBottom: '1px solid var(--border)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      padding: '0 16px',
      zIndex: 100,
      position: 'relative',
      boxShadow: '0 4px 20px rgba(0,0,0,0.6)',
      userSelect: 'none',
    }}>
      {/* Left: Brand + Classification */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: 8,
          background: 'rgba(0, 255, 157, 0.06)',
          border: '1px solid var(--border-glow)',
          boxShadow: '0 0 12px rgba(0, 255, 157, 0.15)',
          padding: '5px 10px',
          borderRadius: 4,
        }}>
          <img src={logoUrl} alt="OnionEye" title="OnionEye" style={{ height: 56, width: 56, borderRadius: 6, display: 'block', filter: 'drop-shadow(0 0 6px rgba(0, 255, 157, 0.6)) drop-shadow(0 0 14px rgba(0, 255, 157, 0.25))' }} />
          <div>
            <div style={{
              fontFamily: 'var(--font-mono)',
              fontWeight: 800,
              fontSize: 12,
              letterSpacing: '0.12em',
              color: 'var(--green)',
              textTransform: 'uppercase',
              display: 'flex',
              alignItems: 'center',
              gap: 4,
            }}>
              ONIONEYE
              <span style={{ fontSize: 9, color: 'var(--text-muted)', fontWeight: 400 }}>v2.4</span>
            </div>
            <div style={{ fontSize: 8, fontFamily: 'var(--font-mono)', color: 'var(--text-secondary)', letterSpacing: '0.04em' }}>
              DARK WEB LINK INTELLIGENCE
            </div>
          </div>
        </div>

        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: 4,
          background: 'rgba(255, 145, 0, 0.08)',
          border: '1px solid rgba(255, 145, 0, 0.3)',
          padding: '2px 6px',
          borderRadius: 3,
          fontFamily: 'var(--font-mono)',
          fontSize: 9,
          fontWeight: 700,
          color: 'var(--orange)',
          letterSpacing: '0.06em',
        }}>
          <span>🔒 TLP:AMBER</span>
        </div>

        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: 4,
          background: 'rgba(255, 23, 68, 0.08)',
          border: '1px solid rgba(255, 23, 68, 0.35)',
          padding: '2px 6px',
          borderRadius: 3,
          fontFamily: 'var(--font-mono)',
          fontSize: 9,
          fontWeight: 700,
          color: 'var(--red)',
        }}>
          <span style={{ width: 5, height: 5, borderRadius: '50%', background: 'var(--red)', animation: 'pulse-red 1.5s infinite' }} />
          <span>CRITICAL THREATS</span>
        </div>
      </div>

      {/* Right: Actions & Telemetry */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
        <div style={{
          fontFamily: 'var(--font-mono)',
          fontSize: 10,
          color: 'var(--cyan)',
          background: 'rgba(0, 229, 255, 0.05)',
          padding: '4px 8px',
          borderRadius: 3,
          border: '1px solid rgba(0, 229, 255, 0.2)',
          letterSpacing: '0.04em',
        }}>
          ⏱ {timeStr}
        </div>

        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: 6,
          fontSize: 9,
          fontFamily: 'var(--font-mono)',
          padding: '4px 8px',
          background: 'rgba(12, 20, 32, 0.8)',
          borderRadius: 3,
          border: '1px solid var(--border)',
        }}>
          <span style={{
            width: 6,
            height: 6,
            borderRadius: '50%',
            background: isLiveApi ? 'var(--green)' : 'var(--orange)',
            display: 'inline-block',
            boxShadow: isLiveApi ? '0 0 8px #00ff9d' : '0 0 6px #ff9100',
          }} />
          <span style={{ color: isLiveApi ? 'var(--green)' : 'var(--orange)', fontWeight: 600 }}>
            {isLiveApi ? 'FEED: LIVE' : 'FEED: LOCAL'}
          </span>
        </div>
      </div>
    </header>
  )
}
