// Command sandbox is the render-sandbox HTTP service.
//
// Endpoints:
//
//	POST /render    body: {"job_id","source_url","html","viewport":{...}}
//	                -> 200 image/png (real pixel screenshot)
//	POST /sanitize  body: {"html"} -> 200 {"sanitized_html": "..."}
//	                (lets other services reuse the deep sanitizer for staging)
//	GET  /health    liveness
//
// The service performs its own defensive re-sanitization on every render,
// regardless of what the caller sent, and drives a network-isolated
// JS-disabled Chromium instance.
package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/onioneye/render-sandbox/internal/render"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n := 0
	for _, r := range v {
		if r < '0' || r > '9' {
			return def
		}
		n = n*10 + int(r-'0')
	}
	if n == 0 {
		return def
	}
	return n
}

func envDur(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

// requireToken returns a gin middleware that validates the X-OnionEye-Token
// header against an expected value. Fail-closed: when expected is empty all
// requests are rejected (operator must set RENDER_SANDBOX_TOKEN).
func requireToken(expected string) gin.HandlerFunc {
	if expected == "" {
		return func(c *gin.Context) {
			log.Printf("[auth] missing RENDER_SANDBOX_TOKEN; rejecting %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "sandbox token not configured"})
			c.Abort()
		}
	}
	return func(c *gin.Context) {
		got := c.GetHeader("X-OnionEye-Token")
		if subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid or missing token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	port := env("RENDER_SANDBOX_PORT", "8084")
	token := env("RENDER_SANDBOX_TOKEN", "")
	cfg := &render.Config{
		ChromePath:     env("CHROME_BINARY", ""),
		ViewportWidth:  envInt("RENDER_VIEWPORT_WIDTH", 1280),
		ViewportHeight: envInt("RENDER_VIEWPORT_HEIGHT", 1024),
		RenderTimeout:  envDur("RENDER_TIMEOUT", 12*time.Second),
		MaxConcurrent:  envInt("RENDER_MAX_CONCURRENT", 2),
	}
	sb := render.New(cfg, cfg.ChromePath)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "render-sandbox",
			"chrome":  cfg.ChromePath,
		})
	})
	r.POST("/render", requireToken(token), handleRender(sb))
	r.POST("/sanitize", requireToken(token), handleSanitize())

	log.Printf("render-sandbox listening on :%s (chrome=%q width=%d height=%d timeout=%s max_concurrent=%d token=%v)",
		port, cfg.ChromePath, cfg.ViewportWidth, cfg.ViewportHeight, cfg.RenderTimeout, cfg.MaxConcurrent, token != "")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}