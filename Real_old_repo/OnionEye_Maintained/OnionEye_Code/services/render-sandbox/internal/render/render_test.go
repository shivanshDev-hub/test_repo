package render

import (
"context"
"testing"
"time"
"os"
"strings"
)

func TestRenderSmoke(t *testing.T) {
	chrome := os.Getenv("CHROME_BINARY")
	if chrome == "" {
		chrome = "/usr/bin/chromium"
	}
	if _, err := os.Stat(chrome); err != nil {
		t.Skipf("chromium not found at %s; skipping smoke test", chrome)
	}
	cfg := &Config{
		ChromePath:     chrome,
		ViewportWidth:  800,
		ViewportHeight: 600,
		RenderTimeout:  15 * time.Second,
		MaxConcurrent:  1,
	}
	sb := New(cfg, chrome)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	html := `<html><head></head><body style="background:#111;color:#fff;font-family:sans-serif">
		<h1>Hello OnionEye</h1>
		<script>alert('XSS')</script>
		<img onerror="alert(1)" src="http://evil.onion/">
		<div style="background:url(http://evil.onion/)">safe div</div>
	</body></html>`
	res, err := sb.Render(ctx, html)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if len(res.PNG) < 1024 {
		t.Fatalf("PNG too small: %d bytes", len(res.PNG))
	}
	if strings.Contains(res.Sanitized, "<script") {
		t.Fatal("script tag survived into sanitized DOM")
	}
	if strings.Contains(res.Sanitized, "onerror") {
		t.Fatal("onerror survived into sanitized DOM")
	}
	if strings.Contains(res.Sanitized, "background:url") || strings.Contains(res.Sanitized, "background-image:url") {
		t.Fatal("url() survived into sanitized DOM")
	}
	if strings.Contains(res.Sanitized, "alert") {
		t.Fatal("alert() survived into sanitized DOM")
	}
	t.Logf("render smoke: PNG %d bytes, sanitized %d bytes", len(res.PNG), len(res.Sanitized))
}
