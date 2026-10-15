package e2e

import (
"net/http"
"testing"
"time"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

// TestDashboardAPI_HealthCheck verifies dashboard API health.
func TestDashboardAPI_HealthCheck(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Dashboard API", config.DashboardHealth())

	t.Run("Dashboard health endpoint responds", func(t *testing.T) {
resp, err := client.client.R().Get(config.DashboardHealth())
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		assert.Contains(t, resp.String(), "ok")
	})
}

// TestDashboardAPI_StatsEndpoint tests the stats endpoint.
func TestDashboardAPI_StatsEndpoint(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Dashboard API", config.DashboardHealth())

	resp, err := client.client.R().Get(config.DashboardAPIURL + "/api/v1/dashboard/stats")
	require.NoError(t, err)

	if resp.StatusCode() == http.StatusOK {
		stats := ParseObject(t, resp)
		expectedFields := []string{"actors", "links", "jobs_done", "wallets", "pgp_keys", "flagged"}
		for _, field := range expectedFields {
			assert.Contains(t, stats, field, "Stats should contain %s", field)
		}
		t.Logf("Stats: %+v", stats)
	}
}

// TestDashboardAPI_GraphEndpoint tests the graph endpoint.
func TestDashboardAPI_GraphEndpoint(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Dashboard API", config.DashboardHealth())

	resp, err := client.client.R().Get(config.DashboardAPIURL + "/api/v1/dashboard/graph")
	require.NoError(t, err)

	if resp.StatusCode() == http.StatusOK {
		graph := ParseObject(t, resp)
		assert.Contains(t, graph, "nodes", "Graph should contain nodes")
		assert.Contains(t, graph, "links", "Graph should contain links (or edges)")
		t.Logf("Graph has %d nodes", len(graph["nodes"].([]interface{})))
	}
}

// TestDashboardAPI_JobsEndpoint tests the jobs endpoint.
func TestDashboardAPI_JobsEndpoint(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Dashboard API", config.DashboardHealth())

	resp, err := client.client.R().Get(config.DashboardAPIURL + "/api/v1/dashboard/jobs")
	require.NoError(t, err)

	if resp.StatusCode() == http.StatusOK {
		jobs := ParseArray(t, resp)
		t.Logf("Found %d jobs", len(jobs))
		if len(jobs) > 0 {
			job := jobs[0].(map[string]interface{})
			assert.True(t, HasAnyKey(job, "job_id", "target", "status"), "Job should have core fields")
		}
	}
}

// TestDashboardAPI_ResponseTime tests response times.
func TestDashboardAPI_ResponseTime(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Dashboard API", config.DashboardHealth())

	endpoints := []string{
		"/api/v1/dashboard/stats",
		"/api/v1/dashboard/graph",
		"/api/v1/dashboard/jobs",
	}

	for _, endpoint := range endpoints {
		t.Run("Response time: "+endpoint, func(t *testing.T) {
start := time.Now()
			_, err := client.client.R().Get(config.DashboardAPIURL + endpoint)
			elapsed := time.Since(start)

			require.NoError(t, err)
			t.Logf("Endpoint %s responded in %v", endpoint, elapsed)
			assert.Less(t, elapsed.Seconds(), 10.0, "Response should be under 10 seconds")
		})
	}
}
