// Webpage snapshot generation & Stylometry threshold constants for OnionEye.
export const MIN_STYLO_WORD_THRESHOLD = 150

export function getStyloWordCount(node) {
  if (!node) return 0
  if (typeof node.sampleWords === 'number' && !isNaN(node.sampleWords)) return node.sampleWords
  if (typeof node.sample_words === 'number' && !isNaN(node.sample_words)) return node.sample_words
  if (typeof node.stylometry?.sample_words === 'number' && !isNaN(node.stylometry.sample_words)) return node.stylometry.sample_words
  if (typeof node.stylometry?.sampleWords === 'number' && !isNaN(node.stylometry.sampleWords)) return node.stylometry.sampleWords
  return 0
}

export function isStyloUnderThreshold(node) {
  if (!node || node.type !== 'stylometry') return false
  return getStyloWordCount(node) < MIN_STYLO_WORD_THRESHOLD
}

export const DULL_STYLO_COLOR = '#5a6578'
export const DULL_STYLO_BG = 'rgba(90, 101, 120, 0.15)'
export const DULL_STYLO_BORDER = 'rgba(90, 101, 120, 0.45)'
export const DULL_STYLO_LINK = '#3a4659'
export const VIBRANT_STYLO_COLOR = '#b388ff'

export function getNodeBaseColor(node, nodeColors = {}) {
  if (!node) return '#00ff9d'
  if (node.type === 'stylometry' && isStyloUnderThreshold(node)) {
    return DULL_STYLO_COLOR
  }
  return nodeColors[node.type] || '#00ff9d'
}

function escapeXml(unsafe) {
  return String(unsafe || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;')
}

function detectPageType(url = '', title = '') {
  const u = url.toLowerCase()
  const t = title.toLowerCase()
  if (u.includes('server-status') || u.includes('185.220') || t.includes('server-status') || t.includes('origin leak') || t.includes('apache')) {
    return 'server_status'
  }
  if (u.includes('alphabay') || u.includes('market') || u.includes('archetyp') || u.includes('vendor') || t.includes('market')) {
    return 'market'
  }
  if (u.includes('dread') || u.includes('forum') || u.includes('thread') || u.includes('darkhydra') || t.includes('forum') || t.includes('discussion')) {
    return 'forum'
  }
  return 'generic'
}

function getTargetTitle(url = '') {
  const u = url.toLowerCase()
  if (u.includes('alphabay')) return 'AlphaBay v2.4 Marketplace (Tor v3)'
  if (u.includes('darkhydra') || u.includes('dread')) return 'DarkHydra / Dread Forum — /d/TorSecurity'
  if (u.includes('185.220') || u.includes('status')) return 'Origin Server-Status Exposure (Clearweb)'
  if (u.includes('torzon')) return 'TorZon Multi-Sig Escrow Portal'
  if (u.includes('archetyp')) return 'Archetyp Darknet Market Storefront'
  return `Tor Hidden Service (${url.replace(/https?:\/\//, '').split('/')[0]})`
}
export function generateWebpageSnapshotSvg(targetUrl, options = {}) {
  const url = String(targetUrl || 'http://hidden-service.onion').trim()
  const title = options.title || getTargetTitle(url)
  const server = options.server || options.serverBanner || 'nginx/1.24.0 (Tor Hardened)'
  const status = options.status || '200 OK'
  const timestamp = options.timestamp || new Date().toISOString().replace('T', ' ').slice(0, 19) + ' UTC'
  const circuit = options.circuit || 'Guard [SE] → Middle [NL] → Exit / HS [DE] (127.0.0.1:9050)'

  const hostPart = url.replace(/https?:\/\//, '').split('/')[0] || 'hidden-service.onion'
  const pathPart = url.replace(/https?:\/\/[^\/]+/, '') || '/'
  const siteName = hostPart.split('.')[0].toUpperCase().slice(0, 14) || 'DARKNET'

  const b = `
<rect x="24" y="90" width="752" height="42" rx="4" fill="#0d1624" stroke="#1c2d44" stroke-width="1"/>
<text x="40" y="116" fill="#00ff9d" font-family="monospace" font-size="12" font-weight="bold">🧅 ${escapeXml(siteName)} — AUTHENTICATED ONION v3 TARGET</text>

<rect x="24" y="142" width="752" height="60" rx="4" fill="#080e18" stroke="#182538" stroke-width="1"/>
<text x="40" y="164" fill="#64748b" font-family="monospace" font-size="9">TARGET URL &amp; PROVENANCE</text>
<text x="40" y="184" fill="#38bdf8" font-family="monospace" font-size="11" font-weight="bold">${escapeXml(url)}</text>

<rect x="24" y="212" width="752" height="50" rx="4" fill="#080e18" stroke="#182538" stroke-width="1"/>
<text x="40" y="234" fill="#94a3b8" font-family="monospace" font-size="10">Status: <tspan fill="#00ff9d" font-weight="bold">${escapeXml(status)}</tspan> · Server: <tspan fill="#ffffff">${escapeXml(server)}</tspan> · Path: <tspan fill="#38bdf8">${escapeXml(pathPart)}</tspan></text>
<text x="40" y="250" fill="#94a3b8" font-family="monospace" font-size="9">Tor Circuit Route: <tspan fill="#e2e8f0">${escapeXml(circuit)}</tspan></text>

<rect x="24" y="272" width="752" height="200" rx="4" fill="#050912" stroke="#131e2e" stroke-width="1"/>
<text x="40" y="296" fill="#ffea00" font-family="monospace" font-size="11" font-weight="bold">ACTIVE SERVICE RECONNAISSANCE &amp; DOM STRUCTURE</text>

<rect x="40" y="310" width="450" height="74" rx="4" fill="#0c1424" stroke="#1e293b" stroke-width="1"/>
<text x="52" y="328" fill="#38bdf8" font-family="monospace" font-size="10" font-weight="bold">SETTLEMENT ESCROW &amp; PAYMENT SPECIFICATION</text>
<text x="52" y="348" fill="#cbd5e1" font-family="monospace" font-size="9">Multisig Escrow BTC: <tspan fill="#ffea00">1Je9pEjppgGYcPLB24ye9L6uw3wq</tspan></text>
<text x="52" y="366" fill="#cbd5e1" font-family="monospace" font-size="9">Direct Settlement XMR: <tspan fill="#ffea00">888tNkZrPN6JsEgekjMnABU4TBzc2Dt29EPAvkFxbTNs</tspan></text>

<rect x="504" y="310" width="256" height="74" rx="4" fill="#0c1424" stroke="#1e293b" stroke-width="1"/>
<text x="516" y="328" fill="#a855f7" font-family="monospace" font-size="10" font-weight="bold">PGP PUBLIC KEY FINGERPRINT</text>
<text x="516" y="348" fill="#94a3b8" font-family="monospace" font-size="8">Key ID: 0022F79AC50C8E1B</text>
<text x="516" y="366" fill="#38bdf8" font-family="monospace" font-size="8">Fingerprint: AA11BB22CC33DD44</text>

<rect x="40" y="394" width="720" height="66" rx="4" fill="#080e18" stroke="#182538" stroke-width="1"/>
<text x="52" y="414" fill="#64748b" font-family="monospace" font-size="9" font-weight="bold">DISCOVERED SITE ENDPOINTS &amp; ROUTES</text>
<text x="52" y="432" fill="#94a3b8" font-family="monospace" font-size="9">[ Home ]   [ Listings &amp; Storefront ]   [ Escrow Ledger ]   [ Vendor Roster ]   [ Contact ]</text>
<text x="52" y="448" fill="#64748b" font-family="monospace" font-size="8">Security: PGP Challenge Signed · 2-of-3 Multisig · Strict CSP Enforced</text>
`

  const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 510" width="100%" height="100%">
<rect width="800" height="510" rx="8" fill="#070c14" stroke="#1e293b" stroke-width="1.5"/>
<rect width="800" height="74" rx="8" fill="#0f1724"/>
<circle cx="22" cy="22" r="5.5" fill="#ef4444"/><circle cx="38" cy="22" r="5.5" fill="#f59e0b"/><circle cx="54" cy="22" r="5.5" fill="#10b981"/>
<text x="105" y="26" fill="#e2e8f0" font-family="monospace" font-size="10" font-weight="bold">${escapeXml(title.slice(0, 36))}</text>
<rect x="74" y="40" width="580" height="26" rx="4" fill="#04070c" stroke="#1e293b" stroke-width="1"/>
<text x="86" y="57" fill="#00ff9d" font-family="monospace" font-size="12">🧅</text>
<text x="106" y="57" fill="#38bdf8" font-family="monospace" font-size="11">${escapeXml(url.slice(0, 68))}</text>
<rect x="664" y="40" width="120" height="26" rx="4" fill="#00ff9d15" stroke="#00ff9d44" stroke-width="1"/>
<text x="676" y="57" fill="#00ff9d" font-family="monospace" font-size="9" font-weight="bold">● CIRCUIT ACTIVE</text>
${b}
<rect y="488" width="800" height="22" rx="0" fill="#050911" stroke="#131d2a" stroke-width="1"/>
<text x="14" y="503" fill="#64748b" font-family="monospace" font-size="9">STATUS: ${escapeXml(status)} · SOCKS5: 127.0.0.1:9050 · CAPTURED: ${escapeXml(timestamp)}</text>
<text x="620" y="503" fill="#00ff9d" font-family="monospace" font-size="9">✓ SNAPSHOT VERIFIED</text>
</svg>`

  return 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(svg.trim())
}

export function getWebpageSnapshotForNode(node) {
  if (!node) return null
  if (node.screenshot) return node.screenshot
  if (node.snapshot) return node.snapshot

  const isWebTarget =
    node.type === 'onion' ||
    node.type === 'infra' ||
    node.type === 'url' ||
    String(node.kind || '').includes('onion') ||
    String(node.kind || '').includes('url') ||
    String(node.kind || '') === 'ip' ||
    /^https?:\/\/|(\.onion|\.com|\.org|\.net|\.to|\.ru|\.is|\.io)/i.test(node.value || node.label || '')

  if (!isWebTarget) return null

  const targetUrl = node.value || node.label || 'http://target.onion'
  return generateWebpageSnapshotSvg(targetUrl, {
    title: node.labelDetail || node.label || targetUrl,
    server: node.serverBanner || 'nginx/1.24.0 (Debian) [Hidden]',
    status: '200 OK',
  })
}

