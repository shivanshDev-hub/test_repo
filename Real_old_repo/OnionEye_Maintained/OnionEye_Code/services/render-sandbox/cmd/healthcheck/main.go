// Command healthcheck is a tiny static binary that verifies the render-sandbox
// HTTP service is alive.  It replaces the former wget-based Docker HEALTHCHECK
// probe so the runtime image can ship without wget or ca-certificates.
//
// Exit 0 = healthy, exit 1 = unhealthy (Docker marks the container as failed).
package main

import (
	"net/http"
	"os"
	"time"
)

func main() {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://127.0.0.1:8084/health")
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
