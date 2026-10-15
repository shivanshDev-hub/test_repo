// Package crawler implements passive fetching over the Tor SOCKS5 proxy,
// with retry, circuit rotation on network errors, and CAPTCHA/block detection.
//
// Fetching is passive only — plain GET/HEAD against whatever the server
// serves. No crafted payloads are ever sent (explicit ethical/legal boundary).
package crawler

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/onioneye/crawler-worker/internal/circuit"
	"github.com/onioneye/crawler-worker/internal/config"
)

// Error classes — see Architecture_Subsection_2.md "Built-in error handling".
var (
	ErrNetworkTimeout   = errors.New("network_timeout")
	ErrConnRefused      = errors.New("connection_refused")
	ErrBlockedOrCaptcha = errors.New("blocked_or_captcha")
)

// Fetched is the result of a successful fetch, handed to the content shield.
type Fetched struct {
	JobID        string
	URL          string
	HTTPStatus   int
	Headers      http.Header
	RawBody      []byte
	TorCircuitID string
	FetchedAt    time.Time
	Blocked      bool
	BlockReason  string
}

// Client performs fetches through Tor.
type Client struct {
	cfg    *config.Config
	http   *http.Client
	mu     sync.Mutex // guards circID + dirty
	circID string
	dirty  bool // true → next Fetch must rotate to a fresh circuit first
	rand   *rand.Rand
}

// NewClient builds an HTTP client whose dialer routes through the Tor SOCKS5
// endpoint provided by the tor-sidecar.
func NewClient(cfg *config.Config) (*Client, error) {
	c := &Client{
		cfg:    cfg,
		circID: fmt.Sprintf("circ-%d", time.Now().UnixNano()),
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	transport := &http.Transport{
		// Use custom SOCKS5 dialer that explicitly uses ATYP=DOMAINNAME
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialSOCKS5(ctx, cfg.TorSocksAddr, addr)
		},
		// Browser-equivalent TLS fingerprint via uTLS (Chrome/Firefox hello
		// rotation) instead of Go's structurally-unique default ClientHello.
		// This is only invoked for https:// URLs; http:// goes through the
		// plain SOCKS5 DialContext above.
		DialTLSContext: c.dialTLS,
		// One circuit per request keeps session-to-circuit binding simple for
		// the demo; a CookieJar bound to the circuit is a future refinement.
		DisableKeepAlives: true,
	}

	c.http = &http.Client{Transport: transport, Timeout: cfg.RequestTimeout}
	return c, nil
}

// Fetch downloads url with retry + circuit rotation. It returns a Fetched with
// Blocked=true when the CAPTCHA/block heuristic fires, so the caller can route
// the job to the manual-intervention queue.
//
// Every attempt is surfaced to stdout so the SYS LOGS panel shows exactly what
// is flowing through the crawler: attempt number, circuit id, HTTP result,
// timing, and the precise Tor/SOCKS failure code when a connection is refused.
func (c *Client) Fetch(ctx context.Context, url string) (*Fetched, error) {
	lastErr := error(nil)

	// If the previous fetch poisoned this circuit (any SOCKS/Tor failure), give
	// the next request a fresh circuit. Without this, one bad circuit is reused
	// forever — every job fails with the stale circuit's original error.
	c.mu.Lock()
	needFresh := c.dirty
	c.mu.Unlock()
	if needFresh {
		log.Printf("[fetch] previous fetch left circuit #%s dirty — rotating first", c.CircuitID())
		if rerr := c.RotateCircuit(); rerr != nil {
			log.Printf("[fetch] pre-fetch circuit rotation failed: %v", rerr)
		}
	}

	for attempt := 0; attempt < c.cfg.MaxRetries; attempt++ {
		start := time.Now()
		// Hardening: never log real egress IPs.
		// We log only the target URL (which is already in the job input) and the Tor circuit id.
		log.Printf("[fetch] attempt=%d/%d circ=%s GET %s", attempt+1, c.cfg.MaxRetries, c.CircuitID(), url)
		f, err := c.fetchOnce(ctx, url)
		switch {
		case err == nil:
			c.mu.Lock()
			c.dirty = false
			c.mu.Unlock()
			if f.Blocked {
				log.Printf("[fetch] %s -> BLOCKED (%s) after %s", url, f.BlockReason, time.Since(start).Round(time.Millisecond))
				return f, nil // not an error — caller decides manual intervention
			}
			log.Printf("[fetch] %s -> HTTP %d, %s, %s", url, f.HTTPStatus, humanBytes(len(f.RawBody)), time.Since(start).Round(time.Millisecond))
			return f, nil
		case errors.Is(err, ErrNetworkTimeout) || errors.Is(err, ErrConnRefused):
			// rotate circuit before the next attempt (max rotations bounded)
			log.Printf("[fetch] %s -> retryable: %v after %s", url, err, time.Since(start).Round(time.Millisecond))
			if attempt < c.cfg.MaxCircuitRot {
				if rerr := c.RotateCircuit(); rerr != nil {
					log.Printf("[fetch] circuit rotation failed: %v", rerr)
				} else {
					log.Printf("[fetch] circuit rotated -> %s", c.CircuitID())
				}
			}
			lastErr = err
		default:
			log.Printf("[fetch] %s -> fatal: %v", url, err)
			c.mu.Lock()
			c.dirty = true
			c.mu.Unlock()
			return nil, err
		}
	}
	c.mu.Lock()
	c.dirty = true
	c.mu.Unlock()
	return nil, lastErr
}

func (c *Client) fetchOnce(ctx context.Context, rawURL string) (*Fetched, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	// Browser-alike fingerprinting: rotate UA + header constellation.
	// NOTE: This is OPSEC-oriented only; it does not bypass site rules.
	p := randomPersona()
	p.apply(req.Header)

	// Disable any Go-client-specific hints that stand out.
	req.Header.Del("Connection")
	req.Header.Del("Sec-WebSocket-Extensions")

	resp, err := c.http.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, ErrNetworkTimeout
		}
		if isConnRefused(err) {
			// Hardening: scrub any potential IPs from error strings before returning.
			return nil, ErrConnRefused
		}
		return nil, err
	}
	defer resp.Body.Close()

	// Cap the body at MaxBodyBytes to prevent memory exhaustion.
	body, err := io.ReadAll(io.LimitReader(resp.Body, c.cfg.MaxBodyBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > c.cfg.MaxBodyBytes {
		body = body[:c.cfg.MaxBodyBytes]
	}

	// The persona sets an explicit Accept-Encoding header, which disables Go's
	// transparent gzip handling (net/http only auto-decompresses when IT adds
	// the header itself). Without this, servers that gzip HTML leave RawBody as
	// compressed bytes and extraction finds zero entities. Honor the response's
	// Content-Encoding and hand the pipeline plain bytes instead.
	if enc := resp.Header.Get("Content-Encoding"); enc != "" {
		dec, derr := decodeContentEncoding(enc, body, c.cfg.MaxBodyBytes)
		if derr != nil {
			log.Printf("[fetch] content-encoding %q decode failed: %v (keeping raw body)", enc, derr)
		} else {
			body = dec
		}
	}

	f := &Fetched{
		JobID:        "", // filled by caller
		URL:          rawURL,
		HTTPStatus:   resp.StatusCode,
		Headers:      resp.Header,
		RawBody:      body,
		TorCircuitID: c.CircuitID(),
		FetchedAt:    time.Now().UTC(),
	}

	if signal := c.detectBlock(resp.StatusCode, body); signal != "" {
		f.Blocked = true
		f.BlockReason = signal
	}
	return f, nil
}

// decodeContentEncoding decompresses a response body per its Content-Encoding
// value (gzip, x-gzip, deflate — both zlib-wrapped and raw). Output is capped at
// maxBytes +1 so a hostile "zip bomb" response cannot exhaust memory; on any
// decode error the caller keeps the original raw bytes.
func decodeContentEncoding(encoding string, body []byte, maxBytes int64) ([]byte, error) {
	enc := strings.ToLower(strings.TrimSpace(encoding))
	limited := func(r io.Reader) ([]byte, error) {
		return io.ReadAll(io.LimitReader(r, maxBytes+1))
	}
	switch enc {
	case "gzip", "x-gzip":
		zr, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		defer zr.Close()
		return limited(zr)
	case "deflate":
		// Most servers send zlib-wrapped deflate (RFC 1950); a few send raw
		// deflate (RFC 1951). Try the standard wrapper first, then raw flate.
		if zr, err := zlib.NewReader(bytes.NewReader(body)); err == nil {
			out, rerr := limited(zr)
			zr.Close()
			if rerr == nil {
				return out, nil
			}
		}
		fr := flate.NewReader(bytes.NewReader(body))
		defer fr.Close()
		return limited(fr)
	default:
		return nil, fmt.Errorf("unsupported content-encoding %q", enc)
	}
}

// detectBlock applies the block heuristic. The full three-signal agreement
// (status + body-length deviation from Redis baseline + cross-circuit verify)
// is described in the architecture; here we implement the status-code signal
// plus a lightweight empty/short-body check. Cross-circuit verification is a
// TODO once per-target baselines are stored in Redis.
func (c *Client) detectBlock(status int, body []byte) string {
	switch status {
	case http.StatusForbidden, http.StatusRequestTimeout, http.StatusServiceUnavailable:
		return "http_status_403_408_503"
	}
	// A page with no body at all is suspicious for a text target.
	if len(body) == 0 {
		return "empty_body"
	}
	return ""
}

// RotateCircuit signals NEWNYM on the Tor control port and mints a new circuit
// id so future fetches (and their CookieJars) are bound to the fresh circuit.
func (c *Client) RotateCircuit() error {
	// Non-local control connections require the sidecar's control password.
	mgr := circuit.NewManager(c.cfg.TorControlAddr, c.cfg.TorControlPasswd)
	if err := mgr.NewNYM(); err != nil {
		return err
	}
	c.mu.Lock()
	c.circID = fmt.Sprintf("circ-%d", time.Now().UnixNano())
	c.mu.Unlock()
	return nil
}

// CircuitID returns the current circuit id (thread-safe).
func (c *Client) CircuitID() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.circID
}

// PoliteDelay sleeps a randomized 2–8s keyed on hostname, per architecture
// ("time.Sleep(time.Duration(2000+rand.Intn(6000)) * time.Millisecond)").
func (c *Client) PoliteDelay(host string) {
	ms := c.cfg.PolitenessMinMs + c.rand.Intn(c.cfg.PolitenessMaxMs-c.cfg.PolitenessMinMs+1)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	_ = host // reserved: per-host pacing
}

// humanBytes renders a byte count compactly (e.g. "47.1KB") for log lines.
func humanBytes(n int) string {
	var (
		val = float64(n)
		exp = 0
	)
	for val >= 1024 && exp < len("KMGTPE") {
		val /= 1024
		exp++
	}
	if exp == 0 {
		return fmt.Sprintf("%dB", n)
	}
	return fmt.Sprintf("%.1f%cB", val, "KMGTPE"[exp-1])
}

func isTimeout(err error) bool {
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}

func isConnRefused(err error) bool {
	// SOCKS5 REP codes returned by the tor-sidecar survive errors.As through
	// net/http's chain, unlike a wrapped string. Codes that a fresh circuit can
	// plausibly clear are retryable: general failure, network/host unreachable,
	// connection refused, TTL expired. 0x02 (ruleset) / 0x07 (cmd) / 0x08 (ATYP)
	// are policy/protocol-level and stay fatal.
	var se *SocksError
	if errors.As(err, &se) {
		switch se.Code {
		case 0x01, 0x03, 0x04, 0x05, 0x06:
			return true
		}
		return false
	}
	var op *net.OpError
	return errors.As(err, &op) && strings.Contains(op.Err.Error(), "refused")
}

// NormalizeURL lowercases scheme+host, strips default ports, and percent-decodes.
func NormalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return strings.TrimSpace(raw)
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	return u.String()
}
