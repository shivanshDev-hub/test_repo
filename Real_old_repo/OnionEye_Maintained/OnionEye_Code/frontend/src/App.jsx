import React, { useState, useRef, useCallback, useEffect, useMemo } from 'react'
import ForceGraph2D from 'react-force-graph-2d'
import { forceCollide } from 'd3-force-3d'
import { fetchStats, fetchGraph, fetchLeads, fetchMedia, deleteNode, deleteAll, fetchAiStatus, fetchTorStatus, restartStackWithAi, fetchNotifications } from './services/api'
import Header from './components/Header'
import StatsBar from './components/StatsBar'
import LeftControlRail from './components/LeftControlRail'
import DetailPanel from './components/DetailPanel'
import NewScanModal from './components/NewScanModal'
import NotificationToast from './components/NotificationToast'
import RightRail from './components/RightRail'
import { isStyloUnderThreshold, DULL_STYLO_COLOR, MIN_STYLO_WORD_THRESHOLD } from './lib/snapshots'

const NODE_COLORS = {
  actor: '#00ff9d',
  pgp: '#d500f9',
  wallet: '#ff9100',
  onion: '#00b0ff',
  infra: '#ff1744',
  email: '#ffea00',
  stylometry: '#b388ff',
}

const TYPE_ICONS = {
  actor: '👤',
  pgp: '🔑',
  wallet: '₿',
  onion: '🧅',
  infra: '🖥',
  email: '✉',
  stylometry: '✍',
}

const LINK_COLORS = {
  persona: '#00ff9d',
  identifier: '#d500f9',
  infra: '#ff1744',
  stylometry: '#b388ff',
}

// ── Label truncation ───────────────────────────────────────────────────────────
// Long node names (onion hosts, handles, pgp fingerprints…) clutter the graph.
// Default rendering shows only a short prefix; the FULL name is drawn only when
// the cursor is hovering the node (or it is selected/focused).
const SHORT_LABEL_MAX = 14
const CANVAS_FONT_STACK = "'JetBrains Mono', 'Fira Code', 'Courier New', monospace"
function truncateLabel(text, max = SHORT_LABEL_MAX) {
  const s = String(text || '')
  return s.length > max ? `${s.slice(0, max)}…` : s
}

// Seed-expansion keywords are investigation search terms, not extracted handles.
// Backend fallback used to mint "TERM:…" actor nodes from them; hide those so they
// never appear as user personas on the radar.
function isSeedKeywordNode(node) {
  if (!node) return false
  const cat = String(node.category || node.kind || '').toUpperCase()
  if (cat === 'TERM_FALLBACK' || cat === 'SEED' || cat === 'SEED_TERM') return true
  const label = String(node.label || node.value || node.id || '')
  if (/^TERM:/i.test(label)) return true
  return false
}

function nodeFootprint(node) {
  const type = node?.type
  if (type === 'actor') return { kind: 'rect', w: 46, h: 26, r: 18, corner: 8 }
  if (type === 'stylometry') return { kind: 'triangle', r: 14 }
  return { kind: 'circle', r: 9 }
}

function drawRoundedRectPath(ctx, cx, cy, w, h, radius) {
  const x = cx - w / 2
  const y = cy - h / 2
  const r = Math.min(radius, w / 2, h / 2)
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.lineTo(x + w - r, y)
  ctx.quadraticCurveTo(x + w, y, x + w, y + r)
  ctx.lineTo(x + w, y + h - r)
  ctx.quadraticCurveTo(x + w, y + h, x + w - r, y + h)
  ctx.lineTo(x + r, y + h)
  ctx.quadraticCurveTo(x, y + h, x, y + h - r)
  ctx.lineTo(x, y + r)
  ctx.quadraticCurveTo(x, y, x + r, y)
  ctx.closePath()
}

function drawTrianglePath(ctx, cx, cy, r) {
  ctx.beginPath()
  ctx.moveTo(cx, cy - r)
  ctx.lineTo(cx + r * 0.96, cy + r * 0.78)
  ctx.lineTo(cx - r * 0.96, cy + r * 0.78)
  ctx.closePath()
}

function drawNodeShapePath(ctx, node, pad = 0) {
  const fp = nodeFootprint(node)
  if (fp.kind === 'rect') {
    drawRoundedRectPath(ctx, node.x, node.y, fp.w + pad * 2, fp.h + pad * 2, fp.corner + pad * 0.4)
    return fp
  }
  if (fp.kind === 'triangle') {
    drawTrianglePath(ctx, node.x, node.y, fp.r + pad)
    return fp
  }
  ctx.beginPath()
  ctx.arc(node.x, node.y, fp.r + pad, 0, 2 * Math.PI)
  return fp
}

function nodeLabelOffset(node) {
  const fp = nodeFootprint(node)
  if (fp.kind === 'rect') return fp.h / 2
  if (fp.kind === 'triangle') return fp.r * 0.82
  return fp.r
}


// ── Connected-component clustering ─────────────────────────────────────────────
// OSINT case data naturally forms islands (one actor + its wallets/onions/pgp per
// ring or market). Splitting the graph into connected components and seeding each
// one as a separate mini-constellation turns “one dense ball” into a set of
// readable, separated case clusters.
function computeClusters(nodes, links) {
  const parent = new Map()
  const find = (id) => {
    while (parent.get(id) !== id) {
      parent.set(id, parent.get(parent.get(id)))
      id = parent.get(id)
    }
    return id
  }
  const union = (a, b) => {
    const ra = find(a)
    const rb = find(b)
    if (ra !== rb) parent.set(rb, ra)
  }
  nodes.forEach(n => parent.set(n.id, n.id))
  ;(links || []).forEach(l => {
    const s = typeof l.source === 'object' ? l.source.id : l.source
    const t = typeof l.target === 'object' ? l.target.id : l.target
    if (parent.has(s) && parent.has(t)) union(s, t)
  })
  const byRoot = new Map()
  nodes.forEach(n => {
    const r = find(n.id)
    if (!byRoot.has(r)) byRoot.set(r, [])
    byRoot.get(r).push(n)
  })
  return [...byRoot.values()].map(members => ({ members }))
}

function linkEndpoints(link) {
  const s = typeof link?.source === 'object' ? link.source.id : link?.source
  const t = typeof link?.target === 'object' ? link.target.id : link?.target
  return [s, t]
}

// Topology fingerprint: nodes + edges. Used to ignore identical 10s poll
// payloads so the force simulation is not re-ingested / reheated.
function graphTopologyKey(data) {
  const nodes = data?.nodes || []
  const links = data?.links || []
  const nodePart = nodes.map(n => n.id).filter(Boolean).sort().join('|')
  const linkPart = links.map(l => {
    const [s, t] = linkEndpoints(l)
    return `${s}>${t}:${l.type || ''}`
  }).sort().join('|')
  return `${nodePart}::${linkPart}`
}

function preserveSimState(target, incoming) {
  const sim = {
    x: target.x,
    y: target.y,
    vx: target.vx,
    vy: target.vy,
    fx: target.fx,
    fy: target.fy,
  }
  Object.assign(target, incoming)
  if (typeof sim.x === 'number' && isFinite(sim.x)) target.x = sim.x
  if (typeof sim.y === 'number' && isFinite(sim.y)) target.y = sim.y
  if (typeof sim.vx === 'number' && isFinite(sim.vx)) target.vx = sim.vx
  if (typeof sim.vy === 'number' && isFinite(sim.vy)) target.vy = sim.vy
  target.fx = sim.fx
  target.fy = sim.fy
  return target
}

// Reuse live node/link objects across polls. react-force-graph-2d restarts
// the d3 simulation whenever it receives a brand-new graphData identity
// with fresh node objects (no x/y). Returning the previous object when
// topology is unchanged keeps the constellation frozen; property updates
// (label, flagged, confidence) are copied onto the existing objects.
function mergeGraphSnapshot(prev, next) {
  const prevNodes = prev?.nodes || []
  const nextNodes = Array.isArray(next?.nodes) ? next.nodes : []
  const prevLinks = prev?.links || []
  const nextLinks = Array.isArray(next?.links) ? next.links : []
  const prevKey = graphTopologyKey(prev)
  const nextKey = graphTopologyKey(next)

  const prevById = new Map(prevNodes.map(n => [n.id, n]))
  const prevLinkByKey = new Map()
  for (const l of prevLinks) {
    const [s, t] = linkEndpoints(l)
    prevLinkByKey.set(`${s}|${t}|${l.type || ''}`, l)
  }

  if (prevKey === nextKey && prevNodes.length) {
    for (const incoming of nextNodes) {
      const existing = prevById.get(incoming.id)
      if (existing) preserveSimState(existing, incoming)
    }
    for (const incoming of nextLinks) {
      const [s, t] = linkEndpoints(incoming)
      const existing = prevLinkByKey.get(`${s}|${t}|${incoming.type || ''}`)
      if (!existing) continue
      const src = existing.source
      const tgt = existing.target
      Object.assign(existing, incoming)
      existing.source = src
      existing.target = tgt
    }
    return prev
  }

  const nodes = nextNodes.map(incoming => {
    const existing = prevById.get(incoming.id)
    if (!existing) return { ...incoming }
    return preserveSimState(existing, incoming)
  })

  const links = nextLinks.map(incoming => {
    const [s, t] = linkEndpoints(incoming)
    const existing = prevLinkByKey.get(`${s}|${t}|${incoming.type || ''}`)
    if (existing) {
      const src = existing.source
      const tgt = existing.target
      Object.assign(existing, incoming)
      existing.source = src
      existing.target = tgt
      return existing
    }
    return { ...incoming, source: s, target: t }
  })

  return { ...next, nodes, links }
}

// ── Stylokey → 2D layout reduction ──────────────────────────────────────────
// The backend projects every actor's lexical style vector + posting rhythm
// histograms into a fixed 64-d deterministic key (correlation-engine's
// pkg/stylokey: normalize → weighted fusion → ±1 random projection → L2).
// Summing the even/odd components back into a 2D offset is a stable linear
// reduction: actors whose writing style + tempo are similar produce nearby
// keys, and therefore nearby seeded positions. Non-actor nodes have no key
// (they orbit their owning actor instead).
function stylokeyTo2D(sk) {
  if (!Array.isArray(sk) || sk.length < 2) return null
  let x = 0
  let y = 0
  for (let j = 0; j < sk.length; j++) {
    if (j % 2 === 0) x += sk[j]
    else y += sk[j]
  }
  return { x, y }
}

// ── Node-type derivation for delete ────────────────────────────────────────────
// storage-api DELETE /delete/node requires a non-empty node_type plus the UI
// graph node id. The /graph endpoint normally populates `node.type`, but some
// nodes (e.g. freshly scraped/ingested rows) can arrive without it, so we fall
// back to parsing the id. Graph ids are built by storage-api as:
//   actor:      bare UUID
//   identifier: "id_" + kind + "_" + value   (kind = wallet:BTC / email / ...)
//   infra:      "infra_" + indicatorType + "_" + onionHost
const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

// kind segment (after "id_") -> graph node type used by the UI + storage-api.
const IDENTIFIER_KIND_TO_TYPE = {
  email: 'email',
  pgp: 'pgp',
  pgp_key: 'pgp',
  onion: 'onion',
  site: 'onion',
  wallet: 'wallet',
  btc: 'wallet',
  xmr: 'wallet',
  ltc: 'wallet',
  eth: 'wallet',
  evm: 'wallet',
  usdt: 'wallet',
  bch: 'wallet',
  zec: 'wallet',
  stylometry: 'stylometry',
  stylometry_vector: 'stylometry',
  stylo: 'stylometry',
}

function deriveNodeType(node) {
  if (!node) return null
  if (node.type && typeof node.type === 'string') return node.type

  const id = typeof node.id === 'string' ? node.id : ''
  if (!id) return null

  // Stylometry profile nodes: "stylo_" or "stylometry_" or "s1..s9"
  if (id.startsWith('stylo_') || id.startsWith('stylometry_') || /^s\d+$/i.test(id)) return 'stylometry'

  // Actors are bare UUIDs (or "actor_" for legacy mock/fallback ids).
  if (UUID_RE.test(id) || id.startsWith('actor_') || /^a\d+$/i.test(id)) return 'actor'

  // Infra indicators: "infra_" + indicatorType + "_" + onionHost.
  if (id.startsWith('infra_')) return 'infra'

  // Identifiers: "id_" + kind + "_" + value.
  if (id.startsWith('id_')) {
    const kind = id.slice(3).split(/[:,_]/)[0].toLowerCase()
    return IDENTIFIER_KIND_TO_TYPE[kind] || null
  }

  // Legacy prefix-only ids.
  if (id.startsWith('site_') || id.startsWith('onion_')) return 'onion'
  if (id.startsWith('ip_') || id.startsWith('host_')) return 'infra'
  if (id.startsWith('email_')) return 'email'
  if (id.startsWith('pgp_')) return 'pgp'
  if (id.startsWith('btc_') || id.startsWith('xmr_') || id.startsWith('ltc_') || id.startsWith('eth_')) return 'wallet'

  return null
}

export default function App() {
  const SIDEBAR_PREF_KEY = 'onioneye.sidebarCollapsed.v1'
  const [sidebarCollapsed, setSidebarCollapsed] = useState(() => {
    try {
      const raw = localStorage.getItem(SIDEBAR_PREF_KEY)
      if (raw === 'true') return true
      if (raw === 'false') return false
    } catch {}
    return false
  })

  useEffect(() => {
    try {
      localStorage.setItem(SIDEBAR_PREF_KEY, String(sidebarCollapsed))
    } catch {}
  }, [sidebarCollapsed])

  const graphRef = useRef()
  const graphContainerRef = useRef()

  // Graph-layout bookkeeping (used by the “Force-directed layout tuning” below):
  //  - nodePosCacheRef:   frozen x/y per node id. Every 10s poll brings back new
  //    node objects with no position; d3 would drop them all onto a tiny default
  //    circle (the classic “blob”). Restoring cached positions keeps the layout
  //    perfectly still between polls.
  //  - clustersRef:       connected components of the current graph, recomputed
  //    in applyGraphLayout on every data change.
  //  - graphDataKeyRef:   stable fingerprint (sorted node ids). We only re-heat
  //    the simulation when it actually changes — identical polls stay frozen.
  //  - autoFitPendingRef: set on first layout + RE-SIMULATE; onEngineStop fits
  //    the constellation when true, then clears (never hijacks the camera later).
  const nodePosCacheRef = useRef(new Map())
  const clustersRef = useRef([])
  const graphDataKeyRef = useRef('')
  const autoFitPendingRef = useRef(false)
  // Text-measure cache: keyed by rendered label string. Invalidated
  // automatically whenever the label or its font size changes; no
  // time-based purge so layout stays stable across the 10s poll.
  const textMeasurementCacheRef = useRef(new Map())
  // layoutModeRef mirrors the layoutMode state so callbacks like
  // seedNodePositions can read the current mode without being
  // stale-closure-captured by applyGraphLayout's dependency array.
  const layoutModeRef = useRef('organic')

  const [selectedNode, setSelectedNode] = useState(null)
  const [hoveredNode, setHoveredNode] = useState(null)
  const [hoveredLink, setHoveredLink] = useState(null)
  const [tab, setTab] = useState('graph')
  const [layoutMode, setLayoutMode] = useState('organic') // 'organic' | 'hierarchical'
  useEffect(() => { layoutModeRef.current = layoutMode }, [layoutMode])
  const [isScanModalOpen, setScanModalOpen] = useState(false)
  const [isLiveApi, setIsLiveApi] = useState(false)
  const [aiMode, setAiMode] = useState(false)
  const [aiRestarting, setAiRestarting] = useState(false)
  const [torStatus, setTorStatus] = useState(null)
  const [torLive, setTorLive] = useState(false)

  // Shared search query (so RightRail can search/clear it)
  const [searchQuery, setSearchQuery] = useState('')

  // Filters
  const [filterType, setFilterType] = useState('all')
  const [confidenceThreshold, setConfidenceThreshold] = useState(0)

  // Data
  const [stats, setStats] = useState(null)
  const [leads, setLeads] = useState([])
  const [jobs, setJobs] = useState([]) // optimistic local crawl-job list (dispatch feedback)
  const [rawGraphData, setRawGraphData] = useState({ nodes: [], links: [] })
  // Overlay of nodes/links imported via OPEN RECORD — survives the 10s graph poll.
  const [importedOverlay, setImportedOverlay] = useState({ nodes: [], links: [] })

  // Extracted evidence gallery (image_hashes → actors). Drives the MEDIA tab.
  const [media, setMedia] = useState([])
  // Full-screen overlay for orphan media rows that have no owning actor node.
  const [mediaViewer, setMediaViewer] = useState(null)

  const [dim, setDim] = useState({
    w: typeof window !== 'undefined' ? window.innerWidth - 300 : 800,
    h: typeof window !== 'undefined' ? window.innerHeight - 120 : 600,
  })

  const loadData = useCallback(async () => {
    const [s, g, l, m] = await Promise.all([
      fetchStats(), fetchGraph(), fetchLeads(), fetchMedia()
    ])
    if (s.live || g.live || l.live || m.live) setIsLiveApi(true)
    if (s.data) setStats(s.data)
    if (g.data?.nodes) {
      setRawGraphData(prev => mergeGraphSnapshot(prev, g.data))
    }
    if (Array.isArray(l.data)) setLeads(l.data)
    if (m.live && Array.isArray(m.data)) setMedia(m.data)
  }, [])

  useEffect(() => {
    loadData()
    const interval = setInterval(loadData, 10000)
    return () => clearInterval(interval)
  }, [loadData])

  // Poll AI-mode status so the AI log feeds + restart controls stay in sync
  // with whether the Ollama engine is actually up.
  const checkAiStatus = useCallback(async () => {
    const res = await fetchAiStatus()
    setAiMode(!!(res.live && res.data?.ai_enabled))
  }, [])

  useEffect(() => {
    checkAiStatus()
    const interval = setInterval(checkAiStatus, 10000)
    return () => clearInterval(interval)
  }, [checkAiStatus])

  const checkTorStatus = useCallback(async () => {
    const res = await fetchTorStatus()
    setTorStatus(res.data || null)
    setTorLive(!!(res.live && res.data))
  }, [])

  useEffect(() => {
    checkTorStatus()
    const interval = setInterval(checkTorStatus, 5000)
    return () => clearInterval(interval)
  }, [checkTorStatus])

  // Restart the whole stack with or without the AI profile. The backend takes
  // ~15s to bring everything back up, so we flip the UI into a "restarting"
  // state and let the user reload when it's done.
  const handleRestartStack = useCallback(async (withAi) => {
    if (aiRestarting) return
    setAiRestarting(true)
    try {
      await restartStackWithAi(withAi)
    } catch (e) {
      // non-fatal; polling will reconcile actual state
    }
    setTimeout(() => { setAiRestarting(false); checkAiStatus() }, 15000)
  }, [aiRestarting, checkAiStatus])

  // Poll the dashboard notification ring for crawler-worker events (rotations +
// anomalies) and show a transient 3s popup for each NEW event. Never repeat an
// event, even after the in-memory ring on the backend is overwritten: each item
// is keyed by its unique id (rotation_id for rotations, anomaly_id otherwise).
  const [notice, setNotice] = useState(null)
  const seenNoticeKeys = useRef(new Set())

  // notificationKey returns a stable dedup key for one ring item. Anomaly events
  // always carry anomaly_id; rotations carry rotation_id. Without a recognized
  // id the item is unpollable (skipped).
  function notificationKey(item) {
    if (!item?.event) return null
    const id = item.anomaly_id || item.rotation_id
    return id ? `${item.event}:${id}` : null
  }

  useEffect(() => {
    let cancelled = false
    let seeded = false
    async function checkNotifications() {
      const res = await fetchNotifications()
      if (cancelled || !res.live || !res.data?.items) return
      // First poll: seed the seen set so we don't replay historical events as a
      // popup on page load. Subsequent polls show only NEW events.
      if (!seeded) {
        for (const item of res.data.items) {
          const key = notificationKey(item)
          if (key) seenNoticeKeys.current.add(key)
        }
        seeded = true
        return
      }
      for (const item of res.data.items) {
        const key = notificationKey(item)
        if (!key) continue
        if (seenNoticeKeys.current.has(key)) continue
        seenNoticeKeys.current.add(key)
        setNotice(item)
        break // one popup at a time
      }
    }
    checkNotifications()
    const interval = setInterval(checkNotifications, 3000)
    return () => {
      cancelled = true
      clearInterval(interval)
    }
  }, [])

  useEffect(() => {
    function updateSize() {
      if (graphContainerRef.current) {
        const rect = graphContainerRef.current.getBoundingClientRect()
        if (rect.width > 50 && rect.height > 50) {
          setDim({ w: Math.floor(rect.width), h: Math.floor(rect.height) })
        }
      }
    }
    updateSize()
    const observer = new ResizeObserver(updateSize)
    if (graphContainerRef.current) observer.observe(graphContainerRef.current)
    window.addEventListener('resize', updateSize)
    return () => {
      observer.disconnect()
      window.removeEventListener('resize', updateSize)
    }
  }, [tab, selectedNode])

  // Filtered nodes and links based on UI state
  const graphData = useMemo(() => {
    const overlayNodes = importedOverlay?.nodes || []
    const overlayLinks = importedOverlay?.links || []
    const liveNodes = rawGraphData?.nodes || []
    const liveLinks = rawGraphData?.links || []
    const liveIds = new Set(liveNodes.map(n => n.id))
    const extraNodes = overlayNodes.filter(n => n?.id && !liveIds.has(n.id))
    const combinedNodes = [...liveNodes, ...extraNodes]
    if (!combinedNodes.length) return { nodes: [], links: [] }

    let validNodes = combinedNodes.filter(n => !isSeedKeywordNode(n))
    if (filterType === 'flagged') {
      validNodes = validNodes.filter(n => n.flagged)
    } else if (filterType !== 'all') {
      validNodes = validNodes.filter(n => n.type === filterType)
    }
    const nodeSet = new Set(validNodes.map(n => n.id))

    const liveLinkKeys = new Set()
    const validLinks = (liveLinks || []).filter(l => {
      const s = typeof l.source === 'object' ? l.source.id : l.source
      const t = typeof l.target === 'object' ? l.target.id : l.target
      liveLinkKeys.add(`${s}|${t}`)
      const hasNodes = nodeSet.has(s) && nodeSet.has(t)
      const meetsConfidence = (l.confidence ?? 1.0) >= confidenceThreshold
      return hasNodes && meetsConfidence
    })
    const extraLinks = overlayLinks.filter(l => {
      const s = typeof l.source === 'object' ? l.source.id : l.source
      const t = typeof l.target === 'object' ? l.target.id : l.target
      if (!nodeSet.has(s) || !nodeSet.has(t)) return false
      if (liveLinkKeys.has(`${s}|${t}`)) return false
      return (l.confidence ?? 1.0) >= confidenceThreshold
    })

    // Keep output stable: if topology doesn't change, reusing the same
    // underlying node/link objects (mergeGraphSnapshot) will avoid reheat.
    return { nodes: validNodes, links: [...validLinks, ...extraLinks] }
  }, [rawGraphData, importedOverlay, filterType, confidenceThreshold])

  // ── Memoized topology key to avoid downstream re-layout on identical polls ──
  const graphTopologyKeyFiltered = useMemo(() => {
    return graphTopologyKey(graphData)
  }, [graphData])
  // ── Force-directed layout tuning ──────────────────────────────────────────
  // Default d3 setup (charge −30, link 30, center 1.0, no collision, every
  // node born at the origin) collapses large graphs into one dense mid-screen
  // blob. We seed nodes on a golden-angle spiral and re-tune the forces so the
  // graph relaxes into a readable, spread-out constellation instead.
  const seedNodePositions = useCallback((nodes, { reshuffleAll = false } = {}) => {
    if (!nodes.length) return
    const cache = nodePosCacheRef.current
    const clusters = clustersRef.current
    const clusterCount = clusters.length || 1
    const goldenAngle = Math.PI * (3 - Math.sqrt(5))

    // Outer spiral radius grows with the number of clusters, but is bounded
    // so multi-case graphs remain in view rather than dispersing into deep space.
    const outerR = Math.min(620, Math.max(280, Math.sqrt(clusterCount) * 170))

    // On a full re-layout, scramble the cluster order so the picture changes.
    const clusterOrder = reshuffleAll
      ? clusters.map((_, i) => i).sort(() => Math.random() - 0.5)
      : clusters.map((_, i) => i)

    clusters.forEach((c, i) => {
      const a = (clusterOrder[i] + 0.5) * goldenAngle
      const cx = outerR * Math.cos(a)
      const cy = outerR * Math.sin(a)
      const count = c.members.length
      // Each cluster is itself a compact spiral so connected nodes start out tight.
      const localR = Math.min(210, Math.max(70, Math.sqrt(count) * 28))

      c.members.forEach((node, j) => {
        // Reuse live or cached positions unless reshuffling — this is what keeps
        // the graph perfectly still across the 10s polls.
        if (!reshuffleAll) {
          if (typeof node.x === 'number' && typeof node.y === 'number' && isFinite(node.x) && isFinite(node.y)) {
            cache.set(node.id, { x: node.x, y: node.y })
            return
          }
          const cached = cache.get(node.id)
          if (cached && typeof cached.x === 'number' && typeof cached.y === 'number' && isFinite(cached.x) && isFinite(cached.y)) {
            node.x = cached.x
            node.y = cached.y
            node.vx = node.vy = 0
            return
          }
        }

        // Stylokey-aware seed for actor nodes: use a stable 2D reduction of
        // the backend key so similar actors start closer.
        const sk2 = node.type === 'actor' ? stylokeyTo2D(node.stylokey) : null
        if (sk2 && typeof sk2.x === 'number' && typeof sk2.y === 'number' && isFinite(sk2.x) && isFinite(sk2.y)) {
          // Keep scaling bounded so large key magnitudes don't explode layout.
          const mag = Math.hypot(sk2.x, sk2.y) || 1
          const nx = sk2.x / mag
          const ny = sk2.y / mag

          // Place actors toward the component center, but still allow the
          // cluster spiral to separate multiple actors.
          const jr = localR * Math.sqrt((j + 0.5) / count)
          node.x = cx + jr * nx
          node.y = cy + jr * ny
        } else {
          const ja = j * goldenAngle
          const jr = localR * Math.sqrt((j + 0.5) / count)
          node.x = cx + jr * Math.cos(ja)
          node.y = cy + jr * Math.sin(ja)
        }

        node.vx = 0
        node.vy = 0
        cache.set(node.id, { x: node.x, y: node.y })
      })
    })
  }, [])

  // ── Hierarchical tree layout ────────────────────────────────────────────
  // Builds a rigid top-to-bottom tree from the current graph links:
  //  - roots = nodes with zero incoming degree  (e.g. actor / user nodes)
  //  - then BFS layers by link direction
  //  - each layer gets a fixed Y spacing; nodes within a layer are spread
  //    evenly across the X axis so the structure looks like a rigid body.
  // Only nodes that are actually connected by the graph's links are
  // positioned this way — orphans fall through to a fallback cluster spiral.
  const seedHierarchicalPositions = useCallback((nodes, links, { reshuffleAll = false } = {}) => {
    if (!nodes.length) return
    const nodeMap = new Map(nodes.map(n => [n.id, n]))
    const inDegree = new Map()
    const outMap = new Map()
    for (const n of nodes) { inDegree.set(n.id, 0); outMap.set(n.id, []) }
    for (const l of links || []) {
      const s = typeof l.source === 'object' ? l.source.id : l.source
      const t = typeof l.target === 'object' ? l.target.id : l.target
      if (nodeMap.has(s) && nodeMap.has(t)) {
        inDegree.set(t, (inDegree.get(t) || 0) + 1)
        outMap.get(s).push(t)
      }
    }
    // roots = zero in-degree; if everything has in-degree, pick the
    // node(s) with minimum in-degree so the tree still has a top.
    let roots = nodes.filter(n => (inDegree.get(n.id) || 0) === 0)
    if (!roots.length) {
      const minIn = Math.min(...Array.from(inDegree.values()))
      roots = nodes.filter(n => (inDegree.get(n.id) || 0) === minIn)
    }

    // BFS layering (preserve insertion order of roots)
    const layerOf = new Map()
    const queue = []
    for (const r of roots) { layerOf.set(r.id, 0); queue.push(r.id) }
    let head = 0
    while (head < queue.length) {
      const cur = queue[head++]
      const lay = layerOf.get(cur) + 1
      for (const nxt of outMap.get(cur) || []) {
        if (!layerOf.has(nxt)) { layerOf.set(nxt, lay); queue.push(nxt) }
      }
    }
    // Any node not reached by BFS gets its own layer after the deepest
    const maxLayer = queue.length ? Math.max(...Array.from(layerOf.values())) : 0
    nodes.forEach(n => { if (!layerOf.has(n.id)) layerOf.set(n.id, maxLayer + 1) })

    // Group nodes by layer
    const layers = new Map()
    for (const n of nodes) {
      const ly = layerOf.get(n.id) || 0
      if (!layers.has(ly)) layers.set(ly, [])
      layers.get(ly).push(n)
    }

    const totalLayers = layers.size || 1
    const ySpacing = Math.min(160, Math.max(80, 620 / totalLayers))
    const H_PAD = 50
    const V_PAD = 40

    layers.forEach((members, ly) => {
      const y = V_PAD + ly * ySpacing
      const layerWidth = Math.min(900, Math.max(300, nodes.length * 42))
      const step = layerWidth / (members.length + 1)
      members.forEach((node, j) => {
        const cache = nodePosCacheRef.current
        // Reuse live or cached positions unless reshuffling — this is what keeps
        // the tree still across the 10s polls instead of collapsing to a blob.
        if (!reshuffleAll) {
          if (typeof node.x === 'number' && typeof node.y === 'number' && isFinite(node.x) && isFinite(node.y)) {
            cache.set(node.id, { x: node.x, y: node.y })
            return
          }
          const cached = cache.get(node.id)
          if (cached && typeof cached.x === 'number' && typeof cached.y === 'number' && isFinite(cached.x) && isFinite(cached.y)) {
            node.x = cached.x
            node.y = cached.y
            node.vx = 0
            node.vy = 0
            return
          }
        }

        const x = H_PAD + step * (j + 1)
        node.x = x
        node.y = y
        node.vx = 0
        node.vy = 0
        cache.set(node.id, { x, y })
      })
    })
  }, [])
  const applyGraphLayout = useCallback(({ reshuffleAll = false } = {}) => {
    const nodes = graphData.nodes || []
    if (!nodes.length) {
      graphDataKeyRef.current = ''
      return
    }

    // Recompute connected components for this data snapshot, then seed any node
    // that has no live/cached position (new arrivals get placed in their cluster).
    clustersRef.current = computeClusters(graphData.nodes, graphData.links)
    // Hierarchical mode: rigid tree — root→users/urls→img/images/bitcoin/pgp.
    if (layoutModeRef.current === 'hierarchical') {
      seedHierarchicalPositions(nodes, graphData.links || [], { reshuffleAll })
      // Pin the tree so charge/link/collision cannot collapse it between polls.
      nodes.forEach(n => {
        if (typeof n.x === 'number' && isFinite(n.x)) n.fx = n.x
        if (typeof n.y === 'number' && isFinite(n.y)) n.fy = n.y
      })
    } else {
      seedNodePositions(nodes, { reshuffleAll })
      if (reshuffleAll) {
        nodes.forEach(n => {
          n.fx = undefined
          n.fy = undefined
        })
      }
    }

    const fg = graphRef.current
    if (!fg) return
    const n = nodes.length

    // Stable repulsion tuned to keep graph bounded and readable as nodes scale
    const chargeForce = fg.d3Force('charge')
    if (chargeForce) {
      chargeForce.strength(-(140 + Math.min(220, n * 1.6)))
      if (typeof chargeForce.distanceMax === 'function') {
        chargeForce.distanceMax(900)
      }
    }

    // Balanced center pull → clusters stay within the radar view instead of drifting off-canvas
    const centerForce = fg.d3Force('center')
    if (centerForce) {
      centerForce.strength(0.05)
    }

    // Type-aware link lengths — keeps clusters well-proportioned without hairball tangles
    const linkForce = fg.d3Force('link')
    if (linkForce) {
      linkForce
        .distance(l => {
          const s = typeof l.source === 'object' ? l.source : null
          const t = typeof l.target === 'object' ? l.target : null
          const st = s?.type || 'unknown'
          const tt = t?.type || 'unknown'
          if (st === 'actor' && tt === 'actor') return 210
          if (st === 'actor' || tt === 'actor') return 150
          if (st === 'stylometry' || tt === 'stylometry') return 110
          if (st === 'onion' || tt === 'onion') return 95
          return 80
        })
        .strength(l => (l.type === 'persona' ? 0.18 : 0.55))
        .iterations(2)
    }

    // Collision force → rounded-rect users / triangle stylo vectors never overlap
    fg.d3Force(
      'collision',
      forceCollide()
        .radius(d => {
          const fp = nodeFootprint(d)
          if (fp.kind === 'rect') return Math.max(fp.w, fp.h) / 2 + 22
          if (fp.kind === 'triangle') return fp.r + 18
          return fp.r + 16
        })
        .strength(1)
        .iterations(3),
    )

    // Cluster cohesion force → each connected component stays an island instead
    // of being shredded apart by the charge force: members are pulled gently
    // toward their own component's centroid.
    fg.d3Force('cluster', (alpha) => {
      const k = alpha * 0.08
      const clusters = clustersRef.current || []
      for (const c of clusters) {
        const m = c.members || []
        if (m.length < 2) continue
        let cx = 0
        let cy = 0
        let validCount = 0
        for (const node of m) {
          if (typeof node.x === 'number' && typeof node.y === 'number' && isFinite(node.x) && isFinite(node.y)) {
            cx += node.x
            cy += node.y
            validCount++
          }
        }
        if (validCount < 2) continue
        cx /= validCount
        cy /= validCount
        for (const node of m) {
          if (typeof node.x === 'number' && typeof node.y === 'number' && isFinite(node.x) && isFinite(node.y)) {
            if (typeof node.vx === 'number' && isFinite(node.vx)) {
              node.vx -= (node.x - cx) * k
            }
            if (typeof node.vy === 'number' && isFinite(node.vy)) {
              node.vy -= (node.y - cy) * k
            }
          }
        }
      }
    })

    // Only re-heat on real topology changes (added/removed nodes or links)
    // or an explicit reshuffle. Identical 10s polls stay frozen.
    const changed = graphTopologyKeyFiltered !== graphDataKeyRef.current
    const key = graphTopologyKeyFiltered
    const firstLayout = changed && !graphDataKeyRef.current
    graphDataKeyRef.current = key

    if (reshuffleAll || changed) {
      if (reshuffleAll || firstLayout) autoFitPendingRef.current = true
      if (typeof fg.d3AlphaTarget === 'function') fg.d3AlphaTarget(0)
      fg.d3ReheatSimulation()
    } else if (typeof fg.d3AlphaTarget === 'function') {
      fg.d3AlphaTarget(0)
    }
  }, [graphData, seedNodePositions, seedHierarchicalPositions])

  // Apply tuned forces whenever the filtered graph data changes
  useEffect(() => {
    applyGraphLayout()
  }, [applyGraphLayout])

  // When layout mode switches, clear position cache and re-layout
  useEffect(() => {
    nodePosCacheRef.current.clear()
    applyGraphLayout()
  }, [layoutMode])

  // Freeze / unfreeze simulation physics when focused on a node so the graph
  // is completely stable and stationary while inspecting an entity dossier.
  useEffect(() => {
    if (!graphData?.nodes) return
    if (selectedNode) {
      // Pin all nodes in place to freeze motion
      graphData.nodes.forEach(n => {
        if (typeof n.x === 'number' && isFinite(n.x)) n.fx = n.x
        if (typeof n.y === 'number' && isFinite(n.y)) n.fy = n.y
      })
      if (graphRef.current && typeof graphRef.current.d3AlphaTarget === 'function') {
        graphRef.current.d3AlphaTarget(0)
      }
    } else {
      // Unpin nodes when focus is cleared
      graphData.nodes.forEach(n => {
        n.fx = undefined
        n.fy = undefined
      })
    }
  }, [selectedNode, graphData])

  // RE-SIMULATE button → reshuffle every cluster + member, re-heat, auto-fit
  const handleReLayout = useCallback(() => {
    if (graphData?.nodes) {
      graphData.nodes.forEach(n => {
        n.fx = undefined
        n.fy = undefined
      })
    }
    applyGraphLayout({ reshuffleAll: true })
  }, [applyGraphLayout, graphData])

  // Pre-seed positions during render — NOT in an effect. react-force-graph-2d
  // pushes graphData into the sim synchronously while rendering, and d3 would
  // stamp fresh poll objects (which arrive with no x/y) onto a tiny default
  // circle — the classic “blob”. Seeding right here, before the child renders,
  // guarantees force-graph's ingest sees valid cluster/cached coordinates.
  const seededGraphData = useMemo(() => {
    const nodes = graphData.nodes || []
    if (!nodes.length) return graphData
    clustersRef.current = computeClusters(nodes, graphData.links || [])
    if (layoutMode === 'hierarchical') {
      seedHierarchicalPositions(nodes, graphData.links || [], { reshuffleAll: false })
      // Do not set fx/fy here; it is controlled by the selection-focused
      // “Freeze / unfreeze” effect to avoid fighting user interactions.
    } else {
      seedNodePositions(nodes)
    }
    return graphData
  }, [graphData, layoutMode, seedNodePositions, seedHierarchicalPositions])

  const neighbourSet = useMemo(() => {
    const m = new Map()
    if (!graphData?.nodes) return m
    graphData.nodes.forEach(n => m.set(n.id, new Set()))
    graphData.links?.forEach(l => {
      const s = typeof l.source === 'object' ? l.source.id : l.source
      const t = typeof l.target === 'object' ? l.target.id : l.target
      m.get(s)?.add(t)
      m.get(t)?.add(s)
    })
    return m
  }, [graphData])

  const isHighlighted = useCallback((nodeId) => {
    if (!hoveredNode && !selectedNode) return true
    const focus = hoveredNode || selectedNode
    if (nodeId === focus) return true
    return neighbourSet.get(focus)?.has(nodeId) ?? false
  }, [hoveredNode, selectedNode, neighbourSet])

  // Fit camera view cleanly to encompass the bounded graph constellation
  const fitGraphView = useCallback((duration = 500, padding = 44) => {
    if (!graphRef.current) return
    graphRef.current.zoomToFit(duration, padding)
  }, [])

  // Center camera on node when selected with smooth motion
  const handleSelectNode = useCallback((nodeId) => {
    setSelectedNode(nodeId)
    if (!nodeId || !graphRef.current) return
    const node = graphData.nodes?.find(n => n.id === nodeId)
    if (node && typeof node.x === 'number' && typeof node.y === 'number' && isFinite(node.x) && isFinite(node.y)) {
      const fg = graphRef.current
      const currentZoom = typeof fg.zoom === 'function' ? fg.zoom() : 1
      const targetZoom = Math.min(2.05, Math.max(currentZoom, 1.55))
      fg.centerAt(node.x, node.y, 1100)
      fg.zoom(targetZoom, 1100)
    }
  }, [graphData])

  // Clicking a MEDIA-tab thumbnail behaves exactly like clicking its graph
  // node: it opens the entity dossier (DetailPanel) on the RADAR tab. The type
  // filter is reset so the owning actor node is guaranteed to be in view even
  // when a node-type filter was hiding it. Orphan rows without an actor just
  // open in the full-screen evidence viewer.
  const handleSelectMediaImage = useCallback((img) => {
    if (!img) return
    if (img.actor_id) {
      setFilterType('all')
      setTab('graph')
      handleSelectNode(img.actor_id)
      return
    }
    setMediaViewer(img.src || null)
  }, [handleSelectNode])

  // Merge an imported/saved record into the live graph (dedup by node id).
  // Used by the right-rail OPEN RECORD panel. Stored as a separate overlay so
  // the 10s graph poll cannot wipe imported case data.
  const handleMergeGraph = useCallback((imported) => {
    if (!imported?.nodes?.length) return
    setImportedOverlay(prev => {
      const existingIds = new Set((prev.nodes || []).map(n => n.id))
      const newNodes = imported.nodes.filter(n => n?.id && !existingIds.has(n.id))
      const mergedNodes = [...(prev.nodes || []), ...newNodes]
      const nodeSet = new Set(mergedNodes.map(n => n.id))
      const existingLinks = new Set(
        (prev.links || []).map(l => {
          const s = typeof l.source === 'object' ? l.source.id : l.source
          const t = typeof l.target === 'object' ? l.target.id : l.target
          return `${s}|${t}`
        })
      )
      const newLinks = (imported.links || []).filter(l => {
        const s = typeof l.source === 'object' ? l.source.id : l.source
        const t = typeof l.target === 'object' ? l.target.id : l.target
        if (!s || !t) return false
        if (!nodeSet.has(s) || !nodeSet.has(t)) return false
        const key = `${s}|${t}`
        if (existingLinks.has(key)) return false
        existingLinks.add(key)
        return true
      }).map(l => ({
        source: typeof l.source === 'object' ? l.source.id : l.source,
        target: typeof l.target === 'object' ? l.target.id : l.target,
        type: l.type || 'imported',
        basis: l.basis,
        confidence: l.confidence,
      }))
      return { nodes: mergedNodes, links: [...(prev.links || []), ...newLinks] }
    })
  }, [])


  const paintNode = useCallback((node, ctx, globalScale) => {
    if (!node || typeof node.x !== 'number' || typeof node.y !== 'number' || !isFinite(node.x) || !isFinite(node.y)) return
    const isFocus = node.id === selectedNode || node.id === hoveredNode
    const highlighted = isHighlighted(node.id)
    const isStyloInactive = node.type === 'stylometry' && isStyloUnderThreshold(node)
    const baseColor = isStyloInactive ? DULL_STYLO_COLOR : (NODE_COLORS[node.type] || '#00ff9d')
    const r = node.type === 'actor' ? 14 : node.type === 'stylometry' ? 11 : 9

    ctx.save()
    ctx.globalAlpha = highlighted ? 1.0 : (isStyloInactive ? 0.08 : 0.12)

    // Outer glow aura (type-aware shape)
    if (isFocus || node.flagged || (node.type === 'stylometry' && !isStyloInactive)) {
      const auraColor = node.flagged ? 'rgba(255, 23, 68, 0.4)' : baseColor + (isStyloInactive ? '22' : '44')
      ctx.beginPath()
      // Draw shape slightly larger than the body
      const pad = isStyloInactive ? 4 : 10
      const fp = nodeFootprint(node)
      if (fp.kind === 'rect') {
        drawRoundedRectPath(ctx, node.x, node.y, fp.w + pad * 2, fp.h + pad * 2, fp.corner + pad * 0.4)
      } else if (fp.kind === 'triangle') {
        drawTrianglePath(ctx, node.x, node.y, fp.r + pad)
      } else {
        ctx.arc(node.x, node.y, fp.r + pad, 0, 2 * Math.PI)
      }
      ctx.fillStyle = auraColor
      ctx.fill()
    }

    // Reticle brackets on focused node
    if (isFocus) {
      ctx.strokeStyle = baseColor
      ctx.lineWidth = 1.5
      const bracket = r + 8
      const corner = 5

      // Top-Left
      ctx.beginPath()
      ctx.moveTo(node.x - bracket, node.y - bracket + corner)
      ctx.lineTo(node.x - bracket, node.y - bracket)
      ctx.lineTo(node.x - bracket + corner, node.y - bracket)
      ctx.stroke()

      // Top-Right
      ctx.beginPath()
      ctx.moveTo(node.x + bracket - corner, node.y - bracket)
      ctx.lineTo(node.x + bracket, node.y - bracket)
      ctx.lineTo(node.x + bracket, node.y - bracket + corner)
      ctx.stroke()

      // Bottom-Left
      ctx.beginPath()
      ctx.moveTo(node.x - bracket, node.y + bracket - corner)
      ctx.lineTo(node.x - bracket, node.y + bracket)
      ctx.lineTo(node.x - bracket + corner, node.y + bracket)
      ctx.stroke()

      // Bottom-Right
      ctx.beginPath()
      ctx.moveTo(node.x + bracket - corner, node.y + bracket)
      ctx.lineTo(node.x + bracket, node.y + bracket)
      ctx.lineTo(node.x + bracket, node.y + bracket - corner)
      ctx.stroke()
    }

       // Solid Node body (type-aware shape)
    const fp = nodeFootprint(node)
    const bodyRadiusPad = isStyloInactive ? 0 : 0
    const body = drawNodeShapePath(ctx, { ...node, x: node.x, y: node.y, type: node.type }, bodyRadiusPad)

    ctx.fillStyle = isStyloInactive ? '#0a0f16' : '#060b12'
    ctx.fill()
    ctx.strokeStyle = node.flagged ? '#ff1744' : (isStyloInactive ? '#475569' : baseColor)
    ctx.lineWidth = node.flagged ? 2.5 : (isStyloInactive ? 1.3 : 1.8)
    if (isStyloInactive) {
      ctx.setLineDash([3, 2])
    } else {
      ctx.setLineDash([])
    }
    ctx.stroke()
    ctx.setLineDash([])

    // Inner icon badge
    const icon = TYPE_ICONS[node.type] || '●'
    const fpR = fp.kind === 'rect' ? Math.max(fp.w, fp.h) / 2 : (fp.kind === 'triangle' ? fp.r : fp.r)
    const iconSize = Math.max(9, Math.floor(fpR * 1.05))
    ctx.font = `${iconSize}px ${CANVAS_FONT_STACK}`
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillStyle = isStyloInactive ? '#64748b' : '#ffffff'
    ctx.fillText(icon, node.x, node.y + 1)

    // Flagged pulsing red hazard mark (type-aware offset)
    if (node.flagged) {
      const fp = nodeFootprint(node)
      const dx = fp.kind === 'rect' ? fp.w * 0.32 : (fp.kind === 'triangle' ? fp.r * 0.7 : fp.r * 0.75)
      const dy = fp.kind === 'rect' ? fp.h * 0.32 : (fp.kind === 'triangle' ? fp.r * 0.7 : fp.r * 0.75)
      ctx.beginPath()
      ctx.arc(node.x + dx, node.y - dy, 4, 0, 2 * Math.PI)
      ctx.fillStyle = '#ff1744'
      ctx.fill()
      ctx.strokeStyle = '#020406'
      ctx.lineWidth = 1
      ctx.stroke()
    }

    // Typography Label
    if (globalScale > 0.35 || highlighted || isFocus) {
      const fontSize = Math.max(9, Math.min(12, 11 / Math.max(globalScale, 0.65)))
      const fontKey = `${fontSize}px ${CANVAS_FONT_STACK}`
      ctx.font = fontKey
      ctx.textAlign = 'center'
      ctx.textBaseline = 'top'

      // Label background plate — short label by default, FULL name on hover/focus
      const fullLabel = String(node.label || node.value || node.id || '') + (isStyloInactive ? ' [INACTIVE <150w]' : '')
      const labelText = isFocus ? fullLabel : truncateLabel(fullLabel)

      // Use cached text width if available; invalidate whenever the
      // rendered label or the font stack changes instead of using a
      // fixed time window so the 10s poll keeps the layout still.
      const cacheKey = `${labelText}::${fontKey}`
      let textWidth = textMeasurementCacheRef.current.get(cacheKey)
      if (textWidth === undefined) {
        textWidth = ctx.measureText(labelText).width
        textMeasurementCacheRef.current.set(cacheKey, textWidth)
      }
      ctx.fillStyle = 'rgba(4, 7, 12, 0.85)'
      ctx.fillRect(node.x - textWidth / 2 - 4, node.y + r + 4, textWidth + 8, fontSize + 4)
      ctx.strokeStyle = highlighted ? (isStyloInactive ? '#475569' : baseColor + '66') : 'rgba(21, 34, 51, 0.8)'
      ctx.lineWidth = 1
      ctx.strokeRect(node.x - textWidth / 2 - 4, node.y + r + 4, textWidth + 8, fontSize + 4)

      // Text fill
      ctx.fillStyle = highlighted ? (isStyloInactive ? '#94a3b8' : '#ffffff') : (isStyloInactive ? '#475569' : '#8b9eb3')
      ctx.fillText(labelText, node.x, node.y + r + 6)
    }

    ctx.restore()
  }, [selectedNode, hoveredNode, isHighlighted])

  const paintLink = useCallback((link, ctx) => {
    const s = link.source, t = link.target
    if (!s || !t) return
    if (typeof s !== 'object' || typeof t !== 'object') return
    if (typeof s.x !== 'number' || typeof s.y !== 'number' || typeof t.x !== 'number' || typeof t.y !== 'number') return
    if (!isFinite(s.x) || !isFinite(s.y) || !isFinite(t.x) || !isFinite(t.y)) return

    const sId = s.id || s
    const tId = t.id || t
    const active = isHighlighted(sId) && isHighlighted(tId)

    const isStyloLink = link.type === 'stylometry'
    const isSourceStyloInactive = s.type === 'stylometry' && isStyloUnderThreshold(s)
    const isTargetStyloInactive = t.type === 'stylometry' && isStyloUnderThreshold(t)
    const isUnderThresholdLink = isStyloLink && (isSourceStyloInactive || isTargetStyloInactive)

    ctx.save()
    ctx.globalAlpha = active ? 0.9 : (isUnderThresholdLink ? 0.04 : 0.08)
    ctx.beginPath()
    ctx.moveTo(s.x, s.y)
    ctx.lineTo(t.x, t.y)
    ctx.strokeStyle = isUnderThresholdLink ? '#3a4659' : (LINK_COLORS[link.type] || '#1e3550')
    ctx.lineWidth = active ? (link.type === 'persona' ? 2.5 : 1.5) : (isUnderThresholdLink ? 0.6 : 0.8)

    if (isUnderThresholdLink) {
      ctx.setLineDash([3, 3])
    } else if (link.type === 'infra') {
      ctx.setLineDash([4, 4])
    } else if (link.type === 'identifier') {
      ctx.setLineDash([2, 2])
    } else {
      ctx.setLineDash([])
    }

    ctx.stroke()
    ctx.restore()
  }, [isHighlighted])

  const detail = useMemo(() => {
    if (!selectedNode || !graphData?.nodes) return null
    const node = graphData.nodes.find(n => n.id === selectedNode)
    if (!node) return null
    const edges = (graphData.links || []).filter(l => {
      const s = typeof l.source === 'object' ? l.source.id : l.source
      const t = typeof l.target === 'object' ? l.target.id : l.target
      return s === node.id || t === node.id
    })
    return { node, edges }
  }, [selectedNode, graphData])


  const handleDeleteNode = useCallback(async (node) => {
    if (!node) return

    // Some graph nodes may arrive without `type` populated.
    // Storage-api requires both node_id and node_type.
    const derivedType = deriveNodeType(node)

    if (!derivedType) {
      alert(`Cannot delete node: missing node type for id=${node.id}`)
      return
    }

    if (!window.confirm(`Delete node ${node.label || node.id}? (Also deletes connected persona links)`)) return
    try {
      await deleteNode({ node_id: node.id, node_type: derivedType })
      setSelectedNode(null)
      await loadData()
    } catch (e) {
      alert(e?.message || 'Delete node failed')
    }
  }, [loadData])

  // Refresh only — the LEFT RAIL DELETE ALL / NUKE ALL panels perform the actual
  // destructive calls (deleteAll / nukeAll) with their own ARM→CONFIRM flow.
  const handleDeleteAll = useCallback(async () => {
    setSelectedNode(null)
    await loadData()
  }, [loadData])

  // DetailPanel's quick-access "Delete All" button: single window.confirm then
  // the remote wipe, then refresh the graph.
  const handleDeleteAllWithConfirm = useCallback(async () => {
    if (!window.confirm('Delete ALL correlation graph data (persona links + dependent rows)?')) return
    try {
      await deleteAll()
      setSelectedNode(null)
      await loadData()
    } catch (e) {
      alert(e?.message || 'Delete all failed')
    }
  }, [loadData])

  // Optimistic job-list entry for any scan dispatch — shared by the
  // "INITIATE DEEP-WEB RECON" modal and the dossier "DISPATCH SCAN" button.
  const handleJobCreated = useCallback((job) => {
    setJobs(prev => [
      {
        job_id: job.job_id || `job-${Date.now().toString().slice(-4)}`,
        target: job.target_value,
        status: 'running',
        terms: job.seed_terms.length,
        depth: job.max_depth,
        enqueued_at: new Date().toISOString(),
        pagesScraped: 1,
        safetyShield: 'ACTIVE',
      },
      ...prev,
    ])
  }, [])

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100vh', width: '100vw', background: 'var(--bg-void)' }}>
      <Header
        isLiveApi={isLiveApi}
      />

      <div style={{ display: 'flex', flex: 1, overflow: 'hidden' }}>

        <LeftControlRail
          onOpenScanModal={() => setScanModalOpen(true)}
          graphData={rawGraphData}
          media={media}
          onDeleteAll={handleDeleteAll}
        />

        <main style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden', minWidth: 0, position: 'relative' }}>
          <StatsBar stats={stats} />

          <div
            ref={graphContainerRef}
            style={{
              flex: 1,
              position: 'relative',
              overflow: 'hidden',
              width: '100%',
              height: '100%',
              background: 'radial-gradient(circle at center, #09121f 0%, #03060a 100%)',
            }}
          >
            {/* Scanline CRT Ambient animation */}
            <div style={{ position: 'absolute', inset: 0, pointerEvents: 'none', zIndex: 1 }}>
              <div style={{
                position: 'absolute',
                left: 0,
                right: 0,
                height: 50,
                background: 'linear-gradient(to bottom, transparent, rgba(0, 255, 157, 0.04), transparent)',
                animation: 'scan-line 6s linear infinite',
              }} />
            </div>

            {/* Canvas HUD Controls Overlay */}
            <div style={{
              position: 'absolute',
              top: 14,
              left: 14,
              zIndex: 10,
              display: 'flex',
              gap: 6,
              background: 'rgba(5, 8, 14, 0.85)',
              border: '1px solid var(--border)',
              borderRadius: 5,
              padding: 4,
            }}>
              <button
                onClick={() => fitGraphView(600, 44)}
                style={{
                  padding: '5px 8px',
                  background: 'transparent',
                  border: '1px solid var(--border)',
                  borderRadius: 3,
                  color: 'var(--text-secondary)',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 10,
                  cursor: 'pointer',
                }}
              >
                🎯 FIT VIEW
              </button>
              <button
                onClick={() => setLayoutMode(layoutMode === 'organic' ? 'hierarchical' : 'organic')}
                style={{
                  padding: '5px 8px',
                  background: 'transparent',
                  border: '1px solid var(--border)',
                  borderRadius: 3,
                  color: layoutMode === 'organic' ? 'var(--text-secondary)' : 'var(--green)',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 10,
                  cursor: 'pointer',
                }}
              >
                {layoutMode === 'organic' ? '🌳 HIERARCHICAL' : '◯ ORGANIC'}
              </button>
              <button
                onClick={() => handleReLayout()}
                style={{
                  padding: '5px 8px',
                  background: 'transparent',
                  border: '1px solid var(--border)',
                  borderRadius: 3,
                  color: 'var(--text-secondary)',
                  fontFamily: 'var(--font-mono)',
                  fontSize: 10,
                  cursor: 'pointer',
                }}
              >
                ↺ RE-SIMULATE
              </button>
              {selectedNode && (
                <button
                  onClick={() => setSelectedNode(null)}
                  style={{
                    padding: '5px 8px',
                    background: 'rgba(255,23,68,0.1)',
                    border: '1px solid rgba(255,23,68,0.3)',
                    borderRadius: 3,
                    color: 'var(--red)',
                    fontFamily: 'var(--font-mono)',
                    fontSize: 10,
                    cursor: 'pointer',
                  }}
                >
                  ✕ CLEAR FOCUS
                </button>
              )}
            </div>

            {/* Graph Telemetry HUD (Bottom-Right) */}
            <div style={{
              position: 'absolute',
              bottom: 12,
              right: 16,
              zIndex: 10,
              fontFamily: 'var(--font-mono)',
              fontSize: 9,
              color: 'var(--text-muted)',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'flex-end',
              gap: 2,
              pointerEvents: 'none',
              background: 'rgba(4, 7, 12, 0.7)',
              padding: '4px 8px',
              borderRadius: 4,
              border: '1px solid rgba(21, 34, 51, 0.5)',
            }}>
              <div>GRAPH ENGINE: 2D VECTOR HUD</div>
              <div style={{ color: 'var(--green)' }}>NODES: {graphData.nodes?.length || 0} • EDGES: {graphData.links?.length || 0}</div>
            </div>

            {/* ForceGraph2D Canvas Component */}
            {dim.w > 0 && dim.h > 0 && (
              <ForceGraph2D
                ref={graphRef}
                graphData={seededGraphData}
                width={dim.w}
                height={dim.h}
                minZoom={0.55}
                maxZoom={2.4}
                backgroundColor="transparent"
                nodeCanvasObject={paintNode}
                nodePointerAreaPaint={(node, color, ctx) => {
                  if (!node || typeof node.x !== 'number' || typeof node.y !== 'number' || !isFinite(node.x) || !isFinite(node.y)) return
                  const r = node.type === 'actor' ? 14 : node.type === 'stylometry' ? 11 : 9
                  ctx.beginPath()
                  ctx.arc(node.x, node.y, r + 6, 0, 2 * Math.PI)
                  ctx.fillStyle = color
                  ctx.fill()
                }}
                linkCanvasObject={paintLink}
                onNodeClick={node => {
                  if (node?.id) {
                    handleSelectNode(selectedNode === node.id ? null : node.id)
                  }
                }}
                onNodeHover={node => setHoveredNode(node ? node.id : null)}
                onLinkHover={link => setHoveredLink(link ? link : null)}
                onBackgroundClick={() => setSelectedNode(null)}
                cooldownTicks={200}
                onEngineStop={() => {
                  const nodes = graphData.nodes || []
                  nodes.forEach(n => {
                    if (typeof n.x === 'number' && isFinite(n.x)) n.fx = n.x
                    if (typeof n.y === 'number' && isFinite(n.y)) n.fy = n.y
                    n.vx = 0
                    n.vy = 0
                    if (n.id) nodePosCacheRef.current.set(n.id, { x: n.x, y: n.y })
                  })
                  if (graphRef.current && typeof graphRef.current.d3AlphaTarget === 'function') {
                    graphRef.current.d3AlphaTarget(0)
                  }
                  if (autoFitPendingRef.current && graphRef.current) {
                    fitGraphView(500, 44)
                    autoFitPendingRef.current = false
                  }
                }}
                linkDirectionalParticles={l => l.type === 'persona' ? 3 : 0}
                linkDirectionalParticleWidth={2.5}
                linkDirectionalParticleColor={() => '#00ff9d'}
                linkDirectionalParticleSpeed={0.005}
                d3VelocityDecay={0.1}
              />
            )}
          </div>
        </main>

        <DetailPanel
          detail={detail}
          graphData={graphData}
          onSelectNode={handleSelectNode}
          onClose={() => setSelectedNode(null)}
          onDeleteNode={handleDeleteNode}
          onDeleteAll={handleDeleteAllWithConfirm}
          onScanDispatched={handleJobCreated}
        />

        <RightRail
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          graphData={graphData}
          onSelectNode={handleSelectNode}
          torStatus={torStatus}
          torLive={torLive}
          aiMode={aiMode}
          aiRestarting={aiRestarting}
          onRestartStack={handleRestartStack}
          onMergeGraph={handleMergeGraph}
        />
      </div>

      <NewScanModal
        isOpen={isScanModalOpen}
        onClose={() => setScanModalOpen(false)}
        onJobCreated={handleJobCreated}
      />

      <NotificationToast notification={notice} onDismiss={() => setNotice(null)} />

      {/* Full-screen evidence viewer for orphan media rows (no owning actor). */}
      {mediaViewer && (
        <div
          onClick={() => setMediaViewer(null)}
          style={{
            position: 'fixed',
            inset: 0,
            background: 'rgba(0,0,0,0.7)',
            zIndex: 2500,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: 16,
          }}
        >
          <div
            onClick={(e) => e.stopPropagation()}
            style={{
              width: 'min(880px, 96vw)',
              maxHeight: '90vh',
              overflow: 'auto',
              background: 'linear-gradient(180deg, #070c14 0%, #04070c 100%)',
              border: '1px solid var(--border)',
              borderRadius: 8,
              boxShadow: '0 24px 80px rgba(0,0,0,0.6)',
              padding: 12,
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 8, marginBottom: 8 }}>
              <div style={{ fontFamily: 'var(--font-mono)', fontSize: 10, color: 'var(--text-muted)' }}>EXTRACTED EVIDENCE MEDIA</div>
              <button
                onClick={() => setMediaViewer(null)}
                style={{ background: 'none', border: 'none', color: 'var(--text-muted)', cursor: 'pointer', fontSize: 16 }}
              >
                ✕
              </button>
            </div>
            <img
              src={mediaViewer}
              alt="evidence"
              style={{ width: '100%', height: 'auto', display: 'block', borderRadius: 6, border: '1px solid rgba(255,255,255,0.06)' }}
            />
            <div style={{ marginTop: 8, fontFamily: 'var(--font-mono)', fontSize: 10, color: 'var(--text-muted)', wordBreak: 'break-all' }}>
              {mediaViewer}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

