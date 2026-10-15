# 🚀 Quick Start: Running OnionEye Tests

> **📋 Team Structure:** See [TEAM_STRUCTURE.md](../TEAM_STRUCTURE.md) for current team assignments and contact information.

## **Prerequisites**

```bash
# Check if you have the required tools
go version      # Should be 1.22+
docker --version
docker compose version
```

---

## **🎯 OPTION 1: Quick Test (No Docker Required)**

Test that everything compiles without running services:

```bash
cd tests/e2e
go mod download
go test -v -short ./...
```

**Result:** Tests will skip gracefully (no failures) since services aren't running.

---

## **🐳 OPTION 2: Full Test (With Docker)**

### Step 1: Start Services
```bash
cd OnionEye_Code
docker compose up -d
```

**Wait 30 seconds for services to fully start.**

### Step 2: Run E2E Tests
```bash
cd tests/e2e
go test -v ./...
```

### Step 3: Run Integration Tests
```bash
cd tests/integration
go test -v ./...
```

### Step 4: Stop Services
```bash
cd OnionEye_Code
docker compose down
```

---

## **⚡ One-Command Full Test**

Copy-paste this to run everything:

```bash
cd OnionEye_Code && \
docker compose up -d && \
echo "Waiting for services..." && sleep 30 && \
cd ../tests/e2e && go test -v ./... && \
cd ../integration && go test -v ./... && \
cd ../../OnionEye_Code && docker compose down
```
go test -v ./...
```

### Step 3: Run Integration Tests
```bash
cd tests/integration
go test -v ./...
```

### Step 4: Stop Services
```bash
cd OnionEye_Code
docker compose down
```

---

## **⚡ One-Command Full Test**

Copy-paste this to run everything:

```bash
cd OnionEye_Code && \
docker compose up -d && \
echo "Waiting for services..." && sleep 30 && \
cd ../tests/e2e && go test -v ./... && \
cd ../integration && go test -v ./... && \
cd ../../OnionEye_Code && docker compose down
```

---

## **🔍 Run Specific Tests**

```bash
cd tests/e2e

# Run only health checks
go test -v -run TestIngestionAPIHealth

# Run only investigation tests
go test -v -run TestCreateInvestigation

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## **🐛 Troubleshooting**

### Services Not Starting?
```bash
# Check service status
cd OnionEye_Code
docker compose ps

# View logs
docker compose logs ingestion-api
docker compose logs -f  # Follow all logs

# Test health endpoints
curl http://localhost:8080/health  # Ingestion API
curl http://localhost:8084/health  # Dashboard API
```

### Tests Timing Out?
```bash
# Increase timeout
go test -timeout 5m -v ./...

# Or set environment variable
export TEST_TIMEOUT=60s
go test -v ./...
```

### Port Already in Use?
```bash
# Find what's using port 8080
sudo lsof -i :8080

# Or change ports in docker-compose.yml
```

### Module Errors?
```bash
cd tests/e2e
go mod tidy
go mod download
```

---

## **⚙️ Custom Configuration**

Set custom URLs if your services run on different ports:

```bash
export INGESTION_API_URL="http://localhost:8080"
export DASHBOARD_API_URL="http://localhost:8084"
export STORAGE_API_URL="http://localhost:8081"
export POSTGRES_URL="postgres://onioneye:onioneye@localhost:5432/onioneye?sslmode=disable"
export REDIS_URL="localhost:6379"

# Then run tests
cd tests/e2e
go test -v ./...
```

---

## **📊 Useful Commands**

```bash
# Basic
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -short ./...             # Skip long tests
go test -run TestName ./...      # Run specific test

# Coverage
go test -cover ./...             # Show coverage %
go test -coverprofile=c.out ./...  # Generate report
go tool cover -html=c.out        # View in browser

# Debugging
go test -v -run TestName ./...   # Debug specific test
go test -failfast ./...          # Stop on first failure
go test -count=1 ./...          # Disable test cache
```

---

## **🎯 CI/CD (Automated)**

Tests run automatically on GitHub when you:
- Push to `main` or `develop`
- Create a pull request
- Manually trigger via Actions tab

To manually trigger:
1. Go to your GitHub repository
2. Click **Actions** tab
3. Select **CI Pipeline**
4. Click **Run workflow**

---

## **✅ Verification Checklist**

Before running tests:
- [ ] Go 1.22+ installed
- [ ] Docker installed and running
- [ ] In correct directory
- [ ] Services started (for full tests)

---

## **📚 More Information**

- **Detailed Guide:** `./guide.md`
- **Architecture:** `./README.md`
- **CI/CD Details:** `../summaries/ci-cd.md`

---

**🎉 You're ready to test!**

Need help? Check service logs: `docker compose logs -f`
