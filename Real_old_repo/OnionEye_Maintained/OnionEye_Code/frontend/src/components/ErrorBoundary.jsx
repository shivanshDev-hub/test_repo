import React from 'react'

export default class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props)
    this.state = { hasError: false, error: null, errorInfo: null }
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error }
  }

  componentDidCatch(error, errorInfo) {
    console.error('OnionEye Dashboard Error caught by ErrorBoundary:', error, errorInfo)
    this.setState({ errorInfo })
  }

  handleReload = () => {
    window.location.reload()
  }

  handleReset = () => {
    this.setState({ hasError: false, error: null, errorInfo: null })
  }

  render() {
    if (this.state.hasError) {
      return (
        <div style={{
          height: '100vh',
          width: '100vw',
          background: 'radial-gradient(circle at center, #0d1522 0%, #03060a 100%)',
          color: 'var(--text-primary)',
          fontFamily: 'var(--font-mono)',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          padding: 24,
          boxSizing: 'border-box',
          textAlign: 'center',
        }}>
          <div style={{
            background: 'rgba(5, 8, 14, 0.95)',
            border: '1px solid rgba(255, 23, 68, 0.5)',
            boxShadow: '0 0 40px rgba(255, 23, 68, 0.2)',
            borderRadius: 8,
            padding: '24px 32px',
            maxWidth: 640,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: 16,
          }}>
            <div style={{ fontSize: 32 }}>⚠️</div>
            <div style={{ fontSize: 14, fontWeight: 800, color: 'var(--red)', letterSpacing: '0.1em' }}>
              RUNTIME ANOMALY INTERCEPTED
            </div>
            <div style={{ fontSize: 11, color: 'var(--text-secondary)', lineHeight: 1.6 }}>
              A UI rendering fault occurred. The dashboard engine caught the exception to prevent crash propagation.
            </div>
            {this.state.error && (
              <pre style={{
                background: 'rgba(0,0,0,0.6)',
                border: '1px solid var(--border)',
                borderRadius: 4,
                padding: '8px 12px',
                fontSize: 10,
                color: 'var(--yellow)',
                maxWidth: '100%',
                overflowX: 'auto',
                textAlign: 'left',
              }}>
                {this.state.error.toString()}
              </pre>
            )}
            <div style={{ display: 'flex', gap: 12, marginTop: 8 }}>
              <button
                onClick={this.handleReset}
                style={{
                  padding: '8px 16px',
                  background: 'rgba(0, 255, 157, 0.1)',
                  border: '1px solid var(--green)',
                  borderRadius: 4,
                  color: 'var(--green)',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 11,
                  fontWeight: 800,
                  cursor: 'pointer',
                }}
              >
                ↺ RESUME HUD
              </button>
              <button
                onClick={this.handleReload}
                style={{
                  padding: '8px 16px',
                  background: 'rgba(255, 255, 255, 0.05)',
                  border: '1px solid var(--border)',
                  borderRadius: 4,
                  color: 'var(--text-primary)',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 11,
                  cursor: 'pointer',
                }}
              >
                RELOAD DASHBOARD
              </button>
            </div>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}
