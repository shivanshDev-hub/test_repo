import React from 'react'
import LiveLogs from './LiveLogs'

const NODE_COLORS = {
  actor: '#00ff9d',
  pgp: '#d500f9',
  wallet: '#ff9100',
  onion: '#00b0ff',
  infra: '#ff1744',
  email: '#ffea00',
  stylometry: '#b388ff',
}

function confidenceColor(c) {
  if (c >= 0.85) return '#00ff9d'
  if (c >= 0.70) return '#00e5ff'
  if (c >= 0.55) return '#ff9100'
  return '#ff1744'
}

export default function Sidebar({
  tab,
  setTab,
  leads,
  aiMode,
  onSelectNode,
  onOpenNewScan,
  filterType,
  setFilterType,
  confidenceThreshold,
  setConfidenceThreshold,
  // New: “very left” control rail actions/state
  onSearchClear,
  onTorClear,
  onAiToggle,
  toolUiState,
}) {
  return (
    <aside style={{
      width: 290,
      background: 'linear-gradient(180deg, #070c14 0%, #04070c 100%)',
      borderRight: '1px solid var(--border)',
      display: 'flex',
      flexDirection: 'column',
      flexShrink: 0,
      overflow: 'hidden',
      userSelect: 'none',
      zIndex: 50,
    }}>
      {/* Navigation Tabs */}
      <div style={{
        display: 'flex',
        borderBottom: '1px solid var(--border)',
        background: 'rgba(5, 8, 13, 0.95)',
      }}>
        {[
          { id: 'graph', label: 'RADAR', icon: '🎯' },
          { id: 'leads', label: `LEADS (${leads?.length || 0})`, icon: '⚡' },
          { id: 'logs', label: 'SYS LOGS', icon: '🖥' },
        ].map(t => {
          const active = tab === t.id
          return (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              style={{
                flex: 1,
                padding: '10px 4px',
                background: active ? 'rgba(0, 255, 157, 0.06)' : 'transparent',
                color: active ? 'var(--green)' : 'var(--text-secondary)',
                border: 'none',
                borderBottom: active ? '2px solid var(--green)' : '2px solid transparent',
                fontFamily: 'var(--font-mono)',
                fontSize: 10,
                fontWeight: active ? 700 : 500,
                cursor: 'pointer',
                letterSpacing: '0.06em',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: 4,
              }}
            >
              <span>{t.icon}</span>
              <span>{t.label}</span>
            </button>
          )
        })}
      </div>
      {/* Drawer Content */}
      <div style={{ flex: 1, overflowY: 'auto', padding: '10px', display: 'flex', flexDirection: 'column', gap: 10 }}>
        {tab === 'graph' && (
          <>
            <button
              onClick={onOpenNewScan}
              style={{
                padding: '8px',
                background: 'rgba(0, 255, 157, 0.08)',
                border: '1px dashed var(--green)',
                borderRadius: 4,
                color: 'var(--green)',
                fontFamily: 'var(--font-mono)',
                fontSize: 10,
                fontWeight: 700,
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: 6,
              }}
            >
              <span>⚡</span> + DISPATCH INVESTIGATION
            </button>

            {/* Entity Filters */}
            <div style={{ background: 'var(--bg-card)', border: '1px solid var(--border)', borderRadius: 5, padding: '8px' }}>
              <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: 6 }}>
                ENTITY FILTER TOGGLES
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 5 }}>
                {[
                  { id: 'all', label: 'All Entities', color: 'var(--text-primary)' },
                  { id: 'actor', label: 'Actors / Vendors', color: NODE_COLORS.actor },
                  { id: 'stylometry', label: 'Stylometry (✍)', color: NODE_COLORS.stylometry },
                  { id: 'wallet', label: 'Crypto Wallets', color: NODE_COLORS.wallet },
                  { id: 'pgp', label: 'PGP Keys', color: NODE_COLORS.pgp },
                  { id: 'onion', label: 'Hidden Svcs', color: NODE_COLORS.onion },
                  { id: 'infra', label: 'Infra Leaks', color: NODE_COLORS.infra },
                  { id: 'flagged', label: '⚠ Flagged Only', color: 'var(--red)' },
                ].map(f => {
                  const active = filterType === f.id
                  return (
                    <button
                      key={f.id}
                      onClick={() => setFilterType(f.id)}
                      style={{
                        padding: '4px 6px',
                        background: active ? 'rgba(0, 255, 157, 0.12)' : 'rgba(255, 255, 255, 0.02)',
                        border: active ? '1px solid var(--green)' : '1px solid var(--border)',
                        borderRadius: 3,
                        color: active ? 'var(--green)' : f.color,
                        fontFamily: 'var(--font-mono)',
                        fontSize: 9,
                        fontWeight: active ? 700 : 500,
                        cursor: 'pointer',
                        textAlign: 'left',
                      }}
                    >
                      {f.label}
                    </button>
                  )
                })}
              </div>
            </div>

            {/* Confidence Cutoff Slider */}
            <div style={{ background: 'var(--bg-card)', border: '1px solid var(--border)', borderRadius: 5, padding: '8px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 4 }}>
                <span style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-muted)', textTransform: 'uppercase' }}>
                  CONFIDENCE CUTOFF
                </span>
                <span style={{ fontFamily: 'var(--font-mono)', fontSize: 10, color: 'var(--green)', fontWeight: 700 }}>
                  {(confidenceThreshold * 100).toFixed(0)}%
                </span>
              </div>
              <input
                type="range"
                min="0"
                max="0.95"
                step="0.05"
                value={confidenceThreshold}
                onChange={e => setConfidenceThreshold(parseFloat(e.target.value))}
                style={{ width: '100%', accentColor: 'var(--green)', cursor: 'pointer' }}
              />
            </div>

            {/* Tor Routing Info */}
            <div style={{ background: 'var(--bg-card)', border: '1px solid var(--border)', borderRadius: 5, padding: '8px' }}>
              <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: 4 }}>
                TOR CIRCUIT ROUTING
              </div>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 9, color: 'var(--cyan)', display: 'flex', flexDirection: 'column', gap: 3 }}>
                <div>Guard: 185.220.101.5 [DE]</div>
                <div>Middle: 198.51.100.22 [NL]</div>
                <div>Exit: SOCKS5:9050 [OK]</div>
              </div>
            </div>
          </>
        )}

        {tab === 'leads' && (
          <>
            <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-muted)', textTransform: 'uppercase' }}>
              CORRELATED THREAT LEADS ({leads?.length || 0})
            </div>
            {(leads || []).map(l => (
              <div
                key={l.link_id}
                onClick={() => onSelectNode(l.actor_id_a)}
                style={{
                  background: 'linear-gradient(180deg, #0f1724 0%, #0a1019 100%)',
                  borderRadius: 5,
                  padding: '8px',
                  border: '1px solid var(--border)',
                  cursor: 'pointer',
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 4 }}>
                  <span style={{ fontFamily: 'var(--font-mono)', fontSize: 11, fontWeight: 700, color: confidenceColor(l.attribution_confidence) }}>
                    {(l.attribution_confidence * 100).toFixed(0)}% MATCH
                  </span>
                  <span style={{ fontSize: 8, fontFamily: 'var(--font-mono)', fontWeight: 700, color: l.nc_state === 'PASS' ? 'var(--green)' : 'var(--orange)', padding: '1px 4px', background: 'rgba(0,0,0,0.4)', borderRadius: 2 }}>
                    NC: {l.nc_state}
                  </span>
                </div>
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: 10, fontWeight: 600, color: 'var(--text-primary)' }}>
                  👤 {l.label_a} ⇄ 👤 {l.label_b}
                </div>
                <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--text-secondary)', marginTop: 4, background: 'rgba(0,0,0,0.3)', padding: '3px 5px', borderRadius: 2 }}>
                  {l.basis}
                </div>
              </div>
            ))}
          </>
        )}

        {tab === 'logs' && (
          <div style={{
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            minHeight: 0,
            background: '#020406',
            border: '1px solid var(--border)',
            borderRadius: 5,
            padding: '6px',
          }}>
            {/* Live per-service Docker log viewer — the whole tab is a single log window. */}
            <LiveLogs aiMode={aiMode} />
          </div>
        )}
      </div>
    </aside>
  )
}

