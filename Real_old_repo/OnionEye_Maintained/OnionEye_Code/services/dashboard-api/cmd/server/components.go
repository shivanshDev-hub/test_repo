package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// componentSpec describes one node of the tool-architecture diagram rendered
// in the right-rail TOOL MAP panel.
type componentSpec struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Service   string `json:"service"`  // compose service name
	Category  string `json:"category"` // ui | api | engine | worker | shield | ai | infra | network
	Role      string `json:"role"`     // human-readable role description
	HealthURL string `json:"-"`        // optional HTTP /health URL
	HealthTCP string `json:"-"`        // optional host:port TCP target
}

var componentList = []componentSpec{
	{ID: "nginx", Name: "Nginx Gateway", Service: "nginx", Category: "ui", Role: "Reverse proxy + auth", HealthTCP: "nginx:80"},
	{ID: "frontend", Name: "Frontend SPA", Service: "frontend", Category: "ui", Role: "React SPA"},
	{ID: "dashboard-api", Name: "Dashboard API", Service: "dashboard-api", Category: "api", Role: "Graph / stats / control", HealthURL: "http://dashboard-api:3000/health"},
	{ID: "storage-api", Name: "Storage API", Service: "storage-api", Category: "api", Role: "Postgres-backed CRUD", HealthURL: "http://storage-api:8081/health"},
	{ID: "ingestion-api", Name: "Ingestion API", Service: "ingestion-api", Category: "api", Role: "Job dispatch + AI seed", HealthTCP: "ingestion-api:8080"},
	{ID: "correlation-engine", Name: "Correlation Engine", Service: "correlation-engine", Category: "engine", Role: "Persona-link builder", HealthTCP: "correlation-engine:8083"},
	{ID: "content-filter", Name: "Content Filter", Service: "content-filter", Category: "shield", Role: "ML abusive-content classifier", HealthTCP: "content-filter:8082"},
	{ID: "content-shield", Name: "Content Shield", Service: "content-shield", Category: "shield", Role: "CSAM hash-set guard", HealthTCP: "content-shield:8080"},
	{ID: "crawler-worker", Name: "Crawler Worker", Service: "crawler-worker", Category: "worker", Role: "Tor dark-web crawler pool"},
	{ID: "redis", Name: "Redis", Service: "redis", Category: "infra", Role: "Queue + cache + pub/sub", HealthTCP: "redis:6379"},
	{ID: "postgres", Name: "PostgreSQL", Service: "postgres", Category: "infra", Role: "Primary data store", HealthTCP: "postgres:5432"},
	{ID: "tor-sidecar", Name: "Tor Sidecar", Service: "tor-sidecar", Category: "network", Role: "SOCKS5 proxy", HealthTCP: "tor-sidecar:9050"},
	{ID: "ollama", Name: "Ollama (AI)", Service: "ollama", Category: "ai", Role: "LLM inference engine", HealthTCP: "ollama:11434"},
}

// componentsStatusHandler probes each known OnionEye service and returns a
// summary so the frontend can render the TOOL MAP diagram with green / red
// status indicators.
//
//	GET /api/v1/dashboard/components
func componentsStatusHandler(c *gin.Context) {
	client := dockerClient()
	containers, err := listComposeContainers(client)
	if err != nil {
		containers = nil
	}

	running := map[string]bool{}
	for _, ctr := range containers {
		if svc := ctr.Labels["com.docker.compose.service"]; svc != "" {
			running[svc] = true
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)
	out := make([]gin.H, 0, len(componentList))
	summary := gin.H{"total": len(componentList), "up": 0, "down": 0}

	for _, spec := range componentList {
		status := "down"
		detail := "container not running"

		if running[spec.Service] {
			if spec.HealthTCP != "" {
				if tcpOK(spec.HealthTCP) {
					status, detail = "up", fmt.Sprintf("%s reachable", spec.HealthTCP)
				} else {
					status, detail = "degraded", fmt.Sprintf("%s unreachable", spec.HealthTCP)
				}
			} else if spec.HealthURL != "" {
				if httpOK(client, spec.HealthURL) {
					status, detail = "up", fmt.Sprintf("%s healthy", spec.HealthURL)
				} else {
					status, detail = "degraded", fmt.Sprintf("%s unhealthy", spec.HealthURL)
				}
			} else {
				status, detail = "up", "container running"
			}
		} else if spec.HealthTCP != "" && tcpOK(spec.HealthTCP) {
			status, detail = "up", fmt.Sprintf("%s reachable (no container)", spec.HealthTCP)
		} else if spec.HealthURL != "" && httpOK(client, spec.HealthURL) {
			status, detail = "up", fmt.Sprintf("%s healthy (no container)", spec.HealthURL)
		}

		if status == "up" {
			summary["up"] = summary["up"].(int) + 1
		} else {
			summary["down"] = summary["down"].(int) + 1
		}

		out = append(out, gin.H{
			"id":       spec.ID,
			"name":     spec.Name,
			"service":  spec.Service,
			"category": spec.Category,
			"role":     spec.Role,
			"status":   status,
			"detail":   detail,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i]["category"].(string) < out[j]["category"].(string)
	})

	c.JSON(http.StatusOK, gin.H{
		"components": out,
		"summary":    summary,
		"timestamp":  now,
	})
}

func tcpOK(addr string) bool {
	return tcpDialOK(addr, 1*time.Second)
}

func httpOK(client *http.Client, url string) bool {
	if client == nil {
		return false
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 500
}
