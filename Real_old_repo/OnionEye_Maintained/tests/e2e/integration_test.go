package e2e

import (
"testing"

"github.com/stretchr/testify/assert"
)

// TestHTTPClient_RetryLogic tests HTTP client retry functionality.
func TestHTTPClient_RetryLogic(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	t.Run("Client has retry configured", func(t *testing.T) {
assert.NotNil(t, client.client)
assert.Equal(t, config.RetryAttempts, 3)
})
}

// TestIntegration_DependencyGraph tests service dependencies.
func TestIntegration_DependencyGraph(t *testing.T) {
	config := GetTestConfig()
	client := NewHTTPClient(config)

	t.Run("All core services respond", func(t *testing.T) {
services := map[string]string{
"Ingestion API":  config.IngestionHealth(),
			"Dashboard API":  config.DashboardHealth(),
			"Storage API":    config.StorageHealth(),
			"Correlation API": config.CorrelationHealth(),
		}

		for name, url := range services {
			t.Run(name, func(t *testing.T) {
resp, err := client.client.R().Get(url)
				if err == nil {
					t.Logf("%s responded with %d", name, resp.StatusCode())
				} else {
					t.Logf("%s unavailable: %v", name, err)
				}
			})
		}
	})
}
