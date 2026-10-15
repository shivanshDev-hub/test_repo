import React, { useState } from 'react'
import { submitInvestigation } from '../services/api'

export default function NewScanModal({ isOpen, onClose, onJobCreated }) {
  const [targetType, setTargetType] = useState('onion_address')
  const [targetValue, setTargetValue] = useState('darkhydra7vxkx4d3.onion')
  const [seedTerms, setSeedTerms] = useState('0x_shadow_lord, escrow, exploit, pgp')
  const [depth, setDepth] = useState(3)
  const [loading, setLoading] = useState(false)
  const [resultMsg, setResultMsg] = useState(null)

  if (!isOpen) return null

  async function handleSubmit(e) {
    e.preventDefault()
    setLoading(true)
    const payload = {
      analyst_id: 'analyst_ui_01',
      target_type: targetType,
      target_value: targetValue.trim(),
      seed_terms: seedTerms.split(',').map(s => s.trim()).filter(Boolean),
      max_depth: Number(depth),
    }
    const res = await submitInvestigation(payload)
    setLoading(false)
    if (res.ok) {
      const llmNote = res.data?.llm_unavailable ? ' (regex-only, LLM offline)' : ''
      setResultMsg({ type: 'success', text: `[+] Task Enqueued to Tor Pool: ${res.data.job_id || 'ACTIVE'}${llmNote}` })
    } else {
      setResultMsg({ type: 'info', text: `[!] Backend unreachable — queued locally: ${targetValue}` })
    }
    if (onJobCreated) onJobCreated(payload)
    setTimeout(() => { onClose(); setResultMsg(null) }, 1400)
  }

  return (
    <div style={{
      position: 'fixed',
      inset: 0,
      zIndex: 999,
      background: 'rgba(2, 4, 6, 0.88)',
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
        maxWidth: 460,
        padding: 20,
        boxShadow: '0 20px 60px rgba(0,0,0,0.9), 0 0 30px rgba(0, 255, 157, 0.1)',
        position: 'relative',
      }}>
        {/* Top Warning Banner */}
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
              // TACTICAL INGESTION DISPATCH //
            </div>
            <div style={{
              fontFamily: 'var(--font-mono)',
              fontSize: 13,
              color: 'var(--green)',
              fontWeight: 800,
              display: 'flex',
              alignItems: 'center',
              gap: 6,
              marginTop: 2,
            }}>
              <span>⚡</span> INITIATE DEEP-WEB RECON
            </div>
          </div>
          <button
            onClick={onClose}
            style={{
              background: 'none',
              border: 'none',
              color: 'var(--text-muted)',
              cursor: 'pointer',
              fontSize: 18,
            }}
          >
            ✕
          </button>
        </div>
        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          <div>
            <label style={{ display: 'block', fontSize: 10, fontFamily: 'var(--font-mono)', color: 'var(--text-secondary)', marginBottom: 4 }}>
              TARGET CLASSIFICATION
            </label>
            <select
              value={targetType}
              onChange={e => setTargetType(e.target.value)}
              style={{
                width: '100%',
                padding: '6px 8px',
                background: 'var(--bg-card)',
                border: '1px solid var(--border-bright)',
                borderRadius: 4,
                color: 'var(--text-primary)',
                fontFamily: 'var(--font-mono)',
                fontSize: 11,
                outline: 'none',
              }}
            >
              <option value="onion_address">🧅 Tor v3 Hidden Service (.onion)</option>
              <option value="vendor_handle">👤 Dark Market Vendor Handle</option>
              <option value="btc_wallet">₿ Bitcoin / Monero Crypto Address</option>
              <option value="pgp_key">🔑 PGP Key Fingerprint</option>
            </select>
          </div>

          <div>
            <label style={{ display: 'block', fontSize: 10, fontFamily: 'var(--font-mono)', color: 'var(--text-secondary)', marginBottom: 4 }}>
              TARGET VALUE
            </label>
            <input
              type="text"
              value={targetValue}
              onChange={e => setTargetValue(e.target.value)}
              required
              style={{
                width: '100%',
                padding: '6px 8px',
                background: 'var(--bg-card)',
                border: '1px solid var(--border-bright)',
                borderRadius: 4,
                color: 'var(--green)',
                fontFamily: 'var(--font-mono)',
                fontSize: 11,
                outline: 'none',
              }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: 10, fontFamily: 'var(--font-mono)', color: 'var(--text-secondary)', marginBottom: 4 }}>
              SEED EXPANSION KEYWORDS
            </label>
            <input
              type="text"
              value={seedTerms}
              onChange={e => setSeedTerms(e.target.value)}
              style={{
                width: '100%',
                padding: '6px 8px',
                background: 'var(--bg-card)',
                border: '1px solid var(--border-bright)',
                borderRadius: 4,
                color: 'var(--text-primary)',
                fontFamily: 'var(--font-mono)',
                fontSize: 11,
                outline: 'none',
              }}
            />
          </div>

          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
              <label style={{ fontSize: 10, fontFamily: 'var(--font-mono)', color: 'var(--text-secondary)' }}>
                CRAWLER DEPTH: {depth} HOPS
              </label>
            </div>
            <input
              type="range"
              min="1"
              max="5"
              value={depth}
              onChange={e => setDepth(e.target.value)}
              style={{ width: '100%', accentColor: 'var(--green)', cursor: 'pointer' }}
            />
          </div>

          <div style={{
            background: 'rgba(0, 229, 255, 0.04)',
            border: '1px solid rgba(0, 229, 255, 0.2)',
            borderRadius: 4,
            padding: '6px 8px',
            fontFamily: 'var(--font-mono)',
            fontSize: 9,
            color: 'var(--cyan)',
          }}>
            🛡 ZERO-TRUST: Content Shield short-circuits any prohibited data prior to storage.
          </div>

          {resultMsg && (
            <div style={{
              padding: '6px 8px',
              borderRadius: 4,
              fontSize: 10,
              fontFamily: 'var(--font-mono)',
              background: resultMsg.type === 'success' ? 'rgba(0, 255, 157, 0.1)' : 'rgba(0, 180, 255, 0.1)',
              color: resultMsg.type === 'success' ? 'var(--green)' : 'var(--cyan)',
              border: `1px solid ${resultMsg.type === 'success' ? 'var(--green)' : 'var(--cyan)'}`,
            }}>
              {resultMsg.text}
            </div>
          )}

          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 4 }}>
            <button
              type="button"
              onClick={onClose}
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
              type="submit"
              disabled={loading}
              style={{
                padding: '6px 16px',
                background: 'var(--green)',
                border: 'none',
                borderRadius: 4,
                color: '#020406',
                fontFamily: 'var(--font-mono)',
                fontWeight: 800,
                fontSize: 11,
                cursor: loading ? 'wait' : 'pointer',
              }}
            >
              {loading ? 'ENQUEUING...' : '⚡ DISPATCH'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

