package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/onioneye/render-sandbox/internal/render"
	"github.com/onioneye/render-sandbox/internal/sanitize"
)

// renderRequest is the worker/caller payload. html must already have been
// sanitized upstream; the sandbox re-sanitizes it anyway.
type renderRequest struct {
	JobID     string          `json:"job_id"`
	SourceURL string          `json:"source_url"`
	HTML      string          `json:"html"`
	Viewport  viewportRequest `json:"viewport,omitempty"`
}

type viewportRequest struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

func handleRender(sb *render.Sandbox) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req renderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		if len(req.HTML) > sanitize.MaxInputBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "html too large"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		// Clamp client-supplied viewport values defensively.
		// Even though the sandbox currently uses a fixed viewport override in
		// the renderer config, we accept the field for future compatibility.
		if req.Viewport.Width != 0 {
			if req.Viewport.Width < 16 {
				req.Viewport.Width = 16
			}
			if req.Viewport.Width > 4096 {
				req.Viewport.Width = 4096
			}
		}
		if req.Viewport.Height != 0 {
			if req.Viewport.Height < 16 {
				req.Viewport.Height = 16
			}
			if req.Viewport.Height > 4096 {
				req.Viewport.Height = 4096
			}
		}

		defer cancel()

		res, err := sb.Render(ctx, req.HTML)
		if err != nil {
			log.Printf("render failed job=%s url=%s err=%v", req.JobID, req.SourceURL, err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "render failed", "detail": err.Error()})
			return
		}

		c.Header("Content-Type", "image/png")
		c.Header("X-OnionEye-Sanitized", "true")
		c.Header("X-OnionEye-JobID", req.JobID)
		c.Data(http.StatusOK, "image/png", res.PNG)
	}
}

// sanitizeRequest is the staging payload for /sanitize.
type sanitizeRequest struct {
	HTML string `json:"html"`
}

func handleSanitize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req sanitizeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		out, err := sanitize.HTML(req.HTML)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"sanitized_html": out})
	}
}