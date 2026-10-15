// Package config loads ingestion-api configuration from environment variables.
package config

import (
	"os"
	"strconv"
)

// Config holds all runtime configuration for the ingestion pipeline.
type Config struct {
	Port            string
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	RateLimitMax    int
	RateLimitWindow int
	SessionWindow   int
	AllowlistPath   string
	OllamaURL       string
	ExpanderModel   string
	GuardModel      string
	MaxSeedTerms    int
	MaxDepthDefault int
	MaxDepthMax     int
}

// Load reads configuration from environment with sensible defaults.
func Load() *Config {
	return &Config{
		Port:            getEnv("INGESTION_PORT", "8080"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RateLimitMax:    getEnvInt("RATE_LIMIT_MAX", 5),
		RateLimitWindow: getEnvInt("RATE_LIMIT_WINDOW", 60),
		SessionWindow:   getEnvInt("SESSION_WINDOW_SECONDS", 900),
		AllowlistPath:   getEnv("ALLOWLIST_PATH", "/config/target_allowlist.yaml"),
		OllamaURL:       getEnv("OLLAMA_URL", "http://localhost:11434"),
		ExpanderModel:   getEnv("EXPANDER_MODEL", "qwen2.5:3b"),
		GuardModel:      getEnv("GUARD_MODEL", "llama-guard3:1b"),
		MaxSeedTerms:    getEnvInt("MAX_SEED_TERMS", 10),
		MaxDepthDefault: getEnvInt("MAX_DEPTH_DEFAULT", 3),
		MaxDepthMax:     getEnvInt("MAX_DEPTH_MAX", 5),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
