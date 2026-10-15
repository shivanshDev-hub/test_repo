package classifier

// handler.go — Gin HTTP handler that wraps the Phase 8A Pipeline.
//
// Routes registered by RegisterRoutes:
//
//	POST /classify        — classify a single image supplied as multipart/form-data
//	POST /classify/job    — classify all images in a JSON job request
//	GET  /healthz         — liveness probe
//
// The handler intentionally does not talk to a database; it is a stateless
// inference service.  The caller (crawler-worker) is responsible for
// persisting ClassificationResult rows to the content_flags table.

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler binds the Pipeline to Gin routes.
type Handler struct {
	pipeline *Pipeline
}

// NewHandler constructs a Handler.  pipeline must not be nil.
func NewHandler(pipeline *Pipeline) *Handler {
	if pipeline == nil {
		panic("classifier.NewHandler: pipeline must not be nil")
	}
	return &Handler{pipeline: pipeline}
}

// RegisterRoutes attaches all routes to r.
func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/healthz", h.healthz)
	r.POST("/classify", h.classifySingle)
	r.POST("/classify/job", h.classifyJob)
}

// ---------------------------------------------------------------------------
// GET /healthz
// ---------------------------------------------------------------------------

func (h *Handler) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ---------------------------------------------------------------------------
// POST /classify
//
// Accepts multipart/form-data with the field "image" (raw bytes) and optional
// "src" (image URL/path for tracing) and "actor_id".
//
// Returns ClassificationResult as JSON.
// ---------------------------------------------------------------------------

func (h *Handler) classifySingle(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("missing 'image' field: %v", err)})
		return
	}
	defer file.Close()

	// Guard: reject suspiciously large uploads (32 MiB).
	const maxBytes = 32 << 20
	if header.Size > maxBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": fmt.Sprintf("image too large: %d bytes (max %d)", header.Size, maxBytes),
		})
		return
	}

	imageBytes, err := io.ReadAll(io.LimitReader(file, maxBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("reading image: %v", err)})
		return
	}

	src := c.Request.FormValue("src")
	if src == "" {
		src = header.Filename
	}

	result, err := h.pipeline.Classify(c.Request.Context(), src, imageBytes)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Optionally override actor_id from the form field.
	if actor := c.Request.FormValue("actor_id"); actor != "" {
		result.ActorID = actor
	}

	statusCode := http.StatusOK
	if result.Action == ActionContentDiscarded {
		statusCode = http.StatusForbidden // 403 signals the caller to discard
	}

	c.JSON(statusCode, result)
}

// ---------------------------------------------------------------------------
// POST /classify/job
//
// Accepts JSON body matching JobRequest.  Each image's RawBytes must be
// base64-encoded under the key "raw_b64".
//
// Returns JobResult as JSON.
// ---------------------------------------------------------------------------

// jobRequestWire is the over-the-wire shape; RawBytes arrive as base64.
type jobRequestWire struct {
	JobID         string           `json:"job_id"`
	SourceURL     string           `json:"source_url"`
	ProvenanceTag string           `json:"provenance_tag"`
	Images        []imageInputWire `json:"images"`
	HandlesOnPage []string         `json:"handles_on_page"`
}

type imageInputWire struct {
	Src    string `json:"src"`
	Type   string `json:"type"`
	RawB64 string `json:"raw_b64"` // standard base64-encoded image bytes
}

func (h *Handler) classifyJob(c *gin.Context) {
	var wire jobRequestWire
	if err := c.ShouldBindJSON(&wire); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
		return
	}

	if len(wire.Images) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "images array must not be empty"})
		return
	}

	// Decode each image's base64 bytes.
	req := &JobRequest{
		JobID:         wire.JobID,
		SourceURL:     wire.SourceURL,
		ProvenanceTag: wire.ProvenanceTag,
		HandlesOnPage: wire.HandlesOnPage,
		Images:        make([]ImageInput, len(wire.Images)),
	}

	for i, wi := range wire.Images {
		decoded, err := decodeBase64Image(wi.RawB64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("images[%s]: base64 decode failed: %v", strconv.Itoa(i), err),
			})
			return
		}
		req.Images[i] = ImageInput{
			Src:      wi.Src,
			Type:     wi.Type,
			RawBytes: decoded,
		}
	}

	jobResult, err := h.pipeline.ClassifyJob(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobResult)
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func decodeBase64Image(b64str string) ([]byte, error) {
	if b64str == "" {
		return nil, fmt.Errorf("empty raw_b64")
	}
	decoded, err := base64.StdEncoding.DecodeString(b64str)
	if err != nil {
		// Try URL-safe encoding as a fallback.
		decoded, err = base64.URLEncoding.DecodeString(b64str)
		if err != nil {
			return nil, err
		}
	}
	return decoded, nil
}
