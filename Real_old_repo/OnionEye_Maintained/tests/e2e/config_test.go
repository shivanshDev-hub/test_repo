package e2e

import (
"fmt"
"os"
"time"
)

// TestConfig holds configuration for E2E tests.
// URLs are env-driven so the suite survives host/port/compose changes.
type TestConfig struct {
	BaseURL          string
	IngestionAPIURL  string
	DashboardAPIURL  string
	StorageAPIURL    string
	CorrelationURL   string
	ContentShieldURL string
	ContentFilterURL string
	Timeout          time.Duration
	RetryAttempts    int
	RetryDelay       time.Duration
}

// GetTestConfig returns test configuration from environment or defaults.
func GetTestConfig() *TestConfig {
	return &TestConfig{
		BaseURL:          getEnv("ONIONEYE_BASE_URL", "http://localhost"),
		IngestionAPIURL:  getEnv("INGESTION_API_URL", "http://localhost:8080"),
		DashboardAPIURL:  getEnv("DASHBOARD_API_URL", "http://localhost:3000"),
		StorageAPIURL:    getEnv("STORAGE_API_URL", "http://localhost:8081"),
		CorrelationURL:   getEnv("CORRELATION_API_URL", "http://localhost:8083"),
		ContentShieldURL: getEnv("CONTENT_SHIELD_URL", "http://localhost:8082"),
		ContentFilterURL: getEnv("CONTENT_FILTER_URL", "http://localhost:8085"),
		Timeout:          getDuration("TEST_TIMEOUT", 30*time.Second),
		RetryAttempts:    getInt("RETRY_ATTEMPTS", 3),
		RetryDelay:       getDuration("RETRY_DELAY", 2*time.Second),
	}
}

func (c *TestConfig) IngestionHealth() string   { return c.IngestionAPIURL + "/healthz" }
func (c *TestConfig) DashboardHealth() string   { return c.DashboardAPIURL + "/health" }
func (c *TestConfig) StorageHealth() string     { return c.StorageAPIURL + "/health" }
func (c *TestConfig) CorrelationHealth() string { return c.CorrelationURL + "/healthz" }
func (c *TestConfig) ShieldHealth() string      { return c.ContentShieldURL + "/health" }
func (c *TestConfig) FilterHealth() string      { return c.ContentFilterURL + "/healthz" }

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var i int
		if _, err := fmt.Sscanf(value, "%d", &i); err == nil {
			return i
		}
	}
	return defaultValue
}
