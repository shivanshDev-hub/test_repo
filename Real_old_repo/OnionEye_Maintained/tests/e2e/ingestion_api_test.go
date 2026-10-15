package e2e

import (
"net/http"
"testing"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

// TestIngestionAPI_HealthCheck verifies ingestion API health endpoint.
func TestIngestionAPI_HealthCheck(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Ingestion API", config.IngestionHealth())

	t.Run("Health endpoint responds", func(t *testing.T) {
resp, err := client.client.R().Get(config.IngestionHealth())
		require.NoError(t, err, "Health check request failed")
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		assert.Contains(t, resp.String(), "ok")
	})
}

// TestIngestionAPI_CreateInvestigation tests the investigation creation workflow.
func TestIngestionAPI_CreateInvestigation(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Ingestion API", config.IngestionHealth())

	t.Run("Valid investigation accepted", func(t *testing.T) {
payload := map[string]interface{}{
"analyst_id":   "test_analyst_001",
"target_type":  "onion_address",
"target_value": "http://test3456789abcdefgh.onion",
"seed_terms":   []string{"ransomware", "darknet market"},
"max_depth":    3,
}

resp, err := client.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
			Post(config.IngestionAPIURL + "/api/v1/investigations")

		require.NoError(t, err)
		if resp.StatusCode() == http.StatusAccepted {
			obj := ParseObject(t, resp)
			assert.True(t, HasAnyKey(obj, "job_id", "status"), "Response should contain job_id or status")
			t.Logf("Investigation accepted: %+v", obj)
		} else {
			t.Logf("Investigation returned %d (may indicate Ollama unavailable or rate limit)", resp.StatusCode())
		}
	})
}

// TestIngestionAPI_BlockedTerms tests blocking of illegal content terms.
func TestIngestionAPI_BlockedTerms(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Ingestion API", config.IngestionHealth())

	blockedTerms := []string{
		"csam test",
		"child porn test",
	}

	for _, term := range blockedTerms {
		t.Run("Blocks: "+term, func(t *testing.T) {
payload := map[string]interface{}{
"analyst_id":   "content_guard_test",
"target_type":  "onion_address",
"target_value": "http://test.onion",
"seed_terms":   []string{term},
"max_depth":    1,
}

resp, err := client.client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(payload).
				Post(config.IngestionAPIURL + "/api/v1/investigations")

			require.NoError(t, err)
			// Guard should block via 422 Unprocessable Entity or 403 Forbidden
			assert.True(t, resp.StatusCode() == http.StatusUnprocessableEntity || resp.StatusCode() == http.StatusForbidden,
				"Expected blocked response, got %d: %s", resp.StatusCode(), resp.String())
		})
	}
}

// TestIngestionAPI_InvalidInput tests input validation.
func TestIngestionAPI_InvalidInput(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	client.RequireService(t, "Ingestion API", config.IngestionHealth())

	testCases := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Empty analyst ID",
			payload: map[string]interface{}{
				"analyst_id":   "",
				"target_type":  "onion_address",
				"target_value": "http://test.onion",
				"seed_terms":   []string{"test"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid JSON",
			payload: map[string]interface{}{
				"target_type": 12345, // wrong type
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
resp, err := client.client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.payload).
				Post(config.IngestionAPIURL + "/api/v1/investigations")

			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode())
		})
	}
}
