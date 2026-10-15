package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// dockerSocketPath is the Docker Engine unix socket. The dashboard-api
// container mounts it read-only so it can serve REAL per-service logs
// (the same output `docker logs` prints) to the crawl console.
const dockerSocketPath = "/var/run/docker.sock"

type dockerContainerInfo struct {
	ID     string            `json:"Id"`
	Names  []string          `json:"Names"`
	Labels map[string]string `json:"Labels"`
}

// dockerClient returns an HTTP client that talks to the Docker Engine via
// the docker-sock-proxy (TCP), or falls back to the unix socket if
// DOCKER_API_URL is not set (local dev without compose).
func dockerClient() *http.Client {
	apiURL := os.Getenv("DOCKER_API_URL")
	if apiURL != "" {
		return &http.Client{
			Timeout:   4 * time.Second,
			Transport: &http.Transport{},
		}
	}
	// Fallback: local unix socket (development only).
	return &http.Client{
		Timeout: 4 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "unix", dockerSocketPath)
			},
		},
	}
}

// ── Redis queue / pipeline introspection ──────────────────────────────────────
// The crawler-worker emits structured [job …] / [fetch …] lines to stdout and
// those already flow into docker logs. But what sits INSIDE Redis — the
// crawl:jobs stream, consumer-group pending backlog, and task_status:*
// documents — is invisible unless we query it directly. The helpers below
// synthesise pseudo-log lines from Redis state so the SYS LOGS panel shows the
// exact data-plane flow the user asked for.

var (
	rdbOnce sync.Once
	rdb     *redis.Client
)

// redisClient lazily initialises a single Redis client for the dashboard-api
// lifetime.  Returns nil when REDIS_HOST is not set (offline dev).
func redisClient() *redis.Client {
	rdbOnce.Do(func() {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			return
		}
		port := os.Getenv("REDIS_PORT")
		if port == "" {
			port = "6379"
		}
		rdb = redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: os.Getenv("REDIS_PASSWORD"),
			PoolSize: 2,
		})
	})
	return rdb
}

// queueSnapshot returns log-style lines describing the CURRENT state of the
// crawl pipeline: stream length, consumer-group pending backlog, the newest
// crawl:jobs entries (with target, depth, job_id), and every task_status:<job>
// document.  Timestamps are formatted identically to Docker logs so the lines
// blend seamlessly with the real container output.
func queueSnapshot(client *redis.Client) []string {
	if client == nil {
		return []string{"[pipeline] redis client unavailable — set REDIS_HOST to enable live queue view"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ts := time.Now().UTC().Format("2006/01/02 15:04:05")
	var lines []string

	// ── Stream length ────────────────────────────────────────────────────────
	xlen, err := client.XLen(ctx, "crawl:jobs").Result()
	if err != nil {
		lines = append(lines, fmt.Sprintf("%s [pipeline] XLEN crawl:jobs ERROR %v", ts, err))
		return lines
	}

	// ── Consumer-group pending (claimed but not yet ACKed) ───────────────────
	pend, err := client.XPending(ctx, "crawl:jobs", "crawler-workers").Result()
	pendText := "n/a"
	if err == nil && pend != nil {
		pendText = strconv.FormatInt(pend.Count, 10)
	}
	lines = append(lines, fmt.Sprintf("%s [pipeline] crawl:jobs length=%d pending=%s",
		ts, xlen, pendText))

	// ── Newest stream entries (most recent → oldest, up to 5) ────────────────
	entries, err := client.XRevRange(ctx, "crawl:jobs", "+", "-").Result()
	if err == nil && len(entries) > 0 {
		lines = append(lines, fmt.Sprintf("%s [pipeline] newest %d/%d entries:", ts, minInt(len(entries), 5), len(entries)))
		for i, e := range entries {
			if i >= 5 {
				break
			}
			v := e.Values
			jobID := fmt.Sprintf("%v", v["job_id"])
			target := fmt.Sprintf("%v", v["target"])
			depth := fmt.Sprintf("%v", v["depth_remaining"])
			terms := fmt.Sprintf("%v", v["terms"])
			if len(terms) > 60 {
				terms = terms[:57] + "..."
			}
			lines = append(lines, fmt.Sprintf("%s [pipeline]   %s  job=%s depth=%s target=%s  terms=%s",
				ts, e.ID, jobID, depth, target, terms))
		}
	} else {
		lines = append(lines, fmt.Sprintf("%s [pipeline] crawl:jobs stream empty", ts))
	}

	// ── task_status:* documents (one per job) ────────────────────────────────
	keys, err := client.Keys(ctx, "task_status:*").Result()
	if err == nil && len(keys) > 0 {
		sort.Strings(keys)
		for _, k := range keys {
			val, err := client.Get(ctx, k).Result()
			if err != nil {
				continue
			}
			var doc map[string]any
			status := val
			if json.Unmarshal([]byte(val), &doc) == nil {
				if s, ok := doc["status"].(string); ok {
					status = s
				}
				if ts2, ok := doc["last_updated"].(string); ok && ts2 != "" {
					status += " @" + ts2
				}
			}
			lines = append(lines, fmt.Sprintf("%s [pipeline] %s = %s", ts, k, status))
		}
	} else if err == nil {
		lines = append(lines, fmt.Sprintf("%s [pipeline] no task_status documents", ts))
	}

	return lines
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// dockerAPIBase returns the base URL for Docker Engine API calls.
// In compose, DOCKER_API_URL is http://docker-sock-proxy:2375; in local dev
// without the proxy we fall back to the historic "http://docker" hostname
// (which docker-compose used to provide via the default network).
func dockerAPIBase() string {
	if u := os.Getenv("DOCKER_API_URL"); u != "" {
		return u
	}
	return "http://docker"
}

// listComposeContainers returns every running container belonging to this
// compose project (identified by the com.docker.compose.project label), keyed
// by service name via listComposeServices.
func listComposeContainers(client *http.Client) ([]dockerContainerInfo, error) {
	project := os.Getenv("DOCKER_PROJECT")
	if project == "" {
		project = "onioneye_code"
	}
	filter := fmt.Sprintf(`{"label":["com.docker.compose.project=%s"]}`, project)
	u := dockerAPIBase() + "/containers/json?filters=" + url.QueryEscape(filter)
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var containers []dockerContainerInfo
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, err
	}
	return containers, nil
}

// fetchContainerLogs tails a container's real stdout/stderr stream. Docker
// multiplexes the two streams into 8-byte frame headers, so we unwrap them
// (a dropped frame simply truncates that tail request).
func fetchContainerLogs(client *http.Client, id string, tail int) ([]string, error) {
	u := fmt.Sprintf("%s/containers/%s/logs?stdout=1&stderr=1&tail=%d", dockerAPIBase(), id, tail)
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var lines []string
	hdr := make([]byte, 8)
	for {
		if _, err := io.ReadFull(resp.Body, hdr); err != nil {
			break
		}
		size := binary.BigEndian.Uint32(hdr[4:8])
		payload := make([]byte, size)
		if _, err := io.ReadFull(resp.Body, payload); err != nil {
			break
		}
		line := string(payload)
		for len(line) > 0 && (line[len(line)-1] == '\n' || line[len(line)-1] == '\r') {
			line = line[:len(line)-1]
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// logsHandler serves real Docker logs for every compose service.
//
//	GET /api/v1/dashboard/logs?service=crawler-worker&tail=60  → one service
//	GET /api/v1/dashboard/logs?tail=40                          → all services
func logsHandler(c *gin.Context) {
	client := dockerClient()

	service := c.DefaultQuery("service", "all")
	tail := 60
	if v := c.Query("tail"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			tail = n
		}
	}
	if tail > 2000 {
		tail = 2000
	}

	containers, err := listComposeContainers(client)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "docker daemon unreachable: " + err.Error(),
			"service": service,
			"logs":    []string{},
		})
		return
	}

	byService := map[string]string{}
	for _, ctr := range containers {
		if svc := ctr.Labels["com.docker.compose.service"]; svc != "" {
			byService[svc] = ctr.ID
		}
	}

	// Single-service request.
	if service != "all" {
		// ── Virtual "pipeline" service — live Redis queue + task_status ──────
		if service == "pipeline" {
			c.JSON(http.StatusOK, gin.H{"service": "pipeline", "logs": queueSnapshot(redisClient())})
			return
		}

		// ── Virtual AI feeds — seed expansion + CASM/guard activity. These are
		// only meaningful when the Ollama (AI) engine is running; aiFeed returns
		// an explicit OFFLINE marker otherwise. ───────────────────────────────
		if service == feedAI || service == feedAISeed || service == feedAIGuard {
			c.JSON(http.StatusOK, gin.H{"service": service, "logs": aiFeed(client, byService, tail, service)})
			return
		}

		id, ok := byService[service]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "no container for service: " + service})
			return
		}
		lines, err := fetchContainerLogs(client, id, tail)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "service": service, "logs": []string{}})
			return
		}

		// ── Enrich "redis" feed with live queue snapshot ─────────────────────
		if service == "redis" {
			if extra := queueSnapshot(redisClient()); len(extra) > 0 {
				lines = append(extra, lines...)
			}
		}

		c.JSON(http.StatusOK, gin.H{"service": service, "logs": lines})
		return
	}

	// "all" → every compose service's container logs + virtual feeds.
	names := make([]string, 0, len(byService))
	for svc := range byService {
		names = append(names, svc)
	}
	sort.Strings(names)

	out := gin.H{}
	for _, svc := range names {
		lines, err := fetchContainerLogs(client, byService[svc], tail)
		if err != nil {
			out[svc] = gin.H{"error": err.Error()}
			continue
		}
		// Enrich redis with queue snapshot in "all" mode too.
		if svc == "redis" {
			if extra := queueSnapshot(redisClient()); len(extra) > 0 {
				lines = append(extra, lines...)
			}
		}
		out[svc] = lines
	}

	// Add virtual "pipeline" feed alongside real compose services.
	out["pipeline"] = queueSnapshot(redisClient())

	// Add virtual AI feeds (seed expansion + CASM/guard). Only surfaced when the
	// Ollama AI engine is actually running so the UI can hide them in standard
	// (non-AI) mode.
	if ollamaRunning(client) {
		for _, kind := range []string{feedAI, feedAISeed, feedAIGuard} {
			out[kind] = aiFeed(client, byService, tail, kind)
		}
	}

	c.JSON(http.StatusOK, gin.H{"services": out})
}
