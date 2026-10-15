package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/onioneye/crawler-worker/internal/config"
	"github.com/onioneye/crawler-worker/internal/worker"
)

// torStatus mirrors the dashboard-api's torStatusResponse shape so the UI can
// render the same fields (status / socks_ok / control_ok / note) regardless of
// which service produced it.
type torStatus struct {
	Status       string `json:"status"`        // connected | connecting | unavailable
	SocksOK      bool   `json:"socks_ok"`      // TCP 9050 reachable from this container
	ControlOK    bool   `json:"control_ok"`    // TCP 9051 reachable from this container
	BootstrapPct int    `json:"bootstrap_pct"` // -1 when unknown (no control parse)
	Note         string `json:"note"`
}

func tcpOK(addr string) bool {
	c, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return false
	}
	_ = c.Close()
	return true
}

// startStatusServer exposes a read-only HTTP surface for the UI/orchestrator.
// The worker shares the `crawler` network with the tor-sidecar, so it can
// measure real Tor reachability that dashboard-api — segmented onto the
// frontend/backend networks — cannot probe directly. It is best-effort: if the
// listener fails, the crawl loop is unaffected.
func startStatusServer(cfg *config.Config, pool *worker.Pool) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/tor/status", func(w http.ResponseWriter, _ *http.Request) {
		socksOK := tcpOK(cfg.TorSocksAddr)
		controlOK := tcpOK(cfg.TorControlAddr)

		status := "unavailable"
		switch {
		case socksOK && controlOK:
			status = "connected"
		case socksOK || controlOK:
			status = "connecting"
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(torStatus{
			Status:       status,
			SocksOK:      socksOK,
			ControlOK:    controlOK,
			BootstrapPct: -1,
			Note: fmt.Sprintf("worker probe: socks(%s)=%t control(%s)=%t",
				cfg.TorSocksAddr, socksOK, cfg.TorControlAddr, controlOK),
		})
	})

	if pool != nil {
		mux.Handle("/screenshot", requireWorkerToken(cfg.WorkerToken, pool.ScreenshotHandler()))
	}

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Status surfacing is best-effort; never take the crawl loop down.
		}
	}()
	return srv
}

func requireWorkerToken(expected string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expected == "" {
			http.Error(w, `{"error":"worker token not configured"}`, http.StatusUnauthorized)
			return
		}
		got := r.Header.Get("X-OnionEye-Token")
		if subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
			http.Error(w, `{"error":"invalid or missing token"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
