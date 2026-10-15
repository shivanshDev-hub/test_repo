// Package render runs the air-gapped Chromium instance that turns sanitized
// HTML into real pixel screenshots.
//
// Isolation model (defense in depth):
//  1. The HTML reaching Render() is passed through sanitize.HTML() — an
//     allowlist pass that strips scripts, event handlers, foreign namespaces,
//     and fetch-capable CSS. The sandbox re-runs it regardless (trusts no
//     caller).
//  2. Chromium is launched with scripting disabled, networking disabled
//     (host-resolver-rules pin every DNS query to 0.0.0.0), no extensions, no
//     plugins, no GPU, and no shared-memory device.
//  3. The document is injected with page.SetContent — never navigated to. A
//     strict CSP meta tag forbids subresources and network entirely.
//  4. The surrounding container has no external network at all (the Docker
//     network render-sandbox lives on is internal: true).
//
// Even a fully hostile sanitizer escape therefore has no egress path: CSS
// {background:url(...)} resolves to 0.0.0.0 and fails in milliseconds.
package render

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/onioneye/render-sandbox/internal/sanitize"
)

// Config configures the sandbox.
type Config struct {
	// ChromePath overrides automatic binary discovery. Leave empty to let
	// chromedp find a Chrome/Chromium executable.
	ChromePath string
	// ViewportWidth/Height set the screenshot viewport.
	ViewportWidth  int
	ViewportHeight int
	// RenderTimeout caps one render.
	RenderTimeout time.Duration
	// MaxConcurrent bounds simultaneous Chromium processes; each instance
	// holds several hundred MB of anonymous memory.
	MaxConcurrent int
}

// DefaultConfig returns a hardened default configuration.
func DefaultConfig() *Config {
	return &Config{
		ChromePath:     "",
		ViewportWidth:  1280,
		ViewportHeight: 1024,
		RenderTimeout:  12 * time.Second,
		MaxConcurrent:  2,
	}
}
// Sandbox is a render worker.
type Sandbox struct {
	cfg    *Config
	sem    chan struct{}
	chrome string
}

// New builds a Sandbox, normalizing configuration. chromePath may be empty to
// let chromedp auto-discover the browser binary.
func New(cfg *Config, chromePath string) *Sandbox {
	c := cfg
	if c == nil {
		c = DefaultConfig()
	}
	if c.MaxConcurrent <= 0 {
		c.MaxConcurrent = 1
	}
	if c.ViewportWidth <= 0 {
		c.ViewportWidth = 1280
	}
	if c.ViewportHeight <= 0 {
		c.ViewportHeight = 1024
	}
	if c.RenderTimeout <= 0 {
		c.RenderTimeout = 12 * time.Second
	}
	return &Sandbox{
		cfg:    c,
		sem:    make(chan struct{}, c.MaxConcurrent),
		chrome: chromePath,
	}
}

// Result wraps a successful render.
type Result struct {
	PNG       []byte
	Sanitized string
}

// ssFlags lists the hardening flags, kept as a function so tests can assert
// on the flag set.
func ssFlags() []chromedp.ExecAllocatorOption {
	return []chromedp.ExecAllocatorOption{
		// Rendering mode.
		chromedp.Headless,
		chromedp.Flag("no-pings", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-features", "Translate,OptimizationHints"),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("hide-scrollbars", true),
		// No scripting, never.
		chromedp.Flag("disable-javascript", true),
		chromedp.Flag("disable-webgl", true),
		chromedp.Flag("disable-webrtc", true),
		// No extensions/plugins/helpers.
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("disable-breakpad", true),
		// Chromium 152's crashpad handler requires a crash database directory
		// when the environment provides a HOME. We disable both breakpad and
		// crash reporter entirely to avoid the "--database is required" crash.
		chromedp.Flag("disable-crash-reporter", true),
		chromedp.Flag("disable-crashpad", true),
		chromedp.Flag("disable-domain-reliability", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("metrics-recording-only", true),
		// Egress lockdown: every DNS resolution maps to 0.0.0.0 and no proxy
		// is ever consulted. Combined with the container's internal-only
		// network there is no egress path even for hostile CSS url() probes.
		chromedp.Flag("host-resolver-rules", "MAP * 0.0.0.0, EXCLUDE localhost"),
		chromedp.Flag("proxy-server", "direct://"),
		chromedp.Flag("proxy-bypass-list", "*"),
		// The container supplies seccomp + dropped capabilities; --no-sandbox
		// is required for the bundled Chromium's non-root user.
		chromedp.NoSandbox,
	}
}

// hardeningCSP forbids every subresource class and connection — an
// independent kill switch even if sanitization or the resolver rules failed.
const hardeningCSP = `<meta http-equiv="Content-Security-Policy" content="default-src 'none'; script-src 'none'; style-src 'unsafe-inline'; img-src 'none'; media-src 'none'; font-src 'none'; connect-src 'none'; frame-src 'none'; object-src 'none'; base-uri 'none'; form-action 'none'">`

// Render sanitizes html (again), renders it, and returns a PNG screenshot of
// the current viewport.
func (s *Sandbox) Render(ctx context.Context, html string) (*Result, error) {
	safe, err := sanitize.HTML(html)
	if err != nil {
		return nil, err
	}

	// Acquire a render slot.
	select {
	case s.sem <- struct{}{}:
		defer func() { <-s.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	safe = injectCSP(safe)

	renderCtx, cancelRender := context.WithTimeout(ctx, s.cfg.RenderTimeout)
	defer cancelRender()

	allocCtx, cancelAlloc := chromeAllocator(renderCtx, s.chrome)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	var png []byte
	tasks := chromedp.Tasks{
		// Disable script execution via CDP as well as the command-line switch.
		emulation.SetScriptExecutionDisabled(true),
		// Consistent device metrics for deterministic layout.
		emulation.SetDeviceMetricsOverride(int64(s.cfg.ViewportWidth), int64(s.cfg.ViewportHeight), 1, false),
		// Navigate to the empty, local about:blank document. This loads NO
		// remote or attacker content (the raw bytes of the target site never
		// touch Chromium — only the sanitized DOM does), and it guarantees a
		// main frame exists for SetDocumentContent.
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			tree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(tree.Frame.ID, safe).Do(ctx)
		}),
		// Give layout a short settling window. Fonts have no network to fetch,
		// so a few hundred ms is plenty on a warm page.
		chromedp.Sleep(400 * time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			params := page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithCaptureBeyondViewport(false).
				WithFromSurface(false)
			p, err := params.Do(ctx)
			if err != nil {
				return err
			}
			png = p
			return nil
		}),
	}
	if err := chromedp.Run(browserCtx, tasks...); err != nil {
		return nil, err
	}
	if len(png) == 0 {
		return nil, errors.New("empty screenshot buffer")
	}
	return &Result{PNG: png, Sanitized: safe}, nil
}

// chromeAllocator builds the hardened Chromium allocator.
func chromeAllocator(ctx context.Context, chrome string) (context.Context, context.CancelFunc) {
	if chrome != "" {
		return chromedp.NewExecAllocator(ctx,
			append(chromedp.DefaultExecAllocatorOptions[:],
				append([]chromedp.ExecAllocatorOption{chromedp.ExecPath(chrome)}, ssFlags()...)...)...)
	}
	return chromedp.NewExecAllocator(ctx,
		append(chromedp.DefaultExecAllocatorOptions[:], ssFlags()...)...)
}

// injectCSP slots the hardening CSP meta into the document head. The
// sanitizer guarantees an <html><head> skeleton exists.
func injectCSP(doc string) string {
	idx := strings.Index(doc, "<head>")
	if idx < 0 {
		return doc
	}
	return doc[:idx+len("<head>")] + hardeningCSP + doc[idx+len("<head>"):]
}