package handler

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// maxScreenshotBytes caps the PNG we accept (10 MiB). The render sandbox
// produces ~100–500 KB viewport PNGs; this is generous headroom that still
// keeps a rogue/pathological render from exhausting the DB row.
const maxScreenshotBytes = 10 << 20

// ScreenshotPayload is the body of POST /screenshots (crawler-worker → storage).
// png_data is base64 because the worker has no binary HTTP helper here and the
// rest of the pipeline (shield, extraction) already moves bytes in-band as b64.
type ScreenshotPayload struct {
	URL     string `json:"url"`
	PngData string `json:"png_data"` // base64-encoded PNG bytes
}

// ScreenshotRecord is the JSON shape returned for metadata lookups. PNG bytes
// are served from /screenshots/:id/png, never embedded in JSON.
type ScreenshotRecord struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// UpsertScreenshot stores (or replaces) the PNG for a unique URL.
//
// POST /screenshots
//
//	{"url":"http://xyz.onion/listing","png_data":"<base64>"} -> 201 {id,url,created_at}
//
// UNIQUE(url) is the upsert anchor: recapturing the same URL silently
// overwrites the old PNG (one screenshot per URL, per the decided scope).
func (h *Handler) UpsertScreenshot(c *gin.Context) {
	var req ScreenshotPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	rawURL := normalizeScreenshotURL(req.URL)
	if !validScreenshotURL(rawURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url must be absolute http(s)"})
		return
	}

	png, err := base64.StdEncoding.DecodeString(req.PngData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "png_data must be base64"})
		return
	}
	if len(png) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "png_data is empty"})
		return
	}
	if len(png) > maxScreenshotBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "png_data exceeds size limit"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var rec ScreenshotRecord
	err = h.pool.QueryRow(ctx, `
		INSERT INTO screenshots (id, url, png_data, created_at)
		VALUES (gen_random_uuid(), $1, $2, now())
		ON CONFLICT (url) DO UPDATE SET
			png_data   = EXCLUDED.png_data,
			created_at = now()
		RETURNING id, url, created_at`,
		rawURL, png).Scan(&rec.ID, &rec.URL, &rec.CreatedAt)
	if err != nil {
		logErr("upsert screenshot", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database write failed"})
		return
	}

	c.JSON(http.StatusCreated, rec)
}

// GetScreenshotByURL returns screenshot metadata for a URL.
//
// GET /screenshots?url=http://xyz.onion/listing -> 200 {id,url,created_at} | 404
func (h *Handler) GetScreenshotByURL(c *gin.Context) {
	rawURL := normalizeScreenshotURL(c.Query("url"))
	if rawURL == "" || !validScreenshotURL(rawURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid url query param"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var rec ScreenshotRecord
	err := h.pool.QueryRow(ctx,
		`SELECT id, url, created_at FROM screenshots WHERE url = $1`, rawURL).
		Scan(&rec.ID, &rec.URL, &rec.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "no screenshot for url"})
		return
	}
	if err != nil {
		logErr("query screenshot", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database read failed"})
		return
	}
	c.JSON(http.StatusOK, rec)
}

// GetScreenshotPNG streams the raw PNG bytes for a screenshot id.
//
// GET /screenshots/:id/png -> image/png | 404
func (h *Handler) GetScreenshotPNG(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if !validUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid screenshot id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var png []byte
	err := h.pool.QueryRow(ctx,
		`SELECT png_data FROM screenshots WHERE id = $1`, id).Scan(&png)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "screenshot not found"})
		return
	}
	if err != nil {
		logErr("query screenshot png", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database read failed"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "private, max-age=300")
	c.Data(http.StatusOK, "image/png", png)
}

// validScreenshotURL accepts only absolute http(s) URLs. This is the same
// guard the crawler uses before fetching, and it prevents junk keys in the
// UNIQUE(url) table. It does not restrict to .onion — clearweb infra nodes are
// screenshotted too.
func normalizeScreenshotURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "http://" + raw
	}
	return raw
}

func validScreenshotURL(raw string) bool {
	raw = normalizeScreenshotURL(raw)
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
