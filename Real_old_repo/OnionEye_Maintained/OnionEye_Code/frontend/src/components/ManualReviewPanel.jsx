import React, { useState, useEffect } from 'react'
import { fetchManualReviewQueue, dismissReviewItem } from '../services/api'

export default function ManualReviewPanel({ isOpen, onClose }) {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(false)
  const [dismissing, setDismissing] = useState(null)

  useEffect(() => {
    if (!isOpen) return
    loadQueue()
    const interval = setInterval(loadQueue, 10000) // Refresh every 10s
    return () => clearInterval(interval)
  }, [isOpen])

  async function loadQueue() {
    setLoading(true)
    const res = await fetchManualReviewQueue()
    setLoading(false)
    if (res.data?.items) {
      // Filter out expired items
      setItems(res.data.items.filter(it => !it.expired))
    }
  }

  async function handleDismiss(item) {
    setDismissing(item.id)
    const res = await dismissReviewItem(item.id)
    setDismissing(null)
    if (res.ok) {
      setItems(items.filter(it => it.id !== item.id))
    } else {
      alert(`Failed to dismiss: ${res.error || 'Unknown error'}`)
    }
  }

  if (!isOpen) return null

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      background: 'rgba(0, 0, 0, 0.85)',
      zIndex: 9999,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      backdropFilter: 'blur(4px)',
    }}>
      <div style={{
        background: 'var(--bg-secondary)',
        border: '1px solid var(--border)',
        borderRadius: 6,
        width: '90%',
        maxWidth: 700,
        maxHeight: '80vh',
        display: 'flex',
        flexDirection: 'column',
        boxShadow: '0 8px 32px rgba(0, 0, 0, 0.6)',
      }}>
        <div style={{
          padding: '16px 20px',
          borderBottom: '1px solid var(--border)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
        }}>
          <div>
            <div style={{
              fontSize: 15,
              fontWeight: 700,
              color: 'var(--text)',
              fontFamily: 'var(--font-mono)',
              marginBottom: 4,
            }}>
              🔍 MANUAL REVIEW QUEUE
            </div>
            <div style={{
              fontSize: 11,
              color: 'var(--text-muted)',
              fontFamily: 'var(--font-mono)',
            }}>
              Sites flagged for human intervention (CAPTCHA, bot blocks, scan failures)
            </div>
          </div>
          <button
            onClick={onClose}
            style={{
              background: 'transparent',
              border: '1px solid var(--border)',
              borderRadius: 3,
              padding: '6px 12px',
              color: 'var(--text-secondary)',
              cursor: 'pointer',
              fontFamily: 'var(--font-mono)',
              fontSize: 11,
              fontWeight: 600,
            }}
          >
            ✕ CLOSE
          </button>
        </div>

        <div style={{
          flex: 1,
          overflowY: 'auto',
          padding: 16,
        }}>
          {loading && items.length === 0 && (
            <div style={{
              textAlign: 'center',
              padding: 40,
              color: 'var(--text-muted)',
              fontFamily: 'var(--font-mono)',
              fontSize: 11,
            }}>
              Loading review queue...
            </div>
          )}

          {!loading && items.length === 0 && (
            <div style={{
              textAlign: 'center',
              padding: 40,
              color: 'var(--text-muted)',
              fontFamily: 'var(--font-mono)',
              fontSize: 11,
            }}>
              ✓ No pending items — all clear!
            </div>
          )}

          {items.map(item => (
            <div key={item.id} style={{
              background: 'rgba(12, 20, 32, 0.5)',
              border: '1px solid var(--border)',
              borderRadius: 4,
              padding: 12,
              marginBottom: 12,
              display: 'flex',
              gap: 12,
              alignItems: 'flex-start',
            }}>
              <div style={{ flex: 1 }}>
                <div style={{
                  fontSize: 12,
                  fontWeight: 600,
                  color: 'var(--text)',
                  fontFamily: 'var(--font-mono)',
                  marginBottom: 6,
                  wordBreak: 'break-all',
                }}>
                  {item.target}
                </div>
                <div style={{
                  fontSize: 10,
                  color: 'var(--text-muted)',
                  fontFamily: 'var(--font-mono)',
                  marginBottom: 4,
                }}>
                  <span style={{ color: 'var(--orange)', fontWeight: 600 }}>
                    {item.reason || 'BLOCKED'}
                  </span>
                  {' · '}
                  <span>Job {item.job_id}</span>
                  {' · '}
                  <span>{new Date(item.flagged_at).toLocaleString()}</span>
                </div>
                <div style={{
                  fontSize: 9,
                  color: 'var(--text-muted)',
                  fontFamily: 'var(--font-mono)',
                }}>
                  5-minute expiry slot per architecture
                </div>
              </div>
              <button
                onClick={() => handleDismiss(item)}
                disabled={dismissing === item.id}
                style={{
                  background: dismissing === item.id ? 'rgba(0, 229, 255, 0.1)' : 'rgba(0, 229, 255, 0.15)',
                  border: '1px solid var(--cyan)',
                  borderRadius: 3,
                  padding: '6px 12px',
                  color: 'var(--cyan)',
                  cursor: dismissing === item.id ? 'wait' : 'pointer',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 10,
                  fontWeight: 600,
                  whiteSpace: 'nowrap',
                }}
              >
                {dismissing === item.id ? '...' : '✓ DISMISS'}
              </button>
            </div>
          ))}
        </div>

        <div style={{
          padding: '12px 20px',
          borderTop: '1px solid var(--border)',
          fontSize: 10,
          color: 'var(--text-muted)',
          fontFamily: 'var(--font-mono)',
          textAlign: 'center',
        }}>
          {items.length > 0 && `${items.length} pending ${items.length === 1 ? 'item' : 'items'}`}
          {items.length === 0 && 'Queue refreshes every 10 seconds'}
        </div>
      </div>
    </div>
  )
}
