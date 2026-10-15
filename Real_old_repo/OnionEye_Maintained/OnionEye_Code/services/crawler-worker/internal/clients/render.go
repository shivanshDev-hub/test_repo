package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// maxRenderPNGBytes caps the screenshot bytes the sandbox may hand back.
// A 1280x1024 PNG rarely exceeds 1 MiB; 20 MiB is generous headroom against
// a pathological fully-white or noisy page.
const maxRenderPNGBytes = 20 << 20

// RenderClient talks to the render-sandbox service. The sandbox only accepts
// sanitized-or-unsanitized HTML — it re-sanitizes every payload itself before
// injecting it into the network-isolated Chromium instance (zero trust).
type RenderClient struct {
	baseURL string
	token   string // shared-secret for worker → sandbox auth
	http    *http.Client
}

// NewRenderClient builds a client for the render-sandbox service.
// token is the shared secret sent as X-OnionEye-Token; when non-empty the
// sandbox rejects requests without a matching header.
func NewRenderClient(baseURL, token string) *RenderClient {
	return &RenderClient{
		baseURL: baseURL,
		token:   token,
		// Renders routinely take several seconds (launch Chromium, layout,
		// screenshot). 30s mirrors the sandbox's own request ceiling.
		http: &http.Client{Timeout: 35 * time.Second},
	}
}

type renderRequest struct {
	JobID     string          `json:"job_id,omitempty"`
	SourceURL string          `json:"source_url"`
	HTML      string          `json:"html"`
	Viewport  viewportRequest `json:"viewport,omitempty"`
}

type viewportRequest struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

// Render posts page HTML to the render-sandbox and returns the PNG screenshot.
//
// The sandbox rejects HTML larger than its own MaxInputBytes (2 MiB) with 413;
// we guard here too so the failure is a clean error instead of a huge POST.
func (r *RenderClient) Render(ctx context.Context, sourceURL, html string) ([]byte, error) {
	if len(html) > 2<<20 {
		return nil, fmt.Errorf("html too large for render sandbox: %d bytes", len(html))
	}

	req := renderRequest{
		JobID:     "screenshot-on-demand",
		SourceURL: sourceURL,
		HTML:      html,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/render", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("render request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if r.token != "" {
		httpReq.Header.Set("X-OnionEye-Token", r.token)
	}

	resp, err := r.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("render request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("render failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(detail)))
	}

	png, err := io.ReadAll(io.LimitReader(resp.Body, maxRenderPNGBytes))
	if err != nil {
		return nil, fmt.Errorf("read render response: %w", err)
	}
	if len(png) == 0 {
		return nil, fmt.Errorf("render returned empty screenshot")
	}
	return png, nil
}
