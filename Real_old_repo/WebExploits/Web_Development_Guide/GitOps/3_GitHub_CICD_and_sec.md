# Github CI/Argon CD + Sec setting guide for dummies

# GitHub CI/CD Testing Guide: From Zero to Production-Safe

A practical, copy-pasteable guide to structuring GitHub Actions workflows, handling secrets properly, and writing all 8 major test types using a black-box-first approach — so the templates work regardless of what language your actual app is written in.

**Philosophy used throughout:** wherever possible, tests talk to your app over HTTP (black box) instead of importing its internals. This means the same test template works whether your app is written in Go, Python, Node, Java, or anything else. Go is used for lightweight, dependency-free black-box tools (smoke tests, load tests) because it compiles to a single static binary. Python is used where a rich ecosystem (HTTP libs, DB drivers, assertions) makes tests faster to write.

---

## Table of Contents

1. [Anatomy of a GitHub Actions Workflow](#1-anatomy-of-a-github-actions-workflow)
2. [Secrets: GitHub Secrets vs .env Files](#2-secrets-github-secrets-vs-env-files)
3. [The Testing Pyramid — What Runs When](#3-the-testing-pyramid--what-runs-when)
4. [Unit Tests](#4-unit-tests)
5. [Integration Tests](#5-integration-tests)
6. [End-to-End (E2E) Tests](#6-end-to-end-e2e-tests)
7. [Smoke Tests](#7-smoke-tests)
8. [Security Tests](#8-security-tests)
9. [Performance Tests](#9-performance-tests)
10. [Regression Tests](#10-regression-tests)
11. [Static Analysis](#11-static-analysis)
12. [Putting It All Together — Full Pipeline](#12-putting-it-all-together--full-pipeline)
13. [Using ArgoCD for CD](#13-using-argocd-for-cd)

---

## 1. Anatomy of a GitHub Actions Workflow

Every workflow lives in `.github/workflows/*.yml`. Here's every building block explained:

```yaml
name: CI Pipeline # Shown in the Actions tab

on: # What triggers this workflow
  push:
    branches: [main]
  pull_request:
    branches: [main]
  workflow_dispatch: # Adds a manual "Run workflow" button

env: # Global env vars, available to all jobs
  NODE_ENV: test

jobs:
  lint: # Job id (used in `needs:` references)
    runs-on: ubuntu-latest # Runner OS
    steps:
      - uses: actions/checkout@v4 # Pulls your repo code onto the runner
      - name: Run linter
        run: echo "linting..."

  unit-tests:
    needs: lint # Won't start until `lint` succeeds
    runs-on: ubuntu-latest
    strategy:
      matrix: # Runs this job once per combination
        python-version: ["3.11", "3.12"]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: ${{ matrix.python-version }}
      - run: pip install -r requirements.txt
      - run: pytest

  build:
    needs: unit-tests
    if: github.ref == 'refs/heads/main' # Conditional execution
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Upload build artifact
        uses: actions/upload-artifact@v4
        with:
          name: build-output
          path: dist/
```

**Key concepts:**

| Concept                                                | Purpose                                                                      |
| ------------------------------------------------------ | ---------------------------------------------------------------------------- |
| `on:`                                                  | Trigger — push, PR, schedule (`cron`), manual, or another workflow finishing |
| `jobs.<id>.needs`                                      | Defines job dependency order (parallel by default otherwise)                 |
| `runs-on`                                              | The VM image (`ubuntu-latest`, `windows-latest`, `macos-latest`)             |
| `strategy.matrix`                                      | Runs the same job across multiple versions/OSes in parallel                  |
| `services:`                                            | Spins up sidecar containers (DB, Redis) alongside your job                   |
| `if:`                                                  | Conditionally runs a job or step                                             |
| `secrets.*` / `vars.*`                                 | Injects secret/config values (see Section 2)                                 |
| `actions/upload-artifact`                              | Passes files between jobs or to the summary page                             |
| `environment:`                                         | Ties a job to a GitHub Environment (approval gates, env-scoped secrets)      |
| Reusable workflows (`uses: ./.github/workflows/x.yml`) | Avoid duplicating pipeline logic across repos                                |

---

## 2. Secrets: GitHub Secrets vs .env Files

### GitHub Secrets (for CI/CD — the real thing)

Set these under **Settings → Secrets and variables → Actions**. There are three scopes:

- **Repository secrets** — available to all workflows in the repo
- **Environment secrets** — scoped to a named environment (e.g. `production`), and can require manual approval before a job using them runs
- **Organization secrets** — shared across many repos

Usage in a workflow:

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production # ties this job to the "production" environment
    steps:
      - name: Deploy
        run: ./deploy.sh
        env:
          API_KEY: ${{ secrets.PROD_API_KEY }}
          DB_PASSWORD: ${{ secrets.PROD_DB_PASSWORD }}
```

**Rules of thumb:**

- GitHub automatically masks secret values in logs — but only exact-match strings, so don't `base64` or transform a secret before printing it.
- Never put secrets in `env:` at the top level of a public-facing workflow that runs on `pull_request` from forks — forked PRs can't read secrets by design, but `pull_request_target` can, so be careful which trigger you use.
- Prefer **OIDC federation** to cloud providers (AWS/GCP/Azure) over long-lived access keys stored as secrets — this removes the credential entirely:

```yaml
permissions:
  id-token: write
  contents: read
steps:
  - uses: aws-actions/configure-aws-credentials@v4
    with:
      role-to-assume: arn:aws:iam::123456789012:role/github-actions-deploy
      aws-region: us-east-1
```

### `.env` files (for local development only)

Never commit real secrets to a `.env` file. The pattern:

**`.env.example`** (committed to the repo, no real values):

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=devuser
DB_PASSWORD=devpass
API_KEY=replace_with_your_own_key
```

**`.gitignore`** (blocks the real file from ever being committed):

```
.env
.env.local
.env.*.local
*.secrets
```

Developers copy `.env.example` → `.env` and fill in real local values. Docker Compose or your app's config loader reads `.env` locally; in CI, the equivalent values come from `secrets.*` injected as environment variables — the `.env` file is never used in CI.

**Extra safety net:** run a secret-scanner (like `gitleaks`, covered in Section 8) on every push so an accidentally committed `.env` gets caught before merge, not after.

---

## 3. The Testing Pyramid — What Runs When

```
        /\
       /E2E\          few, slow, run pre-deploy or on a schedule
      /------\
     /Integr. \       moderate, run on every PR
    /----------\
   / Unit Tests \     many, fast, run on every commit
  /--------------\

Cross-cutting, run alongside the pyramid:
  Static analysis  → every push (seconds)
  Security scans   → every push / nightly
  Smoke tests      → immediately after every deploy
  Performance      → on-demand or nightly against staging
  Regression       → bundled into unit/integration/E2E runs via golden files
```

---

## 4. Unit Tests

Unit tests check one function/module in isolation, with no network, disk, or DB. This is the one category that has to be language-specific — you're testing internals, not a black box. Python is used here since it's the fastest to read/write for illustration; the pattern transfers directly to Go's `testing` package, Jest, JUnit, etc.

**`calculator.py`**

```python
def apply_discount(price: float, percent: float) -> float:
    if not 0 <= percent <= 100:
        raise ValueError("percent must be between 0 and 100")
    return round(price * (1 - percent / 100), 2)
```

**`test_calculator.py`**

```python
import pytest
from calculator import apply_discount

def test_apply_discount_normal():
    assert apply_discount(100, 10) == 90.0

def test_apply_discount_zero_percent():
    assert apply_discount(50, 0) == 50.0

def test_apply_discount_full_percent():
    assert apply_discount(50, 100) == 0.0

@pytest.mark.parametrize("price,percent,expected", [
    (200, 25, 150.0),
    (99.99, 10, 89.99),
])
def test_apply_discount_parametrized(price, percent, expected):
    assert apply_discount(price, percent) == expected

def test_apply_discount_invalid_percent_raises():
    with pytest.raises(ValueError):
        apply_discount(100, 150)
```

**Workflow job:**

```yaml
unit-tests:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-python@v5
      with:
        python-version: "3.12"
    - run: pip install pytest
    - run: pytest -v --junitxml=results/unit.xml
    - name: Upload results
      if: always()
      uses: actions/upload-artifact@v4
      with:
        name: unit-test-results
        path: results/unit.xml
```

---

## 5. Integration Tests

These check that your app correctly talks to a real dependency (a database, cache, or another service) — using GitHub Actions' `services:` to spin up a real container alongside your job.

**`test_integration_db.py`**

```python
import os
import psycopg2
import pytest

@pytest.fixture(scope="module")
def db_conn():
    conn = psycopg2.connect(
        host=os.getenv("DB_HOST", "localhost"),
        port=os.getenv("DB_PORT", 5432),
        dbname=os.getenv("DB_NAME", "testdb"),
        user=os.getenv("DB_USER", "testuser"),
        password=os.getenv("DB_PASSWORD", "testpass"),
    )
    yield conn
    conn.close()

def test_can_insert_and_read_user(db_conn):
    cur = db_conn.cursor()
    cur.execute("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)")
    cur.execute("INSERT INTO users (name) VALUES (%s) RETURNING id", ("Alice",))
    user_id = cur.fetchone()[0]
    db_conn.commit()

    cur.execute("SELECT name FROM users WHERE id = %s", (user_id,))
    assert cur.fetchone()[0] == "Alice"
```

**Workflow job:**

```yaml
integration-tests:
  runs-on: ubuntu-latest
  services:
    postgres:
      image: postgres:16
      env:
        POSTGRES_USER: testuser
        POSTGRES_PASSWORD: testpass
        POSTGRES_DB: testdb
      ports: ["5432:5432"]
      options: >-
        --health-cmd pg_isready
        --health-interval 10s
        --health-timeout 5s
        --health-retries 5
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-python@v5
      with: { python-version: "3.12" }
    - run: pip install pytest psycopg2-binary
    - run: pytest test_integration_db.py -v
      env:
        DB_HOST: localhost
        DB_USER: testuser
        DB_PASSWORD: testpass
        DB_NAME: testdb
```

---

## 6. End-to-End (E2E) Tests

E2E tests are fully black box: they hit a running instance of your app over HTTP and walk through a real user flow, exactly like a client would. This template works regardless of what your backend is written in.

**`test_e2e_signup_flow.py`**

```python
import os
import requests

BASE_URL = os.getenv("BASE_URL", "http://localhost:8080")

def test_full_signup_and_login_flow():
    # 1. Sign up
    signup_resp = requests.post(f"{BASE_URL}/api/signup", json={
        "email": "e2e_test_user@example.com",
        "password": "SuperSecret123!"
    })
    assert signup_resp.status_code == 201

    # 2. Log in
    login_resp = requests.post(f"{BASE_URL}/api/login", json={
        "email": "e2e_test_user@example.com",
        "password": "SuperSecret123!"
    })
    assert login_resp.status_code == 200
    token = login_resp.json()["token"]

    # 3. Access a protected resource
    me_resp = requests.get(f"{BASE_URL}/api/me", headers={
        "Authorization": f"Bearer {token}"
    })
    assert me_resp.status_code == 200
    assert me_resp.json()["email"] == "e2e_test_user@example.com"
```

**Workflow job:**

```yaml
e2e-tests:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Start app stack
      run: docker compose -f docker-compose.test.yml up -d --wait
    - uses: actions/setup-python@v5
      with: { python-version: "3.12" }
    - run: pip install requests pytest
    - run: pytest test_e2e_signup_flow.py -v
      env:
        BASE_URL: http://localhost:8080
    - name: Tear down
      if: always()
      run: docker compose -f docker-compose.test.yml down
```

> For UI-heavy apps (not just APIs), swap `requests` for **Playwright** (`pip install playwright pytest-playwright`) which drives a real browser but follows the exact same black-box philosophy.

---

## 7. Smoke Tests

Smoke tests run **immediately after a deploy** and answer one question: "is the thing even up?" They need to be fast and dependency-free — a perfect fit for a small Go binary, since Go compiles to a single static executable with zero runtime dependencies.

**`tools/smoketest/main.go`**

```go
package main

import (
 "fmt"
 "net/http"
 "os"
 "time"
)

func main() {
 url := getEnv("HEALTH_URL", "http://localhost:8080/health")
 client := http.Client{Timeout: 5 * time.Second}

 start := time.Now()
 resp, err := client.Get(url)
 elapsed := time.Since(start)

 if err != nil {
  fmt.Printf("SMOKE TEST FAILED: could not reach %s: %v\n", url, err)
  os.Exit(1)
 }
 defer resp.Body.Close()

 if resp.StatusCode != http.StatusOK {
  fmt.Printf("SMOKE TEST FAILED: got status %d from %s\n", resp.StatusCode, url)
  os.Exit(1)
 }

 if elapsed > 2*time.Second {
  fmt.Printf("SMOKE TEST FAILED: response took %s (limit 2s)\n", elapsed)
  os.Exit(1)
 }

 fmt.Printf("SMOKE TEST PASSED: %s responded 200 OK in %s\n", url, elapsed)
}

func getEnv(key, fallback string) string {
 if v := os.Getenv(key); v != "" {
  return v
 }
 return fallback
}
```

**Workflow job (with an automatic rollback trigger on failure):**

```yaml
smoke-test:
  needs: deploy-staging
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with: { go-version: "1.22" }
    - name: Run smoke test
      id: smoke
      run: go run ./tools/smoketest
      env:
        HEALTH_URL: https://staging.myapp.com/health

  rollback:
    needs: smoke-test
    if: failure()
    runs-on: ubuntu-latest
    steps:
      - name: Revert manifest to last known-good tag
        run: echo "kubectl rollout undo deployment/myapp"
```

---

## 8. Security Tests

For security, the practical approach is **using proven scanners**, not hand-writing checks. Four tools cover most of the surface area: secret leaks, dependency/container vulnerabilities, static code flaws, and a live black-box scan of the running app.

```yaml
security-tests:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
      with: { fetch-depth: 0 }

    - name: Scan for leaked secrets
      uses: gitleaks/gitleaks-action@v2

    - name: Scan filesystem & dependencies for vulnerabilities
      uses: aquasecurity/trivy-action@master
      with:
        scan-type: fs
        scan-ref: .
        severity: CRITICAL,HIGH
        exit-code: 1

    - name: CodeQL static analysis (SAST)
      uses: github/codeql-action/init@v3
      with:
        languages: python # or go, javascript, java, etc.
    - uses: github/codeql-action/analyze@v3

    - name: DAST against the running app
      uses: zaproxy/action-baseline@v0.12.0
      with:
        target: "https://staging.myapp.com"
```

**What each layer catches:**

| Tool      | Catches                                                               |
| --------- | --------------------------------------------------------------------- |
| Gitleaks  | API keys / credentials accidentally committed                         |
| Trivy     | Known CVEs in dependencies and base container images                  |
| CodeQL    | Insecure code patterns (SQL injection, XSS, etc.) via static analysis |
| OWASP ZAP | Real vulnerabilities on the live, running app (black box)             |

---

## 9. Performance Tests

Performance/load testing is inherently black box — you fire concurrent requests at an endpoint and measure latency and error rate. A tiny Go tool is genuinely practical here: no runtime, fast startup, easy concurrency via goroutines.

**`tools/loadtest/main.go`**

```go
package main

import (
 "fmt"
 "net/http"
 "os"
 "sort"
 "strconv"
 "sync"
 "time"
)

func main() {
 url := getEnv("TARGET_URL", "http://localhost:8080/api/health")
 requests, _ := strconv.Atoi(getEnv("TOTAL_REQUESTS", "500"))
 concurrency, _ := strconv.Atoi(getEnv("CONCURRENCY", "50"))
 maxP95Ms, _ := strconv.Atoi(getEnv("MAX_P95_MS", "300"))
 maxErrorRate, _ := strconv.ParseFloat(getEnv("MAX_ERROR_RATE", "0.01"), 64)

 var wg sync.WaitGroup
 var mu sync.Mutex
 latencies := make([]time.Duration, 0, requests)
 errors := 0

 sem := make(chan struct{}, concurrency)
 client := http.Client{Timeout: 10 * time.Second}

 for i := 0; i < requests; i++ {
  wg.Add(1)
  sem <- struct{}{}
  go func() {
   defer wg.Done()
   defer func() { <-sem }()

   start := time.Now()
   resp, err := client.Get(url)
   elapsed := time.Since(start)

   mu.Lock()
   defer mu.Unlock()
   if err != nil || resp == nil || resp.StatusCode >= 400 {
    errors++
   } else {
    latencies = append(latencies, elapsed)
   }
   if resp != nil {
    resp.Body.Close()
   }
  }()
 }
 wg.Wait()

 if len(latencies) == 0 {
  fmt.Println("FAILED: no successful responses")
  os.Exit(1)
 }

 sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
 p95 := latencies[int(float64(len(latencies))*0.95)]
 errorRate := float64(errors) / float64(requests)

 fmt.Printf("Requests: %d | Errors: %d (%.2f%%) | P95 latency: %s\n",
  requests, errors, errorRate*100, p95)

 if errorRate > maxErrorRate {
  fmt.Printf("FAILED: error rate %.2f%% exceeds max %.2f%%\n", errorRate*100, maxErrorRate*100)
  os.Exit(1)
 }
 if p95.Milliseconds() > int64(maxP95Ms) {
  fmt.Printf("FAILED: P95 latency %s exceeds max %dms\n", p95, maxP95Ms)
  os.Exit(1)
 }
 fmt.Println("PASSED")
}

func getEnv(key, fallback string) string {
 if v := os.Getenv(key); v != "" {
  return v
 }
 return fallback
}
```

**Workflow job:**

```yaml
performance-tests:
  needs: deploy-staging
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with: { go-version: "1.22" }
    - run: go run ./tools/loadtest
      env:
        TARGET_URL: https://staging.myapp.com/api/health
        TOTAL_REQUESTS: "1000"
        CONCURRENCY: "100"
        MAX_P95_MS: "300"
        MAX_ERROR_RATE: "0.01"
```

> For more advanced scenarios (ramping load, scripted user journeys), **k6** (JavaScript-based, from Grafana Labs) is the de facto industry tool and has an official `grafana/k6-action` — worth adopting once your load-test needs outgrow a single script.

---

## 10. Regression Tests

"Regression testing" isn't really a separate test _type_ — it's a strategy: re-run your existing unit/integration/E2E suite on every change so nothing that used to work silently breaks. The one regression-specific _technique_ worth its own template is **golden file (snapshot) testing** — freezing a known-good output and diffing against it forever after.

**`test_regression_report_output.py`**

```python
import json
import subprocess
from pathlib import Path

GOLDEN_FILE = Path("testdata/golden_report.json")

def generate_report():
    result = subprocess.run(
        ["python", "generate_report.py", "--input", "testdata/sample_input.csv"],
        capture_output=True, text=True, check=True
    )
    return json.loads(result.stdout)

def test_report_output_matches_golden_file():
    actual = generate_report()
    expected = json.loads(GOLDEN_FILE.read_text())
    assert actual == expected, (
        "Output changed! If this is expected, review the diff and "
        "update testdata/golden_report.json intentionally."
    )
```

**Workflow job:** regression tests typically just ride along inside your normal `unit-tests` or `integration-tests` job — no separate job is usually needed. If you want it isolated:

```yaml
regression-tests:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-python@v5
      with: { python-version: "3.12" }
    - run: pip install pytest
    - run: pytest test_regression_report_output.py -v
```

---

## 11. Static Analysis

Static analysis catches style and correctness issues without executing the code at all — cheapest and fastest check in the whole pipeline, so it should run first and fail fast.

```yaml
static-analysis:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4

    # Python
    - uses: actions/setup-python@v5
      with: { python-version: "3.12" }
    - run: pip install ruff mypy
    - run: ruff check .
    - run: mypy .

    # Go
    - uses: actions/setup-go@v5
      with: { go-version: "1.22" }
    - uses: golangci/golangci-lint-action@v6
      with:
        version: latest
```

---

## 12. Putting It All Together — Full Pipeline

```yaml
name: Full CI/CD Pipeline

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read
  packages: write

jobs:
  static-analysis:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with: { python-version: "3.12" }
      - run: pip install ruff
      - run: ruff check .

  unit-tests:
    needs: static-analysis
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with: { python-version: "3.12" }
      - run: pip install pytest
      - run: pytest -v

  integration-tests:
    needs: unit-tests
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: testuser
          POSTGRES_PASSWORD: testpass
          POSTGRES_DB: testdb
        ports: ["5432:5432"]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with: { python-version: "3.12" }
      - run: pip install pytest psycopg2-binary
      - run: pytest test_integration_db.py -v
        env:
          DB_HOST: localhost
          DB_USER: testuser
          DB_PASSWORD: testpass
          DB_NAME: testdb

  security-tests:
    needs: static-analysis
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: gitleaks/gitleaks-action@v2
      - uses: aquasecurity/trivy-action@master
        with:
          scan-type: fs
          scan-ref: .
          severity: CRITICAL,HIGH
          exit-code: 1

  build-and-push:
    needs: [integration-tests, security-tests]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ghcr.io/${{ github.repository }}:${{ github.sha }}

  deploy-staging:
    needs: build-and-push
    runs-on: ubuntu-latest
    environment: staging
    steps:
      - run: echo "update manifest repo / trigger ArgoCD sync to staging"

  smoke-test:
    needs: deploy-staging
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: "1.22" }
      - run: go run ./tools/smoketest
        env:
          HEALTH_URL: https://staging.myapp.com/health

  e2e-tests:
    needs: smoke-test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with: { python-version: "3.12" }
      - run: pip install requests pytest
      - run: pytest test_e2e_signup_flow.py -v
        env:
          BASE_URL: https://staging.myapp.com

  performance-tests:
    needs: e2e-tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: "1.22" }
      - run: go run ./tools/loadtest
        env:
          TARGET_URL: https://staging.myapp.com/api/health
          MAX_P95_MS: "300"
          MAX_ERROR_RATE: "0.01"

  deploy-production:
    needs: performance-tests
    runs-on: ubuntu-latest
    environment: production # requires manual approval if configured in repo settings
    steps:
      - run: echo "update manifest repo / trigger ArgoCD sync to production"

  production-smoke-test:
    needs: deploy-production
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: "1.22" }
      - run: go run ./tools/smoketest
        env:
          HEALTH_URL: https://myapp.com/health
```

**Why this order:** cheapest checks (lint) run first and fail fast; nothing expensive runs until static analysis and security pass; nothing deploys to production until staging has passed smoke, E2E, and performance checks; and production itself gets its own smoke test immediately after deploy so a bad release is caught before a human is paged at 3am.

---

## 13. Using ArgoCD for CD

Everything above ends at "artifact pushed to GHCR." This section covers what happens after that — ArgoCD taking the artifact and actually running it in your cluster, and everything you need to operate that safely without having to look anything else up.

### 13.1 Core Concepts

ArgoCD is a **GitOps continuous delivery** tool for Kubernetes. The mental model:

- **Git is the single source of truth.** Whatever is declared in a Git repo is what _should_ be running in the cluster.
- **ArgoCD is a reconciliation loop**, not a pipeline step. It continuously compares the live cluster state to the desired state in Git and corrects drift — it doesn't run once and exit like a CI job does.
- Your GitHub Actions pipeline **never talks to the cluster directly**. It only ever changes a file in Git (an image tag, a Helm value). ArgoCD does the actual `kubectl apply`.

Four objects you need to know:

| Object            | What it is                                                                                                                                                 |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Application**   | A CRD that tells ArgoCD: "watch this Git path, and make the cluster match it." One Application per deployable service (usually).                           |
| **AppProject**    | A grouping/permission boundary around Applications — which repos, which clusters, which resource kinds an Application in this project is allowed to touch. |
| **Sync**          | The act of reconciling: applying the Git-declared manifests to the cluster.                                                                                |
| **Health status** | ArgoCD's built-in understanding of whether a resource (Deployment, Pod, Ingress, custom CRD) is actually healthy, not just applied.                        |

Every Application has two independent states shown in the UI/CLI:

- **Sync status**: `Synced` / `OutOfSync` — does live state match Git?
- **Health status**: `Healthy` / `Progressing` / `Degraded` / `Missing` — is the workload actually running correctly?

Both need to be green. A Deployment can be `Synced` but `Degraded` (e.g. crash-looping pods) — ArgoCD applied what Git said, but the app itself is broken.

### 13.2 Installing ArgoCD

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for pods to be ready
kubectl wait --for=condition=Ready pods --all -n argocd --timeout=300s

# Get the initial admin password (change it immediately after first login)
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d

# Access the UI (port-forward for local access; use an Ingress for real use)
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

Install the CLI (used throughout this section):

```bash
# macOS
brew install argocd

# Linux
curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd

argocd login localhost:8080
```

In production, put ArgoCD's server behind a real Ingress with TLS, and immediately:

```bash
argocd account update-password
```

### 13.3 GitOps Repository Structure

The recommended split is **app repo** (source code, Dockerfile, CI workflow) separate from **manifest repo** (Kubernetes YAML / Helm values, watched by ArgoCD). This isn't mandatory, but it keeps "code changed" and "desired cluster state changed" as distinct, auditable events.

```
manifests-repo/
├── apps/
│   ├── my-service/
│   │   ├── base/
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   └── kustomization.yaml
│   │   └── overlays/
│   │       ├── staging/
│   │       │   ├── kustomization.yaml
│   │       │   └── values-staging.yaml   # if using Helm
│   │       └── production/
│   │           ├── kustomization.yaml
│   │           └── values-production.yaml
│   └── another-service/
│       └── ...
└── argocd/
    ├── my-service-staging.yaml       # Application manifests (see 13.4)
    ├── my-service-production.yaml
    └── app-of-apps.yaml              # see 13.6
```

Kustomize (`base` + `overlays`) or Helm (`values-<env>.yaml`) both work fine with ArgoCD — pick whichever your team already knows. Kustomize is shown here since it needs no templating language.

### 13.4 The Application Resource

This is the one YAML file that actually connects ArgoCD to your Git repo. Apply it once; ArgoCD takes over from there.

```yaml
# argocd/my-service-production.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: my-service-production
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io # ensures cluster resources are cleaned up if the Application is deleted
spec:
  project: default # or a named AppProject, see 13.12

  source:
    repoURL: https://github.com/your-org/manifests-repo.git
    targetRevision: main # branch, tag, or exact commit SHA
    path: apps/my-service/overlays/production

  destination:
    server: https://kubernetes.default.svc # target cluster (in-cluster here)
    namespace: my-service-prod

  syncPolicy:
    automated:
      prune: true # delete resources removed from Git
      selfHeal: true # revert any manual kubectl changes back to match Git
    syncOptions:
      - CreateNamespace=true
    retry:
      limit: 3
      backoff:
        duration: 10s
        factor: 2
        maxDuration: 3m
```

Apply it:

```bash
kubectl apply -f argocd/my-service-production.yaml
```

From this point on, **any commit that changes `apps/my-service/overlays/production/` gets auto-deployed to the cluster** — no human runs `kubectl apply` again, ever.

### 13.5 Sync Policies: Manual vs Automated

| Setting                         | Effect                                                                               |
| ------------------------------- | ------------------------------------------------------------------------------------ |
| No `syncPolicy.automated` block | Manual sync only — a human clicks "Sync" in the UI or runs `argocd app sync`         |
| `automated: {}`                 | Auto-syncs on every detected Git change                                              |
| `automated.prune: true`         | Also deletes cluster resources that were removed from Git                            |
| `automated.selfHeal: true`      | Reverts manual `kubectl edit`/`kubectl apply` drift back to match Git within seconds |

**Recommended split for a 3am-safe pipeline:**

- **Staging**: full `automated` with `prune: true` and `selfHeal: true` — fast feedback, nobody should hand-edit staging anyway.
- **Production**: either manual sync, or `automated` combined with a GitHub Environment **manual approval gate** upstream (Section 2) so the Git commit that triggers the prod sync only lands after a human approves the deploy. Most teams keep `selfHeal: true` even in production (so nobody can silently patch prod by hand) but leave `prune: false` initially until they trust the pipeline.

Manual sync via CLI:

```bash
argocd app sync my-service-production
argocd app wait my-service-production --health   # blocks until healthy or times out
```

### 13.6 App of Apps Pattern (Managing Many Services)

Once you have more than a handful of services, don't `kubectl apply` each Application manifest by hand — manage them as Applications-of-Applications: one root Application whose only job is to watch a folder of other Application manifests.

```yaml
# argocd/app-of-apps.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: root-apps
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/your-org/manifests-repo.git
    targetRevision: main
    path: argocd # the folder containing all the individual Application YAMLs
    directory:
      recurse: true
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

Apply this **one** manifest once:

```bash
kubectl apply -f argocd/app-of-apps.yaml
```

From then on, adding a new service to your platform is just: add a new file under `argocd/`, commit, push. The root Application notices the new file and creates the child Application automatically.

### 13.7 Connecting GitHub Actions to ArgoCD

There are two valid patterns. Pick one per-service; don't mix both for the same image.

**Approach A — CI writes the new tag to the manifest repo (push-based update, pull-based deploy).** This is the most common and most explicit approach — every deploy is a visible Git commit.

```yaml
# in your GitHub Actions workflow, after building & pushing the image to GHCR
- name: Update manifest repo with new image tag
  run: |
    git clone https://x-access-token:${{ secrets.MANIFEST_REPO_TOKEN }}@github.com/your-org/manifests-repo.git
    cd manifests-repo
    git config user.name "github-actions[bot]"
    git config user.email "github-actions[bot]@users.noreply.github.com"

    # Kustomize
    cd apps/my-service/overlays/production
    kustomize edit set image ghcr.io/your-org/my-service=ghcr.io/your-org/my-service:${{ github.sha }}

    git commit -am "deploy(my-service): ${{ github.sha }}"
    git push
```

`MANIFEST_REPO_TOKEN` is a fine-grained GitHub PAT with **contents: write** scoped only to the manifest repo, stored as a repository secret in the app repo.

**Approach B — ArgoCD Image Updater (fully pull-based, no CI write-back at all).** ArgoCD itself watches the container registry and updates the running image automatically, optionally still committing the change back to Git for audit purposes. Useful when you don't want CI to need write access to the manifest repo at all.

```bash
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj-labs/argocd-image-updater/stable/manifests/install.yaml
```

Add annotations to your Application:

```yaml
metadata:
  annotations:
    argocd-image-updater.argoproj.io/image-list: myservice=ghcr.io/your-org/my-service
    argocd-image-updater.argoproj.io/myservice.update-strategy: latest # or "digest", "semver"
    argocd-image-updater.argoproj.io/write-back-method: git # commits the bump back to Git for audit trail
```

With Approach B, your GitHub Actions workflow only needs to build, test, and push the image — it never touches the manifest repo. This is simpler to secure (no cross-repo write token needed) but slightly less explicit about exactly when a prod deploy happened, since the commit is made by the Image Updater bot rather than your CI run.

### 13.8 Handling Secrets Properly

**Never commit plaintext Kubernetes `Secret` manifests to the Git repo ArgoCD watches** — Git history is forever, and anyone with repo read access would see every secret ever deployed. Two standard solutions:

**Option 1 — Sealed Secrets** (encrypt before committing; only the cluster's private key can decrypt):

```bash
# One-time: install the controller
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/latest/download/controller.yaml

# Encrypt a secret locally, safe to commit the output
kubectl create secret generic db-credentials \
  --from-literal=password=supersecret \
  --dry-run=client -o yaml | kubeseal --format=yaml > sealed-db-credentials.yaml

git add sealed-db-credentials.yaml   # safe — this file is ciphertext
```

**Option 2 — External Secrets Operator** (Git stores only a _reference_; the real value lives in AWS Secrets Manager / GCP Secret Manager / Vault):

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: db-credentials
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: ClusterSecretStore
  target:
    name: db-credentials
  data:
    - secretKey: password
      remoteRef:
        key: prod/db/password
```

Either way, the file that ends up in the manifest repo is safe to commit — no plaintext credential ever touches Git.

### 13.9 Progressive Delivery: Canary & Blue-Green (Argo Rollouts)

Plain ArgoCD does an all-at-once rolling update by default (standard Kubernetes Deployment behavior). For canary or blue-green, install the companion project **Argo Rollouts** and swap your `Deployment` for a `Rollout`:

```bash
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
```

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: my-service
spec:
  replicas: 5
  strategy:
    canary:
      steps:
        - setWeight: 20 # 20% of traffic to the new version
        - pause: { duration: 5m }
        - analysis:
            templates:
              - templateName: success-rate-check # auto-abort if error rate spikes
        - setWeight: 50
        - pause: { duration: 5m }
        - setWeight: 100
  selector:
    matchLabels:
      app: my-service
  template:
    metadata:
      labels:
        app: my-service
    spec:
      containers:
        - name: my-service
          image: ghcr.io/your-org/my-service:latest
```

ArgoCD syncs a `Rollout` exactly like a `Deployment` — no change to your Application manifest is needed. The rollout controller handles the traffic shifting and can auto-abort a bad release mid-rollout based on metrics (Prometheus, Datadog) before it ever reaches 100% of users.

### 13.10 Sync Waves & Hooks (Ordering Dependent Resources)

When one resource must exist before another (e.g. a database migration Job before the app Deployment starts), use **sync waves**:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: db-migration
  annotations:
    argocd.argoproj.io/sync-wave: "-1" # negative = applied before wave 0 (the default)
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service
  annotations:
    argocd.argoproj.io/sync-wave: "0" # applied after the migration Job succeeds
```

For one-shot actions tied to the sync lifecycle itself (not just ordering), use **hooks**:

```yaml
metadata:
  annotations:
    argocd.argoproj.io/hook: PreSync # options: PreSync, Sync, PostSync, SyncFail
    argocd.argoproj.io/hook-delete-policy: HookSucceeded # clean up the Job once it succeeds
```

### 13.11 Notifications (Slack/Webhook on Sync Events)

```bash
kubectl apply -f https://raw.githubusercontent.com/argoproj-labs/argocd-notifications/release-1.0/manifests/install.yaml
```

```yaml
# argocd-notifications-cm ConfigMap (excerpt)
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-notifications-cm
  namespace: argocd
data:
  service.slack: |
    token: $slack-token
  template.app-deployed: |
    message: |
      🚀 {{.app.metadata.name}} deployed successfully — commit {{.app.status.sync.revision}}
  template.app-sync-failed: |
    message: |
      🔥 {{.app.metadata.name}} sync FAILED — check ArgoCD immediately
  trigger.on-deployed: |
    - when: app.status.operationState.phase in ['Succeeded'] and app.status.health.status == 'Healthy'
      send: [app-deployed]
  trigger.on-sync-failed: |
    - when: app.status.operationState.phase in ['Error', 'Failed']
      send: [app-sync-failed]
```

Subscribe an Application to a trigger by annotating it:

```yaml
metadata:
  annotations:
    notifications.argoproj.io/subscribe.on-sync-failed.slack: prod-deploys-alerts
```

This is what actually prevents the 3am page from being a surprise — a failed prod sync posts to Slack the second it happens, not when someone notices the app is down.

### 13.12 RBAC & Projects (Multi-Team Safety)

An `AppProject` limits what Applications inside it are allowed to do — which Git repos they may pull from, which clusters/namespaces they may deploy to, and which Kubernetes resource kinds they may create. This stops one team's Application from accidentally (or maliciously) modifying another team's namespace or a cluster-wide resource.

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: team-payments
  namespace: argocd
spec:
  sourceRepos:
    - https://github.com/your-org/manifests-repo.git
  destinations:
    - server: https://kubernetes.default.svc
      namespace: payments-* # only namespaces matching this pattern
  clusterResourceWhitelist: [] # no cluster-scoped resources allowed (CRDs, ClusterRoles, etc.)
  namespaceResourceWhitelist:
    - group: "*"
      kind: "*"
  roles:
    - name: deployer
      policies:
        - p, proj:team-payments:deployer, applications, sync, team-payments/*, allow
      groups:
        - your-org:payments-team # tied to your SSO group, see below
```

Combine with SSO (Dex or your existing OIDC provider) so ArgoCD RBAC maps to your real org's group membership rather than manually managed local accounts:

```yaml
# argocd-cm ConfigMap (excerpt)
data:
  url: https://argocd.yourcompany.com
  dex.config: |
    connectors:
      - type: github
        id: github
        name: GitHub
        config:
          clientID: $dex.github.clientId
          clientSecret: $dex.github.clientSecret
          orgs:
            - name: your-org
```

### 13.13 Rollbacks

Because Git is the source of truth, a rollback is just "make Git point at the previous good state" — there's no separate rollback mechanism to learn.

**Fastest — via ArgoCD directly (rolls back to a previous sync, without touching Git):**

```bash
argocd app history my-service-production
argocd app rollback my-service-production <HISTORY_ID>
```

Note: this is a temporary/manual override. On the next Git-triggered sync (or the next `selfHeal` cycle), ArgoCD will re-apply whatever is currently in Git — so treat this as a stop-the-bleeding measure, not a permanent fix.

**Correct — via Git (the durable fix):**

```bash
cd manifests-repo
git revert <bad-commit-sha>
git push
# ArgoCD auto-syncs the revert within seconds if automated sync is on,
# or: argocd app sync my-service-production
```

Because every deploy is a discrete Git commit with the exact image SHA baked in (Section 13.7), you always know exactly what "previous good state" means — no guessing which image tag was running before.

### 13.14 ArgoCD CLI Cheat Sheet

```bash
argocd login <server>                                   # authenticate
argocd app list                                         # list all Applications
argocd app get my-service-production                    # detailed status
argocd app diff my-service-production                   # show drift between Git and live cluster
argocd app sync my-service-production                   # manual sync
argocd app sync my-service-production --prune           # sync and delete removed resources
argocd app wait my-service-production --health          # block until healthy
argocd app rollback my-service-production <ID>          # roll back to a previous sync
argocd app history my-service-production                # list past syncs
argocd app logs my-service-production                   # tail logs of the app's pods
argocd app set my-service-production --sync-policy automated   # enable auto-sync
argocd app unset my-service-production --sync-policy    # disable auto-sync
argocd repo add https://github.com/your-org/manifests-repo.git --username x --password $TOKEN  # private repo access
argocd cluster list                                     # clusters ArgoCD can deploy to
argocd proj list                                         # list AppProjects
```

### 13.15 Full Worked Example (End-to-End)

Putting Sections 12 and 13 together, here's the complete lifecycle of one commit:

1. Dev pushes to `main` in the **app repo**.
2. GitHub Actions (Section 12 pipeline) runs static analysis → unit → integration → security → builds & pushes image to `ghcr.io/your-org/my-service:<sha>`.
3. On success, CI clones the **manifest repo** and bumps the image tag under `apps/my-service/overlays/staging/`, commits, pushes (Section 13.7, Approach A).
4. ArgoCD's `my-service-staging` Application (auto-sync on) detects the new commit within its poll interval (default 3 minutes, or instantly via a repo webhook — see below) and applies it.
5. `smoke-test` and `e2e-tests` jobs in GitHub Actions hit the now-updated staging URL directly (Section 7 and 6).
6. If those pass, a job updates `apps/my-service/overlays/production/` in the manifest repo — but this job runs inside a GitHub **Environment** with required reviewers, so it pauses for human approval first (Section 2).
7. Once approved, the commit lands, ArgoCD's `my-service-production` Application syncs it (manual or automated per Section 13.5), and Argo Rollouts (if configured) canaries the new version.
8. `production-smoke-test` confirms the live site is healthy; ArgoCD Notifications posts the result to Slack either way.
9. If anything is wrong: `argocd app rollback` immediately, then `git revert` in the manifest repo for the durable fix.

**Optional but recommended — webhook instead of polling:** add a GitHub webhook pointed at `https://argocd.yourcompany.com/api/webhook` (Settings → Webhooks in the manifest repo) so ArgoCD reacts to a push in seconds instead of waiting for its next poll cycle.

### 13.16 Troubleshooting Common Issues

| Symptom                                          | Likely cause / fix                                                                                                                                                      |
| ------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Application stuck `OutOfSync` after a sync       | Check `argocd app diff` — often a mutating admission webhook or a controller (e.g. HPA) is changing a field ArgoCD manages. Add that field to `spec.ignoreDifferences`. |
| `selfHeal` keeps reverting a manual hotfix       | That's it working correctly — the fix belongs in Git, not `kubectl edit`. Commit the fix instead.                                                                       |
| Sync succeeds but app is `Degraded`              | Sync status and health status are independent — check `kubectl describe pod` / `argocd app logs`; this is an application-level failure, not a GitOps failure.           |
| ArgoCD can't reach a private Git repo            | `argocd repo add <url> --username x --password <PAT>`, or configure an SSH deploy key via `repo-creds`.                                                                 |
| Image Updater not picking up new tags            | Confirm the annotation's image alias matches exactly, and check `kubectl logs -n argocd deploy/argocd-image-updater`.                                                   |
| Two Applications fighting over the same resource | Each resource should be owned by exactly one Application — split overlapping paths, or use `AppProject` resource whitelists (Section 13.12) to enforce it.              |
| Rollback via CLI didn't "stick"                  | Expected — `argocd app rollback` is temporary. Revert in Git for a permanent fix (Section 13.13).                                                                       |

## GitHub Production Safety Guide: Branch Protection, CODEOWNERS, Environments, Secrets, Deployment Strategy, Rollback & Dependency Management

This is a companion to the CI/CD testing guide and the ArgoCD guide — those cover _what runs_; this one covers _what stops a bad change from ever reaching production in the first place_. Same example layout as before: an app repo (`my-service`) and a separate GitOps `manifests-repo`.

Every section gives you: (1) the exact click-path in the GitHub UI, and (2) the scriptable CLI/API equivalent, so you can set this up once by hand or bake it into a repo-provisioning script for every new service.

---

### Table of Contents

1. [Branch Protection Rules](#1-branch-protection-rules)
2. [CODEOWNERS File](#2-codeowners-file)
3. [Environment Protection Rules](#3-environment-protection-rules)
4. [Secret Hygiene](#4-secret-hygiene)
5. [Deployment Strategy](#5-deployment-strategy)
6. [Rollback Plan — Built In, Not Bolted On](#6-rollback-plan--built-in-not-bolted-on)
7. [Dependency & Vulnerability Management](#7-dependency--vulnerability-management)
8. [Putting It All Together — Repo Setup Checklist](#8-putting-it-all-together--repo-setup-checklist)

---

### 1. Branch Protection Rules

**Goal:** nobody, including admins, can push straight to `main` or merge without CI passing and a review.

#### Exact steps (UI)

1. Go to your repo → **Settings → Branches**.
2. Under "Branch protection rules," click **Add rule**.
3. Branch name pattern: `main`
4. Check these boxes:
   - ☑ **Require a pull request before merging**
     - ☑ Require approvals → set to `1` (or `2` for higher-risk repos)
     - ☑ Dismiss stale pull request approvals when new commits are pushed
     - ☑ Require review from Code Owners (this activates the CODEOWNERS file — Section 2)
   - ☑ **Require status checks to pass before merging**
     - ☑ Require branches to be up to date before merging
     - In the search box, add every required CI job name from your workflow: `static-analysis`, `unit-tests`, `integration-tests`, `security-tests` (these must have run at least once on a PR before they appear in the list)
   - ☑ **Require conversation resolution before merging**
   - ☑ **Require signed commits** (optional — enable if your org mandates commit signing)
   - ☑ **Require linear history** (blocks merge commits, keeps `git log` readable)
   - ☑ **Do not allow bypassing the above settings** (this is what makes it apply to admins too — without it, repo admins can ignore every rule above)
   - ☐ Leave "Allow force pushes" and "Allow deletions" **unchecked**
5. Click **Create**.
6. Repeat for any other long-lived branch (e.g. `release/*`) if you use one.

#### Scriptable equivalent (API / `gh` CLI)

Save this as `branch-protection.json`:

```json
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "static-analysis",
      "unit-tests",
      "integration-tests",
      "security-tests"
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": true,
    "required_approving_review_count": 1
  },
  "restrictions": null,
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true
}
```

Apply it:

```bash
gh api repos/your-org/my-service/branches/main/protection \
  -X PUT \
  --input branch-protection.json
```

Because this is just a JSON file + one command, you can commit it to a `platform-scripts` repo and run it against every new service repo so protection settings are never "forgotten" for a new project.

> **Newer alternative:** GitHub's **Rulesets** (Settings → Rules → Rulesets) do the same job with more flexibility (target multiple branches by pattern, layer multiple rulesets, org-wide rulesets applied to every repo automatically). Classic branch protection above still works fine and is simpler to reason about for a single repo — move to Rulesets once you're managing protection across many repos at once.

---

### 2. CODEOWNERS File

**Goal:** the right humans are automatically requested as reviewers on the paths they actually own — nobody has to remember to tag them.

### Exact steps

1. Create the file at `.github/CODEOWNERS` (also works at repo root or in `docs/`, but `.github/` is the convention).
2. Populate it — most specific path wins, later rules override earlier ones for the same path:

```
# .github/CODEOWNERS

# Default owner for everything not matched below
*                              @your-org/backend-team

# CI/CD and infra-as-code need platform eyes on every change
/.github/workflows/           @your-org/platform-team
/Dockerfile                   @your-org/platform-team
/infra/                       @your-org/platform-team

# Anything touching auth is security-sensitive
/src/auth/                    @your-org/backend-team @your-org/security-team

# Docs can be approved by anyone on the team, no gate needed
/docs/                        @your-org/backend-team
```

For the **manifests-repo** (from the ArgoCD guide), use a stricter file since this repo directly controls what's running in the cluster:

```
# manifests-repo/.github/CODEOWNERS

*                                    @your-org/platform-team

# Production overlays require both platform AND SRE sign-off
apps/*/overlays/production/         @your-org/platform-team @your-org/sre-team
argocd/                             @your-org/platform-team
```

1. Commit and push the file.
2. Confirm branch protection has **"Require review from Code Owners"** checked (Section 1) — without that box checked, the CODEOWNERS file only auto-_requests_ reviewers, it doesn't _require_ their approval to merge.

### Verifying it worked

Open any existing PR that touches a path in the file — GitHub should show the matching team/user automatically added under "Reviewers" within a minute of the file landing on the default branch.

---

### 3. Environment Protection Rules

**Goal:** deploying to production requires a human to click "approve," and only certain branches/tags are even eligible to deploy there.

### Exact steps (UI)

1. Go to repo → **Settings → Environments → New environment**.
2. Name it `production` (also create `staging` the same way, typically with no protection rules — staging should deploy freely).
3. Under **Deployment protection rules**:
   - ☑ **Required reviewers** → add specific people or a team (e.g. `@your-org/sre-team`). Any workflow job with `environment: production` now pauses and waits for one of these reviewers to click **Approve** in the Actions tab.
   - ☑ **Wait timer** → optional; e.g. `10` minutes as a mandatory cooling-off period even after approval.
4. Under **Deployment branches and tags** → change from "All branches" to **"Selected branches and tags"** → add `main` and, if you tag releases, a pattern like `v*`.
5. Under **Environment secrets**, add production-only secrets here (Section 4) instead of at the repo level.
6. Click **Save protection rules**.

Reference it in your workflow exactly as shown in the earlier CI guide:

```yaml
deploy-production:
  needs: performance-tests
  runs-on: ubuntu-latest
  environment: production # this line alone activates all rules configured above
  steps:
    - run: echo "update manifest repo / trigger ArgoCD sync to production"
```

When this job reaches the runner queue, it will sit in a **"Waiting"** state in the Actions tab until a required reviewer approves it — this is the human gate mentioned throughout the previous guides.

#### Scriptable equivalent (API)

```bash
# Create/update the environment with a wait timer and required reviewers
cat > env-production.json <<'EOF'
{
  "wait_timer": 10,
  "reviewers": [
    { "type": "Team", "id": 123456 }
  ],
  "deployment_branch_policy": {
    "protected_branches": false,
    "custom_branch_policies": true
  }
}
EOF

gh api repos/your-org/my-service/environments/production -X PUT --input env-production.json

# Then explicitly allow "main" and tags matching "v*" to deploy to it
gh api repos/your-org/my-service/environments/production/deployment-branch-policies \
  -X POST -f name='main'

gh api repos/your-org/my-service/environments/production/deployment-branch-policies \
  -X POST -f name='v*'
```

To find a team's numeric `id` for the reviewers field: `gh api orgs/your-org/teams/sre-team --jq .id`.

---

### 4. Secret Hygiene

**Goal:** secrets are scoped as narrowly as possible, rotated on a schedule, and an accidental commit is caught before it ever reaches Git history.

### Exact steps

**a. Store secrets at the narrowest scope that works**

- Production credentials → **Environment secrets** on the `production` environment (Section 3), not repo-level. A PR from a branch can never read them even if it tries, because environment secrets are only exposed to jobs that reference `environment: production` _and_ have passed approval.
- Shared, non-sensitive config (feature flags, non-secret URLs) → **Environment variables**, not secrets — keeps the Actions log readable since variables aren't masked.

**b. Turn on GitHub's built-in secret scanning**

1. Repo → **Settings → Code security and analysis**.
2. Enable **Secret scanning**.
3. Enable **Push protection** — this actively _blocks_ a `git push` that contains a recognizable secret pattern (AWS keys, GitHub tokens, Stripe keys, etc.), before it ever lands in history.

```bash
# Same thing via API, useful for provisioning many repos at once
gh api repos/your-org/my-service -X PATCH \
  -f security_and_analysis[secret_scanning][status]=enabled \
  -f security_and_analysis[secret_scanning_push_protection][status]=enabled
```

**c. Add a local pre-commit scan too** (catches it before it even reaches GitHub)

```bash
pip install pre-commit
```

`.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.18.4
    hooks:
      - id: gitleaks
```

```bash
pre-commit install     # one-time, per developer clone
```

**d. Prefer short-lived credentials over long-lived secrets wherever possible**

For cloud deploys, use OIDC federation instead of a stored access key (no secret to leak or rotate at all):

```yaml
permissions:
  id-token: write
  contents: read
steps:
  - uses: aws-actions/configure-aws-credentials@v4
    with:
      role-to-assume: arn:aws:iam::123456789012:role/github-actions-deploy
      aws-region: us-east-1
```

**e. For secrets that must exist, scope and expire tokens tightly**

When creating a PAT for CI use (e.g. the `MANIFEST_REPO_TOKEN` from the ArgoCD guide):

1. Your **user (or bot account) → Settings → Developer settings → Personal access tokens → Fine-grained tokens → Generate new token**.
2. **Repository access** → "Only select repositories" → pick just `manifests-repo`.
3. **Permissions** → grant only `Contents: Read and write` — nothing else.
4. **Expiration** → set a real expiry (90 days), never "No expiration."
5. Store the resulting token as a repository/environment secret; put the renewal date in a team calendar reminder.

**f. Rotation schedule (write this down somewhere your team actually looks)**

| Secret type                                                   | Rotate every                                                |
| ------------------------------------------------------------- | ----------------------------------------------------------- |
| Cloud provider long-lived keys (avoid if possible — use OIDC) | 90 days                                                     |
| Fine-grained PATs                                             | At expiry (set to ≤90 days)                                 |
| Database passwords                                            | 90–180 days, or immediately after any team member offboards |
| Webhook signing secrets                                       | 180 days                                                    |
| Any secret after a suspected leak                             | Immediately                                                 |

**g. Audit what secrets exist, periodically**

```bash
gh secret list --repo your-org/my-service
gh secret list --repo your-org/my-service --env production
```

---

### 5. Deployment Strategy

**Goal:** pick the deployment style that matches how much blast-radius you can tolerate, and configure it once rather than deciding ad hoc on each release.

#### Exact steps — decide, then implement

**Step 1: Pick a strategy per environment using this table:**

| Strategy                                | How it works                                                                                       | Blast radius if broken                                        | Setup effort                                                                  |
| --------------------------------------- | -------------------------------------------------------------------------------------------------- | ------------------------------------------------------------- | ----------------------------------------------------------------------------- |
| **Rolling update** (Kubernetes default) | Old pods replaced by new ones gradually, all-or-nothing per pod                                    | Medium — bad version reaches 100% of traffic within minutes   | None — this is the default `Deployment` behavior                              |
| **Blue-Green**                          | Full second environment stood up, traffic switched all at once                                     | Low if caught before switch, high if not                      | Medium — needs double the resources briefly                                   |
| **Canary**                              | New version gets a small % of real traffic, ramps up gradually with automated checks between steps | Lowest — a bad release is caught at 10–20% of users, not 100% | Medium — needs Argo Rollouts (see ArgoCD guide, Section 13.9)                 |
| **Feature flags**                       | Code is deployed everywhere, but the _feature_ is toggled per-user/segment independent of deploys  | Lowest for the toggled feature specifically                   | Requires a flag system (LaunchDarkly, Unleash, or a simple config-based flag) |

**Step 2: Recommended default for most teams**

- **Staging** → rolling update, fully automated, no gate (Section 3) — fast feedback, low stakes.
- **Production** → canary via Argo Rollouts, behind the environment approval gate from Section 3. This combination means: a human approves the release, then the rollout itself is additionally staged and can auto-abort on bad metrics before reaching all users.
- **Any user-facing risky feature** (new pricing logic, new onboarding flow) → ship behind a feature flag _in addition to_ the above, so it can be turned off in seconds without a redeploy at all.

**Step 3: Implement it** — this is exactly the `Rollout` resource and canary `steps:` block from the ArgoCD guide (Section 13.9); nothing new to configure here beyond what's already documented there. The one new piece worth adding is an automated abort condition tied to your smoke/E2E results:

```yaml
strategy:
  canary:
    steps:
      - setWeight: 20
      - pause: { duration: 5m }
      - analysis:
          templates:
            - templateName: success-rate-check # auto-abort if error rate spikes mid-rollout
      - setWeight: 100
```

**Step 4: Write down the decision.** Add a short `DEPLOYMENT_STRATEGY.md` to the repo root stating which strategy is used for which environment, so a new team member doesn't have to reverse-engineer it from YAML at 3am.

---

### 6. Rollback Plan — Built In, Not Bolted On

**Goal:** rolling back is a single command that anyone on the team can run without having to first figure out what the previous good version even was.

#### Exact steps

**a. Tag every production release automatically, as part of the pipeline**

```yaml
# Add as a job in the CI workflow, gated behind deploy-production succeeding
tag-release:
  needs: deploy-production
  runs-on: ubuntu-latest
  permissions:
    contents: write
  steps:
    - uses: actions/checkout@v4
      with: { fetch-depth: 0 }
    - name: Create and push tag
      run: |
        VERSION="v$(date +%Y.%m.%d)-${GITHUB_SHA::7}"
        echo "VERSION=$VERSION" >> "$GITHUB_ENV"
        git tag -a "$VERSION" -m "Production release $VERSION"
        git push origin "$VERSION"
    - name: Create GitHub Release
      uses: softprops/action-gh-release@v2
      with:
        tag_name: ${{ env.VERSION }}
        generate_release_notes: true
```

Now every production deploy has a permanent, human-readable marker (`v2026.09.07-a1b2c3d`) tied to an exact commit and an exact container image tag — "what was running before" is never a guessing game.

**b. Never let the image registry only keep `latest`**

- Always push both an immutable tag (git SHA) and, optionally, `latest` — but deploy manifests should always reference the SHA, never `latest` (this was already the pattern used in the earlier CI guide's `build-and-push` job).
- Set a **package retention policy** so old SHAs aren't garbage-collected before you might need to roll back to them: repo → **Settings → Packages** → configure retention (e.g. keep last 30 versions, or keep anything referenced by an open Git tag).

**c. Write the actual rollback command down — better yet, script it**

`scripts/rollback.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

APP_NAME="${1:?Usage: rollback.sh <argocd-app-name>}"

echo "Current sync history for $APP_NAME:"
argocd app history "$APP_NAME"

read -rp "Enter the HISTORY ID to roll back to: " HISTORY_ID

argocd app rollback "$APP_NAME" "$HISTORY_ID"
argocd app wait "$APP_NAME" --health --timeout 120

echo "Rolled back. This is a TEMPORARY fix — now revert the bad commit in the manifests repo:"
echo "  cd manifests-repo && git revert <bad-commit-sha> && git push"
```

```bash
chmod +x scripts/rollback.sh
./scripts/rollback.sh my-service-production
```

This mirrors exactly the two-step rollback process from the ArgoCD guide (Section 13.13: immediate `argocd app rollback` to stop the bleeding, then a durable `git revert`) — but as a script instead of something someone has to remember at 3am.

**d. Add a one-page runbook**

Create `ROLLBACK.md` at the repo root with exactly three things: (1) the command above, (2) where to check if the rollback worked (the smoke-test URL), (3) who to notify (the Slack channel wired up in the ArgoCD Notifications section of the previous guide). This file existing and being findable in under 10 seconds is the actual point — the best rollback plan is useless if it's tribal knowledge in one person's head.

---

### 7. Dependency & Vulnerability Management

**Goal:** vulnerable and outdated dependencies get flagged and patched automatically, with low-risk updates merging themselves.

#### Exact steps

**a. Enable the underlying GitHub security features**

1. Repo → **Settings → Code security and analysis**.
2. Enable **Dependency graph** (usually on by default for public repos).
3. Enable **Dependabot alerts**.
4. Enable **Dependabot security updates** (auto-opens a PR the moment a vulnerability is found in a dependency you use, even before your next scheduled scan).

**b. Configure Dependabot version updates** (routine "keep things current," not just security)

Create `.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    labels: ["dependencies"]

  - package-ecosystem: "pip"
    directory: "/"
    schedule:
      interval: "weekly"

  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"

  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "weekly"

  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
```

Only include the `package-ecosystem` blocks relevant to your repo (e.g. drop `gomod` if there's no Go code). Commit this file — Dependabot picks it up automatically, no toggle needed elsewhere.

**c. Auto-merge low-risk Dependabot PRs** so patch/minor bumps don't pile up needing manual review

```yaml
# .github/workflows/dependabot-automerge.yml
name: Dependabot auto-merge

on: pull_request

permissions:
  pull-requests: write
  contents: write

jobs:
  auto-merge:
    if: github.actor == 'dependabot[bot]'
    runs-on: ubuntu-latest
    steps:
      - name: Fetch Dependabot metadata
        id: metadata
        uses: dependabot/fetch-metadata@v2
        with:
          github-token: "${{ secrets.GITHUB_TOKEN }}"

      - name: Auto-merge patch and minor updates
        if: |
          steps.metadata.outputs.update-type == 'version-update:semver-patch' ||
          steps.metadata.outputs.update-type == 'version-update:semver-minor'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

This still requires your normal branch protection checks (Section 1) to pass first — `--auto` just tells GitHub to merge automatically the moment those checks go green, it doesn't skip them. **Major** version bumps are deliberately excluded here so they still get a human look.

**d. Block a PR that introduces a newly vulnerable dependency, even outside Dependabot's own PRs**

```yaml
# add as a job in your main CI workflow
dependency-review:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/dependency-review-action@v4
      with:
        fail-on-severity: high
```

This runs on every PR (not just Dependabot's) and fails the check if a newly added or bumped dependency has a known high/critical vulnerability — catching the case where a human manually adds a bad package.

**e. Alternative worth knowing about: Renovate**

If you need more control than Dependabot offers (grouping multiple packages into one PR, custom scheduling per package, monorepo support), **Renovate** (via the Renovate GitHub App) is a drop-in alternative with a `renovate.json` config file instead of `dependabot.yml`. Most teams don't need to switch — Dependabot covers the common case well — but it's the standard escape hatch once Dependabot's config options feel limiting.

---

### 8. Putting It All Together — Repo Setup Checklist

Run through this once per new production service repo:

- [ ] Branch protection rule on `main` applied (Section 1) — either via UI or `gh api ... --input branch-protection.json`
- [ ] `.github/CODEOWNERS` committed, and "Require review from Code Owners" checked in branch protection (Section 2)
- [ ] `staging` and `production` Environments created; `production` has required reviewers + restricted deploy branches (Section 3)
- [ ] Secret scanning + push protection enabled; production secrets live in the `production` environment, not repo-level (Section 4)
- [ ] Pre-commit gitleaks hook documented in the repo's `CONTRIBUTING.md` (Section 4)
- [ ] Deployment strategy chosen and written down in `DEPLOYMENT_STRATEGY.md` (Section 5)
- [ ] Release tagging job wired into the pipeline; `scripts/rollback.sh` and `ROLLBACK.md` committed (Section 6)
- [ ] `.github/dependabot.yml` committed; Dependabot alerts + security updates enabled; auto-merge workflow and dependency-review job added (Section 7)

Once every box is checked, the failure modes that actually cause a 3am page — an unreviewed change, a bad prod deploy nobody approved, a leaked credential, an unpatched CVE, or a rollback nobody remembers how to run — each have a specific, already-built answer instead of being improvised in the moment.
