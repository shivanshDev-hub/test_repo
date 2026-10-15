## ExoneraTor

### What It Is

ExoneraTor is an official Tor Project web service (hosted at `exonerator.torproject.org`) that answers one narrow, specific question: **"Was there a Tor relay running on this IP address on this date?"** It maintains a database of IP addresses that have been part of the Tor network and answers whether a Tor relay was running on a given IP on a given date.

The name is literal — it was built to *exonerate* people. If law enforcement traces malicious traffic back to an IP and that IP turns out to have been a Tor exit relay at the time, the traffic likely didn't originate from that machine's owner — it was just Tor traffic passing through on its way out to the internet. ExoneraTor lets anyone verify that claim independently instead of taking a relay operator's word for it.

### Data Source — It Doesn't Scan Anything Itself

ExoneraTor has no crawler and does zero live network probing. It's a pure database-lookup service. Architecturally it's small: one part fetches recent descriptors from CollecTor and imports them into a database, and another part listens for web requests and answers them by querying that database.

**CollecTor** is Tor's own metrics data-archival service. It continuously pulls two kinds of official Tor documents:
- **Relay/server descriptors** — self-published by each relay, stating its own IP, capabilities, and exit policy
- **Network status consensus documents** — published periodically by Tor's directory authorities (a small trusted set of relays), which is the authoritative "here's every relay that was verified and running right now" snapshot of the whole network

ExoneraTor imports both into a **PostgreSQL database**, keyed by IP, date, and relay fingerprint.

### What Gets Stored Per Relay

ExoneraTor may store more than one IP address per relay, since a relay can use a different IP for exiting to the internet than the one it registers with the Tor network, and it also stores whether that relay permitted transit of Tor traffic to the open internet at that time. The specific fields collected are IP address, relay type, nickname, and fingerprint.

So a single lookup can tell you not just "yes, Tor relay" but:
- Which relay (nickname + fingerprint — a unique identity across time even if the IP changes)
- Whether it was an **exit relay** specifically (i.e., could it have been the last hop responsible for traffic hitting the open internet) vs. just a guard/middle relay
- Whether a *different* IP was associated with the same relay identity on other dates

### Query Flow

1. You submit an IP address + a date via the web form (or its API).
2. The service queries its imported descriptor/consensus data for that date range.
3. It returns a match/no-match verdict, plus supporting detail (relay fingerprint, nickname, exit-policy status) if there's a hit — this is presented as evidence you can cite, not just a boolean.

### Why This Matters for Our PS

`for our attribution-confidence pipeline: use this as a negative control before scoring`

For a de-anonymization/attribution pipeline like the one PS 26151 wants, ExoneraTor is a **negative control**, not a discovery tool. Before our system assigns an attribution-confidence score linking an IP to a threat actor, we'd want to check: *was this IP actually a Tor relay at the relevant time?* If yes, that connection could just be Tor traffic passing through — attributing it directly to the operator would be a false positive. This maps cleanly onto the "attribution confidence" field our Expected Solution is required to output — ExoneraTor-style verification is exactly the kind of check that should lower a false attribution's confidence score.

### Limitations

- Only knows about the **public relay network** (relays, exit nodes, bridges) — it has nothing to say about hidden services or onion addresses directly, since those aren't part of the consensus/descriptor data it ingests.
- Coverage depends entirely on CollecTor's archive completeness — if a date/IP combination isn't in the archived descriptors, the answer is "no match found," which is not the same as "definitely not a relay."
- It rules things **out**, it doesn't find anything new — it's a verification step you bolt onto a pipeline, not a source of investigative leads on its own.

---

## Kilos

### What It Is

Kilos is a dark web search engine that emerged in **November 2019**, explicitly positioned as a successor to **Grams** — an earlier dark web search engine that shut down in December 2017 after its operator found it too hard to maintain, and whose associated Bitcoin mixer ("Helix") led to a federal money-laundering indictment. Kilos followed the same two-part playbook: a search engine paired with its own mixer service, which is part of why researchers (Digital Shadows) suspected overlapping developers between the two projects.

### How the Search Actually Works — It's an Aggregator, Not a Crawler in the OnionScan/TorBot Sense

This is the key architectural distinction worth noting: Kilos doesn't crawl generic `.onion` pages and parse HTML the way TorBot or Ahmia do. Its predecessor Grams worked by using **custom APIs against individual dark markets** to pull structured product/vendor/listing data directly, rather than scraping raw pages — Kilos is understood to follow the same model. So the pipeline looks more like:

1. **Per-market/forum integration** — for each of the sources it covers (at various points: markets like Empire, Cryptonia, Samsara, Versus, plus multiple forums), it pulls structured data: listings, vendor profiles, forum posts, and reviews.
2. **Centralized indexing** — all of that gets normalized into one searchable index so a single query returns matches across every connected market/forum simultaneously, instead of you visiting each market separately.
3. **Query-time search** — a user types a product, service, or vendor name, and the query is run across every indexed source at once, similar in spirit to how Ahmia centralizes results, but scoped to markets/forums instead of general onion content.

### The CAPTCHA-as-Training-Data Mechanism (the genuinely unique part)

`this dual-purpose CAPTCHA trick is worth stealing for our own quality-ranking layer`

This is Kilos' most architecturally interesting feature. It uses a custom CAPTCHA that asks users to **rank a piece of vendor/product feedback as positive or negative** before they can proceed. This does two jobs at once:
- **Anti-scraping/anti-DoS gate** — the stated reason, keeping competitors and automated scrapers out.
- **Labeled training data collection** — every answer a human gives becomes a labeled example fed into a **sentiment classifier**, which then scores and ranks search results by the sentiment of their reviews, surfacing the highest-quality/most-trusted listings first.

So the CAPTCHA isn't just friction — it's a crowdsourced labeling pipeline that continuously improves the ranking model, at zero labeling cost to the operator.

### Search & Filter Capabilities

- **Product/service search** across all indexed markets and forums
- **Vendor search** — find a specific seller's presence across multiple markets
- **Bitcoin address search** — look up a specific BTC address across indexed listings/posts
- **Contact search** — find posts/listings tied to a specific handle or contact method
- **Filtering** by payment method (BTC/LTC/XMR), shipping origin/destination country, and currency

### Associated Infrastructure

- **Krumble** — Kilos' own Bitcoin mixer/tumbler, randomizing transaction fees, commission fees, and adding a randomized delay to break timing correlation. Runs Tor-only.
- **Live chat** — direct messaging between registered users, and between users and the site admin.

### Why This Matters for Our PS (Pillar 2 specifically)

`this is literally our relationship-graph requirement, just built by the other side`

Kilos is a real-world, functioning example of exactly the capability our PS asks us to build defensively: **aggregating vendor identity signals (handles, contact info, wallet addresses) across multiple independent marketplaces into one queryable index.** The difference is intent — Kilos does this to help buyers *find* vendors; our system needs to do the same cross-market correlation to help investigators *attribute* vendors. Structurally, "one normalized index built from per-source structured pulls, searchable by handle/wallet/contact across all sources at once" is the same shape as what our Expected Solution describes for the relationship-graph requirement.

### What's Genuinely Not Known (don't overstate this in the report)

Unlike Ahmia, Robin, or OnionScan, Kilos is **not open source** and has no published architecture whitepaper — it's a clandestine, likely single-operator or small-team criminal service. Everything above comes from external researcher analysis (Digital Shadows, SecurityWeek) observing its *behavior*, not from inspecting its code. So: the exact backend stack, database choice, and whether it's still live today are **not independently verifiable** the way an open-source tool's would be. Don't cite it in the report as having a known technical architecture — cite it as a documented case study of what cross-market correlation looks like when built by a threat actor instead of a defender.
