# OnionEye Test Suite

> **📋 Team Structure:** See [TEAM_STRUCTURE.md](../TEAM_STRUCTURE.md) for current team assignments and responsibilities.

## 🚨 **Developer Notice** 🚨

**If you're changing the API, database, or service configuration, you MUST update the corresponding tests!**

See **[Maintenance Guide](./maintenance.md)** for detailed instructions on when and how to update tests.

**Contact Information:**
- **AI/ML components** → Contact AI Team
- **Backend/Frontend code** → Contact Backend & Frontend Team
- **Testing issues** → Contact Live Testing & Debugging Team

---

## 📚 Documentation Index

| Guide | Purpose | When to Read |
|-------|---------|-------------|
| **[Maintenance Guide](./maintenance.md)** | When to update tests | ⚠️ Before making any API/schema changes |
| **[Quick Start](./quick-start.md)** | How to run tests | First time running tests |
| **[Testing Guide](./guide.md)** | Testing patterns & best practices | Writing new tests |
| **[CONTRIBUTING.md](../CONTRIBUTING.md)** | Contribution workflow | Before submitting a PR |

---

## ⚡ Quick Commands

### Run Tests Locally (Without Docker)
```bash
cd tests/e2e
go test -v -short ./...
```

### Run Tests Locally (With Docker)
```bash
# Terminal 1: Start services
cd OnionEye_Code
docker compose up

# Terminal 2: Run tests
cd tests/e2e
go test -v ./...
```

### Run Specific Test
```bash
cd tests/e2e
go test -v -run TestIngestionAPIHealth
```

---

## 🧪 Test Organization

```
tests/
├── e2e/                          # End-to-end tests
│   ├── config_test.go            # Test configuration
│   ├── helpers_test.go           # Test utilities
│   ├── ingestion_api_test.go     # Ingestion API tests
│   ├── dashboard_api_test.go     # Dashboard API tests
│   ├── workflow_test.go          # Multi-service workflows
│   └── integration_test.go       # Full stack integration
│
├── integration/                  # Integration tests
│   └── database_test.go          # Database tests
│
└── maintenance.md               # 👈 READ THIS BEFORE CHANGING CODE
```

---

## 🤖 Automated Testing

Tests run automatically on:

1. **Every push** to `main`, `develop`, `feature/*`, `fix/*`
2. **Every pull request** (must pass before merge)
3. **Daily at 2 AM UTC** (health check)
4. **Manual trigger** via GitHub Actions

---

## ⚠️ Common Mistakes

### ❌ Don't Do This
```bash
# Making API changes without updating tests
git commit -m "changed endpoint from /investigations to /cases"
git push
# ❌ Tests will fail in CI and block your PR!
```

### ✅ Do This Instead
```bash
# Update tests in the same commit
vim OnionEye_Code/services/ingestion-api/handler.go  # Change API
vim tests/e2e/ingestion_api_test.go                  # Update test
git add OnionEye_Code/ tests/
git commit -s -m "feat: rename investigations endpoint to cases"
git push
# ✅ Tests pass and PR can merge
```

---

## 📖 Test Coverage

### E2E Tests Cover:
- ✅ Ingestion API (health, create, validate, guard)
- ✅ Dashboard API (stats, graphs, jobs)
- ✅ Full workflows (investigation → correlation → dashboard)
- ✅ Resilience (retries, timeouts, errors)

### Integration Tests Cover:
- ✅ PostgreSQL connectivity
- ✅ Redis operations
- ✅ Database schema validation
- ✅ Actor CRUD operations

---

## Environment Variables

```bash
# E2E Tests
export INGESTION_API_URL="http://localhost:8080"
export DASHBOARD_API_URL="http://localhost:3000"
export STORAGE_API_URL="http://localhost:8081"

# Integration Tests
export POSTGRES_HOST="localhost"
export POSTGRES_USER="onioneye"
export POSTGRES_PASSWORD="testpass123"
export REDIS_HOST="localhost"
```

---

## 🆘 Getting Help

**Tests failing?**
1. Read the error message in CI logs
2. Run tests locally to reproduce: `cd tests/e2e && go test -v ./...`
3. Check [Maintenance Guide](./maintenance.md)
4. Ask in your PR comments

**Don't know which tests to update?**
```bash
# Find tests using a specific endpoint
rg "/api/v1/investigations" tests/

# Find tests checking a specific field
rg "job_id" tests/
```

---

## Troubleshooting

**Tests fail with "connection refused"**
- Ensure all services are running: `docker compose ps`
- Wait for services to be healthy: `docker compose logs`

**Tests timeout**
- Increase timeout: `go test -timeout 10m ./...`
- Check service logs for errors

**Database errors**
- Run migrations: `docker exec -i onioneye_code-postgres-1 psql -U onioneye -d onioneye < db/schema.sql`

---

## 🎯 Next Steps

1. **New here?** → Read [Quick Start](./quick-start.md)
2. **Making changes?** → Read [Maintenance Guide](./maintenance.md) ⚠️
3. **Writing tests?** → Read [Testing Guide](./guide.md)
4. **Opening a PR?** → Read [CONTRIBUTING.md](../CONTRIBUTING.md)

---

**Remember: If you change the API, you MUST change the tests. No exceptions!** 🧪

