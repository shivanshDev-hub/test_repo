# OnionEye — Documentation Index

This folder is the single source of truth for OnionEye documentation. Every
guide, architecture doc, summary, and manual from across the repository has been
consolidated here. All paths are relative and resolve both on GitHub and locally.

---

## 📑 Repository Docs

| Doc | What it covers |
|-----|----------------|
| [Contributing](CONTRIBUTING.md) | Contribution workflow, testing expectations, PR process |
| [Developer Guide](DEVELOPER_GUIDE.md) | Getting started, dev environment, common workflows |
| [Security Policy](SECURITY.md) | Vulnerability reporting and security expectations |
| [Team Structure](TEAM_STRUCTURE.md) | Team assignments, roles, and responsibilities |

---

## 🏛️ Architecture

| Doc | What it covers |
|-----|----------------|
| [Problem Statement & Features](architecture/problem-statement.md) | SIH 2026 PS 26151 — dark web threat actor de-anonymization scope |
| [Decided Architecture — Overview](architecture/decided/overview.md) | Final chosen architecture, at a glance |
| [Decided Architecture — Subsection 1](architecture/decided/subsection-1.md) | Hardened input processing & safety pipeline |
| [Decided Architecture — Subsection 2](architecture/decided/subsection-2.md) | Redis → Postgres flow (Shivansh & Shubham work) |
| [Decided Architecture — Subsection 3](architecture/decided/subsection-3.md) | Phases 8–10: content filter, correlation engine, dashboard/API/export |
| [Decided Architecture — Legacy Spec](architecture/decided/legacy.md) | Original architecture spec (kept for reference) |
| [Suggestion — Reshit](architecture/suggestions/reshit.md) | Dark web threat actor attribution system proposal |
| [Suggestion — Jinesh & Rudra](architecture/suggestions/jinesh-and-rudra.md) | Hardened input processing & safety pipeline proposal |
| [Suggestion — Shubham](architecture/suggestions/shubham.md) | System architecture engineering draft |
| [Research — Tor & dark web](architecture/research/pre-existing-0.md) | How Tor and the dark web work |
| [Research — Pre-existing tools](architecture/research/pre-existing-1.md) | How pre-existing dark web tools work |
| [Research — ExoneraTor & friends](architecture/research/pre-existing-2.md) | ExoneraTor deep dive |
| [ARCHITECTURE.pdf](architecture/research/ARCHITECTURE.pdf) | Original architecture submission (PDF) |

Related assets: [official logo](architecture/assets/official_logo.svg) ·
[architecture diagrams](architecture/assets/)

---

## 🛠️ Development & Operations

| Doc | What it covers |
|-----|----------------|
| [AI Integration & Setup](development/ai-integration.md) | AI/ML components, model setup, and integration |
| [Workflow Setup](development/workflow.md) | Workflow configuration for the project |
| [Complete Tool Usage Manual](guides/tool-usage.md) | Full walkthrough of every OnionEye tool |
| [Troubleshooting Guide](guides/troubleshooting.md) | Common issues and fixes |
| [SQL Commands](database/sql-commands.md) | Reference SQL for the Postgres schema |

---

## 🧪 Testing

| Doc | What it covers |
|-----|----------------|
| [Test Suite Overview](testing/README.md) | Test organization, quick commands, and notices |
| [Quick Start — Running Tests](testing/quick-start.md) | First-time test setup and run instructions |
| [Testing Guide](testing/guide.md) | Patterns & best practices for writing tests |
| [Test Maintenance Guide](testing/maintenance.md) | When/how to update tests on schema/API changes |

---

## 📋 Work Summaries

| Doc | What it covers |
|-----|----------------|
| [CI/CD Pipeline & Test Infrastructure](summaries/ci-cd.md) | CI pipeline, test files created, metrics |
| [Automatic Test Execution](summaries/auto-test.md) | Auto test-run workflow — implementation summary |

---

## 🎯 Sample Target (demo site)

The fictional **MidnightBazaar** scanner demo target lives in
[`sample-target/midnightbazaar/`](../sample-target/midnightbazaar/README.md).
It is **not** part of the real documentation — it is generated test content for
the OnionEye crawler.

---

_Index maintained by hand — add new docs to this table when they land in `docs/`._