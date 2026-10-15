// OSINT "jump" actions for identifier / infra graph nodes.
//
// Every extracted identifier (email, wallet, PGP fingerprint, onion host, IP,
// handle…) can be copied to the clipboard and / or opened directly in a public
// OSINT lookup tool — ipinfo.io, Shodan, Greynoise, haveibeenpwned,
// mempool.space, Ubuntu keyserver, Ahmia, t.me, Google.
//
// The dashboard-api ships each node with a raw `value` + `kind`, but
// label/address/fingerprint/id fallbacks keep mock and recorded graphs working.

const IP_RE = /^(\d{1,3}\.){3}\d{1,3}$/
const ONION_RE = /(^|\.)onion$/i
const PGP_FP_RE = /^[0-9a-f]{16,}$/i
const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[a-z]{2,}$/i
const HTTP_RE = /^https?:\/\//i
const IMAGE_EXT_RE = /\.(png|jpe?g|gif|webp|bmp|svg|ico)(\?|#|$)/i

function stripLabelPrefix(s) {
  // "BTC: 1A1z…" / "PGP: 0x…" / "IP: 1.2.3.4:8080" / "EMAIL: x@y"
  return String(s).replace(/^(BTC|XMR|LTC|EVM|ETH|PGP|IP|EMAIL|ONION|HANDLE|TELEGRAM|WALLET)[:：]\s*/i, '').trim()
}

// normIP pulls the bare IPv4 out of a value that may carry a port
// ("1.2.3.4", "1.2.3.4:8080") or a label prefix.
export function normIP(v) {
  const m = String(v || '').replace(/\s+/g, '').match(/^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::\d+)?$/)
  return m ? m[1] : ''
}

// hostOf extracts the bare host from "http://host.onion/path…" or "host.onion".
export function hostOf(v) {
  let s = String(v || '').trim()
  s = s.replace(/^[a-z][a-z0-9+.-]*:\/\//i, '')
  s = s.split('/')[0] || s
  return s.trim()
}

function ensureProtocol(h) {
  return /^[a-z][a-z0-9+.-]*:\/\//i.test(h) ? h : `http://${h}`
}

// resolveValue returns the raw searchable / copyable value for a node.
export function resolveValue(node = {}) {
  if (node && typeof node.value === 'string' && node.value.trim()) return node.value.trim()
  if (node && typeof node.ip === 'string' && node.ip.trim()) return node.ip.trim()
  if (node && typeof node.address === 'string' && node.address.trim()) return node.address.trim()
  if (node && typeof node.fingerprint === 'string' && node.fingerprint.trim()) return node.fingerprint.trim()
  if (node && typeof node.label === 'string' && node.label.trim()) {
    const cleaned = stripLabelPrefix(node.label)
    if (cleaned) return cleaned
  }
  // Fallback: id_<kind>_<value> / infra_<type>_<value> — value may contain '_'.
  if (node && typeof node.id === 'string' && /^id_|^infra_/.test(node.id)) {
    const rest = node.id.split('_').slice(2).join('_')
    if (rest.trim()) return rest.trim()
  }
  return ''
}

function kindOf(node = {}) {
  return String(node.kind || '').toLowerCase()
}

function isTelegramValue(value) {
  return /^(?:t\.me|telegram\.me|telegram\.dog)\//i.test(value) || /^@[a-z0-9_]{3,32}$/i.test(value)
}

function handleOf(value) {
  return String(value)
    .replace(/^@/, '')
    .replace(/^(?:t\.me|telegram\.me|telegram\.dog)\//i, '')
    .trim()
}

// looksLikeUrl is true for well-formed http(s) URLs and bare host[:port][/path]
// strings (e.g. "darkhydra7vxkx4d3.onion/forums").
function looksLikeUrl(val) {
  const s = String(val || '').trim()
  return s !== '' &&
    (/^https?:\/\//i.test(s) || /^[a-z0-9][a-z0-9.-]+\.[a-z]{2,}(:\d+)?(\/|$)/i.test(s))
}

// buildOSINTActions returns a list of external lookup links for a node.
// Always returns [] when no value is resolvable.
export function buildOSINTActions(node) {
  const value = resolveValue(node)
  if (!value) return []

  const type = String(node.type || '').toLowerCase()
  const kind = kindOf(node)
  const ip = normIP(value)

  const q = (term) => encodeURIComponent(term || value)
  const gsearch = (term) => ({
    id: 'search',
    label: 'GOOGLE',
    hint: 'Open web search',
    href: `https://www.google.com/search?q=${q(term)}`,
  })

  const urlSearch = (val) => ({
    id: 'url-search',
    label: 'URL SEARCH',
    hint: 'Dispatch a web search for this URL',
    href: `https://www.google.com/search?q=${q(val)}`,
  })

  const lensReverseByUrl = (src) => {
    if (!src || typeof src !== 'string') return null
    if (!/^https?:\/\//i.test(src)) return null
    if (/\.onion(:|\/|$)/i.test(src)) return null
    return `https://lens.google.com/uploadbyurl?url=${encodeURIComponent(src)}`
  }

  const blockchairSearch = (val) => ({
    id: 'blockchair',
    label: 'BLOCKCHAIR',
    hint: 'Crypto explorer (search)',
    href: `https://blockchair.com/search?q=${q(val)}`,
  })

  const mailmeteor = (email) => ({
    id: 'mailmeteor',
    label: 'MAILMETEOR',
    hint: 'Mailmeteor email checker',
    href: `https://mailmeteor.com/email-checker?email=${q(email)}`,
  })

  const whatsMyName = (handle) => ({
    id: 'whatsmyname',
    label: "WHATS MY NAME",
    hint: 'Username / alias reuse across sites',
    href: `https://whatsmyname.io/?q=${q(handle)}`,
  })

  // 1. IPv4 — the classic "copy + ipinfo" pivot.
  if (IP_RE.test(ip)) {
    return [
      { id: 'ipinfo', label: 'IPINFO.IO', hint: 'Geo / ASN / network lookup', href: `https://ipinfo.io/${ip}` },
      { id: 'shodan', label: 'SHODAN', hint: 'Port + banner exposure scan', href: `https://www.shodan.io/host/${ip}` },
      { id: 'greynoise', label: 'GREYNOISE', hint: 'Internet noise / ownership context', href: `https://viz.greynoise.io/ip/${ip}` },
      { id: 'virustotal', label: 'VIRUSTOTAL', hint: 'Threat-intel / reputation', href: `https://www.virustotal.com/gui/ip-address/${ip}` },
      { id: 'abuseipdb', label: 'ABUSEIPDB', hint: 'Abuse reports + history', href: `https://www.abuseipdb.com/check/${ip}` },
      { id: 'otx', label: 'ALIENVAULT OTX', hint: 'Open threat intel feed', href: `https://otx.alienvault.com/indicator/ip/${ip}` },
    ]
  }

  const host = hostOf(value)
  const looksOnion = ONION_RE.test(host)
  const isOnion = looksOnion || kind.includes('onion') || type === 'onion'

  // 2. Email — Mailmeteor verifier + breach exposure + domain/source pivots.
  if (kind.includes('email') || type === 'email' || EMAIL_RE.test(value)) {
    const domain = String(value).split('@')[1] || value
    return [
      mailmeteor(value),
      { id: 'hibp', label: 'HIBP', hint: 'HaveIBeenPwned breach exposure', href: `https://haveibeenpwned.com/account/${q(value)}` },
      { id: 'hunter', label: 'HUNTER', hint: 'Find more emails at this domain', href: `https://hunter.io/search/?q=${q(domain)}` },
      { id: 'dehashed', label: 'DEHASHED', hint: 'Search breached accounts', href: `https://dehashed.com/search?query=${q(value)}` },
      gsearch(`"${value}"`),
    ]
  }

  // 3. Wallets (BTC / EVM / XMR / LTC) — Blockchair is the primary explorer.
  const isBTC = kind.includes('btc') || /^(bc1|[13][a-km-zA-HJ-NP-Z1-9]{25,34})/.test(value)
  const isEVM = kind.includes('evm') || kind.includes('eth') || /^0x[0-9a-fA-F]{40}$/.test(value)
  const isXMR = kind.includes('xmr') || /^[48][0-9a-zA-Z]{90,}/.test(value)
  const isLTC = kind.includes('ltc') || /^[LM3][a-km-zA-HJ-NP-Z1-9]{26,33}$/.test(value)
  if (isBTC || isEVM || isXMR || isLTC || type === 'wallet' || kind.includes('wallet')) {
    const out = [blockchairSearch(value)]
    if (isBTC) {
      out.push({ id: 'mempool', label: 'MEMPOOL', hint: 'BTC mempool explorer', href: `https://mempool.space/address/${value}` })
      out.push({ id: 'blockstream', label: 'BLOCKSTREAM', hint: 'Open-source BTC explorer', href: `https://blockstream.info/address/${value}` })
    } else if (isEVM) {
      out.push({ id: 'etherscan', label: 'ETHERSCAN', hint: 'EVM chain explorer', href: `https://etherscan.io/address/${value}` })
      out.push({ id: 'ethplorer', label: 'ETHPLORER', hint: 'Token balances + activity', href: `https://ethplorer.io/address/${value}` })
    } else if (isXMR) {
      out.push({ id: 'xmr', label: 'XMR EXPLORER', hint: 'Monero chain explorer', href: `https://www.xmrchain.net/search?value=${q(value)}` })
    } else if (isLTC) {
      out.push({ id: 'blockchair-ltc', label: 'BLOCKCHAIR (LTC)', hint: 'LTC chain explorer', href: `https://blockchair.com/litecoin/address/${value}` })
    }
    out.push(gsearch(value))
    return out
  }

  // 4. PGP fingerprints — devglan key inspector + keyserver + web search.
  const isPGP = kind.includes('pgp') || type === 'pgp' || PGP_FP_RE.test(value.replace(/\s+/g, '')) || /^0x[0-9a-f]{8,}$/i.test(value)
  if (isPGP) {
    const fp = String(value).replace(/\s+/g, '').replace(/^0x/i, '')
    const out = [
      { id: 'devglan', label: 'DEVGLAN', hint: 'PGP key inspector', href: 'https://www.devglan.com/online-tools/pgp-key-generator' },
      { id: 'keys-openpgp', label: 'KEYS.OPENPGP.ORG', hint: 'No-JS-exfil key directory', href: `https://keys.openpgp.org/search?q=${q(fp || value)}` },
    ]
    if (PGP_FP_RE.test(fp) || /^[0-9a-f]{8,}$/i.test(fp)) {
      out.push({ id: 'keyserver', label: 'KEYSERVER', hint: 'Ubuntu keyserver lookup', href: `https://keyserver.ubuntu.com/pks/lookup?op=get&search=0x${fp}` })
      out.push({ id: 'keybase', label: 'KEYBASE', hint: 'Identity proof linking', href: `https://keybase.io/search?q=${q(fp)}` })
    }
    out.push(gsearch(fp || value))
    return out
  }

  // 5. Image nodes (direct URL or image src) — reverse-image pivots.
  if (kind.includes('image') || type === 'image' || IMAGE_EXT_RE.test(value)) {
    const href = lensReverseByUrl(value)
    const out = []
    if (href) {
      out.push({ id: 'lens', label: 'GOOGLE LENS', hint: 'Reverse image search', href })
      out.push({ id: 'yandex', label: 'YANDEX', hint: 'Yandex reverse image search', href: `https://yandex.com/images/search?rpt=imageview&url=${encodeURIComponent(value)}` })
      out.push({ id: 'tineye', label: 'TINEYE', hint: 'Reverse image search', href: `https://www.tineye.com/search?url=${encodeURIComponent(value)}` })
    }
    out.push(gsearch(value))
    return out
  }

  // 6. Telegram handles — direct profile deep link.
  if (kind.includes('telegram') || isTelegramValue(value)) {
    const h = handleOf(value)
    return [
      { id: 'telegram', label: 'TELEGRAM', hint: 'Open t.me profile', href: `https://t.me/${q(h)}` },
      whatsMyName(h),
      gsearch(`"${h}"`),
    ]
  }

  // 7. Usernames / aliases — whatsmyname.io + direct profile links.
  const isHandle =
    (kind.includes('handle') || kind.includes('alias') || kind.includes('username') ||
    type === 'actor' || type === 'identifier') && !looksLikeUrl(value)
  if (isHandle) {
    const h = handleOf(value)
    const out = [
      whatsMyName(h || value),
      gsearch(`"${h || value}"`),
      { id: 'x', label: 'X (TWITTER)', hint: 'Open profile if it exists', href: `https://x.com/${q(h || value)}` },
      { id: 'reddit', label: 'REDDIT', hint: 'Open Reddit user if it exists', href: `https://www.reddit.com/user/${q(h || value)}` },
      { id: 'github', label: 'GITHUB', hint: 'Open GitHub profile if it exists', href: `https://github.com/${q(h || value)}` },
      { id: 'idcrawl', label: 'IDCRAWL', hint: 'Aggregated persona search', href: `https://idcrawl.com/${q(h || value)}` },
    ]
    return out
  }

  // Onion services — open in Tor browser + Ahmia search.
  if (looksOnion || isOnion) {
    const onionHost = hostOf(value) || host
    return [
      { id: 'open', label: 'OPEN ONION', hint: 'Open in your Tor Browser', href: ensureProtocol(onionHost) },
      { id: 'ahmia', label: 'SEARCH URL', hint: 'Dispatch an Ahmia .onion search', href: `https://ahmia.fi/search/?q=${q(onionHost)}` },
    ]
  }

  // Non-onion URLs — OPEN + domain search + reputation pivots.
  if (looksLikeUrl(value) || type === 'url' || kind.includes('url')) {
    const hostOnly = hostOf(value)
    return [
      { id: 'open', label: 'OPEN URL', hint: 'Open this URL', href: ensureProtocol(value) },
      urlSearch(hostOnly || value),
      { id: 'virustotal', label: 'VIRUSTOTAL', hint: 'Scan / reputation check', href: `https://www.virustotal.com/gui/search/${q(value)}` },
      { id: 'urlscan', label: 'URLSCAN.IO', hint: 'Live page snapshot + related', href: `https://urlscan.io/search/#${q(hostOnly || value)}` },
      { id: 'sucuri', label: 'SUCURI', hint: 'Dispatch a live site scan', href: `https://sitecheck.sucuri.net/results/${q(hostOnly || value)}` },
      { id: 'whois', label: 'WHOIS', hint: 'Domain registration details', href: `https://whois.domaintools.com/${q(hostOnly || value)}` },
    ]
  }

  // Generic fallback
  return [gsearch(value)]
}

// scanDispatchFor returns a payload-ready investigation target when this node's
// value is a crawlable URL — an .onion hidden service or a clearweb URL. The
// dossier "DISPATCH SCAN" button submits this straight to the ingestion-api
// deep-web recon pipeline (same endpoint as the "INITIATE DEEP-WEB RECON"
// modal). Returns null for anything that cannot be handed to the crawler.
export function scanDispatchFor(node) {
  const value = resolveValue(node)
  if (!value) return null

  const type = String(node.type || '').toLowerCase()
  const kind = kindOf(node)
  const host = hostOf(value)
  const isOnion = ONION_RE.test(host) || kind.includes('onion') || type === 'onion'
  const isUrl = isOnion || looksLikeUrl(value) || type === 'url' || kind.includes('url')
  if (!isUrl) return null

  // The ingestion validator's target charset is [a-zA-Z0-9._:@/-]{1,256}, so
  // query strings / fragments are stripped and the value is length-capped.
  const raw = String(value).split(/[?#]/, 1)[0].trim()
  const targetValue = (raw.length > 256 ? raw.slice(0, 256) : raw) || host

  return {
    target_type: isOnion ? 'onion_address' : 'url',
    target_value: targetValue,
    seed_terms: [host || targetValue],
  }
}