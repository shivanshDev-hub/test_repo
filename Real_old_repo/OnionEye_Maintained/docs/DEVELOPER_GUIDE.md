# OnionEye Developer Guide

> **📋 Team Structure:** See [TEAM_STRUCTURE.md](./TEAM_STRUCTURE.md) for current team assignments.

---

## 📚 Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)  
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Team Contacts](#team-contacts)

---

## Getting Started

### Prerequisites

```bash
go version      # 1.22+
docker --version
docker compose version
```

### Quick Setup

```bash
# 1. Clone and setup
git clone <repo-url>
cd OnionEye
git checkout develop

# 2. Start services
cd OnionEye_Code
docker compose up -d

# 3. Verify
curl http://localhost:8080/health
```

---

## Development Workflow

### Branching Strategy

- **main** - Production only
- **develop** - Integration (target for PRs)
- **feature/*** - New features
- **fix/*** - Bug fixes

### Making Changes

```bash
# Create branch from develop
git checkout develop && git pull
git checkout -b feature/your-feature

# Make changes and test
cd tests/e2e && go test -v ./...

# Commit with signing (required)
git commit -s -m "feat: description"

# Push and create PR to develop
git push origin feature/your-feature
```

### Commit Signing (Required)

```bash
# SSH signing
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/id_ed25519.pub  
git config --global commit.gpgsign true
```

Add your public key: GitHub → Settings → SSH and GPG keys

---

## Testing

### Quick Test (No Docker)

```bash
cd tests/e2e
go test -v -short ./...
```

### Full Test Suite

```bash
# Start services
cd OnionEye_Code && docker compose up -d && sleep 30

# Run tests
cd ../tests/e2e && go test -v ./...
cd ../integration && go test -v ./...

# Stop
cd ../../OnionEye_Code && docker compose down
```

### When to Update Tests

**Update tests if you change:**
1. API endpoints → `tests/e2e/*_api_test.go`
2. Response format → Test assertions
3. Database schema → `tests/integration/database_test.go`
4. Service ports → `tests/e2e/config_test.go`
5. Environment vars → Tests + CI workflow
6. Authentication → `tests/e2e/helpers_test.go`

**Find affected tests:**
```bash
rg "/api/v1/endpoint" tests/
rg "field_name" tests/e2e/
```

### CI/CD Pipeline

Tests run automatically on:
- Every push to main/develop/feature/fix branches
- Every PR (must pass to merge)
- Daily at 2 AM UTC
- Manual trigger in GitHub Actions

**Timeline:** ~12-20 minutes total

---

## Troubleshooting

### Services Not Starting

```bash
# Check status
docker compose ps

# View logs
docker compose logs -f ingestion-api

# Restart
docker compose restart

# Clean restart
docker compose down -v && docker compose up -d
```

### Port Conflicts

```bash
# Find what's using port
sudo lsof -i :8080

# Kill process
kill -9 <PID>
```

### Database Issues

```bash
# Connect
docker exec -it onioneye_code-postgres-1 psql -U onioneye -d onioneye

# Run migrations
docker exec -i onioneye_code-postgres-1 psql -U onioneye -d onioneye < OnionEye_Code/db/migrations/001_init_schema.sql

# Seed test data
docker exec -i onioneye_code-postgres-1 psql -U onioneye -d onioneye < OnionEye_Code/db/seed.sql
```

### Test Failures

```bash
# Connection refused → Check services running
docker compose ps

# Timeout → Increase timeout
go test -timeout 10m ./...

# Module errors → Update dependencies
cd tests/e2e && go mod tidy
```

---

## Team Contacts

See [TEAM_STRUCTURE.md](./TEAM_STRUCTURE.md) for current roster:

- **AI/ML components** → AI Team
- **Backend/Frontend** → Backend & Frontend Team  
- **Testing/CI/CD** → Live Testing & Debugging Team

---

## Quick Reference

### Service Ports

| Service | Port | URL |
|---------|------|-----|
| Ingestion API | 8080 | http://localhost:8080 |
| Storage API | 8081 | http://localhost:8081 |
| Correlation Engine | 8082 | http://localhost:8082 |
| Content Shield | 8083 | http://localhost:8083 |
| Dashboard API | 8084 | http://localhost:8084 |
| Frontend | 3000 | http://localhost:3000 |

### Useful Commands

```bash
# Start/stop
docker compose up -d
docker compose down

# Logs
docker compose logs -f <service>

# Tests
cd tests/e2e && go test -v ./...

# Health check
curl http://localhost:8080/health
```

---

**See also:** [README.md](../README.md) • [SECURITY.md](./SECURITY.md) • [testing README](./testing/README.md)
