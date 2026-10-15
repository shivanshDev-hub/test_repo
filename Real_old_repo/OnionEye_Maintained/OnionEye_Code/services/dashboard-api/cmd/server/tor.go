package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type torStatusResponse struct {
	Status       string `json:"status"`        // connected | connecting | unavailable
	SocksOK      bool   `json:"socks_ok"`      // TCP 9050 reachable from dashboard-api container
	ControlOK    bool   `json:"control_ok"`    // TCP 9051 reachable from dashboard-api container
	LastSignal   string `json:"last_signal"`   // best-effort: e.g. "BOOTSTRAP" status
	BootstrapPct int    `json:"bootstrap_pct"` // best-effort: 0-100, -1 when unknown
	Note         string `json:"note"`
}

func torConnAddrs() (socksAddr, controlAddr string) {
	// Allow local overrides; default matches docker-compose envs.
	// This keeps the endpoint usable even if caller doesn't mount crawler-worker config.
	socksAddr = os.Getenv("TOR_SOCKS_ADDR")
	if socksAddr == "" {
		socksAddr = "tor-sidecar:9050"
	}
	controlAddr = os.Getenv("TOR_CONTROL_ADDR")
	if controlAddr == "" {
		controlAddr = "tor-sidecar:9051"
	}
	return
}

func tcpDialOK(addr string, timeout time.Duration) bool {
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	_ = c.Close()
	return true
}

// torWorkerStatusURL returns the crawler-worker's read-only Tor-status surface.
// Network segmentation keeps tor-sidecar on the `crawler` network only; the
// worker (which shares that segment) reports reachability on our behalf.
func torWorkerStatusURL() string {
	if v := os.Getenv("TOR_STATUS_URL"); v != "" {
		return v
	}
	return "http://crawler-worker:8084/tor/status"
}

// torStatusHandler is a lightweight health/status endpoint for the Tor sidecar.
// It asks the crawler-worker (authoritative, same network segment as Tor) for a
// real probe, falling back to the legacy direct TCP check from this container.
func torStatusHandler(c *gin.Context) {
	// Primary: ask the crawler-worker, which shares the crawler network with Tor.
	if body, ok := fetchTorStatusFromWorker(); ok {
		wd := struct {
			Status    string `json:"status"`
			SocksOK   bool   `json:"socks_ok"`
			ControlOK bool   `json:"control_ok"`
			Note      string `json:"note"`
		}{}
		// Best-effort decode; even on decode failure the JSON above is zero and
		// we can still return it to keep the UI populated.
		_ = json.Unmarshal(body, &wd)
		if wd.Note == "" {
			wd.Note = "reported by crawler-worker (Tor segment probe)"
		}
		c.Header("Content-Type", "application/json")
		_ = json.NewEncoder(c.Writer).Encode(torStatusResponse{
			Status:       wd.Status,
			SocksOK:      wd.SocksOK,
			ControlOK:    wd.ControlOK,
			BootstrapPct: -1,
			Note:         wd.Note,
		})
		return
	}

	// Fallback: best-effort TCP reachability from this container (expected to be
	// false under the segmented compose network; kept for standalone deploys).
	socksAddr, controlAddr := torConnAddrs()
	socksOK := tcpDialOK(socksAddr, 1*time.Second)
	controlOK := tcpDialOK(controlAddr, 1*time.Second)

	status := "unavailable"
	bootstrapPct := -1
	lastSignal := ""

	if socksOK && controlOK {
		status = "connected"
	} else if socksOK || controlOK {
		status = "connecting"
	}

	note := fmt.Sprintf("TCP reachability: socks(%s)=%t control(%s)=%t", socksAddr, socksOK, controlAddr, controlOK)

	out := torStatusResponse{
		Status:       status,
		SocksOK:      socksOK,
		ControlOK:    controlOK,
		LastSignal:   lastSignal,
		BootstrapPct: bootstrapPct,
		Note:         note,
	}

	c.Header("Content-Type", "application/json")
	_ = json.NewEncoder(c.Writer).Encode(out)
}

// fetchTorStatusFromWorker GETs the crawler-worker's /tor/status endpoint and
// returns the raw body on success.
func fetchTorStatusFromWorker() ([]byte, bool) {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(torWorkerStatusURL())
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return nil, false
	}
	return body, true
}

// registerTorRoutes is called from main.go.
func registerTorRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/dashboard")
	api.GET("/tor/status", torStatusHandler)

	// Back-compat alias.
	api.GET("/tor", func(c *gin.Context) { torStatusHandler(c) })
}
