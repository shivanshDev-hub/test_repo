import React from 'react'

// AiRestartModal confirms a full docker-compose stack restart (with or without
// the AI profile). Restarting brings the whole tool down for ~15s, so the modal
// clearly warns the user and instructs them to reload the page afterwards.
export default function AiRestartModal({ withAi, aiRestarting, onConfirm, onClose }) {
  // withAi === null => modal closed; true => "restart with AI"; false => "restart without AI".
  if (withAi === null || withAi === undefined) return null

  const enabled = withAi
  return (
    <div style={{
      position: 'fixed',
      inset: 0,
      zIndex: 999,
      background: 'rgba(2, 4, 6, 0.9)',
      backdropFilter: 'blur(8px)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: 20,
    }}>
      <div style={{
        background: 'linear-gradient(180deg, #0a111c 0%, #05080e 100%)',
        border: '1px solid var(--border-bright)',
        borderRadius: 8,
        width: '100%',
        maxWidth: 440,
        padding: 20,
        boxShadow: '0 20px 60px rgba(0,0,0,0.9), 0 0 30px rgba(255, 145, 0, 0.15)',
        position: 'relative',
      }}>
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          borderBottom: '1px solid var(--border)',
          paddingBottom: 10,
          marginBottom: 14,
        }}>
          <div>
            <div style={{ fontSize: 9, fontFamily: 'var(--font-mono)', color: 'var(--orange)', letterSpacing: '0.1em' }}>
              // RUNTIME MODE SWITCH //
            </div>
            <div style={{
              fontFamily: 'var(--font-mono)',
              fontSize: 13,
              color: enabled ? 'var(--purple)' : 'var(--cyan)',
              fontWeight: 800,
              display: 'flex',
              alignItems: 'center',
              gap: 6,
              marginTop: 2,
            }}>
              <span>{enabled ? '🧠' : '🧅'}</span>
              RESTART WITH {enabled ? 'AI' : 'NO AI'}
            </div>
          </div>
          <button
            onClick={onClose}
            disabled={aiRestarting}
            style={{ background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer', fontSize: 18 }}
          >
            ✕
          </button>
        </div>

        <div style={{
          background: 'rgba(255, 145, 0, 0.06)',
          border: '1px solid rgba(255, 145, 0, 0.3)',
          borderRadius: 4,
          padding: '10px 12px',
          marginBottom: 12,
          fontFamily: 'var(--font-mono)',
          fontSize: 10,
          color: 'var(--orange)',
          lineHeight: 1.6,
        }}>
          <div>⚠ This will restart the ENTIRE docker-compose stack{enabled ? ' with the AI profile (Ollama engine) enabled.' : ' without the AI engine (Ollama processing disabled).'}</div>
          <div style={{ marginTop: 6 }}>⏱ Restart takes roughly <b>15 seconds</b>. During this time the dashboard will be unreachable.</div>
        </div>

        <div style={{
          background: 'var(--bg-card)',
          border: '1px solid var(--border)',
          borderRadius: 4,
          padding: '8px 12px',
          marginBottom: 14,
          fontFamily: 'var(--font-mono)',
          fontSize: 10,
          color: 'var(--text-secondary)',
          lineHeight: 1.6,
        }}>
          🔁 <b>After restart you must reload the website</b> to reconnect to the
          freshly started back-end services. Any in-flight scans are re-queued.
        </div>

        {aiRestarting && (
          <div style={{
            padding: '8px 12px',
            borderRadius: 4,
            fontSize: 10,
            fontFamily: 'var(--font-mono)',
            background: 'rgba(0, 229, 255, 0.1)',
            color: 'var(--cyan)',
            border: '1px solid rgba(0, 229, 255, 0.3)',
            marginBottom: 14,
          }}>
            ⟳ Stack restarting... wait ~15s then reload this page.
          </div>
        )}

        <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 4 }}>
          <button
            onClick={onClose}
            disabled={aiRestarting}
            style={{
              padding: '6px 12px',
              background: 'transparent',
              border: '1px solid var(--border)',
              borderRadius: 4,
              color: 'var(--text-secondary)',
              cursor: 'pointer',
              fontFamily: 'var(--font-mono)',
              fontSize: 11,
            }}
          >
            CANCEL
          </button>
          <button
            onClick={() => onConfirm(enabled)}
            disabled={aiRestarting}
            style={{
              padding: '6px 16px',
              background: enabled ? 'var(--purple)' : 'var(--cyan)',
              border: 'none',
              borderRadius: 4,
              color: '#020406',
              fontFamily: 'var(--font-mono)',
              fontWeight: 800,
              fontSize: 11,
              cursor: aiRestarting ? 'wait' : 'pointer',
            }}
          >
            {aiRestarting ? 'RESTARTING...' : (enabled ? '⚡ RESTART WITH AI' : '↺ RESTART WITHOUT AI')}
          </button>
        </div>
      </div>
    </div>
  )
}
