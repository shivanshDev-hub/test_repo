# MidnightBazaar — FICTIONAL Demo Target for OnionEye

A make-believe marketplace **generated purely as an OnionEye scanner test target**.
Every handle, wallet, email, PGP block, and onion address on these pages is
fabricated. There is no real market, vendor, or key.

## What this folder is

During the documentation reorganization this demo site was moved out of the
repository's `docs/` folder (now reserved for real project documentation) into
`sample-target/midnightbazaar/`.

The four HTML pages were originally written to be served by GitHub Pages via the
`/docs` path. That no longer applies, so to scan the target you should serve the
folder locally and point the scanner at `http://localhost:8000`:

```bash
cd sample-target/midnightbazaar
python3 -m http.server 8000
```

Internal links are **absolute** on purpose: the OnionEye crawler only follows
extracted absolute `http(s)://` URLs, so sub-pages get discovered and depth-2
crawling exercises the BFS frontier.

## Scan it

```bash
curl -s -X POST http://localhost:8080/api/v1/investigations \
  -H 'Content-Type: application/json' \
  -d '{
    "analyst_id": "demo-analyst",
    "authorization_ref": "AUTH-DEMO-0001",
    "target_type": "clearnet_site",
    "target_value": "http://localhost:8000/",
    "target_host": "localhost",
    "seed_terms": ["MidnightBazaar", "DarkVault", "CipherMist"],
    "max_depth": 3
  }'
```

Then check the graph: `GET http://localhost:3000/api/v1/dashboard/graph`.

## What scanning it demonstrates

The pages are deliberately wired so OnionEye produces a **connected** entity graph:

- **3 actor handles** extracted from `<span class="username">` + `@mention`:
  `DarkVault`, `CipherMist`, `NightOwl`, plus `QuantumQuill` on the contact page.
- **Shared escrow BTC wallet** `1Je9pEjppgGYcPLB24ye9L6uw3wq` on both the home
  page (`DarkVault`) and the listing page (`CipherMist`)
  → persona link confidence **0.80**.
- **Shared support email** `support@darkvault.pro` on home (`DarkVault`) and the
  thread (`NightOwl`) → persona link confidence **0.65**.
- **Shared PGP key block** on listing (`CipherMist`) and thread (`NightOwl`)
  → persona link confidence **0.90**.
- **Onion v3 address** + **`nginx/1.24.0` server banner** → infra indicator nodes.
- `QuantumQuill` shares nothing → isolated leaf node (proves both dense clusters
  and standalone nodes render).

The person-link triangle `DarkVault — CipherMist — NightOwl` is exactly the
"linkages" output the scanner could not produce for the GitHub-Pages portfolio
page, because that page had **zero extractable identifiers** — this site fixes
that.