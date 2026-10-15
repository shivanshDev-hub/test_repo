# Test Maintenance Guide for Developers

> **📋 Team Structure:** See [TEAM_STRUCTURE.md](../TEAM_STRUCTURE.md) for current team assignments and responsibilities.  
> **Questions?** Contact the relevant team based on the component you're modifying.

## 🚨 **READ THIS BEFORE MAKING CHANGES** 🚨

This document explains when and how to update tests when you modify the OnionEye codebase.

---

## **When Do I Need to Update Tests?**

### ✅ **YES - Update Tests** if you change:

#### 1. **API Endpoints/Routes**
**Examples:**
- Added new endpoint: `/api/v1/reports`
- Renamed endpoint: `/investigations` → `/cases`
- Changed HTTP method: `POST /submit` → `PUT /submit`
- Added/removed query parameters

**What to update:**
- `tests/e2e/ingestion_api_test.go` (for Ingestion API changes)
- `tests/e2e/dashboard_api_test.go` (for Dashboard API changes)
- `tests/e2e/workflow_test.go` (for multi-service workflows)

**How:**
```go
// Before
resp, err := client.R().Get(config.IngestionAPIURL + "/api/v1/investigations")

// After (if endpoint changed to /cases)
resp, err := client.R().Get(config.IngestionAPIURL + "/api/v1/cases")
```

---

#### 2. **Response Format/JSON Structure**
**Examples:**
- Renamed field: `job_id` → `jobId` or `job_id` → `taskId`
- Added required field: response now includes `priority`
- Removed field: `deprecated_field` no longer returned
- Changed data type: `timestamp` from string to integer

**What to update:**
- Test assertions in `tests/e2e/*_test.go`
- Response validation in `tests/e2e/helpers_test.go`

**How:**
```go
// Before
assert.Contains(t, string(resp.Body()), "job_id")

// After (if field renamed)
assert.Contains(t, string(resp.Body()), "jobId")

// Or use flexible assertion
obj := ParseObject(t, resp)
assert.True(t, HasAnyKey(obj, "job_id", "jobId", "taskId"))
```

---

#### 3. **Database Schema**
**Examples:**
- Added table: `reports`
- Added column: `actors.priority`
- Renamed column: `content` → `post_content`
- Changed column type: `created_at` from timestamp to bigint
- Added foreign key or constraint

**What to update:**
- `tests/integration/database_test.go`
- Schema validation tests

**How:**
```go
// Add test for new column
func TestActorPriorityColumn(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    var exists bool
    err := db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'actors' AND column_name = 'priority'
        )
    `).Scan(&exists)
    
    require.NoError(t, err)
    assert.True(t, exists, "priority column should exist in actors table")
}
```

---

#### 4. **Service Ports**
**Examples:**
- Changed ingestion-api from 8080 → 8090
- Changed dashboard-api from 3000 → 8084
- Added new service on port 9000

**What to update:**
- Default URLs in `tests/e2e/config_test.go`
- Port mappings in `.github/workflows/ci-main.yml`

**How:**
```go
// In tests/e2e/config_test.go
func DefaultIngestionAPIURL() string {
    return "http://localhost:8090"  // Changed from 8080
}
```

---

#### 5. **Environment Variables**
**Examples:**
- Added new required env var: `API_KEY`
- Renamed env var: `DB_URL` → `DATABASE_URL`
- Changed default value

**What to update:**
- `tests/e2e/config_test.go` (env var loading)
- `.github/workflows/ci-main.yml` (CI env section)

**How:**
```go
// In tests/e2e/config_test.go - Add new field to config
type TestConfig struct {
    IngestionAPIURL string
    DashboardAPIURL string
    APIKey          string  // Add new field
}

func GetTestConfig() *TestConfig {
    return &TestConfig{
        IngestionAPIURL: getEnv("INGESTION_API_URL", DefaultIngestionAPIURL()),
        DashboardAPIURL: getEnv("DASHBOARD_API_URL", DefaultDashboardAPIURL()),
        APIKey:          getEnv("API_KEY", ""),  // Add new env var
    }
}
```

---

#### 6. **Authentication/Authorization**
**Examples:**
- Added API key requirement
- Added JWT token authentication
- Added rate limiting
- Added IP whitelisting

**What to update:**
- `tests/e2e/helpers_test.go` (client setup)
- All test files that make HTTP requests

**How:**
```go
// In tests/e2e/helpers_test.go
func NewTestClient() *resty.Client {
    config := GetTestConfig()
    return resty.New().
        SetTimeout(30 * time.Second).
        SetRetryCount(3).
        SetHeader("Authorization", "Bearer " + config.APIKey).  // Add auth
        SetHeader("Content-Type", "application/json")
}
```

---

### ❌ **NO - Tests OK as-is** if you only:

- Fixed typos in comments or documentation
- Refactored internal code without changing behavior
- Updated documentation
- Changed log messages
- Optimized performance without changing outputs
- Fixed bugs that make existing behavior match expected behavior

---

## **How to Find Which Tests to Update**

### Quick Search Commands:

```bash
# Find tests that use a specific endpoint
cd /home/shivansh/Documents/my_projects/skill_work/OnionEye-
rg "/api/v1/investigations" tests/

# Find tests that check a specific field
rg "job_id" tests/e2e/

# Find tests that use a specific port
rg "8080" tests/

# Find all test files for a service
ls tests/e2e/*_test.go

# Find database schema tests
rg "CREATE TABLE|information_schema" tests/
```

---

## **Testing Your Test Updates**

### Step 1: Run Tests Locally

```bash
# Start services
cd OnionEye_Code
docker compose up -d
sleep 30

# Run updated tests
cd tests/e2e
go test -v ./...

# Run integration tests
cd tests/integration
go test -v ./...
```

### Step 2: Verify All Tests Pass

```bash
# Should see output like:
# PASS
# ok      github.com/onioneye/tests/e2e    15.234s
```

### Step 3: Push and Check CI

```bash
git add tests/
git commit -s -m "test: update tests for <your change>"
git push
```

CI will automatically run all tests. Check GitHub Actions for results.

---

## **Common Test Update Patterns**

### Pattern 1: API Endpoint Changed

**Change:** Renamed `/api/v1/investigations` → `/api/v1/cases`

**Steps:**
```bash
# 1. Find all occurrences
cd /home/shivansh/Documents/my_projects/skill_work/OnionEye-
rg "/api/v1/investigations" tests/

# 2. Edit each file found
vim tests/e2e/ingestion_api_test.go
# Replace /api/v1/investigations with /api/v1/cases

# 3. Run tests to verify
cd tests/e2e && go test -v ./...
```

### Pattern 2: Response Field Renamed

**Change:** `job_id` → `task_id`

**Steps:**
```go
// Option 1: Update to new field name
assert.Contains(t, string(resp.Body()), "task_id")

// Option 2: Support both (better for gradual migration)
obj := ParseObject(t, resp)
assert.True(t, HasAnyKey(obj, "job_id", "task_id"))
```

### Pattern 3: New Required Field

**Change:** Investigations now require `priority` field

**Steps:**
```go
// Add priority to all test payloads
payload := map[string]interface{}{
    "target_url": "https://example.onion",
    "username": "testuser",
    "priority": "high",  // ADD THIS
}
```

---

## **Automated Test Execution**

Tests run automatically in these scenarios:

### 1. **On Every Push**
- Pushes to `main`, `develop`, `feature/*`, `fix/*` branches
- Runs: lint → unit tests → integration tests → E2E tests

### 2. **On Every Pull Request**
- When you open/update a PR targeting `main` or `develop`
- **Must pass before PR can be merged**

### 3. **Daily Scheduled Run**
- Every day at 2 AM UTC
- Full health check of all services
- Catches issues from environmental drift

### 4. **Manual Trigger**
- Go to GitHub Actions → Select workflow → "Run workflow"

---

## **What Happens When Tests Fail in CI?**

1. ✅ **PR is blocked** from merging
2. 📧 **You get notified** via GitHub
3. 🔍 **Check the logs**:
   - Go to GitHub Actions tab
   - Click on the failed workflow run
   - View logs for the failed job
4. 🔧 **Fix and push again**:
   ```bash
   # Fix the issue
   git add .
   git commit -s -m "fix: update tests for API changes"
   git push
   # CI will automatically re-run
   ```

---

## **Quick Reference Card**

| What Changed | File to Update | How to Find It |
|-------------|----------------|----------------|
| API endpoint | `tests/e2e/*_api_test.go` | `rg "/api/v1/" tests/` |
| Response field | `tests/e2e/*_test.go` | `rg "field_name" tests/` |
| Database schema | `tests/integration/database_test.go` | `ls tests/integration/` |
| Service port | `tests/e2e/config_test.go` | `rg "8080\|3000" tests/` |
| Env variable | `config_test.go` + CI workflow | `rg "getEnv" tests/` |
| Auth/headers | `tests/e2e/helpers_test.go` | `vim tests/e2e/helpers_test.go` |

---

## **Best Practices**

1. ✅ **Update tests in the same commit** as your code changes
2. ✅ **Run tests locally before pushing**
3. ✅ **Use flexible assertions** when possible (e.g., `HasAnyKey`)
4. ✅ **Write descriptive test names** that explain what's being tested
5. ✅ **Add comments** explaining why a test exists
6. ✅ **Keep tests independent** - each test should work alone
7. ✅ **Use test helpers** - don't repeat HTTP setup code

---

## **Need Help?**

- Check examples in existing test files
- Read `./guide.md` for detailed patterns
- Read `./quick-start.md` for setup instructions
- **Contact the relevant team** (see [TEAM_STRUCTURE.md](../TEAM_STRUCTURE.md)):
  - AI/ML tests → **AI Team**
  - Backend/Frontend tests → **Backend & Frontend Team**
  - Testing infrastructure/CI → **Live Testing & Debugging Team**
- Ask in your PR comments
- Check CI logs for specific failure messages

---

## **⚠️ IMPORTANT REMINDER**

**If you change the API, you MUST change the tests. No exceptions!**

Tests are part of the codebase, not an afterthought. Keep them in sync.

**Failure to update tests will block your PR from merging.**

