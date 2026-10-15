package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// manualReviewHandler fetches pending manual review items from the crawl:manual
// Redis stream. Items are pushed by the crawler-worker when it hits CAPTCHA,
// bot blocks, or repeated scan failures that need human intervention.
//
//	GET /api/v1/dashboard/review/pending
func manualReviewHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Read last 100 entries from crawl:manual stream
	result, err := redisClient().XRevRange(ctx, "crawl:manual", "+", "-").Result()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch review queue", "detail": err.Error()})
		return
	}

	items := []gin.H{}
	now := time.Now().UTC()
	for _, msg := range result {
		vals := msg.Values
		flaggedAt := vals["flagged_at"]
		expiresIn := vals["expires_in"]

		// Parse flagged timestamp
		var ts time.Time
		if flaggedAtStr, ok := flaggedAt.(string); ok {
			ts, _ = time.Parse(time.RFC3339, flaggedAtStr)
		}

		// Check expiry (5-minute TTL per architecture)
		var expired bool
		if expiresInStr, ok := expiresIn.(string); ok {
			if ts.IsZero() {
				expired = true
			} else if expiresInStr == "300" {
				expired = now.Sub(ts) > 5*time.Minute
			}
		}

		items = append(items, gin.H{
			"id":         msg.ID,
			"job_id":     vals["job_id"],
			"target":     vals["target"],
			"reason":     vals["reason"],
			"flagged_at": flaggedAt,
			"expired":    expired,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(items),
		"items": items,
	})
}

// dismissReviewRequest is the expected payload for POST /review/dismiss
type dismissReviewRequest struct {
	ID string `json:"id" binding:"required"`
}

// dismissReviewHandler removes a manual review item from the crawl:manual stream
// after the analyst has resolved it (either by manually visiting the site or
// deciding to skip it).
//
//	POST /api/v1/dashboard/review/dismiss  {"id": "stream-entry-id"}
func dismissReviewHandler(c *gin.Context) {
	var req dismissReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "detail": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// XDEL removes the entry from the stream
	deleted, err := redisClient().XDel(ctx, "crawl:manual", req.ID).Result()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to dismiss item", "detail": err.Error()})
		return
	}

	if deleted == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found or already dismissed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "dismissed",
		"id":      req.ID,
		"deleted": deleted,
		"note":    "manual review item removed from queue",
	})
}
