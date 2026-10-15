# How pre-existing tools work

> offer working mechanisms; study of systems that work with the dark web

> I will particularly flag the features that should be added to the project 
---
## TAILS-OS 

markdown

### Tails Architecture Overview

Tails (The Amnesic Incognito Live System) is built around one goal: leave no trace and force all traffic through Tor. Here's how it's put together.

### Boot and Runtime Model

Tails is a live OS — it boots entirely from a USB stick (or DVD) into RAM, rather than installing to and running from a hard drive.

- **Debian base**: Tails is a Linux distribution built on top of Debian, using GNOME as the desktop environment.
- **squashfs + read-only root**: The core system lives in a compressed, read-only squashfs image on the USB. At boot, this is loaded and layered with a **tmpfs/aufs (or overlayfs) union filesystem**. Any writes during a session go to the RAM-backed overlay, not to the USB's squashfs. This is what makes the system read-only by default and easy to reason about — the base image never mutates during use.
- **RAM-only runtime**: Because nearly everything runs in RAM, when you shut down, that memory is wiped. Tails also actively overwrites RAM on shutdown (via a kernel module / `sdmem`-like approach) specifically to defeat cold-boot attacks that try to recover encryption keys from residual RAM charge.
- **No swap**: Tails disables swap partitions entirely so nothing sensitive ever gets paged to disk unencrypted.

### Networking Architecture — the Core of the Design

`should add this to project.`
- **Tor as mandatory transport**: Tails runs a local Tor client (`tor` daemon), and the system is configured so that essentially *all* outbound network connections are forced through it. This isn't just "an app uses a SOCKS proxy" — it's enforced at a lower level.

`A killswitch to prevent leakage should add this`
- **Netfilter/iptables firewall**: Tails ships a strict firewall ruleset that blocks all outgoing traffic except that destined for the Tor process (identified by the `debian-tor` user/UID). Any application, or any malware, trying to phone home directly over clearnet gets blocked at the kernel packet-filtering level, not just at the application config level. This is the safety net for apps that don't know about Tor.
- **Stream isolation**: Different applications (Tor Browser, Thunderbird via Torbirdy, etc.) are routed through separate Tor circuits so that traffic from one app can't be correlated with another via shared circuits.

`should add this too`
- **MAC address spoofing**: By default, Tails randomizes the network interface's MAC address at boot, to avoid a physical network being able to fingerprint your specific hardware across sessions or locations.

`must add this too`
- **No direct DNS leaks**: DNS resolution is also routed through Tor rather than hitting the local ISP resolver.

### Persistence (Optional, Encrypted)

Since the live system is amnesic by design, Tails offers an optional **Persistent Storage** feature:

- A separate partition on the USB stick, encrypted with **LUKS**.
- You opt in to which categories persist (browser bookmarks, GPG keys, SSH keys, dotfiles, additional software, etc.) — nothing persists by default.
- This is mounted and decrypted only after you enter a passphrase at boot; the rest of the system remains stateless.

### Application Sandboxing and Hardening

- **AppArmor**: Confines high-risk applications (Tor Browser, PDF viewers, Thunderbird, etc.) with mandatory access control profiles, limiting what files/resources they can touch even if compromised.
- **Tor Browser**: A hardened Firefox-ESR build (shared with the Tor Browser Bundle) with fingerprinting resistance, uniform screen-resolution letterboxing, and no browsing history persistence.
- **Kernel hardening**: Various kernel hardening flags to reduce exploitability.
- **Verified boot artifacts**: Tails images are signed, and the project encourages verifying signatures/checksums before writing the USB, to prevent a tampered image from being installed in the first place.

### Update Mechanism

- Tails uses an **incremental update system (IUKs — Incremental Update Kits)** so a persistent Tails USB can upgrade itself without re-downloading the entire image, while still cryptographically verifying each update.

### Threat Model, in Short

The architecture is optimized for: leaving no forensic trace on the host machine's disk, forcing all network traffic through Tor to prevent IP/location leaks, and resisting both network-level and local (physical access) adversaries — at the cost of usability and persistence, both of which are only available as deliberate, narrow opt-ins.

---
## WHONIX 

markdown

### Whonix Architecture Overview

Whonix takes a fundamentally different architectural approach than Tails — instead of one live OS doing everything, it splits the system into two isolated virtual machines, each with a single job.

#### The Two-VM Split

- **Whonix-Gateway**: A VM whose only job is running Tor and routing all traffic through it. It has two network interfaces — one connected to the outside internet (through the host's network), and one connected to an *internal, isolated virtual network* that only the Workstation VM can see.
- **Whonix-Workstation**: The VM where you actually do your work — browsing, email, whatever. Critically, this VM has **no direct connection to the outside network at all**. Its only network route is through the internal virtual network to the Gateway.

This means the Workstation is physically incapable of reaching the internet except by going through Tor on the Gateway. There's no firewall rule to misconfigure or bypass — the isolation is structural, enforced by the virtual network topology itself, not by application-level policy.

#### Why This Matters: IP/Leak Protection

`way stronger ip leak prevention architecture`
- Even if an application on the Workstation is compromised, has a bug, or is flat-out malicious, it **cannot discover your real IP address**, because it has no interface that can see the real network. It only knows about the internal 10.x.x.x-style isolated network.
- Root compromise of the Workstation still doesn't leak your IP, because the leak-prevention lives one layer down, in the Gateway/network architecture — not in the Workstation's own OS config.
- Contrast with Tails: if a Tails firewall rule fails or is bypassed by some 0-day, the whole OS is on the same network stack, so a leak is more plausible. Whonix's isolation is at the virtualization boundary, which is a stronger boundary.

#### Base System and Virtualization

- Both VMs run **Debian** under the hood.
- Whonix is designed to run on top of a hypervisor — traditionally **VirtualBox**, though it also supports **KVM/QEMU** and has a **Qubes OS** integration (where Gateway and Workstation become separate Qubes VMs, gaining Xen-level isolation instead of VirtualBox-level).
- Unlike Tails, Whonix is **not amnesic by default** — it's meant to be a persistent, ongoing-use system with regular VM snapshots/state, not a boot-and-forget live environment. (You can still combine it with encrypted host storage, or reset VMs, but persistence is the default posture.)

#### Tor Integration Details

- **Stream isolation**: Similar to Tails, the Gateway can isolate different types of traffic into separate Tor circuits, often orchestrated via `tor` control-port automation so different applications or "personas" don't share circuits.

`assume zero architecture:assume compromise`
- **Control Port Filter Proxy (CPFP)**: The Gateway exposes a filtered, restricted view of the Tor control port to whichever tools need to interact with Tor (e.g., to request new circuits), without giving full raw control-port access, which would be dangerous if the Workstation were compromised.

`Another way to protect the DNS.`
- **No direct DNS leaks**: DNS resolution happens on the Gateway side, transparently proxied through Tor, so the Workstation never does its own DNS resolution.

#### Hardening

- **Kernel hardening / grsecurity-derived patches** (historically) applied to reduce exploitability of both VMs.
- **AppArmor** confinement of individual applications inside the Workstation, similar in spirit to Tails.
- **Time sync obfuscation**: Whonix uses `sdwdate` (secure distributed web date) instead of NTP, since NTP queries can fingerprint or leak timing info; `sdwdate` cross-checks time against multiple HTTPS sources over Tor instead.
- **Software diversity / multiple Workstation flavors**: You can run different Workstation VMs (e.g., one for browsing, one for other tasks) against the same Gateway, so a compromise in one persona doesn't necessarily contaminate another — closer to compartmentalization than Tails' single-OS model.

#### Threat Model Comparison with Tails

| | Tails | Whonix |
|---|---|---|
| Isolation mechanism | Firewall rules within one OS | Separate VMs, isolated virtual network |
| Amnesia | Default (RAM-wiped on shutdown) | Not default (persistent VMs) |
| IP-leak protection if compromised | Depends on firewall holding | Structural — Workstation has no route out |
| Best for | Single-session anonymous use, minimal forensic trace | Long-term compartmentalized anonymous work |
| Physical machine trace | None (if using live USB properly) | VM images live on host disk unless separately encrypted |

They're often used together in practice — running Whonix Gateway+Workstation inside a Tails or an isolated host session, or running Whonix on Qubes for extra hypervisor-level isolation. The two projects solve overlapping but architecturally distinct problems: Tails optimizes for "leave nothing behind," Whonix optimizes for "even if something goes wrong, it can't see your real network."

`whonix also randomises the time so does not leak anything`

--
## Robin (dark web scrapper)

`This is exactly what we are writing, written and expanded by a person to the industry-grade standard.Must review this we can take lot of insperation from this`

Robin is an open-source, AI-powered OSINT tool for dark web threat intelligence, built by security researcher Apurv Singh Gautam. It automates what a human threat researcher would otherwise do manually over hours, via a modular pipeline with clean separation between search, scrape, and LLM stages.

#### Pipeline Stages

1. **Query Refinement**
   - You give it a plain-language query.
   - An LLM step expands/rewrites this into a semantically richer search phrase rather than relying on literal keyword matching, capturing intent rather than exact wording.

2. **Multi-Engine Search over Tor**
   - Simultaneously queries 15 different dark web search engines through the Tor network, aggregating results.
   - Requires a local Tor instance running so the tool can resolve and reach `.onion` search indexes.

3. **AI Filtering**
   - Raw result sets can be huge and mostly noise — from a pool of 900+ results, Robin uses AI to identify the most relevant 20-30 links, discarding the rest.
   - Exists because most of the dark web is dead links, honeypots, or scam duplicates.

4. **Threaded Scraping**
   - The filtered set of links is scraped in parallel.
   - CLI supports a configurable thread count for the scraping step (default: 5).
   - Scraping is routed through Tor.

5. **LLM Summarization / Reporting**
   - Scraped content is fed back through the LLM to generate a consolidated intelligence summary.
   - Supports conversational follow-up questions answered from that investigation's own data, without re-running the search.

#### Deployment / Runtime Model

- **CLI-first design**: Built for terminal usage and automation/threat-intel pipeline integration.
- **Web UI**: Also runs as a Streamlit web application inside a Docker container, interfacing with Tor for `.onion` access and with LLM providers for query processing and analysis.
- **Multi-model support**: Works with OpenAI, Claude, Gemini, Ollama, or any OpenAI-compatible API (LM Studio, llama.cpp, Groq, etc.) — the LLM backend is swappable.
- **Local LLM via Colab + ngrok**: Can run the LLM through Google Colab if local hardware isn't capable, exposing the Colab-hosted model to Robin via an ngrok tunnel, keeping data processing off third-party servers.
- **Docker-ready**: Optional Docker deployment for isolated usage.
- **Custom reporting**: Saves investigation output to file (markdown, e.g. for Obsidian) for reporting or further analysis.

#### Design Rationale

- Roughly 90% of dark web content is either dead, a law-enforcement honeypot, or a scam/duplicate.
- Professional researchers previously spent 6-8 hours manually navigating broken circuits and unreliable sites.
- Robin compresses that into roughly 30 minutes by pushing the "is this link even worth visiting" triage onto an LLM instead of a human.

#### Legal Note

Accessing or interacting with certain dark web content may be illegal

---

`This is literally the tool we are building.`

### Onion Scan 

OnionScan is a security auditing tool for hidden services, built by Sarah Jamie Lewis. Its core purpose is finding operational security mistakes that let researchers (or adversaries) deanonymize or correlate `.onion` services that are supposed to be unlinkable.

#### Core Architectural Idea

Most dark web tools answer "what's on this page." OnionScan answers a different question: **does this service leak anything that ties it to another service or to a real-world identity?** That reframing is what makes its architecture different — it's not a content indexer, it's a **fingerprint-and-correlate** engine.

#### Pipeline

##### 1. Target Fetch
Given a `.onion` address, it connects via a local Tor SOCKS proxy and pulls down the service's exposed surface — not just the homepage, but a set of commonly-misconfigured paths and protocols (HTTP, and historically also checks for exposed SSH, FTP, SMB, and other open services on the same hidden service).

##### 2. Fingerprint Extraction
This is the real engine. For each service, OnionScan pulls out a set of "identity artifacts" — pieces of technical fingerprint that shouldn't normally be exposed, but often are due to sloppy server configuration:

- **Server banners/headers** — Apache/nginx version strings, default error pages, `Server:` headers that reveal software stack and sometimes hostnames.
- **Exposed status pages** — misconfigured Apache `mod_status` or similar diagnostic endpoints that leak internal IPs, other virtual hosts on the same box, or process lists.
- **SSH host keys** — if a hidden service operator also exposes SSH on the same backend server, the SSH host key fingerprint is a strong, near-unique identifier.
- **PGP keys / Bitcoin addresses** embedded in page content — operators often reuse the same crypto address or PGP key across multiple sites.
- **Shared TLS certificates or self-signed cert fingerprints**, when present.
- **Analytics IDs, favicon hashes, file metadata (EXIF in images), directory listing artifacts** — anything left in static content that's structurally identical across sites.
- **Common misconfiguration**: services that are supposed to be Tor-only but accidentally also bind to a clearnet-facing IP (the single biggest deanonymization vector historically).

##### 3. Fingerprint Database / Correlation Store
Every artifact extracted from a scan is written into a **local database** (historically a simple key-value/JSON store), keyed by artifact type and value — SSH key fingerprint, Bitcoin address, banner string, etc. This is the architectural core: it's not per-scan stateless output, it's a **growing correlation index** across every service ever scanned.

##### 4. Cross-Service Correlation
When a *new* service is scanned, OnionScan checks its extracted fingerprints against everything already in the database. If Service B's SSH key fingerprint matches Service A's, or they share a Bitcoin address, or an identical leaked internal IP, OnionScan flags them as **operationally linked** — i.e., probably run by the same person/infrastructure — even though nothing on either site claims that association.

##### 5. Report Generation
Output is a structured report per scan (JSON or human-readable) listing every artifact found, severity of the leak (e.g., "exposed SSH host key" is high severity, "generic Apache banner" is low), and any cross-matches against the correlation database.

#### Why This Counts as a Different Architecture

A normal scraper's data flow:
```
fetch → parse content → extract text → store/index → done
```

OnionScan's data flow:
```
fetch → extract low-level protocol/config artifacts →
write to a persistent global fingerprint store →
diff against all prior scans → flag correlations → report
```
The second pattern is architecturally closer to how **digital forensics tools** or **threat-actor attribution systems** (e.g., linking malware samples by shared code artifacts, or linking phishing kits by shared infrastructure) are built — a growing correlation graph that gets more powerful with every additional data point, rather than a stateless per-target pipeline.

#### Dual-Use Consideration

OnionScan is explicitly built and published as a **defensive/audit tool** — it exists so hidden-service operators can scan their own infrastructure before an adversary does, and so researchers can study which classes of mistakes are common. The same correlation technique is openly acknowledged as dual-use: it's the same method a law enforcement or intelligence deanonymization effort would use. It's presented as "here's how easy it is to accidentally deanonymize yourself, go check your own service" rather than as an offensive tool, but the underlying mechanism doesn't care which side uses it.

---
## Onion Balance 

OnionBalance solves a problem unique to Tor's hidden-service design: how to load-balance traffic to an onion service when there's no reverse proxy, no DNS, and no IP address to point a load balancer at.

#### The Problem It Solves

A normal onion service is tied to a single Tor instance holding a single private key. If that one instance goes down, or gets overwhelmed with traffic, the whole `.onion` address is unreachable — there's no way to just "add more servers behind it" the way you would with a normal web app, because the onion address *is* the public key, not an IP you can route around.

#### How Onion Services Normally Work (the Piece OnionBalance Exploits)

Every onion service periodically publishes a **descriptor** to Tor's distributed hash table (the HSDir system) — a signed document saying "here's how to reach the service currently identified by this key: use these introduction points." Clients look up this descriptor by the onion address to learn which introduction points to use to establish a rendezvous circuit.

#### OnionBalance's Mechanism

Instead of one Tor instance owning the key, OnionBalance introduces a **management layer** that sits above multiple independent backend Tor instances:

1. **Shared master key, independent instances**
   - One onion service key pair is generated as usual.
   - Instead of a single Tor process holding it, several separate backend servers each run their own Tor instance with their own set of introduction points.

2. **Descriptor aggregation**
   - The OnionBalance management process periodically contacts each backend instance, collects the introduction points each has independently established, and merges them into a single descriptor.
   - Descriptors have a size limit (historically a handful of intro points per descriptor), so OnionBalance does randomized selection/rotation across a larger backend pool if there are more backends than fit in one descriptor.

3. **Publishing under the shared key**
   - The merged descriptor is signed with the master private key and published to the HSDir like a normal onion service descriptor — but it now advertises introduction points belonging to multiple, independent physical servers.

4. **Client-side effect**
   - When a client looks up the onion address, it receives the merged descriptor and picks one of the listed introduction points more or less at random.
   - Different clients — or the same client on different lookups — end up establishing rendezvous circuits to different backend servers, achieving load distribution purely through descriptor-level randomization.
   - Backends are never aware of each other beyond publishing their intro points to the OnionBalance manager.

#### Why This Is Architecturally Distinct

- There's no reverse proxy, no shared IP, no anycast, no traditional load-balancer component in the data path at all.
- Load balancing happens entirely at the **descriptor-publication layer of Tor's own directory protocol** — OnionBalance's job is just periodically stitching together intro-point lists and re-signing a descriptor.
- Actual user traffic never touches OnionBalance itself; once a client picks an introduction point, its rendezvous circuit goes straight to whichever backend Tor instance owns that intro point.
- OnionBalance is pure control-plane, zero data-plane — unlike conventional load balancers (HAProxy, nginx, cloud LBs), which all sit inline in the traffic path.

#### Operational Tradeoffs

- **Failure handling is soft, not hard**: if a backend goes down, OnionBalance doesn't detect and remove it in real time the way a health-checked LB would — the descriptor simply stops including that backend's stale intro points on the next refresh cycle.
- **Private key distribution risk**: since it involves the master key being usable to sign descriptors, key management (keeping the master private key only where the OnionBalance manager needs it, not on every backend) is the main operational security concern.
- **Descriptor size limits** cap how many backend introduction points can be advertised at once, which is why larger OnionBalance deployments rotate/select subsets rather than listing every backend every time.

---
`this is what we can implement in case the osint scan mecahnism fails`
## Tor bot

TorBot is a Python-based dark web crawler/OSINT tool. Architecturally, it follows a standard crawler pattern (crawler + parser + Tor proxy) rather than anything novel — it's a lightweight, Tor-aware version of a classic web crawler like Scrapy.

#### Core Purpose

TorBot maps the link structure of `.onion` sites and extracts surface-level metadata (emails, other onion links, page titles). It does not perform correlation or AI-based filtering — it's purely a discovery/enumeration tool.

#### Pipeline

##### 1. Tor Connectivity Layer
TorBot connects through a local Tor SOCKS proxy (typically `127.0.0.1:9050`, the default Tor daemon port) using a Python SOCKS-aware HTTP library. This is the only Tor-specific part of the tool — everything downstream is generic web scraping logic pointed through that proxy.

##### 2. Seed Fetch
Given a starting `.onion` URL, it issues an HTTP GET through the proxy and pulls down the raw HTML.

##### 3. HTML Parsing
Uses a standard HTML parser (BeautifulSoup, in most versions) to extract:
- **All `<a href>` links** on the page — the core crawling mechanic, discovering more onion addresses by following links out from known ones.
- **Email addresses** via regex pattern matching against the raw page text.
- **Page title/meta info** for basic classification of discovered links.

##### 4. Link Classification and Queueing
Discovered links are filtered into categories: valid `.onion` addresses (added to a crawl queue), clearnet links found on the page, and dead/malformed links. Valid onion links get queued for the next crawl depth.

##### 5. Recursive Crawl (Depth-Limited)
Steps 2-4 repeat for each queued link, up to a user-specified depth (how many hops from the seed URL to follow) — a standard breadth-first or depth-first crawl pattern, no different in structure from how any generic web crawler expands a frontier, just with Tor as the transport.

##### 6. Output
Results — the link graph plus extracted emails/metadata — are written to a file or printed, sometimes with basic visualization of the link tree in later versions.

#### Why This Is the "Standard" Architecture, Not a Novel One

Compare to OnionScan: it extracts protocol-level fingerprints (SSH keys, banners, Bitcoin addresses) and correlates them across a persistent database to link services together — solving a graph-correlation/deanonymization problem. TorBot extracts content-level links and just follows them — solving a discovery/enumeration problem.

Data flow:
```
fetch → parse HTML → extract links/emails → queue → repeat → dump results
```

This is structurally identical to `wget --recursive`, Scrapy, or any 2005-era web crawler — the only genuinely Tor-specific piece is the SOCKS proxy configuration at the network layer. It's a useful building block for understanding raw crawling mechanics before studying tools that layer AI filtering (Robin) or correlation (OnionScan) on top.

---
## Ahmia 

### What Ahmia Is

Ahmia is a search engine built specifically for **Tor hidden services** (`.onion` sites), rather than the regular web. It doesn't crawl the surface web like Google or Bing — it focuses entirely on content reachable only through the Tor network. It was founded in 2014 by Finnish security researcher Juha Nurmi, originally as a research project at the University of Helsinki, and remains free and open source.

### How It Works — Overview

### Discovery & crawling
Regular search engines crawl by following hyperlinks and sitemaps. Onion sites are deliberately hard to discover that way, so Ahmia's crawlers use a few different techniques instead:
- Pulling lists of active services from Tor directory services
- Following onion links found inside already-discovered sites' content
- Using trusted "seed" relays within the Tor network as starting points to bootstrap the crawl

### Indexing
Once a hidden service is found, the crawler parses its content and stores the title, description, and keywords alongside the `.onion` address in a searchable database. Like normal search engines, it also respects `robots.txt`, so site operators can opt out of being indexed.

### Access
You can search from a normal browser at ahmia.fi — no Tor required just to query the index. Ahmia itself also runs as an onion service for people already using Tor Browser. Either way, actually *opening* a result requires Tor Browser, since `.onion` addresses aren't reachable from the regular internet.

### Filtering
Ahmia deliberately curates rather than indexing everything it finds. Abuse material is banned from the index, and the project maintains a blacklist plus a reporting mechanism so flagged content gets removed. This is what distinguishes it from other Tor search tools that index everything discoverable regardless of legality.

## Under the Hood: Technical Architecture

Ahmia isn't one monolithic program — it's split into a few components that work together, all maintained as separate open-source repos under the `ahmia` GitHub organization.

**1. The crawler (`ahmia-crawler`)**
Goes out into the Tor network and does the actual discovery work described above — following links, querying directories, bootstrapping from seed relays.

**2. The index (`ahmia-index`)**
Crawled data is fed into **Elasticsearch**, which powers the actual search functionality — full-text indexing of titles, descriptions, and page content so queries return ranked results quickly.

**3. The website (`ahmia-site`)**
The public-facing layer: a **Django** (Python) app that queries Elasticsearch and renders results. It runs as both the clearnet site (ahmia.fi) and its own onion service, served in production via Gunicorn + Nginx.

### Step-by-step: what happens on a search
1. Crawler discovers an onion address and fetches its page over Tor.
2. It extracts text, title, and links — new links go back into the crawl queue; parsed content goes to the index.
3. Elasticsearch stores it as a searchable document.
4. You type a query on ahmia.fi → Django queries Elasticsearch → matching `.onion` results are returned with title/description snippets.
5. You open a result in Tor Browser, since ahmia.fi (Django/Nginx) doesn't host or proxy the actual onion content itself.

`this is what impresses me the most about this tool, we can implement this in our own tool to make it legal and moral`

### Filtering & maintenance layer
- A `banned_terms.map` file is loaded directly into the Nginx config to block serving/indexing of blacklisted material.
- A scheduled job (`deleteonions`) runs periodically to prune dead or removed onion entries, since hidden services churn constantly — addresses go offline or change often.

## Limitations Worth Knowing

- Ahmia doesn't host any of the content itself, doesn't verify who operates a listed service, and doesn't provide a complete map of the dark web.
- Its own project acknowledges results may be inaccurate, stale, or in some cases unsafe — every result is a lead, not a verified catalog entry.
- There have historically been phishing clones impersonating Ahmia's onion address, so the official site publishes the correct address for verification.

---
## DarkOwl

### What DarkOwl Is

DarkOwl is a commercial darknet intelligence platform, not a free public search tool like Ahmia. It's built for security teams, threat researchers, and law enforcement — the goal is actionable intelligence, not casual browsing. The company is Denver-based and markets itself as holding the world's largest index of darknet content, with an archive going back 8+ years.

### What It Collects
`This has a way higher extended search model which we can implement into our own tool, offering wider search than any other.`
DarkOwl casts a wider net than a typical dark web crawler:

- **Darknets**: Tor, I2P, and ZeroNet networks
- **Darknet-adjacent sites**: marketplaces, forums, paste sites, ransomware leak sites, criminal chat platforms
- **Authenticated platforms**: sources that require login credentials rather than just a Tor connection — e.g. Telegram
- **High-risk surface web**: sites outside the darknet proper that still carry relevant threat signal

Its search platform, Vision, lets analysts search by network specifically — Tor, messaging platforms, paste sites, etc. — rather than treating "the darknet" as one undifferentiated blob.

### How It Processes What It Finds

Raw scraped content isn't much use on its own, so DarkOwl runs it through two layers:

1. **Machine learning** — automated classification and structuring at scale
2. **Human analyst review** — verification and refinement on top of the automated pass

The database is updated from tens of thousands of sites daily. This automation-plus-human-review combination is what separates DarkOwl from a purely automated crawler like Ahmia: it's aiming to produce vetted intelligence, not just an index of pages that happen to exist.

### How Users Actually Query It
`Provides direct integration into OSINT tools that the Intel experts already use,something to copy/get inspired`
The main interface, **Vision UI**, offers search functionality well beyond a plain search box:

- **Free-text search** — works like a normal search engine, across up to 47 languages
- **Boolean / regex search** — precise, structured queries rather than plain keywords
- **Entity search** — search specifically for structured data types like email addresses, IP addresses, or cryptocurrency wallet strings
- **Monitors & Alerts** — set up a standing watch (e.g. a company domain or an executive's email) and get notified automatically when it surfaces in newly indexed content, instead of manually re-searching
- **DarkINT Score** — a computed risk score quantifying a given domain's overall dark web exposure, based on all the data points DarkOwl holds on it

### Safety Layer

Analysts can review flagged content without direct exposure to the underlying harmful images, files, or individuals — DarkOwl effectively interposes itself between the raw dark web and the person doing the investigating.

### Access Methods

- **Vision UI** — the primary web interface for analysts
- **APIs** — designed to plug directly into existing SOAR (security orchestration) and threat-detection workflows
- **Maltego integration ("transforms")** — lets investigators pull DarkOwl data directly into link-analysis graphs for tracing threat actors, campaigns, or stolen-data trails

### How It Differs From Ahmia

| | Ahmia | DarkOwl |
|---|---|---|
| Cost | Free, open source | Commercial |
| Scope | Tor `.onion` services only | Tor + I2P + ZeroNet + surface web + chat apps |
| Content policy | Actively filters out abuse material | Broad collection, filtered for analyst safety at review time |
| Purpose | Public search/discovery | Paid intelligence platform for security workflows |
| Output | Keyword search results | Raw data access + structured intelligence (scores, alerts, entity extraction) |

In short: Ahmia tells you an onion address exists and what it contains. DarkOwl aggregates darknet data at much greater scale and breadth, layers analyst judgment and enrichment on top, and hands it to you as an API or dashboard built for continuous security monitoring rather than one-off lookups.

---
## Steam 

`This is a Python library that makes implementing tor easier`

### What Stem Is

Stem is the Python controller library for Tor. It doesn't tunnel traffic or reimplement Tor itself — it lets a program *talk to* an already-running Tor process using Tor's own control protocol. Under the hood, it's a Python implementation of Tor's directory and control specifications, which is why anything Stem does could technically be done by hand over telnet against Tor's control port — Stem just wraps that raw protocol in a clean, higher-level interface.

It's maintained by the Tor Project and is the library behind tools like Nyx (formerly "arm"), the terminal-based Tor relay monitor.

### How You Connect to It

1. Tor must already be running and configured with a control port or control socket enabled.
2. The application gets a `Controller` instance — the most flexible ways are `from_port()` or `from_socket_file()`, though quicker convenience helpers exist for simple scripts.
3. The `Controller` class wraps a lower-level `ControlSocket`, keeping basic methods (connect, close, is_alive) while adding higher-level functionality for actually working with Tor.
4. The application authenticates itself against Tor before it's allowed to issue any commands.

### What You Can Do Once Connected

Once authenticated, Stem exposes Tor's control protocol as normal Python calls instead of raw socket text. Examples:

- **Request a new identity/circuit** — sending `SIGNAL NEWNYM`, the same command you'd type manually over telnet, wrapped in a clean Python call
- **Subscribe to live events** — `add_event_listener()` lets an application react to what Tor is doing in real time
- **Inspect circuits and relays** — `get_circuits()` shows which relays your connection is currently routed through, including identifying the exit node
- **Interactive exploration** — Stem ships its own control prompt, which behaves like a Python interpreter with tab-completion and colorized output, useful for exploring Tor's behavior live

### Why This Matters for Crawlers

Tools that need to actively manage their Tor connections — rotating circuits, avoiding overusing one exit node, or reacting to failed connections — need exactly this kind of low-level control rather than just tunneling blindly through whatever circuit Tor happens to pick. This is the same category of capability a Tor-aware crawler (like the one behind Ahmia) relies on when it needs deliberate control over how it moves through the network, rather than treating Tor as a black-box proxy.

### Where It Sits Relative to Raw Tor

Stem doesn't add capabilities beyond what Tor's control protocol already offers — it's a convenience layer. Its value is entirely in developer ergonomics: structured Python objects and methods instead of hand-parsed control-protocol text, plus utilities (event listeners, an interactive prompt) that make experimenting with a live Tor process practical.

