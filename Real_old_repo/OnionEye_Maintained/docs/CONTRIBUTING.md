# Contributing

> **📋 Team Structure:** See [TEAM_STRUCTURE.md](./TEAM_STRUCTURE.md) for current team assignments and contact information.

Thanks for taking the time to contribute. This document covers how work flows through this repo — branching, commits, reviews, and releases — so contributions land smoothly.

## Getting started

1. Clone the repo and check out `develop` — that's the active integration branch, not `main`.
   ```bash
   git clone <repo-url>
   cd <repo>
   git checkout develop
   ```
2. Set up commit signing (required — see below) before your first commit.

## Branching strategy

This repo uses a **develop/feature branch model**:

- **`main`** — always reflects production. Protected, only receives merges from `develop` (or emergency hotfixes).
- **`develop`** — the integration branch. All finished feature work lands here first.
- **`feature/*`**, **`fix/*`**, **`chore/*`** — one branch per unit of work, branched off `develop`.

```
main       ●────────────────●───────────────►  (releases only)
                │                    ▲
develop    ●────●────●────●────●────●────────►  (integration)
                     │         ▲
feature/x           ●──●──●──●              (your work)
```

**Naming convention:** `feature/short-description`, `fix/short-description`, `chore/short-description` (e.g. `feature/dark-mode`, `fix/login-timeout`). Keep it lowercase, hyphenated, no ticket-number-only names — the description should be readable on its own.

## Making a change

1. Sync and branch from `develop`:
   ```bash
   git checkout develop
   git pull
   git checkout -b feature/short-description
   ```
2. Make your changes, committing in small, logical chunks with signed commits (see below).
3. **Run tests locally before pushing** (see Testing section below).
4. Push your branch and open a pull request **targeting `develop`** — not `main`. Double-check the base branch in the PR compare view; it defaults to `main` and is easy to miss.
5. Fill out the PR description: what changed and why. Link any related issue.
6. Address review feedback. Once approved and checks pass, it gets merged.

## Testing your changes

**Tests run automatically** on every push and pull request. Make sure they pass before requesting review.

### Run tests locally:

```bash
# Quick check (no Docker needed)
cd tests/e2e
go test -v -short ./...

# Full test suite (requires Docker)
cd OnionEye_Code && docker compose up -d
sleep 30
cd ../tests/e2e && go test -v ./...
cd ../tests/integration && go test -v ./...
```

### ⚠️ **IMPORTANT: When to Update Tests**

**You MUST update tests** if you change any of the following:

1. **API Endpoints** (add/remove/rename routes)
   - Update: `tests/e2e/*_api_test.go`
   - Example: Changed `/api/v1/investigations` → Update `ingestion_api_test.go`

2. **Response Format** (change JSON structure)
   - Update: Assertions in `tests/e2e/*_test.go`
   - Example: Renamed `job_id` → `jobId` → Update all assertions checking this field

3. **Database Schema** (add/remove/rename columns)
   - Update: `tests/integration/database_test.go`
   - Example: Added `priority` column → Add test for new field

4. **Service Ports** (change port numbers in docker-compose.yml)
   - Update: `tests/e2e/config_test.go` default values
   - Example: Changed ingestion-api port 8080 → 8090 → Update `DefaultIngestionAPIURL`

5. **Environment Variables** (new required env vars)
   - Update: `tests/e2e/config_test.go` and CI workflow
   - Update: `.github/workflows/ci-main.yml` env section

6. **Authentication/Authorization** (add auth requirements)
   - Update: `tests/e2e/helpers_test.go` client configuration
   - Add auth headers/tokens to test requests

**Failure to update tests will cause CI to fail and block your PR from merging.**

### Quick Test Update Guide:

**Changed API endpoint?**
```bash
# Edit the test file for that service
vim tests/e2e/ingestion_api_test.go
# Update the URL path in the test
```

**Changed response format?**
```bash
# Find tests that check that response
rg "field_name" tests/e2e/
# Update assertions to match new format
```

**Changed database schema?**
```bash
# Update integration tests
vim tests/integration/database_test.go
# Add/modify tests for new schema
```

See `./testing/guide.md` for detailed examples and patterns.

## Commit signing

All commits must be signed. Unsigned commits will be rejected on protected branches.

- **SSH signing** (recommended if you already have an SSH key for GitHub):
  ```bash
  git config --global gpg.format ssh
  git config --global user.signingkey ~/.ssh/id_ed25519.pub
  git config --global commit.gpgsign true
  ```
- **GPG signing**: generate a key with `gpg --full-generate-key`, then:
  ```bash
  git config --global user.signingkey YOUR_KEY_ID
  git config --global commit.gpgsign true
  ```
- Add the corresponding public key to your GitHub account under **Settings → SSH and GPG keys**, using the "Signing Key" type for SSH.
- Commits show a "Verified" badge on GitHub once signing is working correctly.

## Opening an issue

Use the issue templates (bug report or feature request) rather than a blank issue — they ask for the details needed to triage quickly and auto-apply the relevant labels.

| Label | Meaning |
|---|---|
| `bug` | Something isn't working |
| `needs-triage` | Not yet reviewed by a maintainer |
| `enhancement` | New feature or request |
| `good first issue` | Good for newcomers |
| `help wanted` | Extra attention needed |
| `duplicate` | Already reported elsewhere |
| `wontfix` | Won't be worked on |
| `question` | Further info requested |

## Code review

This repo uses a `CODEOWNERS` file — opening a PR that touches a given path automatically requests review from the owner(s) of that area. A PR can't merge into a protected branch without their approval.

### Team Responsibilities

For code ownership and review assignments, see [TEAM_STRUCTURE.md](./TEAM_STRUCTURE.md):
- **AI/ML components** → Reviewed by AI Team
- **Backend/Frontend services** → Reviewed by Backend & Frontend Team  
- **Testing/CI/CD** → Reviewed by Live Testing & Debugging Team

## Releases

Releases follow semantic versioning (`vMAJOR.MINOR.PATCH`) and only happen from `main`, after a `develop → main` PR is merged:

```bash
git checkout main
git pull
git tag v1.1.0
git push origin v1.1.0
```

Pushing the tag triggers the release workflow automatically — it publishes a GitHub Release with notes grouped by label (Features / Bug Fixes / Docs / Maintenance) from the merged PRs since the last release. No manual changelog writing needed.

## Questions

If anything here is unclear:
- Check [TEAM_STRUCTURE.md](./TEAM_STRUCTURE.md) for the right team to contact
- Open a discussion or ask in your PR
- Don't guess and hope it's fine
