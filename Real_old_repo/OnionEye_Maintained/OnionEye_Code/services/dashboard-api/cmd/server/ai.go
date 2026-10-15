package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Compose file is mounted read-only into the dashboard-api container so it can
// drive `docker compose` restarts of the whole OnionEye stack (see Dockerfile
// and docker-compose.yml). Kept as a const here so the restart handlers stay
// self-documenting.
const composeFile = "/compose/docker-compose.yml"

// ollamaRunning reports whether the Ollama container (AI profile) is currently
// UP in this compose project. The Docker container-list endpoint only returns
// running containers by default, so presence on the list implies "running".
func ollamaRunning(client *http.Client) bool {
	containers, err := listComposeContainers(client)
	if err != nil {
		return false
	}
	for _, ctr := range containers {
		if ctr.Labels["com.docker.compose.service"] == "ollama" {
			return true
		}
	}
	return false
}

// aiStatusHandler reports the live AI-profile state so the frontend can decide
// whether to show AI-only log feeds and AI restart controls.
//
//	GET /api/v1/dashboard/ai/status
func aiStatusHandler(c *gin.Context) {
	running := ollamaRunning(dockerClient())
	c.JSON(http.StatusOK, gin.H{
		"ai_enabled":     running,
		"ollama_running": running,
		"ollama_url":     "http://ollama:11434",
	})
}

// aiRestartRequest is the body expected by POST /ai/restart.
type aiRestartRequest struct {
	Enabled bool `json:"enabled"`
}

// aiRestartHandler brings the whole compose stack back up either WITH the AI
// profile (Ollama enabled) or WITHOUT it. Because `--profile ai` is what gates
// the ollama container, "restart with AI" runs `docker compose --profile ai up
// -d`, while "restart without AI" runs plain `up -d` then explicitly stops the
// ollama container.
//
//	POST /api/v1/dashboard/ai/restart  {"enabled": true|false}
func aiRestartHandler(c *gin.Context) {
	var req aiRestartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "detail": err.Error()})
		return
	}

	project := os.Getenv("DOCKER_PROJECT")
	if project == "" {
		project = "onioneye_code"
	}

	mode := "STANDARD"
	aiFlag := ""
	if req.Enabled {
		mode = "AI-ENHANCED"
		aiFlag = "--profile ai"
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	// Fire the compose rebuild in a DETACHED background shell so the response is
	// delivered before the stack cycles. Because `up -d --force-recreate` may
	// recreate THIS very container, we must not block the request on the command
	// — the dashboard-api process would be killed mid-wait and the client would
	// hang until the HTTP timeout instead of getting the "restarting" ack.
	restartScript := fmt.Sprintf(
		"nohup docker compose -f %s -p %s up -d --force-recreate %s > /tmp/ai-restart.log 2>&1 &",
		composeFile, project, aiFlag,
	)
	cmd := exec.CommandContext(ctx, "sh", "-c", restartScript)
	cmd.Env = append(os.Environ(), "DOCKER_PROJECT="+project)
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "compose restart failed to launch",
			"mode":    mode,
			"details": err.Error(),
		})
		return
	}

	// When disabling AI, also stop the Ollama container (detached) so no AI
	// processing happens after the plain `up -d`.
	if !req.Enabled {
		stopScript := fmt.Sprintf(
			"nohup docker compose -f %s -p %s stop ollama > /tmp/ai-restart.log 2>&1 &",
			composeFile, project,
		)
		_ = exec.CommandContext(ctx, "sh", "-c", stopScript).Run()
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "restarting",
		"mode":        mode,
		"ai_enabled":  req.Enabled,
		"restart_eta": "~15s",
		"note":        "docker-compose stack restarting with " + mode + " profile",
		"log":         "restart detached — see /tmp/ai-restart.log in dashboard-api",
	})
}

// ── AI log feeds ──────────────────────────────────────────────────────────────
// The two AI integration points surface as virtual feeds so analysts can watch
// the exact AI activity without wading through every container. Lines are pulled
// from the real ingestion-api + ollama container logs and filtered for the
// relevant keywords — mirroring how the "pipeline" virtual feed synthesises
// Redis state. These feeds are ONLY populated when Ollama is actually running.

const (
	feedAI      = "ai"       // combined AI activity
	feedAISeed  = "ai-seed"  // initial seed / term expansion (qwen2.5:3b)
	feedAIGuard = "ai-guard" // CASM / safety guard block (llama-guard3:1b)
)

// aiKeywords maps a virtual feed to the log substrings that make a line "about"
// that specific AI stage.
var aiKeywords = map[string][]string{
	feedAI: {
		"ollama", "/api/generate", "qwen", "llama", "llm", "expander",
		"expand", "guard", "seed", "tier2", "unsafe", "block", "safe",
	},
	feedAISeed: {
		"ollama", "/api/generate", "qwen", "expand", "expander", "seed", "llm",
	},
	feedAIGuard: {
		"guard", "llama", "unsafe", "block", "safe", "violat", "tier2",
	},
}

// aiFeed builds the virtual AI log feed. It stitches together real container
// logs (ingestion-api drives both expander + guard; ollama is the engine) and
// filters them to AI-relevant lines. When Ollama is not running it returns an
// explicit offline marker so the feed never silently looks broken.
func aiFeed(client *http.Client, byService map[string]string, tail int, kind string) []string {
	out := []string{}
	ts := time.Now().UTC().Format("2006/01/02 15:04:05")

	if !ollamaRunning(client) {
		out = append(out, fmt.Sprintf("%s [%s] AI engine OFFLINE — Ollama not running (restart stack with AI)", ts, kind))
		return out
	}

	keywords := aiKeywords[kind]
	if keywords == nil {
		keywords = aiKeywords[feedAI]
	}

	sources := []struct {
		svc string
		id  string
	}{
		{"ingestion-api", byService["ingestion-api"]},
		{"ollama", byService["ollama"]},
	}
	seen := map[string]bool{}
	for _, src := range sources {
		if src.id == "" {
			continue
		}
		lines, err := fetchContainerLogs(client, src.id, tail)
		if err != nil {
			continue
		}
		for _, ln := range lines {
			low := strings.ToLower(ln)
			matched := false
			for _, kw := range keywords {
				if strings.Contains(low, kw) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
			key := src.svc + "|" + ln
			if seen[key] {
				continue
			}
			seen[key] = true
			kindTag := feedAI
			if strings.Contains(low, "guard") || strings.Contains(low, "llama") ||
				strings.Contains(low, "unsafe") || strings.Contains(low, "block") ||
				strings.Contains(low, "violat") || strings.Contains(low, "safe") {
				kindTag = feedAIGuard
			}
			if strings.Contains(low, "seed") || strings.Contains(low, "expander") ||
				strings.Contains(low, "expand") || strings.Contains(low, "qwen") {
				kindTag = feedAISeed
			}
			if kind != feedAI && kindTag != kind {
				continue
			}
			out = append(out, fmt.Sprintf("%s [%s] %s", ts, src.svc, ln))
		}
	}

	if len(out) == 0 {
		out = append(out, fmt.Sprintf("%s [%s] no AI activity in recent window (tail=%d)", ts, kind, tail))
	}
	return out
}
