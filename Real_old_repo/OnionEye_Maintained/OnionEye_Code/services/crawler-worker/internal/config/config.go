// Package config loads crawler-worker settings from the environment.
// All values have sane defaults so the worker runs out-of-the-box in compose.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
	"time"
)

// Config holds every tunable for the worker. It mirrors the values in
// Architecture_Subsection_2.md (timeouts, politeness, retry budget, etc.).
type Config struct {
	// Redis (job queue + status tracking)
	RedisAddr     string
	RedisPassword string
	Stream        string // crawl:jobs
	Group         string // consumer group
	Consumer      string // this worker's consumer id
	ManualStream  string // crawl:manual — human-intervention queue

	// Tor sidecar
	TorSocksAddr     string // host:port SOCKS5 endpoint
	TorControlAddr   string // host:port control port for SIGNAL NEWNYM
	TorControlPasswd string // control-port password (non-local auth required by Tor)

	// Internal services
	ShieldURL  string // content-shield
	StorageURL string // storage-api

	// Render sandbox (on-demand screenshots)
	RenderSandboxURL   string
	RenderSandboxToken string
	MaxConcurrentScreenshots int

	// Auth for the worker's own HTTP surface (dashboard-api → worker).
	WorkerToken string

	// Pool / crawl behaviour
	WorkerCount      int
	MaxDepth         int
	MaxRetries       int
	PolitenessMinMs  int
	PolitenessMaxMs  int
	RequestTimeout   time.Duration
	IdleReadDeadline time.Duration
	MaxBodyBytes     int64
	MaxCircuitRot    int // max circuit rotations before giving up on a URL

	// Dashboard notifications
	DashboardAPIURL        string
	DashboardNotifyPath    string
	DashboardNotifyTimeout time.Duration

	// Read-only HTTP status surface (Tor reachability probe for the UI).
	// Dashboard-api lives on the frontend/backend networks and cannot reach the
	// tor-sidecar directly (crawler segment only), so the worker — which sits on
	// the same segment — reports Tor health on its behalf. Internal only; no
	// host port is published.
	HTTPAddr string

	// Worker anomaly detection & kill/remake
	AnomalyFastIdleTurns int           // consecutive fast no-job turns before kill+backoff (default 3)
	AnomalyFastTurnDur   time.Duration // a turn finishing quicker than this counts as 'fast' (default 1s)
	AnomalyMaxFailTurns  int           // consecutive network-dead job turns before kill+backoff (default 3)
	AnomalyBackoff       time.Duration // a killed slot waits this long before remake (default 15s)
	AnomalyMaxPerMinute  int           // max anomaly kill+notify events per rolling minute (default 5)
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envDur(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func randomSuffix() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Load reads the process environment into a Config.
func Load() *Config {
	return &Config{
		RedisAddr:     env("REDIS_ADDR", "redis:6379"),
		RedisPassword: env("REDIS_PASSWORD", ""),
		Stream:        env("REDIS_STREAM", "crawl:jobs"),
		Group:         env("REDIS_GROUP", "crawler-workers"),
		Consumer:      env("REDIS_CONSUMER", "worker-"+randomSuffix()),
		ManualStream:  env("MANUAL_STREAM", "crawl:manual"),

		TorSocksAddr:     env("TOR_SOCKS_ADDR", "tor-sidecar:9050"),
		TorControlAddr:   env("TOR_CONTROL_ADDR", "tor-sidecar:9051"),
		TorControlPasswd: env("TOR_CONTROL_PASSWORD", "onioneye"),

		ShieldURL:  env("SHIELD_URL", "http://content-shield:8082"),
		StorageURL: env("STORAGE_API_URL", "http://storage-api:8081"),

		RenderSandboxURL:         env("RENDER_SANDBOX_URL", "http://render-sandbox:8084"),
		RenderSandboxToken:       env("RENDER_SANDBOX_TOKEN", ""),
		MaxConcurrentScreenshots: envInt("MAX_CONCURRENT_SCREENSHOTS", 4),
		WorkerToken:              env("CRAWLER_WORKER_TOKEN", ""),

		WorkerCount:      envInt("WORKER_COUNT", 5),
		MaxDepth:         envInt("MAX_DEPTH", 3),
		MaxRetries:       envInt("MAX_RETRIES", 3),
		PolitenessMinMs:  envInt("POLITENESS_MIN_MS", 2000),
		PolitenessMaxMs:  envInt("POLITENESS_MAX_MS", 8000),
		// Tor onion fetches frequently stall; increase defaults so “deadline exceeded
		// while reading body” doesn’t prevent any crawl persistence.
		RequestTimeout:   envDur("REQUEST_TIMEOUT", 180*time.Second),
		IdleReadDeadline: envDur("IDLE_READ_DEADLINE", 60*time.Second),

		// Ensure we can fetch full HTML responses without truncation.
		// (Still capped by MaxBodyBytes for safety.)
		MaxBodyBytes:     int64(envInt("MAX_BODY_BYTES", 10<<20)),
		MaxCircuitRot:    envInt("MAX_CIRCUIT_ROTATIONS", 5),

		// Dashboard notifications
		DashboardAPIURL:        env("DASHBOARD_API_URL", "http://dashboard-api:3000"),
		DashboardNotifyPath:    env("DASHBOARD_NOTIFY_PATH", "/api/v1/dashboard/notify"),
		DashboardNotifyTimeout: envDur("DASHBOARD_NOTIFY_TIMEOUT", 3*time.Second),

		// Worker anomaly detection & kill/remake
		AnomalyFastIdleTurns: envInt("ANOMALY_FAST_IDLE_TURNS", 3),
		AnomalyFastTurnDur:   envDur("ANOMALY_FAST_TURN_DUR", 1*time.Second),
		AnomalyMaxFailTurns:  envInt("ANOMALY_MAX_FAIL_TURNS", 3),
		AnomalyBackoff:       envDur("ANOMALY_BACKOFF", 15*time.Second),
		AnomalyMaxPerMinute:  envInt("ANOMALY_MAX_PER_MINUTE", 5),

		HTTPAddr: env("HTTP_ADDR", ":8084"),
	}
}
