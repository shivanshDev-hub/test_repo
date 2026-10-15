package e2e

import (
"encoding/json"
"fmt"
"net/http"
"os"
"testing"
"time"

"github.com/go-resty/resty/v2"
"github.com/stretchr/testify/require"
)

// HTTPClient wraps resty client with retry logic.
type HTTPClient struct {
	client *resty.Client
	config *TestConfig
}

// NewHTTPClient creates a new HTTP client for testing.
func NewHTTPClient(config *TestConfig) *HTTPClient {
	client := resty.New().
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryAttempts).
		SetRetryWaitTime(config.RetryDelay).
		SetRetryMaxWaitTime(config.Timeout).
		SetHeader("Accept", "application/json")

	return &HTTPClient{
		client: client,
		config: config,
	}
}

// inCI reports whether tests are running in GitHub Actions / CI.
func inCI() bool {
	return os.Getenv("CI") != ""
}

// RequireService skips (locally) or fails (CI) when a service is unreachable.
// Use this at the top of tests that need a live stack.
func (c *HTTPClient) RequireService(t *testing.T, name, healthURL string) {
	t.Helper()

	resp, err := c.client.R().Get(healthURL)
	if err == nil && resp.StatusCode() < 500 {
		return
	}

	msg := fmt.Sprintf("%s not available at %s", name, healthURL)
	if err != nil {
		msg = fmt.Sprintf("%s: %v", msg, err)
	} else {
		msg = fmt.Sprintf("%s (status %d)", msg, resp.StatusCode())
	}

	if inCI() {
		t.Fatalf("%s — the CI stack should have started this service", msg)
	}
	t.Skipf("%s — skipping (start the OnionEye stack to run this test)", msg)
}

// WaitForService waits for a service to be healthy.
func (c *HTTPClient) WaitForService(t *testing.T, url string, maxWait time.Duration) {
	t.Helper()
	deadline := time.Now().Add(maxWait)

	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := c.client.R().Get(url)
		if err == nil && resp.StatusCode() == http.StatusOK {
			t.Logf("Service at %s is ready", url)
			return
		}
		lastErr = err
		time.Sleep(2 * time.Second)
	}

	if inCI() {
		t.Fatalf("Service at %s did not become healthy within %v: %v", url, maxWait, lastErr)
	}
	t.Skipf("Service at %s did not become healthy within %v: %v", url, maxWait, lastErr)
}

// AssertJSONResponse validates JSON response structure.
func AssertJSONResponse(t *testing.T, resp *resty.Response, expectedStatus int) {
	t.Helper()
	require.Equal(t, expectedStatus, resp.StatusCode(),
		"Expected status %d, got %d. Body: %s", expectedStatus, resp.StatusCode(), resp.String())

	contentType := resp.Header().Get("Content-Type")
	require.Contains(t, contentType, "application/json",
"Expected JSON content type, got %s", contentType)
}

// ParseObject decodes a JSON object response.
func ParseObject(t *testing.T, resp *resty.Response) map[string]interface{} {
	t.Helper()
	body := resp.Body()
	if len(body) == 0 {
		return map[string]interface{}{}
	}
	var out map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &out), "response is not a JSON object: %s", resp.String())
	return out
}

// ParseArray decodes a JSON array response.
func ParseArray(t *testing.T, resp *resty.Response) []interface{} {
	t.Helper()
	body := resp.Body()
	if len(body) == 0 {
		return []interface{}{}
	}
	var out []interface{}
	if err := json.Unmarshal(body, &out); err == nil {
		return out
	}
	// Wrapped arrays: {"leads":[...]} / {"actors":[...]}
	var obj map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &obj), "response is neither array nor object: %s", resp.String())
	for _, key := range []string{"leads", "actors", "scans", "jobs", "nodes", "links", "edges"} {
		if raw, ok := obj[key]; ok {
			if arr, ok := raw.([]interface{}); ok {
				return arr
			}
		}
	}
	return []interface{}{}
}

// HasAnyKey reports whether the object contains at least one of the given keys.
func HasAnyKey(obj map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		if _, ok := obj[k]; ok {
			return true
		}
	}
	return false
}

// RetryOperation retries an operation with exponential backoff.
func RetryOperation(t *testing.T, maxAttempts int, operation func() error) error {
	t.Helper()
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		if err := operation(); err == nil {
			return nil
		} else {
			lastErr = err
			waitTime := time.Duration(i+1) * 2 * time.Second
			t.Logf("Retry %d/%d failed: %v. Waiting %v...", i+1, maxAttempts, err, waitTime)
			time.Sleep(waitTime)
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxAttempts, lastErr)
}

// CleanupTestData removes test data created during tests (best-effort).
func CleanupTestData(t *testing.T, client *HTTPClient, dashboardURL string) {
	t.Helper()
	t.Log("Cleaning up test data...")

	_, err := client.client.R().Delete(dashboardURL + "/api/v1/dashboard/delete/all")
	if err != nil {
		t.Logf("Warning: Failed to cleanup test data: %v", err)
	}
}
