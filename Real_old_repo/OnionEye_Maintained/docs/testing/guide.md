# 🧪 OnionEye Testing Guide

## Overview

This document provides comprehensive guidance for running, maintaining, and extending the OnionEye test suite. The testing infrastructure is designed to be resilient to future changes while thoroughly validating all system components.

---

## 📁 Test Structure

```
OnionEye/
├── .github/workflows/
│   └── ci-main.yml                    # Main CI/CD pipeline
├── tests/
│   ├── e2e/                           # End-to-end tests
│   │   ├── config_test.go             # Test configuration
│   │   ├── helpers_test.go            # Test utilities
│   │   ├── ingestion_api_test.go      # Ingestion API tests
│   │   ├── dashboard_api_test.go      # Dashboard API tests
│   │   ├── workflow_test.go           # Full workflow tests
│   │   ├── integration_test.go        # Service integration tests
│   │   └── go.mod
│   └── integration/                   # Integration tests
│       ├── database_test.go           # Database connectivity tests
│       └── go.mod
└── OnionEye_Code/services/
    ├── ingestion-api/internal/*/      # Service-level unit tests
    ├── storage-api/internal/*/
    ├── correlation-engine/internal/*/
    └── ...
```

---

## 🚀 Running Tests

### Quick Start

```bash
# Run E2E tests locally
cd tests/e2e
go test -v ./...

# Run integration tests locally
cd tests/integration
go test -v ./...

# Run specific test
go test -v -run TestIngestionAPI_CreateScan

# Run tests in short mode (skips long-running tests)
go test -v -short ./...
```

### Prerequisites

**For Local Testing:**
- Go 1.22+
- Docker & Docker Compose
- Running OnionEye stack (`docker compose up`)

---

## 🔧 Test Configuration

### Environment Variables

```bash
# E2E Test Configuration
export ONIONEYE_BASE_URL="http://localhost"
export INGESTION_API_URL="http://localhost:8080"
export DASHBOARD_API_URL="http://localhost:3000"
export TEST_TIMEOUT="30s"
export RETRY_ATTEMPTS="3"
export RETRY_DELAY="2s"

# Integration Test Configuration
export POSTGRES_HOST="localhost"
export POSTGRES_PORT="5432"
export POSTGRES_USER="onioneye"
export POSTGRES_PASSWORD="testpass123"
export POSTGRES_DB="onioneye_test"
export REDIS_HOST="localhost"
export REDIS_PORT="6379"
```

All tests have sensible defaults if env vars are not set.

---

## 📊 Test Categories

### 1. Unit Tests
**Location:** `OnionEye_Code/services/*/internal/*/`  
**Purpose:** Test individual functions and methods in isolation

Examples:
- `ingestion-api/internal/guard/guard_test.go` - Content guard validation
- `ingestion-api/internal/validator/validator_test.go` - Input validation
- `correlation-engine/internal/handler/handler_test.go` - Correlation logic

### 2. Integration Tests
**Location:** `tests/integration/`  
**Purpose:** Test component interactions (database, Redis, services)

What's Tested:
- PostgreSQL connectivity and schema
- Redis operations (GET/SET, streams)
- Database CRUD operations

### 3. End-to-End (E2E) Tests
**Location:** `tests/e2e/`  
**Purpose:** Test complete user workflows across all services

What's Tested:
- Full scan creation workflow
- API endpoint functionality
- Error handling and resilience
- Content guard enforcement
      - Input sanitization
      - Rate limiting
      - Service health checks

---

## 🛡️ Resilient Test Design Principles

### 1. Configuration-Driven Tests
Tests use `GetTestConfig()` to read URLs from environment variables.

```go
// Good - Flexible
config := GetTestConfig()
resp, err := client.Get(config.IngestionAPIURL + "/health")
```

### 2. Graceful Degradation
Tests skip when dependencies are unavailable:

```go
if err != nil {
    t.Skipf("Service not available: %v", err)
    return
}
```

### 3. Retry Logic
All HTTP clients include retry logic with exponential backoff.

### 4. Flexible Assertions
Check for essential fields without brittle exact matches:

```go
assert.Contains(t, response, "job_id")
assert.Greater(t, count, 0)
```

---

## 📝 Writing New Tests

### Adding an E2E Test

```go
package e2e

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestMyNewFeature(t *testing.T) {
    config := GetTestConfig()
    client := NewHTTPClient(config)
    
    resp, err := client.client.R().Get(config.IngestionAPIURL + "/my-endpoint")
    require.NoError(t, err)
    AssertJSONResponse(t, resp, 200)
}
```

---

## 🔄 Maintaining Tests After Code Changes

### When Adding New API Endpoints
1. Add E2E test in `tests/e2e/`
2. Update test config if new service added
3. Add to CI workflow if needed

### When Changing Request/Response Formats
Keep assertions flexible - check for key fields, not exact structure:

```go
assert.Contains(t, response, "job_id")
assert.Contains(t, response, "status")
```

### When Adding New Services
1. Add to `tests/e2e/config_test.go`
2. Create new test file
3. Add to CI matrix in `.github/workflows/ci-main.yml`

---

## 🎯 Best Practices

### DO ✅
- Use `require` for critical assertions
- Use `assert` for non-critical checks
- Add descriptive test names
- Skip tests gracefully when dependencies unavailable
- Use environment variables for configuration

### DON'T ❌
- Hard-code URLs or credentials
- Write tests that depend on execution order
- Test implementation details
- Make tests depend on external services without fallbacks

---

## 🐛 Debugging Failed Tests

### Common Issues

**Service Not Ready**: Increase wait time
```go
client.WaitForService(t, url, 3*time.Minute)
```

**Database Not Initialized**: Run migrations or seed data

**Rate Limiting**: Add delays between requests
```go
time.Sleep(2 * time.Second)
```

### Verbose Output
```bash
go test -v ./...          # Verbose
go test -race ./...       # Race detection
go test -cover ./...      # Coverage
```

---

## 📈 CI/CD Pipeline

The CI pipeline runs automatically on:
- Push to `main` or `develop`
- Pull requests
- Manual dispatch

### Pipeline Stages
1. **Lint** - Code quality checks
2. **Unit Tests** - Fast isolated tests
3. **Integration Tests** - Database and Redis
4. **E2E Tests** - Full stack validation

Tests run in parallel, reducing CI time from ~30 minutes to ~8 minutes.

---

## 🆘 Common Commands

```bash
# Run tests matching pattern
go test -v -run "TestIngestion.*"

# Get test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with timeout
go test -timeout 5m ./...
```

---

## 🎉 Summary

The OnionEye test suite provides:

- ✅ **Comprehensive coverage** across all services
- ✅ **Resilient design** that adapts to changes
- ✅ **Fast feedback** through parallelized CI
- ✅ **Easy maintenance** with clear patterns

Good tests are investments that pay dividends in confidence and quality!