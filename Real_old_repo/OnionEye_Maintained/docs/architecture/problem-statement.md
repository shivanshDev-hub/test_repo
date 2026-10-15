# SIH 2026 — PS 26151: Dark Web Threat Actor De-anonymization

## Quick Facts

| Field | Value |
|---|---|
| PS ID | **26151** |
| Title | Dark Web Threat Actor De-anonymization |
| Organization | National Technical Research Organisation (NTRO) |
| Category | Software |
| Theme | Blockchain & Cybersecurity |
| Idea submission deadline | **20 September 2026** |
| Dataset provided by org | None listed |

*(Source: Smart India Hackathon 2026 official portal, sih.gov.in — cross-verified against independently scraped SIH 2026 datasets, last confirmed 21 Aug 2026)*

## Official Problem Statement (verbatim from sih.gov.in)

**Background**
The dark web has become a preferred operating space for threat actors mainly because Tor hidden services let them hide their identity, making attribution the core challenge in any investigation. These actors run drug/arms sales, stolen data and hacking services, money laundering, terror financing, etc. The goal is to build a system that de-anonymizes dark web threat actors and links them to real-world entities.

**Description**
The system must continuously gather footprints from marketplaces, forums, and deep web sources, and link them to identifying information. Three core capabilities are explicitly demanded:
1. Detect Tor hidden service misconfigurations (exposed server-status pages, SSL certs tied to clearnet domains, default service banners, descriptor inconsistencies) and match them to clearnet infrastructure to trace likely origin servers.
2. Map threat actors across multiple marketplaces into a single relationship graph of handles, PGP keys, wallets, and trust links.
3. Use AI (stylometric persona identification + behavioural profiling) to link rebranded/migrated personas back to known threat actors.

It must also have an analytical front end to query the data across a timeline, and run in autonomous mode.

**Expected Solution**
An end-to-end system for collection, storage, contextualization, and querying (via GUI/dashboard) of dark web threat intel — covering: actor profiles, identifiers (handles/PGP keys/wallets), hidden service infra indicators, persona linkages, attribution confidence score, category, last scan date, and source. Must support export in **CSV, JSON, and report formats**.

## What This Actually Means — Deliverables Breakdown

| # | Component | What it does |
|---|---|---|
| 1 | **Crawler/collector module** | Continuously pulls data from dark web marketplaces, forums, deep web sources (autonomous mode) |
| 2 | **Infra correlation engine** | Finds Tor misconfigs (SSL cert reuse, banners, server-status leaks, descriptor mismatches) → maps to real clearnet servers |
| 3 | **Identity graph builder** | Links handles, PGP keys, crypto wallets, trust relationships across marketplaces into one relational graph |
| 4 | **AI persona-linking layer** | Stylometry (writing-style fingerprinting) + behavioral profiling to catch the *same actor* under a new alias |
| 5 | **Database/data model** | Stores actor profile, identifiers, infra indicators, persona links, confidence score, category, last scan date, source |
| 6 | **Dashboard/GUI** | Timeline-based querying and visualization of the above |
| 7 | **Export module** | Output results as CSV, JSON, and formatted reports |

## What is NOT Specified (don't let your team assume this)
- No dataset is provided by NTRO — you'll need to build/simulate your own for the demo.
- No mandated tech stack, language, or framework.
- No contact/email listed on the PS yet.
- No legal/ethical operating guidelines given — since this touches real dark web scraping, plan to **demo against synthetic/mock dark web data**, not live illicit marketplaces, for both legality and safety.

---

# Features Suggestion Box

**To implement**
1. A kill switch mechanism just like Tails
2. MAC address spoofing just like Tails OS
3. Whonix two VM isolation model for Exploiation prevention
4. Abusive and illegal filtration like: Ahmia
5. A zero trust architecture and compartmentalization so compromise cannot reach the host system
6. Time sync obfuscation just like: Whonix
7. LLM query expanding just like: Robin
8. AI filtration for garbage data like: Robin
9. Misconfig exploitation just like: Onion Scan
10. Expanded scanning to include deep web + other places outside dark web too just like: Dark Owl
11. Enhanced search functionality just like DarkOWL
12. Can directly be intergated into exisiting OSINT tool like: Dark Owl
13. Maintaining a fingerprint database just like:  Onion Scan
14. A correlation graph generation just like: maltego
15. A Negative control filter with things like: ExoneraTor
16. slow crawling + human intervention for solving captchas  

**Legal boundaries to protect**
1. Preventing illegal demand queries within the query expander AI
-  some should also pentest the very AI refer to [my AI pentest guide](https://github.com/shivanshsahajpal22-dev/WebExploits/blob/main/Web_Exploit_Guide/Web_Manual_Testing_3.md)

2. Never touch CSAM, under any circumstance. , we have to prevent this on the scarper level to let scarper just record that CSAM was found but to never touch it

3. Passive OSINT only — no active intrusion. Just using the exploitable server's moconfig to collect what is already it already gives cause of that misconfig not something that will actively exploit the server: take care of this wherever you add the server exploitability part

4. No real-world naming without authority. You're a student team, not a law enforcement agency. Section 69 of the IT Act reserves interception/monitoring/decryption powers specifically for authorized government agencies through a defined legal process. Your system can produce an investigative lead with a confidence score for NTRO's own analysts to act on — it cannot itself publicly "expose" or accuse a real, identifiable person. That's the difference between a threat-intel tool and vigilante doxxing, and the latter carries personal defamation/harassment liability for whoever built it. Keep tool to only assumed correlation NOT "AHH.. THIS IS CORRELATED !"

5.Never actually crawl live illicit marketplaces for your demo. Use synthetic data, a mocked mirror environment, or a vetted academic darknet dataset instead. Beyond the CSAM risk above, actually browsing/interacting with real markets on personal or institutional devices/networks during a judged event creates real exposure for your team members individually — you have no government authorization shielding you. **We will set the tool to strictly scan only our created, with a kill switch if it ever malfunctions to start with real places on its own**

6. Data minimization — only ingest/store data plausibly tied to threat-actor attribution, not bulk-collect everything a crawler happens to see. Bystander data (ordinary anonymous users, unrelated forum chatter) sitting in your database is DPDP exposure for no benefit. **we would add a sheild filter in front of the very database for such filteration**

**Few warnings**
- Every output should carry a confidence score, never a definitive claim — your own PS explicitly asks for this field, so lean into it as your legal safety valve too.

- Build for human-in-the-loop review before any output leads to a real-world action — the tool surfaces leads, a trained analyst decides.

- No feature that auto-publishes or auto-notifies based on an attribution — false positives in stylometry/graph-correlation are common enough that an automated "naming" feature is a lawsuit waiting to happen.
