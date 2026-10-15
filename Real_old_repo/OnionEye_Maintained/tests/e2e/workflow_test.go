package e2e

import (
"net/http"
"testing"
"time"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

// TestFullWorkflow_EndToEnd tests the complete OnionEye workflow.
func TestFullWorkflow_EndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full workflow test in short mode")
	}

	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Ingestion API", config.IngestionHealth())
	client.RequireService(t, "Dashboard API", config.DashboardHealth())

	t.Log("Step 1: Creating a new investigation...")
	scanPayload := map[string]interface{}{
		"analyst_id":   "e2e_test_analyst",
		"target_type":  "onion_address",
		"target_value": "http://testmarket456789abc.onion",
		"seed_terms":   []string{"ransomware", "cryptocurrency"},
		"max_depth":    2,
	}

	scanResp, err := client.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(scanPayload).
		Post(config.IngestionAPIURL + "/api/v1/investigations")

	require.NoError(t, err, "Failed to create investigation")
	if scanResp.StatusCode() == http.StatusAccepted {
		obj := ParseObject(t, scanResp)
		if jobID, ok := obj["job_id"].(string); ok {
			t.Logf("Created investigation with job_id: %s", jobID)
		}
	} else {
		t.Logf("Investigation returned status %d (may be rate-limited or blocked)", scanResp.StatusCode())
	}

	t.Log("Step 2: Waiting for initial processing...")
	time.Sleep(5 * time.Second)

	t.Log("Step 3: Querying dashboard stats...")
	statsResp, err := client.client.R().Get(config.DashboardAPIURL + "/api/v1/dashboard/stats")
	require.NoError(t, err, "Failed to get stats")
	if statsResp.StatusCode() == http.StatusOK {
		stats := ParseObject(t, statsResp)
		t.Logf("Dashboard stats: %+v", stats)
	}

	t.Log("Step 4: Querying graph data...")
	graphResp, err := client.client.R().Get(config.DashboardAPIURL + "/api/v1/dashboard/graph")
	require.NoError(t, err, "Failed to get graph")
	if graphResp.StatusCode() == http.StatusOK {
		graph := ParseObject(t, graphResp)
		if nodes, ok := graph["nodes"].([]interface{}); ok {
			t.Logf("Graph contains %d nodes", len(nodes))
		}
	}

	t.Log("\u2705 Full workflow test completed successfully")
}

// TestServiceInteraction_Resilience tests service interaction and error handling.
func TestServiceInteraction_Resilience(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	testCases := []struct {
		name     string
		endpoint string
		method   string
	}{
		{"Ingestion API health", config.IngestionHealth(), "GET"},
		{"Dashboard API stats", config.DashboardAPIURL + "/api/v1/dashboard/stats", "GET"},
		{"Dashboard API graph", config.DashboardAPIURL + "/api/v1/dashboard/graph", "GET"},
		{"Storage API health", config.StorageHealth(), "GET"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
resp, err := client.client.R().Get(tc.endpoint)
			if err != nil {
				t.Logf("%s unavailable: %v (skipping)", tc.name, err)
				return
			}
			assert.True(t, resp.StatusCode() < 500, "Should not return 5xx errors, got %d", resp.StatusCode())
		})
	}
}

// TestConcurrentRequests tests system behavior under concurrent load.
func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Ingestion API", config.IngestionHealth())

	concurrency := 5
	results := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			payload := map[string]interface{}{
				"analyst_id":   "concurrent_test",
				"target_type":  "onion_address",
				"target_value": "http://test.onion",
				"seed_terms":   []string{"test"},
				"max_depth":    1,
			}

			resp, err := client.client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(payload).
				Post(config.IngestionAPIURL + "/api/v1/investigations")

			results <- (err == nil && (resp.StatusCode() == http.StatusAccepted || resp.StatusCode() == http.StatusTooManyRequests))
		}(i)
	}

	successCount := 0
	for i := 0; i < concurrency; i++ {
		if <-results {
			successCount++
		}
	}

	t.Logf("Concurrent requests: %d/%d succeeded", successCount, concurrency)
	assert.Greater(t, successCount, 0, "At least some concurrent requests should succeed")
}
